package usecase

import (
	"context"

	"github.com/onmetal/metal-api/internal/entity"
	"github.com/onmetal/metal-api/internal/repository"
)

type SyncUseCase struct {
	repo repository.Synchronization
}

func NewSyncUseCase(r repository.Synchronization) *SyncUseCase {
	return &SyncUseCase{repo: r}
}

// Do - synchronize state of two objects.
func (s *SyncUseCase) Do(ctx context.Context, e entity.Reservation) error {
	sync := entity.Synchronization{
		SourceName:      e.RequestName,
		SourceNamespace: e.RequestNamespace,
		TargetName:      e.OrderName,
		TargetNamespace: e.OrderNamespace,
	}
	return s.repo.Do(ctx, sync)
}
