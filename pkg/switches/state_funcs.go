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
	"fmt"
	"reflect"
	"strings"

	"github.com/go-errors/errors"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/pkg/constants"
	v1 "k8s.io/api/core/v1"
)

func Initialize(obj *switchv1beta1.Switch, _ *SwitchEnvironment) *StateUpdateResult {
	res := NewStateUpdateResult().
		SetCondition(constants.ConditionInitialized).
		SetReason(constants.ReasonConditionInitialized).
		SetMessage(constants.MessageConditionInitialized)
	if obj.Uninitialized() {
		obj.Status = switchv1beta1.SwitchStatus{
			Conditions:            make([]*switchv1beta1.ConditionSpec, 0),
			ConfigRef:             v1.LocalObjectReference{},
			RoutingConfigTemplate: v1.LocalObjectReference{},
			ASN:                   0,
			TotalPorts:            0,
			SwitchPorts:           0,
			Role:                  constants.SwitchRoleLeaf,
			Layer:                 255,
			Interfaces:            make(map[string]*switchv1beta1.InterfaceSpec),
			LoopbackAddresses:     make([]*switchv1beta1.IPAddressSpec, 0),
			Subnets:               make([]*switchv1beta1.SubnetSpec, 0),
			Message:               nil,
		}
		if obj.TopSpine() {
			obj.SetLayer(0)
			obj.SetRole(constants.SwitchRoleSpine)
		}
		obj.SetState(constants.SwitchStateInitial)
		obj.SetMessage(constants.EmptyString)
	}
	if obj.TopSpine() && obj.GetLayer() != 0 {
		obj.SetLayer(0)
	}
	return res
}

func UpdateInterfaces(obj *switchv1beta1.Switch, env *SwitchEnvironment) *StateUpdateResult {
	res := NewStateUpdateResult().
		SetCondition(constants.ConditionInterfacesOK).
		SetReason(constants.ReasonConditionInterfacesOK).
		SetMessage(constants.MessageConditionInterfacesOK)
	if env.Inventory == nil {
		res.SetError(errors.Wrap(constants.ErrorUpdateInterfacesFailed, 0)).
			SetReason(constants.ErrorReasonMissingRequirements).
			SetMessage(constants.MessageMissingInventory)
		obj.SetState(constants.SwitchStateInvalid)
		obj.SetMessage(constants.StateMessageMissingRequirements)
		return res
	}
	ApplyInterfacesFromInventory(obj, env.Inventory)
	return res
}

func UpdateNeighbors(obj *switchv1beta1.Switch, env *SwitchEnvironment) *StateUpdateResult {
	res := NewStateUpdateResult().
		SetCondition(constants.ConditionNeighborsOK).
		SetReason(constants.ReasonConditionNeighborsOK).
		SetMessage(constants.MessageConditionNeighborsOK)
	if env.Switches == nil {
		res.SetError(errors.Wrap(constants.ErrorUpdateNeighborsFailed, 0)).
			SetReason(constants.ErrorReasonRequestFailed).
			SetMessage(fmt.Sprintf("%s: %s", constants.MessageRequestFailed, "SwitchList"))
		obj.SetState(constants.SwitchStateInvalid)
		obj.SetMessage(constants.StateMessageRequestRelatedObjectsFailed)
		return res
	}
	for _, item := range env.Switches.Items {
		for _, nicData := range obj.Status.Interfaces {
			if nicData.Peer == nil {
				continue
			}
			if nicData.Peer.PeerInfoSpec == nil {
				continue
			}
			if reflect.DeepEqual(nicData.Peer.PeerInfoSpec, &switchv1beta1.PeerInfoSpec{}) {
				continue
			}
			peerChassisID := nicData.Peer.PeerInfoSpec.GetChassisID()
			if strings.ReplaceAll(peerChassisID, ":", "") != item.Annotations[constants.HardwareChassisIDAnnotation] {
				continue
			}
			nicData.Peer.SetObjectReference(item.Name, item.Namespace)
		}
	}
	obj.SetState(constants.SwitchStateProcessing)
	obj.SetMessage(constants.EmptyString)
	return res
}

