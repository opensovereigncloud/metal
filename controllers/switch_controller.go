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
	"net"
	"reflect"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	netglobalv1alpha1 "github.com/onmetal/k8s-network-global/api/v1alpha1"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	subnetv1alpha1 "github.com/onmetal/k8s-subnet/api/v1alpha1"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
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
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=subnets,verbs=get;list;watch;update
//+kubebuilder:rbac:groups=machine.onmetal.de,resources=subnets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch

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
			log.Info("requested switch resource not found", "name", req.NamespacedName)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		log.Error(err, "failed to get switch resource", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	if switchRes.DeletionTimestamp != nil {
		if controllerutil.ContainsFinalizer(switchRes, switchv1alpha1.CSwitchFinalizer) {
			if err := r.finalizeSwitch(switchRes, ctx); err != nil {
				log.Error(err, "failed to finalize switch")
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(switchRes, switchv1alpha1.CSwitchFinalizer)
			if err := r.Client.Update(ctx, switchRes); err != nil {
				log.Error(err, "failed to update switch resource on finalizer removal")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	oldRes := switchRes.DeepCopy()

	if !controllerutil.ContainsFinalizer(switchRes, switchv1alpha1.CSwitchFinalizer) {
		controllerutil.AddFinalizer(switchRes, switchv1alpha1.CSwitchFinalizer)
	}

	switch switchRes.Spec.State.Role {
	case switchv1alpha1.CUndefinedRole:
		switchRes.Spec.State.Role = switchv1alpha1.CSpineRole
		fallthrough
	case switchv1alpha1.CSpineRole:
		if switchRes.CheckMachinesConnected() {
			switchRes.Spec.State.Role = switchv1alpha1.CLeafRole
		}
	}

	if switchRes.Spec.SouthSubnetV4 == nil || switchRes.Spec.SouthSubnetV6 == nil {
		if err := r.defineSubnets(switchRes, ctx); err != nil {
			log.Info("unable to define subnet(s)", "error", err.Error())
		}
	}

	if switchRes.CheckNorthNeighboursDataUpdateNeeded() || switchRes.CheckSouthNeighboursDataUpdateNeeded() {
		switch switchRes.Spec.State.ConnectionLevel {
		case 0:
			southSwitchList, err := r.findSouthNeighboursSwitches(switchRes, ctx)
			if err != nil {
				log.Error(err, "failed to get south switch neighbours")
				return ctrl.Result{}, err
			}
			neighboursMap := constructNeighboursFromSwitchList(southSwitchList.Items)
			connections := switchRes.Spec.State.SouthConnections.Connections
			for i, neighbour := range connections {
				if _, ok := neighboursMap[neighbour.ChassisID]; ok {
					connections[i] = neighboursMap[neighbour.ChassisID]
				}
			}
			switchRes.Spec.State.SouthConnections.Connections = connections
		default:
			if err := r.updateConnectionLevel(switchRes, ctx); err != nil {
				log.Error(err, "failed to update switch connection level")
				return ctrl.Result{}, err
			}
		}
	}
	updateInterfacesAddresses(switchRes)

	if !reflect.DeepEqual(oldRes, switchRes) {
		if err := r.Client.Update(ctx, switchRes); err != nil {
			log.Error(err, "failed to update switch resource")
			return ctrl.Result{}, err
		}
	}

	if switchRes.Spec.SouthSubnetV4 == nil || switchRes.Spec.SouthSubnetV6 == nil {
		return ctrl.Result{RequeueAfter: switchv1alpha1.CSwitchRequeueInterval}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1alpha1.Switch{}).
		Watches(&source.Kind{Type: &switchv1alpha1.Switch{}}, handler.Funcs{
			UpdateFunc: r.handleSwitchUpdate(mgr.GetScheme(), &switchv1alpha1.SwitchList{}),
		}).
		Complete(r)
}

func (r *SwitchReconciler) handleSwitchUpdate(scheme *runtime.Scheme, ro runtime.Object) func(event.UpdateEvent, workqueue.RateLimitingInterface) {
	return func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
		err := enqueueSwitchReconcileRequest(r.Client, r.Log, scheme, q, ro)
		if err != nil {
			r.Log.Error(err, "error triggering switch reconciliation on connections update")
		}
	}
}

func enqueueSwitchReconcileRequest(c client.Client, log logr.Logger, scheme *runtime.Scheme, q workqueue.RateLimitingInterface, ro runtime.Object) error {
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
		obj := &switchv1alpha1.Switch{}
		err := c.Get(ctx, types.NamespacedName{
			Namespace: item.GetNamespace(),
			Name:      item.GetName(),
		}, obj)
		if err != nil {
			if !apierrors.IsNotFound(err) {
				log.Error(err, "failed to get switch resource", "name", types.NamespacedName{
					Namespace: item.GetNamespace(),
					Name:      item.GetName(),
				})
			}
			continue
		}
		if obj.Spec.State.SouthConnections != nil {
			for _, neighbour := range obj.Spec.State.SouthConnections.Connections {
				if neighbour.Name != "" && neighbour.Namespace != "" && neighbour.Type == switchv1alpha1.CSwitchType {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Namespace: neighbour.Namespace,
						Name:      neighbour.Name,
					}})
				}
			}
		}
	}
	return nil
}

