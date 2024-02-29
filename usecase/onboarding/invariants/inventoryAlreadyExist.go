// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package invariants

import "github.com/ironcore-dev/metal/usecase/onboarding/providers"

type InventoryAlreadyExist struct {
	extractor providers.InventoryExtractor
}

func NewInventoryAlreadyExist(
	inventoryExtractor providers.InventoryExtractor,
) *InventoryAlreadyExist {
	return &InventoryAlreadyExist{
		extractor: inventoryExtractor,
	}
}

func (m *InventoryAlreadyExist) Invoke(inventoryUUID string) bool {
	inv, _ := m.
		extractor.
		ByUUID(inventoryUUID)
	return inv.ID.String() != ""
}
