package usecase

import (
	domain "github.com/onmetal/metal-api/internal/domain/machine"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
)

type GetMachineUseCase interface {
	Execute(request dto.Request) (domain.Machine, error)
}
