// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scenarios

import (
	domain "github.com/ironcore-dev/metal/domain/infrastructure"
	usecase "github.com/ironcore-dev/metal/usecase/onboarding"
	"github.com/ironcore-dev/metal/usecase/onboarding/providers"
)

type GetServerUseCase struct {
	extractor providers.ServerExtractor
}

func NewGetServerUseCase(
	serverExtractor providers.ServerExtractor,
) *GetServerUseCase {
	return &GetServerUseCase{extractor: serverExtractor}
}

func (g *GetServerUseCase) Execute(
	uuid string,
) (domain.Server, error) {
	server, err := g.extractor.ByUUID(uuid)
	if err != nil {
		return domain.Server{}, usecase.ServerNotFound(err)
	}
	return server, nil
}
