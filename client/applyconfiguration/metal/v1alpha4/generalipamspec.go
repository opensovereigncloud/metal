// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha4

// GeneralIPAMSpecApplyConfiguration represents an declarative configuration of the GeneralIPAMSpec type for use
// with apply.
type GeneralIPAMSpecApplyConfiguration struct {
	AddressFamily     *AddressFamiliesMapApplyConfiguration `json:"addressFamily,omitempty"`
	CarrierSubnets    *IPAMSelectionSpecApplyConfiguration  `json:"carrierSubnets,omitempty"`
	LoopbackSubnets   *IPAMSelectionSpecApplyConfiguration  `json:"loopbackSubnets,omitempty"`
	SouthSubnets      *IPAMSelectionSpecApplyConfiguration  `json:"southSubnets,omitempty"`
	LoopbackAddresses *IPAMSelectionSpecApplyConfiguration  `json:"loopbackAddresses,omitempty"`
}

// GeneralIPAMSpecApplyConfiguration constructs an declarative configuration of the GeneralIPAMSpec type for use with
// apply.
func GeneralIPAMSpec() *GeneralIPAMSpecApplyConfiguration {
	return &GeneralIPAMSpecApplyConfiguration{}
}

// WithAddressFamily sets the AddressFamily field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AddressFamily field is set to the value of the last call.
func (b *GeneralIPAMSpecApplyConfiguration) WithAddressFamily(value *AddressFamiliesMapApplyConfiguration) *GeneralIPAMSpecApplyConfiguration {
	b.AddressFamily = value
	return b
}

// WithCarrierSubnets sets the CarrierSubnets field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CarrierSubnets field is set to the value of the last call.
func (b *GeneralIPAMSpecApplyConfiguration) WithCarrierSubnets(value *IPAMSelectionSpecApplyConfiguration) *GeneralIPAMSpecApplyConfiguration {
	b.CarrierSubnets = value
	return b
}

// WithLoopbackSubnets sets the LoopbackSubnets field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LoopbackSubnets field is set to the value of the last call.
func (b *GeneralIPAMSpecApplyConfiguration) WithLoopbackSubnets(value *IPAMSelectionSpecApplyConfiguration) *GeneralIPAMSpecApplyConfiguration {
	b.LoopbackSubnets = value
	return b
}

// WithSouthSubnets sets the SouthSubnets field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SouthSubnets field is set to the value of the last call.
func (b *GeneralIPAMSpecApplyConfiguration) WithSouthSubnets(value *IPAMSelectionSpecApplyConfiguration) *GeneralIPAMSpecApplyConfiguration {
	b.SouthSubnets = value
	return b
}

// WithLoopbackAddresses sets the LoopbackAddresses field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LoopbackAddresses field is set to the value of the last call.
func (b *GeneralIPAMSpecApplyConfiguration) WithLoopbackAddresses(value *IPAMSelectionSpecApplyConfiguration) *GeneralIPAMSpecApplyConfiguration {
	b.LoopbackAddresses = value
	return b
}
