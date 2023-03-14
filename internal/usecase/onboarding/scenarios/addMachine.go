package scenarios

import (
	"github.com/onmetal/metal-api/internal/usecase/onboarding/access"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
)

type AddMachineUseCase struct {
	machineRepository access.MachineRepository
}

func NewAddMachineUseCase(machineRepository access.MachineRepository) *AddMachineUseCase {
	return &AddMachineUseCase{machineRepository: machineRepository}
}

func (a *AddMachineUseCase) Execute(inventory dto.Inventory) error {
	return a.machineRepository.Create(inventory)
}
