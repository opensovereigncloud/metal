// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package dto

import (
	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	domain "github.com/ironcore-dev/metal/domain/inventory"
)

type MachineInfo struct {
	UUID         string
	Namespace    string
	ProductSKU   string
	SerialNumber string
	Sizes        map[string]string
	Interfaces   []metalv1alpha4.Interface
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
