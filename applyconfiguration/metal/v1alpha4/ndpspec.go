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

// NDPSpecApplyConfiguration represents an declarative configuration of the NDPSpec type for use
// with apply.
type NDPSpecApplyConfiguration struct {
	IPAddress  *string `json:"ipAddress,omitempty"`
	MACAddress *string `json:"macAddress,omitempty"`
	State      *string `json:"state,omitempty"`
}

// NDPSpecApplyConfiguration constructs an declarative configuration of the NDPSpec type for use with
// apply.
func NDPSpec() *NDPSpecApplyConfiguration {
	return &NDPSpecApplyConfiguration{}
}

// WithIPAddress sets the IPAddress field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the IPAddress field is set to the value of the last call.
func (b *NDPSpecApplyConfiguration) WithIPAddress(value string) *NDPSpecApplyConfiguration {
	b.IPAddress = &value
	return b
}

// WithMACAddress sets the MACAddress field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MACAddress field is set to the value of the last call.
func (b *NDPSpecApplyConfiguration) WithMACAddress(value string) *NDPSpecApplyConfiguration {
	b.MACAddress = &value
	return b
}

// WithState sets the State field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the State field is set to the value of the last call.
func (b *NDPSpecApplyConfiguration) WithState(value string) *NDPSpecApplyConfiguration {
	b.State = &value
	return b
}