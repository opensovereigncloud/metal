// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha4

import (
	"net/netip"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MachineState string

const (
	// MachineStateHealthy - When State is `Healthy` Machine` is allowed to be booked.
	MachineStateHealthy MachineState = "Healthy"
	// MachineStateUnhealthy - When State is `Unhealthy`` Machine isn't allowed to be booked.
	MachineStateUnhealthy MachineState = "Unhealthy"
)

// MachineSpec - defines the desired spec of Machine.
type MachineSpec struct {
	// Hostname - defines machine domain name
	// +optional
	Hostname string `json:"hostname,omitempty"`
	// Description - summary info about machine
	// +optional
	Description string `json:"description,omitempty"`
	// Identity - defines machine hardware info
	// +optional
	Identity Identity `json:"identity,omitempty"`
	// InventoryRequested - defines if inventory requested or not
	InventoryRequested bool `json:"inventory_requested,omitempty"`
}

// Identity - defines hardware information about machine.
type Identity struct {
	// SKU - stock keeping unit. The label allows vendors automatically track the movement of inventory
	// +optional
	SKU string `json:"sku,omitempty"`
	// SerialNumber - unique machine number
	// +optional
	SerialNumber string `json:"serial_number,omitempty"`
	// +optional
	Asset string `json:"asset,omitempty"`
	// Deprecated
	Internal []Internal `json:"internal,omitempty"`
}

type Internal struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// MachineStatus - defines machine aggregated info.
type MachineStatus struct {
	// Reboot - defines machine reboot status
	// +optional
	Reboot string `json:"reboot,omitempty"`
	// Health - defines machine condition.
	// "healthy" if both OOB and Inventory are presented and "unhealthy" if one of them isn't
	// +optional
	Health MachineState `json:"health,omitempty"`
	// Network - defines machine network status
	// +optional
	Network Network `json:"network,omitempty"`
	// Reservation - defines machine reservation state and reference object.
	// +optional
	Reservation Reservation `json:"reservation,omitempty"`
	// Orphaned - defines machine condition whether OOB or Inventory is missing or not
	// +optional
	Orphaned bool `json:"orphaned,omitempty"`
}

// Peer - contains machine neighbor information collected from LLDP.
type Peer struct {
	// LLDPSystemName - defines switch name obtained from Link Layer Discovery Protocol
	// layer 2 neighbor discovery protocol
	// +optional
	LLDPSystemName string `json:"lldp_system_name,omitempty"`
	// LLDPChassisID - defines switch ID for chassis obtained from Link Layer Discovery Protocol
	// +optional
	LLDPChassisID string `json:"lldp_chassis_id,omitempty"`
	// LLDPPortID - defines switch port ID obtained from Link Layer Discovery Protocol
	// +optional
	LLDPPortID string `json:"lldp_port_id,omitempty"`
	// LLDPPortDescription - defines switch definition obtained from Link Layer Discovery Protocol
	// +optional
	LLDPPortDescription string `json:"lldp_port_description,omitempty"`
	// ResourceReference refers to the related resource definition
	// +optional
	ResourceReference *ResourceReference `json:"resource_reference,omitempty"`
}

// Network - defines machine network status.
type Network struct {
	// ASN - defines calculated Autonomous system Number.
	ASN uint32 `json:"asn,omitempty"`
	// Redundancy - defines machine redundancy status.
	// Available values: "Single", "HighAvailability" or "None"
	// +kubebuilder:validation:Pattern=`^(?:Single|HighAvailability|None)$`
	Redundancy string `json:"redundancy,omitempty"`
	// Ports - defines number of machine ports
	// +kubebuilder:validation:Optional
	Ports int `json:"ports,omitempty"`
	// UnknownPorts - defines number of machine interface without info
	// +kubebuilder:validation:Optional
	UnknownPorts int `json:"unknown_ports,omitempty"`
	// Interfaces - defines machine interfaces info
	// +kubebuilder:validation:Optional
	Interfaces []Interface `json:"interfaces,omitempty"`
	// Loopbacks refers to the switch's loopback addresses
	// +kubebuilder:validation:Optional
	Loopbacks LoopbackAddresses `json:"loopback_addresses,omitempty"`
}

// Interface - defines information about machine interfaces.
type Interface struct {
	// Name - machine interface name
	// +optional
	Name string `json:"name,omitempty"`
	// SwitchReference - defines unique switch identification
	// +optional
	SwitchReference *ResourceReference `json:"switch_reference,omitempty"`
	// IPv4 - defines machine IPv4 address
	// +optional
	Addresses Addresses `json:"addresses,omitempty"`
	// Peer - defines lldp peer info.
	// +optional
	Peer Peer `json:"peer,omitempty"`
	// Lane - defines number of lines per interface
	// +optional
	Lanes uint32 `json:"lanes,omitempty"`
	// Moved  - defines if interface was reconnected to another switch or not
	// +optional
	Moved bool `json:"moved,omitempty"`
	// Unknown - defines information availability about interface
	// +optional
	Unknown bool `json:"unknown,omitempty"`
}

type Addresses struct {
	IPv4 []IPAddrSpec `json:"ipv4,omitempty"`
	IPv6 []IPAddrSpec `json:"ipv6,omitempty"`
}

type LoopbackAddresses struct {
	IPv4 IPAddrSpec `json:"ipv4,omitempty"`
	IPv6 IPAddrSpec `json:"ipv6,omitempty"`
}

// IPAddrSpec defines interface's ip address info.
// +kubebuilder:validation:Type=string
type IPAddrSpec struct {
	// Address refers to the ip address value
	netip.Prefix `json:"-"`
}

type Reservation struct {
	// Status - defines Machine Order state provided by OOB Machine Resources
	// +optional
	Status string `json:"status,omitempty"`
	// Class - defines what class the mahchine was reserved under
	// +optional
	Class string `json:"class,omitempty"`
	// Reference - defines underlying referenced object.
	// +optional
	Reference *ResourceReference `json:"reference,omitempty"`
}

// ResourceReference defines related resource info.
type ResourceReference struct {
	// APIVersion refers to the resource API version
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind refers to the resource kind
	// +optional
	Kind string `json:"kind,omitempty"`
	// Name refers to the resource name
	// +optional
	Name string `json:"name,omitempty"`
	// Namespace refers to the resource namespace
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ASN",type=string,JSONPath=`.status.network.asn`
// +kubebuilder:printcolumn:name="Healthy",type=string,JSONPath=`.status.health`
// +kubebuilder:printcolumn:name="Redundancy",type=string,JSONPath=`.status.network.redundancy`
// +kubebuilder:printcolumn:name="Reservation Status",type=string,JSONPath=`.status.reservation.status`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// Machine - is the data structure for a Machine resource.
// It contains an aggregated information from Inventory and OOB resources.
type Machine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineSpec   `json:"spec,omitempty"`
	Status MachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineList - contains a list of Machine.
type MachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Machine `json:"items"`
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPAddrSpec) DeepCopyInto(out *IPAddrSpec) {
	*out = *in
	if in.String() != "" {
		out.Prefix = in.Prefix
	}
}
