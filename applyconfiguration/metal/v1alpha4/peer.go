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

package v1alpha4

// PeerApplyConfiguration represents an declarative configuration of the Peer type for use
// with apply.
type PeerApplyConfiguration struct {
	LLDPSystemName      *string                              `json:"lldp_system_name,omitempty"`
	LLDPChassisID       *string                              `json:"lldp_chassis_id,omitempty"`
	LLDPPortID          *string                              `json:"lldp_port_id,omitempty"`
	LLDPPortDescription *string                              `json:"lldp_port_description,omitempty"`
	ResourceReference   *ResourceReferenceApplyConfiguration `json:"resource_reference,omitempty"`
}

// PeerApplyConfiguration constructs an declarative configuration of the Peer type for use with
// apply.
func Peer() *PeerApplyConfiguration {
	return &PeerApplyConfiguration{}
}

// WithLLDPSystemName sets the LLDPSystemName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LLDPSystemName field is set to the value of the last call.
func (b *PeerApplyConfiguration) WithLLDPSystemName(value string) *PeerApplyConfiguration {
	b.LLDPSystemName = &value
	return b
}

// WithLLDPChassisID sets the LLDPChassisID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LLDPChassisID field is set to the value of the last call.
func (b *PeerApplyConfiguration) WithLLDPChassisID(value string) *PeerApplyConfiguration {
	b.LLDPChassisID = &value
	return b
}

// WithLLDPPortID sets the LLDPPortID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LLDPPortID field is set to the value of the last call.
func (b *PeerApplyConfiguration) WithLLDPPortID(value string) *PeerApplyConfiguration {
	b.LLDPPortID = &value
	return b
}

// WithLLDPPortDescription sets the LLDPPortDescription field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LLDPPortDescription field is set to the value of the last call.
func (b *PeerApplyConfiguration) WithLLDPPortDescription(value string) *PeerApplyConfiguration {
	b.LLDPPortDescription = &value
	return b
}

// WithResourceReference sets the ResourceReference field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ResourceReference field is set to the value of the last call.
func (b *PeerApplyConfiguration) WithResourceReference(value *ResourceReferenceApplyConfiguration) *PeerApplyConfiguration {
	b.ResourceReference = value
	return b
}