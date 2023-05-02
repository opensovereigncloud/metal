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

package dto

import inventories "github.com/onmetal/metal-api/apis/inventory/v1alpha1"

const (
	machineSizeName = "machine"
)

type CreateInventory struct {
	Name      string
	Namespace string
}

func NewCreateInventory(name string, namespace string) CreateInventory {
	return CreateInventory{Name: name, Namespace: namespace}
}

type Inventory struct {
	UUID         string
	Namespace    string
	ProductSKU   string
	SerialNumber string
	Sizes        map[string]string
	NICs         []inventories.NICSpec
}

func NewInventory(
	UUID string,
	namespace string,
	productSKU string,
	serialNumber string,
	sizes map[string]string,
	NICs []inventories.NICSpec) Inventory {
	return Inventory{
		UUID:         UUID,
		Namespace:    namespace,
		ProductSKU:   productSKU,
		SerialNumber: serialNumber,
		Sizes:        sizes,
		NICs:         NICs}
}

func (i *Inventory) IsMachine() bool {
	_, ok := i.Sizes[inventories.GetSizeMatchLabel(machineSizeName)]
	return ok
}
