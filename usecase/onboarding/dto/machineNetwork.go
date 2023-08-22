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

	inventories "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machine "github.com/onmetal/metal-api/apis/machine/v1alpha3"
	domain "github.com/onmetal/metal-api/domain/machine"
)

func ToMachineInterfaces(nics []inventories.NICSpec) []machine.Interface {
	interfaces := make([]machine.Interface, 0, len(nics))
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
	nicsSpec *inventories.NICSpec,
) machine.Interface {
	if moreThenOneNeighbour(nicsSpec) {
		return machine.Interface{
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

func moreThenOneNeighbour(nicsSpec *inventories.NICSpec) bool {
	return len(nicsSpec.LLDPs) != 1
}

func appendNewAddress(
	addresses []machine.IPAddressSpec,
	ip netip.Addr,
	bits int,
) []machine.IPAddressSpec {
	return append(
		addresses,
		machine.IPAddressSpec{
			Prefix: netip.PrefixFrom(ip, bits),
		})
}
