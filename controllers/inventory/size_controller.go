// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/ironcore-dev/metal/common/types/events"
	domain "github.com/ironcore-dev/metal/domain/inventory"
)

const (
	CSizeFinalizer = "size.metal.ironcore.dev/finalizer"
	CPageLimit     = 1000
)

// SizeReconciler reconciles a Size object.
type SizeReconciler struct {
	client.Client

	Log            logr.Logger
	Scheme         *runtime.Scheme
	EventPublisher events.DomainEventPublisher
}

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=sizes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=sizes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=sizes/finalizers,verbs=update

func (r *SizeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("size", req.NamespacedName)

	size := &metalv1alpha4.Size{}
	err := r.Get(ctx, req.NamespacedName, size)
	if apierrors.IsNotFound(err) {
		log.Error(err, "requested size resource not found", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if err != nil {
		log.Error(err, "unable to get size resource", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	if size.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(size, CSizeFinalizer) {
			if err := r.finalizeSize(ctx, req, log, size); err != nil {
				log.Error(err, "unable to finalize size resource", "name", req.NamespacedName)
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(size, CSizeFinalizer)
			err := r.Update(ctx, size)
			if err != nil {
				log.Error(err, "unable to update size resource on finalizer removal", "name", req.NamespacedName)
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(size, CSizeFinalizer) {
		controllerutil.AddFinalizer(size, CSizeFinalizer)
		err = r.Update(ctx, size)
		if err != nil {
			log.Error(err, "unable to update size resource with finalizer", "name", req.NamespacedName)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	continueToken := ""
	for {
		inventoryList := &metalv1alpha4.InventoryList{}
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

			matches, err := size.Matches(&inventory)
			inventoryNamespacedName := types.NamespacedName{
				Namespace: inventory.Namespace,
				Name:      inventory.Name,
			}
			if err != nil {
				log.Error(err, "unable to check match to provided size; will remove match if present", "size", req.NamespacedName, "inventory", inventoryNamespacedName)
			}

			labelName := size.GetMatchLabel()
			labelPresent := false
			if inventory.Labels != nil {
				_, labelPresent = inventory.Labels[labelName]
			}
			if matches && labelPresent {
				log.Info("match between inventory and size found, label present, will not update", "size", req.NamespacedName, "inventory", inventoryNamespacedName)
				continue
			}
			if !matches && !labelPresent {
				log.Info("match between inventory and size is not found", "size", req.NamespacedName, "inventory", inventoryNamespacedName)
				continue
			}
			if matches && !labelPresent {
				log.Info("match between inventory and size found", "size", req.NamespacedName, "inventory", inventoryNamespacedName)
				if inventory.Labels == nil {
					inventory.Labels = make(map[string]string)
				}
				inventory.Labels[labelName] = "true"
			}
			if !matches && labelPresent {
				log.Info("inventory no longer matches to size, will remove label", "size", req.NamespacedName, "inventory", inventoryNamespacedName)
				delete(inventory.Labels, labelName)
			}

			if err := r.Update(ctx, &inventory); err != nil {
				log.Error(err, "unable to update inventory resource", "inventory", inventoryNamespacedName)
				return ctrl.Result{}, err
			}
			id := domain.NewInventoryID(inventory.Labels["id"])
			r.EventPublisher.Publish(domain.NewInventoryFlavorUpdatedDomainEvent(id))
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

func (r *SizeReconciler) finalizeSize(ctx context.Context, req ctrl.Request, log logr.Logger, size *metalv1alpha4.Size) error {
	continueToken := ""
	sizeLabel := size.GetMatchLabel()

	for {
		inventoryList := &metalv1alpha4.InventoryList{}
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
			_, ok := inventory.Labels[sizeLabel]
			if !ok {
				continue
			}

			inventoryNamespacedName := types.NamespacedName{
				Namespace: inventory.Namespace,
				Name:      inventory.Name,
			}
			log.Info("inventory contains size, removing on finalize", "size", req.NamespacedName, "inventory", inventoryNamespacedName)

			delete(inventory.Labels, sizeLabel)

			if err := r.Update(ctx, &inventory); err != nil {
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

// SetupWithManager sets up the controller with the Manager.
func (r *SizeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metalv1alpha4.Size{}).
		Complete(r)
}
