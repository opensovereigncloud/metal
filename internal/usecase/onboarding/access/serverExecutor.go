package access

import "github.com/onmetal/metal-api/internal/usecase/onboarding/dto"

type ServerExecutor interface {
	Enable(request dto.Request) error
}
