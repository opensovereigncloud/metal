// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha4

// PeerInfoSpecApplyConfiguration represents an declarative configuration of the PeerInfoSpec type for use
// with apply.
type PeerInfoSpecApplyConfiguration struct {
	ChassisID       *string `json:"chassisId,omitempty"`
	SystemName      *string `json:"systemName,omitempty"`
	PortID          *string `json:"portId,omitempty"`
	PortDescription *string `json:"portDescription,omitempty"`
	Type            *string `json:"type,omitempty"`
}

// PeerInfoSpecApplyConfiguration constructs an declarative configuration of the PeerInfoSpec type for use with
// apply.
func PeerInfoSpec() *PeerInfoSpecApplyConfiguration {
	return &PeerInfoSpecApplyConfiguration{}
}

// WithChassisID sets the ChassisID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ChassisID field is set to the value of the last call.
func (b *PeerInfoSpecApplyConfiguration) WithChassisID(value string) *PeerInfoSpecApplyConfiguration {
	b.ChassisID = &value
	return b
}

// WithSystemName sets the SystemName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SystemName field is set to the value of the last call.
func (b *PeerInfoSpecApplyConfiguration) WithSystemName(value string) *PeerInfoSpecApplyConfiguration {
	b.SystemName = &value
	return b
}

// WithPortID sets the PortID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PortID field is set to the value of the last call.
func (b *PeerInfoSpecApplyConfiguration) WithPortID(value string) *PeerInfoSpecApplyConfiguration {
	b.PortID = &value
	return b
}

// WithPortDescription sets the PortDescription field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PortDescription field is set to the value of the last call.
func (b *PeerInfoSpecApplyConfiguration) WithPortDescription(value string) *PeerInfoSpecApplyConfiguration {
	b.PortDescription = &value
	return b
}

// WithType sets the Type field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Type field is set to the value of the last call.
func (b *PeerInfoSpecApplyConfiguration) WithType(value string) *PeerInfoSpecApplyConfiguration {
	b.Type = &value
	return b
}
