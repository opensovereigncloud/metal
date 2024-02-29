// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	"github.com/google/uuid"
	domain "github.com/ironcore-dev/metal/domain/inventory"
)

type KubernetesInventoryIDGenerator struct{}

func NewKubernetesInventoryIDGenerator() *KubernetesInventoryIDGenerator {
	return &KubernetesInventoryIDGenerator{}
}

func (m *KubernetesInventoryIDGenerator) Generate() domain.InventoryID {
	return domain.NewInventoryID(uuid.NewString())
}
