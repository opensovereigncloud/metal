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
	"net"
	"reflect"

	"github.com/go-logr/logr"
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

	subnetv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"

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
//+kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets,verbs=get;list;watch;update
//+kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets/status,verbs=get;update;patch

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

	if switchRes.CheckMachinesConnected() {
		switchRes.Spec.Role = switchv1alpha1.CLeafRole
	}

	if switchRes.Spec.SouthSubnetV4 == nil || switchRes.Spec.SouthSubnetV6 == nil {
		if err := r.defineSubnets(switchRes, ctx); err != nil {
			log.Info("unable to define subnet(s)", "error", err.Error())
		}
	}

	if err := r.fillSouthConnections(switchRes, ctx); err != nil {
		log.Error(err, "failed to get south switch neighbours")
		return ctrl.Result{}, err
	}

	if err := r.updateConnectionLevel(switchRes, ctx); err != nil {
		log.Error(err, "failed to update switch connection level")
		return ctrl.Result{}, err
	}

	switchRes.FlushAddresses()
	switchRes.UpdateSouthInterfacesAddresses()
	r.updateNorthInterfacesAddresses(switchRes, ctx)

	if !reflect.DeepEqual(oldRes, switchRes) {
		if err := r.Client.Update(ctx, switchRes); err != nil {
			log.Error(err, "failed to update switch resource")
			return ctrl.Result{}, err
		}
	}

	if !switchRes.AddressAssigned() {
		return ctrl.Result{RequeueAfter: switchv1alpha1.CSwitchRequeueInterval}, nil
	}
	if !r.NeighboursOk(switchRes, ctx) {
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
			connections := obj.Spec.State.SouthConnections.Connections
			connections = append(connections, obj.Spec.State.NorthConnections.Connections...)
			for _, neighbour := range connections {
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
			connectionsChassisIds = append(connectionsChassisIds, switchv1alpha1.MacToLabel(item.ChassisID))
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
			connectionsChassisIds = append(connectionsChassisIds, switchv1alpha1.MacToLabel(item.ChassisID))
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
	if sw.Spec.ConnectionLevel == 0 {
		sw.MoveNeighbours(swList)
	} else {
		swList.FillConnections(sw)
		connectionsMap, keys := swList.BuildConnectionsMap()
		for _, connLevel := range keys {
			switches := connectionsMap[connLevel]
			switchNorthNeighbours := sw.GetNorthSwitchConnection(switches)
			if len(switchNorthNeighbours.Items) > 0 {
				minConnLevel := switchNorthNeighbours.GetMinConnectionLevel()
				if minConnLevel != 255 && minConnLevel < sw.Spec.ConnectionLevel {
					sw.Spec.ConnectionLevel = minConnLevel + 1
					northNeighboursMap := switchNorthNeighbours.ConstructNeighboursFromSwitchList()
					sw.UpdateNorthConnections(northNeighboursMap)
					ncm := map[string]struct{}{}
					for _, conn := range sw.Spec.State.NorthConnections.Connections {
						if _, ok := ncm[conn.ChassisID]; !ok {
							ncm[conn.ChassisID] = struct{}{}
						}
					}
					sw.RemoveFromSouthConnections(ncm)
					sw.MoveNeighbours(swList)
				}
			}
		}
	}
	sw.Spec.State.NorthConnections.Count = len(sw.Spec.State.NorthConnections.Connections)
	sw.Spec.State.SouthConnections.Count = len(sw.Spec.State.SouthConnections.Connections)
	return nil
}

//defineSubnets process related resources to define switch's
//"south" subnets. Returns an error.
func (r *SwitchReconciler) defineSubnets(sw *switchv1alpha1.Switch, ctx context.Context) error {
	topLevelSwitch := r.getTopLevelSwitch(sw, ctx)
	if topLevelSwitch == nil {
		return nil
	}
	relatedAssignment := r.getRelatedAssignment(topLevelSwitch, ctx)
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
		cidr, sn, err := getSuitableSubnet(sw, subnetList, subnetv1alpha1.CIPv4SubnetType, regions, zones)
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
		cidr, sn, err := getSuitableSubnet(sw, subnetList, subnetv1alpha1.CIPv6SubnetType, regions, zones)
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

//getTopLevelSwitch recursively searches for top level switch
//related to the target switch.
func (r *SwitchReconciler) getTopLevelSwitch(sw *switchv1alpha1.Switch, ctx context.Context) *switchv1alpha1.Switch {
	if sw.Spec.State.NorthConnections != nil && sw.Spec.State.NorthConnections.Count != 0 {
		nextLevelSwitch := &switchv1alpha1.Switch{}
		if err := r.Client.Get(ctx, types.NamespacedName{
			Namespace: sw.Spec.State.NorthConnections.Connections[0].Namespace,
			Name:      sw.Spec.State.NorthConnections.Connections[0].Name,
		}, nextLevelSwitch); err != nil {
			return nil
		}
		return r.getTopLevelSwitch(nextLevelSwitch, ctx)
	}
	if sw.Spec.ConnectionLevel == 0 {
		return sw
	}
	return nil
}

//getRelatedAssignment searches for switch assignment resource
//related to target top level switch.
func (r *SwitchReconciler) getRelatedAssignment(sw *switchv1alpha1.Switch, ctx context.Context) *switchv1alpha1.SwitchAssignment {
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
		cidrToRelease, err := subnetv1alpha1.CIDRFromString(sw.Spec.SouthSubnetV4.CIDR)
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
		cidrToRelease, err := subnetv1alpha1.CIDRFromString(sw.Spec.SouthSubnetV6.CIDR)
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

//getSuitableSubnet finds the subnet resource, that fits address
//type, region, availability zones and addresses count. It returns
//pointers to the CIDR and subnet resource objects or an error.
func getSuitableSubnet(
	sw *switchv1alpha1.Switch,
	subnetList *subnetv1alpha1.SubnetList,
	addressType subnetv1alpha1.SubnetAddressType,
	regions []string,
	zones []string) (*subnetv1alpha1.CIDR, *subnetv1alpha1.Subnet, error) {

	addressesNeeded := sw.GetAddressNeededCount(addressType)
	for _, sn := range subnetList.Items {
		if sn.Spec.NetworkName == "underlay" &&
			sn.Status.Type == addressType &&
			reflect.DeepEqual(sn.Spec.Regions, regions) &&
			reflect.DeepEqual(sn.Spec.AvailabilityZones, zones) {
			addressesLeft := sn.Status.CapacityLeft
			if sn.Status.Type == addressType && addressesLeft.CmpInt64(addressesNeeded) >= 0 {
				minVacantCIDR := switchv1alpha1.GetMinimalVacantCIDR(sn.Status.Vacant, addressType, addressesNeeded)
				mask := sw.GetNeededMask(addressType, float64(addressesNeeded))
				addr := minVacantCIDR.Net.IP
				network := &net.IPNet{
					IP:   addr,
					Mask: mask,
				}
				cidrCandidate := &subnetv1alpha1.CIDR{Net: network}
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

//updateNorthInterfacesAddresses sets IP addresses for switch
//north interfaces according to north switch subnets values.
func (r *SwitchReconciler) updateNorthInterfacesAddresses(sw *switchv1alpha1.Switch, ctx context.Context) {
	if sw.Spec.State.NorthConnections.Count > 0 {
		for _, item := range sw.Spec.State.NorthConnections.Connections {
			northSwitch := &switchv1alpha1.Switch{}
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
			}
			for _, iface := range sw.Spec.Interfaces {
				if iface.LLDPChassisID == northSwitch.Spec.SwitchChassis.ChassisID {
					peerIface := northSwitch.Spec.Interfaces[iface.LLDPPortDescription]
					ipv4Addr := peerIface.RequestAddress(subnetv1alpha1.CIPv4SubnetType)
					ipv6Addr := peerIface.RequestAddress(subnetv1alpha1.CIPv6SubnetType)
					if ipv4Addr != nil {
						iface.IPv4 = fmt.Sprintf("%s/%d", ipv4Addr.String(), switchv1alpha1.CIPv4InterfaceSubnetMask)
					}
					if ipv6Addr != nil {
						iface.IPv6 = fmt.Sprintf("%s/%d", ipv6Addr.String(), switchv1alpha1.CIPv6InterfaceSubnetMask)
					}
				}
			}
		}
	}
}

func (r *SwitchReconciler) fillSouthConnections(obj *switchv1alpha1.Switch, ctx context.Context) error {
	southSwitchList, err := r.findSouthNeighboursSwitches(obj, ctx)
	if err != nil {
		return err
	}
	neighboursMap := southSwitchList.ConstructNeighboursFromSwitchList()
	connections := obj.Spec.State.SouthConnections.Connections
	for i, neighbour := range connections {
		if _, ok := neighboursMap[neighbour.ChassisID]; ok {
			connections[i] = neighboursMap[neighbour.ChassisID]
		}
	}
	obj.Spec.State.SouthConnections.Connections = connections
	return nil
}

func (r *SwitchReconciler) NeighboursOk(sw *switchv1alpha1.Switch, ctx context.Context) bool {
	list := &switchv1alpha1.SwitchList{}
	if err := r.List(ctx, list); err != nil {
		return false
	}
	for _, item := range list.Items {
		for _, conn := range sw.Spec.State.NorthConnections.Connections {
			if conn.ChassisID == item.Spec.SwitchChassis.ChassisID && item.Spec.ConnectionLevel != sw.Spec.ConnectionLevel-1 {
				return false
			}
		}
		for _, conn := range sw.Spec.State.SouthConnections.Connections {
			if conn.ChassisID == item.Spec.SwitchChassis.ChassisID && item.Spec.ConnectionLevel != sw.Spec.ConnectionLevel+1 {
				return false
			}
		}
	}
	return true
}
