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

// nolint
package v1beta1

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	gocidr "github.com/apparentlymart/go-cidr/cidr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/json"
)

const CAPIVersion = "v1beta1"

const (
	CSwitchRoleLeaf     = "leaf"
	CSwitchRoleSpine    = "spine"
	CSwitchRoleEdgeLeaf = "edge-leaf"
)

const (
	CFECNone = "none"
	CFECRS   = "rs"
)

const (
	CNICUp   = "up"
	CNICDown = "down"
)

const (
	CDirectionSouth = "south"
	CDirectionNorth = "north"
)

const (
	CSwitchStateInitial    = "initial"
	CSwitchStateProcessing = "processing"
	CSwitchStateReady      = "ready"
	CSwitchStateInvalid    = "invalid"

	CAgentStateActive = "active"
	CAgentStateFailed = "failed"
)

const (
	CIPAMLabelPrefix            = "ipam.onmetal.de"
	CIPAMObjectPurposeLabel     = "object-purpose"
	CIPAMObjectOwnerLabel       = "object-owner"
	CIPAMPurposeSwitchCarrier   = "switch-carrier"
	CIPAMPurposeSwitchLoopbacks = "switch-loopbacks"
	CIPAMPurposeSouthSubnet     = "south-subnet"
	CIPAMPurposeLoopback        = "loopback"
	CIPAMPurposeInterfaceIP     = "interface-ip"

	CMetalAPILabelPrefix   = "metalapi.onmetal.de"
	CSwitchSizeName        = "switch"
	CInventoriedLabel      = "inventoried"
	CInventoryRefLabel     = "inventory-ref"
	CSwitchConfigTypeLabel = "type-"
	CSwitchTypeLabel       = "type"

	CDefaultIPAMFieldRef = "spec.uuid"
)

const (
	CSoftwareOSAnnotation       = "software/os"
	CSoftwareVersionAnnotation  = "software/version"
	CSoftwareOnieAnnotation     = "software/onie"
	CSoftwareAsicAnnotation     = "software/asic"
	CSoftwareHostnameAnnotation = "software/hostname"

	CHardwareSkuAnnotation          = "hardware/sku"
	CHardwareManufacturerAnnotation = "hardware/manufacturer"
	CHardwareSerialAnnotation       = "hardware/serial"
	CHardwareChassisIDAnnotation    = "hardware/chassis-id"

	CLocationRoomAnnotation = "location/room"
	CLocationRowAnnotation  = "location/row"
	CLocationRackAnnotation = "location/rack"
	CLocationHUAnnotation   = "location/hu"
)

const (
	CPeerTypeMachine   = "machine"
	CPeerTypeSwitch    = "switch"
	CPeerTypeRouter    = "router"
	CPeerTypeUndefined = "undefined"

	CStationCapability = "Station"
	CRouterCapability  = "Router"
	CBridgeCapability  = "Bridge"
	CNDPReachable      = "Reachable"
)

const (
	CDefaultIPv4MaskLength uint8 = 30
	CDefaultIPv6Prefix     uint8 = 127

	CIPv4MaskLengthBits = 32
	CIPv6PrefixBits     = 128
)

const CSwitchPortPrefix = "Ethernet"
const CEmptyString = ""
const CLabelChassisID = "chassis-id"

var (
	InventoriedLabel  = GetLabelKey(CMetalAPILabelPrefix, CInventoriedLabel)
	InventoryRefLabel = GetLabelKey(CMetalAPILabelPrefix, CInventoryRefLabel)

	SwitchTypeLabel       = GetLabelKey(CMetalAPILabelPrefix, CSwitchTypeLabel)
	SwitchConfigTypeLabel = GetLabelKey(CMetalAPILabelPrefix, CSwitchConfigTypeLabel)

	IPAMObjectPurposeLabel = GetLabelKey(CIPAMLabelPrefix, CIPAMObjectPurposeLabel)
	IPAMObjectOwnerLabel   = GetLabelKey(CIPAMLabelPrefix, CIPAMObjectOwnerLabel)

	LabelChassisID = GetLabelKey(CMetalAPILabelPrefix, CLabelChassisID)
)

type ConnectionsMap map[uint8]*SwitchList

// PortParametersSpec contains a set of parameters of switch port
// +kubebuilder:object:generate=true
type PortParametersSpec struct {
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
	// IPv4MaskLength defines prefix of subnet where switch port's IPv4 address should be reserved
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Minimum=1
	//+kubebuilder:validation:Maximum=32
	IPv4MaskLength *uint8 `json:"ipv4MaskLength,omitempty"`
	// IPv6Prefix defines prefix of subnet where switch port's IPv6 address should be reserved
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Minimum=1
	//+kubebuilder:validation:Maximum=128
	IPv6Prefix *uint8 `json:"ipv6Prefix,omitempty"`
	// FEC refers to forward error correction method which should be applied on switch port
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=rs;none
	FEC *string `json:"fec,omitempty"`
	// State defines default state of switch port
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=up;down
	State *string `json:"state,omitempty"`
}

type PortConfigurablesSpec struct {
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
	// FEC refers to forward error correction method which should be applied on switch port
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=rs;none
	FEC *string `json:"fec,omitempty"`
	// State defines default state of switch port
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:Enum=up;down
	State *string `json:"state,omitempty"`
}

// IPAMSpec contains selectors for subnets and loopback IPs and
// definition of address families which should be claimed
// +kubebuilder:object:generate=true
type IPAMSpec struct {
	// SouthSubnets defines selector for subnet object which will be assigned to switch
	//+kubebuilder:validation:Optional
	SouthSubnets *IPAMSelectionSpec `json:"southSubnets,omitempty"`
	// LoopbackAddresses defines selector for IP object which will be assigned to switch's loopback interface
	//+kubebuilder:validation:Optional
	LoopbackAddresses *IPAMSelectionSpec `json:"loopbackAddresses,omitempty"`
}

