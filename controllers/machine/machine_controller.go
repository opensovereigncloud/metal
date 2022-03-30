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
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	machinerr "github.com/onmetal/metal-api/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type MachineReconciler struct { //nolint:revive
	client.Client

	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// SetupWithManager sets up the controller with the Manager.
func (r *MachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&machinev1alpha2.Machine{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *MachineReconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: isSpecOrStatusUpdated,
		DeleteFunc: r.recreateObject,
	}
}

//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/finalizers,verbs=update
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
//+kubebuilder:rbac:groups=oob.onmetal.de,resources=machines,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=oob.onmetal.de,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=oob.onmetal.de,resources=machines/finalizers,verbs=update
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch

func (r *MachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("machine", req.NamespacedName)

	return machinerr.GetResultForError(reqLogger, nil)
}

func isSpecOrStatusUpdated(e event.UpdateEvent) bool {
	oldObj, oldOk := e.ObjectOld.(*machinev1alpha2.Machine)
	newObj, newOk := e.ObjectNew.(*machinev1alpha2.Machine)
	if !oldOk || !newOk {
		return false
	}
	if oldObj.Status.Reboot != newObj.Status.Reboot {
		return false
	}

	return !(reflect.DeepEqual(oldObj.Spec, newObj.Spec)) ||
		!(reflect.DeepEqual(oldObj.Status, newObj.Status)) ||
		!(reflect.DeepEqual(oldObj.Labels, newObj.Labels))
}

func (r *MachineReconciler) recreateObject(e event.DeleteEvent) bool {
	machineObj, ok := e.Object.(*machinev1alpha2.Machine)
	if !ok {
		return false
	}
	machineObj.ResourceVersion = ""

	if err := r.Client.Create(context.Background(), machineObj); err != nil {
		r.Log.Info("failed to revert deletion machine instance", "error", err)
		return false
	}
	return false
}
