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

package scheduler

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

// Reconciler reconciles a Ignition object.
type Reconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha2.MachineAssignment{}).
		Complete(r)
}

//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/finalizers,verbs=update
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machineassignments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machineassignments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machineassignments/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("namespace", req.NamespacedName, "machineAssignment", req.Name)

	reqLogger.Info("reconciling")
	machineAssignment := &v1alpha2.MachineAssignment{}
	if err := r.Client.Get(ctx, req.NamespacedName, machineAssignment); err != nil {
		r.Log.Error(err, "unable to fetch MachineAssignment")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// check prerequisites
	if machineAssignment.Spec.MachineSize == "" {
		err := errors.New("MachineSize is not set")
		return ctrl.Result{}, err
	}

	// examine DeletionTimestamp to determine if object is under deletion
	if machineAssignment.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(machineAssignment, SchedulerFinalizer) {
			controllerutil.AddFinalizer(machineAssignment, SchedulerFinalizer)
			if err := r.Client.Update(ctx, machineAssignment); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(machineAssignment, SchedulerFinalizer) {
			// our finalizer is present, so lets handle any external dependency
			err := r.ignitionCleanup(ctx, machineAssignment)
			if err != nil {
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(machineAssignment, SchedulerFinalizer)
			if err = r.Client.Update(ctx, machineAssignment); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	// if there is no status we need to find a metal machine
	if machineAssignment.Status.State == "" {
		machine, err := r.findAvailableMachine(ctx, machineAssignment)
		if err != nil {
			return ctrl.Result{}, err
		}

		if machine == nil {
			reqLogger.Info("no available metal machines found, requeue after 60 seconds...")
			return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
		}

		// set the machine reservation to pending first
		machine.Status.Reservation.Status = ReservationStatusPending
		err = r.Client.Status().Update(ctx, machine)
		if err != nil {
			return ctrl.Result{}, errors.Wrapf(err, "failed to update available machine status of %s", machine.Name)
		}
		machineAssignment.Status.MetalComputeRef = &v1alpha2.ResourceReference{
			Name:      machine.Name,
			Namespace: machine.Namespace,
		}
		computeName, ok := machineAssignment.Labels[ComputeNameLabel]
		if !ok {
			err = errors.New(fmt.Sprintf("label %s is missing", ComputeNameLabel))
			return ctrl.Result{}, err
		}
		computeNamespace, ok := machineAssignment.Labels[ComputeNamespaceLabel]
		if !ok {
			err = errors.New(fmt.Sprintf("label %s is missing", ComputeNamespaceLabel))
			return ctrl.Result{}, err
		}
		machineAssignment.Status.OnmetalComputeRef = &v1alpha2.ResourceReference{
			Name:      computeName,
			Namespace: computeNamespace,
		}
		machineAssignment.Status.State = ReservationStatusPending
		err = r.Client.Status().Update(ctx, machineAssignment)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "failed to update status")
		}

		// we have machine reserved so check if it is running
	} else if machineAssignment.Status.State == ReservationStatusPending {
		if machineAssignment.Status.MetalComputeRef == nil {
			err := errors.New("MetalComputeRef is not set")
			return ctrl.Result{}, err
		}

		// get the referenced machine
		machine, err := r.getMachine(ctx, machineAssignment.Status.MetalComputeRef)
		if err != nil {
			return ctrl.Result{}, err
		}

		// if machine is running update the machineAssignment status
		if machine.Status.Reservation.Status == ReservationStatusRunning {
			machineAssignment.Status.State = ReservationStatusRunning

			err = r.Client.Status().Update(ctx, machineAssignment)
			if err != nil {
				return ctrl.Result{}, errors.Wrap(err, "failed to update")
			}

			// if machine is available power it on
		} else if machine.Status.Reservation.Status == ReservationStatusAvailable {
			if machine.Status.OOB.Reference == nil {
				err = errors.New("OOB reference is not set")
				return ctrl.Result{}, err
			}

			oobMachine, err := r.getOOBMachine(ctx, machine)
			if err != nil {
				return ctrl.Result{}, err
			}

			if oobMachine.Status.PowerState == "Off" && oobMachine.Spec.PowerState == "Off" {
				//TODO(flpeter) set .Spec.PowerState
				reqLogger.Info("power on", "machine", machineAssignment.Status.MetalComputeRef.Name)
			}
		}
	}

	err := r.reconcileIgnition(ctx, machineAssignment)
	if err != nil {
		return ctrl.Result{}, err
	}

	reqLogger.Info("reconciliation finished")
	return ctrl.Result{}, nil
}