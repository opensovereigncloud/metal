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
	// ChassisID refers to switch chassis id
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Pattern=`^([0-9a-fA-F]{2}[:]){5}([0-9a-fA-F]{2})$`
	ChassisID string `json:"chassisId"`
	//Region refers to the switch's region
	//+kubebuilder:validation:Required
	Region string `json:"region"`
	//AvailabilityZone refers to the switch's availability zone
	//+kubebuilder:validation:Required
	AvailabilityZone string `json:"availabilityZone"`
}

// SwitchAssignmentStatus defines the observed state of SwitchAssignment
type SwitchAssignmentStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=swa
//+kubebuilder:printcolumn:name="Region",type=string,JSONPath=`.spec.region`,description="switch's region"
//+kubebuilder:printcolumn:name="Availability Zone",type=string,JSONPath=`.spec.availabilityZone`,description="switch's AZ"

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