func UpdateLayerAndRole(obj *switchv1beta1.Switch, env *SwitchEnvironment) *StateUpdateResult {
	res := NewStateUpdateResult().
		SetCondition(constants.ConditionLayerAndRoleOK).
		SetReason(constants.ReasonConditionLayerAndRoleOK).
		SetMessage(constants.MessageConditionLayerAndRoleOK)
	if env.Switches == nil {
		res.SetError(errors.Wrap(constants.ErrorUpdateLayerAndRoleFailed, 0)).
			SetReason(constants.ErrorReasonRequestFailed).
			SetMessage(fmt.Sprintf("%s: %s", constants.MessageRequestFailed, "SwitchList"))
		obj.SetState(constants.SwitchStateInvalid)
		obj.SetMessage(constants.StateMessageRequestRelatedObjectsFailed)
		return res
	}
	ComputeLayer(obj, env.Switches)
	if obj.GetLayer() == 255 {
		res.SetError(errors.Wrap(constants.ErrorUpdateLayerAndRoleFailed, 0)).
			SetReason(constants.ErrorReasonFailedToComputeLayer).
			SetMessage(constants.MessageFailedToComputeLayer)
		obj.SetState(constants.SwitchStateInvalid)
		obj.SetMessage(constants.StateMessageRelatedObjectsStateInvalid)
		return res
	}
	SetRole(obj)
	obj.SetState(constants.SwitchStateProcessing)
	obj.SetMessage(constants.EmptyString)
	return res
}

func UpdateConfigRef(obj *switchv1beta1.Switch, env *SwitchEnvironment) *StateUpdateResult {
	res := NewStateUpdateResult().
		SetCondition(constants.ConditionConfigRefOK).
		SetReason(constants.ReasonConditionConfigRefOK).
		SetMessage(constants.MessageConditionConfigRefOK)
	if env.Config == nil {
		res.SetError(errors.Wrap(constants.ErrorUpdateConfigRefFailed, 0)).
			SetReason(constants.ErrorReasonMissingRequirements).
			SetMessage(constants.MessageFailedToDiscoverConfig)
		obj.SetState(constants.SwitchStatePending)
		obj.SetMessage(constants.StateMessageMissingRequirements)
		return res
	}
	obj.SetConfigRef(env.Config.Name)
	if !env.Config.RoutingConfigTemplateIsEmpty() {
		obj.SetRoutingConfigTemplate(env.Config.GetRoutingConfigTemplate())
	}
	obj.SetState(constants.SwitchStateProcessing)
	obj.SetMessage(constants.EmptyString)
	return res
}

func UpdatePortParameters(obj *switchv1beta1.Switch, env *SwitchEnvironment) *StateUpdateResult {
	res := NewStateUpdateResult().
		SetCondition(constants.ConditionPortParametersOK).
		SetReason(constants.ReasonConditionPortParametersOK).
		SetMessage(constants.MessageConditionPortParametersOK)
	if env.Config == nil {
		res.SetError(errors.Wrap(constants.ErrorUpdatePortParametersFailed, 0)).
			SetReason(constants.ErrorReasonMissingRequirements).
			SetMessage(constants.MessageFailedToDiscoverConfig)
		obj.SetState(constants.SwitchStateInvalid)
		obj.SetMessage(constants.StateMessageMissingRequirements)
		return res
	}
	if env.Switches == nil {
		res.SetError(errors.Wrap(constants.ErrorUpdatePortParametersFailed, 0)).
			SetReason(constants.ErrorReasonRequestFailed).
			SetMessage(fmt.Sprintf("%s: %s", constants.MessageRequestFailed, "SwitchList"))
		obj.SetState(constants.SwitchStateInvalid)
		obj.SetMessage(constants.StateMessageRequestRelatedObjectsFailed)
		return res
	}
	ApplyInterfaceParams(obj, env.Config)
	InheritInterfaceParams(obj, env.Switches)
	AlignInterfacesWithParams(obj)
	obj.SetState(constants.SwitchStateProcessing)
	obj.SetMessage(constants.EmptyString)
	return res
}

