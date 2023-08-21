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
	"reflect"

	"github.com/go-logr/logr"
	benchmarkv1alpha3 "github.com/onmetal/metal-api/apis/benchmark/v1alpha3"
	"github.com/onmetal/metal-api/pkg/benchmark"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// Reconciler reconciles a Machine object.
type Reconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&benchmarkv1alpha3.Machine{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *Reconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: r.compareDifference,
		DeleteFunc: func(event.DeleteEvent) bool { return false },
	}
}

// +kubebuilder:rbac:groups=benchmark.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=benchmark.onmetal.de,resources=machines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=benchmark.onmetal.de,resources=machines/finalizers,verbs=update

func (r *Reconciler) Reconcile(_ context.Context, _ ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (r *Reconciler) compareDifference(e event.UpdateEvent) bool {
	oldObj, oldOk := e.ObjectOld.(*benchmarkv1alpha3.Machine)
	newObj, newOk := e.ObjectNew.(*benchmarkv1alpha3.Machine)
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

	newObj.Status.MachineDeviation = dev
	if err := r.Status().Update(context.Background(), newObj); err != nil {
		r.Log.Info("failed to update benchmark status", "error", err)
		return false
	}
	return true
}
