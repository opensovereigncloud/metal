/*
Copyright 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	switchv1alpha1 "github.com/onmetal/metal-api/apis/switches/v1alpha1"
)

const CSwitchAssignmentFinalizer = "switchassignments.switch.onmetal.de/finalizer"

type SwitchAssignmentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switchassignments,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switchassignments/status,verbs=update
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switchassignments/finalizers,verbs=update
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *SwitchAssignmentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("switch assignment", req.NamespacedName)
	res := &switchv1alpha1.SwitchAssignment{}
	if err := r.Get(ctx, req.NamespacedName, res); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("requested resource not found")
		} else {
			log.Error(err, "failed to get requested resource")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if res.DeletionTimestamp != nil {
		return r.finalize(ctx, res)
	}

	switch controllerutil.ContainsFinalizer(res, CSwitchAssignmentFinalizer) {
	case false:
		controllerutil.AddFinalizer(res, CSwitchAssignmentFinalizer)
		if err := r.Update(ctx, res); err != nil {
			log.Error(err, "failed to update switchAssignment resource")
			return ctrl.Result{}, err
		}
	case true:
		if res.Status.State == switchv1alpha1.CEmptyString {
			res.FillStatus(switchv1alpha1.CAssignmentStatePending, &switchv1alpha1.LinkedSwitchSpec{})
			if err := r.Status().Update(ctx, res); err != nil {
				log.Error(err, "failed to set status on resource creation", "name", res.NamespacedName())
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchAssignmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1alpha1.SwitchAssignment{}).
		Complete(r)
}

func (r *SwitchAssignmentReconciler) finalize(ctx context.Context, res *switchv1alpha1.SwitchAssignment) (ctrl.Result, error) {
	if controllerutil.ContainsFinalizer(res, CSwitchAssignmentFinalizer) {
		res.FillStatus(switchv1alpha1.CStateDeleting, &switchv1alpha1.LinkedSwitchSpec{})
		if err := r.Status().Update(ctx, res); err != nil {
			r.Log.Error(err, "failed to finalize resource", "name", res.NamespacedName())
			return ctrl.Result{}, err
		}

		controllerutil.RemoveFinalizer(res, CSwitchAssignmentFinalizer)
		if err := r.Update(ctx, res); err != nil {
			r.Log.Error(err, "failed to update resource on finalizer removal", "name", res.NamespacedName())
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}
