// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/ironcore-dev/metal/pkg/benchmark"
	machinerr "github.com/ironcore-dev/metal/pkg/errors"
)

type OnboardingReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *OnboardingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metalv1alpha4.Inventory{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *OnboardingReconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool { return false },
		DeleteFunc: func(e event.DeleteEvent) bool { return false },
	}
}

func (r *OnboardingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("benchmark onboarding", req.NamespacedName)

	b, exist := benchmark.New(ctx, r.Client, reqLogger, req)
	if exist {
		return machinerr.GetResultForError(reqLogger, nil)
	}

	if err := b.Create(); err != nil {
		return machinerr.GetResultForError(reqLogger, err)
	}

	return machinerr.GetResultForError(reqLogger, nil)
}