// IPAMSelectionSpec contains label selector and address family
// +kubebuilder:object:generate=true
type IPAMSelectionSpec struct {
	// AddressFamilies defines what ip address families should be claimed
	//+kubebuilder:validation:Optional
	AddressFamilies *AddressFamiliesMap `json:"addressFamilies,omitempty"`
	// LabelSelector contains label selector to pick up IPAM objects
	//+kubebuilder:validation:Optional
	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty"`
	// FieldSelector contains label key and field path where to get label value for search.
	// If FieldSelector is used as part of IPAM configuration in SwitchConfig object it will
	// reference to the field path in related object. If FieldSelector is used as part of IPAM
	// configuration in Switch object, it will reference to the field path in the same object
	//+kubebuilder:validation:Optional
	FieldSelector *FieldSelectorSpec `json:"fieldSelector,omitempty"`
}

// FieldSelectorSpec contains label key and field path where to get label value for search
// +kubebuilder:object:generate=true
type FieldSelectorSpec struct {
	// LabelKey contains label key
	//+kubebuilder:validation:Optional
	LabelKey string `json:"labelKey"`
	// FieldRef contains reference to the field of resource where to get label's value
	//+kubebuilder:validation:Optional
	FieldRef *v1.ObjectFieldSelector `json:"fieldRef"`
}

// AddressFamiliesMap contains flags regarding what IP address families should be used
// +kubebuilder:object:generate=true
type AddressFamiliesMap struct {
	// IPv4 is a flag defining whether IPv4 is used or not
	//+kubebuilder:validation:Optional
	//+kubebuilder:default=true
	IPv4 bool `json:"ipv4"`
	// IPv6 is a flag defining whether IPv6 is used or not
	//+kubebuilder:validation:Optional
	//+kubebuilder:default=true
	IPv6 bool `json:"ipv6"`
}

// AdditionalIPSpec defines IP address and selector for subnet where address should be reserved
// +kubebuilder:object:generate=true
type AdditionalIPSpec struct {
	// Address contains additional IP address that should be assigned to the interface
	//+kubebuilder:validation:Required
	Address string `json:"address,omitempty"`
	// ParentSubnet contains label selector to pick up IPAM objects
	//+kubebuilder:validation:Optional
	ParentSubnet *metav1.LabelSelector `json:"parentSubnet,omitempty"`
}

// ObjectReference contains enough information to let you locate the
// referenced object across namespaces.
// +kubebuilder:object:generate=true
type ObjectReference struct {
	// Name contains name of the referenced object
	//+kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// Namespace contains namespace of the referenced object
	//+kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
}

// MetalAPIUint8 converts native Golang uint8 to pointer.
func MetalAPIUint8(in uint8) *uint8 {
	return &in
}

// GoUint8 converts pointer to native Golang uint8.
func GoUint8(in *uint8) uint8 {
	return *in
}

// MetalAPIUint16 converts native Golang uint16 to pointer.
func MetalAPIUint16(in uint16) *uint16 {
	return &in
}

// GoUint16 converts pointer to native Golang uint16.
func GoUint16(in *uint16) uint16 {
	return *in
}

func MetalAPIString(in string) *string {
	return &in
}

func GoString(in *string) string {
	return *in
}

// GetLabelSelector builds labels selector.
func GetLabelSelector(key string, op selection.Operator, values []string) (*labels.Requirement, error) {
	return labels.NewRequirement(key, op, values)
}

// GetLabelKey builds label key from prefix and suffix.
func GetLabelKey(prefix, suffix string) string {
	return fmt.Sprintf("%s/%s", prefix, suffix)
}

// LabelFromFieldRef converts field reference to valid label to use in label selector.
func LabelFromFieldRef(obj interface{}, src *FieldSelectorSpec) (label map[string]string, err error) {
	label = make(map[string]string)
	if src == nil {
		return
	}
	mapRepr, err := interfaceToMap(obj)
	if err != nil {
		return
	}
	apiVersion, ok := mapRepr["apiVersion"]
	if !ok {
		err = errors.New("object is not valid API type")
		return
	}
	if src.FieldRef.APIVersion != "" && apiVersion != src.FieldRef.APIVersion {
		err = errors.New("api versions mismatch")
		return
	}
	nested := strings.Split(src.FieldRef.FieldPath, ".")
	currentSearchObj := mapRepr
	for i, f := range nested {
		v, ok := currentSearchObj[f]
		if !ok {
			err = errors.New("referenced field path is invalid")
			return
		}
		if i == len(nested)-1 {
			switch v.(type) {
			case string:
				label[src.LabelKey] = fmt.Sprintf("%v", v)
				return
			default:
				err = errors.New("referenced field must have a string value")
				return
			}
		}
		currentSearchObj, err = interfaceToMap(v)
		if err != nil {
			return
		}
	}
	return
}

func interfaceToMap(i interface{}) (m map[string]interface{}, err error) {
	var raw []byte
	m = make(map[string]interface{})
	raw, err = json.Marshal(i)
	if err != nil {
		return
	}
	err = json.Unmarshal(raw, &m)
	if err != nil {
		return
	}
	return
}

func GetInterfaceSubnet(name string, ifaceAddrPrefix int, network *net.IPNet) (*net.IPNet, error) {
	index, _ := strconv.Atoi(strings.ReplaceAll(name, CSwitchPortPrefix, CEmptyString))
	ifaceNetPrefix, _ := network.Mask.Size()
	ifaceNet, err := gocidr.Subnet(network, ifaceAddrPrefix-ifaceNetPrefix, index)
	if err != nil {
		return nil, err
	}
	return ifaceNet, nil
}
