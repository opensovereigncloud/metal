// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	domain "github.com/ironcore-dev/metal/domain/infrastructure"
	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
)

type ServerExtractor interface {
	ByUUID(uuid string) (domain.Server, error)
	Get(request dto.Request) (domain.Server, error)
}
