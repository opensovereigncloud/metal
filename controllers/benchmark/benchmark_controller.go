// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"reflect"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/ironcore-dev/metal/pkg/benchmark"
)

// Reconciler reconciles a Benchmark object.
type Reconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metalv1alpha4.Benchmark{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *Reconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: r.compareDifference,
		DeleteFunc: func(event.DeleteEvent) bool { return false },
	}
}

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=benchmarks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=benchmarks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=benchmarks/finalizers,verbs=update

func (r *Reconciler) Reconcile(_ context.Context, _ ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (r *Reconciler) compareDifference(e event.UpdateEvent) bool {
	oldObj, oldOk := e.ObjectOld.(*metalv1alpha4.Benchmark)
	newObj, newOk := e.ObjectNew.(*metalv1alpha4.Benchmark)
	if !oldOk || !newOk {
		r.Log.Info("compare failed")
		return false
	}
	if reflect.DeepEqual(oldObj.Spec, newObj.Spec) {
		return false
	}
	dev := benchmark.CalculateDeviation(oldObj, newObj)
	r.Log.Info("disks deviation between old and new benchmarks", "uuid", newObj.Name,
		"value", dev)

	newObj.Status.BenchmarkDeviations = dev
	if err := r.Status().Update(context.Background(), newObj); err != nil {
		r.Log.Info("failed to update benchmark status", "error", err)
		return false
	}
	return true
}
