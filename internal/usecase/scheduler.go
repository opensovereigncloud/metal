package usecase

import (
	"context"

	"github.com/onmetal/metal-api/internal/entity"
	"github.com/onmetal/metal-api/internal/repository"
)

type SchedulerUseCase struct {
	schedulerRepo repository.SchedulerRepo
	reserverRepo  repository.ReserverRepo
}

func NewSchedulerUseCase(s repository.SchedulerRepo, r repository.ReserverRepo) *SchedulerUseCase {
	return &SchedulerUseCase{
		schedulerRepo: s,
		reserverRepo:  r,
	}
}

// Schedule - looks for vacant device (machine or switch) and creates new reservation.
func (s *SchedulerUseCase) Schedule(ctx context.Context, e entity.Order) error {
	reservation, err := s.schedulerRepo.FindVacantDevice(ctx, e)
	if err != nil {
		return err
	}

	if err := s.schedulerRepo.Schedule(ctx, reservation); err != nil {
		return err
	}
	return nil
}

// IsScheduled - checks if assignment is already done.
func (s *SchedulerUseCase) IsScheduled(ctx context.Context, e entity.Order) bool {
	return s.schedulerRepo.IsScheduled(ctx, e)
}

// DeleteScheduling - removes scheduling from the device.
func (s *SchedulerUseCase) DeleteScheduling(ctx context.Context, e entity.Reservation) error {
	if err := s.schedulerRepo.DeleteSchedule(ctx, e); err != nil {
		return err
	}
	return nil
}
