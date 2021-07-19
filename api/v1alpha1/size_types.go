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
	"errors"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SizeSpec defines the desired state of Size
type SizeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Constraints is a list of selectors based on machine properties
	// +kubebuilder:validation:Optional
	Constraints []ConstraintSpec `json:"constraints,omitempty"`
}

type AggregateType string

const (
	CAverageAggregateType AggregateType = "avg"
	CSumAggregateType     AggregateType = "sum"
	CMinAggregateType     AggregateType = "min"
	CMaxAggregateType     AggregateType = "max"
	CCountAggregateType   AggregateType = "count"
)

// ConstraintSpec contains conditions of contraint that should be applied on resource
type ConstraintSpec struct {
	// Path is a path to the struct field constraint will be applied to
	// +kubebuilder:validation:Optional
	Path string `json:"path,omitempty"`
	// Aggregate defines whether collection values should be aggregated
	// for constraint checks, in case if path defines selector for collection
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=avg;sum;min;max;count
	Aggregate AggregateType `json:"agg,omitempty"`
	// Equal contains an exact expected value
	// +kubebuilder:validation:Optional
	Equal *ConstraintValSpec `json:"eq,omitempty"`
	// NotEqual contains an exact not expected value
	// +kubebuilder:validation:Optional
	NotEqual *ConstraintValSpec `json:"neq,omitempty"`
	// LessThan contains an highest expected value, exclusive
	// +kubebuilder:validation:Optional
	LessThan *resource.Quantity `json:"lt,omitempty"`
	// LessThan contains an highest expected value, inclusive
	// +kubebuilder:validation:Optional
	LessThanOrEqual *resource.Quantity `json:"lte,omitempty"`
	// LessThan contains an lowest expected value, exclusive
	// +kubebuilder:validation:Optional
	GreaterThan *resource.Quantity `json:"gt,omitempty"`
	// GreaterThanOrEqual contains an lowest expected value, inclusive
	// +kubebuilder:validation:Optional
	GreaterThanOrEqual *resource.Quantity `json:"gte,omitempty"`
}

// ConstraintValSpec is a wrapper around value for constraint.
// Since it is not possilbwe to set oneOf/anyOf through kubebuilder
// markers, type is set to number here, and patched with kustomize
// See https://github.com/kubernetes-sigs/kubebuilder/issues/301
// +kubebuilder:validation:Type=number
type ConstraintValSpec struct {
	Literal *string            `json:"-"`
	Numeric *resource.Quantity `json:"-"`
}

func (s *ConstraintValSpec) MarshalJSON() ([]byte, error) {
	if s.Literal != nil && s.Numeric != nil {
		return nil, errors.New("unable to marshal JSON since both numeric and literal fields are set")
	}
	if s.Literal != nil {
		return json.Marshal(s.Literal)
	}
	if s.Numeric != nil {
		return json.Marshal(s.Numeric)
	}
	return nil, nil
}

func (s *ConstraintValSpec) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		s.Literal = nil
		s.Numeric = nil
		return nil
	}
	q := resource.Quantity{}
	err := q.UnmarshalJSON(data)
	if err == nil {
		s.Numeric = &q
		return nil
	}
	var str string
	err = json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	s.Literal = &str
	return nil
}

// SizeStatus defines the observed state of Size
type SizeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Size is the Schema for the sizes API
type Size struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SizeSpec   `json:"spec,omitempty"`
	Status SizeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SizeList contains a list of Size
type SizeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Size `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Size{}, &SizeList{})
}
