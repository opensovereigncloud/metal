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

type DeviceOnboardingUseCase struct {
	repo repository.OnboardingRepo
}

func NewDeviceOnboarding(p repository.OnboardingRepo) *DeviceOnboardingUseCase {
	return &DeviceOnboardingUseCase{p}
}

// Initiate - creates new devices from onboarded server.
func (o *DeviceOnboardingUseCase) Initiate(ctx context.Context, e entity.Onboarding) error {
	if err := o.repo.Prepare(ctx, e); err != nil {
		return err
	}
	return o.repo.Create(ctx)
}

// IsInitialized - checks if device is already initialized.
func (o *DeviceOnboardingUseCase) InitializationStatus(ctx context.Context,
	e entity.Onboarding) entity.Initialization {
	return o.repo.InitializationStatus(ctx, e)
}

// GatherData - retrieves data from raw object into abstract one.
func (o *DeviceOnboardingUseCase) GatherData(ctx context.Context, e entity.Onboarding) error {
	return o.repo.GatherData(ctx, e)
}
