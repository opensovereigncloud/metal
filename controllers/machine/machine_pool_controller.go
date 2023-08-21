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
	"errors"

	"github.com/go-logr/logr"
	"github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machine "github.com/onmetal/metal-api/apis/machine/v1alpha3"
	domain "github.com/onmetal/metal-api/domain/reservation"
	poolv1alpha1 "github.com/onmetal/onmetal-api/api/compute/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const MachineFinalizer = "metal-api.onmetal.de/machine-finalizer"

var (
	errorSizesListIsEmpty         = errors.New("sizes list is empty")
	errorSizesListNotFound        = errors.New("sizes list not found")
	errorMachineClassListIsEmpty  = errors.New("machine_class list is empty")
	errorMachineClassListNotFound = errors.New("machine_class list not found")
)

// MachinePoolReconciler reconciles a MachinePool object.
type MachinePoolReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Parameter Object.
type machinePoolReconcileWrappedCtx struct {
	ctx context.Context
	log logr.Logger
	req ctrl.Request
}

// +kubebuilder:rbac:groups=machine.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/finalizers,verbs=update
// +kubebuilder:rbac:groups=compute.api.onmetal.de,resources=machinepools,verbs=get;list;watch;update;patch;create;delete
// +kubebuilder:rbac:groups=compute.api.onmetal.de,resources=machinepools/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=compute.api.onmetal.de,resources=machinepools/finalizers,verbs=update
// +kubebuilder:rbac:groups=compute.api.onmetal.de,resources=machineclasses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=compute.api.onmetal.de,resources=machineclasses/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=compute.api.onmetal.de,resources=machineclasses/finalizers,verbs=update

func (r *MachinePoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	wCtx := &machinePoolReconcileWrappedCtx{
		ctx: ctx,
		req: req,
		log: r.Log.WithValues("namespace", req.NamespacedName),
	}

	machine := &machine.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(machine), machine); err != nil {
		wCtx.log.Error(err, "could not get machine")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !controllerutil.ContainsFinalizer(machine, MachineFinalizer) {
		controllerutil.AddFinalizer(machine, MachineFinalizer)
		return ctrl.Result{}, r.Client.Update(wCtx.ctx, machine)
	}

	if !machine.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.handleMachineDeletion(wCtx, machine)
	}

	sizes, err := r.resolveSizes(wCtx)
	switch err {
	case nil:
	case errorSizesListIsEmpty, errorSizesListNotFound, errorMachineClassListIsEmpty, errorMachineClassListNotFound:
		wCtx.log.Info(err.Error())
		return ctrl.Result{}, nil
	default:
		wCtx.log.Error(err, "unable to create the machine_pool")
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
		return r.createMachinePool(wCtx, machine, sizes)
	}

	if err != nil {
		wCtx.log.Error(err, "could not get pool")
		return ctrl.Result{}, err
	}

	return r.updateMachinePool(wCtx, machine, machinePool, sizes)
}

func (r *MachinePoolReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&machine.Machine{}).
		Complete(r)
}

func (r *MachinePoolReconciler) handleMachineDeletion(
	wCtx *machinePoolReconcileWrappedCtx,
	machine *machine.Machine,
) (ctrl.Result, error) {
	if controllerutil.ContainsFinalizer(machine, MachineFinalizer) {
		result, err := r.deleteMachinePool(wCtx, machine)
		if err != nil {
			return result, err
		}
	}

	controllerutil.RemoveFinalizer(machine, MachineFinalizer)
	if err := r.Client.Update(wCtx.ctx, machine); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *MachinePoolReconciler) createMachinePool(
	wCtx *machinePoolReconcileWrappedCtx,
	machine *machine.Machine,
	sizes []string,
) (ctrl.Result, error) {
	wCtx.log.Info("creating machine_pool")

	availableMachineClasses := r.getAvailableMachineClasses(machine, sizes)
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

func (r *MachinePoolReconciler) updateMachinePool(
	wCtx *machinePoolReconcileWrappedCtx,
	machine *machine.Machine,
	machinePool *poolv1alpha1.MachinePool,
	sizes []string,
) (ctrl.Result, error) {
	wCtx.log.Info("updating machine_pool")

	// if machine is booked, remove available classes
	if machine.Status.Reservation.Status != domain.ReservationStatusAvailable {
		machinePool.Status.AvailableMachineClasses = make([]corev1.LocalObjectReference, 0)

		if err := r.Status().Update(wCtx.ctx, machinePool); err != nil {
			wCtx.log.Error(err, "could not update machine_pool")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// refresh available classes
	machinePool.Status.AvailableMachineClasses = r.getAvailableMachineClasses(machine, sizes)
	if err := r.Status().Update(wCtx.ctx, machinePool); err != nil {
		wCtx.log.Error(err, "could not update machine_pool")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// nolint:unparam
func (r *MachinePoolReconciler) deleteMachinePool(
	wCtx *machinePoolReconcileWrappedCtx,
	machine *machine.Machine,
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

func (r *MachinePoolReconciler) getAvailableMachineClasses(
	machine *machine.Machine,
	sizes []string,
) []corev1.LocalObjectReference {
	var availableMachineClasses []corev1.LocalObjectReference

	availableMachineClasses = make([]corev1.LocalObjectReference, 0)
	for _, size := range sizes {
		if metav1.HasLabel(machine.ObjectMeta, v1alpha1.GetSizeMatchLabel(size)) {
			machineClass := corev1.LocalObjectReference{Name: size}
			availableMachineClasses = append(availableMachineClasses, machineClass)
		}
	}

	return availableMachineClasses
}

func (r *MachinePoolReconciler) resolveSizes(wCtx *machinePoolReconcileWrappedCtx) ([]string, error) {
	sizes := make([]string, 0)

	sizeList := &v1alpha1.SizeList{}
	err := r.List(wCtx.ctx, sizeList)
	switch {
	case err == nil:
		wCtx.log.Info("sizes list was found")
		if len(sizeList.Items) == 0 {
			wCtx.log.Info("unable to create the machine_pool. sizes list is empty")
			return nil, errorSizesListIsEmpty
		}
	case apierrors.IsNotFound(err):
		wCtx.log.Info("the machine_pool cannot be created or updated. sizes list not found")
		return nil, errorSizesListNotFound
	default:
		wCtx.log.Error(err, "could not get sizes list")
		return nil, err
	}

	machineClassList := &poolv1alpha1.MachineClassList{}
	err = r.List(wCtx.ctx, machineClassList)
	switch {
	case err == nil:
		wCtx.log.Info("machine_class list was found")
		if len(machineClassList.Items) == 0 {
			wCtx.log.Info("unable to create the machine_pool. machine_class list is empty")
			return nil, errorMachineClassListIsEmpty
		}
	case apierrors.IsNotFound(err):
		wCtx.log.Info("the machine_pool cannot be created or updated. machine_class list not found")
		return nil, errorMachineClassListNotFound
	default:
		wCtx.log.Error(err, "could not get machine_class list")
		return nil, err
	}

	machineClassNames := make(map[string]struct{})
	for _, machineClassItem := range machineClassList.Items {
		machineClassNames[machineClassItem.Name] = struct{}{}
	}

	appendedSizes := make(map[string]struct{})
	for _, sizeListItem := range sizeList.Items {
		// avoid duplicates in sizes.
		if _, ok := appendedSizes[sizeListItem.Name]; ok {
			continue
		}

		if _, ok := machineClassNames[sizeListItem.Name]; ok {
			sizes = append(sizes, sizeListItem.Name)
			appendedSizes[sizeListItem.Name] = struct{}{}
		}
	}

	return sizes, nil
}
