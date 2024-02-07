// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"github.com/ironcore-dev/metal/common/types/base"
	"github.com/ironcore-dev/metal/common/types/errors"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
)

const (
	machineClassName = "machine"
)

type Inventory struct {
	base.DomainEntity

	ID           InventoryID
	UUID         string
	Namespace    string
	ProductSKU   string
	SerialNumber string
	Sizes        map[string]string
	NICs         []metalv1alpha4.NICSpec
}

func NewInventory(
	ID InventoryID,
	UUID string,
	namespace string,
	productSKU string,
	serialNumber string,
	sizes map[string]string,
	NICs []metalv1alpha4.NICSpec,
) Inventory {
	domainEntity := base.NewDomainEntity()
	return Inventory{
		DomainEntity: domainEntity,
		ID:           ID,
		UUID:         UUID,
		Namespace:    namespace,
		ProductSKU:   productSKU,
		SerialNumber: serialNumber,
		Sizes:        sizes,
		NICs:         NICs,
	}
}

func CreateInventory(
	inventoryIDGenerator InventoryIDGenerator,
	inventoryAlreadyExist InventoryAlreadyExist,
	UUID string,
	namespace string,
) (Inventory, errors.BusinessError) {
	if inventoryAlreadyExist.Invoke(UUID) {
		return Inventory{}, InventoryAlreadyCreated()
	}
	domainEntity := base.NewDomainEntity()
	inventoryID := inventoryIDGenerator.Generate()
	domainEntity.AddEvent(NewInventoryCreatedDomainEvent(inventoryID))
	return Inventory{
		DomainEntity: domainEntity,
		ID:           inventoryID,
		UUID:         UUID,
		Namespace:    namespace,
	}, nil
}

func (i *Inventory) IsMachine() bool {
	_, ok := i.Sizes[metalv1alpha4.GetSizeMatchLabel(machineClassName)]
	return ok
}

func InventoryAlreadyCreated() errors.BusinessError {
	return errors.NewBusinessError(
		alreadyExist,
		"inventory already exist",
	)
}
