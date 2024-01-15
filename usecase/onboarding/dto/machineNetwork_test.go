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

package dto_test

import (
	"net/netip"
	"testing"

	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
	"github.com/stretchr/testify/assert"
)

func TestFilterAndGenerateAddresses4Success(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	ipv4Address, err := netip.ParsePrefix("192.168.1.1/30")
	a.Nil(err)
	switchInterfaces := switchInterfacesWithIPAddress(ipv4Address)

	addressSpec := dto.FilterAndGenerateAddresses(
		switchInterfaces,
		dto.IsIPv4,
	)
	a.Equal(1, len(addressSpec))
	a.True(addressSpec[0].Addr().Is4())
	a.Equal(ipv4Address.Addr().Next().String(), addressSpec[0].Addr().String())
}

func TestFilterAndGenerateAddresses6Success(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	ipv6Address, err := netip.ParsePrefix("2001:0db8:85a3:0000:0000:8a2e:0370:7334/128")
	a.Nil(err)
	switchInterfaces := switchInterfacesWithIPAddress(ipv6Address)

	addressSpec := dto.FilterAndGenerateAddresses(
		switchInterfaces,
		dto.IsIPv6,
	)
	a.Equal(1, len(addressSpec))
	a.True(addressSpec[0].Addr().Is6())
	a.Equal(ipv6Address.Addr().Next().String(), addressSpec[0].Addr().String())
}

func TestFilterAndGenerateAddressesWrongIPType(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	ipv6Address, err := netip.ParsePrefix("2001:0db8:85a3:0000:0000:8a2e:0370:7334/128")
	a.Nil(err)
	switchInterfaces := switchInterfacesWithIPAddress(ipv6Address)

	addressSpec := dto.FilterAndGenerateAddresses(
		switchInterfaces,
		dto.IsIPv4,
	)
	a.Equal(0, len(addressSpec))
}

func switchInterfacesWithIPAddress(address ...netip.Prefix) dto.Interface {
	prefixes := make([]netip.Prefix, len(address))
	for a := range address {
		spec := address[a]
		prefixes[a] = spec
	}
	return dto.Interface{IP: prefixes}
}
