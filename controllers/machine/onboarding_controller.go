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
	machinev1lpha1 "github.com/onmetal/metal-api/apis/machine/v1alpha1"
	oobv1 "github.com/onmetal/oob-controller/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const UUIDLabel = "machine.onmetal.de/uuid"

// MachineReconciler reconciles a Machine object.
type OnboardingReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *OnboardingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&oobv1.Machine{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *OnboardingReconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: r.onDelete,
	}
}

//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/finalizers,verbs=update
//+kubebuilder:rbac:groups=oob.onmetal.de,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=oob.onmetal.de,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=oob.onmetal.de,resources=machines/finalizers,verbs=update

func (r *OnboardingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("machine", req.NamespacedName)

	oobObj := &oobv1.Machine{}
	if err := r.Get(ctx, req.NamespacedName, oobObj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	m := &machinev1lpha1.Machine{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: oobObj.Spec.UUID}, m); err != nil {
		if apierrors.IsNotFound(err) {
			if err := r.createAndEnableMachine(ctx, oobObj); err != nil {
				if apierrors.IsAlreadyExists(err) {
					return ctrl.Result{}, nil
				}
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
	}
	if !m.Status.OOB {
		m.Status.OOB = true
		if statusUpdErr := r.Client.Status().Update(ctx, m); statusUpdErr != nil {
			return ctrl.Result{}, statusUpdErr
		}
	}
	reqLogger.Info("reconciliation finished")
	return ctrl.Result{}, nil
}

func (r *OnboardingReconciler) createAndEnableMachine(ctx context.Context, oobObj *oobv1.Machine) error {
	obj := prepareMachine(oobObj)
	if err := r.Client.Create(ctx, obj); err != nil {
		return err
	}
	oobObj.Spec.PowerState = "On"
	return r.Client.Update(ctx, oobObj)
}

func prepareMachine(oob *oobv1.Machine) *machinev1lpha1.Machine {
	return &machinev1lpha1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      oob.Spec.UUID,
			Namespace: oob.Namespace,
			Labels:    map[string]string{UUIDLabel: oob.Spec.UUID},
		},
		Spec: machinev1lpha1.MachineSpec{
			Action: machinev1lpha1.Action{PowerState: ""}, InventoryRequested: true,
		},
	}

}

func (r *OnboardingReconciler) onDelete(e event.DeleteEvent) bool {
	obj, ok := e.Object.(*oobv1.Machine)
	if !ok {
		r.Log.Info("machine oob type assertion failed")
		return false
	}
	ctx := context.Background()
	m := &machinev1lpha1.Machine{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: obj.Namespace, Name: obj.Spec.UUID}, m); err != nil {
		return false
	}
	m.Status.OOB = false
	if statusUpdErr := r.Client.Status().Update(ctx, m); statusUpdErr != nil {
		return false
	}
	return false
}
