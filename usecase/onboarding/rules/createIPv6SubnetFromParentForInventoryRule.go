// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
