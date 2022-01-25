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
	"time"

	gocidr "github.com/apparentlymart/go-cidr/cidr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"
	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

const (
	CUnderlayNetwork            = "underlay"
	CLoopbackParentSubnetPrefix = "switches"
	CSouthParentSubnetPrefix    = "switch-ranges"
	CIPAMv4Suffix               = "v4"
	CIPAMv6Suffix               = "v6"

	CLabelLoopbackRel   = "loopback"
	CLabelSwitchPortRel = "switch-port"

	CSwitchFinalizer = "switches.switch.onmetal.de/finalizer"
)

// SwitchReconciler reconciles a Switch object
type SwitchReconciler struct {
	client.Client
	Log        logr.Logger
	Scheme     *runtime.Scheme
	Background *background
}

type background struct {
	switches   *switchv1alpha1.SwitchList
	inventory  *inventoriesv1alpha1.Inventory
	assignment *switchv1alpha1.SwitchAssignment
	loopbacks  []net.IP
	ipv4Used   bool
	ipv6Used   bool
}

type labelsMap struct {
	include map[string][]string
	exclude map[string][]string
}

// FIXME: commented temporary
// type subnetNICdata struct {
// 	nic    string
// 	af     ipamv1alpha1.SubnetAddressType
// 	subnet *ipamv1alpha1.Subnet
// }

// FIXME: commented temporary
// type ipNICdata struct {
// 	nic string
// 	af  ipamv1alpha1.SubnetAddressType
// 	ip  *ipamv1alpha1.IP
// }

//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/finalizers,verbs=update
//+kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets,verbs=get;create;list;watch;delete;deletecollection
//+kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets/status,verbs=get;update
//+kubebuilder:rbac:groups=ipam.onmetal.de,resources=ips,verbs=get;create;list;watch;delete;deletecollection
//+kubebuilder:rbac:groups=ipam.onmetal.de,resources=ips/status,verbs=get;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *SwitchReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	result = ctrl.Result{}
	obj := &switchv1alpha1.Switch{}
	if err = r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !obj.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(obj, CSwitchFinalizer) {
			return
		}

		if err = r.finalize(ctx, obj); err != nil {
			r.Log.Error(err, "failed to finalize resource",
				"gvk", obj.GroupVersionKind().String(),
				"name", obj.NamespacedName())
			return
		}
		controllerutil.RemoveFinalizer(obj, CSwitchFinalizer)
		err = r.Update(ctx, obj)
		if err != nil {
			r.Log.Error(err, "failed to update resource spec on finalizer removal",
				"gvk", obj.GroupVersionKind().String(),
				"name", obj.NamespacedName())
		}
		return
	}

	if !controllerutil.ContainsFinalizer(obj, CSwitchFinalizer) {
		controllerutil.AddFinalizer(obj, CSwitchFinalizer)
		err = r.Update(ctx, obj)
		if err != nil {
			r.Log.Error(err, "failed to update resource spec",
				"gvk", obj.GroupVersionKind().String(),
				"name", obj.NamespacedName())
		}
		return
	}

	if err = r.getBackground(ctx, obj); err != nil {
		r.Log.Error(err, "failed to get background for requested resource",
			"gvk", obj.GroupVersionKind().String(),
			"name", obj.NamespacedName())
		return ctrl.Result{RequeueAfter: CSwitchRequeueInterval}, nil
	}
	stateMachine := r.prepareStateMachine()
	result, err = stateMachine.launch(ctx, obj)
	r.Background.assignment = &switchv1alpha1.SwitchAssignment{}
	r.Background.inventory = &inventoriesv1alpha1.Inventory{}
	r.Background.switches = &switchv1alpha1.SwitchList{}
	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1alpha1.Switch{}).
		Complete(r)
}

func (r *SwitchReconciler) getBackground(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
	if r.Background == nil {
		r.Background = &background{
			assignment: &switchv1alpha1.SwitchAssignment{},
			inventory:  &inventoriesv1alpha1.Inventory{},
			switches:   &switchv1alpha1.SwitchList{},
			loopbacks:  make([]net.IP, 0),
		}
	}
	if err = r.Get(ctx, obj.NamespacedName(), r.Background.inventory); err != nil {
		return
	}
	if err = r.List(ctx, r.Background.switches); err != nil {
		return
	}
	r.Background.assignment, err = r.findAssignment(ctx, obj)
	if err != nil {
		return
	}

	subnet := &ipamv1alpha1.Subnet{}
	err = r.Get(ctx, types.NamespacedName{
		Namespace: obj.Namespace,
		Name:      fmt.Sprintf("%s-%s", CSouthParentSubnetPrefix, CIPAMv4Suffix),
	}, subnet)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return
		}
		r.Background.ipv4Used = false
		err = nil
	} else {
		r.Background.ipv4Used = true
	}
	err = r.Get(ctx, types.NamespacedName{
		Namespace: obj.Namespace,
		Name:      fmt.Sprintf("%s-%s", CSouthParentSubnetPrefix, CIPAMv6Suffix),
	}, subnet)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return
		}
		r.Background.ipv6Used = false
		err = nil
	} else {
		r.Background.ipv6Used = true
	}
	return
}

