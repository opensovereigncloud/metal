package access

import "github.com/onmetal/metal-api/internal/usecase/onboarding/dto"

type InventoryRepository interface {
	Create(inventory dto.CreateInventory) error
	Get(request dto.Request) (dto.Inventory, error)
}
