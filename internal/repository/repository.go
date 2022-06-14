package repository

import (
	"context"

	"github.com/onmetal/metal-api/internal/entity"
)

type (
	OnboardingRepo interface {
		Create(context.Context) error
		InitializationStatus(context.Context, entity.Onboarding) entity.Initialization
		Prepare(context.Context, entity.Onboarding) error
		GatherData(context.Context, entity.Onboarding) error
	}
	ReserverRepo interface {
		CheckIn(context.Context, entity.Reservation) error
		CheckOut(context.Context, entity.Reservation) error
		GetReservation(context.Context, entity.Order) (entity.Reservation, error)
	}
	SchedulerRepo interface {
		Schedule(context.Context, entity.Reservation) error
		DeleteSchedule(context.Context, entity.Reservation) error
		IsScheduled(context.Context, entity.Order) bool
		FindVacantDevice(context.Context, entity.Order) (entity.Reservation, error)
	}
	Synchronization interface {
		Do(context.Context, entity.Synchronization) error
	}
)
