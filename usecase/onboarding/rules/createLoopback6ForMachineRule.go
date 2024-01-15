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
	"github.com/ironcore-dev/metal/common/types/base"
	"github.com/ironcore-dev/metal/common/types/events"
	ipdomain "github.com/ironcore-dev/metal/domain/address"
	domain "github.com/ironcore-dev/metal/domain/inventory"
	"github.com/ironcore-dev/metal/usecase/onboarding/providers"
)

type CreateLoopback6ForMachineRule struct {
	loopbackSubnetExtractor providers.LoopbackSubnetExtractor
	inventoryExtractor      providers.InventoryExtractor
	addressPersister        providers.AddressPersister
	log                     logr.Logger
}

func NewCreateLoopback6ForMachineRule(
	subnetExtractor providers.LoopbackSubnetExtractor,
	addressPersister providers.AddressPersister,
	inventoryExtractor providers.InventoryExtractor,
	log logr.Logger,
) events.DomainEventListener[base.DomainEvent] {
	log = log.WithName("CreateLoopback6ForMachineRule")
	return &CreateLoopback6ForMachineRule{
		loopbackSubnetExtractor: subnetExtractor,
		addressPersister:        addressPersister,
		inventoryExtractor:      inventoryExtractor,
		log:                     log,
	}
}

func (c *CreateLoopback6ForMachineRule) EventType() base.DomainEvent {
	return &domain.InventoryFlavorUpdatedDomainEvent{}
}

func (c *CreateLoopback6ForMachineRule) Handle(event base.DomainEvent) {
	inventory, err := c.inventoryExtractor.ByID(domain.NewInventoryID(event.ID()))
	if err != nil {
		c.log.Info(
			"can't retrieve inventory",
			"id", event.ID(),
			"error", err,
		)
		return
	}
	if !inventory.IsMachine() {
		return
	}

	if err = c.CreateLoopbackIP(inventory.UUID); err != nil {
		c.log.Info(
			"can't create ipv6 loopback for machine",
			"id", event.ID(),
			"error", err,
		)
	}
}

func (c *CreateLoopback6ForMachineRule) CreateLoopbackIP(
	uuid string,
) error {
	subnet, err := c.loopbackSubnetExtractor.IPv6ByName(uuid)
	if err != nil {
		return err
	}
	addressName := fmt.Sprintf("%s-lo-ipv6", uuid)
	address := ipdomain.CreateNewAddress(
		netip.Addr{},
		128,
		addressName,
		subnet.Namespace,
		subnet.Name)
	address.SetConsumerInfo(uuid, "Machine")
	return c.
		addressPersister.
		Save(address)
}
