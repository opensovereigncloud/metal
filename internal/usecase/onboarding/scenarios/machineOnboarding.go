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
	domain "github.com/onmetal/metal-api/internal/domain/machine"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/access"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
)

type MachineOnboardingUseCase struct {
	machineNetwork    access.MachineNetwork
	machineRepository access.MachineRepository
}

func NewMachineOnboardingUseCase(
	machineNetwork access.MachineNetwork,
	machineRepository access.MachineRepository) *MachineOnboardingUseCase {
	return &MachineOnboardingUseCase{
		machineNetwork:    machineNetwork,
		machineRepository: machineRepository}
}

func (o *MachineOnboardingUseCase) Execute(machine domain.Machine, inventory dto.Inventory) error {
	machine.Interfaces = o.machineNetwork.InterfacesFromInventory(inventory)

	machine.MachineSizes(inventory.Sizes)

	machine.SKU = inventory.ProductSKU
	machine.SerialNumber = inventory.SerialNumber

	return o.machineRepository.Update(machine)
}
