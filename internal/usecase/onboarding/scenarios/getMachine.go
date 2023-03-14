package scenarios

import (
	domain "github.com/onmetal/metal-api/internal/domain/machine"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/access"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
)

type GetMachineUseCase struct {
	machineRepository access.MachineRepository
}

func NewGetMachineUseCase(machineRepository access.MachineRepository) *GetMachineUseCase {
	return &GetMachineUseCase{machineRepository: machineRepository}
}

func (g *GetMachineUseCase) Execute(request dto.Request) (domain.Machine, error) {
	return g.machineRepository.Get(request)
}