//common function
func (r *SwitchReconciler) findSouthNeighboursSwitches(switchRes *switchv1alpha1.Switch, ctx context.Context) (*switchv1alpha1.SwitchList, error) {
	swList := &switchv1alpha1.SwitchList{}
	connectionsChassisIds := make([]string, 0, len(switchRes.Spec.State.SouthConnections.Connections))
	for _, item := range switchRes.Spec.State.SouthConnections.Connections {
		if item.Type == switchv1alpha1.CSwitchType {
			connectionsChassisIds = append(connectionsChassisIds, strings.ReplaceAll(item.ChassisID, ":", "-"))
		}
	}
	if len(connectionsChassisIds) == 0 {
		return swList, nil
	}

	labelsReq, err := labels.NewRequirement(switchv1alpha1.LabelChassisId, selection.In, connectionsChassisIds)
	if err != nil {
		return nil, err
	}
	selector := labels.NewSelector()
	selector = selector.Add(*labelsReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         1000,
	}
	if err := r.Client.List(ctx, swList, opts); err != nil {
		return nil, err
	}
	return swList, nil
}

func (r *SwitchReconciler) findNorthNeighboursSwitches(switchRes *switchv1alpha1.Switch, ctx context.Context) (*switchv1alpha1.SwitchList, error) {
	swList := &switchv1alpha1.SwitchList{}
	connectionsChassisIds := make([]string, 0, len(switchRes.Spec.State.NorthConnections.Connections))
	for _, item := range switchRes.Spec.State.NorthConnections.Connections {
		if item.Type == switchv1alpha1.CSwitchType {
			connectionsChassisIds = append(connectionsChassisIds, strings.ReplaceAll(item.ChassisID, ":", "-"))
		}
	}
	if len(connectionsChassisIds) == 0 {
		return swList, nil
	}

	labelsReq, err := labels.NewRequirement(switchv1alpha1.LabelChassisId, selection.In, connectionsChassisIds)
	if err != nil {
		return nil, err
	}
	selector := labels.NewSelector()
	selector = selector.Add(*labelsReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         1000,
	}
	if err := r.Client.List(ctx, swList, opts); err != nil {
		return nil, err
	}
	return swList, nil
}

