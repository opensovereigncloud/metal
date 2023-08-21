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

	ipam "github.com/onmetal/ipam/api/v1alpha1"
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	"k8s.io/apimachinery/pkg/labels"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type SubnetRepository struct {
	client           ctrlclient.Client
	subnetLabelValue string
}

func NewSubnetRepository(
	client ctrlclient.Client,
	subnetLabelValue string,
) *SubnetRepository {
	return &SubnetRepository{
		client:           client,
		subnetLabelValue: subnetLabelValue,
	}
}

func (s *SubnetRepository) LoopbackIPv4Subnet() (dto.SubnetInfo, error) {
	loopbackLabelOptions := &ctrlclient.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{"loopback": s.subnetLabelValue})}
	subnetList, err := s.getSubnet(loopbackLabelOptions)
	if err != nil {
		return dto.SubnetInfo{}, err
	}
	for s := range subnetList.Items {
		if subnetList.Items[s].Status.Type != ipam.CIPv4SubnetType {
			continue
		}
		return dto.SubnetInfo{
			Name:      subnetList.Items[s].Name,
			Namespace: subnetList.Items[s].Namespace,
		}, nil
	}
	return dto.SubnetInfo{}, errNotFound
}

func (s *SubnetRepository) getSubnet(options ctrlclient.ListOption) (*ipam.SubnetList, error) {
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
	return obj, nil
}
