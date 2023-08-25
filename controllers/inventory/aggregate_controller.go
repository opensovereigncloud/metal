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

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	inventories "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
)

const (
	CAggregateFinalizer = "aggregate.machine.onmetal.de/finalizer"
)

// AggregateReconciler reconciles a Aggregate object.
type AggregateReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=machine.onmetal.de,resources=aggregates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=aggregates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=aggregates/finalizers,verbs=update

func (r *AggregateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("aggregate", req.NamespacedName)

	aggregate := &inventories.Aggregate{}
	err := r.Get(ctx, req.NamespacedName, aggregate)
	if apierrors.IsNotFound(err) {
		log.Error(err, "requested aggregate resource not found", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if err != nil {
		log.Error(err, "unable to get aggregate resource", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	if aggregate.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(aggregate, CAggregateFinalizer) {
			if err := r.finalizeAggregate(ctx, req, log, aggregate); err != nil {
				log.Error(err, "unable to finalize aggregate resource", "name", req.NamespacedName)
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(aggregate, CAggregateFinalizer)
			err := r.Update(ctx, aggregate)
			if err != nil {
				log.Error(err, "unable to update aggregate resource on finalizer removal", "name", req.NamespacedName)
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(aggregate, CAggregateFinalizer) {
		controllerutil.AddFinalizer(aggregate, CAggregateFinalizer)
		err = r.Update(ctx, aggregate)
		if err != nil {
			log.Error(err, "unable to update aggregate resource with finalizer", "name", req.NamespacedName)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	continueToken := ""
	for {
		inventoryList := &inventories.InventoryList{}
		opts := &client.ListOptions{
			Namespace: req.Namespace,
			Limit:     CPageLimit,
			Continue:  continueToken,
		}
		err := r.List(ctx, inventoryList, opts)
		if err != nil {
			log.Error(err, "unable to get inventory resource list", "namespace", req.Namespace)
			return ctrl.Result{}, err
		}

		for _, inventory := range inventoryList.Items {
			if inventory.GetDeletionTimestamp() != nil {
				continue
			}

			inventoryNamespacedName := types.NamespacedName{
				Namespace: inventory.Namespace,
				Name:      inventory.Name,
			}

			aggregatedValues, err := aggregate.Compute(&inventory)
			if err != nil {
				log.Error(err, "unable to compute aggregate", "inventory", inventoryNamespacedName)
				return ctrl.Result{}, err
			}
			if inventory.Status.Computed.Object == nil {
				inventory.Status.Computed.Object = make(map[string]interface{})
			}
			inventory.Status.Computed.Object[aggregate.Name] = aggregatedValues

			if err := r.Status().Update(ctx, &inventory); err != nil {
				log.Error(err, "unable to update inventory resource", "inventory", inventoryNamespacedName)
				return ctrl.Result{}, err
			}
		}

		if inventoryList.Continue == "" ||
			inventoryList.RemainingItemCount == nil ||
			*inventoryList.RemainingItemCount == 0 {
			break
		}

		continueToken = inventoryList.Continue
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AggregateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&inventories.Aggregate{}).
		Complete(r)
}

func (r *AggregateReconciler) finalizeAggregate(ctx context.Context, req ctrl.Request, log logr.Logger, aggregate *inventories.Aggregate) error {
	continueToken := ""
	aggregateKey := aggregate.Name

	for {
		inventoryList := &inventories.InventoryList{}
		opts := &client.ListOptions{
			Namespace: req.Namespace,
			Limit:     CPageLimit,
			Continue:  continueToken,
		}
		err := r.List(ctx, inventoryList, opts)
		if err != nil {
			log.Error(err, "unable to get inventory resource list", "namespace", req.Namespace)
			return err
		}

		for _, inventory := range inventoryList.Items {
			_, ok := inventory.Status.Computed.Object[aggregateKey]
			if !ok {
				continue
			}

			inventoryNamespacedName := types.NamespacedName{
				Namespace: inventory.Namespace,
				Name:      inventory.Name,
			}
			log.Info("inventory contains aggregate, removing on finalize", "aggregate", req.NamespacedName, "inventory", inventoryNamespacedName)

			delete(inventory.Status.Computed.Object, aggregateKey)

			if err := r.Status().Update(ctx, &inventory); err != nil {
				log.Error(err, "unable to update inventory resource", "inventory", inventoryNamespacedName)
				return err
			}
		}

		if inventoryList.Continue == "" ||
			inventoryList.RemainingItemCount == nil ||
			*inventoryList.RemainingItemCount == 0 {
			break
		}

		continueToken = inventoryList.Continue
	}

	return nil
}
