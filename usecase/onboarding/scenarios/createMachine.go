// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scenarios

import (
	domain "github.com/ironcore-dev/metal/domain/machine"
	usecase "github.com/ironcore-dev/metal/usecase/onboarding"
	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
	"github.com/ironcore-dev/metal/usecase/onboarding/providers"
)

type CreateMachineUseCase struct {
	machineRepository   providers.MachinePersister
	machineIDGenerator  domain.MachineIDGenerator
	machineAlreadyExist domain.MachineAlreadyExist
}

func NewCreateMachineUseCase(
	machinePersister providers.MachinePersister,
	machineIDGenerator domain.MachineIDGenerator,
	machineAlreadyExist domain.MachineAlreadyExist,
) *CreateMachineUseCase {
	return &CreateMachineUseCase{
		machineRepository:   machinePersister,
		machineIDGenerator:  machineIDGenerator,
		machineAlreadyExist: machineAlreadyExist,
	}
}

func (a *CreateMachineUseCase) Execute(
	machineInfo dto.MachineInfo,
) (domain.MachineID, error) {
	machine, err := domain.CreateMachine(
		a.machineIDGenerator,
		a.machineAlreadyExist,
		machineInfo.UUID,
		machineInfo.Namespace,
		machineInfo.ProductSKU,
		machineInfo.SerialNumber,
		machineInfo.Interfaces,
		domain.Loopbacks{},
		machineInfo.Sizes)
	if err != nil {
		return domain.MachineID{}, usecase.MachineAlreadCreated(machineInfo.UUID)
	}
	return machine.ID, a.machineRepository.Save(machine)
}
