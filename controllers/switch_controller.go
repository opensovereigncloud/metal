/*
Copyright 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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
	"fmt"
	"math/big"
	"net"
	"reflect"
	"sort"
	"strconv"
	"strings"

	gocidr "github.com/apparentlymart/go-cidr/cidr"
	"github.com/go-logr/logr"
	netglobalv1alpha1 "github.com/onmetal/k8s-network-global/api/v1alpha1"
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
	updateSouthInterfacesAddresses(switchRes)
	r.updateNorthInterfacesAddresses(switchRes, ctx)

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

//handleSwitchUpdate handler for UpdateEvent for switch resources
func (r *SwitchReconciler) handleSwitchUpdate(scheme *runtime.Scheme, ro runtime.Object) func(event.UpdateEvent, workqueue.RateLimitingInterface) {
	return func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
		err := enqueueSwitchReconcileRequest(r.Client, r.Log, scheme, q, ro)
		if err != nil {
			r.Log.Error(err, "error triggering switch reconciliation on connections update")
		}
	}
}

//enqueueSwitchReconcileRequest adds related switch resources
//to the reconciliation queue
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

//findSouthNeighboursSwitches returns a SwitchList resource, that includes
//"downstream" switches or an error.
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

//findNorthNeighboursSwitches returns a SwitchList resource, that includes
//"upstream" switches or an error.
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

//updateConnectionLevel calculates switch's connection level,
//basing on neighbours connection levels. Returns an error.
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

//defineSubnets process related resources to define switch's
//"south" subnets. Returns an error.
func (r *SwitchReconciler) defineSubnets(sw *switchv1alpha1.Switch, ctx context.Context) error {
	topLevelSwitch := r.GetTopLevelSwitch(sw, ctx)
	if topLevelSwitch == nil {
		return nil
	}
	relatedAssignment := r.GetRelatedAssignment(topLevelSwitch, ctx)
	if relatedAssignment == nil {
		return nil
	}
	regions := []string{relatedAssignment.Spec.Region}
	zones := []string{relatedAssignment.Spec.AvailabilityZone}
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

//GetTopLevelSwitch recursively searches for top level switch
//related to the target switch.
func (r *SwitchReconciler) GetTopLevelSwitch(sw *switchv1alpha1.Switch, ctx context.Context) *switchv1alpha1.Switch {
	if sw.Spec.State.NorthConnections != nil && sw.Spec.State.NorthConnections.Count != 0 {
		nextLevelSwitch := &switchv1alpha1.Switch{}
		if err := r.Client.Get(ctx, types.NamespacedName{
			Namespace: sw.Spec.State.NorthConnections.Connections[0].Namespace,
			Name:      sw.Spec.State.NorthConnections.Connections[0].Name,
		}, nextLevelSwitch); err != nil {
			return nil
		}
		return r.GetTopLevelSwitch(nextLevelSwitch, ctx)
	}
	if sw.Spec.State.ConnectionLevel == 0 {
		return sw
	}
	return nil
}

//GetRelatedAssignment searches for switch assignment resource
//related to target top level switch.
func (r *SwitchReconciler) GetRelatedAssignment(sw *switchv1alpha1.Switch, ctx context.Context) *switchv1alpha1.SwitchAssignment {
	swaList := &switchv1alpha1.SwitchAssignmentList{}
	selector := labels.SelectorFromSet(labels.Set{switchv1alpha1.LabelChassisId: sw.Labels[switchv1alpha1.LabelChassisId]})
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         1000,
	}
	if err := r.List(ctx, swaList, opts); err != nil {
		r.Log.Error(err, "unable to get switch assignments list")
	}
	if len(swaList.Items) == 0 {
		return nil
	}
	return &swaList.Items[0]
}

//finalizeSwitch prepare environment for switch resource deletion
//by updating related resources. Returns an error.
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

//getMinConnectionLevel calculates the minimum connection level
//value among switches in the list provided as argument.
func getMinConnectionLevel(switchList []switchv1alpha1.Switch) uint8 {
	result := uint8(255)
	for _, item := range switchList {
		if item.Spec.State.ConnectionLevel < result {
			result = item.Spec.State.ConnectionLevel
		}
	}
	return result
}

//constructNeighboursFromSwitchList creates list of neighbours
//specs from the list of switches.
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

//updateNorthConnections updates switch resource north connections list.
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

//updateSouthConnections updates switch resource south connections list.
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

//removeFromSouthConnections removes item from switch resource
//south connections list if this item presents in north connections.
func removeFromSouthConnections(sw *switchv1alpha1.Switch, ncm map[string]struct{}) {
	connections := make([]switchv1alpha1.NeighbourSpec, 0)
	for _, item := range sw.Spec.State.SouthConnections.Connections {
		if _, ok := ncm[item.ChassisID]; !ok {
			connections = append(connections, item)
		}
	}
	sw.Spec.State.SouthConnections.Connections = connections
}

//getSuitableSubnet finds the subnet resource, that fits address
//type, region, availability zones and addresses count. It returns
//pointers to the CIDR and subnet resource objects or an error.
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
			addressesLeft := sn.Status.CapacityLeft
			if sn.Status.Type == addressType && addressesLeft.CmpInt64(addressesNeeded) >= 0 {
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

//getMinimalVacantCIDR calculates the minimal suitable network
//from the networks list provided as argument according to the
//needed addresses count. It returns the pointer to the CIDR object.
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
		if cidr.AddressCapacity().Cmp(minSuitableNet.AddressCapacity()) < 0 &&
			cidr.AddressCapacity().Cmp(new(big.Int).SetInt64(addressesCount)) >= 0 {
			minSuitableNet = &cidr
		}
	}
	return minSuitableNet
}

//updateSouthInterfacesAddresses sets IP addresses for switch
//south interfaces according to switch south subnets values.
func updateSouthInterfacesAddresses(sw *switchv1alpha1.Switch) {
	interfaces := sw.GetSwitchPorts()
	sort.Slice(interfaces, func(i, j int) bool {
		leftIndex, _ := strconv.Atoi(strings.ReplaceAll(interfaces[i].Name, "Ethernet", ""))
		rightIndex, _ := strconv.Atoi(strings.ReplaceAll(interfaces[j].Name, "Ethernet", ""))
		return leftIndex < rightIndex
	})
	for _, iface := range interfaces {
		if sw.Spec.SouthSubnetV4 != nil && iface.IPv4 == "" && iface.LLDPChassisID != "" {
			for _, item := range sw.Spec.State.SouthConnections.Connections {
				if item.ChassisID == iface.LLDPChassisID {
					_, network, _ := net.ParseCIDR(sw.Spec.SouthSubnetV4.CIDR)
					ifaceSubnet := getInterfaceSubnet(network, iface, subnetv1alpha1.CIPv4SubnetType)
					ifaceAddress, _ := gocidr.Host(ifaceSubnet, 1)
					iface.IPv4 = fmt.Sprintf("%s/%d", ifaceAddress.String(), switchv1alpha1.CIPv4InterfaceSubnetMask)
				}
			}
		}
		if sw.Spec.SouthSubnetV6 != nil && iface.IPv6 == "" && iface.LLDPChassisID != "" {
			for _, item := range sw.Spec.State.SouthConnections.Connections {
				if item.ChassisID == iface.LLDPChassisID {
					_, network, _ := net.ParseCIDR(sw.Spec.SouthSubnetV6.CIDR)
					ifaceSubnet := getInterfaceSubnet(network, iface, subnetv1alpha1.CIPv6SubnetType)
					ifaceAddress, _ := gocidr.Host(ifaceSubnet, 0)
					iface.IPv4 = fmt.Sprintf("%s/%d", ifaceAddress.String(), switchv1alpha1.CIPv6InterfaceSubnetMask)
				}
			}
		}
	}
}

//updateNorthInterfacesAddresses sets IP addresses for switch
//north interfaces according to north switch subnets values.
func (r *SwitchReconciler) updateNorthInterfacesAddresses(sw *switchv1alpha1.Switch, ctx context.Context) {
	if sw.Spec.State.NorthConnections.Count > 0 {
		for _, item := range sw.Spec.State.NorthConnections.Connections {
			northSwitch := &switchv1alpha1.Switch{}
			northIPv4Subnet := &net.IPNet{}
			northIPv6Subnet := &net.IPNet{}
			if item.Name != "" && item.Namespace != "" {
				if err := r.Client.Get(ctx, types.NamespacedName{
					Namespace: item.Namespace,
					Name:      item.Name,
				}, northSwitch); err != nil {
					if apierrors.IsNotFound(err) {
						item.Namespace = ""
						item.Name = ""
						continue
					}
					continue
				}
				// got north switch resource
				for _, iface := range northSwitch.Spec.Interfaces {
					if iface.LLDPChassisID == sw.Spec.SwitchChassis.ChassisID {
						if iface.IPv4 != "" {
							_, northIPv4Subnet, _ = net.ParseCIDR(iface.IPv4)
						}
						if iface.IPv6 != "" {
							_, northIPv6Subnet, _ = net.ParseCIDR(iface.IPv6)
						}
						break
					}
				}
			}
			for _, iface := range sw.Spec.Interfaces {
				if iface.LLDPChassisID == northSwitch.Spec.SwitchChassis.ChassisID {
					if addr, err := gocidr.Host(northIPv4Subnet, 2); err == nil {
						iface.IPv4 = fmt.Sprintf("%s/%d", addr.String(), switchv1alpha1.CIPv4InterfaceSubnetMask)
					}
					if addr, err := gocidr.Host(northIPv6Subnet, 1); err == nil {
						iface.IPv6 = fmt.Sprintf("%s/%d", addr.String(), switchv1alpha1.CIPv6InterfaceSubnetMask)
					}
					break
				}
			}
		}
	}
}

func getInterfaceSubnet(network *net.IPNet, iface *switchv1alpha1.InterfaceSpec, addrType subnetv1alpha1.SubnetAddressType) *net.IPNet {
	index, _ := strconv.Atoi(strings.ReplaceAll(iface.Name, "Ethernet", ""))
	prefix, _ := network.Mask.Size()
	ifaceNet, _ := gocidr.Subnet(network, getInterfaceSubnetMaskLength(addrType)-prefix, index)
	return ifaceNet
}

func getInterfaceSubnetMaskLength(addrType subnetv1alpha1.SubnetAddressType) int {
	if addrType == subnetv1alpha1.CIPv4SubnetType {
		return switchv1alpha1.CIPv4InterfaceSubnetMask
	} else {
		return switchv1alpha1.CIPv6InterfaceSubnetMask
	}
}
