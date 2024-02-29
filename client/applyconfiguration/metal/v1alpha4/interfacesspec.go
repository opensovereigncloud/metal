// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha4

import (
	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
)

// InterfacesSpecApplyConfiguration represents an declarative configuration of the InterfacesSpec type for use
// with apply.
type InterfacesSpecApplyConfiguration struct {
	Defaults  *PortParametersSpecApplyConfiguration   `json:"defaults,omitempty"`
	Overrides []*metalv1alpha4.InterfaceOverridesSpec `json:"overrides,omitempty"`
}

// InterfacesSpecApplyConfiguration constructs an declarative configuration of the InterfacesSpec type for use with
// apply.
func InterfacesSpec() *InterfacesSpecApplyConfiguration {
	return &InterfacesSpecApplyConfiguration{}
}

// WithDefaults sets the Defaults field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Defaults field is set to the value of the last call.
func (b *InterfacesSpecApplyConfiguration) WithDefaults(value *PortParametersSpecApplyConfiguration) *InterfacesSpecApplyConfiguration {
	b.Defaults = value
	return b
}

// WithOverrides adds the given value to the Overrides field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Overrides field.
func (b *InterfacesSpecApplyConfiguration) WithOverrides(values ...**metalv1alpha4.InterfaceOverridesSpec) *InterfacesSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithOverrides")
		}
		b.Overrides = append(b.Overrides, *values[i])
	}
	return b
}
