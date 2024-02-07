// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	domain "github.com/ironcore-dev/metal/domain/inventory"
)

//go:generate mockery --name InventoryExtractor
type InventoryExtractor interface {
	ByUUID(uuid string) (domain.Inventory, error)
	ByID(id domain.InventoryID) (domain.Inventory, error)
}
