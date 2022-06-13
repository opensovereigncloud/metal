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

	"github.com/go-logr/logr"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/internal/entity"
	"github.com/onmetal/metal-api/internal/usecase"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MachineReconciler struct {
	client.Client

	Log             logr.Logger
	Scheme          *runtime.Scheme
	Recorder        record.EventRecorder
	Reserver        usecase.Reserver
	Synchronization usecase.Synchronization
}

// SetupWithManager sets up the controller with the Manager.
func (r *MachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&machinev1alpha2.Machine{}).
		Complete(r)
}

//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/finalizers,verbs=update
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machineassignments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machineassignments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machineassignments/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create

func (r *MachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("namespace", req.NamespacedName)

	e := entity.Order{
		Name: req.Name, Namespace: req.Namespace}

	reservation, err := r.Reserver.GetReservation(ctx, e)
	if err != nil {
		reqLogger.Info("machine reservation process failed", "error", err)
		return ctrl.Result{}, nil
	}

	if !reservation.IsReserved() {
		if err := r.Reserver.CheckOut(ctx, reservation); err != nil {
			reqLogger.Info("machine check in failed", "error", err)
		}
		return ctrl.Result{}, nil
	}
	if reservation.IsReserved() {
		if err := r.Reserver.CheckIn(ctx, reservation); err != nil {
			reqLogger.Info("machine check in failed", "error", err)
			return ctrl.Result{}, nil
		}
	}

	if err := r.Synchronization.Do(ctx, reservation); err != nil {
		reqLogger.Info("machine status sync failed", "error", err)
		return ctrl.Result{}, nil
	}

	reqLogger.Info("reconciliation finished")
	return ctrl.Result{}, nil
}
