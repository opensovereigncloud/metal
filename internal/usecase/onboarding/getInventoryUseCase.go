package usecase

import (
	"fmt"

	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
)

type GetInventoryUseCase interface {
	Execute(request dto.Request) (dto.Inventory, error)
}

func MachineNotFound(name string) error {
	return &OnboardingError{
		Reason:  notFound,
		Message: fmt.Sprintf("machine not found: %s", name),
	}
}
