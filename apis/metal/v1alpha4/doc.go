// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// +k8s:deepcopy-gen=package
// +k8s:openapi-gen=true

// Package v1beta1 is a version of the API.
// +groupName=metal.ironcore.dev
//go:generate go run github.com/ahmetb/gen-crd-api-reference-docs -api-dir . -config ../../../hack/api-reference/template.json -template-dir ../../../hack/api-reference/template -out-file ../../../docs/api-reference/metal.md

package v1alpha4
