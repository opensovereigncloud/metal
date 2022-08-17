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
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/common/types/base"
	"github.com/onmetal/metal-api/pkg/provider"
	domain "github.com/onmetal/metal-api/scheduler/domain/order"
	usecase "github.com/onmetal/metal-api/scheduler/usecase/order"
	dto2 "github.com/onmetal/metal-api/scheduler/usecase/order/dto"
	oobv1 "github.com/onmetal/oob-controller/api/v1"
)

type InstanceExtractor struct {
	client provider.Client
}

func NewInstanceExtractor(c provider.Client) *InstanceExtractor {
	return &InstanceExtractor{
		client: c,
	}
}

func (i *InstanceExtractor) GetInstance(instanceMetadata base.Metadata) (dto2.Instance, error) {
	machineInstance := &machinev1alpha2.Machine{}
	if err := i.client.Get(machineInstance, instanceMetadata); err != nil {
		return nil, err
	}
	return &instance{
		client:   i.client,
		instance: machineInstance,
	}, nil
}

type instance struct {
	client   provider.Client
	instance *machinev1alpha2.Machine
}

func (i *instance) GetOrder() (domain.Order, error) {
	instanceMetadata := base.NewInstanceMetadata(i.instance.Name, i.instance.Namespace)
	if i.instance.Status.Reservation.Reference == nil {
		return nil, usecase.OrderForInstanceNotFound(instanceMetadata)
	}
	orderName := i.instance.Status.Reservation.Reference.Name
	orderNamespace := i.instance.Status.Reservation.Reference.Namespace
	orderMetadata := domain.NewOrder(orderName, orderNamespace)

	return orderMetadata, nil
}

func (i *instance) GetServer() (dto2.Server, error) {
	instanceMetadata := base.NewInstanceMetadata(i.instance.Name, i.instance.Namespace)
	if !i.instance.Status.OOB.Exist {
		return nil, usecase.ServerForInstanceNotFound(instanceMetadata)
	}

	server := &oobv1.Machine{}

	serverName := i.instance.Status.OOB.Reference.Name
	serverNamespace := i.instance.Status.OOB.Reference.Namespace
	serverMetadata := base.NewInstanceMetadata(serverName, serverNamespace)

	if err := i.client.Get(server, serverMetadata); err != nil {
		return nil, err
	}
	return &bareMetalServer{
		client:   i.client,
		server:   server,
		reserved: false,
	}, nil
}

func (i *instance) CleanOrderReference() error {
	if i.instance.Status.Reservation.Reference == nil {
		return nil
	}
	i.instance.Status.Reservation.Status = domain.OrderStatusAvailable
	i.instance.Status.Reservation.Reference = nil

	return i.client.Update(i.instance)
}
