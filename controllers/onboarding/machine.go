// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventories "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1lpha3 "github.com/onmetal/metal-api/apis/machine/v1alpha3"
	domain "github.com/onmetal/metal-api/domain/inventory"
	usecase "github.com/onmetal/metal-api/usecase/onboarding"
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
			&machinev1lpha3.Machine{},
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
			&inventories.Inventory{},
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
		For(&inventories.Inventory{}).
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
