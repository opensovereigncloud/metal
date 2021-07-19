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

type AggregateItem struct {
	// SourcePath is a path in Inventory spec aggregate will be applied to
	// +kubebuilder:validation:Required
	SourcePath string `json:"sourcePath"`
	// TargetPath is a path in Inventory status `computed` field
	// +kubebuilder:validation:Required
	TargetPath string `json:"targetPath"`
	// Aggregate defines whether collection values should be aggregated
	// for constraint checks, in case if path defines selector for collection
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=avg;sum;min;max;count
	Aggregate AggregateType `json:"aggregate,omitempty"`
}

// AggregateSpec defines the desired state of Aggregate
type AggregateSpec struct {
	// Aggregates is a list of aggregates required to be computed
	// +kubebuilder:validation:Optional
	Aggregates []AggregateItem `json:"aggregates"`
}

// AggregateStatus defines the observed state of Aggregate
type AggregateStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Aggregate is the Schema for the aggregates API
type Aggregate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AggregateSpec   `json:"spec,omitempty"`
	Status AggregateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AggregateList contains a list of Aggregate
type AggregateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Aggregate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Aggregate{}, &AggregateList{})
}
