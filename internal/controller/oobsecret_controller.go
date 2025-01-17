// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	metalv1alpha1 "github.com/ironcore-dev/metal/api/v1alpha1"
)

// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=oobsecrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=oobsecrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metal.ironcore.dev,resources=oobsecrets/finalizers,verbs=update

func NewOOBSecretReconciler() (*OOBSecretReconciler, error) {
	return &OOBSecretReconciler{}, nil
}

// OOBSecretReconciler reconciles a OOBSecret object
type OOBSecretReconciler struct {
	client.Client
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *OOBSecretReconciler) Reconcile(ctx context.Context, _ ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OOBSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Client = mgr.GetClient()

	return ctrl.NewControllerManagedBy(mgr).
		For(&metalv1alpha1.OOBSecret{}).
		Complete(r)
}
