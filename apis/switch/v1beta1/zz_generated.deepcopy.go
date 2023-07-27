//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AdditionalIPSpec) DeepCopyInto(out *AdditionalIPSpec) {
	*out = *in
	if in.Address != nil {
		in, out := &in.Address, &out.Address
		*out = new(string)
		**out = **in
	}
	if in.ParentSubnet != nil {
		in, out := &in.ParentSubnet, &out.ParentSubnet
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdditionalIPSpec.
func (in *AdditionalIPSpec) DeepCopy() *AdditionalIPSpec {
	if in == nil {
		return nil
	}
	out := new(AdditionalIPSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AddressFamiliesMap) DeepCopyInto(out *AddressFamiliesMap) {
	*out = *in
	if in.IPv4 != nil {
		in, out := &in.IPv4, &out.IPv4
		*out = new(bool)
		**out = **in
	}
	if in.IPv6 != nil {
		in, out := &in.IPv6, &out.IPv6
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AddressFamiliesMap.
func (in *AddressFamiliesMap) DeepCopy() *AddressFamiliesMap {
	if in == nil {
		return nil
	}
	out := new(AddressFamiliesMap)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConditionSpec) DeepCopyInto(out *ConditionSpec) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.State != nil {
		in, out := &in.State, &out.State
		*out = new(bool)
		**out = **in
	}
	if in.LastUpdateTimestamp != nil {
		in, out := &in.LastUpdateTimestamp, &out.LastUpdateTimestamp
		*out = new(string)
		**out = **in
	}
	if in.LastTransitionTimestamp != nil {
		in, out := &in.LastTransitionTimestamp, &out.LastTransitionTimestamp
		*out = new(string)
		**out = **in
	}
	if in.Reason != nil {
		in, out := &in.Reason, &out.Reason
		*out = new(string)
		**out = **in
	}
	if in.Message != nil {
		in, out := &in.Message, &out.Message
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConditionSpec.
func (in *ConditionSpec) DeepCopy() *ConditionSpec {
	if in == nil {
		return nil
	}
	out := new(ConditionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in ConnectionsMap) DeepCopyInto(out *ConnectionsMap) {
	{
		in := &in
		*out = make(ConnectionsMap, len(*in))
		for key, val := range *in {
			var outVal *SwitchList
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(SwitchList)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConnectionsMap.
func (in ConnectionsMap) DeepCopy() ConnectionsMap {
	if in == nil {
		return nil
	}
	out := new(ConnectionsMap)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FieldSelectorSpec) DeepCopyInto(out *FieldSelectorSpec) {
	*out = *in
	if in.LabelKey != nil {
		in, out := &in.LabelKey, &out.LabelKey
		*out = new(string)
		**out = **in
	}
	if in.FieldRef != nil {
		in, out := &in.FieldRef, &out.FieldRef
		*out = new(corev1.ObjectFieldSelector)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FieldSelectorSpec.
func (in *FieldSelectorSpec) DeepCopy() *FieldSelectorSpec {
	if in == nil {
		return nil
	}
	out := new(FieldSelectorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeneralIPAMSpec) DeepCopyInto(out *GeneralIPAMSpec) {
	*out = *in
	if in.AddressFamily != nil {
		in, out := &in.AddressFamily, &out.AddressFamily
		*out = new(AddressFamiliesMap)
		(*in).DeepCopyInto(*out)
	}
	if in.CarrierSubnets != nil {
		in, out := &in.CarrierSubnets, &out.CarrierSubnets
		*out = new(IPAMSelectionSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.LoopbackSubnets != nil {
		in, out := &in.LoopbackSubnets, &out.LoopbackSubnets
		*out = new(IPAMSelectionSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.SouthSubnets != nil {
		in, out := &in.SouthSubnets, &out.SouthSubnets
		*out = new(IPAMSelectionSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.LoopbackAddresses != nil {
		in, out := &in.LoopbackAddresses, &out.LoopbackAddresses
		*out = new(IPAMSelectionSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeneralIPAMSpec.
func (in *GeneralIPAMSpec) DeepCopy() *GeneralIPAMSpec {
	if in == nil {
		return nil
	}
	out := new(GeneralIPAMSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPAMSelectionSpec) DeepCopyInto(out *IPAMSelectionSpec) {
	*out = *in
	if in.LabelSelector != nil {
		in, out := &in.LabelSelector, &out.LabelSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.FieldSelector != nil {
		in, out := &in.FieldSelector, &out.FieldSelector
		*out = new(FieldSelectorSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPAMSelectionSpec.
func (in *IPAMSelectionSpec) DeepCopy() *IPAMSelectionSpec {
	if in == nil {
		return nil
	}
	out := new(IPAMSelectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPAMSpec) DeepCopyInto(out *IPAMSpec) {
	*out = *in
	if in.SouthSubnets != nil {
		in, out := &in.SouthSubnets, &out.SouthSubnets
		*out = new(IPAMSelectionSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.LoopbackAddresses != nil {
		in, out := &in.LoopbackAddresses, &out.LoopbackAddresses
		*out = new(IPAMSelectionSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPAMSpec.
func (in *IPAMSpec) DeepCopy() *IPAMSpec {
	if in == nil {
		return nil
	}
	out := new(IPAMSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPAddressSpec) DeepCopyInto(out *IPAddressSpec) {
	*out = *in
	if in.ObjectReference != nil {
		in, out := &in.ObjectReference, &out.ObjectReference
		*out = new(ObjectReference)
		(*in).DeepCopyInto(*out)
	}
	if in.Address != nil {
		in, out := &in.Address, &out.Address
		*out = new(string)
		**out = **in
	}
	if in.ExtraAddress != nil {
		in, out := &in.ExtraAddress, &out.ExtraAddress
		*out = new(bool)
		**out = **in
	}
	if in.AddressFamily != nil {
		in, out := &in.AddressFamily, &out.AddressFamily
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPAddressSpec.
func (in *IPAddressSpec) DeepCopy() *IPAddressSpec {
	if in == nil {
		return nil
	}
	out := new(IPAddressSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InterfaceOverridesSpec) DeepCopyInto(out *InterfaceOverridesSpec) {
	*out = *in
	if in.PortParametersSpec != nil {
		in, out := &in.PortParametersSpec, &out.PortParametersSpec
		*out = new(PortParametersSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.IP != nil {
		in, out := &in.IP, &out.IP
		*out = make([]*AdditionalIPSpec, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(AdditionalIPSpec)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InterfaceOverridesSpec.
func (in *InterfaceOverridesSpec) DeepCopy() *InterfaceOverridesSpec {
	if in == nil {
		return nil
	}
	out := new(InterfaceOverridesSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InterfaceSpec) DeepCopyInto(out *InterfaceSpec) {
	*out = *in
	if in.PortParametersSpec != nil {
		in, out := &in.PortParametersSpec, &out.PortParametersSpec
		*out = new(PortParametersSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.MACAddress != nil {
		in, out := &in.MACAddress, &out.MACAddress
		*out = new(string)
		**out = **in
	}
	if in.Speed != nil {
		in, out := &in.Speed, &out.Speed
		*out = new(uint32)
		**out = **in
	}
	if in.IP != nil {
		in, out := &in.IP, &out.IP
		*out = make([]*IPAddressSpec, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(IPAddressSpec)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.Direction != nil {
		in, out := &in.Direction, &out.Direction
		*out = new(string)
		**out = **in
	}
	if in.Peer != nil {
		in, out := &in.Peer, &out.Peer
		*out = new(PeerSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InterfaceSpec.
func (in *InterfaceSpec) DeepCopy() *InterfaceSpec {
	if in == nil {
		return nil
	}
	out := new(InterfaceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InterfacesSpec) DeepCopyInto(out *InterfacesSpec) {
	*out = *in
	if in.Defaults != nil {
		in, out := &in.Defaults, &out.Defaults
		*out = new(PortParametersSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Overrides != nil {
		in, out := &in.Overrides, &out.Overrides
		*out = make([]*InterfaceOverridesSpec, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(InterfaceOverridesSpec)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InterfacesSpec.
func (in *InterfacesSpec) DeepCopy() *InterfacesSpec {
	if in == nil {
		return nil
	}
	out := new(InterfacesSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectReference) DeepCopyInto(out *ObjectReference) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.Namespace != nil {
		in, out := &in.Namespace, &out.Namespace
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectReference.
func (in *ObjectReference) DeepCopy() *ObjectReference {
	if in == nil {
		return nil
	}
	out := new(ObjectReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PeerInfoSpec) DeepCopyInto(out *PeerInfoSpec) {
	*out = *in
	if in.ChassisID != nil {
		in, out := &in.ChassisID, &out.ChassisID
		*out = new(string)
		**out = **in
	}
	if in.SystemName != nil {
		in, out := &in.SystemName, &out.SystemName
		*out = new(string)
		**out = **in
	}
	if in.PortID != nil {
		in, out := &in.PortID, &out.PortID
		*out = new(string)
		**out = **in
	}
	if in.PortDescription != nil {
		in, out := &in.PortDescription, &out.PortDescription
		*out = new(string)
		**out = **in
	}
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PeerInfoSpec.
func (in *PeerInfoSpec) DeepCopy() *PeerInfoSpec {
	if in == nil {
		return nil
	}
	out := new(PeerInfoSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PeerSpec) DeepCopyInto(out *PeerSpec) {
	*out = *in
	if in.ObjectReference != nil {
		in, out := &in.ObjectReference, &out.ObjectReference
		*out = new(ObjectReference)
		(*in).DeepCopyInto(*out)
	}
	if in.PeerInfoSpec != nil {
		in, out := &in.PeerInfoSpec, &out.PeerInfoSpec
		*out = new(PeerInfoSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PeerSpec.
func (in *PeerSpec) DeepCopy() *PeerSpec {
	if in == nil {
		return nil
	}
	out := new(PeerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PortParametersSpec) DeepCopyInto(out *PortParametersSpec) {
	*out = *in
	if in.Lanes != nil {
		in, out := &in.Lanes, &out.Lanes
		*out = new(uint32)
		**out = **in
	}
	if in.MTU != nil {
		in, out := &in.MTU, &out.MTU
		*out = new(uint32)
		**out = **in
	}
	if in.IPv4MaskLength != nil {
		in, out := &in.IPv4MaskLength, &out.IPv4MaskLength
		*out = new(uint32)
		**out = **in
	}
	if in.IPv6Prefix != nil {
		in, out := &in.IPv6Prefix, &out.IPv6Prefix
		*out = new(uint32)
		**out = **in
	}
	if in.FEC != nil {
		in, out := &in.FEC, &out.FEC
		*out = new(string)
		**out = **in
	}
	if in.State != nil {
		in, out := &in.State, &out.State
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PortParametersSpec.
func (in *PortParametersSpec) DeepCopy() *PortParametersSpec {
	if in == nil {
		return nil
	}
	out := new(PortParametersSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegionSpec) DeepCopyInto(out *RegionSpec) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.AvailabilityZone != nil {
		in, out := &in.AvailabilityZone, &out.AvailabilityZone
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegionSpec.
func (in *RegionSpec) DeepCopy() *RegionSpec {
	if in == nil {
		return nil
	}
	out := new(RegionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubnetSpec) DeepCopyInto(out *SubnetSpec) {
	*out = *in
	if in.Subnet != nil {
		in, out := &in.Subnet, &out.Subnet
		*out = new(ObjectReference)
		(*in).DeepCopyInto(*out)
	}
	if in.Network != nil {
		in, out := &in.Network, &out.Network
		*out = new(ObjectReference)
		(*in).DeepCopyInto(*out)
	}
	if in.CIDR != nil {
		in, out := &in.CIDR, &out.CIDR
		*out = new(string)
		**out = **in
	}
	if in.Region != nil {
		in, out := &in.Region, &out.Region
		*out = new(RegionSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.AddressFamily != nil {
		in, out := &in.AddressFamily, &out.AddressFamily
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubnetSpec.
func (in *SubnetSpec) DeepCopy() *SubnetSpec {
	if in == nil {
		return nil
	}
	out := new(SubnetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Switch) DeepCopyInto(out *Switch) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Switch.
func (in *Switch) DeepCopy() *Switch {
	if in == nil {
		return nil
	}
	out := new(Switch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Switch) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SwitchConfig) DeepCopyInto(out *SwitchConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SwitchConfig.
func (in *SwitchConfig) DeepCopy() *SwitchConfig {
	if in == nil {
		return nil
	}
	out := new(SwitchConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SwitchConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SwitchConfigList) DeepCopyInto(out *SwitchConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SwitchConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SwitchConfigList.
func (in *SwitchConfigList) DeepCopy() *SwitchConfigList {
	if in == nil {
		return nil
	}
	out := new(SwitchConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SwitchConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SwitchConfigSpec) DeepCopyInto(out *SwitchConfigSpec) {
	*out = *in
	if in.Switches != nil {
		in, out := &in.Switches, &out.Switches
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.PortsDefaults != nil {
		in, out := &in.PortsDefaults, &out.PortsDefaults
		*out = new(PortParametersSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.IPAM != nil {
		in, out := &in.IPAM, &out.IPAM
		*out = new(GeneralIPAMSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.RoutingConfigTemplate != nil {
		in, out := &in.RoutingConfigTemplate, &out.RoutingConfigTemplate
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SwitchConfigSpec.
func (in *SwitchConfigSpec) DeepCopy() *SwitchConfigSpec {
	if in == nil {
		return nil
	}
	out := new(SwitchConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SwitchConfigStatus) DeepCopyInto(out *SwitchConfigStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SwitchConfigStatus.
func (in *SwitchConfigStatus) DeepCopy() *SwitchConfigStatus {
	if in == nil {
		return nil
	}
	out := new(SwitchConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SwitchList) DeepCopyInto(out *SwitchList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Switch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SwitchList.
func (in *SwitchList) DeepCopy() *SwitchList {
	if in == nil {
		return nil
	}
	out := new(SwitchList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SwitchList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SwitchSpec) DeepCopyInto(out *SwitchSpec) {
	*out = *in
	if in.InventoryRef != nil {
		in, out := &in.InventoryRef, &out.InventoryRef
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
	if in.ConfigSelector != nil {
		in, out := &in.ConfigSelector, &out.ConfigSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Managed != nil {
		in, out := &in.Managed, &out.Managed
		*out = new(bool)
		**out = **in
	}
	if in.Cordon != nil {
		in, out := &in.Cordon, &out.Cordon
		*out = new(bool)
		**out = **in
	}
	if in.TopSpine != nil {
		in, out := &in.TopSpine, &out.TopSpine
		*out = new(bool)
		**out = **in
	}
	if in.ScanPorts != nil {
		in, out := &in.ScanPorts, &out.ScanPorts
		*out = new(bool)
		**out = **in
	}
	if in.IPAM != nil {
		in, out := &in.IPAM, &out.IPAM
		*out = new(IPAMSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Interfaces != nil {
		in, out := &in.Interfaces, &out.Interfaces
		*out = new(InterfacesSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SwitchSpec.
func (in *SwitchSpec) DeepCopy() *SwitchSpec {
	if in == nil {
		return nil
	}
	out := new(SwitchSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SwitchStatus) DeepCopyInto(out *SwitchStatus) {
	*out = *in
	if in.ConfigRef != nil {
		in, out := &in.ConfigRef, &out.ConfigRef
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
	if in.RoutingConfigTemplate != nil {
		in, out := &in.RoutingConfigTemplate, &out.RoutingConfigTemplate
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
	if in.ASN != nil {
		in, out := &in.ASN, &out.ASN
		*out = new(uint32)
		**out = **in
	}
	if in.TotalPorts != nil {
		in, out := &in.TotalPorts, &out.TotalPorts
		*out = new(uint32)
		**out = **in
	}
	if in.SwitchPorts != nil {
		in, out := &in.SwitchPorts, &out.SwitchPorts
		*out = new(uint32)
		**out = **in
	}
	if in.Role != nil {
		in, out := &in.Role, &out.Role
		*out = new(string)
		**out = **in
	}
	if in.Layer != nil {
		in, out := &in.Layer, &out.Layer
		*out = new(uint32)
		**out = **in
	}
	if in.Interfaces != nil {
		in, out := &in.Interfaces, &out.Interfaces
		*out = make(map[string]*InterfaceSpec, len(*in))
		for key, val := range *in {
			var outVal *InterfaceSpec
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(InterfaceSpec)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
	if in.Subnets != nil {
		in, out := &in.Subnets, &out.Subnets
		*out = make([]*SubnetSpec, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(SubnetSpec)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.LoopbackAddresses != nil {
		in, out := &in.LoopbackAddresses, &out.LoopbackAddresses
		*out = make([]*IPAddressSpec, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(IPAddressSpec)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.State != nil {
		in, out := &in.State, &out.State
		*out = new(string)
		**out = **in
	}
	if in.Message != nil {
		in, out := &in.Message, &out.Message
		*out = new(string)
		**out = **in
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]*ConditionSpec, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(ConditionSpec)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SwitchStatus.
func (in *SwitchStatus) DeepCopy() *SwitchStatus {
	if in == nil {
		return nil
	}
	out := new(SwitchStatus)
	in.DeepCopyInto(out)
	return out
}
