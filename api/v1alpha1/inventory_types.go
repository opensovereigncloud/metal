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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// InventorySpec contains result of inventorization process on the host
type InventorySpec struct {
	// System contains DMI system information
	// +kubebuilder:validation:Required
	System *SystemSpec `json:"system,omitempty"`
	// IPMIs contains info about IPMI interfaces on the host
	// +kubebuilder:validation:Optional
	IPMIs []IPMISpec `json:"ipmis,omitempty"`
	// Blocks contains info about block devices on the host
	// +kubebuilder:validation:Required
	Blocks *BlockTotalSpec `json:"blocks,omitempty"`
	// Memory contains info block devices on the host
	// +kubebuilder:validation:Required
	Memory *MemorySpec `json:"memory,omitempty"`
	// CPUs contains info about cpus, cores and threads
	// +kubebuilder:validation:Required
	CPUs *CPUTotalSpec `json:"cpus,omitempty"`
	// NICs contains info about network interfaces and network discovery
	// +kubebuilder:validation:Required
	NICs *NICTotalSpec `json:"nics,omitempty"`
	// Virt is a virtualization detected on host
	// +kubebuilder:validation:Optional
	Virt *VirtSpec `json:"virt,omitempty"`
}

// SystemSpec contains DMI system information
type SystemSpec struct {
	// ID is a UUID of a system board
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`
	ID string `json:"id,omitempty"`
	// Manufacturer refers to the company that produced the product
	// +kubebuilder:validation:Required
	Manufacturer string `json:"manufacturer,omitempty"`
	// ProductSKU is a product's Stock Keeping Unit
	// +kubebuilder:validation:Required
	ProductSKU string `json:"productSku,omitempty"`
	// SerialNumber contains serial number of a system
	// +kubebuilder:validation:Required
	SerialNumber string `json:"serialNumber,omitempty"`
}

// IPMISpec contains info about IPMI module
type IPMISpec struct {
	// IPAddress is an IP address assigned to IPMI network interface
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	IPAddress string `json:"ipAddress,omitempty"`
	// MACAddress is a MAC address of IPMI's network interface
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`
	MACAddress string `json:"macAddress,omitempty"`
}

// BlockTotalSpec contains disk aggregates and disk descriptions
type BlockTotalSpec struct {
	// Count is a total disk count on a host
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Count uint64 `json:"count,omitempty"`
	// Capacity is a total disk storage capacity on a host in bytes
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Capacity uint64 `json:"capacity,omitempty"`
	// Blocks contains info about block devices on the host
	// +kubebuilder:validation:Required
	Blocks []BlockSpec `json:"blocks,omitempty"`
}

