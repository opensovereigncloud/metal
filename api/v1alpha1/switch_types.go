/*
Copyright 2021.

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

//SwitchSpec defines the desired state of Switch
//+kubebuilder:object:generate=true
type SwitchSpec struct {
	//Hostname
	//+kubebuilder:validation:Required
	Hostname string `json:"hostname"`
	//Location referring to the switch location
	//+kubebuilder:validation:Optional
	Location *LocationSpec `json:"location,omitempty"`
	//Ports referring to network interfaces total count
	//+kubebuilder:validation:Required
	Ports uint64 `json:"ports"`
	//SwitchPorts referring to non-management network interfaces count
	//+kubebuilder:validation:Required
	SwitchPorts uint64 `json:"switchPorts"`
	//SwitchDistro referring to switch OS information
	//+kubebuilder:validation:Optional
	SwitchDistro *SwitchDistroSpec `json:"switchDistro,omitempty"`
	//SwitchChassis referring to switch hardware information
	//+kubebuilder:validation:Required
	SwitchChassis *SwitchChassisSpec `json:"switchChassis"`
	//SouthSubnet referring to south IPv4 subnet
	//+kubebuilder:validation:Optional
	SouthSubnetV4 string `json:"southSubnetV4,omitempty"`
	//SouthSubnet referring to south IPv6 subnet
	//+kubebuilder:validation:Optional
	SouthSubnetV6 string `json:"southSubnetV6,omitempty"`
	//Role referring to switch's role: leaf or spine
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=Leaf;Spine;Undefined
	Role string `json:"role,omitempty"`
	// ConnectionLevel refers the level of the connection
	//+kubebuilder:validation:Optional
	//+kubebuilder:default=255
	ConnectionLevel uint8 `json:"connectionLevel"`
	//Interfaces referring to details about network interfaces
	//+kubebuilder:validation:Optional
	Interfaces []*InterfaceSpec `json:"interfaces,omitempty"`
	//ScanPorts flag determining whether scanning of ports is requested
	//+kubebuilder:validation:Optional
	//+kubebuilder:default=true
	ScanPorts bool `json:"scanPorts,omitempty"`
}

//LocationSpec defines location details
//+kubebuilder:object:generate=true
type LocationSpec struct {
	//Room referring to room name
	//+kubebuilder:validation:Optional
	Room string `json:"room,omitempty"`
	//Row referring to row number
	//+kubebuilder:validation:Optional
	Row int16 `json:"row,omitempty"`
	//Rack referring to rack number
	//+kubebuilder:validation:Optional
	Rack int16 `json:"rack,omitempty"`
	//HU referring to height in units
	//+kubebuilder:validation:Optional
	HU int16 `json:"hu,omitempty"`
}

//SwitchDistroSpec defines switch OS details
//+kubebuilder:object:generate=true
type SwitchDistroSpec struct {
	//OS referring to switch operating system
	//+kubebuilder:validation:Optional
	OS string `json:"os,omitempty"`
	//Version referring to switch OS version
	//+kubebuilder:validation:Optional
	Version string `json:"version,omitempty"`
	//ASIC
	//+kubebuilder:validation:Optional
	ASIC string `json:"asic,omitempty"`
}

//SwitchChassisSpec defines switch chassis details
//+kubebuilder:object:generate=true
type SwitchChassisSpec struct {
	//Manufacturer referring to switch chassis manufacturer
	//+kubebuilder:validation:Optional
	Manufacturer string `json:"manufacturer,omitempty"`
	//SKU
	//+kubebuilder:validation:Optional
	SKU string `json:"sku,omitempty"`
	//Serial referring to switch chassis serial number
	//+kubebuilder:validation:Optional
	Serial string `json:"serial,omitempty"`
	//ChassisID referring to switch chassis ID advertising via LLDP
	//+kubebuilder:validation:Optional
	ChassisID string `json:"chassisId,omitempty"`
}

//InterfaceSpec defines switch's network interface details
//+kubebuilder:object:generate=true
type InterfaceSpec struct {
	//Lanes referring to how many lanes are used by the interface based on it's speed
	//+kubebuilder:validation:Optional
	Lanes uint8 `json:"lanes,omitempty"`
	//FEC referring to error correction method
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=None;BaseR;RS
	FEC string `json:"fec,omitempty"`
	//Neighbour referring to neighbour type
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=Machine;Switch
	Neighbour string `json:"neighbour,omitempty"`
	//Name referring to interface's name
	//+kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	//MACAddress referring to interface's MAC address
	//+kubebuilder:validation:Optional
	MACAddress string `json:"macAddress,omitempty"`
	//IPv4 referring to interface's IPv4 address
	//+kubebuilder:validation:Optional
	IPv4 string `json:"ipv4,omitempty"`
	//IPv6 referring to interface's IPv6 address
	//+kubebuilder:validation:Optional
	IPv6 string `json:"ipv6,omitempty"`
	//LLDPSystemName
	//+kubebuilder:validation:Optional
	LLDPSystemName string `json:"lldpSystemName,omitempty"`
	//LLDPChassisID
	//+kubebuilder:validation:Optional
	LLDPChassisID string `json:"lldpChassisId,omitempty"`
	//LLDPPortID
	//+kubebuilder:validation:Optional
	LLDPPortID string `json:"lldpPortId,omitempty"`
	//LLDPPortDescription
	//+kubebuilder:validation:Optional
	LLDPPortDescription string `json:"lldpPortDescription,omitempty"`
}

// SwitchStatus defines the observed state of Switch
type SwitchStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=sw
//+kubebuilder:printcolumn:name="Hostname",type=string,JSONPath=`.spec.hostname`,description="Switch's hostname"
//+kubebuilder:printcolumn:name="Role",type=string,JSONPath=`.spec.role`,description="switch's role"
//+kubebuilder:printcolumn:name="OS",type=string,JSONPath=`.spec.switchDistro.os`,description="OS running on switch"
//+kubebuilder:printcolumn:name="SwitchPorts",type=integer,JSONPath=`.spec.switchPorts`,description="Total amount of non-management network interfaces"
//+kubebuilder:printcolumn:name="ConnectionLevel",type=integer,JSONPath=`.spec.connectionLevel`,description="Vertical level of switch connection"
//+kubebuilder:printcolumn:name="SouthSubnet",type=string,JSONPath=`.spec.southSubnet`,description="South subnet"
//+kubebuilder:printcolumn:name="ScanPorts",type=boolean,JSONPath=`.spec.scanPorts`,description="Request for scan ports"

// Switch is the Schema for the switches API
type Switch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwitchSpec   `json:"spec,omitempty"`
	Status SwitchStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SwitchList contains a list of Switch
type SwitchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Switch `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Switch{}, &SwitchList{})
}