func (r *SwitchReconciler) findAssignment(ctx context.Context, obj *switchv1alpha1.Switch) (*switchv1alpha1.SwitchAssignment, error) {
	lbl := labelsMap{include: map[string][]string{switchv1alpha1.LabelChassisId: {switchv1alpha1.MacToLabel(obj.Spec.Chassis.ChassisID)}}}
	opts, err := getListFilter(lbl)
	if err != nil {
		r.Log.Error(err, "failed to construct list options")
		return nil, err
	}
	list := &switchv1alpha1.SwitchAssignmentList{}
	if err := r.List(ctx, list, opts); err != nil {
		r.Log.Error(err, "failed to list resources", "gvk", list.GroupVersionKind().String())
		return nil, err
	}
	if len(list.Items) == 0 {
		return nil, nil
	}
	return &list.Items[0], nil
}

func getListFilter(labelsMap labelsMap) (*client.ListOptions, error) {
	selector := labels.NewSelector()
	if labelsMap.include != nil {
		for k, v := range labelsMap.include {
			if len(v) == 0 {
				continue
			}
			req, err := getLabelsRequirementIncluded(k, v)
			if err != nil {
				return nil, err
			}
			selector = selector.Add(*req)
		}
	}
	if labelsMap.exclude != nil {
		for k, v := range labelsMap.exclude {
			if len(v) == 0 {
				continue
			}
			req, err := getLabelsRequirementExcluded(k, v)
			if err != nil {
				return nil, err
			}
			selector = selector.Add(*req)
		}
	}
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         100,
	}
	return opts, nil
}

func getLabelsRequirementIncluded(label string, values []string) (*labels.Requirement, error) {
	labelsReq, err := labels.NewRequirement(label, selection.In, values)
	if err != nil {
		return nil, err
	}
	return labelsReq, nil
}

func getLabelsRequirementExcluded(label string, values []string) (*labels.Requirement, error) {
	labelsReq, err := labels.NewRequirement(label, selection.NotIn, values)
	if err != nil {
		return nil, err
	}
	return labelsReq, nil
}

func (r *SwitchReconciler) prepareStateMachine() *stateMachine {
	updateConfState := newStep(r.configUnmanaged, r.setConfigState, r.updateResStatus, nil)
	updateStatus := newStep(r.stateReadyOk, r.completeProcessing, r.updateResStatus, updateConfState)
	// FIXME: commented temporary
	// updateIPAMResources := newStep(r.ipamResOk, r.createNetworkResources, r.updateResStatus, updateStatus)
	// updateNICsAddresses := newStep(r.nicsAddressesOk, r.updateNICsAddresses, r.updateResStatus, updateIPAMResources)
	updateNICsAddresses := newStep(r.nicsAddressesOk, r.updateNICsAddresses, r.updateResStatus, updateStatus)
	updateLoopbackAddresses := newStep(r.loopbackAddressesOk, r.updateLoopbacks, r.updateResStatus, updateNICsAddresses)
	updateSubnets := newStep(r.subnetsOk, r.setSubnets, r.updateResStatus, updateLoopbackAddresses)
	updateConnectionLevel := newStep(r.connectionLevelOk, r.updateConnectionLevel, r.updateResStatus, updateSubnets)
	updateRole := newStep(r.roleOk, r.setRole, r.updateResStatus, updateConnectionLevel)
	updatePeers := newStep(r.peersInfoOk, r.fillPeersInfo, r.updateResStatus, updateRole)
	updateInterfaces := newStep(r.interfacesOk, r.setInterfaces, r.updateResStatus, updatePeers)
	updateInitialStatus := newStep(r.stateOk, r.setInitialStatus, r.updateResStatus, updateInterfaces)
	return newStateMachine(updateInitialStatus)
}

func (r *SwitchReconciler) updateResStatus(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
	if err = r.Status().Update(ctx, obj); err != nil {
		r.Log.Error(err, "failed to update resource status", "gvk", obj.GroupVersionKind().String(), "name", obj.NamespacedName())
	}
	return
}

