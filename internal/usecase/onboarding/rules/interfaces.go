package rules

import "github.com/onmetal/metal-api/internal/usecase/onboarding/dto"

type ServerMustBeEnabledOnFirstTime interface {
	Execute(request dto.Request) error
}
