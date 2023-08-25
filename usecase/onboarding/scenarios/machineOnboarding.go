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
	"strings"

	"github.com/go-logr/logr"
	machine "github.com/onmetal/metal-api/apis/machine/v1alpha3"
	invdomain "github.com/onmetal/metal-api/domain/inventory"
	domain "github.com/onmetal/metal-api/domain/machine"
	"github.com/onmetal/metal-api/pkg/network/bgp"
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	"github.com/onmetal/metal-api/usecase/onboarding/providers"
)

type MachineOnboardingUseCase struct {
	machinePersister   providers.MachinePersister
	machineExtractor   providers.MachineExtractor
	switchExtractor    providers.SwitchExtractor
	loopbackRepository providers.LoopbackAddressExtractor
	log                logr.Logger
}

func NewMachineOnboardingUseCase(
	machinePersister providers.MachinePersister,
	machineExtractor providers.MachineExtractor,
	switchExtractor providers.SwitchExtractor,
	loopbackRepository providers.LoopbackAddressExtractor,
	log logr.Logger,
) *MachineOnboardingUseCase {
	log = log.WithName("MachineOnboardingUseCase")
	return &MachineOnboardingUseCase{
		machinePersister:   machinePersister,
		machineExtractor:   machineExtractor,
		switchExtractor:    switchExtractor,
		loopbackRepository: loopbackRepository,
		log:                log,
	}
}

func (m *MachineOnboardingUseCase) Execute(machine domain.Machine, inventory invdomain.Inventory) error {
	machine.Interfaces = m.MachineInterfacesWithSwitchInfo(inventory)
	loopbacks, err := m.LoopbackAddress(machine)
	if err != nil {
		m.log.Info("can't get loopback addresses", "error", err)
	}
	machine.Loopbacks = loopbacks
	machine.SetMachineSizes(inventory.Sizes)

	machine.ASN = bgp.CalculateAutonomousSystemNumberFromAddress(loopbacks.IPv4.Prefix.Addr())
	machine.SKU = inventory.ProductSKU
	machine.SerialNumber = inventory.SerialNumber

	return m.machinePersister.Save(machine)
}

func (m *MachineOnboardingUseCase) MachineInterfacesWithSwitchInfo(
	inventory invdomain.Inventory,
) []machine.Interface {
	interfaces := dto.ToMachineInterfaces(inventory.NICs)
	for i := range interfaces {
		updatedInterfaceInfo, err := m.FindSwitchAndAddInfo(interfaces[i])
		if err != nil {
			continue
		}
		interfaces[i] = updatedInterfaceInfo
	}
	return interfaces
}

func (m *MachineOnboardingUseCase) FindSwitchAndAddInfo(
	machineInterface machine.Interface,
) (machine.Interface, error) {
	chassisID := strings.ReplaceAll(machineInterface.Peer.LLDPChassisID, ":", "-")
	sw, err := m.switchExtractor.ByChassisID(chassisID)
	if err != nil {
		m.log.Info(
			"unable to extract switch info",
			"chassis-id", chassisID,
			"error", err,
		)
		return machine.Interface{}, err
	}
	return sw.AddSwitchInfoToMachineInterfaces(machineInterface), nil
}

func (m *MachineOnboardingUseCase) LoopbackAddress(machine domain.Machine) (domain.Loopbacks, error) {
	var loopbacks domain.Loopbacks
	ipv4LoopbackAddress, err := m.loopbackRepository.Try(3).IPv4ByMachineUUID(machine.UUID)
	if err != nil {
		return loopbacks, err
	}
	loopbacks.IPv4 = ipv4LoopbackAddress
	ipv6LoopbackAddress, err := m.loopbackRepository.Try(3).IPv6ByMachineUUID(machine.UUID)
	if err != nil {
		return loopbacks, err
	}
	loopbacks.IPv6 = ipv6LoopbackAddress
	return loopbacks, nil
}
