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

package v1beta1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

func TestSpecGetters(t *testing.T) {
	t.Parallel()
	obj := &Switch{
		Spec: SwitchSpec{
			InventoryRef: &v1.LocalObjectReference{},
			IPAM:         &IPAMSpec{},
		},
	}
	assert.Equal(t, false, obj.GetManaged())
	assert.Equal(t, false, obj.GetCordon())
	assert.Equal(t, false, obj.GetTopSpine())
	assert.Equal(t, false, obj.GetScanPorts())
	assert.Equal(t, "", obj.GetInventoryRef())
}

func TestSpecSetters(t *testing.T) {
	t.Parallel()
	obj := &Switch{
		Spec: SwitchSpec{
			InventoryRef: &v1.LocalObjectReference{},
			IPAM:         &IPAMSpec{},
		},
	}
	obj.SetCordon(true)
	obj.SetManaged(true)
	obj.SetTopSpine(true)
	obj.SetScanPorts(true)
	obj.SetInventoryRef("inventory")
	assert.Equal(t, true, obj.GetManaged())
	assert.Equal(t, true, obj.GetCordon())
	assert.Equal(t, true, obj.GetTopSpine())
	assert.Equal(t, true, obj.GetScanPorts())
	assert.Equal(t, "inventory", obj.GetInventoryRef())
}

func TestStatusGetters(t *testing.T) {
	t.Parallel()
	obj := &Switch{
		Status: SwitchStatus{
			ConfigRef: &v1.LocalObjectReference{},
		},
	}
	assert.Equal(t, uint32(255), obj.GetLayer())
	assert.Equal(t, uint32(0), obj.GetASN())
	assert.Equal(t, uint32(0), obj.GetTotalPorts())
	assert.Equal(t, uint32(0), obj.GetSwitchPorts())
	assert.Equal(t, "", obj.GetRole())
	assert.Equal(t, "", obj.GetState())
	assert.Equal(t, "", obj.GetMessage())
	assert.Nil(t, obj.GetCondition("Initialized"))
}

func TestStatusSetters(t *testing.T) {
	t.Parallel()
	obj := &Switch{
		Status: SwitchStatus{
			ConfigRef: &v1.LocalObjectReference{},
		},
	}
	obj.SetLayer(2)
	obj.SetASN(4_204_194_305)
	obj.SetTotalPorts(35)
	obj.SetSwitchPorts(32)
	obj.SetRole("leaf")
	obj.SetState("Ready")
	obj.SetMessage("test")
	obj.SetCondition("LayerAndRoleOK", true)
	assert.Equal(t, uint32(2), obj.GetLayer())
	assert.Equal(t, uint32(4_204_194_305), obj.GetASN())
	assert.Equal(t, uint32(35), obj.GetTotalPorts())
	assert.Equal(t, uint32(32), obj.GetSwitchPorts())
	assert.Equal(t, "leaf", obj.GetRole())
	assert.Equal(t, "Ready", obj.GetState())
	assert.Equal(t, "test", obj.GetMessage())
	assert.NotNil(t, obj.GetCondition("LayerAndRoleOK"))
}

func TestConditionSpecGetters(t *testing.T) {
	t.Parallel()
	obj := &ConditionSpec{}
	assert.Equal(t, "", obj.GetName())
	assert.Equal(t, false, obj.GetState())
	assert.Equal(t, "", obj.GetReason())
	assert.Equal(t, "", obj.GetMessage())
	assert.Equal(t, "", obj.GetLastTransitionTimestamp())
}

func TestConditionSpecSetters(t *testing.T) {
	t.Parallel()
	ts := time.Now()
	obj := &ConditionSpec{Name: pointer.String("sample")}
	obj.SetState(true)
	obj.SetReason("test")
	obj.SetMessage("testing")
	obj.SetLastUpdateTimestamp(ts.String())
	obj.SetLastTransitionTimestamp(ts.String())
	assert.Equal(t, "sample", obj.GetName())
	assert.Equal(t, true, obj.GetState())
	assert.Equal(t, "test", obj.GetReason())
	assert.Equal(t, "testing", obj.GetMessage())
	assert.Equal(t, ts.String(), obj.GetLastTransitionTimestamp())
	assert.Equal(t, ts.String(), obj.GetLastUpdateTimestamp())
	obj.FlushReason()
	obj.FlushMessage()
	assert.Equal(t, "", obj.GetReason())
	assert.Equal(t, "", obj.GetMessage())
}

