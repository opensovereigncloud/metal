// Copyright 2023 OnMetal authors
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

package domain

import (
	"slices"

	"github.com/onmetal/metal-api/common/types/errors"
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
