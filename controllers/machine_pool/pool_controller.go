/*
Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
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
	"github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/controllers/scheduler"
	poolv1alpha1 "github.com/onmetal/onmetal-api/api/compute/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const machineFinalizer = "metal-api.onmetal.de/machine-finalizer"

// PoolReconciler reconciles a MachinePool object
type PoolReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Parameter Object
type machinePoolReconcileWrappedCtx struct {
	ctx context.Context
	log logr.Logger
}

//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/finalizers,verbs=update
//+kubebuilder:rbac:groups=compute.api.onmetal.de,resources=machinepools,verbs=get;list;watch;update;patch;create
//+kubebuilder:rbac:groups=compute.api.onmetal.de,resources=machinepools/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=compute.api.onmetal.de,resources=machinepools/finalizers,verbs=update

func (r *PoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	wCtx := &machinePoolReconcileWrappedCtx{
		ctx: ctx,
		log: r.Log.WithValues("namespace", req.NamespacedName),
	}

	machine := &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(machine), machine); err != nil {
		wCtx.log.Error(err, "could not get machine")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !controllerutil.ContainsFinalizer(machine, machineFinalizer) {
		controllerutil.AddFinalizer(machine, machineFinalizer)
		return ctrl.Result{}, r.Client.Update(wCtx.ctx, machine)
	}

	if !machine.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.handleMachineDeletion(wCtx, machine)
	}

	sizeList := &v1alpha1.SizeList{}
	err := r.List(ctx, sizeList, client.InNamespace(req.Namespace))
	switch {
	case err == nil:
		wCtx.log.Info("sizes list was found")
		if len(sizeList.Items) == 0 {
			wCtx.log.Info("unable to create the machine_pool. sizes list is empty")
			return ctrl.Result{}, nil
		}
	case apierrors.IsNotFound(err):
		wCtx.log.Info("the machine_pool cannot be created or updated. valid sizes not found")
		return ctrl.Result{}, nil
	default:
		wCtx.log.Error(err, "could not get sizes list")
		return ctrl.Result{}, err
	}

	machinePool := &poolv1alpha1.MachinePool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	err = r.Client.Get(ctx, client.ObjectKeyFromObject(machinePool), machinePool)

	if apierrors.IsNotFound(err) {
		return r.createMachinePool(wCtx, machine, sizeList)
	}

	if err != nil {
		wCtx.log.Error(err, "could not get pool")
		return ctrl.Result{}, err
	}

	return r.updateMachinePool(wCtx, machine, machinePool, sizeList)
}

func (r *PoolReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&machinev1alpha2.Machine{}).
		Complete(r)
}

func (r *PoolReconciler) handleMachineDeletion(
	wCtx *machinePoolReconcileWrappedCtx,
	machine *machinev1alpha2.Machine,
) (ctrl.Result, error) {
	if controllerutil.ContainsFinalizer(machine, machineFinalizer) {
		result, err := r.deleteMachinePool(wCtx, machine)
		if err != nil {
			return result, err
		}
	}

	controllerutil.RemoveFinalizer(machine, machineFinalizer)
	if err := r.Client.Update(wCtx.ctx, machine); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PoolReconciler) createMachinePool(
	wCtx *machinePoolReconcileWrappedCtx,
	machine *machinev1alpha2.Machine,
	sizeList *v1alpha1.SizeList,
) (ctrl.Result, error) {
	wCtx.log.Info("creating machine_pool")

	availableMachineClasses := r.getAvailableMachineClasses(wCtx, machine, sizeList)
	if len(availableMachineClasses) == 0 {
		wCtx.log.Info("failed to create machine_pool. no available machine classes")
		return ctrl.Result{}, nil
	}

	machinePool := &poolv1alpha1.MachinePool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      machine.Name,
			Namespace: machine.Namespace,
		},
	}
	if err := r.Create(wCtx.ctx, machinePool); err != nil {
		wCtx.log.Error(err, "could not create machine_pool")
		return ctrl.Result{}, err
	}

	machinePool.Status.AvailableMachineClasses = availableMachineClasses
	if err := r.Status().Update(wCtx.ctx, machinePool); err != nil {
		wCtx.log.Error(err, "could not update machine_pool status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PoolReconciler) updateMachinePool(
	wCtx *machinePoolReconcileWrappedCtx,
	machine *machinev1alpha2.Machine,
	machinePool *poolv1alpha1.MachinePool,
	sizeList *v1alpha1.SizeList,
) (ctrl.Result, error) {
	wCtx.log.Info("updating machine_pool")

	// if machine is booked, remove available classes
	if machine.Status.Reservation.Status != scheduler.ReservationStatusAvailable {
		machinePool.Status.AvailableMachineClasses = make([]corev1.LocalObjectReference, 0)

		if err := r.Status().Update(wCtx.ctx, machinePool); err != nil {
			wCtx.log.Error(err, "could not update machine_pool")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// refresh available classes
	machinePool.Status.AvailableMachineClasses = r.getAvailableMachineClasses(wCtx, machine, sizeList)
	if err := r.Status().Update(wCtx.ctx, machinePool); err != nil {
		wCtx.log.Error(err, "could not update machine_pool")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PoolReconciler) deleteMachinePool(
	wCtx *machinePoolReconcileWrappedCtx,
	machine *machinev1alpha2.Machine,
) (ctrl.Result, error) {
	wCtx.log.Info("deleting machine_pool")

	machinePool := &poolv1alpha1.MachinePool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      machine.Name,
			Namespace: machine.Namespace,
		},
	}

	if err := r.Delete(wCtx.ctx, machinePool); err != nil {
		wCtx.log.Error(err, "could not delete machine_pool")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PoolReconciler) getAvailableMachineClasses(
	wCtx *machinePoolReconcileWrappedCtx,
	machine *machinev1alpha2.Machine,
	sizeList *v1alpha1.SizeList,
) []corev1.LocalObjectReference {
	var availableMachineClasses []corev1.LocalObjectReference

	availableMachineClasses = make([]corev1.LocalObjectReference, 0)
	for _, sizeListItem := range sizeList.Items {
		if metav1.HasLabel(machine.ObjectMeta, v1alpha1.GetSizeMatchLabel(sizeListItem.Name)) {
			machineClass := corev1.LocalObjectReference{Name: sizeListItem.Name}
			availableMachineClasses = append(availableMachineClasses, machineClass)
		}
	}

	wCtx.log.Info("matched available machine classes", "data", availableMachineClasses)

	return availableMachineClasses
}
