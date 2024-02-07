// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scenarios

import (
	domain "github.com/ironcore-dev/metal/domain/inventory"
	usecase "github.com/ironcore-dev/metal/usecase/onboarding"
	"github.com/ironcore-dev/metal/usecase/onboarding/providers"
)

type GetInventoryUseCase struct {
	extractor providers.InventoryExtractor
}

func NewGetInventoryUseCase(
	inventoryExtractor providers.InventoryExtractor,
) *GetInventoryUseCase {
	return &GetInventoryUseCase{extractor: inventoryExtractor}
}

func (g *GetInventoryUseCase) Execute(
	inventoryUUID string,
) (domain.Inventory, error) {
	inventory, err := g.extractor.ByUUID(inventoryUUID)
	if err != nil {
		return inventory, usecase.InventoryNotFound(err)
	}
	return inventory, nil
}
