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
	"strings"

	gocidr "github.com/apparentlymart/go-cidr/cidr"
	subnetv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type State string
type Role string
type PeerType string

//ConnectionsMap
//+kubebuilder:object:generate=false
type ConnectionsMap map[uint8][]Switch

//SwitchSpec defines the desired state of Switch
//+kubebuilder:object:generate=true
type SwitchSpec struct {
	//Hostname
	//+kubebuilder:validation:Required
	Hostname string `json:"hostname"`
	//Location referring to the switch location
	//+kubebuilder:validation:Optional
	Location *LocationSpec `json:"location,omitempty"`
	//TotalPorts referring to network interfaces total count
	//+kubebuilder:validation:Required
	TotalPorts uint64 `json:"totalPorts"`
	//SwitchPorts referring to non-management network interfaces count
	//+kubebuilder:validation:Required
	SwitchPorts uint64 `json:"switchPorts"`
	//SwitchDistro referring to switch OS information
	//+kubebuilder:validation:Optional
	Distro *SwitchDistroSpec `json:"distro,omitempty"`
	//SwitchChassis referring to switch hardware information
	//+kubebuilder:validation:Required
	Chassis *SwitchChassisSpec `json:"chassis"`
	//Interfaces referring to details about network interfaces
	//+kubebuilder:validation:Optional
	Interfaces map[string]*InterfaceSpec `json:"interfaces,omitempty"`
	//SouthSubnet referring to south IPv4 subnet
	//+kubebuilder:validation:Optional
	//+nullable
	SouthSubnetV4 *SwitchSubnetSpec `json:"southSubnetV4,omitempty"`
	//SouthSubnet referring to south IPv6 subnet
	//+kubebuilder:validation:Optional
	//+nullable
	SouthSubnetV6 *SwitchSubnetSpec `json:"southSubnetV6,omitempty"`
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
	//MacAddress referring to interface's MAC address
	//+kubebuilder:validation:Optional
	MacAddress string `json:"macAddress,omitempty"`
	//IPv4 referring to interface's IPv4 address
	//+kubebuilder:validation:Optional
	IPv4 string `json:"ipv4,omitempty"`
	//IPv6 referring to interface's IPv6 address
	//+kubebuilder:validation:Optional
	IPv6 string `json:"ipv6,omitempty"`
	//PeerType referring to neighbour type
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=Machine;Switch
	PeerType PeerType `json:"peerType,omitempty"`
	//PeerSystemName
	//+kubebuilder:validation:Optional
	PeerSystemName string `json:"peerSystemName,omitempty"`
	//PeerChassisID
	//+kubebuilder:validation:Optional
	PeerChassisID string `json:"peerChassisId,omitempty"`
	//PeerPortID
	//+kubebuilder:validation:Optional
	PeerPortID string `json:"peerPortId,omitempty"`
	//PeerPortDescription
	//+kubebuilder:validation:Optional
	PeerPortDescription string `json:"peerPortDescription,omitempty"`
	//Ndp flag defines whether ndp data received from neighbour and it is in Reachable state
	//+kubebuilder:validation:default=false
	Ndp bool `json:"ndp"`
}

// SwitchStatus defines the observed state of Switch
type SwitchStatus struct {
	//Role referring to switch's role: leaf or spine
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Enum=Leaf;Spine
	Role Role `json:"role"`
	// ConnectionLevel refers the level of the connection
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Minimum=0
	//+kubebuilder:validation:Maximum=255
	ConnectionLevel uint8 `json:"connectionLevel"`
	// NorthSwitches refers to up-level switch
	//+kubebuilder:validation:Optional
	NorthConnections *ConnectionsSpec `json:"northConnections,omitempty"`
	// SouthSwitches refers to down-level switch
	//+kubebuilder:validation:Optional
	SouthConnections *ConnectionsSpec `json:"southConnections,omitempty"`
	//State referring to current switch state
	//kubebuilder:validation:Enum=Finished;Deleting;Define peers;Define addresses
	State State `json:"state"`
	//ScanPorts flag determining whether scanning of ports is requested
	//+kubebuilder:validation:Required
	ScanPorts bool `json:"scanPorts"`
}

