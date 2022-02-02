/*
Copyright 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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
	"fmt"
	"net"
	"reflect"
	"sort"
	"strconv"
	"strings"

	gocidr "github.com/apparentlymart/go-cidr/cidr"
	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type SwitchState string

const (
	CSwitchStateReady      SwitchState = "ready"
	CSwitchStateInProgress SwitchState = "in progress"
	CSwitchStateInitial    SwitchState = "initial"
)

type SwitchRole string

const (
	CSwitchRoleLeaf  SwitchRole = "leaf"
	CSwitchRoleSpine SwitchRole = "spine"
)

type SwitchConfState string

const (
	CSwitchConfInitial    SwitchConfState = "initial"
	CSwitchConfApplied    SwitchConfState = "applied"
	CSwitchConfPending    SwitchConfState = "pending"
	CSwitchConfInProgress SwitchConfState = "in progress"
)

type FECType string

const (
	CFECNone FECType = "none"
	CFECRS   FECType = "rs"
	CFECFC   FECType = "fc"
)

type NICState string

const (
	CNICUp   NICState = "up"
	CNICDown NICState = "down"
)

type NICDirection string

const (
	CDirectionSouth NICDirection = "south"
	CDirectionNorth NICDirection = "north"
)

const (
	CPeerTypeMachine   string = "machine"
	CPeerTypeSwitch    string = "switch"
	CPeerTypeRouter    string = "router"
	CPeerTypeUndefined string = "undefined"
)

type ConfManagerType string

const (
	CConfManagerTLocal  ConfManagerType = "local"
	CConfManagerTRemote ConfManagerType = "remote"
)

type ConfManagerState string

const (
	CConfManagerSActive ConfManagerState = "active"
	CConfManagerSFailed ConfManagerState = "failed"

	CSwitchPortMTU    = 9100
	CEmptyString      = ""
	CSwitchPortPrefix = "Ethernet"

	CNamespace = "onmetal"

	CIPv4AddressesPerLane    = uint8(4)
	CIPv6AddressesPerLane    = uint8(2)
	CIPv4InterfaceSubnetMask = 30
	CIPv6InterfaceSubnetMask = 127

	CSonicSwitchOs     = "SONiC"
	CStationCapability = "Station"
	CRouterCapability  = "Router"
	CBridgeCapability  = "Bridge"
	CNDPReachable      = "Reachable"

	CLabelPrefix    = "switch.onmetal.de/"
	CLabelChassisId = "chassisId"
	CLabelName      = "name"
	CLabelInterface = "interface"
	CLabelRelation  = "relation"
)

var LabelChassisId = CLabelPrefix + CLabelChassisId
var LabelSwitchName = CLabelPrefix + CLabelName
var LabelInterfaceName = CLabelPrefix + CLabelInterface
var LabelResourceRelation = CLabelPrefix + CLabelRelation

func MacToLabel(mac string) string {
	return strings.ReplaceAll(mac, ":", "-")
}

type ConnectionsMap map[uint8]*SwitchList

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SwitchSpec defines the desired state of Switch
//+kubebuilder:object:generate=true
type SwitchSpec struct {
	//Hostname refers to switch hostname
	//+kubebuilder:validation:Required
	Hostname string `json:"hostname"`
	//Chassis refers to baremetal box info
	//+kubebuilder:validation:Required
	Chassis *ChassisSpec `json:"chassis"`
	//SoftwarePlatform refers to software info
	//+kubebuilder:validation:Required
	SoftwarePlatform *SoftwarePlatformSpec `json:"softwarePlatform"`
	//Location refers to the switch's location
	//+kubebuilder:validation:Optional
	Location *LocationSpec `json:"location,omitempty"`
}

// ChassisSpec defines switch's chassis info
//+kubebuilder:object:generate=true
type ChassisSpec struct {
	//ChassisID refers to the chassis identificator - either MAC-address or system uuid
	//+kubebuilder:validation:Required
	//validation pattern
	ChassisID string `json:"chassisId"`
	//Manufactirer refers to the switch's manufacturer
	//+kubebuilder:validation:Optional
	Manufacturer string `json:"manufacturer,omitempty"`
	//SerialNumber refers to the switch's serial number
	//+kubebuilder:validation:Optional
	SerialNumber string `json:"serialNumber,omitempty"`
	//SKU refers to the switch's stock keeping unit
	//+kubebuilder:validation:Optional
	SKU string `json:"sku,omitempty"`
}

// SoftwarePlatformSpec defines switch's software base
//+kubebuilder:object:generate=true
type SoftwarePlatformSpec struct {
	//ONIE refers to whether open network installation environment is used
	//+kubebuilder:validation:Optional
	ONIE bool `json:"onie,omitempty"`
	//OperatingSystem refers to switch's operating system
	//+kubebuilder:validation:Required
	OperatingSystem string `json:"operatingSystem"`
	//Version refers to the operating system version
	//+kubebuilder:validation:Optional
	Version string `json:"version,omitempty"`
	//ASIC refers to the switch's ASIC manufacturer
	//+kubebuilder:validation:Optional
	ASIC string `json:"asic,omitempty"`
}

// LocationSpec defines switch's location
//+kubebuilder:object:generate=true
type LocationSpec struct {
	//Room refers to room name
	//+kubebuilder:validation:Optional
	Room string `json:"room,omitempty"`
	//Row refers to row number
	//+kubebuilder:validation:Optional
	Row int16 `json:"row,omitempty"`
	//Rack refers to rack number
	//+kubebuilder:validation:Optional
	Rack int16 `json:"rack,omitempty"`
	//HU refers to height in units
	//+kubebuilder:validation:Optional
	HU int16 `json:"hu,omitempty"`
}

// SwitchStatus defines the observed state of Switch
type SwitchStatus struct {
	//TotalPorts refers to total number of ports
	//+kubebuilder:validation:Required
	TotalPorts uint16 `json:"totalPorts"`
	//SwitchPorts refers to the number of ports excluding management interfaces, loopback etc.
	//+kubebuilder:validation:Required
	SwitchPorts uint16 `json:"switchPorts"`
	//Role refers to switch's role
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Enum=spine;leaf
	Role SwitchRole `json:"role"`
	//ConnectionLevel refers to switch's current position in connection hierarchy
	//+kubebuilder:validation:Required
	ConnectionLevel uint8 `json:"connectionLevel"`
	//Interfaces refers to switch's interfaces configuration
	//+kubebuilder:validation:Required
	Interfaces map[string]*InterfaceSpec `json:"interfaces"`
	//SubnetV4 refers to the switch's south IPv4 subnet
	//+kubebuilder:validation:Optional
	SubnetV4 *SubnetSpec `json:"subnetV4,omitempty"`
	//SubnetV6 refers to the switch's south IPv6 subnet
	//+kubebuilder:validation:Optional
	SubnetV6 *SubnetSpec `json:"subnetV6,omitempty"`
	//LoopbackV4 refers to the switch's loopback IPv4 address
	//+kubebuilder:validation:Optional
	LoopbackV4 *IPAddressSpec `json:"loopbackV4,omitempty"`
	//LoopbackV6 refers to the switch's loopback IPv6 address
	//+kubebuilder:validation:Optional
	LoopbackV6 *IPAddressSpec `json:"loopbackV6,omitempty"`
	//Configuration refers to how switch's configuration manager is defined
	//+kubebuilder:validation:Required
	Configuration *ConfigurationSpec `json:"configuration"`
	//State refers to current switch's processing state
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Enum=initial;in progress;ready
	State SwitchState `json:"state"`
}

// InterfaceSpec defines the state of switch's interface
//+kubebuilder:object:generate=true
type InterfaceSpec struct {
	//MACAddress refers to the interface's hardware address
	//+kubebuilder:validation:Required
	//validation pattern
	MACAddress string `json:"macAddress"`
	//FEC refers to the current interface's forward error correction type
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Enum=none;rs;fc
	FEC FECType `json:"fec"`
	//MTU refers to the current value of interface's MTU
	//+kubebuilder:validation:Required
	MTU uint16 `json:"mtu"`
	//Speed refers to interface's speed
	//+kubebuilder:validation:Required
	Speed uint32 `json:"speed"`
	//Lanes refers to the number of lanes used by interface
	//+kubebuilder:validation:Required
	Lanes uint8 `json:"lanes"`
	//State refers to the current interface's operational state
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Enum=up;down
	State NICState `json:"state"`
	//IPv4 refers to the interface's IPv4 address
	//+kubebuilder:validation:Optional
	IPv4 *IPAddressSpec `json:"ipV4,omitempty"`
	//IPv6 refers to the interface's IPv6 address
	//+kubebuilder:validation:Optional
	IPv6 *IPAddressSpec `json:"ipV6,omitempty"`
	//Direction refers to the interface's connection 'direction'
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Enum=north;south
	Direction NICDirection `json:"direction"`
	//Peer refers to the info about device connected to current switch port
	//+kubebuilder:validation:Optional
	Peer *PeerSpec `json:"peer,omitempty"`
}

// PeerSpec defines peer info
//+kubebuilder:object:generate=true
type PeerSpec struct {
	//ChassisID refers to the chassis identificator - either MAC-address or system uuid
	//+kubebuilder:validation:Optional
	//validation pattern
	ChassisID string `json:"chassisId,omitempty"`
	//SystemName refers to the advertised peer's name
	//+kubebuilder:validation:Optional
	SystemName string `json:"systemName,omitempty"`
	//PortID refers to the advertised peer's port ID
	//+kubebuilder:validation:Optional
	PortID string `json:"portId,omitempty"`
	//PortDescription refers to the advertised peer's port description
	//+kubebuilder:validation:Optional
	PortDescription string `json:"portDescription,omitempty"`
	//Type refers to the peer type
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=machine;switch;router;undefined
	Type string `json:"type,omitempty"`
	//ResourceReference refers to the related resource definition
	//+kubebuilder:validation:Optional
	ResourceReference *ResourceReferenceSpec `json:"resourceReference,omitempty"`
}

// SubnetSpec defines switch's subnet info
//+kubebuilder:object:generate=true
type SubnetSpec struct {
	//CIDR refers to subnet CIDR
	//+kubebuilder:validation:Optional
	//validation pattern
	CIDR string `json:"cidr,omitempty"`
	//Region refers to switch's region
	//+kubebuilder:validation:Optional
	Region *RegionSpec `json:"region,omitempty"`
	//ResourceReference refers to the related resource definition
	//+kubebuilder:validation:Optional
	ResourceReference *ResourceReferenceSpec `json:"resourceReference,omitempty"`
}

// IPAddressSpec defines interface's ip address info
//+kubebuilder:object:generate=true
type IPAddressSpec struct {
	//Address refers to the ip address value
	//+kubebuilder:validation:Optional
	//validation pattern
	Address string `json:"address,omitempty"`
	//ResourceReference refers to the related resource definition
	//+kubebuilder:validation:Optional
	ResourceReference *ResourceReferenceSpec `json:"resourceReference,omitempty"`
}

// ResourceReferenceSpec defines related resource info
//+kubebuilder:object:generate=true
type ResourceReferenceSpec struct {
	//APIVersion refers to the resource API version
	//+kubebuilder:validation:Optional
	APIVersion string `json:"apiVersion,omitempty"`
	//Kind refers to the resource kind
	//+kubebuilder:validation:Optional
	Kind string `json:"kind,omitempty"`
	//Name refers to the resource name
	//+kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	//Namespace refers to the resource namespace
	//+kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
}

// ConfigurationSpec defines switch's computed configuration
//+kubebuilder:object:generate=true
type ConfigurationSpec struct {
	//Managed refers to whether switch configuration is managed or not
	//+kubebuilder:validation:Required
	Managed bool `json:"managed"`
	//State refers to current switch's configuration processing state
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Enum=initial;applied;in progress;pending
	State SwitchConfState `json:"state"`
	//Type refers to configuration manager type
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=local;remote
	ManagerType ConfManagerType `json:"managerType,omitempty"`
	//State refers to configuration manager state
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=active;failed
	ManagerState ConfManagerState `json:"managerState,omitempty"`
	//LastCheck refers to the last timestamp when configuration was applied
	//+kubebuilder:validation:Optional
	LastCheck string `json:"lastCheck,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=sw
//+kubebuilder:printcolumn:name="Hostname",type=string,JSONPath=`.spec.hostname`,description="Switch's hostname"
//+kubebuilder:printcolumn:name="OS",type=string,JSONPath=`.spec.softwarePlatform.operatingSystem`,description="OS running on switch"
//+kubebuilder:printcolumn:name="Ports",type=integer,JSONPath=`.status.switchPorts`,description="Total amount of non-management network interfaces"
//+kubebuilder:printcolumn:name="Role",type=string,JSONPath=`.status.role`,description="switch's role"
//+kubebuilder:printcolumn:name="Conn Level",type=integer,JSONPath=`.status.connectionLevel`,description="Vertical level of switch connection"
//+kubebuilder:printcolumn:name="Subnet V4",type=string,priority=1,JSONPath=`.status.subnetV4.cidr`,description="South IPv4 subnet"
//+kubebuilder:printcolumn:name="Subnet V6",type=string,priority=1,JSONPath=`.status.subnetV6.cidr`,description="South IPv6 subnet"
//+kubebuilder:printcolumn:name="Switch State",type=string,JSONPath=`.status.state`,description="Switch state"
//+kubebuilder:printcolumn:name="Conf State",type=string,priority=1,JSONPath=`.status.configuration.state`,description="Switch configuration processing state"
//+kubebuilder:printcolumn:name="Manager Type",type=string,priority=1,JSONPath=`.status.configuration.managerType`,description="Switch manager type"
//+kubebuilder:printcolumn:name="Manager State",type=string,priority=1,JSONPath=`.status.configuration.managerState`,description="Switch manager state"

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

// NamespacedName returns referenced resource name and namespaced converted
// to native types.NamespacedName format
func (in *ResourceReferenceSpec) NamespacedName() types.NamespacedName {
	return types.NamespacedName{Namespace: in.Namespace, Name: in.Name}
}

// NamespacedName returns switch's name and namespaced converted to native
// types.NamespacedName format
func (in *Switch) NamespacedName() types.NamespacedName {
	return types.NamespacedName{Namespace: in.Namespace, Name: in.Name}
}

// SetSwitchState sets certain state for switch
func (in *Switch) SetSwitchState(state SwitchState) {
	in.Status.State = state
}

// SetConfState sets certain state for switch config
func (in *Switch) SetConfState(state SwitchConfState) {
	in.Status.Configuration.State = state
}

// SwitchFromInventory builds switch resource from inventory resource data
func (in *Switch) SwitchFromInventory(src *inventoriesv1alpha1.Inventory) {
	in.ObjectMeta = metav1.ObjectMeta{
		Name:      src.Name,
		Namespace: CNamespace,
	}
	in.Spec = SwitchSpec{
		Hostname: src.Spec.Host.Name,
		Chassis: &ChassisSpec{
			ChassisID: func(nics []inventoriesv1alpha1.NICSpec) string {
				var chassisID string
				for _, nic := range nics {
					if nic.Name == "eth0" {
						chassisID = nic.MACAddress
					}
				}
				return chassisID
			}(src.Spec.NICs),
			Manufacturer: src.Spec.System.Manufacturer,
			SerialNumber: src.Spec.System.SerialNumber,
			SKU:          src.Spec.System.ProductSKU,
		},
		SoftwarePlatform: &SoftwarePlatformSpec{
			ONIE:            false,
			OperatingSystem: CSonicSwitchOs,
			Version:         src.Spec.Distro.CommitId,
			ASIC:            src.Spec.Distro.AsicType,
		},
		Location: &LocationSpec{},
	}
}

// FillInitialStatus fills switch's status with initial values
func (in *Switch) FillInitialStatus(src *inventoriesv1alpha1.Inventory, switches *SwitchList) {
	in.Status = SwitchStatus{
		TotalPorts:      uint16(len(src.Spec.NICs)),
		SwitchPorts:     0,
		Role:            CSwitchRoleSpine,
		ConnectionLevel: 255,
		Interfaces:      InterfacesFromInventory(src.Spec.NICs, switches),
		SubnetV4:        &SubnetSpec{},
		SubnetV6:        &SubnetSpec{},
		LoopbackV4:      &IPAddressSpec{},
		LoopbackV6:      &IPAddressSpec{},
		Configuration: &ConfigurationSpec{
			Managed: false,
			State:   CSwitchConfInitial,
		},
		State: CSwitchStateInitial,
	}
	in.Status.SwitchPorts = uint16(len(in.Status.Interfaces))
}

// InterfacesFromInventory builds interfaces map based on inventory resource data
func InterfacesFromInventory(nics []inventoriesv1alpha1.NICSpec, switches *SwitchList) map[string]*InterfaceSpec {
	interfaces := make(map[string]*InterfaceSpec)
	for _, nic := range nics {
		if !(strings.HasPrefix(nic.Name, CSwitchPortPrefix)) {
			continue
		}
		iface := &InterfaceSpec{
			MACAddress: nic.MACAddress,
			FEC: func(nicFEC string) FECType {
				if nicFEC == CEmptyString {
					nicFEC = "none"
				}
				return FECType(nicFEC)
			}(nic.ActiveFEC),
			MTU:       CSwitchPortMTU,
			Speed:     nic.Speed,
			Lanes:     nic.Lanes,
			State:     CNICUp, //TODO: extend inventory CR with NIC state field
			IPv4:      &IPAddressSpec{},
			IPv6:      &IPAddressSpec{},
			Direction: CDirectionSouth,
			Peer:      &PeerSpec{Type: CPeerTypeUndefined},
		}
		for _, lldpData := range nic.LLDPs {
			var lldpEmpty inventoriesv1alpha1.LLDPSpec
			if reflect.DeepEqual(lldpData, lldpEmpty) {
				continue
			}
			iface.Peer.ChassisID = lldpData.ChassisID
			iface.Peer.SystemName = lldpData.SystemName
			iface.Peer.PortID = lldpData.PortID
			iface.Peer.PortDescription = lldpData.PortDescription
			iface.Peer.Type = func(caps []inventoriesv1alpha1.LLDPCapabilities) string {
				if len(caps) == 0 {
					return CPeerTypeMachine
				}
				for _, capacity := range caps {
					if capacity == CStationCapability {
						return CPeerTypeMachine
					}
				}
				return CPeerTypeSwitch
			}(lldpData.Capabilities)
			iface.Peer.ResourceReference = &ResourceReferenceSpec{}
		}

		// FIXME: if there is no lldp data on nics connected to another switch we have a problem:
		//  no way to determine what interfaces were use for switches interconnection

		// for _, ndpData := range nic.NDPs {
		// 	if iface.Peer.ChassisID != CEmptyString {
		// 		break
		// 	}
		// 	if ndpData.State != CNDPReachable {
		// 		continue
		// 	}
		// 	iface.Peer.ChassisID = ndpData.MACAddress
		// 	iface.Peer.Type = func(switches *SwitchList) string {
		// 		for _, sw := range switches.Items {
		// 			if sw.Spec.Chassis.ChassisID == ndpData.MACAddress {
		// 				return CPeerTypeSwitch
		// 			}
		// 		}
		// 		return CPeerTypeMachine
		// 	}(switches)
		// 	iface.Peer.ResourceReference = &ResourceReferenceSpec{}
		// }
		interfaces[nic.Name] = iface
	}
	return interfaces
}

// InterfacesDataOk checks stored interfaces data, that can be changed, is
// equal to received from inventory
func (in *Switch) InterfacesDataOk(src *inventoriesv1alpha1.Inventory, switches *SwitchList) bool {
	receivedInterfaces := InterfacesFromInventory(src.Spec.NICs, switches)
	storedInterfaces := in.Status.Interfaces
	for nicName, nicData := range receivedInterfaces {
		storedData, ok := storedInterfaces[nicName]
		if !ok {
			return false
		}
		if nicData.MACAddress != storedData.MACAddress {
			return false
		}
		if nicData.FEC != storedData.FEC {
			return false
		}
		if nicData.Lanes != storedData.Lanes {
			return false
		}
		if nicData.State != storedData.State {
			return false
		}
		if nicData.Speed != storedData.Speed {
			return false
		}
		if nicData.Peer.ChassisID != storedData.Peer.ChassisID {
			return false
		}
		if nicData.Peer.Type != storedData.Peer.Type {
			return false
		}
	}
	for iface := range storedInterfaces {
		if _, ok := receivedInterfaces[iface]; !ok {
			return false
		}
	}
	return true
}

// StateEqualTo checks whether switch resource state is equal to state
// passed as argument
func (in *Switch) StateEqualTo(state string) bool {
	return string(in.Status.State) == state
}

// PeersDefined checks whether peers data is filled and match to the
// existing resources
func (in *Switch) PeersDefined(list *SwitchList) bool {
	for _, data := range in.Status.Interfaces {
		if data.Direction == CDirectionNorth && in.Status.ConnectionLevel == 0 {
			return false
		}
	}
	for _, item := range list.Items {
		for _, nicData := range item.Status.Interfaces {
			if nicData.Peer.ChassisID != in.Spec.Chassis.ChassisID {
				continue
			}
			nic := in.Status.Interfaces[nicData.Peer.PortDescription] // portDescription may be absent!
			if nic.Peer.Type != CPeerTypeSwitch {
				return false
			}
			if nic.Peer.ResourceReference.Name != item.Name {
				return false
			}
			if nic.Peer.ResourceReference.Namespace != item.Namespace {
				return false
			}
		}
	}
	return true
}

// FillPeerSwitches fills references info in interfaces data for
// interfaces connected to another switches
func (in *Switch) FillPeerSwitches(list *SwitchList) {
	for _, item := range list.Items {
		for _, nicData := range item.Status.Interfaces {
			if nicData.Peer.ChassisID != in.Spec.Chassis.ChassisID {
				continue
			}
			nic := in.Status.Interfaces[nicData.Peer.PortDescription]
			nic.Peer.Type = CPeerTypeSwitch
			nic.Peer.ResourceReference.APIVersion = item.APIVersion
			nic.Peer.ResourceReference.Kind = item.Kind
			nic.Peer.ResourceReference.Name = item.Name
			nic.Peer.ResourceReference.Namespace = item.Namespace
		}
	}
}

// RoleMatchPeers checks whether switch's role match stored peers info
func (in *Switch) RoleMatchPeers() bool {
	machinesInPeers := false
	for _, nicData := range in.Status.Interfaces {
		if nicData.Peer.Type == CPeerTypeMachine {
			machinesInPeers = true
		}
	}
	if machinesInPeers && in.Status.Role == CSwitchRoleSpine {
		return false
	}
	if !machinesInPeers && in.Status.Role == CSwitchRoleLeaf {
		return false
	}
	return true
}

// ConnectionLevelMatchPeers checks whether switch's interfaces' directions
// are defined correct and match connection levels of peers
func (in *Switch) ConnectionLevelMatchPeers(list *SwitchList) bool {
	for _, item := range list.Items {
		for _, nicData := range in.Status.Interfaces {
			if nicData.Peer.ChassisID != item.Spec.Chassis.ChassisID {
				continue
			}
			if nicData.Direction == CDirectionNorth && in.Status.ConnectionLevel != item.Status.ConnectionLevel+1 {
				return false
			}
			if nicData.Direction == CDirectionSouth && in.Status.ConnectionLevel != item.Status.ConnectionLevel-1 {
				return false
			}
		}
	}
	return true
}

// ComputeConnectionLevel calculates switch's connection level
// according to peers connection levels
func (in *Switch) ComputeConnectionLevel(list *SwitchList) {
	connectionsMap, keys := list.buildConnectionMap()
	if _, ok := connectionsMap[0]; !ok {
		return
	}
	if in.Status.ConnectionLevel == 0 {
		for _, nicData := range in.Status.Interfaces {
			nicData.Direction = CDirectionSouth
		}
		return
	}
	for _, connectionLevel := range keys {
		if connectionLevel == 255 {
			continue
		}
		if connectionLevel >= in.Status.ConnectionLevel {
			continue
		}
		switches := connectionsMap[connectionLevel]
		northPeers := in.getPeers(switches)
		if len(northPeers.Items) == 0 {
			continue
		}
		in.Status.ConnectionLevel = connectionLevel + 1
		in.fillNorthPeers(northPeers)
		in.setNICsDirections(list)
	}
}

// Determines what Switch resources are the known
// peers for current Switch.
func (in *Switch) getPeers(list *SwitchList) (result *SwitchList) {
	result = &SwitchList{Items: make([]Switch, 0)}
	for _, item := range list.Items {
		for _, data := range in.Status.Interfaces {
			if data.Peer.ChassisID == item.Spec.Chassis.ChassisID {
				result.Items = append(result.Items, item)
			}
		}
	}
	return
}

// fillNorthPeers fills resource reference info for switch's
// north peers
func (in *Switch) fillNorthPeers(list *SwitchList) {
	for _, item := range list.Items {
		for _, nicData := range in.Status.Interfaces {
			if nicData.Peer.ChassisID == item.Spec.Chassis.ChassisID {
				nicData.Peer.ResourceReference.APIVersion = item.APIVersion
				nicData.Peer.ResourceReference.Kind = item.Kind
				nicData.Peer.ResourceReference.Name = item.Name
				nicData.Peer.ResourceReference.Namespace = item.Namespace
			}
		}
	}
}

// setNICsDirections updates NICs' direction field according to
// the computed connection levels
func (in *Switch) setNICsDirections(list *SwitchList) {
	if in.Status.ConnectionLevel == 0 {
		for _, nicData := range in.Status.Interfaces {
			nicData.Direction = CDirectionSouth
		}
		return
	}
	for _, item := range list.Items {
		for _, nicData := range in.Status.Interfaces {
			peerFound := nicData.Peer.ChassisID == item.Spec.Chassis.ChassisID
			peerIsNorth := in.Status.ConnectionLevel > item.Status.ConnectionLevel
			peerIsSouth := in.Status.ConnectionLevel < item.Status.ConnectionLevel
			if peerFound && peerIsNorth {
				nicData.Direction = CDirectionNorth
			}
			if peerFound && peerIsSouth {
				nicData.Direction = CDirectionSouth
			}
		}
	}
}

// Creates map with switches' connection levels as keys
// and slices of Switch resources as values.
// Return ConnectionsMap and sorted slice of existing
// connection levels.
func (in *SwitchList) buildConnectionMap() (ConnectionsMap, []uint8) {
	connectionsMap := make(ConnectionsMap)
	keys := make([]uint8, 0)
	for _, item := range in.Items {
		if item.Status.State == CEmptyString {
			continue
		}
		list, ok := connectionsMap[item.Status.ConnectionLevel]
		if !ok {
			list = &SwitchList{}
			list.Items = append(list.Items, item)
			connectionsMap[item.Status.ConnectionLevel] = list
			keys = append(keys, item.Status.ConnectionLevel)
			continue
		}
		list.Items = append(list.Items, item)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return connectionsMap, keys
}

// SubnetsDefined checks whether switch subnets are defined
func (in *Switch) SubnetsDefined(ipv4Used, ipv6Used bool) bool {
	ipv4SubnetOk := false
	if (ipv4Used && in.Status.SubnetV4.CIDR != CEmptyString) ||
		(!ipv4Used && in.Status.SubnetV4.CIDR == CEmptyString) {
		ipv4SubnetOk = true
	}
	ipv6SubnetOk := false
	if (ipv6Used && in.Status.SubnetV6.CIDR != CEmptyString) ||
		(!ipv6Used && in.Status.SubnetV6.CIDR == CEmptyString) {
		ipv6SubnetOk = true
	}
	return ipv4SubnetOk && ipv6SubnetOk
}

// SwitchSubnetName returns the switch south subnet resource name
// depending on address family
func (in *Switch) SwitchSubnetName(af ipamv1alpha1.SubnetAddressType) string {
	suffix := "v4"
	if af == ipamv1alpha1.CIPv6SubnetType {
		suffix = "v6"
	}
	return fmt.Sprintf("%s-%s", in.Name, suffix)
}

// LoopbackIPResourceName returns the switch loopback IP resource name
// depending on address family
func (in *Switch) LoopbackIPResourceName(af ipamv1alpha1.SubnetAddressType) string {
	suffix := "ipv4"
	if af == ipamv1alpha1.CIPv6SubnetType {
		suffix = "ipv6"
	}
	return fmt.Sprintf("%s-lo-%s", in.Name, suffix)
}

// InterfaceSubnetName returns the interface subnet resource name
// depending on address family
func (in *Switch) InterfaceSubnetName(nic string, af ipamv1alpha1.SubnetAddressType) string {
	suffix := "v4"
	if af == ipamv1alpha1.CIPv6SubnetType {
		suffix = "v6"
	}
	return fmt.Sprintf("%s-%s-%s", in.Name, strings.ToLower(nic), suffix)
}

// InterfaceIPName returns the interface subnet resource name
// depending on address family
func (in *Switch) InterfaceIPName(nic string, af ipamv1alpha1.SubnetAddressType) string {
	suffix := "ipv4"
	if af == ipamv1alpha1.CIPv6SubnetType {
		suffix = "ipv6"
	}
	return fmt.Sprintf("%s-%s-%s", in.Name, strings.ToLower(nic), suffix)
}

// GetAddressCount defines the amount of needed ip addresses according to the
// number of switch ports, used lanes and address type (IPv4 or IPv6).
func (in *Switch) GetAddressCount(af ipamv1alpha1.SubnetAddressType) (count int64) {
	multiplier := CIPv4AddressesPerLane
	if af == ipamv1alpha1.CIPv6SubnetType {
		multiplier = CIPv6AddressesPerLane
	}
	for nic, nicData := range in.Status.Interfaces {
		if strings.HasPrefix(nic, CSwitchPortPrefix) {
			count += int64(nicData.Lanes * multiplier)
		}
	}
	return
}

// LoopbackAddressesDefined checks whether loopback addresses are filled
// for the switch
func (in *Switch) LoopbackAddressesDefined(ipv4Used, ipv6Used bool) bool {
	loopbackV4Ok := false
	if (ipv4Used && in.Status.LoopbackV4.Address != CEmptyString) ||
		(!ipv4Used && in.Status.LoopbackV4.Address == CEmptyString) {
		loopbackV4Ok = true
	}
	loopbackV6Ok := false
	if (ipv6Used && in.Status.LoopbackV6.Address != CEmptyString) ||
		(!ipv6Used && in.Status.LoopbackV6.Address == CEmptyString) {
		loopbackV6Ok = true
	}
	return loopbackV4Ok && loopbackV6Ok
}

// UndefinedLoopbackAF returns the list of address families for which
// loopback addresses are not defined
func (in *Switch) UndefinedLoopbackAF() (afs []ipamv1alpha1.SubnetAddressType) {
	if in.Status.LoopbackV4.Address == CEmptyString {
		afs = append(afs, ipamv1alpha1.CIPv4SubnetType)
	}
	if in.Status.LoopbackV6.Address == CEmptyString {
		afs = append(afs, ipamv1alpha1.CIPv6SubnetType)
	}
	return
}

// NICsAddressesDefined checks whether switch interfaces addresses
// are defined and match interface's direction
func (in *Switch) NICsAddressesDefined(ipv4Used, ipv6Used bool, list *SwitchList) bool {
	nicIPsFilled := in.nicsIPsFilled(ipv4Used, ipv6Used)
	nicIPsCorrect := in.nicsIPsCorrect(ipv4Used, ipv6Used, list)
	return nicIPsFilled && nicIPsCorrect
}

func (in *Switch) nicsIPsFilled(ipv4Used, ipv6Used bool) bool {
	for _, nicData := range in.Status.Interfaces {
		if ipv4Used && nicData.IPv4.Address == CEmptyString {
			return false
		}
		if ipv6Used && nicData.IPv6.Address == CEmptyString {
			return false
		}
	}
	return true
}

func (in *Switch) nicsIPsCorrect(ipv4Used, ipv6Used bool, list *SwitchList) bool {
	for _, nicData := range in.Status.Interfaces {
		if nicData.Direction == CDirectionNorth {
			if !in.nicIPsMatchNorthPeers(ipv4Used, ipv6Used, nicData, list) {
				return false
			}
			continue
		}
		if ipv4Used {
			_, subnetV4, _ := net.ParseCIDR(in.Status.SubnetV4.CIDR)
			nicIPv4, _, _ := net.ParseCIDR(nicData.IPv4.Address)
			if !subnetV4.Contains(nicIPv4) {
				return false
			}
		}
		if ipv6Used {
			_, subnetV6, _ := net.ParseCIDR(in.Status.SubnetV6.CIDR)
			nicIPv6, _, _ := net.ParseCIDR(nicData.IPv6.Address)
			if !subnetV6.Contains(nicIPv6) {
				return false
			}
		}
	}
	return true
}

func (in *Switch) nicIPsMatchNorthPeers(ipv4Used, ipv6Used bool, nicData *InterfaceSpec, list *SwitchList) bool {
	for _, peer := range list.Items {
		if peer.NamespacedName() != nicData.Peer.ResourceReference.NamespacedName() {
			continue
		}
		if ipv4Used {
			peerSubnetV4Defined := peer.Status.SubnetV4.CIDR != CEmptyString
			nicAddressV4Defined := nicData.IPv4.Address != CEmptyString
			if !peerSubnetV4Defined || !nicAddressV4Defined {
				return false
			}
			_, subnetV4, _ := net.ParseCIDR(peer.Status.SubnetV4.CIDR)
			nicIPv4, _, _ := net.ParseCIDR(nicData.IPv4.Address)
			if !subnetV4.Contains(nicIPv4) {
				return false
			}
		}

		if ipv6Used {
			peerSubnetV6Defined := peer.Status.SubnetV6.CIDR != CEmptyString
			nicAddressV6Defined := nicData.IPv6.Address != CEmptyString
			if !peerSubnetV6Defined || !nicAddressV6Defined {
				return false
			}
			_, subnetV6, _ := net.ParseCIDR(peer.Status.SubnetV6.CIDR)
			nicIPv6, _, _ := net.ParseCIDR(nicData.IPv6.Address)
			if !subnetV6.Contains(nicIPv6) {
				return false
			}
		}
		break
	}
	return true
}

// UpdateNorthNICsIP updates ipv4 and ipv6 addresses of interfaces
//// that considered to be north
func (in *Switch) UpdateNorthNICsIP(ipv4Used, ipv6Used bool, list *SwitchList) (err error) {
	for nic, nicData := range in.Status.Interfaces {
		if !strings.HasPrefix(nic, CSwitchPortPrefix) {
			continue
		}
		if nicData.Direction == CDirectionSouth {
			continue
		}
		for _, item := range list.Items {
			if nicData.Peer.ResourceReference.NamespacedName() != item.NamespacedName() {
				continue
			}
			peerNICData := item.Status.Interfaces[nicData.Peer.PortDescription]
			if ipv4Used {
				if nicAddressV4 := peerNICData.RequestAddress(ipamv1alpha1.CIPv4SubnetType); nicAddressV4 != nil {
					nicData.IPv4.Address = fmt.Sprintf("%s/%d", nicAddressV4.String(), CIPv4InterfaceSubnetMask)
				}
			}
			if ipv6Used {
				if nicAddressV6 := peerNICData.RequestAddress(ipamv1alpha1.CIPv6SubnetType); nicAddressV6 != nil {
					nicData.IPv6.Address = fmt.Sprintf("%s/%d", nicAddressV6.String(), CIPv6InterfaceSubnetMask)
				}
			}
		}
	}
	return
}

// UpdateSouthNICsIP updates ipv4 and ipv6 addresses of interfaces
// that considered to be south
func (in *Switch) UpdateSouthNICsIP(ipv4Used, ipv6Used bool) (err error) {
	for nic, nicData := range in.Status.Interfaces {
		if ipv4Used {
			_, switchSubnetV4, err := net.ParseCIDR(in.Status.SubnetV4.CIDR)
			if err != nil {
				return err
			}
			nicSubnetV4 := getInterfaceSubnet(nic, CSwitchPortPrefix, switchSubnetV4, ipamv1alpha1.CIPv4SubnetType)
			nicAddressV4, err := gocidr.Host(nicSubnetV4, 1)
			if err != nil {
				return err
			}
			nicData.IPv4.Address = fmt.Sprintf("%s/%d", nicAddressV4.String(), CIPv4InterfaceSubnetMask)
		}

		if ipv6Used {
			_, switchSubnetV6, err := net.ParseCIDR(in.Status.SubnetV6.CIDR)
			if err != nil {
				return err
			}
			nicSubnetV6 := getInterfaceSubnet(nic, CSwitchPortPrefix, switchSubnetV6, ipamv1alpha1.CIPv6SubnetType)
			nicAddressV6, err := gocidr.Host(nicSubnetV6, 0)
			if err != nil {
				return err
			}
			nicData.IPv6.Address = fmt.Sprintf("%s/%d", nicAddressV6.String(), CIPv6InterfaceSubnetMask)
		}
	}
	return
}

func getInterfaceSubnet(name string, namePrefix string, network *net.IPNet, af ipamv1alpha1.SubnetAddressType) *net.IPNet {
	index, _ := strconv.Atoi(strings.ReplaceAll(name, namePrefix, CEmptyString))
	prefix, _ := network.Mask.Size()
	ifaceNet, _ := gocidr.Subnet(network, getInterfaceSubnetMaskLength(af)-prefix, index)
	return ifaceNet
}

func getInterfaceSubnetMaskLength(af ipamv1alpha1.SubnetAddressType) int {
	if af == ipamv1alpha1.CIPv4SubnetType {
		return CIPv4InterfaceSubnetMask
	}
	return CIPv6InterfaceSubnetMask
}

// RequestAddress returns the IP address next for the
// IP address of the interface.
func (in *InterfaceSpec) RequestAddress(af ipamv1alpha1.SubnetAddressType) (ip net.IP) {
	switch af {
	case ipamv1alpha1.CIPv4SubnetType:
		if in.IPv4.Address == CEmptyString {
			return
		}
		_, cidr, _ := net.ParseCIDR(in.IPv4.Address)
		ip, _ = gocidr.Host(cidr, 2)
	case ipamv1alpha1.CIPv6SubnetType:
		if in.IPv6.Address == CEmptyString {
			return
		}
		_, cidr, _ := net.ParseCIDR(in.IPv6.Address)
		ip, _ = gocidr.Host(cidr, 1)
	}
	return
}

func (in *Switch) IPAMResourcesCreated() bool {
	for _, nicData := range in.Status.Interfaces {
		if nicData.Direction == CDirectionNorth {
			continue
		}
		if nicData.IPv4.ResourceReference == nil || nicData.IPv6.ResourceReference == nil {
			return false
		}
	}
	return true
}
