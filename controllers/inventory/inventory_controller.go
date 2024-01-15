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
	"strings"

	"github.com/d4l3k/messagediff"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
)

const (
	CMACAddressLabelPrefix = "metal.ironcore.dev/mac-address-"
	CDefaultAggregateName  = "default"
)

// InventoryReconciler reconciles an Inventory object.
type InventoryReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=metalv1alpha4,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=metalv1alpha4/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=metalv1alpha4/finalizers,verbs=update

func (r *InventoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("inventory", req.NamespacedName)

	inv := &metalv1alpha4.Inventory{}

	err := r.Get(ctx, req.NamespacedName, inv)
	if apierrors.IsNotFound(err) {
		log.Error(err, "requested inventory resource not found", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if err != nil {
		log.Error(err, "unable to get inventory resource", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	if inv.GetDeletionTimestamp() != nil {
		return ctrl.Result{}, nil
	}

	defaultAggregateNamespacedName := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      CDefaultAggregateName,
	}
	defaultAggregate := &metalv1alpha4.Aggregate{}
	err = r.Get(ctx, defaultAggregateNamespacedName, defaultAggregate)
	if apierrors.IsNotFound(err) {
		log.Info("trying to create default aggregate", "name", defaultAggregateNamespacedName)
		if err := r.createDefaultAggregate(ctx, req.Namespace); err != nil {
			log.Error(err, "unable to create default aggregate", "agg", defaultAggregateNamespacedName)
			return ctrl.Result{}, err
		}
		// Assuming that aggregate will be calculated and inventory will be reconciled on change
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, "unable to check whether default aggregate exists", "agg", defaultAggregateNamespacedName)
		return ctrl.Result{}, err
	}

	if inv.Status.Computed.Object == nil {
		inv.Status.Computed.Object = make(map[string]interface{})
	}

	continueToken := ""
	limit := int64(1000)

	for {
		aggregateList := &metalv1alpha4.AggregateList{}
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
			if aggregate.GetDeletionTimestamp() != nil {
				continue
			}

			aggregatedValues, err := aggregate.Compute(inv)
			if err != nil {
				log.Error(err, "unable to compute aggregate", "inventory", req.NamespacedName)
			}
			inv.Status.Computed.Object[aggregate.Name] = aggregatedValues
		}

		if aggregateList.Continue == "" ||
			aggregateList.RemainingItemCount == nil ||
			*aggregateList.RemainingItemCount == 0 {
			break
		}

		continueToken = aggregateList.Continue
	}

	if err = r.Status().Update(ctx, inv); err != nil {
		log.Error(err, "unable to update inventory status resource", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	if inv.Labels == nil {
		inv.Labels = make(map[string]string)
	}

	// Due to k8s validation which allows labels to consist of alphanumeric characters, '-', '_' or '.' need to replace
	// colons in nic's MAC address
	for _, nic := range inv.Spec.NICs {
		inv.Labels[CMACAddressLabelPrefix+strings.ReplaceAll(nic.MACAddress, ":", "")] = ""
	}

	continueToken = ""

	for {
		sizeList := &metalv1alpha4.SizeList{}
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
			if size.GetDeletionTimestamp() != nil {
				continue
			}

			labelName := size.GetMatchLabel()
			matches, err := size.Matches(inv)
			sizeNamespacedName := types.NamespacedName{
				Namespace: size.Namespace,
				Name:      size.Name,
			}
			if err != nil {
				log.Error(err, "unable to check match to provided size; will remove match if present", "size", sizeNamespacedName)
			}
			// for some reason tests panicked due to assignment to nil map
			// despite the check above
			if inv.Labels == nil {
				inv.Labels = map[string]string{}
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

	if err = r.Update(ctx, inv); err != nil {
		log.Error(err, "unable to update inventory resource", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InventoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metalv1alpha4.Inventory{}).
		WithEventFilter(r.constructPredicates()).
		Complete(r)
}

func (r *InventoryReconciler) constructPredicates() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: r.printDiffOnUpdate,
	}
}

func (r *InventoryReconciler) printDiffOnUpdate(event event.UpdateEvent) bool {
	old, oldOK := event.ObjectOld.(*metalv1alpha4.Inventory)
	upd, updOK := event.ObjectNew.(*metalv1alpha4.Inventory)
	if !oldOK || !updOK {
		r.Log.Info("printDiffOnUpdate: type assertions failed")
		return false
	}
	nsName := types.NamespacedName{
		Namespace: old.Namespace,
		Name:      old.Name,
	}

	l := r.Log.WithValues("inventory", nsName)

	specMsg, specEq := messagediff.PrettyDiff(old.Spec, upd.Spec)
	aggMsg, aggEq := messagediff.PrettyDiff(old.Status.Computed, upd.Status.Computed)

	if specEq && aggEq {
		l.Info("new version is the same")
		return false
	}

	l.Info("found a difference on update", "spec", specMsg, "agg", aggMsg)
	return true
}

func (r *InventoryReconciler) createDefaultAggregate(ctx context.Context, namespace string) error {
	agg := &metalv1alpha4.Aggregate{
		ObjectMeta: ctrl.ObjectMeta{
			Namespace: namespace,
			Name:      CDefaultAggregateName,
		},
		Spec: metalv1alpha4.AggregateSpec{
			Aggregates: []metalv1alpha4.AggregateItem{
				{
					SourcePath: *metalv1alpha4.JSONPathFromString("spec.blocks[*]"),
					TargetPath: *metalv1alpha4.JSONPathFromString("blocks.count"),
					Aggregate:  metalv1alpha4.CCountAggregateType,
				},
				{
					SourcePath: *metalv1alpha4.JSONPathFromString("spec.blocks[*].size"),
					TargetPath: *metalv1alpha4.JSONPathFromString("blocks.capacity"),
					Aggregate:  metalv1alpha4.CSumAggregateType,
				},
				{
					SourcePath: *metalv1alpha4.JSONPathFromString("spec.cpus[*]"),
					TargetPath: *metalv1alpha4.JSONPathFromString("cpus.sockets"),
					Aggregate:  metalv1alpha4.CCountAggregateType,
				},
				{
					SourcePath: *metalv1alpha4.JSONPathFromString("spec.cpus[*].cores"),
					TargetPath: *metalv1alpha4.JSONPathFromString("cpus.cores"),
					Aggregate:  metalv1alpha4.CSumAggregateType,
				},
				{
					SourcePath: *metalv1alpha4.JSONPathFromString("spec.cpus[*].siblings"),
					TargetPath: *metalv1alpha4.JSONPathFromString("cpus.threads"),
					Aggregate:  metalv1alpha4.CSumAggregateType,
				},
				{
					SourcePath: *metalv1alpha4.JSONPathFromString("spec.nics[*]"),
					TargetPath: *metalv1alpha4.JSONPathFromString("nics.count"),
					Aggregate:  metalv1alpha4.CCountAggregateType,
				},
			},
		},
	}

	if err := r.Create(ctx, agg); err != nil {
		return err
	}

	return nil
}
