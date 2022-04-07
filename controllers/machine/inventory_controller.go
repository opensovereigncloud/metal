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
	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinerr "github.com/onmetal/metal-api/pkg/errors"
	"github.com/onmetal/metal-api/pkg/inventory"
	"github.com/onmetal/metal-api/pkg/machine"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type InventoryReconciler struct {
	client.Client

	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
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

//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/finalizers,verbs=update

func (r *InventoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("inventory", req.NamespacedName)

	i, err := inventory.New(ctx, r.Client, reqLogger, r.Recorder, req)
	if err != nil {
		return machinerr.GetResultForError(reqLogger, err)
	}

	machineObj, getErr := i.Machiner.GetMachine(i.Spec.System.ID, i.Namespace)
	if getErr != nil {
		return machinerr.GetResultForError(reqLogger, getErr)
	}

	if err := i.UpdateMachine(machineObj); err != nil {
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
	invObj, ok := e.Object.(*inventoriesv1alpha1.Inventory)
	if !ok {
		r.Log.Info("inventory cast failed")
		return false
	}

	mm := machine.New(ctx, r.Client, r.Log, r.Recorder)
	machineObj, err := mm.GetMachine(invObj.Spec.System.ID, invObj.Namespace)
	if err != nil {
		r.Log.Info("failed to retrieve machine obkect from cluster", "error", err)
		return false
	}

	machineObj.Status.Inventory.Exist = false
	machineObj.Status.Inventory.Reference = nil
	if updErr := mm.UpdateStatus(machineObj); updErr != nil {
		r.Log.Info("can't update machine status for inventory", "error", updErr)
		return false
	}
	return false
}
