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
