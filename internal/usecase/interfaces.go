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
)

type (
	Onboarding interface {
		Initiate(context.Context, entity.Onboarding) error
		InitializationStatus(context.Context, entity.Onboarding) entity.Initialization
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
