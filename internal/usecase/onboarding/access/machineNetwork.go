package access

import (
	machine "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
)

type MachineNetwork interface {
	InterfacesFromInventory(inventory dto.Inventory) []machine.Interface
}
