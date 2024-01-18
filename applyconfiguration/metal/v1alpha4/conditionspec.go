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

// ConditionSpecApplyConfiguration represents an declarative configuration of the ConditionSpec type for use
// with apply.
type ConditionSpecApplyConfiguration struct {
	Name                    *string `json:"name,omitempty"`
	State                   *bool   `json:"state,omitempty"`
	LastUpdateTimestamp     *string `json:"lastUpdateTimestamp,omitempty"`
	LastTransitionTimestamp *string `json:"lastTransitionTimestamp,omitempty"`
	Reason                  *string `json:"reason,omitempty"`
	Message                 *string `json:"message,omitempty"`
}

// ConditionSpecApplyConfiguration constructs an declarative configuration of the ConditionSpec type for use with
// apply.
func ConditionSpec() *ConditionSpecApplyConfiguration {
	return &ConditionSpecApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *ConditionSpecApplyConfiguration) WithName(value string) *ConditionSpecApplyConfiguration {
	b.Name = &value
	return b
}

// WithState sets the State field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the State field is set to the value of the last call.
func (b *ConditionSpecApplyConfiguration) WithState(value bool) *ConditionSpecApplyConfiguration {
	b.State = &value
	return b
}

// WithLastUpdateTimestamp sets the LastUpdateTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LastUpdateTimestamp field is set to the value of the last call.
func (b *ConditionSpecApplyConfiguration) WithLastUpdateTimestamp(value string) *ConditionSpecApplyConfiguration {
	b.LastUpdateTimestamp = &value
	return b
}

// WithLastTransitionTimestamp sets the LastTransitionTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LastTransitionTimestamp field is set to the value of the last call.
func (b *ConditionSpecApplyConfiguration) WithLastTransitionTimestamp(value string) *ConditionSpecApplyConfiguration {
	b.LastTransitionTimestamp = &value
	return b
}

// WithReason sets the Reason field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Reason field is set to the value of the last call.
func (b *ConditionSpecApplyConfiguration) WithReason(value string) *ConditionSpecApplyConfiguration {
	b.Reason = &value
	return b
}

// WithMessage sets the Message field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Message field is set to the value of the last call.
func (b *ConditionSpecApplyConfiguration) WithMessage(value string) *ConditionSpecApplyConfiguration {
	b.Message = &value
	return b
}