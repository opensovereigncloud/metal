// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/ironcore-dev/metal/common/types/base"
	"github.com/ironcore-dev/metal/common/types/errors"
	ipdomain "github.com/ironcore-dev/metal/domain/address"
)

type Machine struct {
	base.DomainEntity

	ID           MachineID
	UUID         string
	Namespace    string
	ASN          uint32
	SKU          string
	SerialNumber string
	Interfaces   []metalv1alpha4.Interface
	Loopbacks    Loopbacks
	Size         map[string]string
}

func NewMachine(
	ID MachineID,
	UUID string,
	namespace string,
	ASN uint32,
	SKU string,
	serialNumber string,
	interfaces []metalv1alpha4.Interface,
	loopbacks Loopbacks,
	size map[string]string,
) Machine {
	domainEntity := base.NewDomainEntity()
	return Machine{
		DomainEntity: domainEntity,
		ID:           ID,
		UUID:         UUID,
		Namespace:    namespace,
		ASN:          ASN,
		SKU:          SKU,
		SerialNumber: serialNumber,
		Interfaces:   interfaces,
		Loopbacks:    loopbacks,
		Size:         size,
	}
}

type Loopbacks struct {
	IPv4 ipdomain.Address
	IPv6 ipdomain.Address
}

func CreateMachine(
	idGenerator MachineIDGenerator,
	machineAlreadyExist MachineAlreadyExist,
	UUID string,
	namespace string,
	SKU string,
	serialNumber string,
	interfaces []metalv1alpha4.Interface,
	loopbacks Loopbacks,
	size map[string]string,
) (Machine, errors.BusinessError) {
	if machineAlreadyExist.Invoke(UUID) {
		return Machine{}, MachineAlreadyCreated()
	}
	machineID := idGenerator.Generate()
	domainEntity := base.NewDomainEntity()
	domainEntity.AddEvent(NewMachineCreatedDomainEvent(machineID))
	return Machine{
		DomainEntity: domainEntity,
		ID:           machineID,
		UUID:         UUID,
		Namespace:    namespace,
		SKU:          SKU,
		SerialNumber: serialNumber,
		Interfaces:   interfaces,
		Loopbacks:    loopbacks,
		Size:         size,
	}, nil
}

func (m *Machine) SetMachineSizes(sizes map[string]string) { m.Size = sizes }

func MachineAlreadyCreated() errors.BusinessError {
	return errors.NewBusinessError(alreadyExist, "metalv1alpha4 already exist")
}
