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
	"github.com/onmetal/metal-api/common/types/base"
	usecase "github.com/onmetal/metal-api/scheduler/usecase/order"
	ctrl "sigs.k8s.io/controller-runtime"
)

type InstanceScheduler struct {
	log               logr.Logger
	instanceScheduler usecase.InstanceSchedulerUseCase
}

func NewInstanceScheduler(log logr.Logger,
	instanceScheduler usecase.InstanceSchedulerUseCase) *InstanceScheduler {
	return &InstanceScheduler{
		log:               log,
		instanceScheduler: instanceScheduler,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceScheduler) SetupWithManager(mgr ctrl.Manager) error {
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

func (r *InstanceScheduler) Reconcile(_ context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.log.WithValues("namespace", req.NamespacedName)

	instanceMetadata := base.NewInstanceMetadata(req.Name, req.Namespace)
	if err := r.instanceScheduler.Execute(instanceMetadata); err != nil {
		reqLogger.Info("InstanceSchedulerUseCase failed", "error", err)
		return ctrl.Result{}, err
	}
	reqLogger.Info("reconciliation finished")
	return ctrl.Result{}, nil
}
