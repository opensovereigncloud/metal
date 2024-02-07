// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
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
