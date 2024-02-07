// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Package v1alpha4 contains API Schema definitions for the switch v1beta1 API group
// +kubebuilder:object:generate=true
// +groupName=metal.ironcore.dev
package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// SchemeGroupVersion is group version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{Group: "metal.ironcore.dev", Version: "v1alpha4"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme.
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Aggregate{},
		&AggregateList{},
		&Benchmark{},
		&BenchmarkList{},
		&Inventory{},
		&InventoryList{},
		&Machine{},
		&MachineList{},
		&Size{},
		&SizeList{},
		&NetworkSwitch{},
		&NetworkSwitchList{},
		&SwitchConfig{},
		&SwitchConfigList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
