package usecase

import (
	"context"

	"github.com/onmetal/metal-api/internal/entity"
)

type (
	Onboarding interface {
		Initiate(context.Context, entity.Onboarding) error
		IsInitialized(context.Context, entity.Onboarding) bool
		GatherData(context.Context, entity.Onboarding) error
	}
	Reserver interface {
		GetReservation(context.Context, entity.Order) (entity.Reservation, error)
		CheckIn(context.Context, entity.Reservation) error
		CheckOut(context.Context, entity.Reservation) error
	}
	Scheduler interface {
		Schedule(context.Context, entity.Order) error
		IsScheduled(context.Context, entity.Order) bool
		DeleteScheduling(context.Context, entity.Reservation) error
	}
	Synchronization interface {
		Do(context.Context, entity.Reservation) error
	}
)
