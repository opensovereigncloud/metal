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

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	oobv1 "github.com/onmetal/oob-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MachinePowerReconciler reconciles a MachineReservation object.
type MachinePowerReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/finalizers,verbs=update
//+kubebuilder:rbac:groups=onmetal.de,resources=oobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=onmetal.de,resources=oobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=onmetal.de,resources=oobs/finalizers,verbs=update

func (r *MachinePowerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("namespace", req.NamespacedName)

	machine := &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(machine), machine); err != nil {
		log.Error(err, "could not get machine")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	oob := &oobv1.OOB{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(oob), oob); err != nil {
		log.Error(err, "could not get OOB")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if machine.Status.Reservation.Reference == nil {
		log.Info("machine has no reservation, turn off OOB if it's turned on")

		if oob.Spec.Power == "Off" {
			return ctrl.Result{}, nil
		}

		oob.Spec.Power = "Off"
		if err := r.Update(ctx, oob); err != nil {
			log.Error(err, "could not turn off OOB")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	oob.Spec.Power = "On"
	if err := r.Update(ctx, oob); err != nil {
		log.Error(err, "could not turn on OOB")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *MachinePowerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&machinev1alpha2.Machine{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *MachinePowerReconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: r.handleMachineDeletion,
	}
}

func (r *MachinePowerReconciler) handleMachineDeletion(e event.DeleteEvent) bool {
	ctx := context.Background()
	machine, ok := e.Object.(*machinev1alpha2.Machine)
	if !ok {
		r.Log.Info("machine cast failed")
		return false
	}

	oob := &oobv1.OOB{
		ObjectMeta: metav1.ObjectMeta{
			Name:      machine.Name,
			Namespace: machine.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(oob), oob); err != nil {
		r.Log.Error(err, "could not get OOB")
		return false
	}

	oob.Spec.Power = "Off"
	if err := r.Update(ctx, oob); err != nil {
		r.Log.Error(err, "could not turn off OOB")
		return false
	}
	return false
}
