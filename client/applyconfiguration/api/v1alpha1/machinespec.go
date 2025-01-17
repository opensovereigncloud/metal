// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/ironcore-dev/metal/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

// MachineSpecApplyConfiguration represents an declarative configuration of the MachineSpec type for use
// with apply.
type MachineSpecApplyConfiguration struct {
	UUID               *string                  `json:"uuid,omitempty"`
	OOBRef             *v1.LocalObjectReference `json:"oobRef,omitempty"`
	InventoryRef       *v1.LocalObjectReference `json:"inventoryRef,omitempty"`
	MachineClaimRef    *v1.ObjectReference      `json:"machineClaimRef,omitempty"`
	LoopbackAddressRef *v1.LocalObjectReference `json:"loopbackAddressRef,omitempty"`
	ASN                *string                  `json:"asn,omitempty"`
	Power              *v1alpha1.Power          `json:"power,omitempty"`
	LocatorLED         *v1alpha1.LocatorLED     `json:"locatorLED,omitempty"`
}

// MachineSpecApplyConfiguration constructs an declarative configuration of the MachineSpec type for use with
// apply.
func MachineSpec() *MachineSpecApplyConfiguration {
	return &MachineSpecApplyConfiguration{}
}

// WithUUID sets the UUID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UUID field is set to the value of the last call.
func (b *MachineSpecApplyConfiguration) WithUUID(value string) *MachineSpecApplyConfiguration {
	b.UUID = &value
	return b
}

// WithOOBRef sets the OOBRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the OOBRef field is set to the value of the last call.
func (b *MachineSpecApplyConfiguration) WithOOBRef(value v1.LocalObjectReference) *MachineSpecApplyConfiguration {
	b.OOBRef = &value
	return b
}

// WithInventoryRef sets the InventoryRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the InventoryRef field is set to the value of the last call.
func (b *MachineSpecApplyConfiguration) WithInventoryRef(value v1.LocalObjectReference) *MachineSpecApplyConfiguration {
	b.InventoryRef = &value
	return b
}

// WithMachineClaimRef sets the MachineClaimRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MachineClaimRef field is set to the value of the last call.
func (b *MachineSpecApplyConfiguration) WithMachineClaimRef(value v1.ObjectReference) *MachineSpecApplyConfiguration {
	b.MachineClaimRef = &value
	return b
}

// WithLoopbackAddressRef sets the LoopbackAddressRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LoopbackAddressRef field is set to the value of the last call.
func (b *MachineSpecApplyConfiguration) WithLoopbackAddressRef(value v1.LocalObjectReference) *MachineSpecApplyConfiguration {
	b.LoopbackAddressRef = &value
	return b
}

// WithASN sets the ASN field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ASN field is set to the value of the last call.
func (b *MachineSpecApplyConfiguration) WithASN(value string) *MachineSpecApplyConfiguration {
	b.ASN = &value
	return b
}

// WithPower sets the Power field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Power field is set to the value of the last call.
func (b *MachineSpecApplyConfiguration) WithPower(value v1alpha1.Power) *MachineSpecApplyConfiguration {
	b.Power = &value
	return b
}

// WithLocatorLED sets the LocatorLED field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LocatorLED field is set to the value of the last call.
func (b *MachineSpecApplyConfiguration) WithLocatorLED(value v1alpha1.LocatorLED) *MachineSpecApplyConfiguration {
	b.LocatorLED = &value
	return b
}