// ConnectionsSpec defines upstream switches count and properties
//+kubebuilder:object:generate=true
type ConnectionsSpec struct {
	// Count refers to upstream switches count
	//+kubebuilder:validation:Optional
	Count int `json:"count"`
	// Peers refers to connected upstream switches
	//+kubebuilder:validation:Optional
	Peers map[string]*PeerSpec `json:"peers"`
}

// PeerSpec defines switch connected to another switch
//+kubebuilder:object:generate=true
type PeerSpec struct {
	// Name refers to switch's name
	//+kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// Namespace refers to switch's namespace
	//+kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
	// ChassisID refers to switch's chassis id
	//+kubebuilder:validation:Required
	ChassisID string `json:"chassisId"`
	//Type referring to neighbour type
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=Machine;Switch
	Type PeerType `json:"type,omitempty"`
	//PortName
	//+kubebuilder:validation:Optional
	PortName string `json:"portName,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=sw
//+kubebuilder:printcolumn:name="Hostname",type=string,JSONPath=`.spec.hostname`,description="Switch's hostname"
//+kubebuilder:printcolumn:name="OS",type=string,JSONPath=`.spec.distro.os`,description="OS running on switch"
//+kubebuilder:printcolumn:name="SwitchPorts",type=integer,JSONPath=`.spec.switchPorts`,description="Total amount of non-management network interfaces"
//+kubebuilder:printcolumn:name="Role",type=string,JSONPath=`.status.role`,description="switch's role"
//+kubebuilder:printcolumn:name="ConnectionLevel",type=integer,JSONPath=`.status.connectionLevel`,description="Vertical level of switch connection"
//+kubebuilder:printcolumn:name="SouthSubnetV4",type=string,JSONPath=`.spec.southSubnetV4.cidr`,description="South IPv4 subnet"
//+kubebuilder:printcolumn:name="SouthSubnetV6",type=string,JSONPath=`.spec.southSubnetV6.cidr`,description="South IPv6 subnet"
//+kubebuilder:printcolumn:name="ScanPorts",type=boolean,JSONPath=`.status.scanPorts`,description="Request for scan ports"
//+kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`,description="Switch processing state"

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

// Determines whether Switch resources with connection
// level equals to zero already exists.
// Return true if exists, false otherwise.
func (in *ConnectionsMap) topLevelSpinesDefined() bool {
	if switches, ok := (*in)[0]; !ok {
		return false
	} else {
		for _, sw := range switches {
			if sw.Status.State == EmptyString {
				return false
			}
		}
	}
	return true
}

// Creates map with switches' connection levels as keys
// and slices of Switch resources as values.
// Return ConnectionsMap and sorted slice of existing
// connection levels.
func (in *SwitchList) buildConnectionMap() (ConnectionsMap, []uint8) {
	connectionsMap := make(ConnectionsMap)
	keys := make([]uint8, 0)
	for _, item := range in.Items {
		if _, ok := connectionsMap[item.Status.ConnectionLevel]; !ok {
			connectionsMap[item.Status.ConnectionLevel] = []Switch{item}
			keys = append(keys, item.Status.ConnectionLevel)
		} else {
			connectionsMap[item.Status.ConnectionLevel] = append(connectionsMap[item.Status.ConnectionLevel], item)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return connectionsMap, keys
}

// Returns minimal existing connection level value.
func (in *SwitchList) minimumConnectionLevel() uint8 {
	result := uint8(255)
	for _, item := range in.Items {
		if item.Status.ConnectionLevel < result {
			result = item.Status.ConnectionLevel
		}
	}
	return result
}

// GetTopLevelSwitch searches for Switch resource with
// connection level equals to zero in ConnectionsMap.
// Return nil in case Switch was not found.
func (in *SwitchList) GetTopLevelSwitch() *Switch {
	connectionsMap, _ := in.buildConnectionMap()
	if switches, ok := connectionsMap[0]; ok {
		return &switches[0]
	}
	return nil
}

// AllConnectionsOk checks whether all Switch resources
// in SwitchList has determined their peer connections.
// Return false if not.
func (in *SwitchList) AllConnectionsOk() bool {
	for _, sw := range in.Items {
		if sw.Status.State == StateDefinePeers {
			return false
		}
	}
	return true
}

// Constructs map of PeerSpec objects based on stored
// Switch interfaces.
// Return map where interface name is a key and PeerSpec
// is a value.
func (in *Switch) getBaseConnections() map[string]*PeerSpec {
	result := make(map[string]*PeerSpec)
	for name, data := range in.Spec.Interfaces {
		if strings.HasPrefix(name, "Ethernet") && data.PeerChassisID != EmptyString {
			result[name] = &PeerSpec{
				Name:      EmptyString,
				Namespace: EmptyString,
				ChassisID: data.PeerChassisID,
				Type:      data.PeerType,
				PortName:  data.PeerPortDescription,
			}
		}
	}
	return result
}

// Defines the Switch role according to existing peers
func (in *Switch) getRole(peers map[string]*PeerSpec) Role {
	for _, data := range peers {
		if data.Type == MachineType {
			return LeafRole
		}
	}
	return SpineRole
}

// Moves peers between south and north peer lists
// according to changes in peers connection levels.
func (in *Switch) movePeers(list *SwitchList) {
	if in.Status.ConnectionLevel == 0 {
		for name, data := range in.Status.NorthConnections.Peers {
			in.Status.SouthConnections.Peers[name] = data
			delete(in.Status.NorthConnections.Peers, name)
		}
	} else {
		for _, item := range list.Items {
			for name, data := range in.Status.SouthConnections.Peers {
				if data.ChassisID == item.Spec.Chassis.ChassisID {
					if item.Status.ConnectionLevel < in.Status.ConnectionLevel {
						in.Status.NorthConnections.Peers[name] = data
						delete(in.Status.SouthConnections.Peers, name)
					}
				}
			}
			for name, data := range in.Status.NorthConnections.Peers {
				if data.ChassisID == item.Spec.Chassis.ChassisID {
					if item.Status.ConnectionLevel > in.Status.ConnectionLevel {
						in.Status.SouthConnections.Peers[name] = data
						delete(in.Status.NorthConnections.Peers, name)
					}
				}
			}
		}
	}
}

// Determines what Switch resources are the known
// peers for current Switch.
// Return SwitchList containing peers.
func (in *Switch) getPeers(list []Switch) *SwitchList {
	result := &SwitchList{}
	for _, item := range list {
		for _, data := range in.Spec.Interfaces {
			if data.PeerChassisID == item.Spec.Chassis.ChassisID {
				result.Items = append(result.Items, item)
			}
		}
	}
	return result
}

// Rewrites stored north peers info with fully defined PeerSpec
func (in *Switch) updateNorthPeers(list *SwitchList) {
	for _, item := range list.Items {
		for name, data := range item.Spec.Interfaces {
			if data.PeerChassisID == in.Spec.Chassis.ChassisID {
				in.Status.NorthConnections.Peers[data.PeerPortDescription] = &PeerSpec{
					Name:      item.Name,
					Namespace: item.Namespace,
					ChassisID: item.Spec.Chassis.ChassisID,
					Type:      SwitchType,
					PortName:  name,
				}
			}
		}
	}
}

// Rewrites stored south peers info with fully defined PeerSpec
func (in *Switch) updateSouthPeers(list *SwitchList) {
	for _, item := range list.Items {
		for name, data := range item.Spec.Interfaces {
			if data.PeerChassisID == in.Spec.Chassis.ChassisID {
				in.Status.SouthConnections.Peers[data.PeerPortDescription] = &PeerSpec{
					Name:      item.Name,
					Namespace: item.Namespace,
					ChassisID: item.Spec.Chassis.ChassisID,
					Type:      SwitchType,
					PortName:  name,
				}
			}
		}
	}
}

// Checks whether all stored peers are unique and fully defined.
// Return true if so, false otherwise.
func (in *Switch) peersOk(swl *SwitchList) bool {
	for _, sw := range swl.Items {
		for inf, peer := range in.Status.SouthConnections.Peers {
			if peer.ChassisID == sw.Spec.Chassis.ChassisID && peer.Name == EmptyString {
				return false
			}
			if _, ok := in.Status.NorthConnections.Peers[inf]; ok {
				return false
			}
		}
		for inf, peer := range in.Status.NorthConnections.Peers {
			if peer.ChassisID == sw.Spec.Chassis.ChassisID && peer.Name == EmptyString {
				return false
			}
			if _, ok := in.Status.SouthConnections.Peers[inf]; ok {
				return false
			}
		}
	}
	return true
}

// Checks whether all stored peers defined correctly: north
// peers should have connection level less by one than current
// Switch, south peers should have connection level greater by
// one than current switch.
// Return true if peers defined correctly, false otherwise.
func (in *Switch) connectionsOk(list *SwitchList) bool {
	if in.Status.ConnectionLevel == 0 && in.Status.NorthConnections.Count != 0 {
		return false
	}
	for _, item := range list.Items {
		for _, peer := range in.Status.NorthConnections.Peers {
			if peer.ChassisID == item.Spec.Chassis.ChassisID {
				if !item.switchInSouthPeers(in) {
					return false
				}
				if in.Status.ConnectionLevel != item.Status.ConnectionLevel+1 {
					return false
				}
			}
		}
		for _, peer := range in.Status.SouthConnections.Peers {
			if peer.ChassisID == item.Spec.Chassis.ChassisID {
				if !item.switchInNorthPeers(in) {
					return false
				}
				if in.Status.ConnectionLevel != item.Status.ConnectionLevel-1 {
					return false
				}
			}
		}
	}
	return true
}

// Checks whether Switch provided as argument is in the south
// peers of current Switch.
func (in *Switch) switchInSouthPeers(tgt *Switch) bool {
	if in.Status.SouthConnections == nil {
		return false
	}
	for _, peer := range in.Status.SouthConnections.Peers {
		if peer.ChassisID == tgt.Spec.Chassis.ChassisID {
			return true
		}
	}
	return false
}

// Checks whether Switch provided as argument is in the north
// peers of current Switch.
func (in *Switch) switchInNorthPeers(tgt *Switch) bool {
	if in.Status.NorthConnections == nil {
		return false
	}
	for _, peer := range in.Status.NorthConnections.Peers {
		if peer.ChassisID == tgt.Spec.Chassis.ChassisID {
			return true
		}
	}
	return false
}

// Defines the amount of needed ip addresses according to the
// number of switch ports and address type (IPv4 or IPv6).
func (in *Switch) getAddressNeededCount(addrType subnetv1alpha1.SubnetAddressType) int64 {
	if addrType == subnetv1alpha1.CIPv4SubnetType {
		return int64(in.Spec.SwitchPorts * CIPv4AddressesPerPort)
	} else {
		return int64(in.Spec.SwitchPorts * CIPv6AddressesPerPort)
	}
}

// NamespacedName returns switch's name and namespace as
// built-in type.
func (in *Switch) NamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: in.Namespace,
		Name:      in.Name,
	}
}

// Prepare constructs Switch resource for creation from
// provided Inventory resource.
func (in *Switch) Prepare(src *inventoriesv1alpha1.Inventory) {
	interfaces, switchPorts := PrepareInterfaces(src.Spec.NICs.NICs)
	in.ObjectMeta = metav1.ObjectMeta{
		Name:      src.Name,
		Namespace: CNamespace,
	}
	in.Spec = SwitchSpec{
		Hostname:    src.Spec.Host.Name,
		Location:    &LocationSpec{},
		TotalPorts:  src.Spec.NICs.Count,
		SwitchPorts: switchPorts,
		Distro: &SwitchDistroSpec{
			OS:      CSonicSwitchOs,
			Version: src.Spec.Distro.CommitId,
			ASIC:    src.Spec.Distro.AsicType,
		},
		Chassis: &SwitchChassisSpec{
			Manufacturer: src.Spec.System.Manufacturer,
			SKU:          src.Spec.System.ProductSKU,
			Serial:       src.Spec.System.SerialNumber,
			ChassisID:    getChassisId(src.Spec.NICs.NICs),
		},
		Interfaces: interfaces,
	}
}

// UpdateInterfaces fills in missing, due to absent LLDPs,
// but existing NDPs, data in switch interfaces, if it's
// peer has LLDP data about current Switch.
func (in *Switch) UpdateInterfaces(swl *SwitchList) {
	for inf, data := range in.Spec.Interfaces {
		if strings.HasPrefix(inf, "Ethernet") && data.PeerChassisID == EmptyString && data.Ndp {
			for _, item := range swl.Items {
				for peerInf, peerInfData := range item.Spec.Interfaces {
					if peerInfData.PeerChassisID == in.Spec.Chassis.ChassisID && peerInfData.PeerPortDescription == inf {
						data.PeerChassisID = item.Spec.Chassis.ChassisID
						data.PeerType = SwitchType
						data.PeerPortDescription = peerInf
						data.PeerSystemName = item.Spec.Hostname
					}
				}
			}
		}
	}
}

// InterfacesUpdated checks whether interfaces data is
// missing or not.
func (in *Switch) InterfacesUpdated(swl *SwitchList) bool {
	for inf, data := range in.Spec.Interfaces {
		if strings.HasPrefix(inf, "Ethernet") && data.PeerChassisID == EmptyString && data.Ndp {
			for _, item := range swl.Items {
				for _, peerInfData := range item.Spec.Interfaces {
					if peerInfData.PeerChassisID == in.Spec.Chassis.ChassisID && peerInfData.PeerPortDescription == inf {
						return false
					}
				}
			}
		}
	}
	return true
}

// FillStatusOnCreate fills Switch status on resource creation.
func (in *Switch) FillStatusOnCreate() {
	peers := in.getBaseConnections()
	in.Status = SwitchStatus{
		Role:            in.getRole(peers),
		ConnectionLevel: 255,
		NorthConnections: &ConnectionsSpec{
			Count: 0,
			Peers: make(map[string]*PeerSpec),
		},
		SouthConnections: &ConnectionsSpec{
			Count: len(peers),
			Peers: peers,
		},
		State:     StateDefinePeers,
		ScanPorts: false,
	}
}

// FlushStatusOnDelete flushes Switch status on resource deletion.
func (in *Switch) FlushStatusOnDelete() {
	in.Status.ConnectionLevel = 255
	in.Status.NorthConnections = &ConnectionsSpec{
		Count: 0,
		Peers: make(map[string]*PeerSpec),
	}
	in.Status.SouthConnections = &ConnectionsSpec{
		Count: 0,
		Peers: make(map[string]*PeerSpec),
	}
	in.Status.State = StateDeleting
	in.Status.ScanPorts = false
}

// GetListFilter builds list options object
func (in *Switch) GetListFilter() (*client.ListOptions, error) {
	labelsReq, err := labels.NewRequirement(LabelChassisId, selection.In, []string{MacToLabel(in.Spec.Chassis.ChassisID)})
	if err != nil {
		return nil, err
	}
	selector := labels.NewSelector()
	selector = selector.Add(*labelsReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         100,
	}
	return opts, nil
}

// UpdatePeersData rewrites existing south peers data
// with fully filled PeerSpec according to info stored
// in interfaces specs.
func (in *Switch) UpdatePeersData(list *SwitchList) {
	for name, data := range in.Spec.Interfaces {
		if strings.HasPrefix(name, "Ethernet") && data.PeerChassisID != EmptyString {
			_, found := in.Status.NorthConnections.Peers[name]
			if !found {
				for _, item := range list.Items {
					if item.Spec.Chassis.ChassisID == data.PeerChassisID {
						in.Status.SouthConnections.Peers[name] = &PeerSpec{
							Name:      item.Name,
							Namespace: item.Namespace,
							ChassisID: data.PeerChassisID,
							Type:      data.PeerType,
							PortName:  data.PeerPortDescription,
						}
					}
				}
			}
		}
	}
}

// UpdateConnectionLevel updates switch's connection level
// and peers info.
func (in *Switch) UpdateConnectionLevel(list *SwitchList) {
	connectionsMap, keys := list.buildConnectionMap()
	if !connectionsMap.topLevelSpinesDefined() {
		return
	}
	if in.Status.ConnectionLevel == 0 {
		in.movePeers(list)
	} else {
		for _, connectionLevel := range keys {
			switches := connectionsMap[connectionLevel]
			northPeers := in.getPeers(switches)
			if len(northPeers.Items) > 0 {
				minConnectionLevel := northPeers.minimumConnectionLevel()
				if minConnectionLevel != 255 && minConnectionLevel < in.Status.ConnectionLevel {
					in.Status.ConnectionLevel = minConnectionLevel + 1
					in.updateNorthPeers(northPeers)
					in.movePeers(list)
				}
			}
		}
	}
	in.Status.NorthConnections.Count = len(in.Status.NorthConnections.Peers)
	in.Status.SouthConnections.Count = len(in.Status.SouthConnections.Peers)
}

// PeersProcessingFinished checks whether peers are
// correctly determined for all existing switches.
func (in *Switch) PeersProcessingFinished(swl *SwitchList, swa *SwitchAssignment) bool {
	if swa != nil && in.Status.ConnectionLevel != 0 {
		return false
	}
	if in.Status.ConnectionLevel == 255 {
		return false
	}
	if !in.peersOk(swl) || !in.connectionsOk(swl) {
		return false
	}
	return true
}

// GetSuitableSubnet looks up for subnet resource, that fits
// with number of parameters as region, availability zone,
// amount of available addresses and address type, to reserve
// CIDR for Switch.
// Returns suitable CIDR and subnet or error.
func (in *Switch) GetSuitableSubnet(
	subnets *subnetv1alpha1.SubnetList,
	addressType subnetv1alpha1.SubnetAddressType,
	regions []string,
	zones []string) (*subnetv1alpha1.CIDR, *subnetv1alpha1.Subnet, error) {

	addressesNeeded := in.getAddressNeededCount(addressType)
	for _, sn := range subnets.Items {
		if sn.Spec.NetworkName == "underlay" &&
			sn.Status.Type == addressType &&
			reflect.DeepEqual(sn.Spec.Regions, regions) &&
			reflect.DeepEqual(sn.Spec.AvailabilityZones, zones) {
			addressesLeft := sn.Status.CapacityLeft
			if sn.Status.Type == addressType && addressesLeft.CmpInt64(addressesNeeded) >= 0 {
				minVacantCIDR := getMinimalVacantCIDR(sn.Status.Vacant, addressType, addressesNeeded)
				mask := getNeededMask(addressType, float64(addressesNeeded))
				addr := minVacantCIDR.Net.IP
				network := &net.IPNet{
					IP:   addr,
					Mask: mask,
				}
				cidrCandidate := &subnetv1alpha1.CIDR{Net: network}
				if sn.CanReserve(cidrCandidate) {
					if err := sn.Reserve(cidrCandidate); err != nil {
						return nil, nil, err
					} else {
						return cidrCandidate, &sn, nil
					}
				}
			}
		}
	}
	return nil, nil, nil
}

// UpdateSouthInterfacesAddresses defines addresses for
// switch interfaces according to the switch's south subnets.
func (in *Switch) UpdateSouthInterfacesAddresses() {
	if in.Spec.SouthSubnetV4 != nil {
		_, network, _ := net.ParseCIDR(in.Spec.SouthSubnetV4.CIDR)
		for inf := range in.Status.SouthConnections.Peers {
			iface := in.Spec.Interfaces[inf]
			ifaceSubnet := getInterfaceSubnet(inf, network, subnetv1alpha1.CIPv4SubnetType)
			ifaceAddress, _ := gocidr.Host(ifaceSubnet, 1)
			iface.IPv4 = fmt.Sprintf("%s/%d", ifaceAddress.String(), CIPv4InterfaceSubnetMask)
		}
	}
	if in.Spec.SouthSubnetV6 != nil {
		_, network, _ := net.ParseCIDR(in.Spec.SouthSubnetV6.CIDR)
		for inf := range in.Status.SouthConnections.Peers {
			iface := in.Spec.Interfaces[inf]
			ifaceSubnet := getInterfaceSubnet(inf, network, subnetv1alpha1.CIPv6SubnetType)
			ifaceAddress, _ := gocidr.Host(ifaceSubnet, 0)
			iface.IPv6 = fmt.Sprintf("%s/%d", ifaceAddress.String(), CIPv6InterfaceSubnetMask)
		}
	}
}

// UpdateNorthInterfacesAddresses defines addresses for
// switch interfaces, that are connected to upstream switches,
// according to the peers' interfaces addresses.
func (in *Switch) UpdateNorthInterfacesAddresses(swl *SwitchList) {
	for inf, peer := range in.Status.NorthConnections.Peers {
		for _, item := range swl.Items {
			if peer.ChassisID == item.Spec.Chassis.ChassisID {
				iface := in.Spec.Interfaces[inf]
				peerIface := item.Spec.Interfaces[peer.PortName]
				ipv4Addr := peerIface.RequestAddress(subnetv1alpha1.CIPv4SubnetType)
				ipv6Addr := peerIface.RequestAddress(subnetv1alpha1.CIPv6SubnetType)
				if ipv4Addr != nil {
					iface.IPv4 = fmt.Sprintf("%s/%d", ipv4Addr.String(), CIPv4InterfaceSubnetMask)
				}
				if ipv6Addr != nil {
					iface.IPv6 = fmt.Sprintf("%s/%d", ipv6Addr.String(), CIPv6InterfaceSubnetMask)
				}
			}
		}
	}
}

// AddressesDefined checks whether ip addresses defined
// for all used switch interfaces.
func (in *Switch) AddressesDefined() bool {
	for inf, data := range in.Spec.Interfaces {
		if strings.HasPrefix(inf, "Ethernet") && data.PeerChassisID != EmptyString {
			if data.IPv4 == EmptyString || data.IPv6 == EmptyString {
				return false
			}
		}
	}
	return true
}

// UpdateInterfacesFromInventory fulfills switch's interfaces
// data according to updated inventory data
func (in *Switch) UpdateInterfacesFromInventory(updated map[string]*InterfaceSpec) {
	for inf := range in.Spec.Interfaces {
		if _, ok := updated[inf]; !ok {
			delete(in.Spec.Interfaces, inf)
			delete(in.Status.SouthConnections.Peers, inf)
			delete(in.Status.NorthConnections.Peers, inf)
		}
	}
	for inf, data := range updated {
		stored, ok := in.Spec.Interfaces[inf]
		if !ok {
			in.Spec.Interfaces[inf] = data
		} else {
			stored.PeerType = data.PeerType
			stored.PeerChassisID = data.PeerChassisID
			stored.PeerSystemName = data.PeerSystemName
			stored.PeerPortID = data.PeerPortID
			stored.PeerPortDescription = data.PeerPortDescription
			stored.Ndp = data.Ndp
		}
	}
}

// PeersUpdateNeeded checks whether interfaces data
// was updated and peers info needed update.
func (in *Switch) PeersUpdateNeeded() bool {
	for name, data := range in.Spec.Interfaces {
		if strings.HasPrefix(name, "Ethernet") && data.PeerChassisID != EmptyString {
			_, northPeer := in.Status.NorthConnections.Peers[name]
			_, southPeer := in.Status.SouthConnections.Peers[name]
			if !northPeer && !southPeer {
				return true
			}
		}
	}
	return false
}

// UpdatePeersInfo updates peers data and switch role
// according to connected peers.
func (in *Switch) UpdatePeersInfo() {
	machinesConnected := false
	for name, data := range in.Spec.Interfaces {
		_, northPeer := in.Status.NorthConnections.Peers[name]
		_, southPeer := in.Status.SouthConnections.Peers[name]
		if northPeer || southPeer {
			continue
		}
		if strings.HasPrefix(name, "Ethernet") && data.PeerChassisID != EmptyString {
			in.Status.SouthConnections.Peers[name] = &PeerSpec{
				Name:      EmptyString,
				Namespace: EmptyString,
				ChassisID: data.PeerChassisID,
				Type:      data.PeerType,
				PortName:  data.PeerPortDescription,
			}
			if data.PeerType == MachineType {
				machinesConnected = true
			}
		}
	}
	if machinesConnected {
		in.Status.Role = LeafRole
	} else {
		in.Status.Role = SpineRole
	}
}

// RequestAddress returns the IP address next for the
// IP address of the interface.
func (in *InterfaceSpec) RequestAddress(addrType subnetv1alpha1.SubnetAddressType) net.IP {
	addr := net.IP{}
	switch addrType {
	case subnetv1alpha1.CIPv4SubnetType:
		if in.IPv4 == EmptyString {
			return nil
		}
		_, cidr, _ := net.ParseCIDR(in.IPv4)
		addr, _ = gocidr.Host(cidr, 2)
	case subnetv1alpha1.CIPv6SubnetType:
		if in.IPv6 == EmptyString {
			return nil
		}
		_, cidr, _ := net.ParseCIDR(in.IPv6)
		addr, _ = gocidr.Host(cidr, 1)
	}
	return addr
}
