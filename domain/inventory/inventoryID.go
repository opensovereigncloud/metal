// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package domain

type InventoryIDGenerator interface {
	Generate() InventoryID
}
type InventoryID struct {
	value string
}

func NewInventoryID(id string) InventoryID {
	return InventoryID{
		value: id,
	}
}

func (m *InventoryID) String() string {
	return m.value
}
