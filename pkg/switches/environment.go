/*
Copyright (c) 2023 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package switches

import (
	"context"

	"github.com/go-logr/logr"
	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SwitchEnvironment struct {
	Inventory         *inventoryv1alpha1.Inventory
	Switches          *switchv1beta1.SwitchList
	Config            *switchv1beta1.SwitchConfig
	LoopbackIPs       *ipamv1alpha1.IPList
	SouthSubnets      *ipamv1alpha1.SubnetList
	SwitchPortSubnets *ipamv1alpha1.SubnetList
}

type SwitchEnvironmentSvc struct {
	client.Client
	Log logr.Logger
	Env *SwitchEnvironment
}

func NewSwitchEnvironmentSvc(cl client.Client, log logr.Logger) *SwitchEnvironmentSvc {
	return &SwitchEnvironmentSvc{cl, log, nil}
}

func (in *SwitchEnvironmentSvc) GetEnvironment(ctx context.Context, obj *switchv1beta1.Switch) *SwitchEnvironment {
	in.Log.Info("gathering info about environment")
	inventory := in.GetInventory(ctx, obj)
	switches := in.GetSwitches(ctx, obj)
	config := in.GetSwitchConfig(ctx, obj)
	loopbacks := in.GetLoopbacks(ctx, obj, config)
	subnets := in.GetSubnets(ctx, obj, config)
	switchPortsSubnets := in.GetSwitchPortsSubnets(ctx, obj, config)
	in.Env = &SwitchEnvironment{
		Inventory:         inventory,
		Switches:          switches,
		Config:            config,
		LoopbackIPs:       loopbacks,
		SouthSubnets:      subnets,
		SwitchPortSubnets: switchPortsSubnets,
	}
	return in.Env
}

func (in *SwitchEnvironmentSvc) GetInventory(
	ctx context.Context,
	obj *switchv1beta1.Switch,
) *inventoryv1alpha1.Inventory {
	in.Log.Info("requesting for related Inventory object")
	if obj.GetInventoryRef() == "" {
		return nil
	}
	inventory := &inventoryv1alpha1.Inventory{}
	key := client.ObjectKeyFromObject(obj)
	err := in.Get(ctx, key, inventory)
	if err != nil {
		return nil
	}
	return inventory
}

func (in *SwitchEnvironmentSvc) GetSwitches(
	ctx context.Context,
	obj *switchv1beta1.Switch,
) *switchv1beta1.SwitchList {
	in.Log.Info("requesting for list of existing switches")
	switches := &switchv1beta1.SwitchList{}
	inventoriedLabelReq, _ := labels.NewRequirement(constants.InventoriedLabel, selection.Exists, []string{})
	selector := labels.NewSelector().Add(*inventoriedLabelReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     obj.Namespace,
		Limit:         100,
	}
	if err := in.List(ctx, switches, opts); err != nil {
		return nil
	}
	return switches
}

func (in *SwitchEnvironmentSvc) GetSwitchConfig(
	ctx context.Context,
	obj *switchv1beta1.Switch,
) *switchv1beta1.SwitchConfig {
	in.Log.Info("requesting for related SwitchConfig object")
	switchConfigs := &switchv1beta1.SwitchConfigList{}
	selector, err := metav1.LabelSelectorAsSelector(obj.GetConfigSelector())
	if err != nil {
		return nil
	}
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     obj.Namespace,
		Limit:         100,
	}
	if err = in.List(ctx, switchConfigs, opts); err != nil {
		return nil
	}
	if len(switchConfigs.Items) == 0 {
		return nil
	}
	if len(switchConfigs.Items) > 1 {
		return nil
	}
	return &switchConfigs.Items[0]
}

func (in *SwitchEnvironmentSvc) GetLoopbacks(
	ctx context.Context,
	obj *switchv1beta1.Switch,
	cfg *switchv1beta1.SwitchConfig,
) *ipamv1alpha1.IPList {
	in.Log.Info("requesting for list of related IP objects")
	if cfg == nil {
		return nil
	}
	params := cfg.Spec.IPAM.LoopbackAddresses
	if obj.Spec.IPAM != nil && obj.Spec.IPAM.LoopbackAddresses != nil {
		params = obj.Spec.IPAM.LoopbackAddresses
	}
	loopbacks := &ipamv1alpha1.IPList{}
	err := in.ListIPAMObjects(ctx, obj, params, loopbacks)
	if err != nil {
		return nil
	}
	addressFamiliesMap := map[ipamv1alpha1.SubnetAddressType]*bool{
		ipamv1alpha1.CIPv4SubnetType: nil,
		ipamv1alpha1.CIPv6SubnetType: nil,
	}
	afEnabledFlag := 0
	for _, item := range loopbacks.Items {
		if item.Status.State != ipamv1alpha1.CFinishedIPState {
			continue
		}
		if !cfg.Spec.IPAM.AddressFamily.GetIPv6() && item.Status.Reserved.Net.Is6() {
			continue
		}
		switch {
		case item.Status.Reserved.Net.Is4():
			addressFamiliesMap[ipamv1alpha1.CIPv4SubnetType] = pointer.Bool(true)
			afEnabledFlag = afEnabledFlag | 1
		case item.Status.Reserved.Net.Is6():
			addressFamiliesMap[ipamv1alpha1.CIPv6SubnetType] = pointer.Bool(true)
			afEnabledFlag = afEnabledFlag | 2
		}
	}
	afOK := AddressFamiliesMatchConfig(true, cfg.Spec.IPAM.AddressFamily.GetIPv6(), afEnabledFlag)
	if len(loopbacks.Items) == 0 || !afOK {
		return nil
	}
	return loopbacks
}

func (in *SwitchEnvironmentSvc) GetSubnets(
	ctx context.Context,
	obj *switchv1beta1.Switch,
	cfg *switchv1beta1.SwitchConfig,
) *ipamv1alpha1.SubnetList {
	in.Log.Info("requesting for list of related Subnet objects")
	if cfg == nil {
		return nil
	}
	c := obj.GetCondition(constants.ConditionPortParametersOK)
	if !c.GetState() {
		return nil
	}
	af := cfg.Spec.IPAM.AddressFamily
	params := cfg.Spec.IPAM.SouthSubnets
	if obj.Spec.IPAM != nil && obj.Spec.IPAM.SouthSubnets != nil {
		params = obj.Spec.IPAM.SouthSubnets
	}
	subnets := &ipamv1alpha1.SubnetList{}
	err := in.ListIPAMObjects(ctx, obj, params, subnets)
	if err != nil {
		return nil
	}
	addressFamiliesMap := map[ipamv1alpha1.SubnetAddressType]*bool{
		ipamv1alpha1.CIPv4SubnetType: nil,
		ipamv1alpha1.CIPv6SubnetType: nil,
	}
	afEnabledFlag := 0
	for _, item := range subnets.Items {
		if item.Status.State == ipamv1alpha1.CFailedSubnetState {
			continue
		}
		if item.Status.State == ipamv1alpha1.CProcessingSubnetState {
			addressFamiliesMap[item.Status.Type] = pointer.Bool(true)
			continue
		}
		if (!af.GetIPv4() && item.Status.Reserved.IsIPv4()) || (!af.GetIPv6() && item.Status.Reserved.IsIPv6()) {
			continue
		}
		requiredCapacity := GetTotalAddressesCount(obj.Status.Interfaces, item.Status.Type)
		if requiredCapacity.IsZero() {
			continue
		}
		if requiredCapacity.Cmp(item.Status.Capacity) > 0 {
			continue
		}
		addressFamiliesMap[item.Status.Type] = pointer.Bool(true)
		switch item.Status.Type {
		case ipamv1alpha1.CIPv4SubnetType:
			afEnabledFlag = afEnabledFlag | 1
		case ipamv1alpha1.CIPv6SubnetType:
			afEnabledFlag = afEnabledFlag | 2
		}
	}
	afOK := AddressFamiliesMatchConfig(af.GetIPv4(), af.GetIPv6(), afEnabledFlag)
	if len(subnets.Items) == 0 || !afOK {
		return nil
	}
	return subnets
}

func (in *SwitchEnvironmentSvc) GetSwitchPortsSubnets(
	ctx context.Context,
	obj *switchv1beta1.Switch,
	cfg *switchv1beta1.SwitchConfig,
) *ipamv1alpha1.SubnetList {
	in.Log.Info("requesting for list of related switch ports Subnet objects")
	if cfg == nil {
		return nil
	}
	subnets := &ipamv1alpha1.SubnetList{}
	selector := labels.NewSelector()
	purposeReq, _ := labels.NewRequirement(constants.IPAMObjectPurposeLabel, selection.In, []string{constants.IPAMSwitchPortPurpose})
	ownerReq, _ := labels.NewRequirement(constants.IPAMObjectOwnerLabel, selection.In, []string{obj.Name})
	selector = selector.Add(*purposeReq).Add(*ownerReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     obj.Namespace,
	}
	if err := in.List(ctx, subnets, opts); err != nil {
		return nil
	}
	subnetsCountPerInterface := make(map[string]int)
	afEnabledCount := func(af *switchv1beta1.AddressFamiliesMap) int {
		count := 0
		if af.GetIPv4() {
			count = count | 1
		}
		if af.GetIPv6() {
			count = count | 2
		}
		return count
	}(cfg.Spec.IPAM.AddressFamily)
	for _, item := range subnets.Items {
		if item.Status.State == ipamv1alpha1.CFailedSubnetState {
			return nil
		}
		nicName := ParseInterfaceNameFromSubnet(item.Name)
		counter, ok := subnetsCountPerInterface[nicName]
		if !ok {
			counter = 0
		}
		if item.Status.Type == ipamv1alpha1.CIPv4SubnetType {
			subnetsCountPerInterface[nicName] = counter | 1
		}
		if item.Status.Type == ipamv1alpha1.CIPv6SubnetType {
			subnetsCountPerInterface[nicName] = counter | 2
		}
	}
	if len(subnetsCountPerInterface) != len(obj.Status.Interfaces) {
		return nil
	}
	for _, v := range subnetsCountPerInterface {
		if v != afEnabledCount {
			return nil
		}
	}
	return subnets
}

func (in *SwitchEnvironmentSvc) ListIPAMObjects(
	ctx context.Context,
	obj *switchv1beta1.Switch,
	params *switchv1beta1.IPAMSelectionSpec,
	list client.ObjectList,
) error {
	selector, err := GetSelectorFromIPAMSpec(obj, params)
	if err != nil {
		return err
	}
	opts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     obj.Namespace,
		Limit:         100,
	}
	if err := in.List(ctx, list, opts); err != nil {
		return err
	}
	return nil
}
