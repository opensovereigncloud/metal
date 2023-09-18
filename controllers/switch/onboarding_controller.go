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

	"github.com/go-logr/logr"
	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/pkg/constants"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/reference"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// OnboardingReconciler reconciles Switch object corresponding
// to given Inventory object.
type OnboardingReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;create
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch

func (r *OnboardingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	obj := &inventoryv1alpha1.Inventory{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	ref, err := reference.GetReference(r.Scheme, obj)
	if err != nil {
		r.Log.Error(err, "failed to construct reference", "request", req.NamespacedName)
		return ctrl.Result{}, err
	}

	log := r.Log.WithValues("object", *ref)
	log.Info("reconciliation started")
	requestCtx := logr.NewContext(ctx, log)
	return r.reconcileRequired(requestCtx, obj)
}

func (r *OnboardingReconciler) reconcileRequired(ctx context.Context, obj *inventoryv1alpha1.Inventory) (ctrl.Result, error) {
	if !obj.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, nil
	}
	return r.onboarding(ctx, obj)
}

// SetupWithManager sets up the controller with the Manager.
func (r *OnboardingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	labelPredicate, err := predicate.LabelSelectorPredicate(metav1.LabelSelector{
		MatchLabels: map[string]string{constants.SizeLabel: ""},
	})
	if err != nil {
		r.Log.Error(err, "failed to setup predicates")
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&inventoryv1alpha1.Inventory{}).
		WithOptions(controller.Options{
			RecoverPanic: pointer.Bool(true),
		}).
		WithEventFilter(predicate.And(labelPredicate)).
		Complete(r)
}

func (r *OnboardingReconciler) onboarding(ctx context.Context, inv *inventoryv1alpha1.Inventory) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	key := client.ObjectKeyFromObject(inv)
	targetSwitch := &switchv1beta1.Switch{}
	if err := r.Get(ctx, key, targetSwitch); err != nil {
		switch {
		case apierrors.IsNotFound(err):
			log.Info("onboarding started")
			return r.onboardNewSwitch(ctx, inv)
		default:
			return ctrl.Result{}, err
		}
	}
	log.Info("onboarding finished")
	return ctrl.Result{}, nil
}

func (r *OnboardingReconciler) onboardNewSwitch(ctx context.Context, inv *inventoryv1alpha1.Inventory) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	targetSwitch := &switchv1beta1.Switch{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Switch",
			APIVersion: "switch.onmetal.de/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      inv.Name,
			Namespace: inv.Namespace,
		},
		Spec: switchv1beta1.SwitchSpec{
			Managed:   pointer.Bool(true),
			Cordon:    pointer.Bool(false),
			TopSpine:  pointer.Bool(false),
			ScanPorts: pointer.Bool(true),
		},
	}
	targetSwitch.SetInventoryRef(inv.Name)
	targetSwitch.UpdateSwitchAnnotations(inv)
	targetSwitch.UpdateSwitchLabels(inv)
	log.Info("creating Switch object")
	err := r.Create(ctx, targetSwitch)
	return ctrl.Result{}, err
}
