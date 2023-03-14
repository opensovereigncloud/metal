package usecase

import (
	domain "github.com/onmetal/metal-api/internal/domain/machine"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
)

type MachineOnboardingUseCase interface {
	Execute(machine domain.Machine, inventory dto.Inventory) error
}
