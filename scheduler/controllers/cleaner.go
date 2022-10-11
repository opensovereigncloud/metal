// /*
// Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
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
	"time"

	"github.com/go-logr/logr"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	usecase "github.com/onmetal/metal-api/scheduler/usecase/order"
	"github.com/onmetal/metal-api/types/common"
	ctrl "sigs.k8s.io/controller-runtime"
)

type InstanceCleaner struct {
	log             logr.Logger
	instanceCleaner usecase.OrderCleanerUseCase
}

func NewInstanceCleaner(log logr.Logger,
	instanceCleaner usecase.OrderCleanerUseCase) *InstanceCleaner {
	return &InstanceCleaner{
		log:             log,
		instanceCleaner: instanceCleaner,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceCleaner) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&machinev1alpha2.Machine{}).
		Complete(r)
}

func (r *InstanceCleaner) Reconcile(_ context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.log.WithValues("namespace", req.NamespacedName)

	instance := common.NewObjectMetadata(req.Name, req.Namespace)
	if err := r.instanceCleaner.Execute(instance); err != nil {
		reqLogger.Info("InstanceSchedulerUseCase failed", "error", err)
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: 60 * time.Minute}, nil
}
