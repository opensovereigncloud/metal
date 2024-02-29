// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package dto

import (
	"net/netip"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	domain "github.com/ironcore-dev/metal/domain/machine"
)

func ToMachineInterfaces(nics []metalv1alpha4.NICSpec) []metalv1alpha4.Interface {
	interfaces := make([]metalv1alpha4.Interface, 0, len(nics))
	for nic := range nics {
		if len(nics[nic].LLDPs) == 0 {
			continue
		}
		interfaces = append(
			interfaces,
			toMachineInterface(&nics[nic]),
		)
	}
	return interfaces
}

func toMachineInterface(
	nicsSpec *metalv1alpha4.NICSpec,
) metalv1alpha4.Interface {
	if moreThenOneNeighbour(nicsSpec) {
		return metalv1alpha4.Interface{
			Name:    nicsSpec.Name,
			Unknown: true,
		}
	}
	return domain.NewMachineInterface(
		nicsSpec.Name,
		nicsSpec.LLDPs[0].SystemName,
		nicsSpec.LLDPs[0].ChassisID,
		nicsSpec.LLDPs[0].PortID,
		nicsSpec.LLDPs[0].PortDescription,
	)
}

func moreThenOneNeighbour(nicsSpec *metalv1alpha4.NICSpec) bool {
	return len(nicsSpec.LLDPs) != 1
}

func appendNewAddress(
	addresses []metalv1alpha4.IPAddrSpec,
	ip netip.Addr,
	bits int,
) []metalv1alpha4.IPAddrSpec {
	return append(
		addresses,
		metalv1alpha4.IPAddrSpec{
			Prefix: netip.PrefixFrom(ip, bits),
		})
}
