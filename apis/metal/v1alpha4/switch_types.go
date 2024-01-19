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

package v1alpha4

import (
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"

	"github.com/ironcore-dev/metal/pkg/constants"
)

// NetworkSwitchSpec contains desired state of resulting NetworkSwitch configuration.
type NetworkSwitchSpec struct {
	// InventoryRef contains reference to corresponding inventory object
	// Empty InventoryRef means that there is no corresponding Inventory object
	// +kubebuilder:validation:Optional
	InventoryRef *v1.LocalObjectReference `json:"inventoryRef,omitempty"`
	// ConfigSelector contains selector to filter out corresponding SwitchConfig.
	// If the selector is not defined, it will be populated by defaulting webhook
	// with MatchLabels item, containing 'metal.ironcore.dev/layer' key with value
	// equals to object's .status.layer.
	// +kubebuilder:validation:Optional
	ConfigSelector *metav1.LabelSelector `json:"configSelector,omitempty"`
	// Managed is a flag defining whether NetworkSwitch object would be processed during reconciliation
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	Managed *bool `json:"managed"`
	// Cordon is a flag defining whether NetworkSwitch object is taken offline
	// +kubebuilder:validation:Required
	// +kubebuilder:default=false
	Cordon *bool `json:"cordon"`
	// TopSpine is a flag defining whether NetworkSwitch is a top-level spine switch
	// +kubebuilder:validation:Required
	// +kubebuilder:default=false
	TopSpine *bool `json:"topSpine"`
	// ScanPorts is a flag defining whether to run periodical scanning on switch ports
	// +kubebuilder:validation:Required
	// +kubebuilder:default=true
	ScanPorts *bool `json:"scanPorts"`
	// IPAM refers to selectors for subnets which will be used for NetworkSwitch object
	// +kubebuilder:validation:Optional
	IPAM *IPAMSpec `json:"ipam,omitempty"`
	// Interfaces contains general configuration for all switch ports
	// +kubebuilder:validation:Optional
	Interfaces *InterfacesSpec `json:"interfaces,omitempty"`
}

// InterfacesSpec contains definitions for general switch ports' configuration.
type InterfacesSpec struct {
	// Defaults contains switch port parameters which will be applied to all ports of the switches
	// +kubebuilder:validation:Optional
	Defaults *PortParametersSpec `json:"defaults,omitempty"`
	// Overrides contains set of parameters which should be overridden for listed switch ports
	// +kubebuilder:validation:Optional
	Overrides []*InterfaceOverridesSpec `json:"overrides,omitempty"`
}

// InterfaceOverridesSpec contains overridden parameters for certain switch port.
type InterfaceOverridesSpec struct {
	// Contains port parameters overrides
	// +kubebuilder:validation:Required
	*PortParametersSpec `json:",inline"`
	// Name refers to switch port name
	// +kubebuilder:validation:Optional
	Name *string `json:"name,omitempty"`
	// IP contains a list of additional IP addresses for interface
	// +kubebuilder:validation:Optional
	IP []*AdditionalIPSpec `json:"ip,omitempty"`
}

// NetworkSwitchStatus contains observed state of NetworkSwitch.
type NetworkSwitchStatus struct {
	// ConfigRef contains reference to corresponding SwitchConfig object
	// Empty ConfigRef means that there is no corresponding SwitchConfig object
	// +kubebuilder:validation:Optional
	ConfigRef v1.LocalObjectReference `json:"configRef,omitempty"`
	// RoutingConfigTemplate contains the reference to the ConfigMap object which contains go-template for FRR config.
	// This field reflects the corresponding field of the related SwitchConfig object.
	// +kubebuilder:validation:Optional
	RoutingConfigTemplate v1.LocalObjectReference `json:"routingConfigTemplate,omitempty"`
	// ASN contains current autonomous system number defined for switch
	// +kubebuilder:validation:Optional
	ASN uint32 `json:"asn,omitempty"`
	// TotalPorts refers to total number of ports
	// +kubebuilder:validation:Optional
	TotalPorts uint32 `json:"totalPorts,omitempty"`
	// SwitchPorts refers to the number of ports excluding management interfaces, loopback etc.
	// +kubebuilder:validation:Optional
	SwitchPorts uint32 `json:"switchPorts,omitempty"`
	// Role refers to switch's role
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=spine;leaf;edge-leaf
	Role string `json:"role,omitempty"`
	// Layer refers to switch's current position in connection hierarchy
	// +kubebuilder:validation:Optional
	Layer uint32 `json:"layer"`
	// Interfaces refers to switch's interfaces configuration
	// +kubebuilder:validation:Optional
	Interfaces map[string]*InterfaceSpec `json:"interfaces,omitempty"`
	// Subnets refers to the switch's south subnets
	// +kubebuilder:validation:Optional
	Subnets []*SubnetSpec `json:"subnets,omitempty"`
	// LoopbackAddresses refers to the switch's loopback addresses
	// +kubebuilder:validation:Optional
	LoopbackAddresses []*IPAddressSpec `json:"loopbackAddresses,omitempty"`
	// State is the current state of corresponding object or process
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=Initial;Processing;Ready;Invalid;Pending
	State *string `json:"state,omitempty"`
	// Message contains a brief description of the current state
	// +kubebuilder:validation:Optional
	Message *string `json:"message,omitempty"`
	// Condition contains state of port parameters
	// +kubebuilder:validation:Optional
	Conditions []*ConditionSpec `json:"conditions,omitempty"`
}

