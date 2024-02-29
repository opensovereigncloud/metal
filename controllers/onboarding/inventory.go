// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	usecase "github.com/ironcore-dev/metal/usecase/onboarding"
	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
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