func TestInterfaceSpecGetters(t *testing.T) {
	t.Parallel()
	obj := &InterfaceSpec{}
	assert.Equal(t, "", obj.GetMACAddress())
	assert.Equal(t, "", obj.GetDirection())
	assert.Equal(t, uint32(0), obj.GetSpeed())
}

func TestInterfaceSpecSetters(t *testing.T) {
	t.Parallel()
	obj := &InterfaceSpec{}
	obj.SetMACAddress("00:00:00:00:00:01")
	obj.SetDirection("north")
	obj.SetSpeed(1_000)
	obj.SetIPEmpty()
	obj.SetPortParametersEmpty()
	assert.Equal(t, "00:00:00:00:00:01", obj.GetMACAddress())
	assert.Equal(t, "north", obj.GetDirection())
	assert.Equal(t, uint32(1_000), obj.GetSpeed())
	assert.NotNil(t, obj.IP)
	assert.Empty(t, obj.IP)
	assert.NotNil(t, obj.PortParametersSpec)
}

func TestPortParametersSpecGetters(t *testing.T) {
	t.Parallel()
	obj := &PortParametersSpec{}
	assert.Equal(t, "", obj.GetState())
	assert.Equal(t, "", obj.GetFEC())
	assert.Equal(t, uint32(0), obj.GetLanes())
	assert.Equal(t, uint32(0), obj.GetMTU())
	assert.Equal(t, uint32(0), obj.GetIPv4MaskLength())
	assert.Equal(t, uint32(0), obj.GetIPv6Prefix())
}

func TestPortParametersSpecSetters(t *testing.T) {
	t.Parallel()
	obj := &PortParametersSpec{}
	obj.SetState("up")
	obj.SetFEC("rs")
	obj.SetLanes(4)
	obj.SetMTU(1500)
	obj.SetIPv4MaskLength(30)
	obj.SetIPv6Prefix(112)
	assert.Equal(t, "up", obj.GetState())
	assert.Equal(t, "rs", obj.GetFEC())
	assert.Equal(t, uint32(4), obj.GetLanes())
	assert.Equal(t, uint32(1500), obj.GetMTU())
	assert.Equal(t, uint32(30), obj.GetIPv4MaskLength())
	assert.Equal(t, uint32(112), obj.GetIPv6Prefix())
}

func TestIPAddressSpecGetters(t *testing.T) {
	t.Parallel()
	obj := &IPAddressSpec{}
	assert.Equal(t, "", obj.GetAddress())
	assert.Equal(t, "", obj.GetAddressFamily())
	assert.Equal(t, "", obj.GetObjectReferenceName())
	assert.Equal(t, "", obj.GetObjectReferenceNamespace())
	assert.False(t, obj.GetExtraAddress())
}

func TestIPAddressSpecSetters(t *testing.T) {
	t.Parallel()
	obj := &IPAddressSpec{}
	obj.SetAddress("100.64.0.1")
	obj.SetAddressFamily("IPv4")
	obj.SetObjectReference("sample", "default")
	obj.SetExtraAddress(true)
	assert.Equal(t, "100.64.0.1", obj.GetAddress())
	assert.Equal(t, "IPv4", obj.GetAddressFamily())
	assert.Equal(t, "sample", obj.GetObjectReferenceName())
	assert.Equal(t, "default", obj.GetObjectReferenceNamespace())
	assert.True(t, obj.GetExtraAddress())
}

