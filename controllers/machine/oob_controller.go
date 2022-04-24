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
	"github.com/onmetal/metal-api/pkg/machine"
	"github.com/onmetal/metal-api/pkg/provider"
	oobv1 "github.com/onmetal/oob-controller/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const maintainedMachineLabel = "onmetal.de/oob-ignore"

// MachineReconciler reconciles a Machine object.
type OOBReconciler struct {
	client.Client

	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// SetupWithManager sets up the controller with the Manager.
func (r *OOBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&oobv1.Machine{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *OOBReconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: r.onDelete,
	}
}

//+kubebuilder:rbac:groups=oob.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=oob.onmetal.de,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=oob.onmetal.de,resources=machines/finalizers,verbs=update

func (r *OOBReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("machine-oob", req.NamespacedName)

	oobObj := &oobv1.Machine{}
	if err := provider.GetObject(ctx, req.Name, req.Namespace, r.Client, oobObj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	machineObj := &machinev1alpha2.Machine{}
	if err := provider.GetObject(ctx, oobObj.Spec.UUID, req.Namespace, r.Client, machineObj); err != nil {
		if apierrors.IsNotFound(err) {
			if err := r.createAndEnableMachine(ctx, oobObj); err != nil {
				if apierrors.IsAlreadyExists(err) {
					return ctrl.Result{}, nil
				}
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, err
	}

	mm := machine.New(ctx, r.Client, r.Log, r.Recorder)
	if _, ok := oobObj.Labels[machinev1alpha2.UUIDLabel]; !ok {
		oobObj.Labels = setUpLabels(oobObj)
		if err := r.Client.Update(ctx, oobObj); err != nil {
			return ctrl.Result{}, err
		}
	}

	updateTaints(oobObj, machineObj)

	if specUpdErr := mm.UpdateSpec(machineObj); specUpdErr != nil {
		return ctrl.Result{}, specUpdErr
	}

	if !machineObj.Status.OOB.Exist {
		machineObj.Status.OOB = prepareRefenceSpec(oobObj)
	}

	updateStatusState(oobObj, machineObj)

	if statusUpdErr := mm.UpdateStatus(machineObj); statusUpdErr != nil {
		return ctrl.Result{}, statusUpdErr
	}

	reqLogger.Info("reconciliation finished")
	return ctrl.Result{}, nil
}

func (r *OOBReconciler) createAndEnableMachine(ctx context.Context, oobObj *oobv1.Machine) error {
	obj := prepareMachine(oobObj)
	if err := r.Client.Create(ctx, obj); err != nil {
		return err
	}

	oobObj.Spec.PowerState = getPowerState(oobObj.Spec.PowerState)
	oobObj.Labels = setUpLabels(oobObj)
	return r.Client.Update(ctx, oobObj)
}

func prepareMachine(oob *oobv1.Machine) *machinev1alpha2.Machine {
	return &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      oob.Spec.UUID,
			Namespace: oob.Namespace,
			Labels:    map[string]string{machinev1alpha2.UUIDLabel: oob.Spec.UUID},
		},
		Spec: machinev1alpha2.MachineSpec{InventoryRequested: true},
	}

}

func prepareRefenceSpec(oob *oobv1.Machine) machinev1alpha2.ObjectReference {
	return machinev1alpha2.ObjectReference{
		Exist: true,
		Reference: &machinev1alpha2.ResourceReference{
			Kind: oob.Kind, APIVersion: oob.APIVersion,
			Name: oob.Name, Namespace: oob.Namespace},
	}
}

func (r *OOBReconciler) onDelete(e event.DeleteEvent) bool {
	ctx := context.Background()
	obj, ok := e.Object.(*oobv1.Machine)
	if !ok {
		r.Log.Info("machine oob type assertion failed")
		return false
	}

	machineObj := &machinev1alpha2.Machine{}
	if err := provider.GetObject(ctx, obj.Name, obj.Namespace, r.Client, machineObj); err != nil {
		r.Log.Info("failed to retrieve machine object from cluster", "error", err)
		return false
	}

	machineObj.Status.OOB.Exist = false
	machineObj.Status.OOB.Reference = nil

	mm := machine.New(ctx, r.Client, r.Log, r.Recorder)
	if updErr := mm.UpdateStatus(machineObj); updErr != nil {
		r.Log.Info("can't update machine status for inventory", "error", updErr)
		return false
	}
	return false
}

func getPowerState(state string) string {
	switch state {
	case "On":
		// In case when machine already running Reset is required.
		// Machine should be started from scratch.
		return "Reset"
	default:
		return "On"
	}
}

func setUpLabels(oobObj *oobv1.Machine) map[string]string {
	if oobObj.Labels == nil {
		return map[string]string{machinev1alpha2.UUIDLabel: oobObj.Spec.UUID}
	}
	if _, ok := oobObj.Labels[machinev1alpha2.UUIDLabel]; !ok {
		oobObj.Labels[machinev1alpha2.UUIDLabel] = oobObj.Spec.UUID
	}
	return oobObj.Labels
}

func updateTaints(oobObj *oobv1.Machine, machineObj *machinev1alpha2.Machine) {
	if len(machineObj.Spec.Taints) == 0 {
		return
	}
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

func updateStatusState(oobObj *oobv1.Machine, machineObj *machinev1alpha2.Machine) {
	switch {
	case oobObj.Status.SystemStateReadTimeout:
		machineObj.Status.RequestState = machinev1alpha2.RequestStateError
	case oobObj.Status.SystemState == "Ok" || oobObj.Status.SystemState == "Unknown":
		machineObj.Status.RequestState = machinev1alpha2.RequestStateRunning
	default:
		machineObj.Status.RequestState = machinev1alpha2.RequestStatePending
	}
}
