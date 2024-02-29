// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha4

// SizeSpecApplyConfiguration represents an declarative configuration of the SizeSpec type for use
// with apply.
type SizeSpecApplyConfiguration struct {
	Constraints []ConstraintSpecApplyConfiguration `json:"constraints,omitempty"`
}

// SizeSpecApplyConfiguration constructs an declarative configuration of the SizeSpec type for use with
// apply.
func SizeSpec() *SizeSpecApplyConfiguration {
	return &SizeSpecApplyConfiguration{}
}

// WithConstraints adds the given value to the Constraints field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Constraints field.
func (b *SizeSpecApplyConfiguration) WithConstraints(values ...*ConstraintSpecApplyConfiguration) *SizeSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithConstraints")
		}
		b.Constraints = append(b.Constraints, *values[i])
	}
	return b
}
