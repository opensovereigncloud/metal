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
	"github.com/go-logr/logr"
	"github.com/ironcore-dev/metal/common/types/base"
	domain "github.com/ironcore-dev/metal/domain/inventory"
	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
	"github.com/ironcore-dev/metal/usecase/onboarding/providers"
)

type CreateIPv6SubnetFromParentForInventoryRule struct {
	loopbackSubnetExtractor providers.LoopbackSubnetExtractor
	subnetPersister         providers.SubnetPersister
	inventoryExtractor      providers.InventoryExtractor
	log                     logr.Logger
}

func NewCreateIPv6SubnetFromParentForInventoryRule(
	loopbackSubnetExtractor providers.LoopbackSubnetExtractor,
	subnetPersister providers.SubnetPersister,
	inventoryExtractor providers.InventoryExtractor,
	log logr.Logger,
) *CreateIPv6SubnetFromParentForInventoryRule {
	log = log.WithName("CreateIPv6SubnetFromParentForInventoryRule")
	return &CreateIPv6SubnetFromParentForInventoryRule{
		loopbackSubnetExtractor: loopbackSubnetExtractor,
		subnetPersister:         subnetPersister,
		inventoryExtractor:      inventoryExtractor,
		log:                     log,
	}
}

func (c *CreateIPv6SubnetFromParentForInventoryRule) EventType() base.DomainEvent {
	return &domain.InventoryCreatedDomainEvent{}
}

func (c *CreateIPv6SubnetFromParentForInventoryRule) Handle(event base.DomainEvent) {
	parentSubnet, err := c.loopbackSubnetExtractor.ByType(providers.IPv6)
	if err != nil {
		c.log.Info(
			"can't retrieve loopback subnet",
			"id", event.ID(),
			"error", err,
		)
		return
	}
	inv, err := c.inventoryExtractor.ByID(domain.NewInventoryID(event.ID()))
	if err != nil {
		c.log.Info(
			"can't retrieve inventory",
			"id", event.ID(),
			"error", err,
		)
		return
	}
	subnet := dto.SubnetInfo{
		Name:             inv.UUID,
		Namespace:        parentSubnet.Namespace,
		Prefix:           64,
		ParentSubnetName: parentSubnet.Name,
	}
	if err = c.subnetPersister.Save(subnet); err != nil {
		c.log.Info(
			"can't create loopback subnet",
			"id", event.ID(),
			"error", err,
		)
	}
}
