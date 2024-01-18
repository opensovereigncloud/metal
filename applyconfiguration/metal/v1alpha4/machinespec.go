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

// MachineSpecApplyConfiguration represents an declarative configuration of the MachineSpec type for use
// with apply.
type MachineSpecApplyConfiguration struct {
	Hostname           *string                     `json:"hostname,omitempty"`
	Description        *string                     `json:"description,omitempty"`
	Identity           *IdentityApplyConfiguration `json:"identity,omitempty"`
	InventoryRequested *bool                       `json:"inventory_requested,omitempty"`
}

// MachineSpecApplyConfiguration constructs an declarative configuration of the MachineSpec type for use with
// apply.
func MachineSpec() *MachineSpecApplyConfiguration {
	return &MachineSpecApplyConfiguration{}
}

// WithHostname sets the Hostname field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Hostname field is set to the value of the last call.
func (b *MachineSpecApplyConfiguration) WithHostname(value string) *MachineSpecApplyConfiguration {
	b.Hostname = &value
	return b
}

// WithDescription sets the Description field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Description field is set to the value of the last call.
func (b *MachineSpecApplyConfiguration) WithDescription(value string) *MachineSpecApplyConfiguration {
	b.Description = &value
	return b
}

// WithIdentity sets the Identity field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Identity field is set to the value of the last call.
func (b *MachineSpecApplyConfiguration) WithIdentity(value *IdentityApplyConfiguration) *MachineSpecApplyConfiguration {
	b.Identity = value
	return b
}

// WithInventoryRequested sets the InventoryRequested field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the InventoryRequested field is set to the value of the last call.
func (b *MachineSpecApplyConfiguration) WithInventoryRequested(value bool) *MachineSpecApplyConfiguration {
	b.InventoryRequested = &value
	return b
}