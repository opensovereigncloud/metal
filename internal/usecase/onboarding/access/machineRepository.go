package access

import (
	domain "github.com/onmetal/metal-api/internal/domain/machine"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
)

type MachineRepository interface {
	Create(inventory dto.Inventory) error
	Update(machine domain.Machine) error
	Get(request dto.Request) (domain.Machine, error)
}
