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

package rules

import (
	"fmt"
	"net/netip"

	"github.com/go-logr/logr"
	"github.com/onmetal/metal-api/common/types/base"
	"github.com/onmetal/metal-api/common/types/events"
	ipdomain "github.com/onmetal/metal-api/domain/address"
	domain "github.com/onmetal/metal-api/domain/machine"
	"github.com/onmetal/metal-api/usecase/onboarding/providers"
)

type CreateLoopbackForMachineRule struct {
	subnetExtractor   providers.SubnetExtractor
	machineExtractor  providers.MachineExtractor
	loopbackPersister providers.LoopbackPersister
	log               logr.Logger
}

func NewCreateLoopbackForMachineRule(
	subnetExtractor providers.SubnetExtractor,
	loopbackPersister providers.LoopbackPersister,
	machineExtractor providers.MachineExtractor,
	log logr.Logger,
) events.DomainEventListener[base.DomainEvent] {
	log = log.WithName("CreateLoopbackForMachineRule")
	return &CreateLoopbackForMachineRule{
		subnetExtractor:   subnetExtractor,
		loopbackPersister: loopbackPersister,
		machineExtractor:  machineExtractor,
		log:               log,
	}
}

func (c *CreateLoopbackForMachineRule) EventType() base.DomainEvent {
	return &domain.MachineCreatedDomainEvent{}
}

func (c *CreateLoopbackForMachineRule) Handle(event base.DomainEvent) {
	subnet, err := c.subnetExtractor.LoopbackIPv4Subnet()
	if err != nil {
		c.log.Info(
			"can't retrieve loopback subnet",
			"id", event.ID(),
			"error", err,
		)
		return
	}
	machine, err := c.machineExtractor.ByID(domain.NewMachineID(event.ID()))
	if err != nil {
		c.log.Info(
			"can't retrieve machine",
			"id", event.ID(),
			"error", err,
		)
		return
	}
	loopbackAddressName := fmt.Sprintf("%s-lo-ipv4", machine.UUID)
	address := ipdomain.CreateNewAddress(
		netip.Addr{},
		32,
		loopbackAddressName,
		subnet.Namespace,
		subnet.Name)
	address.SetConsumerInfo(machine.UUID, "Machine")
	if err = c.loopbackPersister.Save(address); err != nil {
		c.log.Info(
			"can't save loopback ip for machine",
			"id", event.ID(),
			"error", err)
	}
}
