// /*
// Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
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
// */

package order

import (
	"context"

	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	domain "github.com/onmetal/metal-api/scheduler/domain/order"
	usecase "github.com/onmetal/metal-api/scheduler/usecase/order"
	dto "github.com/onmetal/metal-api/scheduler/usecase/order/dto"
	"github.com/onmetal/metal-api/types/common"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type InstanceExtractor struct {
	client ctrlclient.Client
}

func NewInstanceExtractor(c ctrlclient.Client) *InstanceExtractor {
	return &InstanceExtractor{
		client: c,
	}
}

func (i *InstanceExtractor) GetInstance(instanceMetadata common.Metadata) (dto.Instance, error) {
	machineInstance := &machinev1alpha2.Machine{}
	if err := i.client.
		Get(
			context.Background(),
			types.NamespacedName{
				Namespace: instanceMetadata.Namespace(),
				Name:      instanceMetadata.Name(),
			},
			machineInstance); err != nil {
		return nil, err
	}
	return &instance{
		client:   i.client,
		instance: machineInstance,
	}, nil
}

type instance struct {
	client   ctrlclient.Client
	instance *machinev1alpha2.Machine
}

func (i *instance) GetOrder() (domain.Order, error) {
	instanceMetadata := common.NewObjectMetadata(i.instance.Name, i.instance.Namespace)
	if i.instance.Status.Reservation.Reference.Name == "" {
		return nil, usecase.OrderForInstanceNotFound(instanceMetadata)
	}
	orderName := i.instance.Status.Reservation.Reference.Name
	orderNamespace := i.instance.Status.Reservation.Reference.Namespace
	orderMetadata := domain.NewOrder(orderName, orderNamespace)

	return orderMetadata, nil
}

func (i *instance) Reserved() bool {
	return i.instance.Status.Reservation.Reference.Name != ""
}

func (i *instance) CleanOrderReference() error {
	if i.instance.Status.Reservation.Reference.Name == "" {
		return nil
	}
	i.instance.Status.Reservation.Status = domain.OrderStatusAvailable
	i.instance.Status.Reservation.Reference = common.ResourceReference{}
	return i.client.Update(context.Background(), i.instance)
}
