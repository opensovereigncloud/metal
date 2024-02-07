// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scenarios

import (
	"strings"

	"github.com/go-logr/logr"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	invdomain "github.com/ironcore-dev/metal/domain/inventory"
	domain "github.com/ironcore-dev/metal/domain/machine"
	"github.com/ironcore-dev/metal/pkg/network/bgp"
	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
	"github.com/ironcore-dev/metal/usecase/onboarding/providers"
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
) []metalv1alpha4.Interface {
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
	machineInterface metalv1alpha4.Interface,
) (metalv1alpha4.Interface, error) {
	chassisID := strings.ReplaceAll(machineInterface.Peer.LLDPChassisID, ":", "-")
	sw, err := m.switchExtractor.ByChassisID(chassisID)
	if err != nil {
		m.log.Info(
			"unable to extract switch info",
			"chassis-id", chassisID,
			"error", err,
		)
		return metalv1alpha4.Interface{}, err
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
