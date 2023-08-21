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

package providers

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	ipam "github.com/onmetal/ipam/api/v1alpha1"
	"github.com/onmetal/metal-api/common/types/common"
	"github.com/onmetal/metal-api/usecase/onboarding/providers"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ipV4LoopbackBitLen = 32
	ipV6LoopbackBitLen = 128
)

type LoopbackRepository struct {
	tryCount int
	client   ctrlclient.Client
}

func NewLoopbackRepository(
	client ctrlclient.Client,
) *LoopbackRepository {
	return &LoopbackRepository{
		tryCount: 1,
		client:   client,
	}
}

func (l *LoopbackRepository) Save(address common.Address) error {
	ip := prepareIP(address)
	return l.
		client.
		Create(
			context.Background(),
			ip)
}

func (l *LoopbackRepository) Try(times int) providers.LoopbackExtractor {
	l.tryCount = times
	return l
}
func (l *LoopbackRepository) IPv4ByMachineUUID(
	uuid string,
) (common.Address, error) {
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
	if err != nil {
		return common.Address{}, err
	}
	addr, parseErr := netip.ParseAddr(address.Status.Reserved.String())
	if parseErr != nil {
		return common.Address{}, parseErr
	}
	return common.CreateNewAddress(
		addr,
		prefixBitsFromType(addr),
		address.Name,
		address.Namespace,
		address.Spec.Subnet.Name,
	), nil
}

func (l *LoopbackRepository) getIPAMLoopbackIP(
	uuid string,
) (*ipam.IP, error) {
	ipamListData := &ipam.IPList{}
	err := l.
		client.
		List(
			context.Background(),
			ipamListData,
			ctrlclient.MatchingFields{
				"metadata.name": fmt.Sprintf("%s-lo-ipv4", uuid),
			})
	if len(ipamListData.Items) == 0 {
		return nil, fmt.Errorf("%s: %s", errIPNotFound, err)
	}
	return &ipamListData.Items[0], err
}

func prepareIP(address common.Address) *ipam.IP {
	return &ipam.IP{
		ObjectMeta: meta.ObjectMeta{
			Name:      address.Name,
			Namespace: address.Namespace,
		},
		Spec: ipam.IPSpec{
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
