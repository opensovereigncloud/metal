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

// SwitchConnectionSpec defines the desired state of SwitchConnection
//+kubebuilder:object:generate=true
type SwitchConnectionSpec struct {
	// Switch refers to current switch parameters
	//+kubebuilder:validation:Required
	Switch *ConnectedSwitchSpec `json:"switch"`
	// UpstreamSwitch refers to up-level switch
	//+kubebuilder:validation:Optional
	UpstreamSwitches *UpstreamSwitchesSpec `json:"upstreamSwitches,omitempty"`
	// DownstreamSwitch refers to down-level switch
	//+kubebuilder:validation:Optional
	DownstreamSwitches *DownstreamSwitchesSpec `json:"downstreamSwitches,omitempty"`
	// ConnectionLevel refers to vertical level of switch connection
	//+kubebuilder:validation:Optional
	//+kubebuilder:default=255
	ConnectionLevel uint8 `json:"connectionLevel"`
}

// UpstreamSwitchesSpec defines upstream switches count and properties
//+kubebuilder:object:generate=true
type UpstreamSwitchesSpec struct {
	// Count refers to upstream switches count
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:default=0
	Count int `json:"count"`
	// Switches refers to connected upstream switches
	//+kubebuilder:validation:Optional
	Switches []ConnectedSwitchSpec `json:"switches"`
}

// DownstreamSwitchesSpec defines downstream switches count and properties
//+kubebuilder:object:generate=true
type DownstreamSwitchesSpec struct {
	// Count refers to upstream switches count
	//+kubebuilder:validation:Optional
	//+kubebuilder:validation:default=0
	Count int `json:"count"`
	// Switches refers to connected upstream switches
	//+kubebuilder:validation:Optional
	Switches []ConnectedSwitchSpec `json:"switches"`
}

// ConnectedSwitchSpec defines switch connected to another switch
//+kubebuilder:object:generate=true
type ConnectedSwitchSpec struct {
	// Name refers to switch's name
	//+kubebuilder:validation:Optional
	Name string `json:"name"`
	// Namespace refers to switch's namespace
	//+kubebuilder:validation:Optional
	Namespace string `json:"namespace"`
	// ChassisID refers to switch's chassis id
	//+kubebuilder:validation:Required
	ChassisID string `json:"chassisId"`
}

// SwitchConnectionStatus defines the observed state of SwitchConnection
type SwitchConnectionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=swc
//+kubebuilder:printcolumn:name="ConnectionLevel",type=integer,JSONPath=`.spec.connectionLevel`,description="Vertical level of switch connection"

// SwitchConnection is the Schema for the switchconnections API
type SwitchConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwitchConnectionSpec   `json:"spec,omitempty"`
	Status SwitchConnectionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SwitchConnectionList contains a list of SwitchConnection
type SwitchConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SwitchConnection `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SwitchConnection{}, &SwitchConnectionList{})
}
