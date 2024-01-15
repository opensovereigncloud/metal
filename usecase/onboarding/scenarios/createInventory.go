// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
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
// */

package scenarios

import (
	domain "github.com/ironcore-dev/metal/domain/inventory"
	usecase "github.com/ironcore-dev/metal/usecase/onboarding"
	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
	"github.com/ironcore-dev/metal/usecase/onboarding/providers"
)

type CreateInventoryUseCase struct {
	inventoryAlreadyExist domain.InventoryAlreadyExist
	inventoryIDGenerator  domain.InventoryIDGenerator
	inventoryRepository   providers.InventoryPersister
}

func NewCreateInventoryUseCase(
	inventoryAlreadyExist domain.InventoryAlreadyExist,
	inventoryIDGenerator domain.InventoryIDGenerator,
	inventoryRepository providers.InventoryPersister,
) *CreateInventoryUseCase {
	return &CreateInventoryUseCase{
		inventoryIDGenerator:  inventoryIDGenerator,
		inventoryAlreadyExist: inventoryAlreadyExist,
		inventoryRepository:   inventoryRepository,
	}
}

func (o *CreateInventoryUseCase) Execute(
	inventoryInfo dto.InventoryInfo,
) error {
	inv, err := domain.CreateInventory(
		o.inventoryIDGenerator,
		o.inventoryAlreadyExist,
		inventoryInfo.UUID,
		inventoryInfo.Namespace)
	if err != nil {
		return usecase.InventoryAlreadyCreated(inventoryInfo.UUID)
	}
	if err := o.inventoryRepository.Save(inv); err != nil {
		return err
	}
	return nil
}
