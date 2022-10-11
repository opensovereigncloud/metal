/*
Copyright 2022.

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
	"time"

	"github.com/go-logr/logr"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	domain "github.com/onmetal/metal-api/scheduler/domain/order"
	usecase "github.com/onmetal/metal-api/scheduler/usecase/order"
	"github.com/onmetal/metal-api/types/common"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// Scheduler reconciles an Order object.
type Scheduler struct {
	log                     logr.Logger
	orderAlreadyScheduled   usecase.OrderAlreadyScheduledUseCase
	cancelOrder             usecase.CancelOrderUseCase
	instanceForOrderUseCase usecase.FindInstanceForOrderUseCase
}

func NewSchedulerController(
	log logr.Logger,
	orderAlreadyScheduled usecase.OrderAlreadyScheduledUseCase,
	cancelOrder usecase.CancelOrderUseCase,
	instanceForOrderUseCase usecase.FindInstanceForOrderUseCase) *Scheduler {
	return &Scheduler{
		log:                     log,
		orderAlreadyScheduled:   orderAlreadyScheduled,
		cancelOrder:             cancelOrder,
		instanceForOrderUseCase: instanceForOrderUseCase,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *Scheduler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&machinev1alpha2.MachineAssignment{}).
		WithEventFilter(r.constructEventPredicates()).
		Complete(r)
}

func (r *Scheduler) constructEventPredicates() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: r.AlreadyOrdered,
		DeleteFunc: r.CancelOrder,
	}
}

func (r *Scheduler) Reconcile(_ context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.log.WithValues("namespace", req.NamespacedName)

	order := domain.NewOrder(req.Name, req.Namespace)

	if r.orderAlreadyScheduled.Invoke(order) {
		reqLogger.Info("reconciliation not needed, order already scheduled")
		return ctrl.Result{}, nil
	}

	orderScheduler, err := r.instanceForOrderUseCase.Execute(order)
	if err != nil {
		reqLogger.Info("instanceForOrderUseCase failed", "error", err)
		if usecase.IsVacantInstanceNotFound(err) {
			return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
		}
		return ctrl.Result{}, err
	}

	if err := orderScheduler.Schedule(); err != nil {
		reqLogger.Info("OrderScheduler failed", "error", err)
		return ctrl.Result{}, err
	}

	reqLogger.Info("reconciliation finished")
	return ctrl.Result{}, nil
}

func (r *Scheduler) AlreadyOrdered(e event.UpdateEvent) bool {
	newObj, ok := e.ObjectNew.(*machinev1alpha2.MachineAssignment)
	if !ok {
		r.log.Info("request delete event cast failed")
		return false
	}
	if newObj.Status.MachineRef.Name != "" {
		return false
	}
	return true
}

func (r *Scheduler) CancelOrder(e event.DeleteEvent) bool {
	obj, ok := e.Object.(*machinev1alpha2.MachineAssignment)
	if !ok {
		r.log.Info("cancelOrder: delete event cast failed", "object", e.Object.GetName())
		return false
	}

	if obj.Status.MachineRef.Name == "" {
		return false
	}

	instanceName := obj.Status.MachineRef.Name
	instanceNamespace := obj.Status.MachineRef.Namespace
	instanceMeta := common.NewObjectMetadata(instanceName, instanceNamespace)

	if err := r.cancelOrder.Cancel(instanceMeta); err != nil {
		r.log.Info("cancelOrder: failed", "error", err)
		return false
	}
	return false
}
