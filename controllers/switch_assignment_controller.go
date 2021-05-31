package controllers

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

type SwitchAssignmentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switchassignments,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=list;update
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/status,verbs=update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *SwitchAssignmentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("switchAssignment", req.NamespacedName)
	assignmentRes := &switchv1alpha1.SwitchAssignment{}
	err := r.Get(ctx, req.NamespacedName, assignmentRes)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Error(err, "requested switch assignment resource not found", "name", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to get switchAssignment resource", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	// if webhooks are not configured need to set labels in controller
	if assignmentRes.Labels == nil {
		assignmentRes.Labels = map[string]string{}
		assignmentRes.Labels[switchv1alpha1.LabelSerial] = assignmentRes.Spec.Serial
		assignmentRes.Labels[switchv1alpha1.LabelChassisId] = strings.ReplaceAll(assignmentRes.Spec.ChassisID, ":", "-")
		if err := r.Update(ctx, assignmentRes); err != nil {
			log.Error(err, "unable to set labels for switchAssignment resource", "name", req.NamespacedName)
			return ctrl.Result{}, err
		}
	}

	// find and update dependent switch
	selector := labels.SelectorFromSet(labels.Set{switchv1alpha1.LabelChassisId: strings.ReplaceAll(assignmentRes.Spec.ChassisID, ":", "-")})
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         1000,
	}
	switchesList := &switchv1alpha1.SwitchList{}
	if err := r.List(ctx, switchesList, opts); err != nil {
		log.Error(err, "unable to get switches list")
	}
	if len(switchesList.Items) == 0 {
		return ctrl.Result{RequeueAfter: switchv1alpha1.CRequeueInterval}, nil
	} else {
		targetSwitch := &switchesList.Items[0]
		targetSwitch.Spec.Role = switchv1alpha1.CSpineRole
		targetSwitch.Spec.ConnectionLevel = 0
		if err := r.Update(ctx, targetSwitch); err != nil {
			log.Error(err, "unable to update switch resource", "name", types.NamespacedName{
				Namespace: targetSwitch.Namespace,
				Name:      targetSwitch.Name,
			})
			return ctrl.Result{}, err
		}
		switchConn := &switchv1alpha1.SwitchConnection{}
		err := r.Client.Get(ctx, types.NamespacedName{
			Namespace: targetSwitch.Namespace,
			Name:      targetSwitch.Name,
		}, switchConn)
		if err != nil {
			if !apierrors.IsNotFound(err) {
				log.Error(err, "unable to get switchConnection resource")
				return ctrl.Result{}, err
			}
		} else {
			switchConn.Spec.ConnectionLevel = 0
			if err := r.Update(ctx, switchConn); err != nil {
				log.Error(err, "unable to update switchConnection resource", "name", types.NamespacedName{
					Namespace: targetSwitch.Namespace,
					Name:      targetSwitch.Name,
				})
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchAssignmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1alpha1.SwitchAssignment{}).
		Complete(r)
}
