// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
