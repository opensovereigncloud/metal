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
func (o *DeviceOnboardingUseCase) IsInitialized(ctx context.Context, e entity.Onboarding) bool {
	return o.repo.IsInitialized(ctx, e)
}

// GatherData - retrieves data from raw object into abstract one.
func (o *DeviceOnboardingUseCase) GatherData(ctx context.Context, e entity.Onboarding) error {
	return o.repo.GatherData(ctx, e)
}
