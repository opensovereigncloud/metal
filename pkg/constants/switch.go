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
	NICDown string = "down"
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

	LocationRoomAnnotation string = "location/room"
	LocationRowAnnotation  string = "location/row"
	LocationRackAnnotation string = "location/rack"
	LocationHUAnnotation   string = "location/hu"
)

const SwitchFinalizer string = "switch.onmetal.de/finalizer"

const (
	InventoriedLabel            string = "switch.onmetal.de/inventoried"
	LabelChassisID              string = "metalapi.onmetal.de/chassis-id"
	SwitchTypeLabel             string = "switch.onmetal.de/type"
	SwitchConfigTypeLabelPrefix string = "switch.onmetal.de/type-"
	SwitchConfigLayerLabel      string = "switch.onmetal.de/layer"

	IPAMObjectPurposeLabel string = "ipam.onmetal.de/object-purpose"
	IPAMObjectOwnerLabel   string = "ipam.onmetal.de/object-owner"
)

const (
	EmptyString string = ""

	SwitchStateInvalid    string = "Invalid"
	SwitchStateInitial    string = "Initial"
	SwitchStateProcessing string = "Processing"
	SwitchStateReady      string = "Ready"
	SwitchStatePending    string = "Pending"

	SizeLabel string = "machine.onmetal.de/size-switch"

	IPAMObjectNICNameLabel string = "ipam.onmetal.de/interface-name"

	IPAMLoopbackPurpose    string = "loopback"
	IPAMSouthSubnetPurpose string = "south-subnet"
	IPAMSwitchPortPurpose  string = "switch-port"

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

	IPv4LoopbackBits uint32 = 32
	IPv6LoopbackBits uint32 = 128

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
)