func (r *SwitchReconciler) finalize(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
	assignment, err := r.findAssignment(ctx, obj)
	if err != nil {
		r.Log.Error(err, "failed to get related assignment for resource",
			"gvk", obj.GroupVersionKind().String(),
			"name", obj.NamespacedName())
	}
	if assignment != nil {
		assignment.FillStatus(switchv1alpha1.CAssignmentStatePending, &switchv1alpha1.LinkedSwitchSpec{})
		err = r.Status().Update(ctx, assignment)
		if err != nil {
			r.Log.Error(err, "failed to update resource status",
				"gvk", assignment.GroupVersionKind().String(),
				"name", assignment.NamespacedName())
		}
	}

	// here we'll filter only IPs and subnets related to switch ports.
	// Subnets related to switch and switch loopback IPs will be left as is.
	lbl := labelsMap{
		include: map[string][]string{
			switchv1alpha1.LabelSwitchName:       {obj.Name},
			switchv1alpha1.LabelResourceRelation: {CLabelSwitchPortRel},
		},
		exclude: map[string][]string{switchv1alpha1.LabelResourceRelation: {CLabelLoopbackRel}},
	}
	filter, err := getListFilter(lbl)
	if err != nil {
		r.Log.Error(err, "failed to build list filter for finalizer")
	}
	err = r.DeleteAllOf(ctx, &ipamv1alpha1.IP{}, client.InNamespace(obj.Namespace), client.MatchingLabelsSelector{Selector: filter.LabelSelector})
	if err != nil {
		r.Log.Error(err, "failed to delete related ip resources")
	}
	err = r.DeleteAllOf(ctx, &ipamv1alpha1.Subnet{}, client.InNamespace(obj.Namespace), client.MatchingLabelsSelector{Selector: filter.LabelSelector})
	if err != nil {
		r.Log.Error(err, "failed to delete related subnet resources")
	}

	return
}

func (r *SwitchReconciler) stateOk(obj *switchv1alpha1.Switch) bool {
	return !obj.StateEqualTo(switchv1alpha1.CEmptyString)
}

func (r *SwitchReconciler) setInitialStatus(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
	obj.FillInitialStatus(r.Background.inventory, r.Background.switches)
	return
}

func (r *SwitchReconciler) interfacesOk(obj *switchv1alpha1.Switch) bool {
	return obj.InterfacesDataOk(r.Background.inventory, r.Background.switches)
}

func (r *SwitchReconciler) setInterfaces(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
	obj.SetSwitchState(switchv1alpha1.CSwitchStateInProgress)
	obj.SetConfState(switchv1alpha1.CSwitchConfInProgress)
	receivedInterfaces := switchv1alpha1.InterfacesFromInventory(r.Background.inventory.Spec.NICs, r.Background.switches)
	for iface, data := range obj.Status.Interfaces {
		receivedIface, ok := receivedInterfaces[iface]
		if !ok {
			delete(obj.Status.Interfaces, iface)
			continue
		}
		data.MACAddress = receivedIface.MACAddress
		data.FEC = receivedIface.FEC
		data.Lanes = receivedIface.Lanes
		data.State = receivedIface.State
		data.MTU = receivedIface.MTU
		data.Speed = receivedIface.Speed

		peerStored := data.Peer.ChassisID != switchv1alpha1.CEmptyString
		peerRemoved := receivedIface.Peer.ChassisID == switchv1alpha1.CEmptyString
		peerDisconnected := peerStored && peerRemoved
		if peerDisconnected {
			if data.IPv4.ResourceReference != nil {
				_ = r.removeRelatedIP(ctx, data.IPv4.ResourceReference.NamespacedName())
			}
			if data.IPv6.ResourceReference != nil {
				_ = r.removeRelatedIP(ctx, data.IPv6.ResourceReference.NamespacedName())
			}
		}
		if peerDisconnected && data.Direction == switchv1alpha1.CDirectionSouth {
			if obj.Status.SubnetV4.ResourceReference != nil {
				_ = r.removeRelatedSubnet(ctx, obj.Status.SubnetV4.ResourceReference.NamespacedName())
			}
			if obj.Status.SubnetV6.ResourceReference != nil {
				_ = r.removeRelatedSubnet(ctx, obj.Status.SubnetV6.ResourceReference.NamespacedName())
			}
		}
		data.Peer.ChassisID = receivedIface.Peer.ChassisID
		data.Peer.Type = receivedIface.Peer.Type
		delete(receivedInterfaces, iface)
	}
	for iface, data := range receivedInterfaces {
		obj.Status.Interfaces[iface] = data
	}
	return
}

func (r *SwitchReconciler) removeRelatedIP(ctx context.Context, namespacedName types.NamespacedName) (err error) {
	ip := &ipamv1alpha1.IP{}
	err = r.Get(ctx, namespacedName, ip)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return
	}
	err = r.Delete(ctx, ip)
	return
}

func (r *SwitchReconciler) removeRelatedSubnet(ctx context.Context, namespacedName types.NamespacedName) (err error) {
	subnet := &ipamv1alpha1.Subnet{}
	err = r.Get(ctx, namespacedName, subnet)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return
	}
	err = r.Delete(ctx, subnet)
	return
}

