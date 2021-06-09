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
	"math"
	"net"

	subnetv1alpha1 "github.com/onmetal/k8s-subnet/api/v1alpha1"
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
	SouthSubnetV4 *SwitchSubnetSpec `json:"southSubnetV4,omitempty"`
	//SouthSubnet referring to south IPv6 subnet
	//+kubebuilder:validation:Optional
	SouthSubnetV6 *SwitchSubnetSpec `json:"southSubnetV6,omitempty"`
	//Interfaces referring to details about network interfaces
	//+kubebuilder:validation:Optional
	Interfaces []*InterfaceSpec `json:"interfaces,omitempty"`
	//ScanPorts flag determining whether scanning of ports is requested
	//+kubebuilder:validation:Optional
	//+kubebuilder:default=true
	ScanPorts bool `json:"scanPorts,omitempty"`
	//State referring to current switch state
	//kubebuilder:validation:Optional
	State *SwitchStateSpec `json:"state"`
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

// SwitchSubnetSpec defines switch subnet details
//+kubebuilder:object:generate=true
type SwitchSubnetSpec struct {
	// ParentSubnet referring to the subnet resource namespaced name where CIDR was booked
	//+kubebuilder:validation:Optional
	ParentSubnet *ParentSubnetSpec `json:"parentSubnet"`
	// CIDR referring to the assigned subnet
	//+kubebuilder:validation:Optional
	CIDR string `json:"cidr"`
}

// ParentSubnetSpec defines switch subnet name and namespace
//+kubebuilder:object:generate=true
type ParentSubnetSpec struct {
	// Name referring to the subnet resource name where CIDR was booked
	//+kubebuilder:validation:Optional
	Name string `json:"name"`
	// Namespace referring to the subnet resource name where CIDR was booked
	//+kubebuilder:validation:Optional
	Namespace string `json:"namespace"`
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

// SwitchStateSpec defines current connection state of the Switch
type SwitchStateSpec struct {
	//Role referring to switch's role: leaf or spine
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=Leaf;Spine;Undefined
	Role string `json:"role,omitempty"`
	// ConnectionLevel refers the level of the connection
	//+kubebuilder:validation:Optional
	//+kubebuilder:default=255
	ConnectionLevel uint8 `json:"connectionLevel"`
	// NorthSwitches refers to up-level switch
	//+kubebuilder:validation:Optional
	NorthConnections *NorthConnectionsSpec `json:"northSwitches,omitempty"`
	// SouthSwitches refers to down-level switch
	//+kubebuilder:validation:Optional
	SouthConnections *SouthConnectionsSpec `json:"southSwitches,omitempty"`
}

// SwitchStatus defines the observed state of Switch
type SwitchStatus struct{}

// NorthConnectionsSpec defines upstream switches count and properties
//+kubebuilder:object:generate=true
type NorthConnectionsSpec struct {
	// Count refers to upstream switches count
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:default=0
	Count int `json:"count"`
	// Switches refers to connected upstream switches
	//+kubebuilder:validation:Optional
	Connections []NeighbourSpec `json:"switches"`
}

// SouthConnectionsSpec defines downstream switches count and properties
//+kubebuilder:object:generate=true
type SouthConnectionsSpec struct {
	// Count refers to upstream switches count
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:default=0
	Count int `json:"count"`
	// Switches refers to connected upstream switches
	//+kubebuilder:validation:Optional
	Connections []NeighbourSpec `json:"switches"`
}

// NeighbourSpec defines switch connected to another switch
//+kubebuilder:object:generate=true
type NeighbourSpec struct {
	// Name refers to switch's name
	//+kubebuilder:validation:Optional
	Name string `json:"name"`
	// Namespace refers to switch's namespace
	//+kubebuilder:validation:Optional
	Namespace string `json:"namespace"`
	// ChassisID refers to switch's chassis id
	//+kubebuilder:validation:Required
	ChassisID string `json:"chassisId"`
	//Type referring to neighbour type
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=Machine;Switch
	Type string `json:"type,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=sw
//+kubebuilder:printcolumn:name="Hostname",type=string,JSONPath=`.spec.hostname`,description="Switch's hostname"
//+kubebuilder:printcolumn:name="Role",type=string,JSONPath=`.spec.state.role`,description="switch's role"
//+kubebuilder:printcolumn:name="OS",type=string,JSONPath=`.spec.switchDistro.os`,description="OS running on switch"
//+kubebuilder:printcolumn:name="SwitchPorts",type=integer,JSONPath=`.spec.switchPorts`,description="Total amount of non-management network interfaces"
//+kubebuilder:printcolumn:name="ConnectionLevel",type=integer,JSONPath=`.spec.state.connectionLevel`,description="Vertical level of switch connection"
//+kubebuilder:printcolumn:name="SouthSubnetV4",type=string,JSONPath=`.spec.southSubnetV4.cidr`,description="South IPv4 subnet"
//+kubebuilder:printcolumn:name="SouthSubnetV6",type=string,JSONPath=`.spec.southSubnetV6.cidr`,description="South IPv6 subnet"
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

func (sw *Switch) GetNorthSwitchConnection(swList []Switch) []Switch {
	result := make([]Switch, 0)
	for _, obj := range swList {
		if sw.checkSwitchInSouthConnections(obj) {
			result = append(result, obj)
		}
	}
	return result
}

func (sw *Switch) checkSwitchInSouthConnections(obj Switch) bool {
	for _, conn := range obj.Spec.State.SouthConnections.Connections {
		if sw.Spec.SwitchChassis.ChassisID == conn.ChassisID {
			return true
		}
	}
	return false
}

func (sw *Switch) CheckMachinesConnected() bool {
	for _, iface := range sw.Spec.Interfaces {
		if iface.Neighbour == CMachineType {
			return true
		}
	}
	return false
}

func (sw *Switch) CheckSouthNeighboursDataUpdateNeeded() bool {
	for _, item := range sw.Spec.State.SouthConnections.Connections {
		if item.Name == "" || item.Namespace == "" {
			return true
		}
	}
	return false
}

func (sw *Switch) CheckNorthNeighboursDataUpdateNeeded() bool {
	for _, item := range sw.Spec.State.NorthConnections.Connections {
		if item.Name == "" || item.Namespace == "" {
			return true
		}
	}
	return false
}

func (sw *Switch) GetNeededMask(addrType subnetv1alpha1.SubnetAddressType, addressesCount float64) net.IPMask {
	bits := uint8(0)
	if addrType == subnetv1alpha1.CIPv4SubnetType {
		addressesCount = float64(sw.GetAddressNeededCount(addrType))
		bits = 32
	}
	if addrType == subnetv1alpha1.CIPv6SubnetType {
		addressesCount = float64(sw.GetAddressNeededCount(addrType))
		bits = 128
	}
	pow := 2.0
	for math.Pow(2, pow) < addressesCount {
		pow++
	}
	maskLength := bits - uint8(pow)
	mask := (0xFFFFFFFF << (bits - maskLength)) & 0xFFFFFFFF
	netMask := make([]byte, 0, 4)
	for i := 1; i <= 4; i++ {
		tmp := byte(mask >> (bits - 8) & 0xFF)
		netMask = append(netMask, tmp)
		bits -= 8
	}
	return netMask
}

func (sw *Switch) GetAddressNeededCount(addrType subnetv1alpha1.SubnetAddressType) int64 {
	if addrType == subnetv1alpha1.CIPv4SubnetType {
		return int64(sw.Spec.SwitchPorts * CIPv4AddressesPerPort)
	} else {
		return int64(sw.Spec.SwitchPorts * CIPv6AddressesPerPort)
	}
}
