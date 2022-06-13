package usecase

import (
	"context"

	"github.com/onmetal/metal-api/internal/entity"
	"github.com/onmetal/metal-api/internal/repository"
)

type ReserverUseCase struct {
	repo repository.ReserverRepo
}

func NewReserverUseCase(r repository.ReserverRepo) *ReserverUseCase {
	return &ReserverUseCase{
		repo: r,
	}
}

// GetReservation - retrieve reservation status of Order.
func (r *ReserverUseCase) GetReservation(ctx context.Context, e entity.Order) (entity.Reservation, error) {
	return r.repo.GetReservation(ctx, e)
}

// CheckIn - enables scheduled device. Device will be powered on.
func (r *ReserverUseCase) CheckIn(ctx context.Context, e entity.Reservation) error {
	return r.repo.CheckIn(ctx, e)

}

// CheckOut - disables scheduled device. Device will be powered off.
func (r *ReserverUseCase) CheckOut(ctx context.Context, e entity.Reservation) error {
	return r.repo.CheckOut(ctx, e)
}
