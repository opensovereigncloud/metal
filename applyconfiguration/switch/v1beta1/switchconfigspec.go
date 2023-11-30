/*
Copyright (c) 2023 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SwitchConfigSpecApplyConfiguration represents an declarative configuration of the SwitchConfigSpec type for use
// with apply.
type SwitchConfigSpecApplyConfiguration struct {
	Switches              *v1.LabelSelector                     `json:"switches,omitempty"`
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
func (b *SwitchConfigSpecApplyConfiguration) WithSwitches(value v1.LabelSelector) *SwitchConfigSpecApplyConfiguration {
	b.Switches = &value
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