// InterfaceSpec defines the state of switch's interface.
type InterfaceSpec struct {
	// Contains port parameters
	// +kubebuilder:validation:Required
	*PortParametersSpec `json:",inline"`
	// MACAddress refers to the interface's hardware address
	// +kubebuilder:validation:Required
	// validation pattern
	MACAddress *string `json:"macAddress"`
	// Speed refers to interface's speed
	// +kubebuilder:validation:Required
	Speed *uint32 `json:"speed"`
	// IP contains a list of IP addresses that are assigned to interface
	// +kubebuilder:validation:Optional
	IP []*IPAddressSpec `json:"ip,omitempty"`
	// Direction refers to the interface's connection 'direction'
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=north;south
	Direction *string `json:"direction"`
	// Peer refers to the info about device connected to current switch port
	// +kubebuilder:validation:Optional
	Peer *PeerSpec `json:"peer,omitempty"`
}

// PeerSpec defines peer info.
type PeerSpec struct {
	// Contains information to locate the referenced object
	// +kubebuilder:validation:Optional
	*ObjectReference `json:",inline"`
	// Contains LLDP info about peer
	// +kubebuilder:validation:Optional
	*PeerInfoSpec `json:",inline"`
}

// PeerInfoSpec contains LLDP info about peer.
type PeerInfoSpec struct {
	// ChassisID refers to the chassis identificator - either MAC-address or system uuid
	// +kubebuilder:validation:Optional
	// validation pattern
	ChassisID *string `json:"chassisId,omitempty"`
	// SystemName refers to the advertised peer's name
	// +kubebuilder:validation:Optional
	SystemName *string `json:"systemName,omitempty"`
	// PortID refers to the advertised peer's port ID
	// +kubebuilder:validation:Optional
	PortID *string `json:"portId,omitempty"`
	// PortDescription refers to the advertised peer's port description
	// +kubebuilder:validation:Optional
	PortDescription *string `json:"portDescription,omitempty"`
	// Type refers to the peer type
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=machine;switch;router;undefined
	Type *string `json:"type,omitempty"`
}

// SubnetSpec defines switch's subnet info.
type SubnetSpec struct {
	// Contains information to locate the referenced object
	// +kubebuilder:validation:Optional
	Subnet *ObjectReference `json:"subnet"`
	// Contains information to locate the referenced object
	// +kubebuilder:validation:Optional
	Network *ObjectReference `json:"network"`
	// CIDR refers to subnet CIDR
	// +kubebuilder:validation:Optional
	// validation pattern
	CIDR *string `json:"cidr,omitempty"`
	// Region refers to switch's region
	// +kubebuilder:validation:Optional
	Region *RegionSpec `json:"region,omitempty"`
	// AddressFamily refers to the AF of subnet
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=IPv4;IPv6
	AddressFamily *string `json:"addressFamily,omitempty"`
}

// RegionSpec defines region info.
type RegionSpec struct {
	// Name refers to the switch's region
	// +kubebuilder:validation:Pattern=^[a-z0-9]([-./a-z0-9]*[a-z0-9])?$
	// +kubebuilder:validation:Required
	Name *string `json:"name"`
	// AvailabilityZone refers to the switch's availability zone
	// +kubebuilder:validation:Required
	AvailabilityZone *string `json:"availabilityZone"`
}

// IPAddressSpec defines interface's ip address info.
type IPAddressSpec struct {
	// Contains information to locate the referenced object
	// +kubebuilder:validation:Optional
	*ObjectReference `json:",inline"`
	// Address refers to the ip address value
	// +kubebuilder:validation:Optional
	Address *string `json:"address,omitempty"`
	// ExtraAddress is a flag defining whether address was added as additional by user
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	ExtraAddress *bool `json:"extraAddress,omitempty"`
	// AddressFamily refers to the AF of IP address
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=IPv4;IPv6
	AddressFamily *string `json:"addressFamily,omitempty"`
}

