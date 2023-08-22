// Copyright 2023 OnMetal authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dto

import (
	"net/netip"

	machine "github.com/onmetal/metal-api/apis/machine/v1alpha3"
	switches "github.com/onmetal/metal-api/apis/switch/v1beta1"
)

type SwitchInfo struct {
	Name           string
	Lanes          uint32
	InterfacesInfo map[string]Interface
	Interfaces     *switches.InterfaceSpec
}

type Interface struct {
	IP []netip.Prefix
}

func (s *SwitchInfo) AddSwitchInfoToMachineInterfaces(
	machineInterface machine.Interface,
) machine.Interface {
	machineInterface.Lanes = s.Lanes
	machineInterface.Unknown = false
	machineInterface.SwitchReference = &machine.ResourceReference{Kind: "Switch", Name: s.Name}
	switchInterface := s.InterfacesInfo[machineInterface.Peer.LLDPPortDescription]
	machineInterface.Addresses = CalculateMachineAddressFromSwitchInterface(switchInterface)
	return machineInterface
}

func CalculateMachineAddressFromSwitchInterface(
	switchNIC Interface,
) machine.Addresses {
	addresses4 := FilterAndGenerateAddresses(switchNIC, IsIPv4)
	addresses6 := FilterAndGenerateAddresses(switchNIC, IsIPv6)
	return machine.Addresses{
		IPv4: addresses4,
		IPv6: addresses6,
	}
}

func FilterAndGenerateAddresses(
	switchInterface Interface,
	ipTypeFunc func(addr netip.Addr) bool,
) []machine.IPAddressSpec {
	var addresses []machine.IPAddressSpec
	for _, ip := range switchInterface.IP {
		if !ipTypeFunc(ip.Addr()) {
			continue
		}
		addresses = appendNewAddress(addresses, ip.Addr().Next(), ip.Addr().BitLen())
	}
	return addresses
}

func IsIPv4(addr netip.Addr) bool {
	return addr.Is4()
}

func IsIPv6(addr netip.Addr) bool {
	return addr.Is6()
}
