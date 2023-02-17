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

package v1beta1

import (
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/internal/constants"
)

// OnboardingReconciler reconciles Switch object corresponding
// to given Inventory object.
type OnboardingReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch

func (r *OnboardingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	nestedCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	inventory := &inventoryv1alpha1.Inventory{}
	if err := r.Get(nestedCtx, req.NamespacedName, inventory); err != nil {
		switch {
		case apierrors.IsNotFound(err):
			r.Log.Info("requested Inventory object not found", "name", req.NamespacedName)
		default:
			r.Log.Error(err, "failed to get requested Inventory object", "name", req.NamespacedName)
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	return r.reconcile(nestedCtx, inventory)
}

// SetupWithManager sets up the controller with the Manager.
func (r *OnboardingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	log := mgr.GetLogger().WithName("onboarding-controller-setup")

	// setting up the label selector predicate to filter inventories related to switches
	labelSelectorPredicate, err := predicate.LabelSelectorPredicate(metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      constants.SizeLabel,
				Operator: metav1.LabelSelectorOpExists,
				Values:   []string{},
			},
		},
	})
	if err != nil {
		log.Error(err, "failed to setup predicates")
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&inventoryv1alpha1.Inventory{}).
		WithEventFilter(predicate.And(labelSelectorPredicate)).
		Complete(r)
}

func (r *OnboardingReconciler) reconcile(ctx context.Context, inv *inventoryv1alpha1.Inventory) (ctrl.Result, error) {
	if !inv.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, nil
	}
	return r.onboarding(ctx, inv)
}

func (r *OnboardingReconciler) onboarding(ctx context.Context, inv *inventoryv1alpha1.Inventory) (ctrl.Result, error) {
	var err error
	key := client.ObjectKeyFromObject(inv)
	targetSwitch := &switchv1beta1.Switch{}
	if err = r.Get(ctx, key, targetSwitch); err != nil {
		if !apierrors.IsNotFound(err) {
			r.Log.Error(err, "failed to get Switch object", "name", key)
			return ctrl.Result{}, err
		}
		targetSwitch = &switchv1beta1.Switch{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: switchv1beta1.SwitchSpec{
				Managed:   pointer.Bool(true),
				Cordon:    pointer.Bool(false),
				TopSpine:  pointer.Bool(false),
				ScanPorts: pointer.Bool(true),
			},
		}
	}
	targetSwitch.SetInventoryRef(key.Name)
	if targetSwitch.Labels == nil {
		targetSwitch.Labels = map[string]string{}
	}
	targetSwitch.UpdateSwitchAnnotations(inv)
	targetSwitch.UpdateSwitchLabels(inv)
	targetSwitch.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
	targetSwitch.SetCondition(constants.ConditionInterfacesOK, false).SetReason(constants.ReasonSourceUpdated)
	err = r.Patch(ctx, targetSwitch, client.Apply, patchOpts)
	return ctrl.Result{}, err
}