func (r *SwitchReconciler) peersInfoOk(obj *switchv1alpha1.Switch) bool {
	return obj.PeersDefined(r.Background.switches)
}

func (r *SwitchReconciler) fillPeersInfo(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
	obj.SetSwitchState(switchv1alpha1.CSwitchStateInProgress)
	obj.SetConfState(switchv1alpha1.CSwitchConfInProgress)
	obj.FillPeerSwitches(r.Background.switches)
	return
}

func (r *SwitchReconciler) roleOk(obj *switchv1alpha1.Switch) bool {
	return obj.RoleMatchPeers()
}

func (r *SwitchReconciler) setRole(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
	if obj.Status.Role == switchv1alpha1.CSwitchRoleLeaf {
		obj.Status.Role = switchv1alpha1.CSwitchRoleSpine
	}
	if obj.Status.Role == switchv1alpha1.CSwitchRoleSpine {
		obj.Status.Role = switchv1alpha1.CSwitchRoleLeaf
	}
	return
}

func (r *SwitchReconciler) connectionLevelOk(obj *switchv1alpha1.Switch) (result bool) {
	switch {
	case r.Background.assignment != nil:
		connectionLevelOk := obj.Status.ConnectionLevel == 0
		assignmentFinished := r.Background.assignment.Status.State == switchv1alpha1.CAssignmentStateFinished
		result = connectionLevelOk && assignmentFinished
	case r.Background.assignment == nil:
		connectionLevelOk := obj.Status.ConnectionLevel != 0 && obj.Status.ConnectionLevel != 255
		matchPeers := obj.ConnectionLevelMatchPeers(r.Background.switches)
		result = connectionLevelOk && matchPeers
	}
	return
}

func (r *SwitchReconciler) updateConnectionLevel(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
	obj.SetSwitchState(switchv1alpha1.CSwitchStateInProgress)
	obj.SetConfState(switchv1alpha1.CSwitchConfInProgress)
	if r.Background.assignment != nil && obj.Status.ConnectionLevel != 0 {
		obj.Status.ConnectionLevel = 0
		err = r.processAssignment(ctx, obj)
	}
	obj.ComputeConnectionLevel(r.Background.switches)
	return
}

func (r *SwitchReconciler) processAssignment(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
	r.Background.assignment.FillStatus(switchv1alpha1.CAssignmentStateFinished, &switchv1alpha1.LinkedSwitchSpec{
		Name:      obj.Name,
		Namespace: obj.Namespace,
	})
	if err = r.Status().Update(ctx, r.Background.assignment); err != nil {
		r.Log.Error(err, "failed to update resource",
			"gvk", r.Background.assignment.GroupVersionKind().String(),
			"name", r.Background.assignment.NamespacedName())
		return
	}
	return
}

func (r *SwitchReconciler) subnetsOk(obj *switchv1alpha1.Switch) bool {
	return obj.SubnetsDefined(r.Background.ipv4Used, r.Background.ipv6Used)
}

func (r *SwitchReconciler) setSubnets(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
	obj.SetSwitchState(switchv1alpha1.CSwitchStateInProgress)
	obj.SetConfState(switchv1alpha1.CSwitchConfInProgress)
	if r.Background.ipv4Used {
		if err = r.setV4Subnet(ctx, obj); err != nil {
			r.Log.Error(err, "failed to setup switch V4 subnet")
			return
		}
	}
	if r.Background.ipv6Used {
		if err = r.setV6Subnet(ctx, obj); err != nil {
			r.Log.Error(err, "failed to setup switch V6 subnet")
			return
		}
	}
	return
}

func (r *SwitchReconciler) setV4Subnet(ctx context.Context, obj *switchv1alpha1.Switch) error {
	subnet := &ipamv1alpha1.Subnet{}
	err := r.Get(ctx, types.NamespacedName{Namespace: obj.Namespace, Name: obj.SwitchSubnetName(ipamv1alpha1.CIPv4SubnetType)}, subnet)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		region, err := r.getRegion(ctx, obj)
		if err != nil {
			r.Log.Error(err, "failed to get regions for subnet search")
			return err
		}
		addressesCount := obj.GetAddressCount(ipamv1alpha1.CIPv4SubnetType)
		spec := ipamv1alpha1.SubnetSpec{
			Capacity:    resource.NewQuantity(addressesCount, resource.DecimalSI),
			NetworkName: CUnderlayNetwork,
			Regions:     region.ConvertToSubnetRegion(),
		}
		subnetName := obj.SwitchSubnetName(ipamv1alpha1.CIPv4SubnetType)
		lbl := map[string]string{switchv1alpha1.LabelSwitchName: obj.Name}
		subnet, err = r.createSubnet(ctx, obj, CSouthParentSubnetPrefix, subnetName, ipamv1alpha1.CIPv4SubnetType, spec, lbl)
		if err != nil {
			return err
		}
	}
	if subnet.Status.State == ipamv1alpha1.CFinishedSubnetState {
		obj.Status.SubnetV4.CIDR = subnet.Status.Reserved.String()
		obj.Status.SubnetV4.Region = switchv1alpha1.ConvertFromSubnetRegion(subnet.Spec.Regions)
		obj.Status.SubnetV4.ResourceReference = &switchv1alpha1.ResourceReferenceSpec{
			APIVersion: subnet.APIVersion,
			Kind:       subnet.Kind,
			Namespace:  subnet.Namespace,
			Name:       subnet.Name,
		}
	}
	return nil
}

