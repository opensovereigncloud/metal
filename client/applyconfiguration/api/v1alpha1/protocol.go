// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/ironcore-dev/metal/api/v1alpha1"
)

// ProtocolApplyConfiguration represents an declarative configuration of the Protocol type for use
// with apply.
type ProtocolApplyConfiguration struct {
	Name *v1alpha1.ProtocolName `json:"name,omitempty"`
	Port *int32                 `json:"port,omitempty"`
}

// ProtocolApplyConfiguration constructs an declarative configuration of the Protocol type for use with
// apply.
func Protocol() *ProtocolApplyConfiguration {
	return &ProtocolApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *ProtocolApplyConfiguration) WithName(value v1alpha1.ProtocolName) *ProtocolApplyConfiguration {
	b.Name = &value
	return b
}

// WithPort sets the Port field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Port field is set to the value of the last call.
func (b *ProtocolApplyConfiguration) WithPort(value int32) *ProtocolApplyConfiguration {
	b.Port = &value
	return b
}
