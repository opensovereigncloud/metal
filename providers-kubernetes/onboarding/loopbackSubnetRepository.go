// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	"context"
	"fmt"

	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
	"github.com/ironcore-dev/metal/usecase/onboarding/providers"
	ipam "github.com/onmetal/ipam/api/v1alpha1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type LoopbackSubnetRepository struct {
	client           ctrlclient.Client
	subnetLabelValue string
}

func NewLoopbackSubnetRepository(
	client ctrlclient.Client,
	subnetLabelValue string,
) *LoopbackSubnetRepository {
	return &LoopbackSubnetRepository{
		client:           client,
		subnetLabelValue: subnetLabelValue,
	}
}

func (s *LoopbackSubnetRepository) Save(info dto.SubnetInfo) error {
	subnet := prepareSubnet(info)
	return s.
		client.
		Create(
			context.Background(),
			subnet,
		)
}

func (s *LoopbackSubnetRepository) ByType(ipType string) (dto.SubnetInfo, error) {
	loopbackLabelOptions := &ctrlclient.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{
			"loopback": s.subnetLabelValue})}

	ipamIPAddressType := ipAddressType(ipType)
	subnetByType, err := s.getSubnetByType(loopbackLabelOptions, ipamIPAddressType)
	if err != nil {
		return dto.SubnetInfo{}, err
	}
	return dto.NewSubnetInfo(
		subnetByType.Name,
		subnetByType.Namespace,
		prefixBitsDereference(subnetByType.Spec.PrefixBits),
		subnetByType.Spec.ParentSubnet.Name,
	), nil
}

func (s *LoopbackSubnetRepository) IPv6ByName(name string) (dto.SubnetInfo, error) {
	name = fmt.Sprintf("%s-lo-ipv6", name)
	nameOptions := ctrlclient.MatchingFields{
		"metadata.name": name,
	}
	subnetByType, err := s.getSubnetByType(nameOptions, ipam.CIPv6SubnetType)
	if err != nil {
		return dto.SubnetInfo{}, err
	}
	return dto.NewSubnetInfo(
		subnetByType.Name,
		subnetByType.Namespace,
		int(*subnetByType.Spec.PrefixBits),
		subnetByType.Spec.ParentSubnet.Name,
	), nil
}

func (s *LoopbackSubnetRepository) getSubnetByType(
	options ctrlclient.ListOption,
	subnetAddressType ipam.SubnetAddressType,
) (*ipam.Subnet, error) {
	obj := &ipam.SubnetList{}
	if err := s.
		client.
		List(
			context.Background(),
			obj,
			options,
		); err != nil {
		return nil, err
	}

	if len(obj.Items) == 0 {
		return nil, errNotFound
	}

	for s := range obj.Items {
		if obj.Items[s].Status.Type != subnetAddressType {
			continue
		}
		return &obj.Items[s], nil
	}
	return nil, errNotFound
}

func prepareSubnet(subnetInfo dto.SubnetInfo) *ipam.Subnet {
	prefix := byte(subnetInfo.Prefix)
	return &ipam.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      subnetInfo.Name,
			Namespace: subnetInfo.Namespace,
		},
		Spec: ipam.SubnetSpec{
			PrefixBits: &prefix,
			ParentSubnet: core.LocalObjectReference{
				Name: subnetInfo.ParentSubnetName,
			},
		},
	}
}

func ipAddressType(ipType string) ipam.SubnetAddressType {
	if ipType == providers.IPv4 {
		return ipam.CIPv4SubnetType
	}
	return ipam.CIPv6SubnetType
}

func prefixBitsDereference(bits *byte) int {
	if bits != nil {
		return int(*bits)
	}
	return 0
}
