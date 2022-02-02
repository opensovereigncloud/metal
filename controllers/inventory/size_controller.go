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
	"github.com/onmetal/k8s-inventory/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	CSizeFinalizer = "size.machine.onmetal.de/finalizer"
	CPageLimit     = 1000
)

// SizeReconciler reconciles a Size object
type SizeReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=machine.onmetal.de,resources=sizes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=sizes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=machine.onmetal.de,resources=sizes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Size object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *SizeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("size", req.NamespacedName)

	size := &v1alpha1.Size{}
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
		inventoryList := &v1alpha1.InventoryList{}
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

func (r *SizeReconciler) finalizeSize(ctx context.Context, req ctrl.Request, log logr.Logger, size *v1alpha1.Size) error {
	continueToken := ""
	sizeLabel := size.GetMatchLabel()

	for {
		inventoryList := &v1alpha1.InventoryList{}
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
		For(&v1alpha1.Size{}).
		Complete(r)
}
