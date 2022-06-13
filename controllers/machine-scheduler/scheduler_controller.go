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

	"github.com/go-logr/logr"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/internal/entity"
	"github.com/onmetal/metal-api/internal/usecase"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// SchedulerReconciler reconciles a Request object
type SchedulerReconciler struct {
	ctrlclient.Client

	Log       logr.Logger
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Scheduler usecase.Scheduler
}

// SetupWithManager sets up the controller with the Manager.
func (r *SchedulerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&machinev1alpha2.MachineAssignment{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *SchedulerReconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: r.onUpdate,
		DeleteFunc: r.onDelete,
	}
}

//+kubebuilder:rbac:groups=machine.onmetal.de,resources=requests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=requests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=requests/finalizers,verbs=update

func (r *SchedulerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("namespace", req.NamespacedName)

	e := entity.Order{
		Name: req.Name, Namespace: req.Namespace,
	}
	if !r.Scheduler.IsScheduled(ctx, e) {
		if err := r.Scheduler.Schedule(ctx, e); err != nil {
			reqLogger.Info("machine assignment failed", "error", err)
			return ctrl.Result{}, nil
		}
	}

	reqLogger.Info("reconciliation finished")
	return ctrl.Result{}, nil
}

func (r *SchedulerReconciler) onUpdate(e event.UpdateEvent) bool {
	newObj, ok := e.ObjectNew.(*machinev1alpha2.MachineAssignment)
	if !ok {
		r.Log.Info("request delete event cast failed")
		return false
	}
	if newObj.Status.Reference != nil {
		return false
	}
	return true
}

func (r *SchedulerReconciler) onDelete(e event.DeleteEvent) bool {
	obj, ok := e.Object.(*machinev1alpha2.MachineAssignment)
	if !ok {
		r.Log.Info("request delete event cast failed")
		return false
	}

	if obj.Status.Reference == nil {
		return false
	}

	ctx := context.Background()

	reservation := entity.Reservation{
		OrderName:        obj.Name,
		OrderNamespace:   obj.Namespace,
		RequestName:      obj.Status.Reference.Name,
		RequestNamespace: obj.Status.Reference.Namespace,
	}

	if err := r.Scheduler.DeleteScheduling(ctx, reservation); err != nil {
		r.Log.Info("scheduling deletion unsuccessful", "error", err)
		return false
	}

	return false
}
