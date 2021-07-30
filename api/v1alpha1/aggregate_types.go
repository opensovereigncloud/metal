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
)

type AggregateType string

const (
	CAverageAggregateType AggregateType = "avg"
	CSumAggregateType     AggregateType = "sum"
	CMinAggregateType     AggregateType = "min"
	CMaxAggregateType     AggregateType = "max"
	CCountAggregateType   AggregateType = "count"
)

type AggregateItem struct {
	// SourcePath is a path in Inventory spec aggregate will be applied to
	// +kubebuilder:validation:Required
	SourcePath JSONPath `json:"sourcePath"`
	// TargetPath is a path in Inventory status `computed` field
	// +kubebuilder:validation:Required
	TargetPath JSONPath `json:"targetPath"`
	// Aggregate defines whether collection values should be aggregated
	// for constraint checks, in case if path defines selector for collection
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=avg;sum;min;max;count
	Aggregate AggregateType `json:"aggregate,omitempty"`
}

// AggregateSpec defines the desired state of Aggregate
type AggregateSpec struct {
	// Aggregates is a list of aggregates required to be computed
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
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

func (in *Aggregate) Compute(inventory *Inventory) (interface{}, error) {
	resultMap := make(map[string]interface{})

	for _, ai := range in.Spec.Aggregates {
		jp, err := ai.SourcePath.ToK8sJSONPath()
		if err != nil {
			return nil, err
		}

		jp.AllowMissingKeys(true)
		data, err := jp.FindResults(inventory.Spec)
		if err != nil {
			return nil, err
		}

		var aggregatedValue interface{} = nil
		tokenizedPath := ai.TargetPath.Tokenize()

		dataLen := len(data)
		if dataLen == 0 {
			if err := setValueToPath(resultMap, tokenizedPath, aggregatedValue); err != nil {
				return nil, err
			}
			continue
		}
		if dataLen > 1 {
			return nil, errors.New("expected only one value collection to be returned for aggregation")
		}

		values := data[0]
		valuesLen := len(values)

		if valuesLen == 0 {
			if err := setValueToPath(resultMap, tokenizedPath, aggregatedValue); err != nil {
				return nil, err
			}
			continue
		}

		if ai.Aggregate == "" {
			interfacedValues := make([]interface{}, valuesLen)
			for i, value := range values {
				interfacedValues[i] = value.Interface()
			}
			aggregatedValue = interfacedValues
		} else {
			aggregatedValue, err = makeAggregate(ai.Aggregate, values)
		}

		if err := setValueToPath(resultMap, tokenizedPath, aggregatedValue); err != nil {
			return nil, err
		}
	}

	return resultMap, nil
}
