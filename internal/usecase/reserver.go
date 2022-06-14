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
