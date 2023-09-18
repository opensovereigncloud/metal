// Copyright 2023 OnMetal authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package constants

import (
	"github.com/onmetal/metal-api/apis/inventory/v1alpha1"
)

const APIVersion string = "v1beta1"

const (
	FECNone string = "none"
	FECRS   string = "rs"
	NICUp   string = "up"
)

const (
	DefaultIPAMFieldRef string = "metadata.name"
)

const (
	SoftwareOSAnnotation       string = "software/os"
	SoftwareVersionAnnotation  string = "software/version"
	SoftwareOnieAnnotation     string = "software/onie"
	SoftwareAsicAnnotation     string = "software/asic"
	SoftwareHostnameAnnotation string = "software/hostname"

	HardwareSkuAnnotation          string = "hardware/sku"
	HardwareManufacturerAnnotation string = "hardware/manufacturer"
	HardwareSerialAnnotation       string = "hardware/serial"
	HardwareChassisIDAnnotation    string = "hardware/chassis-id"
)

//const SwitchFinalizer string = "switch.onmetal.de/finalizer"

const (
	InventoriedLabel       string = "switch.onmetal.de/inventoried"
	LabelChassisID         string = "metalapi.onmetal.de/chassis-id"
	SwitchTypeLabel        string = "switch.onmetal.de/type"
	SwitchConfigLayerLabel string = "switch.onmetal.de/layer"

	IPAMObjectPurposeLabel     string = "ipam.onmetal.de/object-purpose"
	IPAMObjectOwnerLabel       string = "ipam.onmetal.de/object-owner"
	IPAMObjectGeneratedByLabel string = "ipam.onmetal.de/generated-by"
	IPAMObjectNICNameLabel     string = "ipam.onmetal.de/interface-name"

	SizeLabel string = "machine.onmetal.de/size-switch"

	IPAMLoopbackPurpose    string = "loopback"
	IPAMSouthSubnetPurpose string = "south-subnet"
	IPAMSwitchPortPurpose  string = "switch-port"
)

const (
	EmptyString string = ""

	SwitchManager string = "metal-api-controller-manager"

	SwitchStateInvalid    string = "Invalid"
	SwitchStateInitial    string = "Initial"
	SwitchStateProcessing string = "Processing"
	SwitchStateReady      string = "Ready"
	SwitchStatePending    string = "Pending"

	LLDPCapabilityStation v1alpha1.LLDPCapabilities = "Station"

	NeighborTypeMachine string = "machine"
	NeighborTypeSwitch  string = "switch"

	DirectionSouth string = "south"
	DirectionNorth string = "north"

	SwitchRoleSpine string = "spine"
	SwitchRoleLeaf  string = "leaf"

	IPv4AF string = "IPv4"
	IPv6AF string = "IPv6"

	ASNBase uint32 = 4_200_000_000

	SwitchPortNamePrefix string = "Ethernet"

	IPv4MaskLength   uint32 = 32
	IPv6PrefixLength uint32 = 128

	ConditionInitialized      string = "Initialized"
	ConditionInterfacesOK     string = "InterfacesOK"
	ConditionConfigRefOK      string = "ConfigRefOK"
	ConditionPortParametersOK string = "PortParametersOK"
	ConditionNeighborsOK      string = "NeighborsOK"
	ConditionLayerAndRoleOK   string = "LayerAndRoleOK"
	ConditionLoopbacksOK      string = "LoopbacksOK"
	ConditionAsnOK            string = "AsnOK"
	ConditionSubnetsOK        string = "SubnetsOK"
	ConditionIPAddressesOK    string = "IPAddressesOK"
	ConditionReady            string = "Ready"

	ReasonConditionInitialized      string = "InitializationCompleted"
	ReasonConditionInterfacesOK     string = "InterfacesUpdated"
	ReasonConditionConfigRefOK      string = "ConfigReferenceUpdated"
	ReasonConditionPortParametersOK string = "PortParametersUpdated"
	ReasonConditionNeighborsOK      string = "NeighborsUpdated"
	ReasonConditionLayerAndRoleOK   string = "LayerAndRoleUpdated"
	ReasonConditionLoopbacksOK      string = "LoopbacksUpdated"
	ReasonConditionAsnOK            string = "AsnUpdated"
	ReasonConditionSubnetsOK        string = "SubnetsUpdated"
	ReasonConditionIPAddressesOK    string = "IPAddressesUpdated"
	ReasonConditionReady            string = "ReconciliationCompleted"
	ReasonUnmanagedSwitch           string = "UnmanagedSwitch"
	ReasonObjectUnchanged           string = "ObjectUnchanged"

	MessageConditionInitialized      string = "Switch object initialized successfully"
	MessageConditionInterfacesOK     string = "Interfaces updated successfully"
	MessageConditionConfigRefOK      string = "SwitchConfig reference updated successfully"
	MessageConditionPortParametersOK string = "Port parameters updated successfully"
	MessageConditionNeighborsOK      string = "Neighbors data updated successfully"
	MessageConditionLayerAndRoleOK   string = "Layer and role updated successfully"
	MessageConditionLoopbacksOK      string = "Loopback IP addresses updated successfully"
	MessageConditionAsnOK            string = "ASN updated successfully"
	MessageConditionSubnetsOK        string = "Switch's south subnets updated successfully"
	MessageConditionIPAddressesOK    string = "Ports' IP addresses updated successfully"
	MessageConditionReady            string = "Configuration updated successfully"
)

