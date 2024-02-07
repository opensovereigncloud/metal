// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package domain

import "github.com/ironcore-dev/metal/common/types/base"

type MachineCreatedDomainEvent struct {
	id MachineID
}

func (m *MachineCreatedDomainEvent) ID() string {
	return m.id.String()
}

func (m *MachineCreatedDomainEvent) Type() string {
	return "machine created"
}

func NewMachineCreatedDomainEvent(
	machineID MachineID,
) base.DomainEvent {
	return &MachineCreatedDomainEvent{
		id: machineID,
	}
}
