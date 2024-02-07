// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package domain_test

import (
	"net/netip"
	"testing"

	domain "github.com/ironcore-dev/metal/domain/address"
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
