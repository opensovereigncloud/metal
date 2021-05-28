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

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
	"github.com/onmetal/switch-operator/util"
)

// SwitchReconciler reconciles a Switch object
type SwitchReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/finalizers,verbs=update
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=machines/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *SwitchReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("switch", req.NamespacedName)
	switchRes := &switchv1alpha1.Switch{}
	if err := r.Get(ctx, req.NamespacedName, switchRes); err != nil {
		if apierrors.IsNotFound(err) {
			log.Error(err, "requested switch resource not found", "name", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to get switch resource", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	// create switch connection resource if not exist
	connRes := &switchv1alpha1.SwitchConnection{}
	err := r.Get(ctx, req.NamespacedName, connRes)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			log.Error(err, "failed to get switch connection resource for switch", "name", req.NamespacedName)
			return ctrl.Result{}, err
		} else {
			// create connection
			connectedSwitchesIds := getDownstreamSwitches(switchRes)
			getPreparedSwitchConnection(connRes, switchRes, connectedSwitchesIds)
			if err := r.Client.Create(ctx, connRes); err != nil {
				log.Error(err, "unable to create switchConnection resource")
				return ctrl.Result{}, err
			}
		}
	} else {
		updateNeeded := false
		if switchRes.Spec.ConnectionLevel != connRes.Spec.ConnectionLevel {
			switchRes.Spec.ConnectionLevel = connRes.Spec.ConnectionLevel
			updateNeeded = true
		}
		switch switchRes.Spec.Role {
		case util.CUndefinedRole:
			if checkMachinesConnected(switchRes) {
				switchRes.Spec.Role = util.CLeafRole
			} else {
				switchRes.Spec.Role = util.CSpineRole
			}
			updateNeeded = true
		case util.CSpineRole:
			if checkMachinesConnected(switchRes) {
				switchRes.Spec.Role = util.CLeafRole
				updateNeeded = true
			}
		}
		if updateNeeded {
			if err := r.Update(ctx, switchRes); err != nil {
				log.Error(err, "unable to update switch resource", "name", req.NamespacedName)
				return ctrl.Result{}, err
			}
		}
	}

	//todo define subnet
	//if switchRes.Spec.SouthSubnetV4 == "" {
	//	_, err := r.findV4Subnet(ctx)
	//	if err != nil {
	//		log.Error(err, "unable to find suitable IPv4 subnet")
	//	}
	//}
	//if switchRes.Spec.SouthSubnetV4 == "" {
	//	_, err := r.findV6Subnet(ctx)
	//	if err != nil {
	//		log.Error(err, "unable to find suitable IPv6 subnet")
	//	}
	//}

	return ctrl.Result{RequeueAfter: util.CRequeueInterval}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1alpha1.Switch{}).
		Watches(&source.Kind{Type: &switchv1alpha1.SwitchConnection{}}, handler.Funcs{
			CreateFunc:  nil,
			UpdateFunc:  handleConnectionUpdate(r.Client, r.Log, mgr.GetScheme(), &switchv1alpha1.SwitchConnectionList{}),
			DeleteFunc:  nil,
			GenericFunc: nil,
		}).
		Complete(r)
}

func handleConnectionUpdate(c client.Client, log logr.Logger, scheme *runtime.Scheme, ro runtime.Object) func(event.UpdateEvent, workqueue.RateLimitingInterface) {
	return func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
		err := enqueueSwitchReconcileRequest(c, log, scheme, q, ro)
		if err != nil {
			log.Error(err, "error triggering switch reconciliation on connection update")
		}
	}
}

func enqueueSwitchReconcileRequest(c client.Client, log logr.Logger, scheme *runtime.Scheme, q workqueue.RateLimitingInterface, ro runtime.Object) error {
	list := &unstructured.UnstructuredList{}
	gvk, err := apiutil.GVKForObject(ro, scheme)
	if err != nil {
		log.Error(err, "unable to get gvk")
		return err
	}
	list.SetGroupVersionKind(gvk)
	if err := c.List(context.Background(), list); err != nil {
		log.Error(err, "unable to get list of items")
		return err
	}
	for _, item := range list.Items {
		data, err := json.Marshal(item)
		if err != nil {
			log.Error(err, "unable to marshal data")
			return err
		}
		obj := &switchv1alpha1.SwitchConnection{}
		err = json.Unmarshal(data, obj)
		if err != nil {
			log.Error(err, "unable to unmarshal data")
			return err
		}
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
	}
	return nil
}

func getDownstreamSwitches(sw *switchv1alpha1.Switch) []string {
	connMap := make(map[string]struct{})
	downstreamSwitches := make([]string, 0)
	for _, iface := range sw.Spec.Interfaces {
		if iface.Neighbour == util.CSwitchType {
			if _, ok := connMap[iface.LLDPChassisID]; !ok {
				connMap[iface.LLDPChassisID] = struct{}{}
				downstreamSwitches = append(downstreamSwitches, iface.LLDPChassisID)
			}
		}
	}
	return downstreamSwitches
}

func getPreparedSwitchConnection(conn *switchv1alpha1.SwitchConnection, sw *switchv1alpha1.Switch, connectedSwitches []string) {
	connectedSwitchesSpecs := make([]switchv1alpha1.ConnectedSwitchSpec, 0)
	for _, id := range connectedSwitches {
		connectedSwitchesSpecs = append(connectedSwitchesSpecs, switchv1alpha1.ConnectedSwitchSpec{ChassisID: id})
	}

	conn.ObjectMeta = metav1.ObjectMeta{
		Name:      sw.Name,
		Namespace: sw.Namespace,
		Labels: map[string]string{
			util.ConnectionLabelChassisId: strings.ReplaceAll(sw.Spec.SwitchChassis.ChassisID, ":", "-"),
		},
	}
	conn.Spec = switchv1alpha1.SwitchConnectionSpec{
		Switch: &switchv1alpha1.ConnectedSwitchSpec{
			Name:      sw.Name,
			Namespace: sw.Namespace,
			ChassisID: sw.Spec.SwitchChassis.ChassisID,
		},
		UpstreamSwitches: &switchv1alpha1.UpstreamSwitchesSpec{
			Count:    0,
			Switches: nil,
		},
		DownstreamSwitches: &switchv1alpha1.DownstreamSwitchesSpec{
			Count:    len(connectedSwitchesSpecs),
			Switches: connectedSwitchesSpecs,
		},
		ConnectionLevel: sw.Spec.ConnectionLevel,
	}
}

func checkMachinesConnected(sw *switchv1alpha1.Switch) bool {
	for _, nic := range sw.Spec.Interfaces {
		if nic.Neighbour == util.CMachineType {
			return true
		}
	}
	return false
}

//func (r *SwitchReconciler) findV4Subnet(ctx context.Context) (*subnetv1alpha1.Subnet, error) {
//	subnetsList := &subnetv1alpha1.SubnetList{}
//	err := r.Client.List(ctx, subnetsList)
//	if err != nil {
//		return nil, err
//	}
//	for _, item := range subnetsList.Items {
//		if item.Spec.Type == "IPv4" {
//			return &item, nil
//		}
//	}
//	return nil, nil
//}
//
//func (r *SwitchReconciler) findV6Subnet(ctx context.Context) (*subnetv1alpha1.Subnet, error) {
//	subnetsList := &subnetv1alpha1.SubnetList{}
//	err := r.Client.List(ctx, subnetsList)
//	if err != nil {
//		return nil, err
//	}
//	for _, item := range subnetsList.Items {
//		if item.Spec.Type == "IPv6" {
//			return &item, nil
//		}
//	}
//	return nil, nil
//}
