// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha4

// PeerSpecApplyConfiguration represents an declarative configuration of the PeerSpec type for use
// with apply.
type PeerSpecApplyConfiguration struct {
	ObjectReferenceApplyConfiguration `json:",inline"`
	PeerInfoSpecApplyConfiguration    `json:",inline"`
}

// PeerSpecApplyConfiguration constructs an declarative configuration of the PeerSpec type for use with
// apply.
func PeerSpec() *PeerSpecApplyConfiguration {
	return &PeerSpecApplyConfiguration{}
}
