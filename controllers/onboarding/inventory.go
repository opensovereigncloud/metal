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
	usecase "github.com/onmetal/metal-api/usecase/onboarding"
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	oob "github.com/onmetal/oob-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type InventoryOnboardingReconciler struct {
	log             logr.Logger
	createInventory usecase.CreateInventory
	getServer       usecase.GetServer
}

func NewInventoryOnboardingReconciler(
	log logr.Logger,
	onboardingUseCase usecase.CreateInventory,
	getServer usecase.GetServer,
) *InventoryOnboardingReconciler {
	return &InventoryOnboardingReconciler{
		log:             log,
		createInventory: onboardingUseCase,
		getServer:       getServer,
	}
}

// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=ips,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=ips/status,verbs=get
// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets/status,verbs=get

// SetupWithManager sets up the controller with the Manager.
func (r *InventoryOnboardingReconciler) SetupWithManager(
	mgr ctrl.Manager,
) error {
	r.log.Info("reconciler started")
	if err := mgr.
		GetFieldIndexer().
		IndexField(
			context.Background(),
			&oob.OOB{},
			"metadata.name",
			oobIndex,
		); err != nil {
		r.log.Error(err, "unable to setup oob index field")
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&oob.OOB{}).
		Complete(r)
}

func (r *InventoryOnboardingReconciler) Reconcile(
	_ context.Context,
	req ctrl.Request,
) (ctrl.Result, error) {
	reqLogger := r.log.WithValues("namespace", req.NamespacedName)

	server, err := r.getServer.Execute(req.Name)
	if err != nil {
		reqLogger.Info("get server failed", "error", err)
		return ctrl.Result{}, nil
	}
	if !server.HasPowerCapabilities() {
		reqLogger.Info("no power caps", "server", server.UUID)
		return ctrl.Result{}, nil
	}

	inventoryInfo := dto.InventoryInfo{
		UUID:      server.UUID,
		Namespace: server.Namespace,
	}
	err = r.createInventory.Execute(inventoryInfo)
	if usecase.IsAlreadyCreated(err) {
		return ctrl.Result{}, nil
	}
	if err != nil {
		reqLogger.Info("inventory onboarding failed", "error", err)
		return ctrl.Result{}, err
	}

	reqLogger.Info("reconciliation finished")
	return ctrl.Result{}, nil
}
