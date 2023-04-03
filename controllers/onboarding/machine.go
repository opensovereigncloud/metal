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
	inventories "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	usecase "github.com/onmetal/metal-api/internal/usecase/onboarding"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OnboardingMachineReconciler struct {
	log                   logr.Logger
	getMachineUseCase     usecase.GetMachineUseCase
	getInventoryUseCase   usecase.GetInventoryUseCase
	addMachineUseCase     usecase.AddMachineUseCase
	machineOnboardUseCase usecase.MachineOnboardingUseCase
}

func NewOnboardingMachineReconciler(
	log logr.Logger,
	getMachineUseCase usecase.GetMachineUseCase,
	getInventoryUseCase usecase.GetInventoryUseCase,
	addMachineUseCase usecase.AddMachineUseCase,
	machineOnboardUseCase usecase.MachineOnboardingUseCase) *OnboardingMachineReconciler {
	return &OnboardingMachineReconciler{
		log:                   log,
		getMachineUseCase:     getMachineUseCase,
		getInventoryUseCase:   getInventoryUseCase,
		addMachineUseCase:     addMachineUseCase,
		machineOnboardUseCase: machineOnboardUseCase}
}

// SetupWithManager sets up the controller with the Manager.
func (r *OnboardingMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.log.Info("reconciler started")
	return ctrl.NewControllerManagedBy(mgr).
		For(&inventories.Inventory{}).
		Complete(r)
}

func (r *OnboardingMachineReconciler) Reconcile(_ context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.log.WithValues("namespace", req.NamespacedName)

	request := dto.Request{
		Name:      req.Name,
		Namespace: req.Namespace,
	}

	inventory, err := r.getInventoryUseCase.Execute(request)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if !inventory.IsMachine() {
		return ctrl.Result{}, nil
	}

	machine, err := r.getMachineUseCase.Execute(request)
	if usecase.IsNotFound(err) {
		if err := r.AddMachineWhenNotExist(inventory); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
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

func (r *OnboardingMachineReconciler) AddMachineWhenNotExist(
	inventory dto.Inventory) error {
	return r.addMachineUseCase.Execute(inventory)
}
