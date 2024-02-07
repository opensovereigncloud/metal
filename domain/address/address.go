// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"net/netip"
)

type Address struct {
	Consumer

	Prefix    netip.Prefix
	Name      string
	Namespace string
	Subnet    string
}

func CreateNewAddress(
	address netip.Addr,
	prefixBits int,
	name string,
	namespace string,
	subnetName string,
) Address {
	return Address{
		Prefix:    netip.PrefixFrom(address, prefixBits),
		Name:      name,
		Namespace: namespace,
		Subnet:    subnetName,
	}
}

func (a *Address) SetConsumerInfo(name, consumerType string) {
	a.Consumer = Consumer{
		Name: name,
		Type: consumerType,
	}
}
