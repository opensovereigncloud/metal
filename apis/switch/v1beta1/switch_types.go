/*
 * Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1beta1

import (
	"fmt"
	"math"
	"net"
	"reflect"
	"sort"
	"strconv"
	"strings"

	gocidr "github.com/apparentlymart/go-cidr/cidr"
	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

//todo: clientset for v1beta1

// SwitchSpec contains desired state of resulting Switch configuration
// +kubebuilder:object:generate=true
type SwitchSpec struct {
	// UUID is a unique system identifier
	//+kubebuilder:validation:Required
	UUID string `json:"uuid"`
	// Managed is a flag defining whether Switch object would be processed during reconciliation
	//+kubebuilder:validation:Required
	//+kubebuilder:default=true
	Managed bool `json:"managed"`
	// Cordon is a flag defining whether Switch object is taken offline
	//+kubebuilder:validation:Required
	//+kubebuilder:default=false
	Cordon bool `json:"cordon"`
	// TopSpine is a flag defining whether Switch is a top-level spine switch
	//+kubebuilder:validation:Required
	//+kubebuilder:default=false
	TopSpine bool `json:"topSpine"`
	// IPAM refers to selectors for subnets which will be used for Switch object
	//+kubebuilder:validation:Optional
	IPAM *IPAMSpec `json:"ipam,omitempty"`
	// Interfaces contains general configuration for all switch ports
	//+kubebuilder:validation:Optional
	Interfaces *InterfacesSpec `json:"interfaces,omitempty"`
}

// InterfacesSpec contains definitions for general switch ports' configuration
// +kubebuilder:object:generate=true
type InterfacesSpec struct {
	// Scan is a flag defining whether to run periodical scanning on switch ports
	//+kubebuilder:validation:Optional
	//+kubebuilder:default=true
	Scan bool `json:"scan"`
	// Defaults contains switch port parameters which will be applied to all ports of the switches
	//+kubebuilder:validation:Optional
	Defaults *PortParametersSpec `json:"defaults,omitempty"`
	// Overrides contains set of parameters which should be overridden for listed switch ports
	//+kubebuilder:validation:Optional
	Overrides []*InterfaceOverridesSpec `json:"overrides,omitempty"`
}

// InterfaceOverridesSpec contains overridden parameters for certain switch port
// +kubebuilder:object:generate=true
type InterfaceOverridesSpec struct {
	// Lanes refers to a number of lanes used by switch port
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Minimum=1
	//+kubebuilder:validation:Maximum=8
	Lanes *uint8 `json:"lanes,omitempty"`
	// MTU refers to maximum transmission unit value which should be applied on switch port
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Minimum=84
	//+kubebuilder:validation:Maximum=65535
	MTU *uint16 `json:"mtu,omitempty"`
	// Name refers to switch port name
	//+kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// State defines default state of switch port
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=up;down
	State *string `json:"state,omitempty"`
	// FEC refers to forward error correction method which should be applied on switch port
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=rs;none
	FEC *string `json:"fec,omitempty"`
	// IP contains a list of additional IP addresses for interface
	//+kubebuilder:validation:Optional
	IP []*AdditionalIPSpec `json:"ip,omitempty"`
}

// SwitchStatus contains observed state of Switch
// +kubebuilder:object:generate=true
type SwitchStatus struct {
	// TotalPorts refers to total number of ports
	//+kubebuilder:validation:Optional
	TotalPorts uint16 `json:"totalPorts,omitempty"`
	// SwitchPorts refers to the number of ports excluding management interfaces, loopback etc.
	//+kubebuilder:validation:Optional
	SwitchPorts uint16 `json:"switchPorts,omitempty"`
	// Role refers to switch's role
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Enum=spine;leaf;edge-leaf
	Role *string `json:"role,omitempty"`
	// ConnectionLevel refers to switch's current position in connection hierarchy
	//+kubebuilder:validation:Optional
	ConnectionLevel uint8 `json:"connectionLevel"`
	// Interfaces refers to switch's interfaces configuration
	//+kubebuilder:validation:Optional
	Interfaces map[string]*InterfaceSpec `json:"interfaces,omitempty"`
	// Subnets refers to the switch's south subnets
	//+kubebuilder:validation:Optional
	Subnets []*SubnetSpec `json:"subnets,omitempty"`
	// LoopbackAddresses refers to the switch's loopback addresses
	//+kubebuilder:validation:Optional
	LoopbackAddresses []*IPAddressSpec `json:"loopbackAddresses,omitempty"`
	// SwitchState contains information about current Switch object's processing state
	//+kubebuilder:validation:Optional
	SwitchState *SwitchStateSpec `json:"switch,omitempty"`
	// ConfigAgent contains information about current state of configuration agent
	// running on the switch
	//+kubebuilder:validation:Optional
	ConfigAgent *ConfigAgentStateSpec `json:"agent,omitempty"`
}

// InterfaceSpec defines the state of switch's interface
// +kubebuilder:object:generate=true
type InterfaceSpec struct {
	// MACAddress refers to the interface's hardware address
	//+kubebuilder:validation:Required
	// validation pattern
	MACAddress string `json:"macAddress"`
	// FEC refers to the current interface's forward error correction type
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Enum=none;rs;fc
	FEC string `json:"fec"`
	// MTU refers to the current value of interface's MTU
	//+kubebuilder:validation:Required
	MTU uint16 `json:"mtu"`
	// Speed refers to interface's speed
	//+kubebuilder:validation:Required
	Speed uint32 `json:"speed"`
	// Lanes refers to the number of lanes used by interface
	//+kubebuilder:validation:Required
	Lanes uint8 `json:"lanes"`
	// State refers to the current interface's operational state
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Enum=up;down
	State string `json:"state"`
	// IP contains a list of IP addresses that are assigned to interface
	//+kubebuilder:validation:Optional
	IP []*IPAddressSpec `json:"ip,omitempty"`
	// Direction refers to the interface's connection 'direction'
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Enum=north;south
	Direction string `json:"direction"`
	// Peer refers to the info about device connected to current switch port
	//+kubebuilder:validation:Optional
	Peer *PeerSpec `json:"peer,omitempty"`
}

// PeerSpec defines peer info
// +kubebuilder:object:generate=true
type PeerSpec struct {
	// Contains information to locate the referenced object
	//+kubebuilder:validation:Optional
	*ObjectReference `json:",inline"`
	// Contains LLDP info about peer
	//+kubebuilder:validation:Optional
	*PeerInfoSpec `json:",inline"`
}

// PeerInfoSpec contains LLDP info about peer
// +kubebuilder:object:generate=true
type PeerInfoSpec struct {
	// ChassisID refers to the chassis identificator - either MAC-address or system uuid
	//+kubebuilder:validation:Optional
	// validation pattern
	ChassisID string `json:"chassisId,omitempty"`
	// SystemName refers to the advertised peer's name
	//+kubebuilder:validation:Optional
	SystemName string `json:"systemName,omitempty"`
	// PortID refers to the advertised peer's port ID
	//+kubebuilder:validation:Optional
	PortID string `json:"portId,omitempty"`
	// PortDescription refers to the advertised peer's port description
	//+kubebuilder:validation:Optional
	PortDescription string `json:"portDescription,omitempty"`
	// Type refers to the peer type
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=machine;switch;router;undefined
	Type string `json:"type,omitempty"`
}

// SubnetSpec defines switch's subnet info
// +kubebuilder:object:generate=true
type SubnetSpec struct {
	// Contains information to locate the referenced object
	//+kubebuilder:validation:Optional
	*ObjectReference `json:",inline"`
	// CIDR refers to subnet CIDR
	//+kubebuilder:validation:Optional
	// validation pattern
	CIDR string `json:"cidr,omitempty"`
	// Region refers to switch's region
	//+kubebuilder:validation:Optional
	Region *RegionSpec `json:"region,omitempty"`
}

// RegionSpec defines region info
// +kubebuilder:object:generate=true
type RegionSpec struct {
	// Name refers to the switch's region
	//+kubebuilder:validation:Pattern=^[a-z0-9]([-./a-z0-9]*[a-z0-9])?$
	//+kubebuilder:validation:Required
	Name string `json:"name"`
	// AvailabilityZone refers to the switch's availability zone
	//+kubebuilder:validation:Required
	AvailabilityZone string `json:"availabilityZone"`
}

// IPAddressSpec defines interface's ip address info
// +kubebuilder:object:generate=true
type IPAddressSpec struct {
	// Contains information to locate the referenced object
	//+kubebuilder:validation:Optional
	*ObjectReference `json:",inline"`
	// Address refers to the ip address value
	//+kubebuilder:validation:Optional
	Address string `json:"address,omitempty"`
	// ExtraAddress is a flag defining whether address was added as additional by user
	//+kubebuilder:validation:Optional
	//+kubebuilder:default=false
	ExtraAddress bool `json:"extraAddress,omitempty"`
}

// SwitchStateSpec contains current Switch object state.
type SwitchStateSpec struct {
	// State is the current state of corresponding object or process
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=initial;processing;ready;invalid
	State *string `json:"state,omitempty"`
	// Message contains a brief description of the current state
	//+kubebuilder:validation:Optional
	Message *string `json:"message,omitempty"`
}

// ConfigAgentStateSpec contains current configuration agent's state
// +kubebuilder:object:generate=true
type ConfigAgentStateSpec struct {
	// State is the current state of corresponding object or process
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=active;failed
	State *string `json:"state,omitempty"`
	// Message contains a brief description of the current state
	//+kubebuilder:validation:Optional
	Message *string `json:"message,omitempty"`
	// LastCheck refers to the last timestamp when configuration was applied
	//+kubebuilder:validation:Optional
	LastCheck *string `json:"lastCheck,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=sw
//+kubebuilder:storageversion
//+kubebuilder:printcolumn:name="Ports",type=integer,JSONPath=`.status.switchPorts`,description="Total amount of non-management network interfaces"
//+kubebuilder:printcolumn:name="Role",type=string,JSONPath=`.status.role`,description="switch's role"
//+kubebuilder:printcolumn:name="Connection Level",type=integer,JSONPath=`.status.connectionLevel`,description="Vertical level of switch connection"
//+kubebuilder:printcolumn:name="Switch State",type=string,JSONPath=`.status.switch.state`,description="Switch state"
//+kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.switch.message`,description="Switch state message. Reports about any issues duiring reconciliation process"

// Switch is the Schema for switches API.
type Switch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwitchSpec   `json:"spec"`
	Status SwitchStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SwitchList contains a list of Switch.
type SwitchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Switch `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Switch{}, &SwitchList{})
}

// GetNamespacedName returns object's name and namespace as types.NamespacedName.
func (in *Switch) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: in.Namespace,
		Name:      in.Name,
	}
}

func (in *Switch) SetState(state string) {
	in.Status.SwitchState.State = MetalAPIString(state)
}

func (in *Switch) SetRole() {
	in.Status.Role = MetalAPIString(CSwitchRoleSpine)
	for _, data := range in.Status.Interfaces {
		if data.Peer == nil {
			continue
		}
		if data.Peer.Type == CPeerTypeMachine {
			in.Status.Role = MetalAPIString(CSwitchRoleLeaf)
			break
		}
	}
}

func (in *Switch) StateEqualsTo(state string) bool {
	return GoString(in.Status.SwitchState.State) == state
}

func (in *Switch) SetInitialStatus(inv *inventoryv1alpha1.Inventory) {
	in.Status = SwitchStatus{
		TotalPorts:  uint16(len(inv.Spec.NICs)),
		SwitchPorts: 0,
		Role:        nil,
		ConnectionLevel: func() uint8 {
			if in.Spec.TopSpine {
				return 0
			}
			return 255
		}(),
		Interfaces:        nil,
		Subnets:           nil,
		LoopbackAddresses: nil,
		SwitchState: &SwitchStateSpec{
			State:   MetalAPIString(CSwitchStateInitial),
			Message: nil,
		},
		ConfigAgent: nil,
	}
	in.Status.Interfaces = interfacesFromInventory(inv)
	in.Status.SwitchPorts = uint16(len(in.Status.Interfaces))
}

func (in *Switch) UpdateInterfacesParameters(conf *SwitchConfig, list *SwitchList) {
	var (
		resultFEC, resultState *string
		resultLanes            *uint8
		resultMTU              *uint16
	)
	switchesMap := make(map[string]SwitchStatus)
	for _, item := range list.Items {
		switchesMap[item.Name] = item.Status
	}

	if conf != nil {
		if conf.Spec.PortsDefaults != nil {
			if conf.Spec.PortsDefaults.State != nil {
				resultState = conf.Spec.PortsDefaults.State
			}
			if conf.Spec.PortsDefaults.FEC != nil {
				resultFEC = conf.Spec.PortsDefaults.FEC
			}
			if conf.Spec.PortsDefaults.MTU != nil {
				resultMTU = conf.Spec.PortsDefaults.MTU
			}
			if conf.Spec.PortsDefaults.Lanes != nil {
				resultLanes = conf.Spec.PortsDefaults.Lanes
			}
		}
	}
	if in.Spec.Interfaces != nil {
		if in.Spec.Interfaces.Defaults != nil {
			if in.Spec.Interfaces.Defaults.State != nil {
				resultState = in.Spec.Interfaces.Defaults.State
			}
			if in.Spec.Interfaces.Defaults.FEC != nil {
				resultFEC = in.Spec.Interfaces.Defaults.FEC
			}
			if in.Spec.Interfaces.Defaults.MTU != nil {
				resultMTU = in.Spec.Interfaces.Defaults.MTU
			}
			if in.Spec.Interfaces.Defaults.Lanes != nil {
				resultLanes = in.Spec.Interfaces.Defaults.Lanes
			}
		}
	}
	overridden := map[string]struct{}{}
	if in.Spec.Interfaces != nil {
		if in.Spec.Interfaces.Overrides != nil {
			for _, nic := range in.Spec.Interfaces.Overrides {
				stored, ok := in.Status.Interfaces[nic.Name]
				if !ok {
					continue
				}
				overridden[nic.Name] = struct{}{}
				if nic.State != nil {
					stored.State = GoString(nic.State)
				}
				if nic.FEC != nil {
					stored.FEC = GoString(nic.FEC)
				}
				if nic.MTU != nil {
					stored.MTU = GoUint16(nic.MTU)
				}
				if nic.Lanes != nil {
					stored.Lanes = GoUint8(nic.Lanes)
				}
			}
		}
	}
	for nic, params := range in.Status.Interfaces {
		if _, ok := overridden[nic]; ok {
			continue
		}
		if params.Direction == CDirectionNorth {
			peerNICs, ok := switchesMap[params.Peer.Name]
			if !ok {
				continue
			}
			peerNIC, ok := peerNICs.Interfaces[params.Peer.PortDescription]
			params.FEC = peerNIC.FEC
			params.MTU = peerNIC.MTU
			params.Lanes = peerNIC.Lanes
			continue
		}
		if resultState != nil {
			params.State = GoString(resultState)
		}
		if resultFEC != nil {
			params.FEC = GoString(resultFEC)
		}
		if resultMTU != nil {
			params.MTU = GoUint16(resultMTU)
		}
		if resultLanes != nil {
			params.Lanes = GoUint8(resultLanes)
		}
	}
}

func (in *Switch) InterfacesMatchInventory(inv *inventoryv1alpha1.Inventory) bool {
	interfaces := interfacesFromInventory(inv)
	if len(interfaces) != len(in.Status.Interfaces) {
		return false
	}
	for name := range interfaces {
		_, ok := in.Status.Interfaces[name]
		if !ok {
			return false
		}
	}
	for name := range in.Status.Interfaces {
		_, ok := interfaces[name]
		if !ok {
			return false
		}
	}
	return true
}

func (in *Switch) UpdatePeers(inv *inventoryv1alpha1.Inventory) {
	interfaces := interfacesFromInventory(inv)
	for nic, nicData := range interfaces {
		switchPort := in.Status.Interfaces[nic]
		switchPort.Peer = nicData.Peer.DeepCopy()
	}
}

func interfacesFromInventory(inv *inventoryv1alpha1.Inventory) map[string]*InterfaceSpec {
	interfaces := make(map[string]*InterfaceSpec)
	for _, nic := range inv.Spec.NICs {
		if !strings.HasPrefix(nic.Name, CSwitchPortPrefix) {
			continue
		}
		inf := &InterfaceSpec{
			MACAddress: nic.MACAddress,
			FEC:        nic.ActiveFEC,
			MTU:        nic.MTU,
			Speed:      nic.Speed,
			Lanes:      nic.Lanes,
			State:      CNICUp,
			IP:         nil,
			Direction:  CDirectionSouth,
			Peer:       nil,
		}
		for _, data := range nic.LLDPs {
			var emptyLLDP inventoryv1alpha1.LLDPSpec
			if reflect.DeepEqual(data, emptyLLDP) {
				continue
			}
			inf.Peer = &PeerSpec{
				nil,
				&PeerInfoSpec{
					ChassisID:       data.ChassisID,
					SystemName:      data.SystemName,
					PortID:          data.PortID,
					PortDescription: data.PortDescription,
					Type: func() string {
						if len(data.Capabilities) == 0 {
							return CPeerTypeMachine
						}
						for _, c := range data.Capabilities {
							if c == CStationCapability {
								return CPeerTypeMachine
							}
						}
						return CPeerTypeSwitch
					}(),
				},
			}
			break
		}
		interfaces[nic.Name] = inf
	}
	return interfaces
}

func (in *Switch) ConnectionsOK(list *SwitchList) bool {
	return in.peersOK(list) && in.connectionLevelOK(list)
}

func (in *Switch) peersOK(list *SwitchList) bool {
	for _, item := range list.Items {
		if item.Name == in.Name {
			continue
		}
		for _, nicData := range item.Status.Interfaces {
			if item.Status.ConnectionLevel == 0 && nicData.Direction == CDirectionNorth {
				return false
			}
			if nicData.Peer == nil {
				continue
			}
			if reflect.DeepEqual(nicData.Peer.PeerInfoSpec, &PeerInfoSpec{}) {
				continue
			}
			if strings.ReplaceAll(nicData.Peer.PeerInfoSpec.ChassisID, ":", "") != in.Annotations[CHardwareChassisIDAnnotation] {
				continue
			}
			if nicData.Peer.PeerInfoSpec.PortDescription == "" {
				continue
			}
			nic, ok := in.Status.Interfaces[nicData.Peer.PeerInfoSpec.PortDescription]
			if !ok {
				nic, ok = in.Status.Interfaces[nicData.Peer.PeerInfoSpec.PortID]
				if !ok {
					return false
				}
			}
			if nic.Peer == nil {
				return false
			}
			if reflect.DeepEqual(nic.Peer.PeerInfoSpec, &PeerInfoSpec{}) {
				return false
			}
			if strings.ReplaceAll(nic.Peer.PeerInfoSpec.ChassisID, ":", "") != item.Annotations[CHardwareChassisIDAnnotation] {
				return false
			}
			if nic.Peer.ObjectReference == nil {
				return false
			}
			if !(nic.Peer.ObjectReference.Name == item.Name) || !(nic.Peer.ObjectReference.Namespace == item.Namespace) {
				return false
			}
		}
	}
	return true
}

func (in *Switch) connectionLevelOK(list *SwitchList) bool {
	if in.Status.ConnectionLevel == 255 {
		return false
	}
	if in.Spec.TopSpine && in.Status.ConnectionLevel != 0 {
		return false
	}
	if !in.Spec.TopSpine && in.Status.ConnectionLevel == 0 {
		return false
	}
	for _, item := range list.Items {
		for _, nicData := range in.Status.Interfaces {
			if in.Status.ConnectionLevel == 0 && nicData.Direction == CDirectionNorth {
				return false
			}
			if nicData.Peer == nil {
				continue
			}
			if nicData.Peer.ObjectReference == nil {
				continue
			}
			if nicData.Peer.ObjectReference.Name != item.Name {
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

func (in *Switch) SetConnections(list *SwitchList) {
	in.fillPeersInfo(list)
	in.computeConnectionLevel(list)
}

func (in *Switch) fillPeersInfo(list *SwitchList) {
	for _, item := range list.Items {
		for _, nicData := range in.Status.Interfaces {
			if nicData.Peer == nil {
				continue
			}
			if nicData.Peer.PeerInfoSpec == nil {
				continue
			}
			if reflect.DeepEqual(nicData.Peer.PeerInfoSpec, &PeerInfoSpec{}) {
				continue
			}
			if strings.ReplaceAll(nicData.Peer.PeerInfoSpec.ChassisID, ":", "") != item.Annotations[CHardwareChassisIDAnnotation] {
				continue
			}
			nicData.Peer.ObjectReference = &ObjectReference{
				Name:      item.Name,
				Namespace: item.Namespace,
			}
		}
	}
}

func (in *Switch) computeConnectionLevel(list *SwitchList) {
	connectionsMap, keys := list.buildConnectionMap()
	if _, ok := connectionsMap[0]; !ok {
		return
	}

	switch in.Spec.TopSpine {
	case true:
		in.Status.ConnectionLevel = 0
		for _, nicData := range in.Status.Interfaces {
			nicData.Direction = CDirectionSouth
		}
		return
	case false:
		if in.Status.ConnectionLevel != 0 {
			break
		}
		in.Status.ConnectionLevel = 255
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
		in.setNICsDirections(list)
	}
}

func (in *SwitchList) buildConnectionMap() (map[uint8]*SwitchList, []uint8) {
	connectionsMap := make(map[uint8]*SwitchList)
	keys := make([]uint8, 0)
	for _, item := range in.Items {
		if item.Status.SwitchState == nil {
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

func (in *Switch) getPeers(list *SwitchList) *SwitchList {
	result := &SwitchList{Items: make([]Switch, 0)}
	for _, item := range list.Items {
		for _, nicData := range in.Status.Interfaces {
			if nicData.Peer == nil {
				continue
			}
			if nicData.Peer.PeerInfoSpec == nil {
				continue
			}
			if reflect.DeepEqual(nicData.Peer.PeerInfoSpec, &PeerInfoSpec{}) {
				continue
			}
			if strings.ReplaceAll(nicData.Peer.PeerInfoSpec.ChassisID, ":", "") == item.Annotations[CHardwareChassisIDAnnotation] {
				result.Items = append(result.Items, item)
			}
		}
	}
	return result
}

func (in *Switch) setNICsDirections(list *SwitchList) {
	if in.Status.ConnectionLevel == 0 {
		for _, nicData := range in.Status.Interfaces {
			nicData.Direction = CDirectionSouth
		}
		return
	}
	for _, item := range list.Items {
		for _, nicData := range in.Status.Interfaces {
			if nicData.Peer == nil {
				nicData.Direction = CDirectionSouth
				continue
			}
			if nicData.Peer.ObjectReference == nil {
				nicData.Direction = CDirectionSouth
				continue
			}
			if reflect.DeepEqual(nicData.Peer.PeerInfoSpec, &PeerInfoSpec{}) {
				nicData.Direction = CDirectionSouth
				continue
			}
			peerFound := strings.ReplaceAll(nicData.Peer.PeerInfoSpec.ChassisID, ":", "") == item.Annotations[CHardwareChassisIDAnnotation]
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

func (in *Switch) SubnetSelectorsExist() bool {
	if in.Spec.IPAM == nil {
		return false
	}
	if in.Spec.IPAM.SouthSubnets == nil {
		return false
	}
	if in.Spec.IPAM.SouthSubnets.LabelSelector == nil && in.Spec.IPAM.SouthSubnets.FieldSelector == nil {
		return false
	}
	return true
}

func (in *Switch) LoopbackSelectorsExist() bool {
	if in.Spec.IPAM == nil {
		return false
	}
	if in.Spec.IPAM.LoopbackAddresses == nil {
		return false
	}
	if in.Spec.IPAM.LoopbackAddresses.LabelSelector == nil && in.Spec.IPAM.LoopbackAddresses.FieldSelector == nil {
		return false
	}
	return true
}

func (in *Switch) GetAddressesCount(bits uint8, af ipamv1alpha1.SubnetAddressType) int64 {
	var addressesCount int64
	addressesPerPort := int64(math.Pow(float64(2), float64(CIPv4MaskLengthBits-bits)))
	if af == ipamv1alpha1.CIPv6SubnetType {
		addressesPerPort = int64(math.Pow(float64(2), float64(CIPv6PrefixBits-bits)))
	}
	for _, nic := range in.Status.Interfaces {
		addressesCount += addressesPerPort * int64(nic.Lanes)
	}
	return addressesCount
}

func (in *Switch) SouthSubnetsAFUsage(cfg *SwitchConfig) (v4used bool, v6used bool) {
	if in.Spec.IPAM == nil {
		if cfg == nil {
			return
		}
		if cfg.Spec.IPAM.SouthSubnets == nil {
			return true, true
		}
		if cfg.Spec.IPAM.SouthSubnets.AddressFamilies == nil {
			return true, true
		}
		v4used = cfg.Spec.IPAM.SouthSubnets.AddressFamilies.IPv4
		v6used = cfg.Spec.IPAM.SouthSubnets.AddressFamilies.IPv6
		return
	}
	if in.Spec.IPAM.SouthSubnets == nil {
		return true, true
	}
	if in.Spec.IPAM.SouthSubnets.AddressFamilies == nil {
		return true, true
	}
	v4used = in.Spec.IPAM.SouthSubnets.AddressFamilies.IPv4
	v6used = in.Spec.IPAM.SouthSubnets.AddressFamilies.IPv6
	return
}

func (in *Switch) LoopbacksAFUsage(cfg *SwitchConfig) (v4used bool, v6used bool) {
	if in.Spec.IPAM == nil {
		if cfg == nil {
			return
		}
		if cfg.Spec.IPAM.LoopbackAddresses == nil {
			return true, true
		}
		if cfg.Spec.IPAM.LoopbackAddresses.AddressFamilies == nil {
			return true, true
		}
		v4used = cfg.Spec.IPAM.LoopbackAddresses.AddressFamilies.IPv4
		v6used = cfg.Spec.IPAM.LoopbackAddresses.AddressFamilies.IPv6
		return
	}
	if in.Spec.IPAM.LoopbackAddresses == nil {
		return true, true
	}
	if in.Spec.IPAM.LoopbackAddresses.AddressFamilies == nil {
		return true, true
	}
	v4used = in.Spec.IPAM.LoopbackAddresses.AddressFamilies.IPv4
	v6used = in.Spec.IPAM.LoopbackAddresses.AddressFamilies.IPv6
	return
}

func (in *Switch) ResultingIPAMConfig(cfg *SwitchConfig) error {
	//todo: looks ugly from my perspective. Need to rethink and rewrite
	ipamSelectorsResult := &IPAMSpec{}
	portDefaultsResult := &PortParametersSpec{}
	if cfg != nil {
		ipamSelectorsResult = &IPAMSpec{
			SouthSubnets: &IPAMSelectionSpec{
				AddressFamilies: &AddressFamiliesMap{
					IPv4: cfg.Spec.IPAM.SouthSubnets.AddressFamilies.IPv4,
					IPv6: cfg.Spec.IPAM.SouthSubnets.AddressFamilies.IPv6,
				},
				LabelSelector: cfg.Spec.IPAM.SouthSubnets.LabelSelector.DeepCopy(),
			},
			LoopbackAddresses: &IPAMSelectionSpec{
				AddressFamilies: &AddressFamiliesMap{
					IPv4: cfg.Spec.IPAM.LoopbackAddresses.AddressFamilies.IPv4,
					IPv6: cfg.Spec.IPAM.LoopbackAddresses.AddressFamilies.IPv6,
				},
				LabelSelector: cfg.Spec.IPAM.LoopbackAddresses.LabelSelector.DeepCopy(),
			},
		}
		southSubnetLabelFromFiledRef, err := LabelFromFieldRef(*in, cfg.Spec.IPAM.SouthSubnets.FieldSelector)
		if err != nil {
			return err
		}
		for k, v := range southSubnetLabelFromFiledRef {
			ipamSelectorsResult.SouthSubnets.LabelSelector.MatchLabels[k] = v
		}
		loopbacksLabelFromFieldRef, err := LabelFromFieldRef(*in, cfg.Spec.IPAM.LoopbackAddresses.FieldSelector)
		if err != nil {
			return err
		}
		for k, v := range loopbacksLabelFromFieldRef {
			ipamSelectorsResult.LoopbackAddresses.LabelSelector.MatchLabels[k] = v
		}
		portDefaultsResult = cfg.Spec.PortsDefaults.DeepCopy()
	}

	if in.Spec.IPAM == nil {
		in.Spec.IPAM = ipamSelectorsResult.DeepCopy()
	}
	if in.Spec.IPAM.SouthSubnets == nil {
		in.Spec.IPAM.SouthSubnets = ipamSelectorsResult.SouthSubnets
	}
	if in.Spec.IPAM.LoopbackAddresses == nil {
		in.Spec.IPAM.LoopbackAddresses = ipamSelectorsResult.LoopbackAddresses
	}

	if in.Spec.Interfaces == nil {
		in.Spec.Interfaces = &InterfacesSpec{
			Defaults: portDefaultsResult.DeepCopy(),
		}
	}
	if in.Spec.Interfaces.Defaults.IPv4MaskLength == nil {
		in.Spec.Interfaces.Defaults.IPv4MaskLength = portDefaultsResult.IPv4MaskLength
	}
	if in.Spec.Interfaces.Defaults.IPv6Prefix == nil {
		in.Spec.Interfaces.Defaults.IPv6Prefix = portDefaultsResult.IPv6Prefix
	}
	return nil
}

func (in *Switch) LoopbackIPsMatchStoredIPs(list *ipamv1alpha1.IPList) bool {
	v4used := in.Spec.IPAM.LoopbackAddresses.AddressFamilies.IPv4
	v6used := in.Spec.IPAM.LoopbackAddresses.AddressFamilies.IPv6

receivedAddressesLoop:
	for _, item := range list.Items {
		if item.Status.State != ipamv1alpha1.CFinishedIPState {
			continue
		}
		if !v4used && item.Status.Reserved.Net.Is4() {
			continue
		}
		if !v6used && item.Status.Reserved.Net.Is6() {
			continue
		}
		for _, lo := range in.Status.LoopbackAddresses {
			if item.Status.Reserved.String() == lo.Address {
				continue receivedAddressesLoop
			}
		}
		return false
	}

storedAddressesLoop:
	for _, lo := range in.Status.LoopbackAddresses {
		for _, item := range list.Items {
			if item.Status.Reserved.String() == lo.Address {
				continue storedAddressesLoop
			}
		}
		return false
	}

	return true
}

func (in *Switch) SubnetsMatchStored(list *ipamv1alpha1.SubnetList) bool {
	v4used := in.Spec.IPAM.SouthSubnets.AddressFamilies.IPv4
	v6used := in.Spec.IPAM.SouthSubnets.AddressFamilies.IPv6

receivedSubnetsLoop:
	for _, item := range list.Items {
		if item.Status.State != ipamv1alpha1.CFinishedSubnetState {
			continue
		}
		if !v4used && item.Status.Reserved.IsIPv4() {
			continue
		}
		if !v6used && item.Status.Reserved.IsIPv6() {
			continue
		}
		for _, sn := range in.Status.Subnets {
			if item.Status.Reserved.String() == sn.CIDR {
				continue receivedSubnetsLoop
			}
		}
		return false
	}

storedSubnetsLoop:
	for _, sn := range in.Status.Subnets {
		for _, item := range list.Items {
			if item.Status.Reserved.String() == sn.CIDR {
				continue storedSubnetsLoop
			}
		}
		return false
	}
	return true
}

func (in *Switch) IPaddressesOK(list *SwitchList) bool {
	return in.ipsMatchSubnets() && in.ipsMatchPeers(list) && in.ipsMatchOverrides()
}

func (in *Switch) ipsMatchSubnets() bool {
	for _, nic := range in.Status.Interfaces {
		if nic.Direction == CDirectionNorth {
			continue
		}
		if len(in.Status.Subnets) != 0 && len(nic.IP) == 0 {
			return false
		}
	ipsLoop:
		for _, ip := range nic.IP {
			if ip.ExtraAddress {
				continue
			}
			for _, subnet := range in.Status.Subnets {
				_, cidr, _ := net.ParseCIDR(subnet.CIDR)
				addr, _, _ := net.ParseCIDR(ip.Address)
				if cidr.Contains(addr) {
					continue ipsLoop
				}
			}
			return false
		}
	}
	return true
}

func (in *Switch) ipsMatchPeers(list *SwitchList) bool {
	for nic, nicData := range in.Status.Interfaces {
		if !strings.HasPrefix(nic, CSwitchPortPrefix) {
			continue
		}
		if nicData.Direction == CDirectionSouth {
			continue
		}
		for _, item := range list.Items {
			if nicData.Peer.Name != item.Name {
				continue
			}
		peerSubnetsLoop:
			for _, subnet := range item.Status.Subnets {
				for _, ip := range nicData.IP {
					_, cidr, _ := net.ParseCIDR(subnet.CIDR)
					addr, _, _ := net.ParseCIDR(ip.Address)
					if cidr.Contains(addr) {
						continue peerSubnetsLoop
					}
				}
				return false
			}
		}
	}
	return true
}

func (in *Switch) ipsMatchOverrides() bool {
	if in.Spec.Interfaces == nil {
		return true
	}
	if in.Spec.Interfaces.Overrides == nil {
		return true
	}
overriddenNICsLoop:
	for _, override := range in.Spec.Interfaces.Overrides {
		if override.IP == nil {
			continue
		}
	extraIPsLoop:
		for _, ip := range override.IP {
			stored, ok := in.Status.Interfaces[override.Name]
			if !ok {
				continue overriddenNICsLoop
			}
			for _, storedIP := range stored.IP {
				if storedIP.Address == ip.Address {
					continue extraIPsLoop
				}
			}
			return false
		}
	}
	return true
}

func (in *Switch) GetExtraNICsIPs() map[string][]*IPAddressSpec {
	ipsToApply := make(map[string][]*IPAddressSpec)
	if in.Spec.Interfaces != nil && in.Spec.Interfaces.Overrides != nil {
		for _, nic := range in.Spec.Interfaces.Overrides {
			_, ok := in.Status.Interfaces[nic.Name]
			if !ok {
				continue
			}
			nicIPs := make([]*IPAddressSpec, 0)
			for _, ip := range nic.IP {
				nicIPs = append(nicIPs, &IPAddressSpec{
					Address:      ip.Address,
					ExtraAddress: true,
				})
			}
			ipsToApply[nic.Name] = nicIPs
		}
	}
	return ipsToApply
}

func (in *Switch) GetSouthNICsIP() (map[string][]*IPAddressSpec, error) {
	ipsToApply := make(map[string][]*IPAddressSpec)
	for nic, nicData := range in.Status.Interfaces {
		if !strings.HasPrefix(nic, CSwitchPortPrefix) {
			continue
		}
		if nicData.Direction == CDirectionNorth {
			continue
		}
		nicIPs := make([]*IPAddressSpec, 0)
		for _, subnet := range in.Status.Subnets {
			cidr, _ := ipamv1alpha1.CIDRFromString(subnet.CIDR)
			mask := GoUint8(in.Spec.Interfaces.Defaults.IPv4MaskLength)
			addrIndex := 1
			if cidr.IsIPv6() {
				mask = GoUint8(in.Spec.Interfaces.Defaults.IPv6Prefix)
				addrIndex = 0
			}
			nicSubnet := getInterfaceSubnet(nic, CSwitchPortPrefix, cidr.Net.IPNet(), mask)
			nicAddr, err := gocidr.Host(nicSubnet, addrIndex)
			if err != nil {
				return nil, err
			}
			nicIPs = append(nicIPs, &IPAddressSpec{
				Address:      fmt.Sprintf("%s/%d", nicAddr.String(), mask),
				ExtraAddress: false,
			})
		}
		ipsToApply[nic] = nicIPs
	}
	return ipsToApply, nil
}

func (in *Switch) GetNorthNICsIP(list *SwitchList) map[string][]*IPAddressSpec {
	ipsToApply := make(map[string][]*IPAddressSpec)
	for nic, nicData := range in.Status.Interfaces {
		if !strings.HasPrefix(nic, CSwitchPortPrefix) {
			continue
		}
		if nicData.Direction == CDirectionSouth {
			continue
		}
		nicIPs := make([]*IPAddressSpec, 0)
		for _, item := range list.Items {
			if nicData.Peer.Name != item.Name {
				continue
			}
			peerNICdata, ok := item.Status.Interfaces[nicData.Peer.PortDescription]
			if !ok {
				peerNICdata, ok = item.Status.Interfaces[nicData.Peer.PortID]
				if !ok {
					continue
				}
			}
			nicIPs = append(nicIPs, func() []*IPAddressSpec {
				requestedAddresses := make([]*IPAddressSpec, 0)
				for _, ip := range peerNICdata.RequestAddress() {
					requestedAddresses = append(requestedAddresses, &IPAddressSpec{
						Address:      ip.String(),
						ExtraAddress: false,
					})
				}
				return requestedAddresses
			}()...)
		}
		ipsToApply[nic] = nicIPs
	}
	return ipsToApply
}

func getInterfaceSubnet(name string, namePrefix string, network *net.IPNet, mask uint8) *net.IPNet {
	index, _ := strconv.Atoi(strings.ReplaceAll(name, namePrefix, CEmptyString))
	prefix, _ := network.Mask.Size()
	ifaceNet, _ := gocidr.Subnet(network, int(mask)-prefix, index)
	return ifaceNet
}

func (in *InterfaceSpec) RequestAddress() (ips []net.IPNet) {
	ips = make([]net.IPNet, 0)
	for _, addr := range in.IP {
		_, cidr, _ := net.ParseCIDR(addr.Address)
		ip, _ := gocidr.Host(cidr, 1)
		ips = append(ips, net.IPNet{IP: ip, Mask: cidr.Mask})
	}
	return
}

func (in *Switch) LabelsOK() bool {
	if in.Labels == nil {
		return false
	}
	if _, ok := in.Labels[InventoriedLabel]; !ok {
		return false
	}
	if _, ok := in.Labels[InventoryRefLabel]; !ok {
		return false
	}
	return true
}

func (in *Switch) UpdateSwitchLabels(inv *inventoryv1alpha1.Inventory) {
	appliedLabels := map[string]string{
		InventoriedLabel:  "true",
		InventoryRefLabel: inv.Name,
		LabelChassisID: strings.ReplaceAll(
			func() string {
				var chassisID string
				for _, nic := range inv.Spec.NICs {
					if nic.Name == "eth0" {
						chassisID = nic.MACAddress
					}
				}
				return chassisID
			}(), ":", "-",
		),
	}
	if in.Labels == nil {
		in.Labels = make(map[string]string)
	}
	for k, v := range appliedLabels {
		in.Labels[k] = v
	}
}

func (in *Switch) UpdateSwitchAnnotations(inv *inventoryv1alpha1.Inventory) {
	hardwareAnnotations := make(map[string]string)
	softwareAnnotations := make(map[string]string)
	if inv.Spec.System != nil {
		hardwareAnnotations[CHardwareSerialAnnotation] = inv.Spec.System.SerialNumber
		hardwareAnnotations[CHardwareManufacturerAnnotation] = inv.Spec.System.Manufacturer
		hardwareAnnotations[CHardwareSkuAnnotation] = inv.Spec.System.ProductSKU
	}
	if inv.Spec.Distro != nil {
		softwareAnnotations[CSoftwareOnieAnnotation] = "false"
		softwareAnnotations[CSoftwareAsicAnnotation] = inv.Spec.Distro.AsicType
		softwareAnnotations[CSoftwareVersionAnnotation] = inv.Spec.Distro.CommitID
		softwareAnnotations[CSoftwareOSAnnotation] = "sonic"
		softwareAnnotations[CSoftwareHostnameAnnotation] = inv.Spec.Host.Name
	}
	if in.Annotations == nil {
		in.Annotations = make(map[string]string)
	}
	in.Annotations[CHardwareChassisIDAnnotation] = strings.ReplaceAll(
		func() string {
			var chassisID string
			for _, nic := range inv.Spec.NICs {
				if nic.Name == "eth0" {
					chassisID = nic.MACAddress
				}
			}
			return chassisID
		}(), ":", "",
	)
	for k, v := range hardwareAnnotations {
		in.Annotations[k] = v
	}
	for k, v := range softwareAnnotations {
		in.Annotations[k] = v
	}
}
