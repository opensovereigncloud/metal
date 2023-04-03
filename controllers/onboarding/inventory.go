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
	usecase "github.com/onmetal/metal-api/internal/usecase/onboarding"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
	oobv1 "github.com/onmetal/oob-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type InventoryOnboardingReconciler struct {
	log               logr.Logger
	onboardingUseCase usecase.OnboardingUseCase
}

func NewInventoryOnboardingReconciler(
	log logr.Logger,
	onboardingUseCase usecase.OnboardingUseCase) *InventoryOnboardingReconciler {
	return &InventoryOnboardingReconciler{log: log, onboardingUseCase: onboardingUseCase}
}

// SetupWithManager sets up the controller with the Manager.
func (r *InventoryOnboardingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.log.Info("reconciler started")
	return ctrl.NewControllerManagedBy(mgr).
		For(&oobv1.OOB{}).
		Complete(r)
}

func (r *InventoryOnboardingReconciler) Reconcile(_ context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.log.WithValues("namespace", req.NamespacedName)

	request := dto.Request{
		Name:      req.Name,
		Namespace: req.Namespace,
	}
	err := r.onboardingUseCase.Execute(request)
	if usecase.IsAlreadyOnboarded(err) {
		return ctrl.Result{}, nil
	}
	if err != nil {
		reqLogger.Info("inventory onboarding failed", "error", err)
		return ctrl.Result{}, err
	}

	reqLogger.Info("reconciliation finished")
	return ctrl.Result{}, nil
}
