package usecase

import (
	"context"

	"github.com/onmetal/metal-api/internal/entity"
	"github.com/onmetal/metal-api/internal/repository"
)

type ServerOnboardingUseCase struct {
	repo repository.OnboardingRepo
}

func NewServerOnboarding(p repository.OnboardingRepo) *ServerOnboardingUseCase {
	return &ServerOnboardingUseCase{p}
}

// Initiate - creates new raw object from real server.
func (o *ServerOnboardingUseCase) Initiate(ctx context.Context, e entity.Onboarding) error {
	if err := o.repo.Prepare(ctx, e); err != nil {
		return err
	}

	return o.repo.Create(ctx)
}

// IsInitialized - checks if raw object is already initialized.
func (o *ServerOnboardingUseCase) InitializationStatus(ctx context.Context,
	e entity.Onboarding) entity.Initialization {
	return o.repo.InitializationStatus(ctx, e)
}

// GatherData - retrieves data from real server.
func (o *ServerOnboardingUseCase) GatherData(ctx context.Context, e entity.Onboarding) error {
	return o.repo.GatherData(ctx, e)
}
