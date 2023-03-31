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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
)

type InventoryReconciler struct {
	client.Client

	Log       logr.Logger
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Namespace string
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

// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/finalizers,verbs=update

func (r *InventoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, nil
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

	machine := &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      invObj.Spec.System.ID,
			Namespace: invObj.Namespace,
		},
	}
	err := r.Client.Get(ctx, client.ObjectKeyFromObject(machine), machine)
	if err != nil {
		r.Log.Info("failed to retrieve machine object from cluster", "error", err)
		return false
	}

	machine.Status.Inventory.Exist = false
	machine.Status.Inventory.Reference = nil

	if updErr := r.Client.Status().Update(ctx, machine); updErr != nil {
		r.Log.Info("can't update machine status for inventory", "error", updErr)
		return false
	}
	return false
}
