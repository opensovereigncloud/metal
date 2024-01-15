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
