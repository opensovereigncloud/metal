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

package v1alpha3

// ReservationApplyConfiguration represents an declarative configuration of the Reservation type for use
// with apply.
type ReservationApplyConfiguration struct {
	Status    *string                              `json:"status,omitempty"`
	Class     *string                              `json:"class,omitempty"`
	Reference *ResourceReferenceApplyConfiguration `json:"reference,omitempty"`
}

// ReservationApplyConfiguration constructs an declarative configuration of the Reservation type for use with
// apply.
func Reservation() *ReservationApplyConfiguration {
	return &ReservationApplyConfiguration{}
}

// WithStatus sets the Status field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Status field is set to the value of the last call.
func (b *ReservationApplyConfiguration) WithStatus(value string) *ReservationApplyConfiguration {
	b.Status = &value
	return b
}

// WithClass sets the Class field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Class field is set to the value of the last call.
func (b *ReservationApplyConfiguration) WithClass(value string) *ReservationApplyConfiguration {
	b.Class = &value
	return b
}

// WithReference sets the Reference field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Reference field is set to the value of the last call.
func (b *ReservationApplyConfiguration) WithReference(value *ResourceReferenceApplyConfiguration) *ReservationApplyConfiguration {
	b.Reference = value
	return b
}
