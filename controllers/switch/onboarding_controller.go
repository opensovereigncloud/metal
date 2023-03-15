/*
Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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
	"time"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/pkg/constants"
	"github.com/onmetal/metal-api/pkg/errors"
	switchespkg "github.com/onmetal/metal-api/pkg/switches"
)

// OnboardingReconciler reconciles Switch object corresponding
// to given Inventory object.
type OnboardingReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch

func (r *OnboardingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("reconciliation started", "name", req.NamespacedName)
	nestedCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	inventory := &inventoryv1alpha1.Inventory{}
	if err := r.Get(nestedCtx, req.NamespacedName, inventory); err != nil {
		switch {
		case apierrors.IsNotFound(err):
			r.Log.Info("requested Inventory object not found", "name", req.NamespacedName)
		default:
			r.Log.Info(
				"failed to get requested Inventory object",
				"name", req.NamespacedName,
				"error", err.Error(),
			)
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	result, err := r.reconcile(nestedCtx, inventory)
	r.Log.Info("reconciliation finished", "name", req.NamespacedName)
	return result, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *OnboardingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&inventoryv1alpha1.Inventory{}).
		WithOptions(controller.Options{
			RateLimiter:  workqueue.NewItemExponentialFailureRateLimiter(time.Millisecond*500, time.Minute*5),
			RecoverPanic: pointer.Bool(true),
		}).
		WithEventFilter(setupPredicates()).
		Watches(&source.Kind{Type: &switchv1beta1.Switch{}}, handler.Funcs{
			CreateFunc: switchCreateEventHandler,
			UpdateFunc: switchUpdateEventHandler,
		}).
		Complete(r)
}

func setupPredicates() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: labelsCreateEventPredicate,
		UpdateFunc: labelsUpdateEventPredicate,
	}
}

func labelsCreateEventPredicate(e event.CreateEvent) bool {
	result := false
	inventoryObject, sourceIsInventory := e.Object.(*inventoryv1alpha1.Inventory)
	switchObject, sourceIsSwitch := e.Object.(*switchv1beta1.Switch)
	if sourceIsInventory {
		result = checkObjectMetadata(inventoryObject)
	}
	if sourceIsSwitch {
		result = checkObjectMetadata(switchObject)
	}
	return result
}

func labelsUpdateEventPredicate(e event.UpdateEvent) bool {
	inventoryObject, sourceIsInventory := e.ObjectNew.(*inventoryv1alpha1.Inventory)
	switchObject, sourceIsSwitch := e.ObjectNew.(*switchv1beta1.Switch)
	if sourceIsInventory {
		inventoryOld, ok := e.ObjectOld.(*inventoryv1alpha1.Inventory)
		if !ok {
			return false
		}
		specChanged := !reflect.DeepEqual(inventoryOld.Spec, inventoryObject.Spec)
		metaMatchRequirements := checkObjectMetadata(inventoryObject)
		return specChanged && metaMatchRequirements
	}
	if sourceIsSwitch {
		return checkObjectMetadata(switchObject)
	}
	return false
}

func checkObjectMetadata(obj client.Object) bool {
	switch obj.(type) {
	case *inventoryv1alpha1.Inventory:
		labels := obj.GetLabels()
		_, ok := labels[constants.SizeLabel]
		return ok
	case *switchv1beta1.Switch:
		labels := obj.GetLabels()
		annotations := obj.GetAnnotations()
		_, inventoriedLabelExist := labels[constants.InventoriedLabel]
		_, chassisIDAnnotationExist := annotations[constants.HardwareChassisIDAnnotation]
		return !inventoriedLabelExist || !chassisIDAnnotationExist
	default:
		return false
	}
}

func switchCreateEventHandler(e event.CreateEvent, q workqueue.RateLimitingInterface) {
	obj, ok := e.Object.(*switchv1beta1.Switch)
	if !ok {
		return
	}
	enqueue(obj, q)
}

func switchUpdateEventHandler(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	obj, ok := e.ObjectNew.(*switchv1beta1.Switch)
	if !ok {
		return
	}
	enqueue(obj, q)
}

func enqueue(obj *switchv1beta1.Switch, q workqueue.RateLimitingInterface) {
	namespace := obj.Namespace
	if name := obj.GetInventoryRef(); name != constants.EmptyString {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: name}})
		return
	}
	q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(obj)})
}

func (r *OnboardingReconciler) reconcile(ctx context.Context, inv *inventoryv1alpha1.Inventory) (ctrl.Result, error) {
	if !inv.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, nil
	}
	return r.onboarding(ctx, inv)
}

func (r *OnboardingReconciler) onboarding(ctx context.Context, inv *inventoryv1alpha1.Inventory) (ctrl.Result, error) {
	key := client.ObjectKeyFromObject(inv)
	r.Log.Info("onboarding started", "name", key)
	targetSwitch := &switchv1beta1.Switch{}
	if err := r.Get(ctx, key, targetSwitch); err != nil {
		if !apierrors.IsNotFound(err) {
			r.Log.Error(err, "failed to get Switch object", "name", key)
			return ctrl.Result{}, err
		}
		r.Log.Info("prepare switch object to create", "name", key)
		targetSwitch = &switchv1beta1.Switch{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Switch",
				APIVersion: "switch.onmetal.de/v1beta1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: switchv1beta1.SwitchSpec{
				Managed:   pointer.Bool(true),
				Cordon:    pointer.Bool(false),
				TopSpine:  pointer.Bool(false),
				ScanPorts: pointer.Bool(true),
			},
		}
	}
	r.Log.Info("set switch properties, labels, annotations", "name", key)
	targetSwitch.SetInventoryRef(key.Name)
	if targetSwitch.Labels == nil {
		targetSwitch.Labels = map[string]string{}
	}
	targetSwitch.UpdateSwitchAnnotations(inv)
	targetSwitch.UpdateSwitchLabels(inv)
	targetSwitch.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
	targetSwitch.SetCondition(constants.ConditionInterfacesOK, false).
		SetReason(errors.ErrorReasonDataOutdated.String())
	r.Log.Info("apply changes for switch", "name", key)
	result := ctrl.Result{}
	err := r.Patch(ctx, targetSwitch, client.Apply, switchespkg.PatchOpts)
	if apierrors.IsConflict(err) {
		r.Log.Info("failed to patch Switch, object is outdated")
		result.Requeue = true
		return result, nil
	}
	return result, err
}
