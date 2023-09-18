/*
Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"context"
	"reflect"

	"github.com/go-logr/logr"
	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/pkg/constants"
	switchespkg "github.com/onmetal/metal-api/pkg/switches"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// SwitchReconciler reconciles Switch object corresponding
// to given Inventory object.
type SwitchReconciler struct {
	client.Client

	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/finalizers,verbs=update
// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switchconfigs,verbs=get;list;watch
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *SwitchReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	obj := &switchv1beta1.Switch{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	ref, err := reference.GetReference(r.Scheme, obj)
	if err != nil {
		r.Log.WithValues("request", req.NamespacedName).Error(err, "failed to construct reference")
		return ctrl.Result{}, err
	}

	log := r.Log.WithValues("object", *ref)
	log.Info("reconciliation started")
	requestCtx := logr.NewContext(ctx, log)
	return r.reconcileRequired(requestCtx, obj)
}

func (r *SwitchReconciler) reconcileRequired(ctx context.Context, obj *switchv1beta1.Switch) (ctrl.Result, error) {
	if !obj.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, nil
	}
	return r.reconcileManaged(ctx, obj)
}

func (r *SwitchReconciler) reconcileManaged(ctx context.Context, obj *switchv1beta1.Switch) (ctrl.Result, error) {
	if !obj.Managed() {
		log := logr.FromContextOrDiscard(ctx)
		log.WithValues("reason", constants.ReasonUnmanagedSwitch).Info("reconciliation finished")
		return ctrl.Result{}, nil
	}
	return r.reconcile(ctx, obj)
}

func (r *SwitchReconciler) reconcile(ctx context.Context, obj *switchv1beta1.Switch) (ctrl.Result, error) {
	if ok, err := r.mapToInventory(ctx, obj); !ok {
		return ctrl.Result{}, err
	}
	if ok, err := r.configSelectorValid(ctx, obj); !ok {
		return ctrl.Result{}, err
	}
	return r.reconcileState(ctx, obj)
}

func (r *SwitchReconciler) reconcileState(ctx context.Context, obj *switchv1beta1.Switch) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	svc := switchespkg.NewSwitchEnvironmentSvc(r.Client, log)
	env := svc.GetEnvironment(ctx, obj)

	// stop reconciliation if state is "Pending" and SwitchConfig
	// object matching defined selector was not found.
	if obj.StatePending() && env.Config == nil {
		return ctrl.Result{}, nil
	}

	snapshot := obj.DeepCopy()

	writer := switchespkg.NewSwitchStateWriter(r.Recorder /*, log*/)
	writer.SetEnvironment(env).
		RegisterStateFunc(switchespkg.Initialize).
		RegisterStateFunc(switchespkg.UpdateInterfaces).
		RegisterStateFunc(switchespkg.UpdateNeighbors).
		RegisterStateFunc(switchespkg.UpdateLayerAndRole).
		RegisterStateFunc(switchespkg.UpdateConfigRef).
		RegisterStateFunc(switchespkg.UpdatePortParameters).
		RegisterStateFunc(switchespkg.UpdateLoopbacks).
		RegisterStateFunc(switchespkg.UpdateASN).
		RegisterStateFunc(switchespkg.UpdateSubnets).
		RegisterStateFunc(switchespkg.UpdateSwitchPortIPs).
		RegisterStateFunc(switchespkg.SetStateReady)
	_, result, msg, err := writer.WriteState(obj).Unwrap()
	if err == nil {
		log.WithValues("reason", result).Info(msg)
	} else {
		log.WithValues("reason", result).Error(err, msg)
	}

	if !switchespkg.ObjectChanged(snapshot.DeepCopy(), obj.DeepCopy()) {
		log.WithValues("reason", constants.ReasonObjectUnchanged).Info("reconciliation finished")
		return ctrl.Result{}, nil
	}
	return r.updateState(ctx, obj)
}

