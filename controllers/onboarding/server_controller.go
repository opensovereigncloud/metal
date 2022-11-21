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
	"github.com/onmetal/metal-api/internal/entity"
	"github.com/onmetal/metal-api/internal/usecase"
	oobv1 "github.com/onmetal/oob-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type InventoryOnboardingReconciler struct {
	client.Client

	Log                  logr.Logger
	Scheme               *runtime.Scheme
	OnboardingRepo       usecase.Onboarding
	DestinationNamespace string
}

// SetupWithManager sets up the controller with the Manager.
func (r *InventoryOnboardingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&oobv1.OOB{}).
		Complete(r)
}

func (r *InventoryOnboardingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("namespace", req.NamespacedName)

	e := entity.Onboarding{
		RequestName:                   req.Name,
		RequestNamespace:              req.Namespace,
		InitializationObjectNamespace: r.DestinationNamespace,
	}
	is := r.OnboardingRepo.InitializationStatus(ctx, e)
	if is.Error != nil {
		reqLogger.Info("inventory initialization status", "error", is.Error)
		return ctrl.Result{}, nil
	}
	if is.Require {
		if err := r.OnboardingRepo.Initiate(ctx, e); err != nil {
			reqLogger.Info("inventory initialization failed", "error", err)
			return ctrl.Result{}, nil
		}
	}

	if err := r.OnboardingRepo.GatherData(ctx, e); err != nil {
		reqLogger.Info("can't gather the information", "error", err)
	}

	reqLogger.Info("reconciliation finished")
	return ctrl.Result{}, nil
}
