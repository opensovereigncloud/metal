/*
Copyright (c) 2023 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1beta1

import (
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	internal "github.com/onmetal/metal-api/applyconfiguration/internal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	managedfields "k8s.io/apimachinery/pkg/util/managedfields"
	v1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

// SwitchConfigApplyConfiguration represents an declarative configuration of the SwitchConfig type for use
// with apply.
type SwitchConfigApplyConfiguration struct {
	v1.TypeMetaApplyConfiguration    `json:",inline"`
	*v1.ObjectMetaApplyConfiguration `json:"metadata,omitempty"`
	Spec                             *SwitchConfigSpecApplyConfiguration `json:"spec,omitempty"`
	Status                           *switchv1beta1.SwitchConfigStatus   `json:"status,omitempty"`
}

// SwitchConfig constructs an declarative configuration of the SwitchConfig type for use with
// apply.
func SwitchConfig(name, namespace string) *SwitchConfigApplyConfiguration {
	b := &SwitchConfigApplyConfiguration{}
	b.WithName(name)
	b.WithNamespace(namespace)
	b.WithKind("SwitchConfig")
	b.WithAPIVersion("switch.onmetal.de/v1beta1")
	return b
}

// ExtractSwitchConfig extracts the applied configuration owned by fieldManager from
// switchConfig. If no managedFields are found in switchConfig for fieldManager, a
// SwitchConfigApplyConfiguration is returned with only the Name, Namespace (if applicable),
// APIVersion and Kind populated. It is possible that no managed fields were found for because other
// field managers have taken ownership of all the fields previously owned by fieldManager, or because
// the fieldManager never owned fields any fields.
// switchConfig must be a unmodified SwitchConfig API object that was retrieved from the Kubernetes API.
// ExtractSwitchConfig provides a way to perform a extract/modify-in-place/apply workflow.
// Note that an extracted apply configuration will contain fewer fields than what the fieldManager previously
// applied if another fieldManager has updated or force applied any of the previously applied fields.
// Experimental!
func ExtractSwitchConfig(switchConfig *switchv1beta1.SwitchConfig, fieldManager string) (*SwitchConfigApplyConfiguration, error) {
	return extractSwitchConfig(switchConfig, fieldManager, "")
}

// ExtractSwitchConfigStatus is the same as ExtractSwitchConfig except
// that it extracts the status subresource applied configuration.
// Experimental!
func ExtractSwitchConfigStatus(switchConfig *switchv1beta1.SwitchConfig, fieldManager string) (*SwitchConfigApplyConfiguration, error) {
	return extractSwitchConfig(switchConfig, fieldManager, "status")
}

func extractSwitchConfig(switchConfig *switchv1beta1.SwitchConfig, fieldManager string, subresource string) (*SwitchConfigApplyConfiguration, error) {
	b := &SwitchConfigApplyConfiguration{}
	err := managedfields.ExtractInto(switchConfig, internal.Parser().Type("com.github.onmetal.metal-api.apis.switch.v1beta1.SwitchConfig"), fieldManager, b, subresource)
	if err != nil {
		return nil, err
	}
	b.WithName(switchConfig.Name)
	b.WithNamespace(switchConfig.Namespace)

	b.WithKind("SwitchConfig")
	b.WithAPIVersion("switch.onmetal.de/v1beta1")
	return b, nil
}

// WithKind sets the Kind field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Kind field is set to the value of the last call.
func (b *SwitchConfigApplyConfiguration) WithKind(value string) *SwitchConfigApplyConfiguration {
	b.Kind = &value
	return b
}

// WithAPIVersion sets the APIVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the APIVersion field is set to the value of the last call.
func (b *SwitchConfigApplyConfiguration) WithAPIVersion(value string) *SwitchConfigApplyConfiguration {
	b.APIVersion = &value
	return b
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *SwitchConfigApplyConfiguration) WithName(value string) *SwitchConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Name = &value
	return b
}

// WithGenerateName sets the GenerateName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GenerateName field is set to the value of the last call.
func (b *SwitchConfigApplyConfiguration) WithGenerateName(value string) *SwitchConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.GenerateName = &value
	return b
}

// WithNamespace sets the Namespace field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Namespace field is set to the value of the last call.
func (b *SwitchConfigApplyConfiguration) WithNamespace(value string) *SwitchConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Namespace = &value
	return b
}

// WithUID sets the UID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UID field is set to the value of the last call.
func (b *SwitchConfigApplyConfiguration) WithUID(value types.UID) *SwitchConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.UID = &value
	return b
}

// WithResourceVersion sets the ResourceVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ResourceVersion field is set to the value of the last call.
func (b *SwitchConfigApplyConfiguration) WithResourceVersion(value string) *SwitchConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ResourceVersion = &value
	return b
}

// WithGeneration sets the Generation field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Generation field is set to the value of the last call.
func (b *SwitchConfigApplyConfiguration) WithGeneration(value int64) *SwitchConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Generation = &value
	return b
}

// WithCreationTimestamp sets the CreationTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CreationTimestamp field is set to the value of the last call.
func (b *SwitchConfigApplyConfiguration) WithCreationTimestamp(value metav1.Time) *SwitchConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.CreationTimestamp = &value
	return b
}

// WithDeletionTimestamp sets the DeletionTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DeletionTimestamp field is set to the value of the last call.
func (b *SwitchConfigApplyConfiguration) WithDeletionTimestamp(value metav1.Time) *SwitchConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.DeletionTimestamp = &value
	return b
}

// WithDeletionGracePeriodSeconds sets the DeletionGracePeriodSeconds field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DeletionGracePeriodSeconds field is set to the value of the last call.
func (b *SwitchConfigApplyConfiguration) WithDeletionGracePeriodSeconds(value int64) *SwitchConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.DeletionGracePeriodSeconds = &value
	return b
}

// WithLabels puts the entries into the Labels field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Labels field,
// overwriting an existing map entries in Labels field with the same key.
func (b *SwitchConfigApplyConfiguration) WithLabels(entries map[string]string) *SwitchConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	if b.Labels == nil && len(entries) > 0 {
		b.Labels = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Labels[k] = v
	}
	return b
}

// WithAnnotations puts the entries into the Annotations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Annotations field,
// overwriting an existing map entries in Annotations field with the same key.
func (b *SwitchConfigApplyConfiguration) WithAnnotations(entries map[string]string) *SwitchConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	if b.Annotations == nil && len(entries) > 0 {
		b.Annotations = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Annotations[k] = v
	}
	return b
}

// WithOwnerReferences adds the given value to the OwnerReferences field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the OwnerReferences field.
func (b *SwitchConfigApplyConfiguration) WithOwnerReferences(values ...*v1.OwnerReferenceApplyConfiguration) *SwitchConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithOwnerReferences")
		}
		b.OwnerReferences = append(b.OwnerReferences, *values[i])
	}
	return b
}

// WithFinalizers adds the given value to the Finalizers field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Finalizers field.
func (b *SwitchConfigApplyConfiguration) WithFinalizers(values ...string) *SwitchConfigApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	for i := range values {
		b.Finalizers = append(b.Finalizers, values[i])
	}
	return b
}

func (b *SwitchConfigApplyConfiguration) ensureObjectMetaApplyConfigurationExists() {
	if b.ObjectMetaApplyConfiguration == nil {
		b.ObjectMetaApplyConfiguration = &v1.ObjectMetaApplyConfiguration{}
	}
}

// WithSpec sets the Spec field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Spec field is set to the value of the last call.
func (b *SwitchConfigApplyConfiguration) WithSpec(value *SwitchConfigSpecApplyConfiguration) *SwitchConfigApplyConfiguration {
	b.Spec = value
	return b
}

// WithStatus sets the Status field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Status field is set to the value of the last call.
func (b *SwitchConfigApplyConfiguration) WithStatus(value switchv1beta1.SwitchConfigStatus) *SwitchConfigApplyConfiguration {
	b.Status = &value
	return b
}