func (r *SwitchReconciler) updateState(ctx context.Context, obj *switchv1beta1.Switch) (ctrl.Result, error) {
	obj.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
	err := r.Status().Patch(ctx, obj, client.Apply, switchespkg.PatchOpts)
	if apierrors.IsConflict(err) {
		return ctrl.Result{Requeue: true}, nil
	}
	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// predicate to filter switch object update which was not caused
	// by conditions lastUpdateTimestamp change.
	discoverObjectChangesPredicate := predicate.Funcs{
		UpdateFunc: detectChangesPredicate,
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1beta1.Switch{}).
		WithOptions(controller.Options{
			RecoverPanic: pointer.Bool(true),
		}).
		WithEventFilter(predicate.And(discoverObjectChangesPredicate)).
		// watches for Switch objects required, because switches are
		// interconnected and changes in configuration of one object
		// might affect another objects.
		Watches(&switchv1beta1.Switch{}, &handler.Funcs{
			UpdateFunc: r.handleSwitchUpdateEvent,
			DeleteFunc: r.handleSwitchDeleteEvent,
		}).
		// watches for SwitchConfig objects required, because
		// Switch objects' configuration is based on config defined
		// in SwitchConfig objects, so changes must be tracked.
		Watches(&switchv1beta1.SwitchConfig{}, &handler.Funcs{
			UpdateFunc: r.handleSwitchConfigUpdateEvent,
		}).
		// watches for Inventory objects required, because
		// changes in hardware, especially discovering new
		// neighbors connected to switch ports must be tracked.
		Watches(&inventoryv1alpha1.Inventory{}, &handler.Funcs{
			CreateFunc: r.handleInventoryCreateEvent,
			UpdateFunc: r.handleInventoryUpdateEvent,
		}).
		Complete(r)
}

func detectChangesPredicate(e event.UpdateEvent) bool {
	var (
		switchOld, switchNew       *switchv1beta1.Switch
		configOld, configNew       *switchv1beta1.SwitchConfig
		inventoryOld, inventoryNew *inventoryv1alpha1.Inventory
		castOldOk, castNewOk       bool
	)
	switchOld, castOldOk = e.ObjectOld.(*switchv1beta1.Switch)
	switchNew, castNewOk = e.ObjectNew.(*switchv1beta1.Switch)
	if castOldOk && castNewOk {
		return switchespkg.ObjectChanged(switchOld.DeepCopy(), switchNew.DeepCopy())
	}
	configOld, castOldOk = e.ObjectOld.(*switchv1beta1.SwitchConfig)
	configNew, castNewOk = e.ObjectNew.(*switchv1beta1.SwitchConfig)
	if castOldOk && castNewOk {
		specChanged := !reflect.DeepEqual(configOld.Spec, configNew.Spec)
		labelsChanged := !reflect.DeepEqual(configOld.Labels, configNew.Labels)
		return specChanged || labelsChanged
	}
	inventoryOld, castOldOk = e.ObjectOld.(*inventoryv1alpha1.Inventory)
	inventoryNew, castNewOk = e.ObjectNew.(*inventoryv1alpha1.Inventory)
	if castOldOk && castNewOk {
		return !reflect.DeepEqual(inventoryOld.Spec, inventoryNew.Spec)
	}
	return false
}

func (r *SwitchReconciler) handleSwitchUpdateEvent(_ context.Context, e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	r.Log.WithValues("handler", "SwitchUpdateEvent").Info("enqueueing switches")
	objOld, okOld := e.ObjectOld.(*switchv1beta1.Switch)
	objNew, okNew := e.ObjectNew.(*switchv1beta1.Switch)
	if !okOld || !okNew {
		return
	}
	// if switch object has no changes, which affect neighbors, then there is no need to
	// enqueue it's neighbors for reconciliation.
	if !switchespkg.ObjectChanged(objOld.DeepCopy(), objNew.DeepCopy()) {
		return
	}
	switchesQueue := make(map[string]struct{})
	for _, nicData := range objOld.Status.Interfaces {
		if !switchespkg.NeighborIsSwitch(nicData) {
			continue
		}
		switchesQueue[nicData.Peer.GetObjectReferenceName()] = struct{}{}
	}
	for _, nicData := range objNew.Status.Interfaces {
		if !switchespkg.NeighborIsSwitch(nicData) {
			continue
		}
		switchesQueue[nicData.Peer.GetObjectReferenceName()] = struct{}{}
	}
	for name := range switchesQueue {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Namespace: objNew.Namespace,
			Name:      name,
		}})
	}
}

