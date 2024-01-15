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
