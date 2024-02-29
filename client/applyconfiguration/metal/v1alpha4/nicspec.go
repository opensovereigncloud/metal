// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha4

// NICSpecApplyConfiguration represents an declarative configuration of the NICSpec type for use
// with apply.
type NICSpecApplyConfiguration struct {
	Name       *string                      `json:"name,omitempty"`
	PCIAddress *string                      `json:"pciAddress,omitempty"`
	MACAddress *string                      `json:"macAddress,omitempty"`
	MTU        *uint16                      `json:"mtu,omitempty"`
	Speed      *uint32                      `json:"speed,omitempty"`
	Lanes      *byte                        `json:"lanes,omitempty"`
	ActiveFEC  *string                      `json:"activeFEC,omitempty"`
	LLDPs      []LLDPSpecApplyConfiguration `json:"lldps,omitempty"`
	NDPs       []NDPSpecApplyConfiguration  `json:"ndps,omitempty"`
}

// NICSpecApplyConfiguration constructs an declarative configuration of the NICSpec type for use with
// apply.
func NICSpec() *NICSpecApplyConfiguration {
	return &NICSpecApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *NICSpecApplyConfiguration) WithName(value string) *NICSpecApplyConfiguration {
	b.Name = &value
	return b
}

// WithPCIAddress sets the PCIAddress field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PCIAddress field is set to the value of the last call.
func (b *NICSpecApplyConfiguration) WithPCIAddress(value string) *NICSpecApplyConfiguration {
	b.PCIAddress = &value
	return b
}

// WithMACAddress sets the MACAddress field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MACAddress field is set to the value of the last call.
func (b *NICSpecApplyConfiguration) WithMACAddress(value string) *NICSpecApplyConfiguration {
	b.MACAddress = &value
	return b
}

// WithMTU sets the MTU field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MTU field is set to the value of the last call.
func (b *NICSpecApplyConfiguration) WithMTU(value uint16) *NICSpecApplyConfiguration {
	b.MTU = &value
	return b
}

// WithSpeed sets the Speed field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Speed field is set to the value of the last call.
func (b *NICSpecApplyConfiguration) WithSpeed(value uint32) *NICSpecApplyConfiguration {
	b.Speed = &value
	return b
}

// WithLanes sets the Lanes field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Lanes field is set to the value of the last call.
func (b *NICSpecApplyConfiguration) WithLanes(value byte) *NICSpecApplyConfiguration {
	b.Lanes = &value
	return b
}

// WithActiveFEC sets the ActiveFEC field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ActiveFEC field is set to the value of the last call.
func (b *NICSpecApplyConfiguration) WithActiveFEC(value string) *NICSpecApplyConfiguration {
	b.ActiveFEC = &value
	return b
}

// WithLLDPs adds the given value to the LLDPs field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the LLDPs field.
func (b *NICSpecApplyConfiguration) WithLLDPs(values ...*LLDPSpecApplyConfiguration) *NICSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithLLDPs")
		}
		b.LLDPs = append(b.LLDPs, *values[i])
	}
	return b
}

// WithNDPs adds the given value to the NDPs field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the NDPs field.
func (b *NICSpecApplyConfiguration) WithNDPs(values ...*NDPSpecApplyConfiguration) *NICSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNDPs")
		}
		b.NDPs = append(b.NDPs, *values[i])
	}
	return b
}
