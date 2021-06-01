package controllers

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
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
		if r.checkDownstreamConnectionsRebuildNeeded(*switchConnection, ctx) || r.checkUpstreamConnectionsRebuildNeeded(*switchConnection, ctx) {
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
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchConnectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1alpha1.SwitchConnection{}).
		Watches(&source.Kind{Type: &switchv1alpha1.SwitchConnection{}}, handler.Funcs{
			UpdateFunc: r.handleConnectionUpdate(mgr.GetScheme(), &switchv1alpha1.SwitchConnectionList{}),
		}).
		Complete(r)
}

func (r *SwitchConnectionReconciler) handleConnectionUpdate(scheme *runtime.Scheme, ro runtime.Object) func(event.UpdateEvent, workqueue.RateLimitingInterface) {
	return func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
		err := enqueueSwitchConnectionReconcileRequest(r.Client, r.Log, scheme, q, ro)
		if err != nil {
			r.Log.Error(err, "error triggering switch connections reconciliation on connection update")
		}
	}
}

func enqueueSwitchConnectionReconcileRequest(c client.Client, log logr.Logger, scheme *runtime.Scheme, q workqueue.RateLimitingInterface, ro runtime.Object) error {
	ctx := context.Background()
	list := &unstructured.UnstructuredList{}
	gvk, err := apiutil.GVKForObject(ro, scheme)
	if err != nil {
		log.Error(err, "unable to get gvk")
		return err
	}
	list.SetGroupVersionKind(gvk)
	if err := c.List(ctx, list); err != nil {
		log.Error(err, "unable to get list of items")
		return err
	}
	for _, item := range list.Items {
		obj := &switchv1alpha1.SwitchConnection{}
		err := c.Get(ctx, types.NamespacedName{
			Namespace: item.GetNamespace(),
			Name:      item.GetName(),
		}, obj)
		if err != nil {
			log.Error(err, "failed to get switchConnection resource", "name", types.NamespacedName{
				Namespace: item.GetNamespace(),
				Name:      item.GetName(),
			})
			continue
		}
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Namespace: obj.Spec.Switch.Namespace,
			Name:      obj.Spec.Switch.Name,
		}})
		if obj.Spec.DownstreamSwitches != nil {
			for _, sw := range obj.Spec.DownstreamSwitches.Switches {
				if sw.Name != "" && sw.Namespace != "" {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Namespace: sw.Namespace,
						Name:      sw.Name,
					}})
				}
			}
		}
		if obj.Spec.UpstreamSwitches != nil {
			for _, sw := range obj.Spec.UpstreamSwitches.Switches {
				if sw.Name != "" && sw.Namespace != "" {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Namespace: sw.Namespace,
						Name:      sw.Name,
					}})
				}
			}
		}
	}
	return nil
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

	labelsReq, err := labels.NewRequirement(switchv1alpha1.LabelChassisId, selection.In, chassisIdsForLabels)
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
	connectionsList := &switchv1alpha1.SwitchConnectionList{}
	if err := r.Client.List(ctx, connectionsList); err != nil {
		return false, err
	}
	connectionsMap := map[uint8][]switchv1alpha1.SwitchConnection{}
	for _, item := range connectionsList.Items {
		if _, ok := connectionsMap[item.Spec.ConnectionLevel]; !ok {
			connectionsMap[item.Spec.ConnectionLevel] = []switchv1alpha1.SwitchConnection{item}
		} else {
			connectionsMap[item.Spec.ConnectionLevel] = append(connectionsMap[item.Spec.ConnectionLevel], item)
		}
	}

	for connLevel, connList := range connectionsMap {
		for _, item := range connList {
			if item.Spec.ConnectionLevel == connLevel { // loop through connection with connectionLevel == connLevel
				if checkSwitchInDownstreams(conn.Spec.Switch.ChassisID, item.Spec.DownstreamSwitches.Switches) {
					if item.Spec.ConnectionLevel == 255 {
						continue
					}
					if conn.Spec.ConnectionLevel != item.Spec.ConnectionLevel+1 {
						conn.Spec.ConnectionLevel = item.Spec.ConnectionLevel + 1
						update = true
					}
					if checkSwitchInDownstreams(item.Spec.Switch.ChassisID, conn.Spec.DownstreamSwitches.Switches) {
						conn.Spec.DownstreamSwitches.Switches = updateDownstreamSwitches(conn, item.Spec.Switch.ChassisID)
						conn.Spec.DownstreamSwitches.Count = len(conn.Spec.DownstreamSwitches.Switches)
						conn.Spec.UpstreamSwitches.Switches = updateUpstreamSwitches(conn, *item.Spec.Switch)
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

func checkSwitchInDownstreams(switchChassisId string, switches []switchv1alpha1.ConnectedSwitchSpec) bool {
	for _, item := range switches {
		if switchChassisId == item.ChassisID {
			return true
		}
	}
	return false
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
		newUpstreams := []switchv1alpha1.ConnectedSwitchSpec{switchToAdd}
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

func (r *SwitchConnectionReconciler) checkDownstreamConnectionsRebuildNeeded(conn switchv1alpha1.SwitchConnection, ctx context.Context) bool {
	if len(conn.Spec.DownstreamSwitches.Switches) == 0 {
		return false
	}
	downstreamConnections := conn.Spec.DownstreamSwitches.Switches
	downstreamSwitchesList := &switchv1alpha1.SwitchList{}
	chassisIdsForLabels := make([]string, 0)
	for _, item := range downstreamConnections {
		chassisIdsForLabels = append(chassisIdsForLabels, strings.ReplaceAll(item.ChassisID, ":", "-"))
	}

	labelsReq, err := labels.NewRequirement(switchv1alpha1.LabelChassisId, selection.In, chassisIdsForLabels)
	if err != nil {
		r.Log.Error(err, "unable to build label selector requirements")
		return true
	}
	selector := labels.NewSelector()
	selector = selector.Add(*labelsReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         1000,
	}
	err = r.Client.List(ctx, downstreamSwitchesList, opts)
	if err != nil {
		r.Log.Error(err, "unable to get related switches list")
		return true
	}
	for _, sw := range downstreamSwitchesList.Items {
		if sw.Spec.ConnectionLevel != conn.Spec.ConnectionLevel+1 {
			return true
		}
	}
	return false
}

func (r *SwitchConnectionReconciler) checkUpstreamConnectionsRebuildNeeded(conn switchv1alpha1.SwitchConnection, ctx context.Context) bool {
	if conn.Spec.UpstreamSwitches.Switches == nil {
		return false
	}
	if len(conn.Spec.UpstreamSwitches.Switches) == 0 {
		return false
	}
	upstreamConnections := conn.Spec.UpstreamSwitches.Switches
	upstreamSwitchesList := &switchv1alpha1.SwitchList{}
	chassisIdsForLabels := make([]string, 0)
	for _, item := range upstreamConnections {
		chassisIdsForLabels = append(chassisIdsForLabels, strings.ReplaceAll(item.ChassisID, ":", "-"))
	}

	labelsReq, err := labels.NewRequirement(switchv1alpha1.LabelChassisId, selection.In, chassisIdsForLabels)
	if err != nil {
		r.Log.Error(err, "unable to build label selector requirements")
		return true
	}
	selector := labels.NewSelector()
	selector = selector.Add(*labelsReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         1000,
	}
	err = r.Client.List(ctx, upstreamSwitchesList, opts)
	if err != nil {
		r.Log.Error(err, "unable to get related switches list")
		return true
	}
	for _, sw := range upstreamSwitchesList.Items {
		if sw.Spec.ConnectionLevel != conn.Spec.ConnectionLevel-1 {
			return true
		}
	}
	return false
}
