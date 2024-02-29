// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha4

import (
	v1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
)

// MachineStatusApplyConfiguration represents an declarative configuration of the MachineStatus type for use
// with apply.
type MachineStatusApplyConfiguration struct {
	Reboot      *string                        `json:"reboot,omitempty"`
	Health      *v1alpha4.MachineState         `json:"health,omitempty"`
	Network     *NetworkApplyConfiguration     `json:"network,omitempty"`
	Reservation *ReservationApplyConfiguration `json:"reservation,omitempty"`
	Orphaned    *bool                          `json:"orphaned,omitempty"`
}

// MachineStatusApplyConfiguration constructs an declarative configuration of the MachineStatus type for use with
// apply.
func MachineStatus() *MachineStatusApplyConfiguration {
	return &MachineStatusApplyConfiguration{}
}

// WithReboot sets the Reboot field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Reboot field is set to the value of the last call.
func (b *MachineStatusApplyConfiguration) WithReboot(value string) *MachineStatusApplyConfiguration {
	b.Reboot = &value
	return b
}

// WithHealth sets the Health field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Health field is set to the value of the last call.
func (b *MachineStatusApplyConfiguration) WithHealth(value v1alpha4.MachineState) *MachineStatusApplyConfiguration {
	b.Health = &value
	return b
}

// WithNetwork sets the Network field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Network field is set to the value of the last call.
func (b *MachineStatusApplyConfiguration) WithNetwork(value *NetworkApplyConfiguration) *MachineStatusApplyConfiguration {
	b.Network = value
	return b
}

// WithReservation sets the Reservation field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Reservation field is set to the value of the last call.
func (b *MachineStatusApplyConfiguration) WithReservation(value *ReservationApplyConfiguration) *MachineStatusApplyConfiguration {
	b.Reservation = value
	return b
}

// WithOrphaned sets the Orphaned field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Orphaned field is set to the value of the last call.
func (b *MachineStatusApplyConfiguration) WithOrphaned(value bool) *MachineStatusApplyConfiguration {
	b.Orphaned = &value
	return b
}