// BlockSpec contains info about block device
type BlockSpec struct {
	// Name is a name of the device registered by Linux Kernel
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`
	// Type refers to data carrier form-factor
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=Floppy;CD-ROM;SCSI;IDE;NVMe;USB;MMC;VirtIO;Xen
	Type string `json:"type,omitempty"`
	// Rotational shows whether disk is solid state or not
	// +kubebuilder:validation:Required
	Rotational bool `json:"rotational"`
	// Bus is a type of hardware interface used to connect the disk to the system
	// +kubebuilder:validation:Optional
	Bus string `json:"system,omitempty"`
	// Model is a unique hardware part identifier
	// +kubebuilder:validation:Required
	Model string `json:"model,omitempty"`
	// Size is a disk space available in bytes
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Size uint64 `json:"size,omitempty"`
	// PartitionTable is a partition table currently written to the disk
	// +kubebuilder:validation:Optional
	PartitionTable *PartitionTableSpec `json:"partitionTable,omitempty"`
}

// PartitionTableSpec contains info about partition table on block device
type PartitionTableSpec struct {
	// Type is a format of partition table
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=MBR;GPT
	Type string `json:"type,omitempty"`
	// Partitions are active partition records on disk
	// +kubebuilder:validation:Optional
	Partitions []PartitionSpec `json:"partitions,omitempty"`
}

// PartitionSpec contains info about partition
type PartitionSpec struct {
	// ID is a GUID of GPT partition or number for MBR partition
	// +kubebuilder:validation:Required
	ID string `json:"id,omitempty"`
	// Name is a human readable name given to the partition
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// Size is a size of partition in bytes
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Size uint64 `json:"size,omitempty"`
}

// MemorySpec contains info about RAM on host
type MemorySpec struct {
	// Total is a total amount of RAM on host
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Total uint64 `json:"total,omitempty"`
}

// CPUTotalSpec contains an overview of CPUs and calculated total
type CPUTotalSpec struct {
	// Sockets represents a total amount of physical processors
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Sockets uint64 `json:"sockets,omitempty"`
	// Cores is a total amount of physical cores
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Cores uint64 `json:"cores,omitempty"`
	// Cores is a total amount of logical cores
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Threads uint64 `json:"threads,omitempty"`
	// CPUs is a collection of specs for physical CPUs
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	CPUs []CPUSpec `json:"cpus,omitempty"`
}

// CPUSpec contains info about CPUs on hsot machine
type CPUSpec struct {
	// PhysicalID is an ID of physical CPU
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	PhysicalID uint64 `json:"physicalId,omitempty"`
	// LogicalIDs is a collection of logical CPU nums related to the physical CPU (required for NUMA)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	LogicalIDs []uint64 `json:"logicalIds,omitempty"`
	// Cores is a number of physical cores
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Cores uint64 `json:"cores,omitempty"`
	// Siblings is a number of logical CPUs/threads
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Siblings uint64 `json:"siblings,omitempty"`
	// VendorID is a manufacturer identifire
	// +kubebuilder:validation:Required
	VendorID string `json:"vendorId,omitempty"`
	// Family refers to processor type
	// +kubebuilder:validation:Optional
	Family string `json:"family,omitempty"`
	// Model is a reference id of the model
	// +kubebuilder:validation:Required
	Model string `json:"model,omitempty"`
	// ModelName is a common name of the processor
	// +kubebuilder:validation:Required
	ModelName string `json:"modelName,omitempty"`
	// Stepping is an iteration of the architecture
	// +kubebuilder:validation:Optional
	Stepping string `json:"stepping,omitempty"`
	// Microcode is a firmware reference
	// +kubebuilder:validation:Optional
	Microcode string `json:"microcode,omitempty"`
	// MHz is a logical core frequency
	// +kubebuilder:validation:Optional
	MHz resource.Quantity `json:"mhz,omitempty"`
	// CacheSize is an L2 cache size
	// +kubebuilder:validation:Optional
	CacheSize string `json:"cacheSize,omitempty"`
	// FPU defines if CPU has a Floating Point Unit
	// +kubebuilder:validation:Optional
	FPU bool `json:"fpu"`
	// FPUException
	// +kubebuilder:validation:Optional
	FPUException bool `json:"fpuException"`
	// CPUIDLevel
	// +kubebuilder:validation:Optional
	CPUIDLevel uint64 `json:"cpuIdLevel,omitempty"`
	// WP tells if WP bit is present
	// +kubebuilder:validation:Optional
	WP bool `json:"wp"`
	// Flags defines a list of low-level computing capabilities
	// +kubebuilder:validation:Optional
	Flags []string `json:"flags,omitempty"`
	// VMXFlags defines a list of virtualization capabilities
	// +kubebuilder:validation:Optional
	VMXFlags []string `json:"vmxFlags,omitempty"`
	// Bugs contains a list of known hardware bugs
	// +kubebuilder:validation:Optional
	Bugs []string `json:"bugs,omitempty"`
	// BogoMIPS is a synthetic performance metric
	// +kubebuilder:validation:Required
	BogoMIPS resource.Quantity `json:"bogoMips,omitempty"`
	// CLFlushSize size for cache line flushing feature
	// +kubebuilder:validation:Optional
	CLFlushSize uint64 `json:"clFlushSize,omitempty"`
	// CacheAlignment is a cache size
	// +kubebuilder:validation:Optional
	CacheAlignment uint64 `json:"cacheAlignment,omitempty"`
	// AddressSizes is an info about address transition system
	// +kubebuilder:validation:Optional
	AddressSizes string `json:"addressSizes,omitempty"`
	// PowerManagement
	// +kubebuilder:validation:Optional
	PowerManagement string `json:"powerManagement,omitempty"`
}

// NICSpec contains info about network interfaces and aggregates
type NICTotalSpec struct {
	// Count is a total amount of hardware NICs on host
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Count uint64 `json:"count,omitempty"`
	// NICs contains info about network interfaces and network discovery
	// +kubebuilder:validation:Required
	NICs []NICSpec `json:"nics,omitempty"`
}

// NICSpec contains info about network interfaces
type NICSpec struct {
	// Name is a name of the device registered by Linux Kernel
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`
	// PCIAddress is the PCI bus address network interface is connected to
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([0-9a-fA-F]{4}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}.[0-9]{1})|([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})$`
	PCIAddress string `json:"pciAddress,omitempty"`
	// MACAddress is the MAC address of network interface
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([0-9A-Fa-f]{2}[:.-]){5}([0-9A-Fa-f]{2})$`
	MACAddress string `json:"macAddress,omitempty"`
	// MTU is refers to Maximum Transmission Unit
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	MTU uint16 `json:"mtu,omitempty"`
	// Speed is a speed of network interface in Mbits/s
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Speed uint32 `json:"speed,omitempty"`
	// LLDP is a collection of LLDP messages received by the network interface
	// +kubebuilder:validation:Optional
	LLDPs []LLDPSpec `json:"lldps,omitempty"`
	// NDP is a collection of NDP messages received by the network interface
	// +kubebuilder:validation:Optional
	NDPs []NDPSpec `json:"ndps,omitempty"`
}

