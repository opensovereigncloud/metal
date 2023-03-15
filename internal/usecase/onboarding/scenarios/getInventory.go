package scenarios

import (
	"github.com/onmetal/metal-api/internal/usecase/onboarding/access"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
)

type GetInventoryUseCase struct {
	inventoryRepository access.InventoryRepository
}

func NewGetInventoryUseCase(inventoryRepository access.InventoryRepository) *GetInventoryUseCase {
	return &GetInventoryUseCase{inventoryRepository: inventoryRepository}
}

func (g *GetInventoryUseCase) Execute(request dto.Request) (dto.Inventory, error) {
	return g.inventoryRepository.Get(request)
}
