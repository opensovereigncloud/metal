// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package bgp_test

import (
	"net/netip"
	"testing"

	"github.com/ironcore-dev/metal/pkg/network/bgp"
	"github.com/stretchr/testify/assert"
)

func TestCalculateAutonomousSystemNumberFromAddress(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	address := netip.AddrFrom4([4]byte{0, 0, 0, 1})
	var expectedASN uint32 = 4_200_000_01
	asn := bgp.CalculateAutonomousSystemNumberFromAddress(address)
	a.Equal(expectedASN, asn)
}
