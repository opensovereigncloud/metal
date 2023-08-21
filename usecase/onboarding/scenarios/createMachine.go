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
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	"github.com/onmetal/metal-api/usecase/onboarding/providers"
)

type CreateMachineUseCase struct {
	machineRepository   providers.MachinePersister
	machineIDGenerator  domain.MachineIDGenerator
	machineAlreadyExist domain.MachineAlreadyExist
}

func NewCreateMachineUseCase(
	machinePersister providers.MachinePersister,
	machineIDGenerator domain.MachineIDGenerator,
	machineAlreadyExist domain.MachineAlreadyExist,
) *CreateMachineUseCase {
	return &CreateMachineUseCase{
		machineRepository:   machinePersister,
		machineIDGenerator:  machineIDGenerator,
		machineAlreadyExist: machineAlreadyExist,
	}
}

func (a *CreateMachineUseCase) Execute(
	machineInfo dto.MachineInfo,
) (domain.MachineID, error) {
	machine, err := domain.CreateMachine(
		a.machineIDGenerator,
		a.machineAlreadyExist,
		machineInfo.UUID,
		machineInfo.Namespace,
		machineInfo.ProductSKU,
		machineInfo.SerialNumber,
		machineInfo.Interfaces,
		domain.Loopbacks{},
		machineInfo.Sizes)
	if err != nil {
		return domain.MachineID{}, usecase.MachineAlreadCreated(machineInfo.UUID)
	}
	return machine.ID, a.machineRepository.Save(machine)
}
