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
	"strings"

	"github.com/go-logr/logr"
	subnetv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

type background struct {
	switches   *switchv1alpha1.SwitchList
	assignment *switchv1alpha1.SwitchAssignment
	inventory  *inventoriesv1alpha1.Inventory
	ctx        context.Context
}

type labelsMap struct {
	include map[string][]string
	exclude map[string][]string
}

const (
	CUnderlayNetwork        = "underlay"
	CSwitchLoopbackV4Subnet = "switches-v4"
	CSwitchLoopbackV6Subnet = "switches-v6"
)

// SwitchReconciler reconciles a Switch object
type SwitchReconciler struct {
	client.Client
	Log        logr.Logger
	Scheme     *runtime.Scheme
	Recorder   record.EventRecorder
	Background *background
}

//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=switches/finalizers,verbs=update
//+kubebuilder:rbac:groups=switch.onmetal.de,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets,verbs=get;create;list;watch;update;delete
//+kubebuilder:rbac:groups=ipam.onmetal.de,resources=subnets/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *SwitchReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	res := &switchv1alpha1.Switch{}
	if err := r.Get(ctx, req.NamespacedName, res); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.prepareBackground(ctx, res); err != nil {
		return ctrl.Result{}, err
	}

	switchStateMachine := r.prepareStateMachine(res)
	return switchStateMachine.launch(res)
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&switchv1alpha1.Switch{}).
		Complete(r)
}

func (r *SwitchReconciler) finalize(obj *switchv1alpha1.Switch) error {
	ctx := r.Background.ctx
	if !controllerutil.ContainsFinalizer(obj, switchv1alpha1.CSwitchFinalizer) {
		return nil
	}

	swa, err := r.findAssignment(ctx, obj)
	if err != nil {
		r.Log.Error(err, "failed to lookup for related switch assignment resource",
			"gvk", obj.GroupVersionKind(), "name", obj.NamespacedName())
	}
	if swa != nil {
		swa.FillStatus(switchv1alpha1.CAssignmentStatePending, &switchv1alpha1.LinkedSwitchSpec{})
		if err = r.Status().Update(ctx, swa); err != nil {
			r.Log.Error(err, "failed to set status on resource creation",
				"gvk", swa.GroupVersionKind(), "name", swa.NamespacedName())
		}
	}

	usedNICs := obj.GetUsedNICs()
	for i, nic := range usedNICs {
		usedNICs[i] = strings.ToLower(nic)
	}
	labels := labelsMap{
		include: map[string][]string{
			switchv1alpha1.LabelSwitchName:    {obj.Name},
			switchv1alpha1.LabelInterfaceName: usedNICs,
		},
	}
	opts, err := getListFilter(labels)
	if err != nil {
		return err
	}
	snl := &subnetv1alpha1.SubnetList{}
	err = r.List(ctx, snl, opts)
	if err != nil {
		return err
	}
	for _, sn := range snl.Items {
		if err := r.Delete(ctx, &sn); err != nil {
			return err
		}
	}

	controllerutil.RemoveFinalizer(obj, switchv1alpha1.CSwitchFinalizer)
	if err := r.Update(ctx, obj); err != nil {
		r.Log.Error(err, "failed to update resource on finalizer removal",
			"gvk", obj.GroupVersionKind(), "name", obj.NamespacedName())
		return err
	}

	return nil
}

