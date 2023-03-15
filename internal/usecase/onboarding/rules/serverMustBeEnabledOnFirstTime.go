package rules

import (
	"github.com/go-logr/logr"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/access"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
)

type ServerMustBeEnabledOnFirstTimeRule struct {
	serverExecutor access.ServerExecutor
	log            logr.Logger
}

func NewServerMustBeEnabledOnFirstTimeRule(
	serverExecutor access.ServerExecutor,
	log logr.Logger) *ServerMustBeEnabledOnFirstTimeRule {
	return &ServerMustBeEnabledOnFirstTimeRule{
		serverExecutor: serverExecutor, log: log}
}

func (s *ServerMustBeEnabledOnFirstTimeRule) Execute(request dto.Request) error {
	return s.serverExecutor.Enable(request)
}
