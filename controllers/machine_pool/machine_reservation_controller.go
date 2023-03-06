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
	computev1alpha1 "github.com/onmetal/onmetal-api/api/compute/v1alpha1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MachineReservationReconciler reconciles a MachineReservation object.
type MachineReservationReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/finalizers,verbs=update
//+kubebuilder:rbac:groups=compute.api.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=compute.api.onmetal.de,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=compute.api.onmetal.de,resources=machines/finalizers,verbs=update

func (r *MachineReservationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("namespace", req.NamespacedName)

	computeMachine := &computev1alpha1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(computeMachine), computeMachine); err != nil {
		log.Error(err, "could not get compute machine")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if computeMachine.Spec.MachinePoolRef == nil {
		log.Info("compute machine has empty machine pool ref. skip reservation update")
		return ctrl.Result{}, nil
	}

	metalMachine := &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(metalMachine), metalMachine); err != nil {
		log.Error(err, "could not get metal machine")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if metalMachine.Status.Health != machinev1alpha2.MachineStateHealthy {
		log.Info("could not update reservation. metal machine is unhealthy")
		return ctrl.Result{}, nil
	}

	metalMachine.Status.Reservation.Reference = &machinev1alpha2.ResourceReference{
		APIVersion: computeMachine.APIVersion,
		Kind:       computeMachine.Kind,
		Name:       computeMachine.Name,
		Namespace:  computeMachine.Namespace,
	}

	if err := r.Client.Status().Update(ctx, metalMachine); err != nil {
		log.Error(err, "could not update metal machine status")
		return ctrl.Result{}, errors.Wrap(err, "failed to update machine status")
	}

	return ctrl.Result{}, nil
}

func (r *MachineReservationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&computev1alpha1.Machine{}).
		Complete(r)
}