func (r *SwitchReconciler) findAssignment(ctx context.Context, sw *switchv1alpha1.Switch) (*switchv1alpha1.SwitchAssignment, error) {
	labels := labelsMap{include: map[string][]string{switchv1alpha1.LabelChassisId: {switchv1alpha1.MacToLabel(sw.Spec.Chassis.ChassisID)}}}
	opts, err := getListFilter(labels)
	if err != nil {
		r.Log.Error(err, "failed to construct list options object")
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
			req, err := getLabelsRequirementIncluded(k, v)
			if err != nil {
				return nil, err
			}
			selector = selector.Add(*req)
		}
	}
	if labelsMap.exclude != nil {
		for k, v := range labelsMap.exclude {
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

func (r *SwitchReconciler) prepareBackground(ctx context.Context, sw *switchv1alpha1.Switch) error {
	if r.Background == nil {
		r.Background = &background{
			switches:   nil,
			assignment: nil,
			ctx:        ctx,
		}
	}

	list := &switchv1alpha1.SwitchList{}
	if err := r.List(ctx, list); err != nil {
		r.Log.Error(err, "failed to list resources", "gvk", list.GroupVersionKind().String())
		return err
	}
	r.Background.switches = list

	swa, err := r.findAssignment(ctx, sw)
	if err != nil {
		r.Log.Error(err, "failed to get related assignment resource",
			"gvk", sw.GroupVersionKind().String(),
			"name", sw.NamespacedName())
		return err
	}
	r.Background.assignment = swa

	invFound := false
	invList := &inventoriesv1alpha1.InventoryList{}
	err = r.List(ctx, invList)
	if err != nil {
		r.Log.Error(err, "failed to list resources", "gvk", invList.GroupVersionKind().String())
		return err
	}
	if len(invList.Items) == 0 {
		err = errors.New("empty inventories list")
		return err
	}
	inv := &inventoriesv1alpha1.Inventory{}
	for _, item := range invList.Items {
		if item.Name == sw.Name {
			inv = &item
			invFound = true
			break
		}
	}
	if !invFound {
		return errors.New("related inventory not found")
	}
	r.Background.inventory = inv
	return nil
}

func (r *SwitchReconciler) createSubnetWithCIDR(
	name string,
	namespace string,
	cidr *subnetv1alpha1.CIDR,
	parentSubnet string,
	regions []subnetv1alpha1.Region,
	labels map[string]string) (*subnetv1alpha1.Subnet, error) {
	sn := &subnetv1alpha1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: subnetv1alpha1.SubnetSpec{
			CIDR:             cidr,
			ParentSubnetName: parentSubnet,
			NetworkName:      CUnderlayNetwork,
			Regions:          regions,
		},
	}
	if err := r.Create(r.Background.ctx, sn); err != nil {
		return nil, err
	}
	return sn, nil
}

func (r *SwitchReconciler) createSubnetWithCapacity(
	name string,
	namespace string,
	capacity *resource.Quantity,
	parentSubnet string,
	regions []subnetv1alpha1.Region,
	labels map[string]string) (*subnetv1alpha1.Subnet, error) {
	sn := &subnetv1alpha1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: subnetv1alpha1.SubnetSpec{
			Capacity:         capacity,
			ParentSubnetName: parentSubnet,
			NetworkName:      CUnderlayNetwork,
			Regions:          regions,
		},
	}
	if err := r.Create(r.Background.ctx, sn); err != nil {
		return nil, err
	}
	return sn, nil
}

func (r *SwitchReconciler) removeInterfaceSubnet(name string, namespace string) error {
	infSubnet := &subnetv1alpha1.Subnet{}
	if err := r.Get(r.Background.ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, infSubnet); err == nil {
		if err := r.Delete(r.Background.ctx, infSubnet); err != nil {
			return err
		}
	}
	return nil
}

