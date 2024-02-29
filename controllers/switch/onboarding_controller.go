// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/reference"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/ironcore-dev/metal/pkg/constants"
)

// OnboardingReconciler reconciles NetworkSwitch object corresponding
// to given Inventory object.
type OnboardingReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=networkswitches,verbs=get;list;create
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=inventories,verbs=get;list;watch

func (r *OnboardingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	obj := &metalv1alpha4.Inventory{}
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

func (r *OnboardingReconciler) reconcileRequired(ctx context.Context, obj *metalv1alpha4.Inventory) (ctrl.Result, error) {
	if !obj.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, nil
	}
	return r.onboarding(ctx, obj)
}

// SetupWithManager sets up the controller with the Manager.
func (r *OnboardingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	labelPredicate, err := predicate.LabelSelectorPredicate(metav1.LabelSelector{
		MatchLabels: map[string]string{constants.SizeLabel: "true"},
	})
	if err != nil {
		r.Log.Error(err, "failed to setup predicates")
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&metalv1alpha4.Inventory{}).
		WithOptions(controller.Options{
			RecoverPanic: ptr.To(true),
		}).
		WithEventFilter(predicate.And(labelPredicate)).
		Complete(r)
}

func (r *OnboardingReconciler) onboarding(ctx context.Context, inv *metalv1alpha4.Inventory) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	key := client.ObjectKeyFromObject(inv)
	targetSwitch := &metalv1alpha4.NetworkSwitch{}
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

func (r *OnboardingReconciler) onboardNewSwitch(ctx context.Context, inv *metalv1alpha4.Inventory) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)
	targetSwitch := &metalv1alpha4.NetworkSwitch{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NetworkSwitch",
			APIVersion: "metal.ironcore.dev/v1alpha4",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      inv.Name,
			Namespace: inv.Namespace,
		},
		Spec: metalv1alpha4.NetworkSwitchSpec{
			Managed:   ptr.To(true),
			Cordon:    ptr.To(false),
			TopSpine:  ptr.To(false),
			ScanPorts: ptr.To(true),
		},
	}
	targetSwitch.SetInventoryRef(inv.Name)
	targetSwitch.UpdateSwitchAnnotations(inv)
	targetSwitch.UpdateSwitchLabels(inv)
	log.Info("creating NetworkSwitch object")
	err := r.Create(ctx, targetSwitch)
	return ctrl.Result{}, err
}