func (r *SwitchReconciler) setV6Subnet(ctx context.Context, obj *switchv1alpha1.Switch) error {
	subnet := &ipamv1alpha1.Subnet{}
	err := r.Get(ctx, types.NamespacedName{Namespace: obj.Namespace, Name: obj.SwitchSubnetName(ipamv1alpha1.CIPv6SubnetType)}, subnet)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		region, err := r.getRegion(ctx, obj)
		if err != nil {
			r.Log.Error(err, "failed to get regions for subnet search")
			return err
		}
		addressesCount := obj.GetAddressCount(ipamv1alpha1.CIPv6SubnetType)
		spec := ipamv1alpha1.SubnetSpec{
			Capacity:    resource.NewQuantity(addressesCount, resource.DecimalSI),
			NetworkName: CUnderlayNetwork,
			Regions:     region.ConvertToSubnetRegion(),
		}
		subnetName := obj.SwitchSubnetName(ipamv1alpha1.CIPv6SubnetType)
		lbl := map[string]string{switchv1alpha1.LabelSwitchName: obj.Name}
		subnet, err = r.createSubnet(ctx, obj, CSouthParentSubnetPrefix, subnetName, ipamv1alpha1.CIPv6SubnetType, spec, lbl)
		if err != nil {
			return err
		}
	}
	if subnet.Status.State == ipamv1alpha1.CFinishedSubnetState {
		obj.Status.SubnetV6.CIDR = subnet.Status.Reserved.String()
		obj.Status.SubnetV6.Region = switchv1alpha1.ConvertFromSubnetRegion(subnet.Spec.Regions)
		obj.Status.SubnetV6.ResourceReference = &switchv1alpha1.ResourceReferenceSpec{
			APIVersion: subnet.APIVersion,
			Kind:       subnet.Kind,
			Namespace:  subnet.Namespace,
			Name:       subnet.Name,
		}
	}
	return nil
}

func (r *SwitchReconciler) getRegion(ctx context.Context, obj *switchv1alpha1.Switch) (region *switchv1alpha1.RegionSpec, err error) {
	topLevelSpine, err := r.getTopLevelSpine(ctx, obj)
	if err != nil {
		return
	}
	assignment, err := r.findAssignment(ctx, topLevelSpine)
	if err != nil {
		return
	}
	if assignment == nil {
		err = fmt.Errorf("failed to get region")
		return
	}
	region = assignment.Spec.Region
	return
}

func (r *SwitchReconciler) getTopLevelSpine(ctx context.Context, obj *switchv1alpha1.Switch) (*switchv1alpha1.Switch, error) {
	if obj.Status.ConnectionLevel == 0 {
		return obj, nil
	}
	upstreamSwitch := &switchv1alpha1.Switch{}
	key := func(nics map[string]*switchv1alpha1.InterfaceSpec) (result types.NamespacedName) {
		for _, nicData := range nics {
			if nicData.Direction == switchv1alpha1.CDirectionNorth {
				result = nicData.Peer.ResourceReference.NamespacedName()
				break
			}
		}
		return
	}(obj.Status.Interfaces)
	if err := r.Get(ctx, key, upstreamSwitch); err != nil {
		r.Log.Error(err, "failed to get top-level switch resource")
		return nil, err
	}
	return r.getTopLevelSpine(ctx, upstreamSwitch)
}

func (r *SwitchReconciler) findParentSubnet(
	prefix string,
	af ipamv1alpha1.SubnetAddressType,
	list *ipamv1alpha1.SubnetList) (*ipamv1alpha1.Subnet, error) {
	var subnet *ipamv1alpha1.Subnet
	suffix := CIPAMv4Suffix
	if af == ipamv1alpha1.CIPv6SubnetType {
		suffix = CIPAMv6Suffix
	}
	name := fmt.Sprintf("%s-%s", prefix, suffix)
	err := fmt.Errorf("failed to find parent subnet %s", name)
	for _, item := range list.Items {
		if item.Name == name {
			subnet = &item
			return subnet, nil
		}
	}
	return nil, err
}

