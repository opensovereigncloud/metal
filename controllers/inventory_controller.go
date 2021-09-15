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
	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

// InventoryReconciler reconciles a Switch object
type InventoryReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *InventoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("inventory", req.NamespacedName)

	res := &inventoriesv1alpha1.Inventory{}
	if err := r.Get(ctx, req.NamespacedName, res); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("requested resource not found")
		} else {
			log.Error(err, "failed to get requested resource")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if res.Spec.Host.Type == string(switchv1alpha1.CMachineType) {
		return ctrl.Result{}, nil
	}
	sw := &switchv1alpha1.Switch{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: switchv1alpha1.CNamespace, Name: res.Name}, sw); err != nil {
		if apierrors.IsNotFound(err) {
			sw.Prepare(res)
			if err := r.Client.Create(ctx, sw); err != nil {
				r.Log.Error(err, "failed to create switch resource", "name", sw.NamespacedName())
				return ctrl.Result{}, err
			}
		} else {
			log.Error(err, "failed to lookup for switch resource")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InventoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&inventoriesv1alpha1.Inventory{}).
		Complete(r)
}
