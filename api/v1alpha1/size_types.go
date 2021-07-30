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
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/jsonpath"
)

const (
	CLabelPrefix = "machine.onmetal.de/size-"
)

// SizeSpec defines the desired state of Size
type SizeSpec struct {
	// Constraints is a list of selectors based on machine properties
	// +kubebuilder:validation:Optional
	Constraints []ConstraintSpec `json:"constraints,omitempty"`
}

// SizeStatus defines the observed state of Size
type SizeStatus struct {
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

func (in *Size) GetMatchLabel() string {
	return CLabelPrefix + in.Name
}

func (in *Size) Matches(inventory *Inventory) (bool, error) {
	for _, constraint := range in.Spec.Constraints {
		jp := jsonpath.New(constraint.Path)
		// Do not return errors if data is not found
		jp.AllowMissingKeys(true)
		err := jp.Parse(normalizeJSONPath(constraint.Path))
		if err != nil {
			return false, err
		}

		data, err := jp.FindResults(&inventory.Spec)
		if err != nil {
			return false, err
		}

		dataLen := len(data)
		// If validation data is empty, return "does not match"
		if dataLen == 0 {
			return false, nil
		}
		// If data has more than 2 arrays, multiple result sets were returned
		// we do not support that case
		if dataLen > 1 {
			return false, errors.New("multiple selection results are not supported")
		}

		validationData := data[0]
		validationDataLen := len(validationData)
		// If result array is empty for some reason, return "does not match"
		if validationDataLen == 0 {
			return false, nil
		}

		var valid bool
		// If result set has only one value, validate it as a single value
		// even if it is an aggregate, since result will be the same
		if validationDataLen == 1 {
			valid, err = constraint.MatchSingleValue(&validationData[0])
		} else {
			valid, err = constraint.MatchMultipleValues(constraint.Aggregate, validationData)
		}

		if err != nil {
			return false, err
		}

		if !valid {
			return false, nil
		}
	}

	return true, nil
}