func UpdateLoopbacks(obj *switchv1beta1.Switch, env *SwitchEnvironment) *StateUpdateResult {
	res := NewStateUpdateResult().
		SetCondition(constants.ConditionLoopbacksOK).
		SetReason(constants.ReasonConditionLoopbacksOK).
		SetMessage(constants.MessageConditionLoopbacksOK)
	if env.LoopbackIPs == nil {
		res.SetError(errors.Wrap(constants.ErrorUpdateLoopbacksFailed, 0)).
			SetReason(constants.ErrorReasonMissingRequirements).
			SetMessage(constants.MessageMissingLoopbacks)
		obj.SetState(constants.SwitchStateInvalid)
		obj.SetMessage(constants.StateMessageMissingRequirements)
		return res
	}
	loopbacksToApply := make([]*switchv1beta1.IPAddressSpec, 0)
	for _, item := range env.LoopbackIPs.Items {
		var af string
		switch {
		case item.Status.Reserved.Net.Is4():
			af = constants.IPv4AF
		case item.Status.Reserved.Net.Is6():
			af = constants.IPv6AF
		}
		ip := &switchv1beta1.IPAddressSpec{}
		ip.SetObjectReference(item.Name, item.Namespace)
		ip.SetAddress(item.Status.Reserved.String())
		ip.SetAddressFamily(af)
		ip.SetExtraAddress(false)
		loopbacksToApply = append(loopbacksToApply, ip)
	}
	obj.Status.LoopbackAddresses = make([]*switchv1beta1.IPAddressSpec, len(loopbacksToApply))
	copy(obj.Status.LoopbackAddresses, loopbacksToApply)
	obj.SetState(constants.SwitchStateProcessing)
	obj.SetMessage(constants.EmptyString)
	return res
}

func UpdateASN(obj *switchv1beta1.Switch, _ *SwitchEnvironment) *StateUpdateResult {
	res := NewStateUpdateResult().
		SetCondition(constants.ConditionAsnOK).
		SetReason(constants.ReasonConditionAsnOK).
		SetMessage(constants.MessageConditionAsnOK)
	asn, err := CalculateASN(obj.Status.LoopbackAddresses)
	if err != nil {
		res.SetError(err).
			SetReason(constants.ErrorReasonASNCalculationFailed).
			SetMessage(constants.ErrorUpdateASNFailed)
		obj.SetState(constants.SwitchStateInvalid)
		obj.SetMessage(constants.ErrorUpdateASNFailed)
		return res
	}
	obj.SetASN(asn)
	obj.SetState(constants.SwitchStateProcessing)
	obj.SetMessage(constants.EmptyString)
	return res
}

func UpdateSubnets(obj *switchv1beta1.Switch, env *SwitchEnvironment) *StateUpdateResult {
	res := NewStateUpdateResult().
		SetCondition(constants.ConditionSubnetsOK).
		SetReason(constants.ReasonConditionSubnetsOK).
		SetMessage(constants.MessageConditionSubnetsOK)
	if env.SouthSubnets == nil {
		res.SetError(errors.Wrap(constants.ErrorUpdateSubnetsFailed, 0)).
			SetReason(constants.ErrorReasonMissingRequirements).
			SetMessage(constants.MessageMissingSouthSubnets)
		obj.SetState(constants.SwitchStateInvalid)
		obj.SetMessage(constants.StateMessageMissingRequirements)
		return res
	}
	subnetsToApply := make([]*switchv1beta1.SubnetSpec, 0)
	for _, item := range env.SouthSubnets.Items {
		subnet := &switchv1beta1.SubnetSpec{}
		subnet.SetSubnetObjectRef(item.Name, item.Namespace)
		subnet.SetNetworkObjectRef(item.Spec.Network.Name, item.Namespace)
		subnet.SetCIDR(item.Status.Reserved.Net.String())
		subnet.SetAddressFamily(string(item.Status.Type))
		subnetsToApply = append(subnetsToApply, subnet)
	}
	obj.Status.Subnets = make([]*switchv1beta1.SubnetSpec, len(subnetsToApply))
	copy(obj.Status.Subnets, subnetsToApply)
	obj.SetState(constants.SwitchStateProcessing)
	obj.SetMessage(constants.EmptyString)
	return res
}

