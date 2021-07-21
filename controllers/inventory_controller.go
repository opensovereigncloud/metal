/*
Copyright 2021.

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

	"github.com/d4l3k/messagediff"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	machinev1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
)

// InventoryReconciler reconciles a Inventory object
type InventoryReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=inventories/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Inventory object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *InventoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("inventory", req.NamespacedName)

	inv := &machinev1alpha1.Inventory{}

	err := r.Get(ctx, req.NamespacedName, inv)
	if apierrors.IsNotFound(err) {
		log.Error(err, "requested inventory resource not found", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if err != nil {
		log.Error(err, "unable to get inventory resource", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	continueToken := ""
	limit := int64(1000)

	for {
		sizeList := &machinev1alpha1.SizeList{}
		opts := &client.ListOptions{
			Namespace: req.Namespace,
			Limit:     limit,
			Continue:  continueToken,
		}
		err := r.List(ctx, sizeList, opts)
		if err != nil {
			log.Error(err, "unable to get size resource list", "namespace", req.Namespace)
			return ctrl.Result{}, err
		}

		for _, size := range sizeList.Items {
			labelName := size.GetMatchLabel()
			matches, err := size.Matches(inv)
			sizeNamespacedName := types.NamespacedName{
				Namespace: size.Namespace,
				Name:      size.Name,
			}
			if err != nil {
				log.Error(err, "unable to check match to provided size; will remove match if present", "size", sizeNamespacedName)
			}
			if matches {
				log.Info("match between inventory and size found", "inventory", req.NamespacedName, "size", sizeNamespacedName)
				inv.Labels[labelName] = "true"
			} else {
				if _, ok := inv.Labels[labelName]; ok {
					log.Info("inventory no longer matches to size", "inventory", req.NamespacedName, "size", sizeNamespacedName)
					delete(inv.Labels, labelName)
				}
			}
		}

		if sizeList.Continue == "" ||
			sizeList.RemainingItemCount == nil ||
			*sizeList.RemainingItemCount == 0 {
			break
		}

		continueToken = sizeList.Continue
	}

	continueToken = ""

	for {
		aggregateList := &machinev1alpha1.AggregateList{}
		opts := &client.ListOptions{
			Namespace: req.Namespace,
			Limit:     limit,
			Continue:  continueToken,
		}

		err := r.List(ctx, aggregateList, opts)
		if err != nil {
			log.Error(err, "unable to get aggregate resource list", "namespace", req.Namespace)
			return ctrl.Result{}, err
		}

		for _, aggregate := range aggregateList.Items {
			aggregatedValues := aggregate.Compute(inv)
			inv.Status.Computed[aggregate.Name] = aggregatedValues
		}

		if aggregateList.Continue == "" ||
			aggregateList.RemainingItemCount == nil ||
			*aggregateList.RemainingItemCount == 0 {
			break
		}

		continueToken = aggregateList.Continue
	}

	err = r.Update(ctx, inv)
	if err != nil {
		log.Error(err, "unable to update inventory resource", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InventoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&machinev1alpha1.Inventory{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *InventoryReconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: r.printDiffOnUpdate,
	}
}

func (r *InventoryReconciler) printDiffOnUpdate(event event.UpdateEvent) bool {
	old := event.ObjectOld.(*machinev1alpha1.Inventory)
	upd := event.ObjectNew.(*machinev1alpha1.Inventory)

	nsName := types.NamespacedName{
		Namespace: old.Namespace,
		Name:      old.Name,
	}

	l := r.Log.WithValues("inventory", nsName)

	msg, eq := messagediff.PrettyDiff(old.Spec, upd.Spec)
	if eq {
		l.Info("new version is the same")
		return false
	}

	l.Info("found a difference on update", "diff", msg)
	return true
}
