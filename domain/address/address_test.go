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

package domain_test

import (
	"net/netip"
	"testing"

	domain "github.com/onmetal/metal-api/domain/address"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewSuccess(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	address, err := netip.ParsePrefix("192.168.0.1/24")
	a.Nil(err)

	name := "test"
	result := domain.CreateNewAddress(address.Addr(), address.Bits(), name, "", "")
	a.Equal(address, result.Prefix)
	a.Equal(name, result.Name)
}

func TestCreateNewZeroPrefix(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	address, err := netip.ParseAddr("192.168.0.1")
	a.Nil(err)

	name := "test"
	result := domain.CreateNewAddress(address, 0, name, "", "")
	a.Equal(address, result.Prefix.Addr())
	a.Equal(name, result.Name)
}
