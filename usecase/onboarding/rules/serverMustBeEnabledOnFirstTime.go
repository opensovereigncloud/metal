// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"github.com/go-logr/logr"
	"github.com/ironcore-dev/metal/common/types/base"
	infradomain "github.com/ironcore-dev/metal/domain/infrastructure"
	domain "github.com/ironcore-dev/metal/domain/inventory"
	"github.com/ironcore-dev/metal/usecase/onboarding/providers"
)

type ServerMustBeEnabledOnFirstTimeRule struct {
	inventoryExtractor providers.InventoryExtractor
	serverExecutor     providers.ServerExecutor
	log                logr.Logger
}

func NewServerMustBeEnabledOnFirstTimeRule(
	serverExecutor providers.ServerExecutor,
	inventoryExtractor providers.InventoryExtractor,
	log logr.Logger,
) *ServerMustBeEnabledOnFirstTimeRule {
	return &ServerMustBeEnabledOnFirstTimeRule{
		inventoryExtractor: inventoryExtractor,
		serverExecutor:     serverExecutor,
		log:                log,
	}
}

func (c *ServerMustBeEnabledOnFirstTimeRule) EventType() base.DomainEvent {
	return &domain.InventoryCreatedDomainEvent{}
}

func (c *ServerMustBeEnabledOnFirstTimeRule) Handle(event base.DomainEvent) {
	inv, err := c.inventoryExtractor.ByID(domain.NewInventoryID(event.ID()))
	if err != nil {
		c.log.Info("")
	}
	serverInfo := infradomain.Server{UUID: inv.UUID}
	if err = c.serverExecutor.Enable(serverInfo); err != nil {
		c.log.Info("")
	}
}
