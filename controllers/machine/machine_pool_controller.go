// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"errors"

	"github.com/go-logr/logr"
	poolv1alpha1 "github.com/ironcore-dev/ironcore/api/compute/v1alpha1"
	corev1alpha1 "github.com/ironcore-dev/ironcore/api/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	domain "github.com/ironcore-dev/metal/domain/reservation"
)

const (
	MachineFinalizer                  = "metal.ironcore.dev/machine-finalizer"
	MachinePoolOOBNameAnnotation      = "metal.ironcore.dev/oob-name"
	MachinePoolOOBNamespaceAnnotation = "metal.ironcore.dev/oob-namespace"
)

var OOBServiceAddr = poolv1alpha1.MachinePoolAddress{
	Type:    poolv1alpha1.MachinePoolInternalDNS,
	Address: "oob-console.oob.svc.cluster.local",
}

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

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machines/finalizers,verbs=update
// +kubebuilder:rbac:groups=compute.ironcore.dev,resources=machinepools,verbs=get;list;watch;update;patch;create;delete
// +kubebuilder:rbac:groups=compute.ironcore.dev,resources=machinepools/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=compute.ironcore.dev,resources=machinepools/finalizers,verbs=update
// +kubebuilder:rbac:groups=compute.ironcore.dev,resources=machineclasses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=compute.ironcore.dev,resources=machineclasses/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=compute.ironcore.dev,resources=machineclasses/finalizers,verbs=update

func (r *MachinePoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	wCtx := &machinePoolReconcileWrappedCtx{
		ctx: ctx,
		req: req,
		log: r.Log.WithValues("namespace", req.NamespacedName),
	}

	metalMachine := &metalv1alpha4.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(metalMachine), metalMachine); err != nil {
		wCtx.log.Error(err, "could not get machine")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !controllerutil.ContainsFinalizer(metalMachine, MachineFinalizer) {
		controllerutil.AddFinalizer(metalMachine, MachineFinalizer)
		return ctrl.Result{}, r.Client.Update(wCtx.ctx, metalMachine)
	}

	if !metalMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.handleMachineDeletion(wCtx, metalMachine)
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
		return r.createMachinePool(wCtx, metalMachine, sizes)
	}

	if err != nil {
		wCtx.log.Error(err, "could not get pool")
		return ctrl.Result{}, err
	}

	return r.updateMachinePool(wCtx, metalMachine, machinePool, sizes)
}

func (r *MachinePoolReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metalv1alpha4.Machine{}).
		Complete(r)
}

func (r *MachinePoolReconciler) handleMachineDeletion(
	wCtx *machinePoolReconcileWrappedCtx,
	machine *metalv1alpha4.Machine,
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
	machine *metalv1alpha4.Machine,
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
			Annotations: map[string]string{
				MachinePoolOOBNameAnnotation:      machine.Name,
				MachinePoolOOBNamespaceAnnotation: machine.Namespace,
			},
		},
	}
	if err := r.Create(wCtx.ctx, machinePool); err != nil {
		wCtx.log.Error(err, "could not create machine_pool")
		return ctrl.Result{}, err
	}

	capacity := corev1alpha1.ResourceList{}
	for _, class := range availableMachineClasses {
		capacity[corev1alpha1.ClassCountFor(corev1alpha1.ClassTypeMachineClass, class.Name)] = resource.MustParse("1")
	}

	machinePool.Status.AvailableMachineClasses = availableMachineClasses
	machinePool.Status.Addresses = []poolv1alpha1.MachinePoolAddress{OOBServiceAddr}
	machinePool.Status.Capacity = capacity
	machinePool.Status.Allocatable = capacity.DeepCopy()
	if err := r.Status().Update(wCtx.ctx, machinePool); err != nil {
		wCtx.log.Error(err, "could not update machine_pool status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *MachinePoolReconciler) updateMachinePool(
	wCtx *machinePoolReconcileWrappedCtx,
	machine *metalv1alpha4.Machine,
	machinePool *poolv1alpha1.MachinePool,
	sizes []string,
) (ctrl.Result, error) {
	wCtx.log.Info("updating machine_pool")

	// ensure correct annotations
	if machinePool.Annotations == nil {
		machinePool.Annotations = make(map[string]string, 2)
	}
	if machinePool.Annotations[MachinePoolOOBNameAnnotation] != machine.Name ||
		machinePool.Annotations[MachinePoolOOBNamespaceAnnotation] != machine.Namespace {
		machinePool.Annotations[MachinePoolOOBNameAnnotation] = machine.Name
		machinePool.Annotations[MachinePoolOOBNamespaceAnnotation] = machine.Namespace

		if err := r.Update(wCtx.ctx, machinePool); err != nil {
			wCtx.log.Error(err, "could not update machine_pool")
			return ctrl.Result{}, err
		}
	}

	machinePool.Status.Addresses = r.ensureAddress(machinePool.Status.Addresses, OOBServiceAddr)

	// if machine is booked, remove available classes
	if machine.Status.Reservation.Status != domain.ReservationStatusAvailable {
		allocatable := machinePool.Status.Capacity.DeepCopy()
		for name := range allocatable {
			allocatable[name] = resource.MustParse("0")
		}
		machinePool.Status.Allocatable = allocatable
		machinePool.Status.AvailableMachineClasses = nil
		if err := r.Status().Update(wCtx.ctx, machinePool); err != nil {
			wCtx.log.Error(err, "could not update machine_pool status")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// refresh available classes
	machinePool.Status.Allocatable = machinePool.Status.Capacity.DeepCopy()
	machinePool.Status.AvailableMachineClasses = r.getAvailableMachineClasses(machine, sizes)
	if err := r.Status().Update(wCtx.ctx, machinePool); err != nil {
		wCtx.log.Error(err, "could not update machine_pool status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// nolint:unparam
func (r *MachinePoolReconciler) deleteMachinePool(
	wCtx *machinePoolReconcileWrappedCtx,
	machine *metalv1alpha4.Machine,
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
	machine *metalv1alpha4.Machine,
	sizes []string,
) []corev1.LocalObjectReference {
	var availableMachineClasses []corev1.LocalObjectReference

	availableMachineClasses = make([]corev1.LocalObjectReference, 0)
	for _, size := range sizes {
		if metav1.HasLabel(machine.ObjectMeta, metalv1alpha4.GetSizeMatchLabel(size)) {
			machineClass := corev1.LocalObjectReference{Name: size}
			availableMachineClasses = append(availableMachineClasses, machineClass)
		}
	}

	return availableMachineClasses
}

func (r *MachinePoolReconciler) resolveSizes(wCtx *machinePoolReconcileWrappedCtx) ([]string, error) {
	sizes := make([]string, 0)

	sizeList := &metalv1alpha4.SizeList{}
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

func (r *MachinePoolReconciler) ensureAddress(addrs []poolv1alpha1.MachinePoolAddress, addr poolv1alpha1.MachinePoolAddress) []poolv1alpha1.MachinePoolAddress {
	for _, a := range addrs {
		if a == addr {
			return addrs
		}
	}
	return append(addrs, addr)
}
