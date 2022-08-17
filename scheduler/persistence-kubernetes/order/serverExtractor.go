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
	usecase "github.com/onmetal/metal-api/scheduler/usecase/order"
	"github.com/onmetal/metal-api/scheduler/usecase/order/dto"
	oobv1 "github.com/onmetal/oob-controller/api/v1"
)

type ServerExtractor struct {
	client provider.Client
}

func NewServerExtractor(c provider.Client) *ServerExtractor {
	return &ServerExtractor{
		client: c,
	}
}

func (e *ServerExtractor) Get(instanceMetadata base.Metadata) (dto.Server, error) {
	instance := &machinev1alpha2.Machine{}
	if err := e.client.Get(instance, instanceMetadata); err != nil {
		return nil, err
	}
	notReserved := instance.Status.Reservation.Status != "Available"

	if !instance.Status.OOB.Exist {
		return nil, usecase.ServerForInstanceNotFound(instanceMetadata)
	}

	serverName := instance.Status.OOB.Reference.Name
	serverNamespace := instance.Status.OOB.Reference.Namespace
	serverMeta := base.NewInstanceMetadata(serverName, serverNamespace)

	server := &oobv1.Machine{}
	if err := e.client.Get(server, serverMeta); err != nil {
		return nil, err
	}
	return &bareMetalServer{
		client:   e.client,
		server:   server,
		reserved: notReserved,
	}, nil
}

type bareMetalServer struct {
	client   provider.Client
	server   *oobv1.Machine
	reserved bool
}

func (b *bareMetalServer) Reserved() bool {
	return b.reserved
}

func (b *bareMetalServer) SetPowerState(powerState string) error {
	if b.server.Spec.PowerState == machinev1alpha2.MachinePowerStateOFF &&
		powerState == machinev1alpha2.MachinePowerStateOFF {
		return nil
	}
	if powerState == b.server.Spec.PowerState && powerState == machinev1alpha2.MachinePowerStateON {
		powerState = machinev1alpha2.MachinePowerStateReset
	}
	b.server.Spec.PowerState = powerState
	return b.client.Update(b.server)
}
