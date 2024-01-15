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

package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BenchmarkSpec contains machine benchmark results.
type BenchmarkSpec struct {
	// Benchmarks is the collection of benchmarks.
	Benchmarks map[string]Benchmarks `json:"benchmarks,omitempty"`
}

type Benchmarks []BenchmarkResult

type BenchmarkResult struct {
	// Name is the specific benchmark name. e.g. `fio-1k`.
	Name string `json:"name"`
	// Value is the exact result of specific benchmark.
	Value uint64 `json:"value,omitempty"`
}

// BenchmarkStatus contains machine benchmarks deviations.
type BenchmarkStatus struct {
	// BenchmarkDeviations shows the difference between last and current benchmark results.
	BenchmarkDeviations map[string]BenchmarkDeviations `json:"machine_deviation,omitempty"`
}

type BenchmarkDeviations []BenchmarkDeviation

// BenchmarkDeviation is a deviation between old value and the new one.
type BenchmarkDeviation struct {
	// Name is the specific benchmark name. e.g. `fio-1k`.
	Name string `json:"name"`
	// Value is the exact result of specific benchmark.
	Value string `json:"value,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +genclient

// Benchmark is the Schema for the machines API.
type Benchmark struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BenchmarkSpec   `json:"spec,omitempty"`
	Status BenchmarkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BenchmarkList contains a list of Benchmark.
type BenchmarkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Benchmark `json:"items"`
}
