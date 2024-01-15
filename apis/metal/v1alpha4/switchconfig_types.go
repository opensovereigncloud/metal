/*
 * Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha4

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ironcore-dev/metal/pkg/constants"
)

// SwitchConfigSpec contains desired configuration for selected switches.
type SwitchConfigSpec struct {
	// Switches contains label selector to pick up NetworkSwitch objects
	// +kubebuilder:validation:Required
	Switches *metav1.LabelSelector `json:"switches,omitempty"`
	// PortsDefaults contains switch port parameters which will be applied to all ports of the switches
	// which fit selector conditions
	// +kubebuilder:validation:Required
	PortsDefaults *PortParametersSpec `json:"portsDefaults"`
	// IPAM refers to selectors for subnets which will be used for NetworkSwitch object
	// +kubebuilder:validation:Required
	IPAM *GeneralIPAMSpec `json:"ipam"`
	// RoutingConfigTemplate contains the reference to the ConfigMap object which contains go-template for FRR config
	// +kubebuilder:validation:Optional
	RoutingConfigTemplate *v1.LocalObjectReference `json:"routingConfigTemplate,omitempty"`
}

// GeneralIPAMSpec contains definition of selectors, used to filter
// required IPAM objects.
type GeneralIPAMSpec struct {
	// AddressFamily contains flags to define which address families are used for switch subnets
	// +kubebuilder:validation:Required
	AddressFamily *AddressFamiliesMap `json:"addressFamily"`
	// CarrierSubnets contains label selector for Subnet object where switch's south subnet
	// should be reserved
	// +kubebuilder:validation:Required
	CarrierSubnets *IPAMSelectionSpec `json:"carrierSubnets"`
	// LoopbackSubnets contains label selector for Subnet object where switch's loopback
	// IP addresses should be reserved
	// +kubebuilder:validation:Required
	LoopbackSubnets *IPAMSelectionSpec `json:"loopbackSubnets"`
	// SouthSubnets defines selector for subnets object which will be assigned to switch
	// +kubebuilder:validation:Optional
	SouthSubnets *IPAMSelectionSpec `json:"southSubnets,omitempty"`
	// LoopbackAddresses defines selector for IP objects which should be referenced as switch's loopback addresses
	// +kubebuilder:validation:Optional
	LoopbackAddresses *IPAMSelectionSpec `json:"loopbackAddresses,omitempty"`
}

// SwitchConfigStatus contains observed state of SwitchConfig.
type SwitchConfigStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=sc
// +genclient

// SwitchConfig is the Schema for switch config API.
type SwitchConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwitchConfigSpec   `json:"spec,omitempty"`
	Status SwitchConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SwitchConfigList contains a list of SwitchConfig.
type SwitchConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SwitchConfig `json:"items"`
}

// GetRoutingConfigTemplate returns value of spec.routingConfigTemplate.name field if
// routingConfigTemplate is not nil, otherwise empty string.
func (in *SwitchConfig) GetRoutingConfigTemplate() string {
	if in.Spec.RoutingConfigTemplate == nil {
		return constants.EmptyString
	}
	return in.Spec.RoutingConfigTemplate.Name
}

// RoutingConfigTemplateIsEmpty checks whether the spec.routingConfigTemplate contains
// value or not.
func (in *SwitchConfig) RoutingConfigTemplateIsEmpty() bool {
	return in.GetRoutingConfigTemplate() == constants.EmptyString
}
