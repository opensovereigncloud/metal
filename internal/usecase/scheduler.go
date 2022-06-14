// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

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
