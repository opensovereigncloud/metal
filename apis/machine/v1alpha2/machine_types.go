/*
Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TolerationOperator - is the set of operators that can be used in a toleration.
type TolerationOperator string

const (
	TolerationOpEqual  TolerationOperator = "Equal"
	TolerationOpExists TolerationOperator = "Exists"
)

type TaintEffect string

const (
	// When Machine taint effect is NotAvailable that's mean that Inventory or OOB not exist.
	TaintEffectNotAvailable TaintEffect = "NotAvailable"
	// When Machine taint effect is Suspended.
	TaintEffectSuspended TaintEffect = "Suspended"
	// When Machine taint effect is NoSchedule.
	TaintEffectNoSchedule TaintEffect = "NoSchedule"
	// When Machine taint effect is Error it's impossible to order machine. And it requires to run stresstest.
	TaintEffectError TaintEffect = "Error"
)

type MachineState string

const (
	// MachineStateHealthy - When State is `Healthy` Machine` is allowed to be booked.
	MachineStateHealthy MachineState = "Healthy"
	// MachineStateUnhealthy - When State is `Unhealthy`` Machine isn't allowed to be booked.
	MachineStateUnhealthy MachineState = "Unhealthy"
)

const (
	InterfaceRedundancySingle           = "Single"
	InterfaceRedundancyHighAvailability = "HighAvailability"
	InterfaceRedundancyNone             = "None"
)

const (
	ReservationStatusAvailable = "Available"
	ReservationStatusReserved  = "Reserved"
	ReservationStatusPending   = "Pending"
	ReservationStatusError     = "Error"
	ReservationStatusRunning   = "Running"
)

// MachineSpec - defines the desired spec of Machine.
type MachineSpec struct {
	// Taints - defines list of Taint that applied on the Machine
	// +optional
	Taints []Taint `json:"taints,omitempty"`
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

// Taint represents taint that can be applied to the machine.
// The machine this Taint is attached to has the "effect" on
// any pod that does not tolerate the Taint.
type Taint struct {
	// Key - applied to the machine.
	// +required
	Key string `json:"key"`
	// Value - corresponding to the taint key.
	// +optional
	Value string `json:"value,omitempty"`
	// Effect - defines taint effect on the Machine.
	// Valid effects are NotAvailable and Suspended.
	// +required
	Effect TaintEffect `json:"status"`
	// TimeAdded represents the time at which the taint was added.
	// It is only written for NoExecute taints.
	// +optional
	TimeAdded *metav1.Time `json:"time_added,omitempty"`
}

// The resource this Toleration is attached to tolerates any taint that matches
// the triple <key,value,effect> using the matching operator <operator>.
type Toleration struct {
	// Key is the taint key that the toleration applies to. Empty means match all taint keys.
	// If the key is empty, operator must be Exists; this combination means to match all values and all keys.
	Key string `json:"key,omitempty"`
	// Operator represents a key's relationship to the value.
	// Valid operators are Exists and Equal. Defaults to Equal.
	// Exist is equivalent to wildcard for value, so that a resource can
	// tolerate all taints of a particular category.
	Operator TolerationOperator `json:"operator,omitempty"`
	// Value is the taint value the toleration matches to.
	// If the operator is Exists, the value should be empty, otherwise just a regular string.
	Value string `json:"value,omitempty"`
	// Effect indicates the taint effect to match. Empty means match all taint effects.
	// When specified, allowed values are NoSchedule.
	Effect TaintEffect `json:"effect,omitempty"`
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
	Asset    string     `json:"asset,omitempty"`
	Internal []Internal `json:"internal,omitempty"`
}

type Internal struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// MachineStatus - defines machine aggregated info.
type MachineStatus struct {
	// Interfaces - defines machine interfaces info
	// +optional
	Interfaces []Interface `json:"interfaces,omitempty"`
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
	// // RequestState - defines Machine Request state provided by OOB Machine Resources
	// // +optional
	// RequestState RequestState `json:"request_state,omitempty"`
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
	IPv4 *IPAddressSpec `json:"ipv4,omitempty"`
	// IPv6 - defines machine IPv6 address
	// +optional
	IPv6 *IPAddressSpec `json:"ipv6,omitempty"`
	// Peer - defines lldp peer info.
	// +optional
	Peer *Peer `json:"peer,omitempty"`
	// Lane - defines number of lines per interface
	// +optional
	Lanes uint8 `json:"lanes,omitempty"`
	// Moved  - defines if interface was reconnected to another switch or not
	// +optional
	Moved bool `json:"moved,omitempty"`
	// Unknown - defines information availability about interface
	// +optional
	Unknown bool `json:"unknown,omitempty"`
}

// IPAddressSpec defines interface's ip address info.
type IPAddressSpec struct {
	// Address refers to the ip address value
	Address string `json:"address,omitempty"`
	// ResourceReference refers to the related resource definition
	ResourceReference *ResourceReference `json:"resource_reference,omitempty"`
}

// Peer - contains machine neighbor information collected from LLDP.
type Peer struct {
	// LLDPSystemName - defines switch name obtained from Link Layer Discovery Protocol
	// layer 2 neighbor discovery protocol
	// +optional
	LLDPSystemName string `json:"lldp_system_name,omitempty"`
	// LLDPChassisID - defines switch ID for chassis obtained from Link Layer Discovery Protocol
	// +optional
	LLDPChassisID string `json:"lldp_chassi_id,omitempty"`
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
	// Redundancy - defines machine redundancy status.
	// Available values: "Single", "HighAvailability" or "None"
	// +kubebuilder:validation:Pattern=`^(?:Single|HighAvailability|None)$`
	// +optional
	Redundancy string `json:"redundancy,omitempty"`
	// Ports - defines number of machine ports
	// +optional
	Ports int `json:"ports,omitempty"`
	// UnknownPorts - defines number of machine interface without info
	// +optional
	UnknownPorts int `json:"unknown_ports,omitempty"`
}

// ObjectReference - defines object reference status and additional information.
type ObjectReference struct {
	// Exist - defines where referenced object exist or not
	// +optional
	Exist bool `json:"exist,omitempty"`
	// Reference - defines underlaying referenced object e.g. Inventory or OOB kind.
	// +optional
	Reference *ResourceReference `json:"reference,omitempty"`
}

type Reservation struct {
	// Status - defines Machine Order state provided by OOB Machine Resources
	// +optional
	Status string `json:"status,omitempty"`
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
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Healthy",type=string,JSONPath=`.status.health`
// +kubebuilder:printcolumn:name="Redundancy",type=string,JSONPath=`.status.network.redundancy`
// +kubebuilder:printcolumn:name="Reservation Status",type=string,JSONPath=`.status.reservation.status`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Machine - is the data structure for a Machine resource.
// It contains an aggregated information from Inventory and OOB resources.
type Machine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineSpec   `json:"spec,omitempty"`
	Status MachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MachineList - contains a list of Machine.
type MachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Machine `json:"items"`
}
