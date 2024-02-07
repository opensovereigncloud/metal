// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	"github.com/google/uuid"
	domain "github.com/ironcore-dev/metal/domain/machine"
)

type KubernetesMachineIDGenerator struct{}

func NewKubernetesMachineIDGenerator() *KubernetesMachineIDGenerator {
	return &KubernetesMachineIDGenerator{}
}

func (m *KubernetesMachineIDGenerator) Generate() domain.MachineID {
	return domain.NewMachineID(uuid.NewString())
}