// LLDPSpec is an entry received by network interface by Link Layer Discovery Protocol
type LLDPSpec struct {
	// ChassisID is a neighbour box identifier
	// +kubebuilder:validation:Required
	ChassisID string `json:"chassisId,omitempty"`
	// SystemName is given name to the neighbour box
	// +kubebuilder:validation:Optional
	SystemName string `json:"systemName,omitempty"`
	// SystemDescription is a short description of the neighbour box
	// +kubebuilder:validation:Optional
	SystemDescription string `json:"systemDescription,omitempty"`
	// PortID is a hardware identifier of the link port
	// +kubebuilder:validation:Required
	PortID string `json:"portId,omitempty"`
	// PortDescription is a short description of the link port
	// +kubebuilder:validation:Optional
	PortDescription string `json:"portDescription,omitempty"`
}

// NDPSpec is an entry received by IPv6 Neighbour Discovery Protocol
type NDPSpec struct {
	// IPAddress is an IPv6 address of a neighbour
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
	IPAddress string `json:"ipAddress,omitempty"`
	// MACAddress is an MAC address of a neighbour
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([0-9A-Fa-f]{2}[:.-]){5}([0-9A-Fa-f]{2})$`
	MACAddress string `json:"macAddress,omitempty"`
	// State is a state of discovery
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=None;Incomplete;Reachable;Stale;Delay;Probe;Failed;No ARP;Permanent
	State string `json:"state,omitempty"`
}

// VirtSpec contains info about detected host virtualization
type VirtSpec struct {
	// VMType is a type of virtual machine engine
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=none;kvm;qemu;bochs;xen;uml;vmware;oracle;microsoft;zvm;parallels;bhyve;qnx;acrn;powervm;other
	VMType string `json:"vmType,omitempty"`
}

// InventoryStatus defines the observed state of Inventory
type InventoryStatus struct {
	// No additional state required for now
}

// Inventory is the Schema for the inventories API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cores",type=integer,JSONPath=`.spec.cpus.cores`,description="Total amount of cores"
// +kubebuilder:printcolumn:name="Memory",type=integer,JSONPath=`.spec.memory.total`,description="RAM amount in bytes"
// +kubebuilder:printcolumn:name="Disks",type=integer,JSONPath=`.spec.blocks.count`,description="Hardware disk count"
// +kubebuilder:printcolumn:name="Storage",type=integer,JSONPath=`.spec.blocks.capacity`,description="Total amount of disk capacity"
// +kubebuilder:printcolumn:name="NICs",type=integer,JSONPath=`.spec.nics.count`,description="Total amount of hardware network interfaces"
type Inventory struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InventorySpec   `json:"spec,omitempty"`
	Status InventoryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InventoryList contains a list of Inventory
type InventoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Inventory `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Inventory{}, &InventoryList{})
}
