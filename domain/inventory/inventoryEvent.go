// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
