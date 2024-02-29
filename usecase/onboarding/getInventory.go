// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package usecase

import (
	"fmt"

	"github.com/ironcore-dev/metal/common/types/errors"
	domain "github.com/ironcore-dev/metal/domain/inventory"
)

type GetInventory interface {
	Execute(inventoryUUID string) (domain.Inventory, error)
}

func InventoryNotFound(err error) error {
	return errors.NewBusinessError(
		notFound,
		fmt.Sprintf("inventory not found: %s", err.Error()),
	)
}
