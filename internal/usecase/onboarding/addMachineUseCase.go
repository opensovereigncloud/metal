package usecase

import "github.com/onmetal/metal-api/internal/usecase/onboarding/dto"

type AddMachineUseCase interface {
	Execute(request dto.Inventory) error
}
