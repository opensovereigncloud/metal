// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package usecase

import (
	"fmt"

	"github.com/ironcore-dev/metal/common/types/errors"
	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
)

type CreateInventory interface {
	Execute(inventoryInfo dto.InventoryInfo) error
}

func InventoryAlreadyCreated(name string) error {
	return errors.NewBusinessError(
		alreadyCreated,
		fmt.Sprintf("Inventory Already Created. Name: %s", name),
	)
}
