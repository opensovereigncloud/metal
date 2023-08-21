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

package scenarios

import (
	domain "github.com/onmetal/metal-api/domain/inventory"
	usecase "github.com/onmetal/metal-api/usecase/onboarding"
	"github.com/onmetal/metal-api/usecase/onboarding/providers"
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