func (r *SwitchReconciler) prepareStateMachine(obj *switchv1alpha1.Switch) *stateMachine {
	if obj.DeletionTimestamp != nil {
		return newStateMachine(newStep(nil, r.finalize, nil, nil))
	}
	setConfManagerState := newStep(r.configManagerTimeoutChecker, r.configManagerStatusFailed, r.statusUpdater, nil)
	setConfManagerStatus := newStep(r.configManagerStatusChecker, r.configManagerStatusSetter, r.statusUpdater, setConfManagerState)
	setReadyState := newStep(r.readyStatusChecker, r.readyStateSetter, r.statusUpdater, setConfManagerStatus)
	setLoopbackAddresses := newStep(r.switchAddressesChecker, r.switchAddressesSetter, r.specUpdater, setReadyState)
	setInterfacesSubnets := newStep(nil, r.interfacesSubnetsSetter, nil, setLoopbackAddresses)
	setIpAddresses := newStep(r.ipAddressesChecker, r.ipAddressesSetter, r.specUpdater, setInterfacesSubnets)
	setSubnets := newStep(r.subnetsChecker, r.subnetsSetter, r.statusUpdater, setIpAddresses)
	setConnectionLevel := newStep(r.connectionLevelChecker, r.connectionLevelSetter, r.statusUpdater, setSubnets)
	setRole := newStep(r.roleChecker, r.roleSetter, r.statusUpdater, setConnectionLevel)
	setPeers := newStep(r.peersInfoChecker, r.peersInfoSetter, r.statusUpdater, setRole)
	setAssignment := newStep(r.assignmentChecker, r.assignmentSetter, r.statusUpdater, setPeers)
	setInterfaces := newStep(r.interfacesChecker, r.interfacesSetter, r.specUpdater, setAssignment)
	setInitialStatus := newStep(r.initialStatusChecker, r.initialStatusSetter, r.statusUpdater, setInterfaces)
	setFinalizer := newStep(r.finalizerChecker, r.finalizerSetter, r.specUpdater, setInitialStatus)
	return newStateMachine(setFinalizer)
}

func (r *SwitchReconciler) statusUpdater(obj *switchv1alpha1.Switch) error {
	return r.Status().Update(r.Background.ctx, obj)
}

func (r *SwitchReconciler) specUpdater(obj *switchv1alpha1.Switch) error {
	return r.Update(r.Background.ctx, obj)
}

func (r *SwitchReconciler) finalizerChecker(obj *switchv1alpha1.Switch) bool {
	return obj.FinalizerOk()
}

func (r *SwitchReconciler) finalizerSetter(obj *switchv1alpha1.Switch) error {
	controllerutil.AddFinalizer(obj, switchv1alpha1.CSwitchFinalizer)
	return nil
}

func (r *SwitchReconciler) initialStatusChecker(obj *switchv1alpha1.Switch) bool {
	return !obj.CheckState(switchv1alpha1.CEmptyString)
}

func (r *SwitchReconciler) initialStatusSetter(obj *switchv1alpha1.Switch) error {
	obj.FillStatusOnCreate()
	return nil
}

func (r *SwitchReconciler) assignmentChecker(obj *switchv1alpha1.Switch) bool {
	if r.Background.assignment != nil {
		if r.Background.assignment.Status.State != switchv1alpha1.CAssignmentStateFinished {
			return false
		}
		if obj.Status.ConnectionLevel != 0 {
			return false
		}
	}
	return true
}

func (r *SwitchReconciler) assignmentSetter(obj *switchv1alpha1.Switch) error {
	obj.SetState(switchv1alpha1.CSwitchStateInProgress)
	r.Background.assignment.FillStatus(switchv1alpha1.CAssignmentStateFinished, &switchv1alpha1.LinkedSwitchSpec{
		Name:      obj.Name,
		Namespace: obj.Namespace,
	})
	if err := r.Status().Update(r.Background.ctx, r.Background.assignment); err != nil {
		r.Log.Error(err, "failed to update resource",
			"gvk", r.Background.assignment.GroupVersionKind().String(),
			"name", r.Background.assignment.NamespacedName())
		return err
	}
	obj.Status.ConnectionLevel = 0
	return nil
}

func (r *SwitchReconciler) interfacesChecker(obj *switchv1alpha1.Switch) bool {
	return obj.InterfacesMatchInventory(r.Background.inventory)
}