func (r *SwitchReconciler) createSubnet(
	ctx context.Context,
	obj *switchv1alpha1.Switch,
	prefix string,
	subnetName string,
	af ipamv1alpha1.SubnetAddressType,
	spec ipamv1alpha1.SubnetSpec,
	labels map[string]string) (subnet *ipamv1alpha1.Subnet, err error) {
	subnets := &ipamv1alpha1.SubnetList{}
	if err = r.List(ctx, subnets); err != nil {
		r.Log.Error(err, "failed to list resources", "gvk", subnets.GroupVersionKind().String())
		return
	}
	parentSubnet, err := r.findParentSubnet(prefix, af, subnets)
	if err != nil {
		return
	}
	if parentSubnet == nil {
		return
	}
	subnet = &ipamv1alpha1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: obj.Namespace,
			Name:      subnetName,
			Labels:    labels,
		},
		Spec: spec,
	}
	subnet.Spec.ParentSubnetName = parentSubnet.Name
	if err = r.Create(ctx, subnet); err != nil {
		r.Log.Error(err, "failed to create resource",
			"gvk", subnet.GroupVersionKind().String(),
			"name", types.NamespacedName{Namespace: parentSubnet.Namespace, Name: subnetName})
		return
	}
	return
}

func (r *SwitchReconciler) loopbackAddressesOk(obj *switchv1alpha1.Switch) bool {
	return obj.LoopbackAddressesDefined(r.Background.ipv4Used, r.Background.ipv6Used)
}

func (r *SwitchReconciler) updateLoopbacks(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
	obj.SetSwitchState(switchv1alpha1.CSwitchStateInProgress)
	obj.SetConfState(switchv1alpha1.CSwitchConfInProgress)
	subnets := &ipamv1alpha1.SubnetList{}
	if err = r.List(ctx, subnets); err != nil {
		r.Log.Error(err, "failed to list resources", "gvk", subnets.GroupVersionKind().String())
		return
	}
	if r.Background.ipv4Used && obj.Status.LoopbackV4.Address == switchv1alpha1.CEmptyString {
		err = r.setLoopbackIP(ctx, obj, subnets, ipamv1alpha1.CIPv4SubnetType)
		if err != nil {
			return
		}
	}
	if r.Background.ipv6Used && obj.Status.LoopbackV6.Address == switchv1alpha1.CEmptyString {
		err = r.setLoopbackIP(ctx, obj, subnets, ipamv1alpha1.CIPv6SubnetType)
		if err != nil {
			return
		}
	}
	return
}

func (r *SwitchReconciler) setLoopbackIP(
	ctx context.Context,
	obj *switchv1alpha1.Switch,
	subnets *ipamv1alpha1.SubnetList,
	af ipamv1alpha1.SubnetAddressType) error {
	loopbackIP := &ipamv1alpha1.IP{}
	err := r.Get(ctx, types.NamespacedName{Namespace: obj.Namespace, Name: obj.LoopbackIPResourceName(af)}, loopbackIP)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		parentSubnet, err := r.findParentSubnet(CLoopbackParentSubnetPrefix, af, subnets)
		if err != nil {
			return err
		}
		lbl := map[string]string{
			switchv1alpha1.LabelSwitchName:       obj.Name,
			switchv1alpha1.LabelResourceRelation: CLabelLoopbackRel,
		}
		loopbackIP, err = r.createIP(ctx, obj, parentSubnet, obj.LoopbackIPResourceName(af), af, 0, true, lbl)
		if err != nil {
			return err
		}
	}
	ref := &switchv1alpha1.ResourceReferenceSpec{
		APIVersion: loopbackIP.APIVersion,
		Kind:       loopbackIP.Kind,
		Namespace:  loopbackIP.Namespace,
		Name:       loopbackIP.Name,
	}
	switch af {
	case ipamv1alpha1.CIPv4SubnetType:
		obj.Status.LoopbackV4.Address = loopbackIP.Spec.IP.String()
		obj.Status.LoopbackV4.ResourceReference = ref
	case ipamv1alpha1.CIPv6SubnetType:
		obj.Status.LoopbackV6.Address = loopbackIP.Spec.IP.String()
		obj.Status.LoopbackV6.ResourceReference = ref
	}
	return nil
}

func (r *SwitchReconciler) createIP(
	ctx context.Context,
	obj *switchv1alpha1.Switch,
	parentSubnet *ipamv1alpha1.Subnet,
	name string,
	af ipamv1alpha1.SubnetAddressType,
	offset int,
	loopback bool,
	labels map[string]string) (ip *ipamv1alpha1.IP, err error) {
	address, err := r.getIPAddress(obj, parentSubnet, af, offset, loopback)
	if err != nil {
		return
	}
	ipAddr, _ := ipamv1alpha1.IPAddrFromString(address.String())
	ip = &ipamv1alpha1.IP{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: obj.Namespace,
			Name:      name,
			Labels:    labels,
		},
		Spec: ipamv1alpha1.IPSpec{
			SubnetName: parentSubnet.Name,
			ResourceReference: &ipamv1alpha1.ResourceReference{
				APIVersion: obj.APIVersion,
				Kind:       obj.Kind,
				Name:       obj.Name,
			},
			IP: ipAddr,
		},
	}
	err = r.Create(ctx, ip)
	return
}