// ConditionSpec contains current condition of port parameters.
type ConditionSpec struct {
	// Name reflects the name of the condition
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=Initialized;InterfacesOK;ConfigRefOK;PortParametersOK;NeighborsOK;LayerAndRoleOK;LoopbacksOK;AsnOK;SubnetsOK;IPAddressesOK
	Name *string `json:"name,omitempty"`
	// State reflects the state of the condition
	// +kubebuilder:validation:Optional
	State *bool `json:"state,omitempty"`
	// LastUpdateTimestamp reflects the last timestamp when condition was updated
	// +kubebuilder:validation:Optional
	LastUpdateTimestamp *string `json:"lastUpdateTimestamp"`
	// LastTransitionTimestamp reflects the last timestamp when condition changed state from one to another
	// +kubebuilder:validation:Optional
	LastTransitionTimestamp *string `json:"lastTransitionTimestamp"`
	// Reason reflects the reason of condition state
	// +kubebuilder:validation:Optional
	Reason *string `json:"reason,omitempty"`
	// Message reflects the verbose message about the reason
	// +kubebuilder:validation:Optional
	Message *string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=nsw
// +kubebuilder:printcolumn:name="ASN",type=integer,JSONPath=`.status.asn`,description="NetworkSwitch ASN"
// +kubebuilder:printcolumn:name="Ports",type=integer,JSONPath=`.status.switchPorts`,description="Number of switch ports"
// +kubebuilder:printcolumn:name="Role",type=string,JSONPath=`.status.role`,description="NetworkSwitch's role"
// +kubebuilder:printcolumn:name="Layer",type=integer,JSONPath=`.status.layer`,description="Vertical level in switches' connections hierarchy"
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`,description="NetworkSwitch state"
// +kubebuilder:printcolumn:name="Message",priority=1,type=string,JSONPath=`.status.message`,description="NetworkSwitch state message reports about any issues during processing"
// +genclient

// NetworkSwitch is the Schema for switches API.
type NetworkSwitch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkSwitchSpec   `json:"spec"`
	Status NetworkSwitchStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NetworkSwitchList contains a list of NetworkSwitch.
type NetworkSwitchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkSwitch `json:"items"`
}

func (in *NetworkSwitch) GetNamespace() string {
	return in.Namespace
}

// NamespacedName returns types.NamespacedName built from
// object's metadata.name and metadata.namespace.
func (in *NetworkSwitch) NamespacedName() types.NamespacedName {
	return types.NamespacedName{Namespace: in.Namespace, Name: in.Name}
}

// ----------------------------------------
// NetworkSwitchSpec getters
// ----------------------------------------

// GetInventoryRef returns value of spec.inventoryRef.name field if
// inventoryRef is not nil, otherwise empty string.
func (in *NetworkSwitch) GetInventoryRef() string {
	if in.Spec.InventoryRef == nil {
		return ""
	}
	return in.Spec.InventoryRef.Name
}

// Managed returns value of spec.managed field if it is not nil,
// otherwise false.
func (in *NetworkSwitch) Managed() bool {
	return ptr.Deref(in.Spec.Managed, false)
}

// Cordon returns value of spec.cordon field if it is not nil,
// otherwise false.
func (in *NetworkSwitch) Cordon() bool {
	return ptr.Deref(in.Spec.Cordon, false)
}

// TopSpine returns value of spec.topSpine field if it is not nil,
// otherwise false.
func (in *NetworkSwitch) TopSpine() bool {
	return ptr.Deref(in.Spec.TopSpine, false)
}

// ScanPorts returns value of spec.topSpine field if it is not nil,
// otherwise false.
func (in *NetworkSwitch) ScanPorts() bool {
	return ptr.Deref(in.Spec.ScanPorts, false)
}

// GetConfigSelector returns LabelSelector if it is not nil and
// matching criteria are defined, otherwise nil.
func (in *NetworkSwitch) GetConfigSelector() *metav1.LabelSelector {
	if in.Spec.ConfigSelector == nil {
		return nil
	}
	if len(in.Spec.ConfigSelector.MatchLabels) == 0 && len(in.Spec.ConfigSelector.MatchExpressions) == 0 {
		return nil
	}
	return in.Spec.ConfigSelector
}

// ----------------------------------------
// NetworkSwitchStatus getters
// ----------------------------------------

// GetConfigRef returns value of status.configRef.name field if
// configRef is not nil, otherwise empty string.
func (in *NetworkSwitch) GetConfigRef() string {
	// if in.Status.ConfigRef == nil {
	//	return ""
	// }
	return in.Status.ConfigRef.Name
}

// GetRoutingConfigTemplate returns value of status.routingConfigTemplate.name field if
// routingConfigTemplate is not nil, otherwise empty string.
func (in *NetworkSwitch) GetRoutingConfigTemplate() string {
	// if in.Status.RoutingConfigTemplate == nil {
	//	return ""
	// }
	return in.Status.RoutingConfigTemplate.Name
}

// RoutingConfigTemplateIsEmpty checks whether the status.routingConfigTemplate contains
// value or not.
func (in *NetworkSwitch) RoutingConfigTemplateIsEmpty() bool {
	return in.GetRoutingConfigTemplate() == constants.EmptyString
}

// GetASN returns value of status.asn field if it is not nil,
// otherwise 0.
func (in *NetworkSwitch) GetASN() uint32 {
	return in.Status.ASN
}

// GetLayer returns value of status.layer field if it is not nil,
// otherwise 255.
func (in *NetworkSwitch) GetLayer() uint32 {
	return in.Status.Layer
}

// GetRole returns value of status.role field if it is not nil,
// otherwise empty string.
func (in *NetworkSwitch) GetRole() string {
	return in.Status.Role
}

// GetTotalPorts returns value of status.totalPorts field if it is not nil,
// otherwise 0.
func (in *NetworkSwitch) GetTotalPorts() uint32 {
	return in.Status.TotalPorts
}

// GetSwitchPorts returns value of status.switchPorts field if it is not nil,
// otherwise 0.
func (in *NetworkSwitch) GetSwitchPorts() uint32 {
	return in.Status.SwitchPorts
}

// GetState returns value of status.state field if it is not nil,
// otherwise empty string.
func (in *NetworkSwitch) GetState() string {
	return ptr.Deref(in.Status.State, "")
}

// GetMessage returns value of status.message field if it is not nil,
// otherwise empty string.
func (in *NetworkSwitch) GetMessage() string {
	return ptr.Deref(in.Status.Message, "")
}

// ----------------------------------------
// NetworkSwitchSpec setters
// ----------------------------------------

// SetInventoryRef sets passed argument as a value of
// spec.inventoryRef.name field.
func (in *NetworkSwitch) SetInventoryRef(value string) {
	in.Spec.InventoryRef = &v1.LocalObjectReference{Name: value}
}

// SetManaged sets passed argument as a value of
// spec.managed field.
func (in *NetworkSwitch) SetManaged(value bool) {
	in.Spec.Managed = ptr.To(value)
}

// SetCordon sets passed argument as a value of
// spec.cordon field.
func (in *NetworkSwitch) SetCordon(value bool) {
	in.Spec.Cordon = ptr.To(value)
}

// SetTopSpine sets passed argument as a value of
// spec.topSpine field.
func (in *NetworkSwitch) SetTopSpine(value bool) {
	in.Spec.TopSpine = ptr.To(value)
}

// SetScanPorts sets passed argument as a value of
// spec.scanPorts field.
func (in *NetworkSwitch) SetScanPorts(value bool) {
	in.Spec.ScanPorts = ptr.To(value)
}

// ----------------------------------------
// NetworkSwitchStatus setters
// ----------------------------------------

// SetConfigRef sets passed argument as a value of
// status.configRef.name field.
func (in *NetworkSwitch) SetConfigRef(value string) {
	// if value == constants.EmptyString {
	//	in.Status.ConfigRef = nil
	// } else {
	//	in.Status.ConfigRef = v1.LocalObjectReference{Name: value}
	// }
	in.Status.ConfigRef = v1.LocalObjectReference{Name: value}
}

// SetRoutingConfigTemplate sets passed argument as a value of
// status.routingConfigTemplate.name field.
func (in *NetworkSwitch) SetRoutingConfigTemplate(value string) {
	// if value == constants.EmptyString {
	//	in.Status.RoutingConfigTemplate = nil
	// } else {
	//	in.Status.RoutingConfigTemplate = &v1.LocalObjectReference{Name: value}
	// }
	in.Status.RoutingConfigTemplate = v1.LocalObjectReference{Name: value}
}

// SetASN sets passed argument as a value of
// status.asn field.
func (in *NetworkSwitch) SetASN(value uint32) {
	in.Status.ASN = value
}

// SetLayer sets passed argument as a value of
// status.layer field.
func (in *NetworkSwitch) SetLayer(value uint32) {
	in.Status.Layer = value
}

// SetRole sets passed argument as a value of
// status.role field. Possible values:
//   - spine
//   - leaf
//   - edge-leaf
func (in *NetworkSwitch) SetRole(value string) {
	// switch value {
	// case "":
	//	in.Status.Role = nil
	// default:
	//	in.Status.Role = ptr.To(value)
	// }
	in.Status.Role = value
}

// SetTotalPorts sets passed argument as a value of
// status.totalPorts field.
func (in *NetworkSwitch) SetTotalPorts(value uint32) {
	in.Status.TotalPorts = value
}

// SetSwitchPorts sets passed argument as a value of
// status.switchPorts field.
func (in *NetworkSwitch) SetSwitchPorts(value uint32) {
	in.Status.SwitchPorts = value
}

// SetState sets passed argument as a value of
// status.state field. If passed argument is equal
// to empty string, then nil will be set as field
// value. Possible not empty values:
//   - Initial
//   - Processing
//   - Ready
//   - Invalid
//   - Pending
func (in *NetworkSwitch) SetState(value string) {
	switch value {
	case "":
		in.Status.State = nil
	default:
		in.Status.State = ptr.To(value)
	}
}

// SetMessage sets passed argument as a value of
// status.message field. If passed argument is equal
// to empty string, then nil will be set as field
// value.
func (in *NetworkSwitch) SetMessage(value string) {
	switch value {
	case "":
		in.Status.Message = nil
	default:
		in.Status.Message = ptr.To(value)
	}
}

// ----------------------------------------
// ConditionSpec getters
// ----------------------------------------

// GetCondition returns the pointer to ConditionSpec if it is
// stored in the list of switch's conditions, otherwise nil.
func (in *NetworkSwitch) GetCondition(name string) *ConditionSpec {
	for _, item := range in.Status.Conditions {
		if item.GetName() == name {
			return item
		}
	}
	return nil
}

// GetState returns value of ConditionSpec.State if it is not nil,
// otherwise false.
func (in *ConditionSpec) GetState() bool {
	if in == nil {
		return false
	}
	return ptr.Deref(in.State, false)
}

// GetName returns value of ConditionSpec.Name if it is not nil,
// otherwise empty string.
func (in *ConditionSpec) GetName() string {
	return ptr.Deref(in.Name, "")
}

// GetLastTransitionTimestamp returns value of ConditionSpec.LastTransitionTimestamp
// if it is not nil, otherwise empty string.
func (in *ConditionSpec) GetLastTransitionTimestamp() string {
	return ptr.Deref(in.LastTransitionTimestamp, "")
}

// GetLastUpdateTimestamp returns value of ConditionSpec.LastUpdateTimestamp
// if it is not nil, otherwise empty string.
func (in *ConditionSpec) GetLastUpdateTimestamp() string {
	return ptr.Deref(in.LastUpdateTimestamp, "")
}

// GetReason returns value of ConditionSpec.Reason if it is not nil,
// otherwise empty string.
func (in *ConditionSpec) GetReason() string {
	return ptr.Deref(in.Reason, "")
}

// GetMessage returns value of ConditionSpec.Message if it is not nil,
// otherwise empty string.
func (in *ConditionSpec) GetMessage() string {
	return ptr.Deref(in.Message, "")
}

// ----------------------------------------
// ConditionSpec setters
// ----------------------------------------

func NewCondition(name string) *ConditionSpec {
	return &ConditionSpec{Name: ptr.To(name)}
}

// UpdateCondition updates the switch object's status.conditions list.
// Using passed "name" argument, it looks up for existing condition
// with provided name. In case condition was found, it will be updated.
// Otherwise, new ConditionSpec object will be created and added to the list.
func (in *NetworkSwitch) UpdateCondition(name, reason, message string, state bool) {
	var conditions []*ConditionSpec
	ts := time.Now()
	c := in.GetCondition(name)
	if c == nil {
		conditions = make([]*ConditionSpec, len(in.Status.Conditions)+1)
		copy(conditions, in.Status.Conditions)
		conditions[len(in.Status.Conditions)] = NewCondition(name).
			SetState(state).
			SetReason(reason).
			SetMessage(message).
			SetLastUpdateTimestamp(ts.String()).
			SetLastTransitionTimestamp(ts.String())
	} else {
		conditions = make([]*ConditionSpec, len(in.Status.Conditions))
		for i, item := range in.Status.Conditions {
			if item.GetName() != c.GetName() {
				conditions[i] = item
				continue
			}
			if c.GetState() != state {
				c.SetState(state).SetReason(reason).SetMessage(message).SetLastTransitionTimestamp(ts.String())
			}
			c.SetLastUpdateTimestamp(ts.String())
			conditions[i] = c
		}
	}
	in.Status.Conditions = conditions
}

// SetState sets passed argument as a value of
// condition.state field.
func (in *ConditionSpec) SetState(value bool) *ConditionSpec {
	in.State = ptr.To(value)
	return in
}

// SetLastUpdateTimestamp sets passed argument as a value of
// condition.lastUpdateTimestamp field.
func (in *ConditionSpec) SetLastUpdateTimestamp(value string) *ConditionSpec {
	in.LastUpdateTimestamp = ptr.To(value)
	return in
}

// SetLastTransitionTimestamp sets passed argument as a value of
// condition.lastTransitionTimestamp field.
func (in *ConditionSpec) SetLastTransitionTimestamp(value string) *ConditionSpec {
	in.LastTransitionTimestamp = ptr.To(value)
	return in
}

// SetReason sets passed argument as a value of
// condition.reason field.
func (in *ConditionSpec) SetReason(value string) *ConditionSpec {
	switch value {
	case "":
		in.Reason = nil
	default:
		in.Reason = ptr.To(value)
	}
	return in
}

// FlushReason sets nil value of condition.reason field.
func (in *ConditionSpec) FlushReason() *ConditionSpec {
	in.Reason = nil
	return in
}

// SetMessage sets passed argument as a value of
// condition.message field.
func (in *ConditionSpec) SetMessage(value string) *ConditionSpec {
	switch value {
	case "":
		in.Message = nil
	default:
		in.Message = ptr.To(value)
	}
	return in
}

// FlushMessage sets nil value of condition.message field.
func (in *ConditionSpec) FlushMessage() *ConditionSpec {
	in.Message = nil
	return in
}

// ----------------------------------------
// InterfaceSpec getters
// ----------------------------------------

// GetMACAddress returns value of macAddress field of
// given InterfaceSpec object if it is not nil, otherwise
// empty string.
func (in *InterfaceSpec) GetMACAddress() string {
	return ptr.Deref(in.MACAddress, "")
}

// GetSpeed returns value of speed field of given
// InterfaceSpec object if it is not nil, otherwise 0.
func (in *InterfaceSpec) GetSpeed() uint32 {
	return ptr.Deref(in.Speed, 0)
}

// GetDirection returns value of direction field of
// given InterfaceSpec object if it is not nil, otherwise
// empty string.
func (in *InterfaceSpec) GetDirection() string {
	return ptr.Deref(in.Direction, "")
}

// ----------------------------------------
// InterfaceSpec setters
// ----------------------------------------

// SetMACAddress sets passed argument as a value of
// macAddress field for given InterfaceSpec object.
func (in *InterfaceSpec) SetMACAddress(value string) {
	in.MACAddress = ptr.To(value)
}

// SetSpeed sets passed argument as a value of
// speed field for given InterfaceSpec object.
func (in *InterfaceSpec) SetSpeed(value uint32) {
	in.Speed = ptr.To(value)
}

// SetDirection sets passed argument as a value of
// direction field for given InterfaceSpec object.
// Possible values:
//   - north
//   - south
func (in *InterfaceSpec) SetDirection(value string) {
	in.Direction = ptr.To(value)
}

// SetIPEmpty empties the list of assigned IP addresses
// for given InterfaceSpec object.
func (in *InterfaceSpec) SetIPEmpty() {
	in.IP = make([]*IPAddressSpec, 0)
}

// SetPortParametersEmpty resets portParameters field
// for given InterfaceSpec object by assigning the empty
// PortParametersSpec struct as a value of this field.
func (in *InterfaceSpec) SetPortParametersEmpty() {
	in.PortParametersSpec = &PortParametersSpec{}
}

// ----------------------------------------
// PortParametersSpec getters
// ----------------------------------------

// GetLanes returns value of lanes field of given
// PortParametersSpec object if it is not nil,
// otherwise 0.
func (in *PortParametersSpec) GetLanes() uint32 {
	return ptr.Deref(in.Lanes, 0)
}

// GetMTU returns value of mtu field of given
// PortParametersSpec object if it is not nil,
// otherwise 0.
func (in *PortParametersSpec) GetMTU() uint32 {
	return ptr.Deref(in.MTU, 0)
}

// GetIPv4MaskLength returns value of ipv4MaskLength
// field of given PortParametersSpec object if it is not nil,
// otherwise 0.
func (in *PortParametersSpec) GetIPv4MaskLength() uint32 {
	return ptr.Deref(in.IPv4MaskLength, 0)
}

// GetIPv6Prefix returns value of ipv6Prefix field of given
// PortParametersSpec object if it is not nil, otherwise 0.
func (in *PortParametersSpec) GetIPv6Prefix() uint32 {
	return ptr.Deref(in.IPv6Prefix, 0)
}

// GetFEC returns value of fec field of given
// PortParametersSpec object if it is not nil,
// otherwise empty string.
func (in *PortParametersSpec) GetFEC() string {
	return ptr.Deref(in.FEC, "")
}

// GetState returns value of state field of given
// PortParametersSpec object if it is not nil,
// otherwise empty string.
func (in *PortParametersSpec) GetState() string {
	return ptr.Deref(in.State, "")
}

// ----------------------------------------
// PortParametersSpec setters
// ----------------------------------------

// SetLanes sets passed argument as value of lanes field
// for given PortParametersSpec object.
func (in *PortParametersSpec) SetLanes(value uint32) {
	in.Lanes = ptr.To(value)
}

// SetMTU sets passed argument as value of mtu field
// for given PortParametersSpec object.
func (in *PortParametersSpec) SetMTU(value uint32) {
	in.MTU = ptr.To(value)
}

// SetIPv4MaskLength sets passed argument as value of
// ipv4MaskLength for given PortParametersSpec object.
func (in *PortParametersSpec) SetIPv4MaskLength(value uint32) {
	in.IPv4MaskLength = ptr.To(value)
}

// SetIPv6Prefix sets passed argument as value of
// ipv6Prefix field for given PortParametersSpec object.
func (in *PortParametersSpec) SetIPv6Prefix(value uint32) {
	in.IPv6Prefix = ptr.To(value)
}

// SetFEC sets passed argument as value of fec field
// for given PortParametersSpec object. Possible values:
//   - rs
//   - none
func (in *PortParametersSpec) SetFEC(value string) {
	in.FEC = ptr.To(value)
}

// SetState sets passed argument as value of state field
// for given PortParametersSpec object. Possible values:
//   - up
//   - down
func (in *PortParametersSpec) SetState(value string) {
	in.State = ptr.To(value)
}

// ----------------------------------------
// IPAddressSpec getters
// ----------------------------------------

// GetAddress returns value of address field of given
// IPAddressSpec object if it is not nil, otherwise empty string.
func (in *IPAddressSpec) GetAddress() string {
	return ptr.Deref(in.Address, "")
}

// GetAddressFamily returns value of addressFamily field of given
// IPAddressSpec object if it is not nil, otherwise empty string.
func (in *IPAddressSpec) GetAddressFamily() string {
	return ptr.Deref(in.AddressFamily, "")
}

// GetExtraAddress returns value of extraAddress field of given
// IPAddressSpec object if it is not nil, otherwise false.
func (in *IPAddressSpec) GetExtraAddress() bool {
	return ptr.Deref(in.ExtraAddress, false)
}

// GetObjectReferenceName returns value of objectReference.name field
// of given IPAddressSpec object if objectReference is not nil,
// otherwise empty string.
func (in *IPAddressSpec) GetObjectReferenceName() string {
	if ptr.AllPtrFieldsNil(in.ObjectReference) {
		return ""
	}
	return ptr.Deref(in.ObjectReference.Name, "")
}

// GetObjectReferenceNamespace returns value of objectReference.namespace
// field of given IPAddressSpec object if objectReference is not nil,
// otherwise empty string.
func (in *IPAddressSpec) GetObjectReferenceNamespace() string {
	if ptr.AllPtrFieldsNil(in.ObjectReference) {
		return ""
	}
	return ptr.Deref(in.ObjectReference.Namespace, "")
}

// ----------------------------------------
// IPAddressSpec setters
// ----------------------------------------

// SetAddress sets passed argument as value of address
// field for given IPAddressSpec object.
func (in *IPAddressSpec) SetAddress(value string) {
	in.Address = ptr.To(value)
}

// SetAddressFamily sets passed argument as value of
// addressFamily field for given IPAddressSpec object.
// Possible values:
//   - IPv4
//   - IPv6
func (in *IPAddressSpec) SetAddressFamily(value string) {
	in.AddressFamily = ptr.To(value)
}

// SetExtraAddress sets passed argument as value of
// extraAddress field for given IPAddressSpec object.
func (in *IPAddressSpec) SetExtraAddress(value bool) {
	in.ExtraAddress = ptr.To(value)
}

// SetObjectReference updates value of objectReference field of
// given IPAddressSpec object with new ObjectReference object
// where Name and Namespace fields are assigned with passed arguments.
func (in *IPAddressSpec) SetObjectReference(name, namespace string) {
	in.ObjectReference = &ObjectReference{
		Name:      ptr.To(name),
		Namespace: ptr.To(namespace),
	}
}

// ----------------------------------------
// PeerSpec getters
// ----------------------------------------

// GetObjectReferenceName returns value of objectReference.name field
// of given PeerSpec object if objectReference is not nil,
// otherwise empty string.
func (in *PeerSpec) GetObjectReferenceName() string {
	if ptr.AllPtrFieldsNil(in.ObjectReference) {
		return ""
	}
	return ptr.Deref(in.ObjectReference.Name, "")
}

// GetObjectReferenceNamespace returns value of objectReference.namespace
// field of given PeerSpec object if objectReference is not nil,
// otherwise empty string.
func (in *PeerSpec) GetObjectReferenceNamespace() string {
	if ptr.AllPtrFieldsNil(in.ObjectReference) {
		return ""
	}
	return ptr.Deref(in.ObjectReference.Namespace, "")
}

// ----------------------------------------
// PeerSpec setters
// ----------------------------------------

// SetObjectReference updates value of objectReference field of
// given PeerSpec object with new ObjectReference object where
// Name and Namespace fields are assigned with passed arguments.
func (in *PeerSpec) SetObjectReference(name, namespace string) {
	in.ObjectReference = &ObjectReference{
		Name:      ptr.To(name),
		Namespace: ptr.To(namespace),
	}
}

// ----------------------------------------
// PeerInfoSpec getters
// ----------------------------------------

// GetChassisID returns value of chassisID field of
// given PeerInfoSpec object if it is not nil, otherwise
// empty string.
func (in *PeerInfoSpec) GetChassisID() string {
	return ptr.Deref(in.ChassisID, "")
}

// GetSystemName returns value of systemName field
// of given PeerInfoSpec object if it is not nil,
// otherwise empty string.
func (in *PeerInfoSpec) GetSystemName() string {
	return ptr.Deref(in.SystemName, "")
}

// GetPortID returns value of portID field of given
// PeerInfoSpec object if it is not nil, otherwise
// empty string.
func (in *PeerInfoSpec) GetPortID() string {
	return ptr.Deref(in.PortID, "")
}

// GetPortDescription returns value of portDescription
// field of given PeerInfoSpec object if it is not nil,
// otherwise empty string.
func (in *PeerInfoSpec) GetPortDescription() string {
	return ptr.Deref(in.PortDescription, "")
}

// GetType returns value of type field of given
// PeerInfoSpec object if it is not nil, otherwise
// empty string.
func (in *PeerInfoSpec) GetType() string {
	return ptr.Deref(in.Type, "")
}

// ----------------------------------------
// PeerInfoSpec setters
// ----------------------------------------

// SetChassisID sets passed argument as value of chassisID
// field for given PeerInfoSpec object.
func (in *PeerInfoSpec) SetChassisID(value string) {
	in.ChassisID = ptr.To(value)
}

// SetSystemName sets passed argument as value of systemName
// field for given PeerInfoSpec object.
func (in *PeerInfoSpec) SetSystemName(value string) {
	in.SystemName = ptr.To(value)
}

// SetPortID sets passed argument as value of portID field
// for given PeerSpecInfo object.
func (in *PeerInfoSpec) SetPortID(value string) {
	in.PortID = ptr.To(value)
}

// SetPortDescription sets passed arguments as value of
// portDescription field for given PeerInfoSpec object.
func (in *PeerInfoSpec) SetPortDescription(value string) {
	in.PortDescription = ptr.To(value)
}

// SetType sets passed argument as value of type field
// for given PeerInfoSpec object. Possible values:
//   - machine
//   - switch
//   - router (for future use)
//   - undefined
func (in *PeerInfoSpec) SetType(value string) {
	in.Type = ptr.To(value)
}

// ----------------------------------------
// SubnetSpec getters
// ----------------------------------------

// GetSubnetObjectRefName returns value of Subnet.name field
// of given SubnetSpec object if Subnet is not nil,
// otherwise empty string.
func (in *SubnetSpec) GetSubnetObjectRefName() string {
	if ptr.AllPtrFieldsNil(in.Subnet) {
		return ""
	}
	return ptr.Deref(in.Subnet.Name, "")
}

// GetSubnetObjectRefNamespace returns value of Subnet.namespace
// field of given SubnetSpec object if Subnet is not nil,
// otherwise empty string.
func (in *SubnetSpec) GetSubnetObjectRefNamespace() string {
	if ptr.AllPtrFieldsNil(in.Subnet) {
		return ""
	}
	return ptr.Deref(in.Subnet.Namespace, "")
}

// GetNetworkObjectRefName returns value of Network.name field
// of given SubnetSpec object if Network is not nil,
// otherwise empty string.
func (in *SubnetSpec) GetNetworkObjectRefName() string {
	if ptr.AllPtrFieldsNil(in.Network) {
		return ""
	}
	return ptr.Deref(in.Network.Name, "")
}

// GetNetworkObjectRefNamespace returns value of Network.namespace
// field of given SubnetSpec object if Network is not nil,
// otherwise empty string.
func (in *SubnetSpec) GetNetworkObjectRefNamespace() string {
	if ptr.AllPtrFieldsNil(in.Network) {
		return ""
	}
	return ptr.Deref(in.Network.Namespace, "")
}

// GetCIDR returns value of cidr field of given SubnetSpec object
// if it is not nil, otherwise empty string.
func (in *SubnetSpec) GetCIDR() string {
	return ptr.Deref(in.CIDR, "")
}

// GetAddressFamily returns value of addressFamily field of given
// SubnetSpec object if it is not nil, otherwise empty string.
func (in *SubnetSpec) GetAddressFamily() string {
	return ptr.Deref(in.AddressFamily, "")
}

// ----------------------------------------
// SubnetSpec setters
// ----------------------------------------

// SetSubnetObjectRef updates value of Subnet field of
// given SubnetSpec object with new ObjectReference object where
// Name and Namespace fields are assigned with passed arguments.
func (in *SubnetSpec) SetSubnetObjectRef(name, namespace string) {
	in.Subnet = &ObjectReference{
		Name:      ptr.To(name),
		Namespace: ptr.To(namespace),
	}
}

// SetNetworkObjectRef updates value of Network field of
// given SubnetSpec object with new ObjectReference object where
// Name and Namespace fields are assigned with passed arguments.
func (in *SubnetSpec) SetNetworkObjectRef(name, namespace string) {
	in.Network = &ObjectReference{
		Name:      ptr.To(name),
		Namespace: ptr.To(namespace),
	}
}

// SetCIDR sets passed argument as value of cidr field of
// given SubnetSpec object.
func (in *SubnetSpec) SetCIDR(value string) {
	in.CIDR = ptr.To(value)
}

// SetAddressFamily sets passed argument as value of
// addressFamily field for given SubnetSpec object.
// Possible values:
//   - IPv4
//   - IPv6
func (in *SubnetSpec) SetAddressFamily(value string) {
	in.AddressFamily = ptr.To(value)
}

// ----------------------------------------
// InterfaceOverridesSpec getters
// ----------------------------------------

// GetName returns value of name field of given
// InterfaceOverridesSpec object if it is not nil,
// otherwise empty string.
func (in *InterfaceOverridesSpec) GetName() string {
	return ptr.Deref(in.Name, "")
}

// ----------------------------------------
// InterfaceOverridesSpec setters
// ----------------------------------------

// SetName sets passed argument as value of name field
// of given InterfaceOverridesSpec object.
func (in *InterfaceOverridesSpec) SetName(value string) {
	in.Name = ptr.To(value)
}

// GetAddress returns value of address field of given
// AdditionalIPSpec object if it is not nil, otherwise empty string.
func (in *AdditionalIPSpec) GetAddress() string {
	return ptr.Deref(in.Address, "")
}

// GetIPv4 returns value of ipv4 field of given AddressFamilyMap
// object if it is not nil, otherwise false.
func (in *AddressFamiliesMap) GetIPv4() bool {
	return ptr.Deref(in.IPv4, false)
}

// GetIPv6 returns value of ipv6 field of given AddressFamilyMap
// object if it is not nil, otherwise false.
func (in *AddressFamiliesMap) GetIPv6() bool {
	return ptr.Deref(in.IPv6, false)
}

// GetLabelKey returns value of labelKey field of given
// FieldSelectorSpec object if it is not nil, otherwise
// empty string.
func (in *FieldSelectorSpec) GetLabelKey() string {
	return ptr.Deref(in.LabelKey, "")
}

// GetLoopbacksSelection helps to get the loopback addresses selection spec
// in safely manner with handling the case when whole IPAMSpec spec equals to nil.
func (in *IPAMSpec) GetLoopbacksSelection() *IPAMSelectionSpec {
	if in == nil {
		return nil
	}
	return in.LoopbackAddresses
}

// GetSubnetsSelection helps to get the south subnets selection spec
// in safely manner with handling the case when whole IPAMSpec spec equals to nil.
func (in *IPAMSpec) GetSubnetsSelection() *IPAMSelectionSpec {
	if in == nil {
		return nil
	}
	return in.SouthSubnets
}

func (in *NetworkSwitch) UpdateSwitchLabels(inv *Inventory) {
	appliedLabels := map[string]string{
		constants.InventoriedLabel: "true",
		constants.LabelChassisID: strings.ReplaceAll(
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

func (in *NetworkSwitch) UpdateSwitchAnnotations(inv *Inventory) {
	hardwareAnnotations := make(map[string]string, 3)
	softwareAnnotations := make(map[string]string, 5)
	if inv.Spec.System != nil {
		hardwareAnnotations[constants.HardwareSerialAnnotation] = inv.Spec.System.SerialNumber
		hardwareAnnotations[constants.HardwareManufacturerAnnotation] = inv.Spec.System.Manufacturer
		hardwareAnnotations[constants.HardwareSkuAnnotation] = inv.Spec.System.ProductSKU
	}
	if inv.Spec.Distro != nil {
		softwareAnnotations[constants.SoftwareOnieAnnotation] = "false"
		softwareAnnotations[constants.SoftwareAsicAnnotation] = inv.Spec.Distro.AsicType
		softwareAnnotations[constants.SoftwareVersionAnnotation] = inv.Spec.Distro.CommitID
		softwareAnnotations[constants.SoftwareOSAnnotation] = "sonic"
		softwareAnnotations[constants.SoftwareHostnameAnnotation] = inv.Spec.Host.Name
	}
	if in.Annotations == nil {
		in.Annotations = make(map[string]string)
	}
	in.Annotations[constants.HardwareChassisIDAnnotation] = strings.ReplaceAll(
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

// StateReady checks actual value of .status.state field and
// returns boolean value whether it matches corresponding constant.
func (in *NetworkSwitch) StateReady() bool {
	return in.GetState() == constants.SwitchStateReady
}

// StateNotReady checks actual value of .status.state field and
// returns boolean value whether it differs from "Ready".
func (in *NetworkSwitch) StateNotReady() bool {
	return in.GetState() != constants.SwitchStateReady
}

// StatePending checks actual value of .status.state field and
// returns boolean value whether it matches corresponding constant.
func (in *NetworkSwitch) StatePending() bool {
	return in.GetState() == constants.SwitchStatePending
}

// StateInvalid checks actual value of .status.state field and
// returns boolean value whether it matches corresponding constant.
func (in *NetworkSwitch) StateInvalid() bool {
	return in.GetState() == constants.SwitchStateInvalid
}

// StateInitial checks actual value of .status.state field and
// returns boolean value whether it matches corresponding constant.
func (in *NetworkSwitch) StateInitial() bool {
	return in.GetState() == constants.SwitchStateInitial
}

// StateProcessing checks actual value of .status.state field and
// returns boolean value whether it matches corresponding constant.
func (in *NetworkSwitch) StateProcessing() bool {
	return in.GetState() == constants.SwitchStateProcessing
}

// Uninitialized checks whether the .status.state field is empty.
func (in *NetworkSwitch) Uninitialized() bool {
	return in.GetState() == constants.EmptyString
}