func (r *SwitchReconciler) handleSwitchDeleteEvent(_ context.Context, e event.DeleteEvent, q workqueue.RateLimitingInterface) {
	r.Log.WithValues("handler", "SwitchDeleteEvent").Info("enqueueing switches")
	obj, ok := e.Object.(*switchv1beta1.Switch)
	if !ok {
		return
	}
	switchesQueue := make(map[string]struct{})
	for _, nicData := range obj.Status.Interfaces {
		if !switchespkg.NeighborIsSwitch(nicData) {
			continue
		}
		switchesQueue[nicData.Peer.GetObjectReferenceName()] = struct{}{}
	}
	for name := range switchesQueue {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Namespace: obj.Namespace,
			Name:      name,
		}})
	}
}

func (r *SwitchReconciler) handleSwitchConfigUpdateEvent(ctx context.Context, e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	r.Log.WithValues("handler", "SwitchConfigUpdateEvent").Info("enqueueing switches")
	_, castOldOk := e.ObjectOld.(*switchv1beta1.SwitchConfig)
	_, castNewOk := e.ObjectNew.(*switchv1beta1.SwitchConfig)
	if !castOldOk || !castNewOk {
		return
	}
	switches := &switchv1beta1.SwitchList{}
	if err := r.List(ctx, switches); err != nil {
		r.Log.Error(err, "failed to list switches while handling SwitchConfig update event")
		return
	}
	for _, item := range switches.Items {
		q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(&item)})
	}
}

func (r *SwitchReconciler) handleInventoryCreateEvent(_ context.Context, e event.CreateEvent, q workqueue.RateLimitingInterface) {
	r.Log.WithValues("handler", "InventoryCreateEvent").Info("enqueueing corresponding switch")
	inventory, castOk := e.Object.(*inventoryv1alpha1.Inventory)
	if !castOk {
		return
	}
	q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(inventory)})
}

func (r *SwitchReconciler) handleInventoryUpdateEvent(_ context.Context, e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	r.Log.WithValues("handler", "InventoryUpdateEvent").Info("enqueueing corresponding switch")
	_, castOldOk := e.ObjectOld.(*inventoryv1alpha1.Inventory)
	inventoryNew, castNewOk := e.ObjectNew.(*inventoryv1alpha1.Inventory)
	if castOldOk && castNewOk {
		return
	}
	q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(inventoryNew)})
}

func (r *SwitchReconciler) mapToInventory(ctx context.Context, obj *switchv1beta1.Switch) (bool, error) {
	inventoryRefDefined := obj.GetInventoryRef() != constants.EmptyString
	_, inventoriedLabel := obj.Labels[constants.InventoriedLabel]
	_, chassisIDLabel := obj.Labels[constants.LabelChassisID]
	if !(inventoryRefDefined && inventoriedLabel && chassisIDLabel) {
		inv := &inventoryv1alpha1.Inventory{}
		if err := r.Get(ctx, client.ObjectKeyFromObject(obj), inv); err != nil {
			return false, err
		}
		obj.SetInventoryRef(inv.Name)
		obj.UpdateSwitchAnnotations(inv)
		obj.UpdateSwitchLabels(inv)
		obj.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
		err := r.Patch(ctx, obj, client.Apply, switchespkg.PatchOpts)
		return false, err
	}
	return true, nil
}

func (r *SwitchReconciler) configSelectorValid(ctx context.Context, obj *switchv1beta1.Switch) (bool, error) {
	if switchespkg.SwitchConfigSelectorInvalid(obj) {
		switchespkg.UpdateSwitchConfigSelector(obj)
		obj.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
		err := r.Patch(ctx, obj, client.Apply, switchespkg.PatchOpts)
		return false, err
	}
	return true, nil
}