func (r *SwitchReconciler) getIPAddress(
	obj *switchv1alpha1.Switch,
	subnet *ipamv1alpha1.Subnet,
	af ipamv1alpha1.SubnetAddressType,
	offset int,
	loopback bool) (ip net.IP, err error) {
	if !loopback && subnet.Spec.CIDR != nil {
		ip, _ = gocidr.Host(subnet.Spec.CIDR.Net, offset)
		return
	}
	ip, _ = gocidr.Host(subnet.Status.Reserved.Net, offset)
	if !loopback {
		return
	}
	if !r.ipAddressUsed(ip) {
		r.Background.loopbacks = append(r.Background.loopbacks, ip)
		return
	}
	capLeft, _ := subnet.Status.CapacityLeft.AsInt64()
	if int64(offset+1) > capLeft {
		err = fmt.Errorf("no free addresses left in subnet %s", subnet.Name)
		return
	}
	return r.getIPAddress(obj, subnet, af, offset+1, loopback)
}

func (r *SwitchReconciler) ipAddressUsed(address net.IP) bool {
	for _, ip := range r.Background.loopbacks {
		if reflect.DeepEqual(ip, address) {
			return true
		}
	}
	return false
}

func (r *SwitchReconciler) nicsAddressesOk(obj *switchv1alpha1.Switch) bool {
	return obj.NICsAddressesDefined(r.Background.ipv4Used, r.Background.ipv6Used, r.Background.switches)
}

func (r *SwitchReconciler) updateNICsAddresses(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
	obj.SetSwitchState(switchv1alpha1.CSwitchStateInProgress)
	obj.SetConfState(switchv1alpha1.CSwitchConfInProgress)
	if err = obj.UpdateSouthNICsIP(r.Background.ipv4Used, r.Background.ipv6Used); err != nil {
		r.Log.Error(err, "failed to setup south NICs subnets")
		return
	}
	for _, nicData := range obj.Status.Interfaces {
		if nicData.Direction == switchv1alpha1.CDirectionSouth {
			continue
		}
		for _, item := range r.Background.switches.Items {
			if item.NamespacedName() != nicData.Peer.ResourceReference.NamespacedName() {
				continue
			}
			peerV4SubnetOk := false
			if !r.Background.ipv4Used {
				peerV4SubnetOk = true
			}
			if r.Background.ipv4Used && item.Status.SubnetV4.CIDR != switchv1alpha1.CEmptyString {
				peerV4SubnetOk = true
			}
			peerV6SubnetOk := false
			if !r.Background.ipv6Used {
				peerV6SubnetOk = true
			}
			if r.Background.ipv6Used && item.Status.SubnetV6.CIDR != switchv1alpha1.CEmptyString {
				peerV6SubnetOk = true
			}
			if !(peerV4SubnetOk && peerV6SubnetOk) {
				err = fmt.Errorf("peer's subnets not defined")
				return
			}
		}
	}
	err = obj.UpdateNorthNICsIP(r.Background.ipv4Used, r.Background.ipv6Used, r.Background.switches)
	return
}

// FIXME: commented temporary
// func (r *SwitchReconciler) ipamResOk(obj *switchv1alpha1.Switch) bool {
// 	return obj.IPAMResourcesCreated()
// }

// FIXME: commented temporary
// func (r *SwitchReconciler) createNetworkResources(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
// 	subnets := r.createNICsSubnets(ctx, obj)
// 	ips := r.createNICsIPs(ctx, obj, subnets)
// 	for data := range ips {
// 		nicData := obj.Status.Interfaces[data.nic]
// 		ref := &switchv1alpha1.ResourceReferenceSpec{
// 			APIVersion: data.ip.APIVersion,
// 			Kind:       data.ip.Kind,
// 			Namespace:  data.ip.Namespace,
// 			Name:       data.ip.Name,
// 		}
// 		switch data.af {
// 		case ipamv1alpha1.CIPv4SubnetType:
// 			nicData.IPv4.ResourceReference = ref
// 		case ipamv1alpha1.CIPv6SubnetType:
// 			nicData.IPv6.ResourceReference = ref
// 		}
// 	}
// 	return
// }

func (r *SwitchReconciler) stateReadyOk(obj *switchv1alpha1.Switch) bool {
	//return r.connectionLevelOk(obj) && obj.IPAMResourcesCreated() && obj.Status.State == switchv1alpha1.CSwitchStateReady
	return r.connectionLevelOk(obj) && obj.Status.State == switchv1alpha1.CSwitchStateReady
}

