// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scenarios

import (
	domain "github.com/ironcore-dev/metal/domain/machine"
	usecase "github.com/ironcore-dev/metal/usecase/onboarding"
	"github.com/ironcore-dev/metal/usecase/onboarding/providers"
)

type GetMachineUseCase struct {
	extractor providers.MachineExtractor
}

func NewGetMachineUseCase(
	machineExtractor providers.MachineExtractor,
) *GetMachineUseCase {
	return &GetMachineUseCase{extractor: machineExtractor}
}

func (g *GetMachineUseCase) Execute(
	machineUUID string,
) (domain.Machine, error) {
	machine, err := g.extractor.ByUUID(machineUUID)
	if err != nil {
		return domain.Machine{}, usecase.MachineNotFound(err)
	}
	return machine, nil
}
