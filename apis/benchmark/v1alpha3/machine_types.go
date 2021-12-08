/*
Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MachineSpec contains machine benchmark results.
type MachineSpec struct {
	Benchmarks map[string][]Benchmark `json:"benchmarks,omitempty"`
}

type Benchmark struct {
	Name  string `json:"name"`
	Value uint64 `json:"value,omitempty"`
}

// MachineStatus contains machine benchmarks deviations.
type MachineStatus struct {
	// Deviation shows the difference between last and current benchmark results.
	*Deviation `json:"deviation,omitempty"`
}

type Deviation struct {
	Disks    []DiskDeviation    `json:"disks,omitempty"`
	Networks []NetworkDeviation `json:"networks,omitempty"`
}

type DiskDeviation struct {
	// Name contains full device name (like "/dev/hda" etc)
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Results contains disk benchmark results.
	Value []DiskValue `json:"value"`
}

// DiskValue contains block (device) changes.
// +kubebuilder:object:generate=true
type DiskValue struct {
	// IOPattern defines type of I/O pattern (like "read/write/readwrite" etc)
	// more types could be found here: https://fio.readthedocs.io/en/latest/fio_doc.html#cmdoption-arg-readwrite
	// +kubebuilder:validation:Required
	IOPattern string `json:"ioPattern"`
	// SmallBlockReadIOPS contains benchmark result for read IOPS with small block size (device specified block size)
	// +kubebuilder:validation:Required
	SmallBlockReadIOPS string `json:"smallBlockReadIops"`
	// SmallBlockWriteIOPS contains benchmark result for write IOPS with small block size (device specified block size)
	// +kubebuilder:validation:Optional
	SmallBlockWriteIOPS string `json:"smallBlockWriteIops"`
	// BandwidthReadIOPS contains benchmark result for read IOPS with large block size (much larger then device specified block size)
	// +kubebuilder:validation:Optional
	BandwidthReadIOPS string `json:"bandwidthReadIops"`
	// BandwidthWriteIOPS contains benchmark result for write IOPS with large block size (much larger then device specified block size)
	// +kubebuilder:validation:Optional
	BandwidthWriteIOPS string `json:"bandwidthWriteIops"`
}

type NetworkDeviation struct {
	// Name defines a name of network device
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Results contains disk benchmark results.
	Value string `json:"value"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Machine is the Schema for the machines API
type Machine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineSpec   `json:"spec,omitempty"`
	Status MachineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineList contains a list of Machine
type MachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Machine `json:"items"`
}
