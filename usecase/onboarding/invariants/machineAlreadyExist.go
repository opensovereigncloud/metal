// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package invariants

import "github.com/ironcore-dev/metal/usecase/onboarding/providers"

type MachineAlreadyExist struct {
	machineExtractor providers.MachineExtractor
}

func NewMachineAlreadyExist(
	machineExtractor providers.MachineExtractor,
) *MachineAlreadyExist {
	return &MachineAlreadyExist{
		machineExtractor: machineExtractor,
	}
}

func (m *MachineAlreadyExist) Invoke(machineUUID string) bool {
	machine, _ := m.
		machineExtractor.
		ByUUID(machineUUID)
	return machine.ID.String() != ""
}
