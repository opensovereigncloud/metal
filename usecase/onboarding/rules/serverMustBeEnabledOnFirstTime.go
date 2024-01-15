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
