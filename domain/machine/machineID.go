// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package domain

type MachineIDGenerator interface {
	Generate() MachineID
}

type MachineID struct {
	value string
}

func NewMachineID(id string) MachineID {
	return MachineID{
		value: id,
	}
}

func (m *MachineID) String() string {
	return m.value
}
