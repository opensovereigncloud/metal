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

package domain

import (
	machine "github.com/onmetal/metal-api/apis/machine/v1alpha3"
	"github.com/onmetal/metal-api/common/types/base"
	"github.com/onmetal/metal-api/common/types/errors"
	ipdomain "github.com/onmetal/metal-api/domain/address"
)

type Machine struct {
	base.DomainEntity

	ID           MachineID
	UUID         string
	Namespace    string
	ASN          uint32
	SKU          string
	SerialNumber string
	Interfaces   []machine.Interface
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
	interfaces []machine.Interface,
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
	interfaces []machine.Interface,
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
	return errors.NewBusinessError(alreadyExist, "machine already exist")
}
