// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