func (r *SwitchReconciler) completeProcessing(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
	obj.SetSwitchState(switchv1alpha1.CSwitchStateReady)
	obj.SetConfState(switchv1alpha1.CSwitchConfPending)
	return
}

func (r *SwitchReconciler) configUnmanaged(obj *switchv1alpha1.Switch) bool {
	if !obj.Status.Configuration.Managed {
		return true
	}
	if obj.Status.Configuration.LastCheck == switchv1alpha1.CEmptyString {
		return true
	}
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	lastCheck, _ := time.Parse(time.UnixDate, obj.Status.Configuration.LastCheck)
	return now.Sub(lastCheck) < time.Second*10
}

func (r *SwitchReconciler) setConfigState(ctx context.Context, obj *switchv1alpha1.Switch) (err error) {
	obj.Status.Configuration.ManagerState = switchv1alpha1.CConfManagerSFailed
	obj.Status.Configuration.State = switchv1alpha1.CSwitchConfPending
	return
}

// FIXME: commented temporary
// func (r *SwitchReconciler) createNICsSubnets(ctx context.Context, obj *switchv1alpha1.Switch) <-chan *subnetNICdata {
// 	resChan := make(chan *subnetNICdata, 10)

// 	go func() {
// 		defer func() {
// 			close(resChan)
// 		}()

// 		for nic, nicData := range obj.Status.Interfaces {
// 			if nicData.Direction == switchv1alpha1.CDirectionNorth {
// 				continue
// 			}
// 			cidrV4, _ := ipamv1alpha1.CIDRFromString(nicData.IPv4.Address)
// 			specV4 := ipamv1alpha1.SubnetSpec{
// 				CIDR:        cidrV4,
// 				NetworkName: CUnderlayNetwork,
// 				Regions:     obj.Status.SubnetV4.Region.ConvertToSubnetRegion(),
// 			}
// 			subnetNameV4 := obj.InterfaceSubnetName(nic, ipamv1alpha1.CIPv4SubnetType)
// 			labels := map[string]string{
// 				switchv1alpha1.LabelSwitchName:       obj.Name,
// 				switchv1alpha1.LabelResourceRelation: CLabelSwitchPortRel,
// 			}
// 			subnetV4, err := r.createSubnet(ctx, obj, obj.Name, subnetNameV4, ipamv1alpha1.CIPv4SubnetType, obj.Status.SubnetV4.Region, specV4, labels)
// 			if err != nil && !apierrors.IsAlreadyExists(err) {
// 				return
// 			}
// 			resChan <- &subnetNICdata{nic, ipamv1alpha1.CIPv4SubnetType, subnetV4}

// 			cidrV6, _ := ipamv1alpha1.CIDRFromString(nicData.IPv6.Address)
// 			specV6 := ipamv1alpha1.SubnetSpec{
// 				CIDR:        cidrV6,
// 				NetworkName: CUnderlayNetwork,
// 				Regions:     obj.Status.SubnetV6.Region.ConvertToSubnetRegion(),
// 			}
// 			subnetNameV6 := obj.InterfaceSubnetName(nic, ipamv1alpha1.CIPv6SubnetType)
// 			subnetV6, err := r.createSubnet(ctx, obj, obj.Name, subnetNameV6, ipamv1alpha1.CIPv6SubnetType, obj.Status.SubnetV6.Region, specV6, labels)
// 			if err != nil && !apierrors.IsAlreadyExists(err) {
// 				return
// 			}
// 			resChan <- &subnetNICdata{nic, ipamv1alpha1.CIPv6SubnetType, subnetV6}
// 		}
// 	}()
// 	return resChan
// }

// FIXME: commented temporary
// func (r *SwitchReconciler) createNICsIPs(ctx context.Context, obj *switchv1alpha1.Switch, subnets <-chan *subnetNICdata) <-chan *ipNICdata {
// 	resChan := make(chan *ipNICdata, 10)

// 	go func() {
// 		defer func() {
// 			close(resChan)
// 		}()

// 		for data := range subnets {
// 			if data.subnet.Spec.CIDR == nil {
// 				continue
// 			}
// 			offset := 1
// 			if data.af == ipamv1alpha1.CIPv6SubnetType {
// 				offset = 0
// 			}
// 			labels := map[string]string{
// 				switchv1alpha1.LabelSwitchName:       obj.Name,
// 				switchv1alpha1.LabelResourceRelation: CLabelSwitchPortRel,
// 			}
// 			ip, err := r.createIP(ctx, obj, data.subnet, obj.InterfaceIPName(data.nic, data.af), data.af, offset, false, labels)
// 			if err != nil && !apierrors.IsAlreadyExists(err) {
// 				return
// 			}
// 			resChan <- &ipNICdata{data.nic, data.af, ip}
// 		}
// 	}()
// 	return resChan
// }
