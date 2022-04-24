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
	"time"

	"github.com/go-logr/logr"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	requestv1alpha1 "github.com/onmetal/metal-api/apis/request/v1alpha1"
	machinerr "github.com/onmetal/metal-api/pkg/errors"
	"github.com/onmetal/metal-api/pkg/provider"
	"github.com/onmetal/metal-api/pkg/reserve"
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
		DeleteFunc: r.recreateObject,
	}
}

//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create

func (r *MachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("machine", req.NamespacedName)

	machineObj := &machinev1alpha2.Machine{}
	if err := provider.GetObject(ctx, req.Name, req.Namespace, r.Client, machineObj); err != nil {
		return machinerr.GetResultForError(reqLogger, err)
	}

	mm := reserve.NewMachineReserver(ctx, r.Client, reqLogger, r.Recorder, machineObj)
	key, ok := machineObj.Labels[machinev1alpha2.LeasedLabel]
	if key == "true" && ok && machineObj.Status.RequestState != machinev1alpha2.RequestStateRunning {
		if err := mm.CheckIn(); err != nil {
			return machinerr.GetResultForError(reqLogger, err)
		}
	} else if !ok {
		if err := mm.CheckOut(); err != nil {
			return machinerr.GetResultForError(reqLogger, err)
		}
	}
	metalName, metalNamespace := machineObj.Labels[machinev1alpha2.MetalRequestLabel], machineObj.Namespace
	request := &requestv1alpha1.Request{}
	if err := provider.GetObject(ctx, metalName, metalNamespace, r.Client, request); err != nil {
		return machinerr.GetResultForError(reqLogger, err)
	}

	if request.Status.State != machineObj.Status.RequestState {
		request.Status.State = machineObj.Status.RequestState
	}

	if updErr := r.Client.Status().Update(ctx, request); updErr != nil {
		reqLogger.Info("unable to update request state", "error", updErr)
	}

	if key == "true" && ok {
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
	}

	return machinerr.GetResultForError(reqLogger, nil)
}

func (r *MachineReconciler) recreateObject(e event.DeleteEvent) bool {
	machineObj, ok := e.Object.(*machinev1alpha2.Machine)
	if !ok {
		return false
	}
	machineObj.ResourceVersion = ""

	if err := r.Client.Create(context.Background(), machineObj); err != nil {
		r.Log.Info("failed to revert deletion of machine instance", "error", err)
		return false
	}
	return false
}
