// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package domain

type InventoryAlreadyExist interface {
	Invoke(inventoryUUID string) bool
}