func (r *SwitchReconciler) updateConnectionLevel(sw *switchv1alpha1.Switch, ctx context.Context) error {
	swList := &switchv1alpha1.SwitchList{}
	if err := r.Client.List(ctx, swList); err != nil {
		return err
	}

	connectionLevelMap := map[uint8][]switchv1alpha1.Switch{}
	keys := make([]uint8, 0)
	for _, item := range swList.Items {
		if _, ok := connectionLevelMap[item.Spec.State.ConnectionLevel]; !ok {
			connectionLevelMap[item.Spec.State.ConnectionLevel] = []switchv1alpha1.Switch{item}
			keys = append(keys, item.Spec.State.ConnectionLevel)
		} else {
			connectionLevelMap[item.Spec.State.ConnectionLevel] = append(connectionLevelMap[item.Spec.State.ConnectionLevel], item)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, connLevel := range keys {
		switches := connectionLevelMap[connLevel]
		switchNorthNeighbours := sw.GetNorthSwitchConnection(switches)
		if len(switchNorthNeighbours) > 0 {
			minConnLevel := getMinConnectionLevel(switchNorthNeighbours)
			if minConnLevel != 255 && minConnLevel < sw.Spec.State.ConnectionLevel {
				sw.Spec.State.ConnectionLevel = minConnLevel + 1
				northNeighboursMap := constructNeighboursFromSwitchList(switchNorthNeighbours)
				updateNorthConnections(sw, northNeighboursMap)
				ncm := map[string]struct{}{}
				for _, conn := range sw.Spec.State.NorthConnections.Connections {
					if _, ok := ncm[conn.ChassisID]; !ok {
						ncm[conn.ChassisID] = struct{}{}
					}
				}
				removeFromSouthConnections(sw, ncm)
				switchSouthNeighbours, err := r.findSouthNeighboursSwitches(sw, ctx)
				if err != nil {
					return err
				}
				southNeighboursMap := constructNeighboursFromSwitchList(switchSouthNeighbours.Items)
				updateSouthConnections(sw, southNeighboursMap)
				sw.Spec.State.NorthConnections.Count = len(sw.Spec.State.NorthConnections.Connections)
				sw.Spec.State.SouthConnections.Count = len(sw.Spec.State.SouthConnections.Connections)
			}
		}
	}
	return nil
}

func (r *SwitchReconciler) defineSubnets(sw *switchv1alpha1.Switch, ctx context.Context) error {
	cm := &v1.ConfigMap{}
	// todo: discuss configMap parameters and data structure
	if err := r.Client.Get(ctx, types.NamespacedName{Name: "partition-cm", Namespace: "default"}, cm); err != nil {
		if apierrors.IsNotFound(err) {
			return errors.New("requested configMap not found")
		}
		return err
	}
	regions := strings.Split(cm.Data["regions"], ",")
	zones := strings.Split(cm.Data["availabilityZones"], ",")
	subnetList := &subnetv1alpha1.SubnetList{}
	if err := r.Client.List(ctx, subnetList); err != nil {
		return err
	}
	if sw.Spec.SouthSubnetV4 == nil {
		addressesCount := sw.GetAddressNeededCount(subnetv1alpha1.CIPv4SubnetType)
		cidr, sn, err := getSuitableSubnet(sw, subnetList, subnetv1alpha1.CIPv4SubnetType, regions, zones, addressesCount)
		if err != nil {
			return err
		}
		if cidr != nil && sn != nil {
			sw.Spec.SouthSubnetV4 = &switchv1alpha1.SwitchSubnetSpec{
				ParentSubnet: &switchv1alpha1.ParentSubnetSpec{
					Namespace: sn.Namespace,
					Name:      sn.Name,
				},
				CIDR: cidr.String(),
			}
			if err := r.Status().Update(ctx, sn); err != nil {
				return err
			}
		}
	}
	if sw.Spec.SouthSubnetV6 == nil {
		addressCount := sw.GetAddressNeededCount(subnetv1alpha1.CIPv4SubnetType)
		cidr, sn, err := getSuitableSubnet(sw, subnetList, subnetv1alpha1.CIPv6SubnetType, regions, zones, addressCount)
		if err != nil {
			return err
		}
		if cidr != nil && sn != nil {
			sw.Spec.SouthSubnetV6 = &switchv1alpha1.SwitchSubnetSpec{
				ParentSubnet: &switchv1alpha1.ParentSubnetSpec{
					Namespace: sn.Namespace,
					Name:      sn.Name,
				},
				CIDR: cidr.String(),
			}
			if err := r.Status().Update(ctx, sn); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *SwitchReconciler) finalizeSwitch(sw *switchv1alpha1.Switch, ctx context.Context) error {
	if sw.Spec.SouthSubnetV4 != nil {
		subnetV4 := &subnetv1alpha1.Subnet{}
		if err := r.Client.Get(ctx, types.NamespacedName{
			Namespace: sw.Spec.SouthSubnetV4.ParentSubnet.Namespace,
			Name:      sw.Spec.SouthSubnetV4.ParentSubnet.Name,
		}, subnetV4); err != nil {
			if apierrors.IsNotFound(err) {
				r.Log.Error(err, "subnet not found")
				return nil
			}
			r.Log.Error(err, "failed to get subnet")
			return err
		}
		cidrToRelease, err := netglobalv1alpha1.CIDRFromString(sw.Spec.SouthSubnetV4.CIDR)
		if err != nil {
			r.Log.Error(err, "failed to get switch south network CIDR to release")
			return err
		}
		if subnetV4.CanRelease(cidrToRelease) {
			err := subnetV4.Release(cidrToRelease)
			if err != nil {
				r.Log.Error(err, "failed to release switch south network CIDR")
				return err
			}
			if err := r.Status().Update(ctx, subnetV4); err != nil {
				r.Log.Error(err, "failed to update subnet status")
				return err
			}
		}
	}
	if sw.Spec.SouthSubnetV6 != nil {
		subnetV6 := &subnetv1alpha1.Subnet{}
		if err := r.Client.Get(ctx, types.NamespacedName{
			Namespace: sw.Spec.SouthSubnetV6.ParentSubnet.Namespace,
			Name:      sw.Spec.SouthSubnetV6.ParentSubnet.Name,
		}, subnetV6); err != nil {
			if apierrors.IsNotFound(err) {
				r.Log.Error(err, "subnet not found")
				return nil
			}
			r.Log.Error(err, "failed to get subnet")
			return err
		}
		cidrToRelease, err := netglobalv1alpha1.CIDRFromString(sw.Spec.SouthSubnetV6.CIDR)
		if err != nil {
			r.Log.Error(err, "failed to get switch south network CIDR to release")
			return err
		}
		if subnetV6.CanRelease(cidrToRelease) {
			err := subnetV6.Release(cidrToRelease)
			if err != nil {
				r.Log.Error(err, "failed to release switch south network CIDR")
				return err
			}
			if err := r.Status().Update(ctx, subnetV6); err != nil {
				r.Log.Error(err, "failed to update subnet status")
				return err
			}
		}
	}
	return nil
}

func getMinConnectionLevel(switchList []switchv1alpha1.Switch) uint8 {
	result := uint8(255)
	for _, item := range switchList {
		if item.Spec.State.ConnectionLevel < result {
			result = item.Spec.State.ConnectionLevel
		}
	}
	return result
}

func constructNeighboursFromSwitchList(swl []switchv1alpha1.Switch) map[string]switchv1alpha1.NeighbourSpec {
	neighbours := map[string]switchv1alpha1.NeighbourSpec{}
	for _, item := range swl {
		neighbours[item.Spec.SwitchChassis.ChassisID] = switchv1alpha1.NeighbourSpec{
			Name:      item.Name,
			Namespace: item.Namespace,
			ChassisID: item.Spec.SwitchChassis.ChassisID,
			Type:      switchv1alpha1.CSwitchType,
		}
	}
	return neighbours
}

func updateNorthConnections(sw *switchv1alpha1.Switch, ncm map[string]switchv1alpha1.NeighbourSpec) {
	connections := make([]switchv1alpha1.NeighbourSpec, 0)
	if sw.Spec.State.NorthConnections.Connections == nil || len(sw.Spec.State.NorthConnections.Connections) == 0 {
		for _, value := range ncm {
			connections = append(connections, value)
		}
	} else {
		connections = sw.Spec.State.NorthConnections.Connections
		for i, neighbour := range connections {
			if _, ok := ncm[neighbour.ChassisID]; ok {
				connections[i] = ncm[neighbour.ChassisID]
			}
		}
	}
	sw.Spec.State.NorthConnections.Connections = connections
}

func updateSouthConnections(sw *switchv1alpha1.Switch, ncm map[string]switchv1alpha1.NeighbourSpec) {
	connections := make([]switchv1alpha1.NeighbourSpec, 0)
	if sw.Spec.State.SouthConnections.Connections == nil || len(sw.Spec.State.SouthConnections.Connections) == 0 {
		for _, value := range ncm {
			connections = append(connections, value)
		}
	} else {
		connections = sw.Spec.State.SouthConnections.Connections
		for i, neighbour := range connections {
			if _, ok := ncm[neighbour.ChassisID]; ok {
				connections[i] = ncm[neighbour.ChassisID]
			}
		}
	}
	sw.Spec.State.SouthConnections.Connections = connections
}

func removeFromSouthConnections(sw *switchv1alpha1.Switch, ncm map[string]struct{}) {
	connections := make([]switchv1alpha1.NeighbourSpec, 0)
	for _, item := range sw.Spec.State.SouthConnections.Connections {
		if _, ok := ncm[item.ChassisID]; !ok {
			connections = append(connections, item)
		}
	}
	sw.Spec.State.SouthConnections.Connections = connections
}

func getSuitableSubnet(
	sw *switchv1alpha1.Switch,
	subnetList *subnetv1alpha1.SubnetList,
	addressType subnetv1alpha1.SubnetAddressType,
	regions []string,
	zones []string,
	addressesNeeded int64) (*netglobalv1alpha1.CIDR, *subnetv1alpha1.Subnet, error) {

	for _, sn := range subnetList.Items {
		if sn.Spec.NetworkGlobalName == "underlay" &&
			reflect.DeepEqual(sn.Spec.Regions, regions) &&
			reflect.DeepEqual(sn.Spec.AvailabilityZones, zones) {
			addressesLeft, _ := sn.Status.CapacityLeft.AsInt64()
			if sn.Status.Type == addressType && addressesLeft >= addressesNeeded {
				minVacantCIDR := getMinimalVacantCIDR(sn.Status.Vacant, addressType, addressesNeeded)
				mask := sw.GetNeededMask(subnetv1alpha1.CIPv4SubnetType, float64(addressesNeeded))
				addr := minVacantCIDR.Net.IP
				network := &net.IPNet{
					IP:   addr,
					Mask: mask,
				}
				cidrCandidate := &netglobalv1alpha1.CIDR{Net: network}
				if sn.CanReserve(cidrCandidate) {
					if err := sn.Reserve(cidrCandidate); err != nil {
						return nil, nil, err
					} else {
						return cidrCandidate, &sn, nil
					}
				}
			}
		}
	}
	return nil, nil, nil
}

func getMinimalVacantCIDR(vacant []netglobalv1alpha1.CIDR, addressType subnetv1alpha1.SubnetAddressType, addressesCount int64) *netglobalv1alpha1.CIDR {
	zeroNetString := ""
	if addressType == subnetv1alpha1.CIPv4SubnetType {
		zeroNetString = switchv1alpha1.CIPv4ZeroNet
	} else {
		zeroNetString = switchv1alpha1.CIPv6ZeroNet
	}
	_, zeroNet, _ := net.ParseCIDR(zeroNetString)
	minSuitableNet := netglobalv1alpha1.CIDRFromNet(zeroNet)
	for _, cidr := range vacant {
		if cidr.AddressCapacity().Int64() < minSuitableNet.AddressCapacity().Int64() &&
			cidr.AddressCapacity().Int64() >= addressesCount {
			minSuitableNet = &cidr
		}
	}
	return minSuitableNet
}

func updateInterfacesAddresses(sw *switchv1alpha1.Switch) {
	// todo: assign addresses for "south" interfaces
}
