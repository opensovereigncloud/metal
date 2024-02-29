// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha4

// InventoryStatusesApplyConfiguration represents an declarative configuration of the InventoryStatuses type for use
// with apply.
type InventoryStatusesApplyConfiguration struct {
	Ready         *bool `json:"ready,omitempty"`
	RequestsCount *int  `json:"requestsCount,omitempty"`
}

// InventoryStatusesApplyConfiguration constructs an declarative configuration of the InventoryStatuses type for use with
// apply.
func InventoryStatuses() *InventoryStatusesApplyConfiguration {
	return &InventoryStatusesApplyConfiguration{}
}

// WithReady sets the Ready field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Ready field is set to the value of the last call.
func (b *InventoryStatusesApplyConfiguration) WithReady(value bool) *InventoryStatusesApplyConfiguration {
	b.Ready = &value
	return b
}

// WithRequestsCount sets the RequestsCount field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RequestsCount field is set to the value of the last call.
func (b *InventoryStatusesApplyConfiguration) WithRequestsCount(value int) *InventoryStatusesApplyConfiguration {
	b.RequestsCount = &value
	return b
}
