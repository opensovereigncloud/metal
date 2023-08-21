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

import (
	machine "github.com/onmetal/metal-api/apis/machine/v1alpha3"
	domain "github.com/onmetal/metal-api/domain/inventory"
)

type MachineInfo struct {
	UUID         string
	Namespace    string
	ProductSKU   string
	SerialNumber string
	Sizes        map[string]string
	Interfaces   []machine.Interface
}

func NewMachineInfoFromInventory(inv domain.Inventory) MachineInfo {
	return MachineInfo{
		UUID:         inv.UUID,
		Namespace:    inv.Namespace,
		ProductSKU:   inv.ProductSKU,
		SerialNumber: inv.SerialNumber,
		Sizes:        inv.Sizes,
		Interfaces:   ToMachineInterfaces(inv.NICs),
	}
}
