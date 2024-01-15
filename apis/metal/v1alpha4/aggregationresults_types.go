// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package v1alpha4

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
)

// +kubebuilder:validation:Type=object
type AggregationResults struct {
	Object map[string]interface{} `json:"-"`
}

func (in AggregationResults) MarshalJSON() ([]byte, error) {
	if in.Object == nil {
		return json.Marshal(map[string]interface{}{})
	}
	return json.Marshal(in.Object)
}

func (in *AggregationResults) UnmarshalJSON(b []byte) error {
	stringVal := string(b)
	if stringVal == "null" {
		in.Object = nil
		return nil
	}
	if err := json.Unmarshal(b, &in.Object); err != nil {
		return err
	}

	return nil
}

func (in *AggregationResults) DeepCopyInto(out *AggregationResults) {
	if in == nil {
		out = nil
	} else if in.Object == nil {
		out.Object = nil
	} else {
		out.Object = runtime.DeepCopyJSON(in.Object)
	}
}

func (in *AggregationResults) DeepCopy() *AggregationResults {
	if in == nil {
		return nil
	}
	out := new(AggregationResults)
	in.DeepCopyInto(out)
	return out
}