const (
	StateMessageRequestRelatedObjectsFailed string = "failed to request related objects, check conditions for details"
	StateMessageMissingRequirements         string = "some requirements are missing, check conditions for details"
	StateMessageRelatedObjectsStateInvalid  string = "some of related objects are not in required state yet, check conditions for details"
)

const (
	ErrorReasonMissingRequirements  string = "MissingRequirements"
	ErrorReasonRequestFailed        string = "APIRequestFailed"
	ErrorReasonASNCalculationFailed string = "ASNCalculationFailed"
	ErrorReasonIPAssignmentFailed   string = "IPAssignmentFailed"
	ErrorReasonFailedToComputeLayer string = "FailedToComputeLayer"
)

const (
	MessageMissingInventory          string = "failed to get corresponding inventory object, check reference at .spec.inventoryRef field"
	MessageMissingLoopbacks          string = "failed to get corresponding ip objects, check loopback selectors"
	MessageMissingSouthSubnets       string = "failed to get corresponding subnet objects, check south subnet selectors"
	MessageFailedIPAddressRequest    string = "failed to request IP address from upstream peer"
	MessageFailedToDiscoverConfig    string = "failed to discover corresponding SwitchConfig object: check labels applied to SwitchConfig objects and selector in Switch .spec.configSelector"
	MessageMissingLoopbackV4IP       string = "missing requirements: IP object of V4 address family to be assigned to loopback interface"
	MessageRequestFailed             string = "failed to get requested object"
	MessageFailedToAssignIPAddresses string = "failed to assign IP addresses to switch ports"
	MessageParseIPFailed             string = "failed to parse IP address"
	MessageParseCIDRFailed           string = "failed to parse CIDR"
	MessageInvalidInputType          string = "invalid input type"
	MessageMissingAPIVersion         string = "missing API version"
	MessageAPIVersionMismatch        string = "API version mismatch"
	MessageFieldSelectorNotDefined   string = "field selector is not defined"
	MessageUnmarshallingFailed       string = "failed to unmarshal bytes to map"
	MessageMarshallingFailed         string = "failed to marshal input to bytes"
	MessageInvalidFieldPath          string = "invalid field path"
	MessageFailedToComputeLayer      string = "failed to compute layer: possibly no top spine switches were discovered yet"
)

const (
	ErrorUpdateInterfacesFailed     string = "failed to update interfaces"
	ErrorUpdateNeighborsFailed      string = "failed to update neighbors"
	ErrorUpdateLayerAndRoleFailed   string = "failed to update layer and role"
	ErrorUpdateConfigRefFailed      string = "failed to update reference to SwitchConfig"
	ErrorUpdatePortParametersFailed string = "failed to update port parameters"
	ErrorUpdateLoopbacksFailed      string = "failed to update loopbacks"
	ErrorUpdateASNFailed            string = "failed to update ASN"
	ErrorUpdateSubnetsFailed        string = "failed to update south subnets"
	ErrorUpdateSwitchPortIPsFailed  string = "failed to update switch port IP addresses"
)
