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

package v1alpha1

import (
	"math"
	"math/big"
	"net"
	"strconv"
	"strings"

	gocidr "github.com/apparentlymart/go-cidr/cidr"
	subnetv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
)

type State string
type Role string
type PeerType string

const (
	StatePending      State = "Pending"
	StateInitializing State = "Initializing"
	StateDiscovery    State = "Discovery"
	StateConfiguring  State = "Configuring"
	StateReady        State = "Ready"
	StateFinished     State = "Finished"
	StateDeleting     State = "Deleting"

	LeafRole  Role = "Leaf"
	SpineRole Role = "Spine"

	MachineType PeerType = "Machine"
	SwitchType  PeerType = "Switch"
)

const (
	EmptyString = ""

	CLabelPrefix    = "switch.onmetal.de/"
	CLabelChassisId = "chassisId"

	CSonicSwitchOs     = "SONiC"
	CStationCapability = "Station"
	CNdpStateReachable = "Reachable"

	CIPv4ZeroNet             = "0.0.0.0/0"
	CIPv6ZeroNet             = "::/0"
	CIPv4AddressesPerPort    = 16
	CIPv6AddressesPerPort    = 8
	CIPv4InterfaceSubnetMask = 30
	CIPv6InterfaceSubnetMask = 127

	CNamespace = "onmetal"

	CSwitchFinalizer           = "switches.switch.onmetal.de/finalizer"
	CSwitchAssignmentFinalizer = "switchassignments.switch.onmetal.de/finalizer"
)

var LabelChassisId = CLabelPrefix + CLabelChassisId

var Lanes = map[uint32]uint8{
	1000:   1,
	10000:  1,
	25000:  1,
	40000:  4,
	50000:  2,
	100000: 4,
}

func getChassisId(nics []inventoriesv1alpha1.NICSpec) string {
	for _, nic := range nics {
		if nic.Name == "eth0" {
			return nic.MACAddress
		}
	}
	return EmptyString
}

func PrepareInterfaces(nics []inventoriesv1alpha1.NICSpec) (map[string]*InterfaceSpec, uint64) {
	result := make(map[string]*InterfaceSpec)
	switchPorts := uint64(0)
	for _, nic := range nics {
		spec := &InterfaceSpec{
			Lanes:      Lanes[nic.Speed],
			MacAddress: nic.MACAddress,
		}
		if len(nic.LLDPs) > 1 {
			data := nic.LLDPs[1]
			spec.PeerType = definePeerType(data)
			spec.PeerChassisID = data.ChassisID
			spec.PeerSystemName = data.SystemName
			spec.PeerPortID = data.PortID
			spec.PeerPortDescription = data.PortDescription
		}
		if len(nic.NDPs) > 1 {
			for _, data := range nic.NDPs {
				if data.State == CNdpStateReachable {
					spec.Ndp = true
				}
			}
		}
		result[nic.Name] = spec
		if strings.HasPrefix(nic.Name, "Ethernet") {
			switchPorts += 1
		}
	}
	return result, switchPorts
}

func definePeerType(data inventoriesv1alpha1.LLDPSpec) PeerType {
	if len(data.Capabilities) == 0 {
		return MachineType
	}
	for _, c := range data.Capabilities {
		if c == CStationCapability {
			return MachineType
		}
	}
	return SwitchType
}

func getInterfaceSubnetMaskLength(addrType subnetv1alpha1.SubnetAddressType) int {
	if addrType == subnetv1alpha1.CIPv4SubnetType {
		return CIPv4InterfaceSubnetMask
	} else {
		return CIPv6InterfaceSubnetMask
	}
}

//getMinimalVacantCIDR calculates the minimal suitable network
//from the networks list provided as argument according to the
//needed addresses count. It returns the pointer to the CIDR object.
func getMinimalVacantCIDR(vacant []subnetv1alpha1.CIDR, addressType subnetv1alpha1.SubnetAddressType, addressesCount int64) *subnetv1alpha1.CIDR {
	zeroNetString := EmptyString
	if addressType == subnetv1alpha1.CIPv4SubnetType {
		zeroNetString = CIPv4ZeroNet
	} else {
		zeroNetString = CIPv6ZeroNet
	}
	_, zeroNet, _ := net.ParseCIDR(zeroNetString)
	minSuitableNet := subnetv1alpha1.CIDRFromNet(zeroNet)
	for _, cidr := range vacant {
		if cidr.AddressCapacity().Cmp(minSuitableNet.AddressCapacity()) < 0 &&
			cidr.AddressCapacity().Cmp(new(big.Int).SetInt64(addressesCount)) >= 0 {
			minSuitableNet = &cidr
		}
	}
	return minSuitableNet
}

func getNeededMask(addrType subnetv1alpha1.SubnetAddressType, addressesCount float64) net.IPMask {
	bits := 32
	pow := 2.0
	for math.Pow(2, pow) < addressesCount {
		pow++
	}
	if addrType == subnetv1alpha1.CIPv6SubnetType {
		bits = 128
	}
	ones := bits - int(pow)
	return net.CIDRMask(ones, bits)
}

func getInterfaceSubnet(name string, network *net.IPNet, addrType subnetv1alpha1.SubnetAddressType) *net.IPNet {
	index, _ := strconv.Atoi(strings.ReplaceAll(name, "Ethernet", EmptyString))
	prefix, _ := network.Mask.Size()
	ifaceNet, _ := gocidr.Subnet(network, getInterfaceSubnetMaskLength(addrType)-prefix, index)
	return ifaceNet
}

func MacToLabel(mac string) string {
	return strings.ReplaceAll(mac, ":", "-")
}
