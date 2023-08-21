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

package bgp

import (
	"math"
	"net/netip"
)

const ASNBase uint32 = 4_200_000_00

func CalculateAutonomousSystemNumberFromAddress(
	address netip.Addr,
) uint32 {
	if address.String() == "" {
		return 0
	}
	asn := ASNBase
	addr := address.As16()
	asn += uint32(addr[13]) * uint32(math.Pow(2, 16))
	asn += uint32(addr[14]) * uint32(math.Pow(2, 8))
	asn += uint32(addr[15])
	return asn
}
