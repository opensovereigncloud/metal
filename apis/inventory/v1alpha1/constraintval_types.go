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

package v1alpha1

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/json"
)

// ConstraintValSpec is a wrapper around value for constraint.
// Since it is not possilble to set oneOf/anyOf through kubebuilder
// markers, type is set to number here, and patched with kustomize
// See https://github.com/kubernetes-sigs/kubebuilder/issues/301
// +kubebuilder:validation:Type=number
type ConstraintValSpec struct {
	Literal *string            `json:"-"`
	Numeric *resource.Quantity `json:"-"`
}

func (in *ConstraintValSpec) MarshalJSON() ([]byte, error) {
	if in.Literal != nil && in.Numeric != nil {
		return nil, errors.New("unable to marshal JSON since both numeric and literal fields are set")
	}
	if in.Literal != nil {
		return json.Marshal(in.Literal)
	}
	if in.Numeric != nil {
		return json.Marshal(in.Numeric)
	}
	return nil, nil
}

func (in *ConstraintValSpec) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		in.Literal = nil
		in.Numeric = nil
		return nil
	}
	q := resource.Quantity{}
	err := q.UnmarshalJSON(data)
	if err == nil {
		in.Numeric = &q
		return nil
	}
	var str string
	err = json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	in.Literal = &str
	return nil
}

func (in *ConstraintValSpec) Compare(value *reflect.Value) (int, error) {
	if in.Literal != nil {
		s, err := valueToString(value)
		if err != nil {
			return 0, err
		}
		return strings.Compare(s, *in.Literal), nil
	}

	if in.Numeric != nil {
		q, err := valueToQuantity(value)
		if err != nil {
			return 0, err
		}
		return q.Cmp(*in.Numeric), nil
	}

	return 0, errors.New("both numeric and literal constraints are nil")
}
