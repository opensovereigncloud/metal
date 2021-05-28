package controllers

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
	"github.com/onmetal/switch-operator/util"
)

type SwitchConnectionReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switchconnections,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *SwitchConnectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("switchConnection", req.NamespacedName)
	switchConnection := &switchv1alpha1.SwitchConnection{}
	err := r.Get(ctx, req.NamespacedName, switchConnection)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Error(err, "requested switch connection resource not found", "name", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to get switchConnection resource", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	if switchConnection.Spec.ConnectionLevel == 0 {
		if len(switchConnection.Spec.DownstreamSwitches.Switches) > 0 {
			if err := r.updateConnectionSwitchesData(switchConnection, ctx); err != nil {
				log.Error(err, "failed to update switchConnection resource")
				return ctrl.Result{}, err
			}
		}
		if err := r.Update(ctx, switchConnection); err != nil {
			log.Error(err, "unable to update switchConnection resource")
			return ctrl.Result{}, err
		}
	} else {
		connectionRebuild, err := r.rebuildConnections(switchConnection, ctx)
		if err != nil {
			log.Error(err, "unable to update connection level", "name", req.NamespacedName)
		}

		if connectionRebuild {
			if err := r.Update(ctx, switchConnection); err != nil {
				log.Error(err, "unable to update switchConnection resource")
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{RequeueAfter: util.CRequeueInterval}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchConnectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1alpha1.SwitchConnection{}).
		Complete(r)
}

func (r *SwitchConnectionReconciler) updateConnectionSwitchesData(switchConnection *switchv1alpha1.SwitchConnection, ctx context.Context) error {
	switchesMap := map[string]switchv1alpha1.ConnectedSwitchSpec{}
	chassisIdsForLabels := make([]string, 0)
	for _, item := range switchConnection.Spec.DownstreamSwitches.Switches {
		chassisIdsForLabels = append(chassisIdsForLabels, strings.ReplaceAll(item.ChassisID, ":", "-"))
	}
	for _, item := range switchConnection.Spec.UpstreamSwitches.Switches {
		chassisIdsForLabels = append(chassisIdsForLabels, strings.ReplaceAll(item.ChassisID, ":", "-"))
	}

	labelsReq, err := labels.NewRequirement(util.LabelChassisId, selection.In, chassisIdsForLabels)
	if err != nil {
		r.Log.Error(err, "unable to build label selector requirements")
		return err
	}
	selector := labels.NewSelector()
	selector = selector.Add(*labelsReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         1000,
	}
	swList := &switchv1alpha1.SwitchList{}
	err = r.List(ctx, swList, opts)
	if err != nil {
		r.Log.Error(err, "unable to get switches list")
		return err
	}

	for _, switchRes := range swList.Items {
		switchesMap[switchRes.Spec.SwitchChassis.ChassisID] = switchv1alpha1.ConnectedSwitchSpec{
			Name:      switchRes.Name,
			Namespace: switchRes.Namespace,
			ChassisID: switchRes.Spec.SwitchChassis.ChassisID,
		}
	}
	for i, item := range switchConnection.Spec.DownstreamSwitches.Switches {
		if value, ok := switchesMap[item.ChassisID]; ok {
			switchConnection.Spec.DownstreamSwitches.Switches[i] = value
		}
	}
	for i, item := range switchConnection.Spec.UpstreamSwitches.Switches {
		if value, ok := switchesMap[item.ChassisID]; ok {
			switchConnection.Spec.UpstreamSwitches.Switches[i] = value
		}
	}

	return nil
}

func (r *SwitchConnectionReconciler) rebuildConnections(conn *switchv1alpha1.SwitchConnection, ctx context.Context) (bool, error) {
	update := false
	connList := &switchv1alpha1.SwitchConnectionList{}
	if err := r.Client.List(ctx, connList); err != nil {
		return false, err
	}
	for _, item := range connList.Items {
		for _, downstreamConn := range item.Spec.DownstreamSwitches.Switches {
			if conn.Spec.Switch.ChassisID == downstreamConn.ChassisID && item.Spec.ConnectionLevel != 255 {
				if conn.Spec.ConnectionLevel != item.Spec.ConnectionLevel+1 {
					conn.Spec.ConnectionLevel = item.Spec.ConnectionLevel + 1
					update = true
				}
				for _, value := range conn.Spec.DownstreamSwitches.Switches {
					if value.ChassisID == item.Spec.Switch.ChassisID {
						conn.Spec.DownstreamSwitches.Switches = updateDownstreamSwitches(conn, item.Spec.Switch.ChassisID)
						conn.Spec.DownstreamSwitches.Count = len(conn.Spec.DownstreamSwitches.Switches)
						conn.Spec.UpstreamSwitches.Switches = updateUpstreamSwitches(conn, value)
						conn.Spec.UpstreamSwitches.Count = len(conn.Spec.UpstreamSwitches.Switches)
						if err := r.updateConnectionSwitchesData(conn, ctx); err != nil {
							return false, err
						}
						update = true
					}
				}
			}
		}
	}
	return update, nil
}

func updateDownstreamSwitches(conn *switchv1alpha1.SwitchConnection, switchToRemove string) []switchv1alpha1.ConnectedSwitchSpec {
	newDownstreams := make([]switchv1alpha1.ConnectedSwitchSpec, 0)
	for _, item := range conn.Spec.DownstreamSwitches.Switches {
		if item.ChassisID != switchToRemove {
			newDownstreams = append(newDownstreams, item)
		}
	}
	return newDownstreams
}

func updateUpstreamSwitches(conn *switchv1alpha1.SwitchConnection, switchToAdd switchv1alpha1.ConnectedSwitchSpec) []switchv1alpha1.ConnectedSwitchSpec {
	if conn.Spec.UpstreamSwitches.Switches == nil {
		newUpstreams := make([]switchv1alpha1.ConnectedSwitchSpec, 0)
		newUpstreams = append(newUpstreams, switchToAdd)
		return newUpstreams
	} else {
		newUpstreams := make([]switchv1alpha1.ConnectedSwitchSpec, 0)
		for _, item := range conn.Spec.UpstreamSwitches.Switches {
			if item != switchToAdd {
				newUpstreams = append(newUpstreams, item)
			}
		}
		newUpstreams = append(newUpstreams, switchToAdd)
		return newUpstreams
	}
}
