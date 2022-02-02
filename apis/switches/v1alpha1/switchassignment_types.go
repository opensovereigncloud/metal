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

package v1alpha1

import (
	subnetv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type State string

const (
	CAssignmentStatePending  State = "pending"
	CAssignmentStateFinished State = "finished"
	CStateDeleting           State = "deleting"
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
	Region *RegionSpec `json:"region"`
}

// RegionSpec defines region info
//+kubebuilder:object:generate=true
type RegionSpec struct {
	// Name refers to the switch's region
	//+kubebuilder:validation:Pattern=^[a-z0-9]([-./a-z0-9]*[a-z0-9])?$
	//+kubebuilder:validation:Required
	Name string `json:"name"`
	// AvailabilityZone refers to the switch's availability zone
	//+kubebuilder:validation:Required
	AvailabilityZone string `json:"availabilityZone"`
}

// SwitchAssignmentStatus defines the observed state of SwitchAssignment
type SwitchAssignmentStatus struct {
	//+kubebuilder:validation:Enum=pending;finished;deleting
	State State `json:"state"`
	//+kubebuilder:validation:Optional
	Switch *LinkedSwitchSpec `json:"switch"`
}

type LinkedSwitchSpec struct {
	//+kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	//+kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=swa
//+kubebuilder:printcolumn:name="Switch Chassis ID",type=string,JSONPath=`.spec.chassisId`,description="switch's chassis Id"
//+kubebuilder:printcolumn:name="Region",type=string,JSONPath=`.spec.region.name`,description="switch's region"
//+kubebuilder:printcolumn:name="Availability Zone",type=string,JSONPath=`.spec.region.availabilityZone`,description="switch's AZ"
//+kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`,description="Assignment state"

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

// NamespacedName returns assignment's name and namespace as
// built-in type.
func (in *SwitchAssignment) NamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: in.Namespace,
		Name:      in.Name,
	}
}

// FillStatus fills resource status with provided values
func (in *SwitchAssignment) FillStatus(state State, relatedSwitch *LinkedSwitchSpec) {
	in.Status.State = state
	in.Status.Switch = relatedSwitch
}

// GetListFilter builds list options object
func (in *SwitchAssignment) GetListFilter() (*client.ListOptions, error) {
	labelsReq, err := labels.NewRequirement(LabelChassisId, selection.In, []string{MacToLabel(in.Spec.ChassisID)})
	if err != nil {
		return nil, err
	}
	selector := labels.NewSelector()
	selector = selector.Add(*labelsReq)
	opts := &client.ListOptions{
		LabelSelector: selector,
		Limit:         100,
	}
	return opts, nil
}

func (in *RegionSpec) ConvertToSubnetRegion() []subnetv1alpha1.Region {
	return []subnetv1alpha1.Region{
		{
			Name:              in.Name,
			AvailabilityZones: []string{in.AvailabilityZone},
		},
	}
}

func ConvertFromSubnetRegion(src []subnetv1alpha1.Region) *RegionSpec {
	if len(src) > 1 {
		return nil
	}
	if len(src[0].AvailabilityZones) > 1 {
		return nil
	}
	return &RegionSpec{
		Name:             src[0].Name,
		AvailabilityZone: src[0].AvailabilityZones[0],
	}
}
