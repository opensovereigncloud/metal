// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package dto

import (
	"net/netip"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
)

type SwitchInfo struct {
	Name           string
	Lanes          uint32
	InterfacesInfo map[string]Interface
	Interfaces     *metalv1alpha4.InterfaceSpec
}

type Interface struct {
	IP []netip.Prefix
}

func (s *SwitchInfo) AddSwitchInfoToMachineInterfaces(
	machineInterface metalv1alpha4.Interface,
) metalv1alpha4.Interface {
	machineInterface.Lanes = s.Lanes
	machineInterface.Unknown = false
	machineInterface.SwitchReference = &metalv1alpha4.ResourceReference{Kind: "NetworkSwitch", Name: s.Name}
	switchInterface := s.InterfacesInfo[machineInterface.Peer.LLDPPortDescription]
	machineInterface.Addresses = CalculateMachineAddressFromSwitchInterface(switchInterface)
	return machineInterface
}

func CalculateMachineAddressFromSwitchInterface(
	switchNIC Interface,
) metalv1alpha4.Addresses {
	addresses4 := FilterAndGenerateAddresses(switchNIC, IsIPv4)
	addresses6 := FilterAndGenerateAddresses(switchNIC, IsIPv6)
	return metalv1alpha4.Addresses{
		IPv4: addresses4,
		IPv6: addresses6,
	}
}

func FilterAndGenerateAddresses(
	switchInterface Interface,
	ipTypeFunc func(addr netip.Addr) bool,
) []metalv1alpha4.IPAddrSpec {
	var addresses []metalv1alpha4.IPAddrSpec
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
