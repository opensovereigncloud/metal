/*
Copyright 2021.

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
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

type SwitchAssignmentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switchassignments,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=list;update
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/status,verbs=update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *SwitchAssignmentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("switchAssignment", req.NamespacedName)
	assignmentRes := &switchv1alpha1.SwitchAssignment{}
	err := r.Get(ctx, req.NamespacedName, assignmentRes)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Error(err, "requested switch assignment resource not found", "name", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to get switchAssignment resource", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	selector := labels.SelectorFromSet(labels.Set{switchv1alpha1.LabelChassisId: assignmentRes.Labels[switchv1alpha1.LabelChassisId]})
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         1000,
	}
	switchesList := &switchv1alpha1.SwitchList{}
	if err := r.List(ctx, switchesList, opts); err != nil {
		log.Error(err, "unable to get switches list")
	}
	if len(switchesList.Items) == 0 {
		return ctrl.Result{RequeueAfter: switchv1alpha1.CAssignmentRequeueInterval}, nil
	} else {
		targetSwitch := &switchesList.Items[0]
		targetSwitch.Spec.State.ConnectionLevel = 0
		targetSwitch.Spec.State.Role = switchv1alpha1.CSpineRole
		if err := r.Update(ctx, targetSwitch); err != nil {
			log.Error(err, "unable to update switch resource status", "name", types.NamespacedName{
				Namespace: targetSwitch.Namespace,
				Name:      targetSwitch.Name,
			})
			return ctrl.Result{}, err
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
