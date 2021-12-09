/*
Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MachineSpec contains machine benchmark results.
type MachineSpec struct {
	Benchmarks map[string][]Benchmark `json:"benchmarks,omitempty"`
}

type Benchmark struct {
	Name  string `json:"name"`
	Value uint64 `json:"value,omitempty"`
}

// MachineStatus contains machine benchmarks deviations.
type MachineStatus struct {
	// Deviation shows the difference between last and current benchmark results.
	*Deviation `json:"deviation,omitempty"`
}

type Deviation struct {
	MachineDeviation map[string][]Benchmark `json:"benchmarks,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Machine is the Schema for the machines API
type Machine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineSpec   `json:"spec,omitempty"`
	Status MachineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineList contains a list of Machine
type MachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Machine `json:"items"`
}
