// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
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
// */

package scenarios

import (
	"github.com/onmetal/metal-api/usecase/onboarding/access"
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
)

type ServerValidationUseCase struct {
	serverRepository access.ServerRepository
}

func NewServerValidationUseCase(serverRepository access.ServerRepository) *ServerValidationUseCase {
	return &ServerValidationUseCase{serverRepository: serverRepository}
}

func (o *ServerValidationUseCase) Execute(request dto.Request) bool {
	server, err := o.serverRepository.Get(request)
	if err != nil {
		return false
	}
	if !server.HasPowerCapabilities() ||
		!server.HasUUID() {
		return false
	}
	return true
}
