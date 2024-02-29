// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	domain "github.com/ironcore-dev/metal/domain/inventory"
)

type InventoryPersister interface {
	Save(inventory domain.Inventory) error
}
