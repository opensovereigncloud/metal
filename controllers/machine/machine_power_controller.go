// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	oobv1 "github.com/onmetal/oob-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
)

// MachinePowerReconciler reconciles a MachineReservation object.
type MachinePowerReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=machines/finalizers,verbs=update
// +kubebuilder:rbac:groups=onmetal.de,resources=oobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=onmetal.de,resources=oobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=onmetal.de,resources=oobs/finalizers,verbs=update

func (r *MachinePowerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("namespace", req.NamespacedName)

	machineObj := &metalv1alpha4.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(machineObj), machineObj); err != nil {
		log.Error(err, "could not get metalv1alpha4")
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

	if machineObj.Status.Reservation.Reference == nil {
		log.Info("metalv1alpha4 has no reservation, turn off OOB if it's turned on")

		if oob.Spec.Power == "Off" {
			return ctrl.Result{}, nil
		}

		oob.Spec.Power = "Off"
		log.Info("OOB is turned off")

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
		For(&metalv1alpha4.Machine{}).
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
	machineObj, ok := e.Object.(*metalv1alpha4.Machine)
	if !ok {
		r.Log.Info("metalv1alpha4 cast failed")
		return false
	}

	oob := &oobv1.OOB{
		ObjectMeta: metav1.ObjectMeta{
			Name:      machineObj.Name,
			Namespace: machineObj.Namespace,
		},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(oob), oob); err != nil {
		r.Log.Error(err, "could not get OOB")
		return false
	}

	oob.Spec.Power = "Off"
	r.Log.Info("OOB is turned off")
	return false
}
