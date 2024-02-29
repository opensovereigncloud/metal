// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha4

// ObjectReferenceApplyConfiguration represents an declarative configuration of the ObjectReference type for use
// with apply.
type ObjectReferenceApplyConfiguration struct {
	Name      *string `json:"name,omitempty"`
	Namespace *string `json:"namespace,omitempty"`
}

// ObjectReferenceApplyConfiguration constructs an declarative configuration of the ObjectReference type for use with
// apply.
func ObjectReference() *ObjectReferenceApplyConfiguration {
	return &ObjectReferenceApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *ObjectReferenceApplyConfiguration) WithName(value string) *ObjectReferenceApplyConfiguration {
	b.Name = &value
	return b
}

// WithNamespace sets the Namespace field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Namespace field is set to the value of the last call.
func (b *ObjectReferenceApplyConfiguration) WithNamespace(value string) *ObjectReferenceApplyConfiguration {
	b.Namespace = &value
	return b
}
