package controllers

import (
	"context"

	"github.com/go-logr/logr"
	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	"github.com/onmetal/metal-api/internal/entity"
	"github.com/onmetal/metal-api/internal/usecase"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OnboardingReconciler struct {
	client.Client

	Log                  logr.Logger
	Scheme               *runtime.Scheme
	OnboardingRepo       usecase.Onboarding
	DestinationNamespace string
}

// SetupWithManager sets up the controller with the Manager.
func (r *OnboardingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&inventoriesv1alpha1.Inventory{}).
		Complete(r)
}

func (r *OnboardingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("namespace", req.NamespacedName)

	e := entity.Onboarding{
		RequestName:                   req.Name,
		RequestNamespace:              req.Namespace,
		InitializationObjectNamespace: r.DestinationNamespace}

	is := r.OnboardingRepo.InitializationStatus(ctx, e)
	if is.Error != nil {
		reqLogger.Info("device initialization status", "error", is.Error)
		return ctrl.Result{}, nil
	}
	if is.Require {
		if err := r.OnboardingRepo.Initiate(ctx, e); err != nil {
			reqLogger.Info("device initialization failed", "error", err)
			return ctrl.Result{}, nil
		}
	}

	if err := r.OnboardingRepo.GatherData(ctx, e); err != nil {
		if apierrors.IsConflict(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		reqLogger.Info("can't gather the information", "error", err)
	}

	reqLogger.Info("reconciliation finished")
	return ctrl.Result{}, nil
}
