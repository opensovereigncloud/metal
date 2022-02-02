/*
Copyright 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
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
	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

const CMachineType = "Machine"

// OnboardingReconciler reconciles a Switch object
type OnboardingReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *OnboardingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&inventoriesv1alpha1.Inventory{}).
		Complete(r)
}

//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch

func (r *OnboardingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("inventory", req.NamespacedName)

	invObj := &inventoriesv1alpha1.Inventory{}
	if err := r.Get(ctx, req.NamespacedName, invObj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if invObj.Spec.Host.Type == CMachineType {
		return ctrl.Result{}, nil
	}
	sw := &switchv1alpha1.Switch{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: switchv1alpha1.CNamespace, Name: invObj.Name}, sw); err != nil {
		if apierrors.IsNotFound(err) {
			// sw.SwitchFromInventory(invObj)
			if err = r.Client.Create(ctx, sw); err != nil {
				if apierrors.IsAlreadyExists(err) {
					return ctrl.Result{}, nil
				}
				r.Log.Error(err, "failed to create switch resource", "name", sw.NamespacedName())
				return ctrl.Result{}, err
			}
		} else {
			log.Error(err, "failed to lookup for switch resource")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}
