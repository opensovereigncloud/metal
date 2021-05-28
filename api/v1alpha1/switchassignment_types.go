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

// SwitchAssignmentSpec defines the desired state of SwitchAssignment
//+kubebuilder:object:generate=true
type SwitchAssignmentSpec struct {
	// Role refers to the role of the switch. Always "Spine"
	//+kubebuilder:validation:default=Spine
	Role string `json:"role"`
	// Serial refers to switch serial number
	//+kubebuilder:validation:Required
	Serial string `json:"serial"`
	// ChassisID refers to switch chassis id
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Pattern='^([0-9a-fA-F]{2}[:]){5}([0-9a-fA-F]{2})$'
	ChassisID string `json:"chassisId"`
}

// SwitchAssignmentStatus defines the observed state of SwitchAssignment
type SwitchAssignmentStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=swa
//+kubebuilder:printcolumn:name="Role",type=string,JSONPath=`.spec.role`,description="switch's role"

// SwitchAssignment is the Schema for the switch assignments API
type SwitchAssignment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwitchAssignmentSpec   `json:"spec,omitempty"`
	Status SwitchAssignmentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SwitchAssignmentList contains a list of SwitchAssignment
type SwitchAssignmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SwitchAssignment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SwitchAssignment{}, &SwitchAssignmentList{})
}
