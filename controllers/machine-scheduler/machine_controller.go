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
	machinerr "github.com/onmetal/metal-api/pkg/errors"
	"github.com/onmetal/metal-api/pkg/provider"
	"github.com/onmetal/metal-api/pkg/reserve"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MachineReconciler struct {
	client.Client

	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
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
	reqLogger := r.Log.WithValues("machine", req.NamespacedName)

	var err error

	machineObj := &machinev1alpha2.Machine{}
	if err = provider.GetObject(ctx, req.Name, req.Namespace, r.Client, machineObj); err != nil {
		return machinerr.GetResultForError(reqLogger, err)
	}

	metalName, metalNamespace := machineObj.Labels[machinev1alpha2.MetalRequestLabel], machineObj.Namespace

	var reserver reserve.Reserver //nolint:gosimple
	reserver = reserve.NewMachineReserver(ctx, r.Client, reqLogger, r.Recorder, machineObj)

	key, ok := machineObj.Labels[machinev1alpha2.LeasedLabel]
	if key == "true" && ok && machineObj.Status.Reservation.RequestState == machinev1alpha2.RequestStateReserved {
		err = reserver.CheckIn()
	} else if machineObj.Status.Reservation.RequestState == machinev1alpha2.RequestStateAvailable {
		if checkoutErr := reserver.CheckOut(); checkoutErr != nil {
			return machinerr.GetResultForError(reqLogger, checkoutErr)
		}
		return machinerr.GetResultForError(reqLogger, nil)
	}

	if err != nil {
		return machinerr.GetResultForError(reqLogger, err)
	}

	request := &machinev1alpha2.MachineAssignment{}
	if err = provider.GetObject(ctx, metalName, metalNamespace, r.Client, request); err != nil {
		if apierrors.IsNotFound(err) {
			return machinerr.GetResultForError(reqLogger, nil)
		}
		return machinerr.GetResultForError(reqLogger, err)
	}

	if syncErr := r.syncStatusState(ctx, request, machineObj); syncErr != nil {
		return machinerr.GetResultForError(reqLogger, syncErr)
	}

	return machinerr.GetResultForError(reqLogger, nil)
}

func (r *MachineReconciler) syncStatusState(ctx context.Context,
	request *machinev1alpha2.MachineAssignment, machineObj *machinev1alpha2.Machine) error {
	if request.Status.State != machineObj.Status.Reservation.RequestState {
		request.Status.State = machineObj.Status.Reservation.RequestState
	}
	if machineObj.Status.Reservation.RequestState == machinev1alpha2.RequestStateError {
		request.Status.Reference = nil
	}
	return r.Client.Status().Update(ctx, request)
}
