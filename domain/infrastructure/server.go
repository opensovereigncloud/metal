// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"slices"

	"github.com/ironcore-dev/metal/common/types/errors"
)

const PowerCapabilities = "power"

type Server struct {
	UUID              string
	Namespace         string
	PowerCapabilities []string
}

func NewServer(
	UUID string,
	namespace string,
	powerCapabilities []string,
) (Server, errors.BusinessError) {
	if EmptyUUID(UUID) {
		return Server{}, EmptyServerUUID()
	}
	return Server{
		UUID:              UUID,
		Namespace:         namespace,
		PowerCapabilities: powerCapabilities,
	}, nil
}

func (s *Server) HasPowerCapabilities() bool {
	return slices.Contains(s.PowerCapabilities, PowerCapabilities)
}

func EmptyUUID(uuid string) bool {
	return uuid == ""
}

func EmptyServerUUID() errors.BusinessError {
	return errors.NewBusinessError(emptyUUID, "no uuid provided")
}
