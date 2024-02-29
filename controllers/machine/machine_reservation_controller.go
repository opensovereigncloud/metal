// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	computev1alpha1 "github.com/ironcore-dev/ironcore/api/compute/v1alpha1"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	domain "github.com/ironcore-dev/metal/domain/reservation"
)

// MachineReservationReconciler reconciles a MachineReservation object.
type MachineReservationReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

var (
	ErrMetalMachineNotMatchedWithComputeMachines = errors.New("metal metalv1alpha4 not matched with compute machines")
	ErrMetalMachineListEmpty                     = errors.New("metal metalv1alpha4 list is empty")
	ErrMetalMachineListNotFound                  = errors.New("metal metalv1alpha4 list not found")
)

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machines/finalizers,verbs=update
// +kubebuilder:rbac:groups=compute.ironcore.dev,resources=machines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=compute.ironcore.dev,resources=machines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=compute.ironcore.dev,resources=machines/finalizers,verbs=update

func (r *MachineReservationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("namespace", req.NamespacedName)

	computeMachine := &computev1alpha1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(computeMachine), computeMachine); err != nil {
		log.Error(err, "could not get compute metalv1alpha4")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if computeMachine.Spec.MachinePoolRef == nil {
		log.Info("compute metalv1alpha4 has empty metalv1alpha4 pool ref. skip reservation update")
		return ctrl.Result{}, nil
	}

	metalMachine, err := r.getMetalMachine(ctx, log, computeMachine)
	if err != nil {
		log.Error(err, "could not get metal metalv1alpha4")
		return ctrl.Result{}, err
	}

	if metalMachine.Status.Health != metalv1alpha4.MachineStateHealthy {
		log.Info("could not update reservation. metal metalv1alpha4 is unhealthy")
		return ctrl.Result{}, nil
	}

	metalMachine.Status.Reservation.Reference = &metalv1alpha4.ResourceReference{
		APIVersion: computeMachine.APIVersion,
		Kind:       computeMachine.Kind,
		Name:       computeMachine.Name,
		Namespace:  computeMachine.Namespace,
	}
	metalMachine.Status.Reservation.Status = domain.ReservationStatusReserved

	if err := r.Client.Status().Update(ctx, metalMachine); err != nil {
		log.Error(err, "could not update metal metalv1alpha4 status")
		return ctrl.Result{}, errors.Wrap(err, "failed to update metalv1alpha4 status")
	}

	return ctrl.Result{}, nil
}

func (r *MachineReservationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&computev1alpha1.Machine{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *MachineReservationReconciler) handleComputeMachineDeletion(e event.DeleteEvent) bool {
	ctx := context.Background()
	computeMachine, ok := e.Object.(*computev1alpha1.Machine)
	if !ok {
		r.Log.Info("compute metalv1alpha4 cast failed")
		return false
	}

	metalMachine, err := r.getMetalMachine(ctx, r.Log, computeMachine)
	if err != nil {
		r.Log.Error(err, "could not get metal metalv1alpha4")
		return false
	}

	metalMachine.Status.Reservation.Reference = nil
	metalMachine.Status.Reservation.Status = domain.ReservationStatusAvailable

	if err := r.Client.Status().Update(ctx, metalMachine); err != nil {
		r.Log.Error(err, "could not update metal metalv1alpha4 status")
		return false
	}

	return false
}

func (r *MachineReservationReconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: r.handleComputeMachineDeletion,
	}
}

func (r *MachineReservationReconciler) getMetalMachine(
	ctx context.Context,
	log logr.Logger,
	computeMachine *computev1alpha1.Machine,
) (*metalv1alpha4.Machine, error) {
	metalMachinesList := &metalv1alpha4.MachineList{}
	err := r.List(ctx, metalMachinesList)
	switch {
	case err == nil:
		log.Info("metal machines list was found")
		if len(metalMachinesList.Items) == 0 {
			log.Info("unable to create metalv1alpha4 reservation. metal machines list is empty")
			return nil, ErrMetalMachineListEmpty
		}
	case apierrors.IsNotFound(err):
		log.Info("metal machines list not found")
		return nil, ErrMetalMachineListNotFound
	default:
		log.Error(err, "could not get metal machines list")
		return nil, err
	}

	for _, metalMachine := range metalMachinesList.Items {
		if metalMachine.Name == computeMachine.Spec.MachinePoolRef.Name {
			log.Info("metal metalv1alpha4 matched with metalv1alpha4 pool name")
			return &metalMachine, nil
		}
	}

	return nil, ErrMetalMachineNotMatchedWithComputeMachines
}
