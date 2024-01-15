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

package domain

import "github.com/ironcore-dev/metal/common/types/base"

type InventoryCreatedDomainEvent struct {
	id InventoryID
}

type InventoryFlavorUpdatedDomainEvent struct {
	id InventoryID
}

func (m *InventoryCreatedDomainEvent) ID() string         { return m.id.String() }
func (m *InventoryFlavorUpdatedDomainEvent) ID() string   { return m.id.String() }
func (m *InventoryCreatedDomainEvent) Type() string       { return "inventory created" }
func (m *InventoryFlavorUpdatedDomainEvent) Type() string { return "inventory flavor update" }

func NewInventoryCreatedDomainEvent(
	inventoryID InventoryID,
) base.DomainEvent {
	return &InventoryCreatedDomainEvent{
		id: inventoryID,
	}
}

func NewInventoryFlavorUpdatedDomainEvent(
	inventoryID InventoryID,
) base.DomainEvent {
	return &InventoryFlavorUpdatedDomainEvent{
		id: inventoryID,
	}
}
