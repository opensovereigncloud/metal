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

package scenarios

import (
	domain "github.com/onmetal/metal-api/domain/machine"
	usecase "github.com/onmetal/metal-api/usecase/onboarding"
	"github.com/onmetal/metal-api/usecase/onboarding/providers"
)

type GetMachineUseCase struct {
	extractor providers.MachineExtractor
}

func NewGetMachineUseCase(
	machineExtractor providers.MachineExtractor,
) *GetMachineUseCase {
	return &GetMachineUseCase{extractor: machineExtractor}
}

func (g *GetMachineUseCase) Execute(
	machineUUID string,
) (domain.Machine, error) {
	machine, err := g.extractor.ByUUID(machineUUID)
	if err != nil {
		return domain.Machine{}, usecase.MachineNotFound(err)
	}
	return machine, nil
}