func (r *SwitchReconciler) interfacesSetter(obj *switchv1alpha1.Switch) error {
	interfaces, _ := switchv1alpha1.PrepareInterfaces(r.Background.inventory.Spec.NICs.NICs)
	stale := obj.UpdateInterfacesFromInventory(interfaces)
	for _, inf := range stale {
		if obj.Status.SouthSubnetV4 != nil {
			subnetName := fmt.Sprintf("%s-%s-v4", obj.Name, strings.ToLower(inf))
			if err := r.removeInterfaceSubnet(subnetName, obj.Status.SouthSubnetV4.ParentSubnet.Namespace); err != nil {
				return err
			}
		}
		if obj.Status.SouthSubnetV6 != nil {
			subnetName := fmt.Sprintf("%s-%s-v6", obj.Name, strings.ToLower(inf))
			if err := r.removeInterfaceSubnet(subnetName, obj.Status.SouthSubnetV6.ParentSubnet.Namespace); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *SwitchReconciler) peersInfoChecker(obj *switchv1alpha1.Switch) bool {
	return obj.PeersProcessingFinished(r.Background.switches)
}

func (r *SwitchReconciler) peersInfoSetter(obj *switchv1alpha1.Switch) error {
	obj.SetState(switchv1alpha1.CSwitchStateInProgress)
	obj.UpdateStoredPeers()
	obj.SetDiscoveredPeers(r.Background.switches)
	return nil
}

func (r *SwitchReconciler) roleChecker(obj *switchv1alpha1.Switch) bool {
	return obj.RoleOk()
}

func (r *SwitchReconciler) roleSetter(obj *switchv1alpha1.Switch) error {
	obj.UpdateRole()
	return nil
}

func (r *SwitchReconciler) connectionLevelChecker(obj *switchv1alpha1.Switch) bool {
	return obj.ConnectionLevelDefined(r.Background.switches, r.Background.assignment)
}

func (r *SwitchReconciler) connectionLevelSetter(obj *switchv1alpha1.Switch) error {
	obj.SetState(switchv1alpha1.CSwitchStateInProgress)
	obj.UpdateConnectionLevel(r.Background.switches)
	return nil
}

func (r *SwitchReconciler) subnetsChecker(obj *switchv1alpha1.Switch) bool {
	return obj.SubnetsOk()
}

func (r *SwitchReconciler) subnetsSetter(obj *switchv1alpha1.Switch) error {
	var assignment *switchv1alpha1.SwitchAssignment
	ctx := r.Background.ctx
	if r.Background.assignment == nil {
		topLevelSwitch := r.Background.switches.GetTopLevelSwitch()
		if topLevelSwitch == nil {
			return nil
		}
		swa, err := r.findAssignment(ctx, topLevelSwitch)
		if err != nil {
			return err
		}
		if swa == nil {
			return nil
		}
		assignment = swa
	} else {
		assignment = r.Background.assignment
	}

	labels := map[string]string{switchv1alpha1.LabelSwitchName: obj.Name}
	region := assignment.Spec.Region
	subnets := &subnetv1alpha1.SubnetList{}
	if err := r.Client.List(ctx, subnets); err != nil {
		return err
	}
	if obj.Status.SouthSubnetV4 == nil || obj.Status.SouthSubnetV4.CIDR == switchv1alpha1.CEmptyString {
		if err := r.updateSouthSubnetV4(obj, subnets, region, labels); err != nil {
			return err
		}
	}
	if obj.Status.SouthSubnetV6 == nil || obj.Status.SouthSubnetV6.CIDR == switchv1alpha1.CEmptyString {
		if err := r.updateSouthSubnetV6(obj, subnets, region, labels); err != nil {
			return err
		}
	}
	return nil
}

func (r *SwitchReconciler) updateSouthSubnetV4(
	obj *switchv1alpha1.Switch,
	subnets *subnetv1alpha1.SubnetList,
	region *switchv1alpha1.RegionSpec,
	labels map[string]string) error {
	var parentSubnet *subnetv1alpha1.Subnet
	var switchSubnet *subnetv1alpha1.Subnet
	var err error

	var noParentSubnetError = errors.New("failed to create v4 south subnet")

	for _, sn := range subnets.Items {
		if sn.Name == switchv1alpha1.CSwitchesParentSubnet+"-v4" && sn.Status.Type == subnetv1alpha1.CIPv4SubnetType {
			parentSubnet = &sn
			break
		}
	}
	if parentSubnet == nil {
		return fmt.Errorf("%w: no v4 subnet to reserve CIDR for south network exists", noParentSubnetError)
	}

	for _, item := range subnets.Items {
		if item.Name == fmt.Sprintf("%s-v4", switchv1alpha1.MacToResName(obj.Spec.Chassis.ChassisID)) {
			switchSubnet = &item
			break
		}
	}

	if switchSubnet == nil {
		addressesCount := obj.GetAddressCount(subnetv1alpha1.CIPv4SubnetType)
		switchSubnet, err = r.createSubnetWithCapacity(
			fmt.Sprintf("%s-v4", switchv1alpha1.MacToResName(obj.Spec.Chassis.ChassisID)),
			parentSubnet.Namespace,
			resource.NewQuantity(addressesCount, resource.DecimalSI),
			parentSubnet.Name,
			region.ConvertToSubnetRegion(),
			labels,
		)
		if err != nil {
			return err
		}
		obj.Status.SouthSubnetV4 = &switchv1alpha1.SwitchSubnetSpec{
			ParentSubnet: &switchv1alpha1.ParentSubnetSpec{
				Namespace: switchSubnet.Namespace,
				Name:      switchSubnet.Name,
				Region:    region,
			},
		}
		return err
	}

	subnetFinished := switchSubnet.Status.State == subnetv1alpha1.CFinishedSubnetState
	CIDRUnfilled := obj.Status.SouthSubnetV4.CIDR == switchv1alpha1.CEmptyString
	if subnetFinished && CIDRUnfilled {
		obj.Status.SouthSubnetV4.CIDR = switchSubnet.Status.Reserved.String()
	}
	return err
}

func (r *SwitchReconciler) updateSouthSubnetV6(
	obj *switchv1alpha1.Switch,
	subnets *subnetv1alpha1.SubnetList,
	region *switchv1alpha1.RegionSpec,
	labels map[string]string) error {
	var parentSubnet *subnetv1alpha1.Subnet
	var switchSubnet *subnetv1alpha1.Subnet
	var err error

	var noParentSubnetError = errors.New("failed to create v6 south subnet")

	for _, sn := range subnets.Items {
		if sn.Name == switchv1alpha1.CSwitchesParentSubnet+"-v6" && sn.Status.Type == subnetv1alpha1.CIPv6SubnetType {
			parentSubnet = &sn
			break
		}
	}
	if parentSubnet == nil {
		return fmt.Errorf("%w: no v6 subnet to reserve CIDR for south network exists", noParentSubnetError)
	}

	for _, item := range subnets.Items {
		if item.Name == fmt.Sprintf("%s-v6", switchv1alpha1.MacToResName(obj.Spec.Chassis.ChassisID)) {
			switchSubnet = &item
			break
		}
	}

	if switchSubnet == nil {
		addressesCount := obj.GetAddressCount(subnetv1alpha1.CIPv6SubnetType)
		switchSubnet, err = r.createSubnetWithCapacity(
			fmt.Sprintf("%s-v6", switchv1alpha1.MacToResName(obj.Spec.Chassis.ChassisID)),
			parentSubnet.Namespace,
			resource.NewQuantity(addressesCount, resource.DecimalSI),
			parentSubnet.Name,
			region.ConvertToSubnetRegion(),
			labels,
		)
		if err != nil {
			return err
		}
		obj.Status.SouthSubnetV6 = &switchv1alpha1.SwitchSubnetSpec{
			ParentSubnet: &switchv1alpha1.ParentSubnetSpec{
				Namespace: switchSubnet.Namespace,
				Name:      switchSubnet.Name,
				Region:    region,
			},
		}
		return err
	}

	subnetFinished := switchSubnet.Status.State == subnetv1alpha1.CFinishedSubnetState
	CIDRUnfilled := obj.Status.SouthSubnetV6.CIDR == switchv1alpha1.CEmptyString
	if subnetFinished && CIDRUnfilled {
		obj.Status.SouthSubnetV6.CIDR = switchSubnet.Status.Reserved.String()
	}
	return err
}

func (r *SwitchReconciler) ipAddressesChecker(obj *switchv1alpha1.Switch) bool {
	return obj.AddressesDefined() && obj.AddressesOk(r.Background.switches)
}

func (r *SwitchReconciler) ipAddressesSetter(obj *switchv1alpha1.Switch) error {
	obj.UpdateSouthInterfacesAddresses()
	obj.UpdateNorthInterfacesAddresses(r.Background.switches)
	return nil
}

func (r *SwitchReconciler) interfacesSubnetsSetter(obj *switchv1alpha1.Switch) error {
	labels := map[string]string{switchv1alpha1.LabelSwitchName: obj.Name}
	for iface, data := range obj.Spec.Interfaces {
		labels[switchv1alpha1.LabelInterfaceName] = strings.ToLower(iface)
		if _, ok := obj.Status.SouthConnections.Peers[iface]; !ok {
			continue
		}
		if data.IPv4 != switchv1alpha1.CEmptyString {
			if err := r.updateInterfaceSubnetV4(obj, iface, data.IPv4, labels); err != nil {
				return err
			}
		}
		if data.IPv6 != switchv1alpha1.CEmptyString {
			if err := r.updateInterfaceSubnetV6(obj, iface, data.IPv6, labels); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *SwitchReconciler) updateInterfaceSubnetV4(obj *switchv1alpha1.Switch, name string, address string, labels map[string]string) error {
	var err error
	if err = r.Get(r.Background.ctx, types.NamespacedName{
		Namespace: obj.Status.SouthSubnetV4.ParentSubnet.Namespace,
		Name:      fmt.Sprintf("%s-%s-v4", switchv1alpha1.MacToResName(obj.Spec.Chassis.ChassisID), strings.ToLower(name)),
	}, &subnetv1alpha1.Subnet{}); err != nil {
		if apierrors.IsNotFound(err) {
			_, cidr, _ := net.ParseCIDR(address)
			if _, err = r.createSubnetWithCIDR(
				fmt.Sprintf("%s-%s-v4", switchv1alpha1.MacToResName(obj.Spec.Chassis.ChassisID), strings.ToLower(name)),
				obj.Status.SouthSubnetV4.ParentSubnet.Namespace,
				subnetv1alpha1.CIDRFromNet(cidr),
				obj.Status.SouthSubnetV4.ParentSubnet.Name,
				obj.Status.SouthSubnetV4.ParentSubnet.Region.ConvertToSubnetRegion(),
				labels); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return err
}

func (r *SwitchReconciler) updateInterfaceSubnetV6(obj *switchv1alpha1.Switch, name string, address string, labels map[string]string) error {
	var err error
	if err = r.Get(r.Background.ctx, types.NamespacedName{
		Namespace: obj.Status.SouthSubnetV6.ParentSubnet.Namespace,
		Name:      fmt.Sprintf("%s-%s-v6", switchv1alpha1.MacToResName(obj.Spec.Chassis.ChassisID), strings.ToLower(name)),
	}, &subnetv1alpha1.Subnet{}); err != nil {
		if apierrors.IsNotFound(err) {
			_, cidr, _ := net.ParseCIDR(address)
			if _, err = r.createSubnetWithCIDR(
				fmt.Sprintf("%s-%s-v6", switchv1alpha1.MacToResName(obj.Spec.Chassis.ChassisID), strings.ToLower(name)),
				obj.Status.SouthSubnetV6.ParentSubnet.Namespace,
				subnetv1alpha1.CIDRFromNet(cidr),
				obj.Status.SouthSubnetV6.ParentSubnet.Name,
				obj.Status.SouthSubnetV6.ParentSubnet.Region.ConvertToSubnetRegion(),
				labels); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return err
}

func (r *SwitchReconciler) switchAddressesChecker(obj *switchv1alpha1.Switch) bool {
	return obj.SwitchAddressesDefined()
}

func (r *SwitchReconciler) switchAddressesSetter(obj *switchv1alpha1.Switch) error {
	labels := map[string]string{switchv1alpha1.LabelSwitchName: obj.Name}
	if obj.Spec.IPv4 == switchv1alpha1.CEmptyString {
		if err := r.updateSwitchV4Address(obj, labels); err != nil {
			return err
		}
	}
	if obj.Spec.IPv6 == switchv1alpha1.CEmptyString {
		if err := r.updateSwitchV6Address(obj, labels); err != nil {
			return err
		}
	}
	return nil
}

func (r *SwitchReconciler) updateSwitchV4Address(obj *switchv1alpha1.Switch, labels map[string]string) error {
	if obj.Status.SouthSubnetV4 == nil {
		return nil
	}
	if err := r.Get(r.Background.ctx, types.NamespacedName{
		Namespace: obj.Status.SouthSubnetV4.ParentSubnet.Namespace,
		Name:      CSwitchLoopbackV4Subnet,
	}, &subnetv1alpha1.Subnet{}); err != nil {
		return err
	}
	sn := &subnetv1alpha1.Subnet{}
	if err := r.Get(r.Background.ctx, types.NamespacedName{
		Namespace: obj.Status.SouthSubnetV4.ParentSubnet.Namespace,
		Name:      fmt.Sprintf("%s-lo-v4", obj.Name),
	}, sn); err != nil {
		if apierrors.IsNotFound(err) {
			sn, err = r.createSubnetWithCapacity(
				fmt.Sprintf("%s-lo-v4", obj.Name),
				obj.Status.SouthSubnetV4.ParentSubnet.Namespace,
				resource.NewQuantity(1, resource.DecimalSI),
				CSwitchLoopbackV4Subnet,
				obj.Status.SouthSubnetV4.ParentSubnet.Region.ConvertToSubnetRegion(),
				labels,
			)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if sn.Status.Reserved != nil {
		obj.Spec.IPv4 = sn.Status.Reserved.String()
	}
	return nil
}

func (r *SwitchReconciler) updateSwitchV6Address(obj *switchv1alpha1.Switch, labels map[string]string) error {
	if obj.Status.SouthSubnetV6 == nil {
		return nil
	}
	if err := r.Get(r.Background.ctx, types.NamespacedName{
		Namespace: obj.Status.SouthSubnetV6.ParentSubnet.Namespace,
		Name:      CSwitchLoopbackV6Subnet,
	}, &subnetv1alpha1.Subnet{}); err != nil {
		return err
	}
	sn := &subnetv1alpha1.Subnet{}
	if err := r.Get(r.Background.ctx, types.NamespacedName{
		Namespace: obj.Status.SouthSubnetV6.ParentSubnet.Namespace,
		Name:      fmt.Sprintf("%s-lo-v6", obj.Name),
	}, sn); err != nil {
		if apierrors.IsNotFound(err) {
			sn, err = r.createSubnetWithCapacity(
				fmt.Sprintf("%s-lo-v6", obj.Name),
				obj.Status.SouthSubnetV6.ParentSubnet.Namespace,
				resource.NewQuantity(1, resource.DecimalSI),
				CSwitchLoopbackV6Subnet,
				obj.Status.SouthSubnetV6.ParentSubnet.Region.ConvertToSubnetRegion(),
				labels,
			)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if sn.Status.Reserved != nil {
		obj.Spec.IPv6 = sn.Status.Reserved.String()
	}
	return nil
}

func (r *SwitchReconciler) readyStatusChecker(obj *switchv1alpha1.Switch) bool {
	return !(r.connectionLevelChecker(obj) &&
		r.subnetsChecker(obj) &&
		r.ipAddressesChecker(obj) &&
		r.switchAddressesChecker(obj))
}

func (r *SwitchReconciler) readyStateSetter(obj *switchv1alpha1.Switch) error {
	obj.SetState(switchv1alpha1.CSwitchStateReady)
	return nil
}

func (r *SwitchReconciler) configManagerStatusChecker(obj *switchv1alpha1.Switch) bool {
	return obj.ConfigManagerStatusOk()
}

func (r *SwitchReconciler) configManagerStatusSetter(obj *switchv1alpha1.Switch) error {
	obj.SetConfigManagerStatus(switchv1alpha1.CEmptyString)
	return nil
}

func (r *SwitchReconciler) configManagerTimeoutChecker(obj *switchv1alpha1.Switch) bool {
	return obj.ConfigManagerTimeoutOk()
}

func (r *SwitchReconciler) configManagerStatusFailed(obj *switchv1alpha1.Switch) error {
	obj.SetConfigManagerStatus(switchv1alpha1.CConfigManagementTypeFailed)
	return nil
}
