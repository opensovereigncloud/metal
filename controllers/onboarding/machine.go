// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	domain "github.com/ironcore-dev/metal/domain/inventory"
	usecase "github.com/ironcore-dev/metal/usecase/onboarding"
	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
)

type OnboardingMachineReconciler struct {
	log                   logr.Logger
	getMachineUseCase     usecase.GetMachine
	getInventoryUseCase   usecase.GetInventory
	createMachine         usecase.CreateMachine
	machineOnboardUseCase usecase.MachineOnboarding
}

func NewOnboardingMachineReconciler(
	log logr.Logger,
	getMachineUseCase usecase.GetMachine,
	getInventoryUseCase usecase.GetInventory,
	addMachineUseCase usecase.CreateMachine,
	machineOnboardUseCase usecase.MachineOnboarding,
) *OnboardingMachineReconciler {
	return &OnboardingMachineReconciler{
		log:                   log,
		getMachineUseCase:     getMachineUseCase,
		getInventoryUseCase:   getInventoryUseCase,
		createMachine:         addMachineUseCase,
		machineOnboardUseCase: machineOnboardUseCase}
}

// SetupWithManager sets up the controller with the Manager.
func (r *OnboardingMachineReconciler) SetupWithManager(
	mgr ctrl.Manager,
) error {
	r.log.Info("reconciler started")
	if err := mgr.
		GetFieldIndexer().
		IndexField(
			context.Background(),
			&metalv1alpha4.Machine{},
			"metadata.name",
			machineIndex,
		); err != nil {
		r.log.Error(err, "unable to setup machine index field")
		return err
	}
	if err := mgr.
		GetFieldIndexer().
		IndexField(
			context.Background(),
			&metalv1alpha4.Inventory{},
			"metadata.name",
			inventoryIndex,
		); err != nil {
		r.log.Error(err, "unable to setup inventory index field")
		return err
	}

	if err := mgr.
		GetFieldIndexer().
		IndexField(
			context.Background(),
			&ipamv1alpha1.IP{},
			"metadata.name",
			ipIndex,
		); err != nil {
		r.log.Error(err, "unable to setup ip index field")
		return err
	}
	if err := mgr.
		GetFieldIndexer().
		IndexField(
			context.Background(),
			&ipamv1alpha1.Subnet{},
			"metadata.name",
			subnetIndex,
		); err != nil {
		r.log.Error(err, "unable to setup subnet index field")
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&metalv1alpha4.Inventory{}).
		Complete(r)
}

func (r *OnboardingMachineReconciler) Reconcile(
	_ context.Context,
	req ctrl.Request,
) (ctrl.Result, error) {
	reqLogger := r.log.WithValues("namespace", req.NamespacedName)

	inventory, err := r.getInventoryUseCase.Execute(req.Name)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if !inventory.IsMachine() {
		return ctrl.Result{}, nil
	}

	machine, err := r.getMachineUseCase.Execute(inventory.UUID)
	if usecase.IsNotFound(err) {
		return ctrl.Result{Requeue: true}, r.CreateMachine(inventory)
	}
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := r.machineOnboardUseCase.Execute(machine, inventory); err != nil {
		reqLogger.Info("can't gather the information", "error", err)
		return ctrl.Result{}, err
	}

	reqLogger.Info("reconciliation finished")
	return ctrl.Result{}, nil
}

func (r *OnboardingMachineReconciler) CreateMachine(inventory domain.Inventory) error {
	machineInfo := dto.NewMachineInfoFromInventory(inventory)
	_, err := r.createMachine.Execute(machineInfo)
	return err
}
