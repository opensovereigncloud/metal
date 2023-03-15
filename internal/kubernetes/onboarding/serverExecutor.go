package persistence

import (
	"github.com/go-logr/logr"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
)

type FakeServerExecutor struct {
	log logr.Logger
}

func NewFakeServerExecutor(log logr.Logger) *FakeServerExecutor {
	return &FakeServerExecutor{log: log}
}

func (f *FakeServerExecutor) Enable(request dto.Request) error {
	f.log.Info("server turned on after inventory onboarding", "server", request)
	return nil
}
