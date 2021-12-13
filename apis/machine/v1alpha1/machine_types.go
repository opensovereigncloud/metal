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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	MachinPowerStateON     = "On"
	MachinePowerStateOFF   = "Off"
	MachinePowerStateReset = "Reset"
)

// MachineSpec - defines the desired spec of Machine.
type MachineSpec struct {
	// Hostname - defines machine domain name
	//+optional
	Hostname string `json:"hostname,omitempty"`
	// Description - summary info about machine
	//+optional
	Description string `json:"description,omitempty"`
	// Location - defines machine location in datacenter
	//+optional
	Location Location `json:"location,omitempty"`
	// Action - defines desired operation on machine
	//+optional
	Action Action `json:"action,omitempty"`
	// Identity - defines machine hardware info
	//+optional
	Identity Identity `json:"identity,omitempty"`
	// ScanPorts - trigger manual port scan
	// +kubebuilder:default:=false
	ScanPorts bool `json:"scan_ports,omitempty"`
	// InventoryRequested - defines if inventory requested or not
	InventoryRequested bool `json:"inventory_requested,omitempty"`
}

// Location - defines information about place where machines are stored.
type Location struct {
	// Datacenter - name of building where machine lies
	//+optional
	Datacenter string `json:"datacenter,omitempty"`
	// DataHall - name of room in Datacenter where machine lies
	//+optional
	DataHall string `json:"data_hall,omitempty"`
	// Shelf - defines place for server in Datacenter (an alternative name of Rack)
	//+optional
	Shelf string `json:"shelf,omitempty"`
	// Slot - defines switch location in rack (an alternative name for Row)
	//+optional
	Slot string `json:"slot,omitempty"`
	// HU - is a unit of measure defined 44.45 mm
	//+optional
	HU string `json:"hu,omitempty"`
	// Row - switch location in rack
	//+optional
	Row int16 `json:"row,omitempty"`
	// Rack - is a place for server in DataCenter
	//+optional
	Rack int16 `json:"rack,omitempty"`
}

// Identity - defines hardware information about machine.
type Identity struct {
	// SKU - stock keeping unit. The label allows vendors automatically track the movement of inventory
	//+optional
	SKU string `json:"sku,omitempty"`
	// SerialNumber - unique machine number
	//+optional
	SerialNumber string `json:"serial_number,omitempty"`
	//+optional
	Asset    string     `json:"asset,omitempty"`
	Internal []Internal `json:"internal,omitempty"`
}

type Internal struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Action struct {
	// PowerState - defines desired machine power state
	//+optional
	//+kubebuilder:validation:Pattern=`^(?:On|Reset|ResetImmediate|Off|OffImmediate)$`
	PowerState string `json:"power_state,omitempty"`
}

// Interface - defines information about machine interfaces.
type Interface struct {
	// Name - machine interface name
	//+optional
	Name string `json:"name,omitempty"`
	// SwitchUUID - defines unique switch identification
	//+optional
	SwitchUUID string `json:"switch_uuid,omitempty"`
	// IPv4 - defines machine IPv4 address
	//+optional
	//+kubebuilder:validation:Optional
	IPv4 string `json:"ipv4,omitempty"`
	// IPv6 - defines machine IPv6 address
	//+optional
	//+kubebuilder:validation:Optional
	IPv6 string `json:"ipv6,omitempty"`
	// LLDPSystemName - defines switch name obtained from Link Layer Discovery Protocol -
	// layer 2 neighbor discovery protocol
	//+optional
	LLDPSystemName string `json:"lldp_system_name,omitempty"`
	// LLDPChassisID - defines switch ID for chassis obtained from Link Layer Discovery Protocol
	//+optional
	LLDPChassisID string `json:"lldp_chassis_id,omitempty"`
	// LLDPPortID - defines switch port ID obtained from Link Layer Discovery Protocol
	//+optional
	LLDPPortID string `json:"lldp_port_id,omitempty"`
	// LLDPPortDescription - defines switch definition obtained from Link Layer Discovery Protocol
	//+optional
	LLDPPortDescription string `json:"lldp_port_description,omitempty"`
	// Lane - defines number of lines per interface
	//+optional
	Lane uint8 `json:"lane,omitempty"`
	// Moved  - defines if interface was reconnected to another switch or not
	//+optional
	Moved bool `json:"moved,omitempty"`
	// Unknown - defines information availability about interface
	//+optional
	Unknown bool `json:"unknown,omitempty"`
}

// MachineStatus - defines machine aggregated info.
type MachineStatus struct {
	// Interfaces - defines machine interfaces info
	//+optional
	Interfaces []Interface `json:"interfaces,omitempty"`
	// Reboot - defines machine reboot status
	//+optional
	Reboot string `json:"reboot,omitempty"`
	// Health - defines machine condition.
	// "healthy" if both OOB and Inventory are presented and "unhealthy" if one of them isn't
	//+optional
	Health string `json:"health,omitempty"`
	// Network - defines machine network status
	//+optional
	Network `json:"network,omitempty"`
	// Orphaned - defines machine condition whether OOB or Inventory is missing or not
	//+optional
	Orphaned bool `json:"orphaned,omitempty"`
	// Inventory - defines status, if Inventory is presented or not
	Inventory bool `json:"inventory,omitempty"`
	// OOB define status, OOB is presented or not
	OOB bool `json:"oob,omitempty"`
}

// Network - defines machine network status.
type Network struct {
	// Redundancy - defines machine redundancy status.
	// Available values: "Single", "High Availability" or "None"
	//+optional
	Redundancy string `json:"redundancy,omitempty"`
	// Ports - defines number of machine ports
	//+optional
	Ports int `json:"ports,omitempty"`
	// UnknownPorts - defines number of machine interface without info
	//+optional
	UnknownPorts int `json:"unknown_ports,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Healthy",type=string,JSONPath=`.status.health`
//+kubebuilder:printcolumn:name="Inventory",type=boolean,JSONPath=`.status.inventory`
//+kubebuilder:printcolumn:name="OOB",type=boolean,JSONPath=`.status.oob`
//+kubebuilder:printcolumn:name="Redundancy",type=string,JSONPath=`.status.network.redundancy`
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Machine - is the data structure for a Machine resource.
// It contains an aggregated information from Inventory and OOB resources.
type Machine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineSpec   `json:"spec,omitempty"`
	Status MachineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MachineList - contains a list of Machine.
type MachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Machine `json:"items"`
}
