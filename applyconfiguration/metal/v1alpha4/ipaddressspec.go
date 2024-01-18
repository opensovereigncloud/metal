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

// IPAddressSpecApplyConfiguration represents an declarative configuration of the IPAddressSpec type for use
// with apply.
type IPAddressSpecApplyConfiguration struct {
	ObjectReferenceApplyConfiguration `json:",inline"`
	Address                           *string `json:"address,omitempty"`
	ExtraAddress                      *bool   `json:"extraAddress,omitempty"`
	AddressFamily                     *string `json:"addressFamily,omitempty"`
}

// IPAddressSpecApplyConfiguration constructs an declarative configuration of the IPAddressSpec type for use with
// apply.
func IPAddressSpec() *IPAddressSpecApplyConfiguration {
	return &IPAddressSpecApplyConfiguration{}
}

// WithAddress sets the Address field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Address field is set to the value of the last call.
func (b *IPAddressSpecApplyConfiguration) WithAddress(value string) *IPAddressSpecApplyConfiguration {
	b.Address = &value
	return b
}

// WithExtraAddress sets the ExtraAddress field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ExtraAddress field is set to the value of the last call.
func (b *IPAddressSpecApplyConfiguration) WithExtraAddress(value bool) *IPAddressSpecApplyConfiguration {
	b.ExtraAddress = &value
	return b
}

// WithAddressFamily sets the AddressFamily field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AddressFamily field is set to the value of the last call.
func (b *IPAddressSpecApplyConfiguration) WithAddressFamily(value string) *IPAddressSpecApplyConfiguration {
	b.AddressFamily = &value
	return b
}