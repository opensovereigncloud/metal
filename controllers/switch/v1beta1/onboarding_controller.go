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
	"time"

	"github.com/go-logr/logr"
	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var RequeuingInterval = time.Second * 15

// OnboardingReconciler reconciles Switch object corresponding
// to given Inventory object.
type OnboardingReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

type wrappedError struct {
	Err       error
	Msg       string
	KeyValues []interface{}
}

//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *OnboardingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	log := r.Log.WithValues("switch-onboarding", req.NamespacedName)
	result = ctrl.Result{}

	obj := &inventoryv1alpha1.Inventory{}
	if err = r.Get(ctx, req.NamespacedName, obj); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info(
				"requested resource not found",
				"name", req.NamespacedName, "gvk", obj.GroupVersionKind().String(),
			)
		} else {
			log.Error(
				err,
				"failed to get requested resource",
				"name", req.NamespacedName, "gvk", obj.GroupVersionKind().String(),
			)
		}
		return result, client.IgnoreNotFound(err)
	}

	existingLabels := obj.GetLabels()
	if len(existingLabels) == 0 {
		return
	}

	sizeLabel := inventoryv1alpha1.GetSizeMatchLabel(switchv1beta1.CSwitchSizeName)
	if _, ok := existingLabels[sizeLabel]; !ok {
		return
	}

	switchObject, switchOnboarded, werr := r.switchOnboarded(ctx, obj)
	if werr != nil {
		err = werr.Err
		log.Error(werr.Err, werr.Msg, werr.KeyValues...)
		result.RequeueAfter = RequeuingInterval
		return
	}
	if switchOnboarded {
		goto updateMetadata
	}

	switchObject, werr = r.onboardSwitch(ctx, obj)
	if werr != nil {
		err = werr.Err
		log.Error(werr.Err, werr.Msg, werr.KeyValues...)
		result.RequeueAfter = RequeuingInterval
		return
	}

updateMetadata:
	switchObject.UpdateSwitchLabels(obj)
	switchObject.UpdateSwitchAnnotations(obj)
	if err = r.Update(ctx, switchObject); err != nil {
		log.Error(err, "failed to onboard switch", "name", switchObject.Name, "gvk", switchObject.GroupVersionKind().String())
		result.RequeueAfter = RequeuingInterval
		return
	}

	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *OnboardingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&inventoryv1alpha1.Inventory{}).
		// watches for switches to handle CREATE events of corresponding objects
		// on CREATE event: enqueue corresponding inventory
		Watches(&source.Kind{Type: &switchv1beta1.Switch{}}, &handler.Funcs{
			CreateFunc: r.handleSwitchCreate,
		}).
		Complete(r)
}

// nolint:forcetypeassert
func (r *OnboardingReconciler) handleSwitchCreate(e event.CreateEvent, q workqueue.RateLimitingInterface) {
	eventSource := e.Object.(*switchv1beta1.Switch)
	if _, ok := eventSource.Labels[switchv1beta1.InventoriedLabel]; ok {
		return
	}
	target := &inventoryv1alpha1.Inventory{}
	key := types.NamespacedName{
		Namespace: eventSource.Namespace,
		Name:      eventSource.Spec.UUID,
	}
	err := r.Get(context.Background(), key, target)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Log.Info(
				"requested resource not found",
				"name", key, "gvk", target.GroupVersionKind().String(),
			)
		} else {
			r.Log.Error(
				err,
				"failed to get requested resource",
				"name", key, "gvk", target.GroupVersionKind().String(),
			)
		}
		return
	}
	q.Add(reconcile.Request{NamespacedName: key})
}

func (r *OnboardingReconciler) switchOnboarded(ctx context.Context, inv *inventoryv1alpha1.Inventory) (*switchv1beta1.Switch, bool, *wrappedError) {
	switches := &switchv1beta1.SwitchList{}
	labelsReq, _ := switchv1beta1.GetLabelSelector(switchv1beta1.InventoryRefLabel, selection.Equals, []string{inv.Name})
	selector := labels.NewSelector().Add(*labelsReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         100,
	}
	if err := r.List(ctx, switches, opts); err != nil {
		return nil, false, &wrappedError{Err: err, Msg: "failed to list resources", KeyValues: []interface{}{"gvk", switches.GroupVersionKind().String()}}
	}
	if len(switches.Items) != 0 {
		return &switches.Items[0], true, nil
	}
	return nil, false, nil
}

func (r *OnboardingReconciler) onboardSwitch(ctx context.Context, inv *inventoryv1alpha1.Inventory) (*switchv1beta1.Switch, *wrappedError) {
	switchObject := &switchv1beta1.Switch{}
	switches := &switchv1beta1.SwitchList{}
	labelsReq, _ := switchv1beta1.GetLabelSelector(switchv1beta1.InventoriedLabel, selection.DoesNotExist, []string{})
	selector := labels.NewSelector().Add(*labelsReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         100,
	}
	if err := r.List(ctx, switches, opts); err != nil {
		return nil, &wrappedError{Err: err, Msg: "failed to list resources", KeyValues: []interface{}{"gvk", switches.GroupVersionKind().String()}}
	}

	switch len(switches.Items) > 0 {
	case true:
		switchObject = getCorrespondingSwitch(inv, switches)
		if switchObject != nil {
			break
		}
		fallthrough
	case false:
		switchObject = &switchv1beta1.Switch{
			ObjectMeta: metav1.ObjectMeta{
				Name:      inv.Name,
				Namespace: inv.Namespace,
				Labels: map[string]string{
					switchv1beta1.InventoriedLabel:  "true",
					switchv1beta1.InventoryRefLabel: inv.Name,
				},
			},
			Spec: switchv1beta1.SwitchSpec{
				UUID:     inv.Name,
				TopSpine: false,
				Managed:  true,
				Cordon:   false,
			},
		}
		err := r.Create(ctx, switchObject)
		if err != nil {
			return nil, &wrappedError{
				Err:       err,
				Msg:       "failed to get or create corresponding resource",
				KeyValues: []interface{}{"name", switchObject.Name, "gvk", switchObject.GroupVersionKind().String()}}
		}
	}
	return switchObject, nil
}

func getCorrespondingSwitch(inv *inventoryv1alpha1.Inventory, list *switchv1beta1.SwitchList) (s *switchv1beta1.Switch) {
	for _, item := range list.Items {
		if item.Spec.UUID == inv.Name {
			s = &item
			break
		}
	}
	return
}
