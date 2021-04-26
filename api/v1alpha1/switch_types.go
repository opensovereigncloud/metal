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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SwitchSpec defines the desired state of Switch
// +kubebuilder:object:generate=true
type SwitchSpec struct {
	// ID referring to switch object id
	// +kubebuilder:validation:Required
	//todo: +kubebuilder:validation:Pattern=`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})$`
	ID string `json:"id"`
	// Partition referring to switch's location
	// +kubebuilder:validation:Optional
	Partition string `json:"partition,omitempty"`
	// Room referring to switch's location
	// +kubebuilder:validation:Optional
	Room string `json:"room,omitempty"`
	// Row referring to switch's location
	// +kubebuilder:validation:Optional
	Row int16 `json:"row,omitempty"`
	// Rack referring to switch's location
	// +kubebuilder:validation:Optional
	Rack int16 `json:"rack,omitempty"`
	// Ports referring to switch port number
	// +kubebuilder:validation:Required
	Ports uint64 `json:"ports"`
	// NeighboursCount referring to switch's neighbours
	// +kubebuilder:validation:Optional
	NeighboursCount uint8 `json:"neighboursCount,omitempty"`
	// Neighbours referring to switch's neighbours
	// +kubebuilder:validation:Optional
	Neighbours []NeighbourSpec `json:"neighbours,omitempty"`
}

// NeighbourSpec defines switch's neighbour
// +kubebuilder:object:generate=true
type NeighbourSpec struct {
	// ID referring to neighbour id
	// +kubebuilder:validation:Required
	//todo: +kubebuilder:validation:Pattern=`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})$`
	ID string `json:"id"`
	// Name referring to neighbour name
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Type referring to neighbour machine type
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=Machine;Switch
	Type string `json:"type"`
	// Port referring to neighbour port name
	// +kubebuilder:validation:Required
	Port string `json:"port"`
	// MACAddress referring to neighbour MAC address
	// +kubebuilder:validation:Optional
	MACAddress string `json:"macAddress,omitempty"`
}

// SwitchStatus defines the observed state of Switch
type SwitchStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Partition",type=string,JSONPath=`.spec.partition`,description="Partition name"
//+kubebuilder:printcolumn:name="Room",type=string,JSONPath=`.spec.room`,description="Room name"
//+kubebuilder:printcolumn:name="Row",type=integer,JSONPath=`.spec.row`,description="Row number"
//+kubebuilder:printcolumn:name="Rack",type=integer,JSONPath=`.spec.rack`,description="Rack number"
//+kubebuilder:printcolumn:name="Ports",type=integer,JSONPath=`.spec.ports`,description="Total amount of hardware network interfaces"
//+kubebuilder:printcolumn:name="Neighbours",type=integer,JSONPath=`.spec.neighboursCount`,description="Total amount of neighbours"

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