func TestPeerSpecGetters(t *testing.T) {
	t.Parallel()
	obj := &PeerSpec{}
	assert.Equal(t, "", obj.GetObjectReferenceName())
	assert.Equal(t, "", obj.GetObjectReferenceNamespace())
}

func TestPeerSpecSetters(t *testing.T) {
	t.Parallel()
	obj := &PeerSpec{}
	obj.SetObjectReference("sample", "default")
	assert.Equal(t, "sample", obj.GetObjectReferenceName())
	assert.Equal(t, "default", obj.GetObjectReferenceNamespace())
}

func TestPeerInfoSpecGetters(t *testing.T) {
	t.Parallel()
	obj := &PeerInfoSpec{}
	assert.Equal(t, "", obj.GetChassisID())
	assert.Equal(t, "", obj.GetSystemName())
	assert.Equal(t, "", obj.GetType())
	assert.Equal(t, "", obj.GetPortDescription())
	assert.Equal(t, "", obj.GetPortID())
}

func TestPeerInfoSpecSetters(t *testing.T) {
	t.Parallel()
	obj := &PeerInfoSpec{}
	obj.SetChassisID("00:00:00:00:00:02")
	obj.SetSystemName("sample-host")
	obj.SetType("machine")
	obj.SetPortDescription("Eth0/1")
	obj.SetPortID("Ethernet0")
	assert.Equal(t, "00:00:00:00:00:02", obj.GetChassisID())
	assert.Equal(t, "sample-host", obj.GetSystemName())
	assert.Equal(t, "machine", obj.GetType())
	assert.Equal(t, "Eth0/1", obj.GetPortDescription())
	assert.Equal(t, "Ethernet0", obj.GetPortID())
}

func TestSubnetSpecGetters(t *testing.T) {
	t.Parallel()
	obj := &SubnetSpec{}
	assert.Equal(t, "", obj.GetObjectReferenceName())
	assert.Equal(t, "", obj.GetObjectReferenceNamespace())
	assert.Equal(t, "", obj.GetCIDR())
	assert.Equal(t, "", obj.GetAddressFamily())
}

func TestSubnetSpecSetters(t *testing.T) {
	t.Parallel()
	obj := &SubnetSpec{}
	obj.SetObjectReference("sample", "default")
	obj.SetCIDR("100.64.0.0/24")
	obj.SetAddressFamily("IPv4")
	assert.Equal(t, "sample", obj.GetObjectReferenceName())
	assert.Equal(t, "default", obj.GetObjectReferenceNamespace())
	assert.Equal(t, "100.64.0.0/24", obj.GetCIDR())
	assert.Equal(t, "IPv4", obj.GetAddressFamily())
}

func TestInterfaceOverridesSpecFuncs(t *testing.T) {
	t.Parallel()
	obj := &InterfaceOverridesSpec{}
	assert.Equal(t, "", obj.GetName())
	obj.SetName("Ethernet1")
	assert.Equal(t, "Ethernet1", obj.GetName())
}

func TestAdditionalIPSpecFuncs(t *testing.T) {
	t.Parallel()
	obj := &AdditionalIPSpec{}
	assert.Equal(t, "", obj.GetAddress())
	obj.Address = pointer.String("fe80:abcd::1")
	assert.Equal(t, "fe80:abcd::1", obj.GetAddress())
}

func TestAddressFamiliesMapFuncs(t *testing.T) {
	t.Parallel()
	obj := &AddressFamiliesMap{}
	assert.False(t, obj.GetIPv4())
	assert.False(t, obj.GetIPv6())
	obj.IPv4 = pointer.Bool(true)
	obj.IPv6 = pointer.Bool(true)
	assert.True(t, obj.GetIPv4())
	assert.True(t, obj.GetIPv6())
}

func TestFieldSelectorSpecFuncs(t *testing.T) {
	t.Parallel()
	obj := FieldSelectorSpec{}
	assert.Equal(t, "", obj.GetLabelKey())
	obj.LabelKey = pointer.String("switch.onmetal.de/type")
	assert.Equal(t, "switch.onmetal.de/type", obj.GetLabelKey())
}
