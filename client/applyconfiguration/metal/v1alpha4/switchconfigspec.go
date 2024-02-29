// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha4

import (
	v1 "github.com/ironcore-dev/metal/client/applyconfiguration/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

// SwitchConfigSpecApplyConfiguration represents an declarative configuration of the SwitchConfigSpec type for use
// with apply.
type SwitchConfigSpecApplyConfiguration struct {
	Switches              *v1.LabelSelectorApplyConfiguration   `json:"switches,omitempty"`
	PortsDefaults         *PortParametersSpecApplyConfiguration `json:"portsDefaults,omitempty"`
	IPAM                  *GeneralIPAMSpecApplyConfiguration    `json:"ipam,omitempty"`
	RoutingConfigTemplate *corev1.LocalObjectReference          `json:"routingConfigTemplate,omitempty"`
}

// SwitchConfigSpecApplyConfiguration constructs an declarative configuration of the SwitchConfigSpec type for use with
// apply.
func SwitchConfigSpec() *SwitchConfigSpecApplyConfiguration {
	return &SwitchConfigSpecApplyConfiguration{}
}

// WithSwitches sets the Switches field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Switches field is set to the value of the last call.
func (b *SwitchConfigSpecApplyConfiguration) WithSwitches(value *v1.LabelSelectorApplyConfiguration) *SwitchConfigSpecApplyConfiguration {
	b.Switches = value
	return b
}

// WithPortsDefaults sets the PortsDefaults field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PortsDefaults field is set to the value of the last call.
func (b *SwitchConfigSpecApplyConfiguration) WithPortsDefaults(value *PortParametersSpecApplyConfiguration) *SwitchConfigSpecApplyConfiguration {
	b.PortsDefaults = value
	return b
}

// WithIPAM sets the IPAM field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the IPAM field is set to the value of the last call.
func (b *SwitchConfigSpecApplyConfiguration) WithIPAM(value *GeneralIPAMSpecApplyConfiguration) *SwitchConfigSpecApplyConfiguration {
	b.IPAM = value
	return b
}

// WithRoutingConfigTemplate sets the RoutingConfigTemplate field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RoutingConfigTemplate field is set to the value of the last call.
func (b *SwitchConfigSpecApplyConfiguration) WithRoutingConfigTemplate(value corev1.LocalObjectReference) *SwitchConfigSpecApplyConfiguration {
	b.RoutingConfigTemplate = &value
	return b
}
