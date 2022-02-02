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

package controllers

import (
	"context"
	"reflect"

	"github.com/go-logr/logr"
	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	machinev1alpha1 "github.com/onmetal/metal-api/apis/machine/v1alpha1"
	machinerr "github.com/onmetal/metal-api/internal/errors"
	"github.com/onmetal/metal-api/internal/inventory"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type InventoryReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *InventoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&inventoriesv1alpha1.Inventory{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *InventoryReconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: isUpdatedOrDeleted,
		DeleteFunc: r.updateMachineStatusOnDelete,
	}
}

func (r *InventoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("inventory", req.NamespacedName)

	i, err := inventory.New(ctx, r.Client, reqLogger, req)
	if err != nil {
		return machinerr.GetResultForError(reqLogger, err)
	}
	if err := i.Update(); err != nil {
		return machinerr.GetResultForError(reqLogger, err)
	}

	return machinerr.GetResultForError(reqLogger, nil)
}

func isUpdatedOrDeleted(e event.UpdateEvent) bool {
	oldObj, oldOk := e.ObjectOld.(*inventoriesv1alpha1.Inventory)
	newObj, newOk := e.ObjectNew.(*inventoriesv1alpha1.Inventory)
	if !oldOk || !newOk {
		return false
	}
	return !reflect.DeepEqual(oldObj.Spec, newObj.Spec) ||
		!reflect.DeepEqual(oldObj.Labels, newObj.Labels) ||
		!newObj.DeletionTimestamp.IsZero()
}

func (r *InventoryReconciler) updateMachineStatusOnDelete(e event.DeleteEvent) bool {
	ctx := context.Background()
	obj := &machinev1alpha1.Machine{}
	if getErr := r.Get(ctx, types.NamespacedName{
		Name: e.Object.GetName(), Namespace: e.Object.GetNamespace(),
	}, obj); getErr != nil {
		r.Log.Info("can't retrieve machine from cluster", "error", getErr)
		return false
	}
	obj.Status.Inventory = false
	if updErr := r.Status().Update(ctx, obj); updErr != nil {
		r.Log.Info("can't update machine status", "error", updErr)
		return false
	}
	return false
}
