// +k8s:deepcopy-gen=package

// Package v1alpha3 is a version of the API.
// +groupName=benchmark.onmetal.de
//go:generate go run github.com/ahmetb/gen-crd-api-reference-docs -api-dir . -config ../../../hack/api-reference/template.json -template-dir ../../../hack/api-reference/template -out-file ../../../docs/api-reference/benchmark/machine.md

package v1alpha3
