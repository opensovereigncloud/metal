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
	"github.com/onmetal/metal-api/controllers/scheduler"
	oobv1 "github.com/onmetal/oob-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const maintainedMachineLabel = "onmetal.de/oob-ignore"

// OOBReconciler reconciles a Machine object.
type OOBReconciler struct {
	client.Client

	Log       logr.Logger
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Namespace string
}

// SetupWithManager sets up the controller with the Manager.
func (r *OOBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&oobv1.OOB{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *OOBReconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: isUUIDExist,
		UpdateFunc: onUpdate,
		DeleteFunc: r.onDelete,
	}
}

// +kubebuilder:rbac:groups=onmetal.de,resources=oobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=onmetal.de,resources=oobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=onmetal.de,resources=oobs/finalizers,verbs=update

func (r *OOBReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("namespace", req.NamespacedName)

	oobObj := &oobv1.OOB{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(oobObj), oobObj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	machineObj := &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      oobObj.Status.UUID,
			Namespace: r.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(machineObj), machineObj); err != nil {
		reqLogger.Info("no machine for oob", "error", err)
		return ctrl.Result{}, nil
	}

	updateTaints(oobObj, machineObj)

	if specUpdErr := r.Client.Update(ctx, machineObj); specUpdErr != nil {
		reqLogger.Info("machine taints update failed", "error", specUpdErr)
		return ctrl.Result{}, nil
	}

	if !machineObj.Status.OOB.Exist {
		machineObj.Status.OOB = prepareReferenceSpec(oobObj)
	}

	previousReservationStatus := machineObj.DeepCopy().Status.Reservation.Status
	syncStatusState(oobObj, machineObj)

	if statusUpdErr := r.Client.Status().Update(ctx, machineObj); statusUpdErr != nil {
		return ctrl.Result{}, statusUpdErr
	}

	// if machine previous reservation status changed or reservation status is available
	// trigger reconcile loop on machine assignment
	if previousReservationStatus != machineObj.Status.Reservation.Status ||
		machineObj.Status.Reservation.Status == scheduler.ReservationStatusAvailable {
		// only get assignment machine if reference is set
		if machineObj.Status.Reservation.Reference != nil {
			machineAssignment := &machinev1alpha2.MachineAssignment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      machineObj.Status.Reservation.Reference.Name,
					Namespace: machineObj.Status.Reservation.Reference.Namespace,
				},
			}

			// get the machine assignment
			err := r.Client.Get(ctx, client.ObjectKeyFromObject(machineAssignment), machineAssignment)
			if err != nil {
				reqLogger.Info("failed to get MachineAssignment",
					"name", machineAssignment.Name, "namespace", machineAssignment.Namespace)
				return ctrl.Result{}, err
			}

			// update the label to force reconcile if machine assignment is pending
			if machineAssignment.Status.State == scheduler.ReservationStatusPending {
				machineAssignment.Labels[scheduler.SchedulerReconcileLabel] = time.Now().String()
				err = r.Client.Update(ctx, machineAssignment)
				if err != nil {
					reqLogger.Info("failed to update MachineAssignment",
						"name", machineAssignment.Name, "namespace", machineAssignment.Namespace)
					return ctrl.Result{}, err
				}
			}
		}
	}

	reqLogger.Info("reconciliation finished")
	return ctrl.Result{}, nil
}

func prepareReferenceSpec(oob *oobv1.OOB) machinev1alpha2.ObjectReference {
	return machinev1alpha2.ObjectReference{
		Exist: true,
		Reference: &machinev1alpha2.ResourceReference{
			Kind: oob.Kind, APIVersion: oob.APIVersion,
			Name: oob.Name, Namespace: oob.Namespace},
	}
}

func (r *OOBReconciler) onDelete(e event.DeleteEvent) bool {
	ctx := context.Background()
	obj, ok := e.Object.(*oobv1.OOB)
	if !ok {
		r.Log.Info("machine oob type assertion failed")
		return false
	}

	machine := &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      obj.Status.UUID,
			Namespace: r.Namespace,
		},
	}
	err := r.Client.Get(ctx, client.ObjectKeyFromObject(machine), machine)
	if err != nil {
		r.Log.Info("failed to retrieve machine object from cluster", "error", err)
		return false
	}

	// The status update below did literally nothing, since there were no changes
	// of status fields. Also, I didn't manage to find where the following fields
	// are updated in case of referenced OOB object deletion. Hence the status fields
	// update was added. Apart from that, predicates should not contain any calls
	// which updates either object - the source of the event, or any other related
	// objects. In current case update logic should be moved to finalizer
	// Artem Bortnikov <artem.bortnikov@t-systems.com>
	machine.Status.OOB.Reference = nil
	machine.Status.OOB.Exist = false

	if updErr := r.Client.Status().Update(ctx, machine); updErr != nil {
		r.Log.Info("can't update machine status for oob", "error", updErr)
		return false
	}
	return false
}

func isUUIDExist(e event.CreateEvent) bool {
	obj, ok := e.Object.(*oobv1.OOB)
	if !ok {
		return false
	}
	return obj.Status.UUID != ""
}

func onUpdate(e event.UpdateEvent) bool {
	obj, ok := e.ObjectNew.(*oobv1.OOB)
	if !ok {
		return false
	}
	return obj.Status.UUID != ""
}

func updateTaints(oobObj *oobv1.OOB, machineObj *machinev1alpha2.Machine) {
	if v, ok := oobObj.Labels[maintainedMachineLabel]; ok && v == "true" {
		if getNoScheduleTaintIdx(machineObj.Spec.Taints) == -1 {
			machineObj.Spec.Taints = append(machineObj.Spec.Taints, machinev1alpha2.Taint{
				Effect: machinev1alpha2.TaintEffectNoSchedule,
				Key:    machinev1alpha2.UnschedulableLabel,
			})
		}
	} else {
		if idx := getNoScheduleTaintIdx(machineObj.Spec.Taints); idx >= 0 {
			machineObj.Spec.Taints = append(machineObj.Spec.Taints[:idx], machineObj.Spec.Taints[idx+1:]...)
		}
	}
}

func getNoScheduleTaintIdx(taints []machinev1alpha2.Taint) int {
	for t := range taints {
		if taints[t].Effect != machinev1alpha2.TaintEffectNoSchedule {
			continue
		}
		return t
	}
	return -1
}

func syncStatusState(oobObj *oobv1.OOB, machineObj *machinev1alpha2.Machine) {
	switch {
	case oobObj.Status.OS == "TimedOut":
		machineObj.Status.Reservation.Status = scheduler.ReservationStatusError
	case oobObj.Status.OS == "Ok" && oobObj.Status.Power != "Off":
		machineObj.Status.Reservation.Status = scheduler.ReservationStatusRunning
	}

	// if machine has no reservation reference and power state is off then the machine is Available
	if machineObj.Status.Reservation.Reference == nil && oobObj.Status.Power == "Off" {
		machineObj.Status.Reservation.Status = scheduler.ReservationStatusAvailable
	}
}
