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