func UpdateSwitchPortIPs(obj *switchv1beta1.Switch, env *SwitchEnvironment) *StateUpdateResult {
	var err error
	res := NewStateUpdateResult().
		SetCondition(constants.ConditionIPAddressesOK).
		SetReason(constants.ReasonConditionIPAddressesOK).
		SetMessage(constants.MessageConditionIPAddressesOK)
	for name, data := range obj.Status.Interfaces {
		if !strings.HasPrefix(name, constants.SwitchPortNamePrefix) {
			continue
		}
		switch data.GetDirection() {
		case constants.DirectionNorth:
			if data.Peer == nil {
				continue
			}
			err = updateNorthIPs(data, env)
		case constants.DirectionSouth:
			err = updateSouthIPs(name, obj, data)
		}
		if err != nil {
			res.SetError(errors.Wrap(err, 0)).
				SetReason(constants.ErrorReasonIPAssignmentFailed).
				SetMessage(constants.ErrorUpdateSwitchPortIPsFailed)
			obj.SetState(constants.SwitchStateInvalid)
			obj.SetMessage(constants.MessageFailedToAssignIPAddresses)
			return res
		}
	}
	obj.SetState(constants.SwitchStateProcessing)
	obj.SetMessage(constants.EmptyString)
	return res
}

func SetStateReady(obj *switchv1beta1.Switch, _ *SwitchEnvironment) *StateUpdateResult {
	res := NewStateUpdateResult().
		SetCondition(constants.ConditionReady).
		SetReason(constants.ReasonConditionReady).
		SetMessage(constants.MessageConditionReady)
	obj.SetState(constants.SwitchStateReady)
	obj.SetMessage(constants.EmptyString)
	return res
}

func updateNorthIPs(data *switchv1beta1.InterfaceSpec, env *SwitchEnvironment) error {
	if env.Switches == nil {
		return errors.New(fmt.Sprintf("%s: %s", constants.MessageRequestFailed, "SwitchList"))
	}
	ipsToApply := make([]*switchv1beta1.IPAddressSpec, 0)
	for _, item := range env.Switches.Items {
		if item.Name != data.Peer.GetObjectReferenceName() {
			continue
		}
		peerNICData := GetPeerData(item.Status.Interfaces, data.Peer.GetPortDescription(), data.Peer.GetPortID())
		if peerNICData == nil {
			continue
		}
		requestedIPs := RequestIPs(peerNICData)
		if len(requestedIPs) == 0 {
			return errors.New(constants.MessageFailedIPAddressRequest)
		}
		ipsToApply = append(ipsToApply, requestedIPs...)
		data.IP = make([]*switchv1beta1.IPAddressSpec, len(ipsToApply))
		copy(data.IP, ipsToApply)
	}
	return nil
}

func updateSouthIPs(nic string, obj *switchv1beta1.Switch, data *switchv1beta1.InterfaceSpec) error {
	ipsToApply := make([]*switchv1beta1.IPAddressSpec, 0)
	extraIPs, err := GetExtraIPs(obj, nic)
	if err != nil {
		return err
	}
	ipsToApply = append(ipsToApply, extraIPs...)
	computedIPs, _, err := GetComputedIPs(obj, nic, data)
	if err != nil {
		return err
	}
	ipsToApply = append(ipsToApply, computedIPs...)
	data.IP = make([]*switchv1beta1.IPAddressSpec, len(ipsToApply))
	copy(data.IP, ipsToApply)
	return nil
}