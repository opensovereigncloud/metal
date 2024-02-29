// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	domain "github.com/ironcore-dev/metal/domain/address"
	"github.com/ironcore-dev/metal/usecase/onboarding/providers"
	ipam "github.com/onmetal/ipam/api/v1alpha1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ipV4LoopbackBitLen = 32
	ipV6LoopbackBitLen = 128
)

type LoopbackAddressRepository struct {
	tryCount int
	client   ctrlclient.Client
}

func NewLoopbackAddressRepository(
	client ctrlclient.Client,
) *LoopbackAddressRepository {
	return &LoopbackAddressRepository{
		tryCount: 1,
		client:   client,
	}
}

func (l *LoopbackAddressRepository) Save(address domain.Address) error {
	ip := prepareIP(address)
	return l.
		client.
		Create(
			context.Background(),
			ip)
}

func (l *LoopbackAddressRepository) Try(times int) providers.LoopbackAddressExtractor {
	l.tryCount = times
	return l
}
func (l *LoopbackAddressRepository) IPv4ByMachineUUID(
	uuid string,
) (domain.Address, error) {
	uuid = fmt.Sprintf("%s-lo-ipv4", uuid)
	address, err := l.tryByUUID(uuid)
	if err != nil {
		return domain.Address{}, err
	}
	addr, parseErr := netip.ParseAddr(address.Status.Reserved.String())
	if parseErr != nil {
		return domain.Address{}, parseErr
	}
	return domain.CreateNewAddress(
		addr,
		prefixBitsFromType(addr),
		address.Name,
		address.Namespace,
		address.Spec.Subnet.Name,
	), nil
}

func (l *LoopbackAddressRepository) IPv6ByMachineUUID(
	uuid string,
) (domain.Address, error) {
	uuid = fmt.Sprintf("%s-lo-ipv6", uuid)
	address, err := l.tryByUUID(uuid)
	if err != nil {
		return domain.Address{}, err
	}
	addr, parseErr := netip.ParseAddr(address.Status.Reserved.String())
	if parseErr != nil {
		return domain.Address{}, parseErr
	}
	return domain.CreateNewAddress(
		addr,
		prefixBitsFromType(addr),
		address.Name,
		address.Namespace,
		address.Spec.Subnet.Name,
	), nil
}

func (l *LoopbackAddressRepository) tryByUUID(uuid string) (*ipam.IP, error) {
	var err error
	var address *ipam.IP
	for i := 0; i < l.tryCount; i++ {
		address, err = l.getIPAMLoopbackIP(uuid)
		if err != nil {
			continue
		}
		if address.Status.Reserved != nil {
			break
		}
		err = errIPNotSet
		time.Sleep(2 * time.Second)
	}
	return address, err
}

func (l *LoopbackAddressRepository) getIPAMLoopbackIP(
	uuid string,
) (*ipam.IP, error) {
	ipamListData := &ipam.IPList{}
	err := l.
		client.
		List(
			context.Background(),
			ipamListData,
			ctrlclient.MatchingFields{
				"metadata.name": uuid,
			})
	if len(ipamListData.Items) == 0 {
		return nil, fmt.Errorf("%s: %s", errIPNotFound, err)
	}
	return &ipamListData.Items[0], err
}

func prepareIP(address domain.Address) *ipam.IP {
	return &ipam.IP{
		ObjectMeta: meta.ObjectMeta{
			Name:      address.Name,
			Namespace: address.Namespace,
		},
		Spec: ipam.IPSpec{
			Consumer: &ipam.ResourceReference{
				Kind: address.Consumer.Type,
				Name: address.Consumer.Name,
			},
			Subnet: core.LocalObjectReference{
				Name: address.Subnet,
			},
		},
	}
}

func prefixBitsFromType(addr netip.Addr) int {
	if addr.Is4() {
		return ipV4LoopbackBitLen
	}
	return ipV6LoopbackBitLen
}
