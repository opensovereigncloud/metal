// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha4

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConnectionsMap map[uint8]*NetworkSwitchList

// PortParametersSpec contains a set of parameters of switch port.
type PortParametersSpec struct {
	// Lanes refers to a number of lanes used by switch port
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=8
	Lanes *uint32 `json:"lanes,omitempty"`
	// MTU refers to maximum transmission unit value which should be applied on switch port
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=84
	// +kubebuilder:validation:Maximum=65535
	MTU *uint32 `json:"mtu,omitempty"`
	// IPv4MaskLength defines prefix of subnet where switch port's IPv4 address should be reserved
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=32
	IPv4MaskLength *uint32 `json:"ipv4MaskLength,omitempty"`
	// IPv6Prefix defines prefix of subnet where switch port's IPv6 address should be reserved
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=128
	IPv6Prefix *uint32 `json:"ipv6Prefix,omitempty"`
	// FEC refers to forward error correction method which should be applied on switch port
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=rs;none
	FEC *string `json:"fec,omitempty"`
	// State defines default state of switch port
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=up;down
	State *string `json:"state,omitempty"`
}

// IPAMSpec contains selectors for subnets and loopback IPs and
// definition of address families which should be claimed.
type IPAMSpec struct {
	// SouthSubnets defines selector for subnet object which will be assigned to switch
	// +kubebuilder:validation:Optional
	SouthSubnets *IPAMSelectionSpec `json:"southSubnets,omitempty"`
	// LoopbackAddresses defines selector for IP object which will be assigned to switch's loopback interface
	// +kubebuilder:validation:Optional
	LoopbackAddresses *IPAMSelectionSpec `json:"loopbackAddresses,omitempty"`
}

// IPAMSelectionSpec contains label selector and address family.
type IPAMSelectionSpec struct {
	// LabelSelector contains label selector to pick up IPAM objects
	// +kubebuilder:validation:Optional
	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty"`
	// FieldSelector contains label key and field path where to get label value for search.
	// If FieldSelector is used as part of IPAM configuration in SwitchConfig object it will
	// reference to the field path in related NetworkSwitch object. If FieldSelector is used as part of IPAM
	// configuration in NetworkSwitch object, it will reference to the field path in the same object
	// +kubebuilder:validation:Optional
	FieldSelector *FieldSelectorSpec `json:"fieldSelector,omitempty"`
}

// FieldSelectorSpec contains label key and field path where to get label value for search.
type FieldSelectorSpec struct {
	// LabelKey contains label key
	// +kubebuilder:validation:Optional
	LabelKey *string `json:"labelKey"`
	// FieldRef contains reference to the field of resource where to get label's value
	// +kubebuilder:validation:Optional
	FieldRef *v1.ObjectFieldSelector `json:"fieldRef"`
}

// AddressFamiliesMap contains flags regarding what IP address families should be used.
type AddressFamiliesMap struct {
	// IPv4 is a flag defining whether IPv4 is used or not
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	IPv4 *bool `json:"ipv4"`
	// IPv6 is a flag defining whether IPv6 is used or not
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	IPv6 *bool `json:"ipv6"`
}

// AdditionalIPSpec defines IP address and selector for subnet where address should be reserved.
type AdditionalIPSpec struct {
	// Address contains additional IP address that should be assigned to the interface
	// +kubebuilder:validation:Required
	Address *string `json:"address,omitempty"`
	// ParentSubnet contains label selector to pick up IPAM objects
	// +kubebuilder:validation:Optional
	ParentSubnet *metav1.LabelSelector `json:"parentSubnet,omitempty"`
}

// ObjectReference contains enough information to let you locate the
// referenced object across namespaces.
type ObjectReference struct {
	// Name contains name of the referenced object
	// +kubebuilder:validation:Optional
	Name *string `json:"name,omitempty"`
	// Namespace contains namespace of the referenced object
	// +kubebuilder:validation:Optional
	Namespace *string `json:"namespace,omitempty"`
}

// GetLabelKey builds label key from prefix and suffix.
// func GetLabelKey(prefix, suffix string) string {
// 	return fmt.Sprintf("%s/%s", prefix, suffix)
// }
