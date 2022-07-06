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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SwitchConfigSpec contains desired configuration for selected switches
//+kubebuilder:object:generate=true
type SwitchConfigSpec struct {
	// Switches contains label selector to pick up Switch objects
	//+kubebuilder:validation:Required
	Switches *metav1.LabelSelector `json:"switches,omitempty"`
	// PortsDefaults contains switch port parameters which will be applied to all ports of the switches
	// which fit selector conditions
	//+kubebuilder:validation:Required
	PortsDefaults *PortParametersSpec `json:"portsDefaults"`
	// IPAM refers to selectors for subnets which will be used for Switch object
	//+kubebuilder:validation:Required
	IPAM *GeneralIPAMSpec `json:"ipam"`
}

type GeneralIPAMSpec struct {
	// CarrierSubnets contains label selector for Subnet object where switch's south subnet
	// should be reserved
	//+kubebuilder:validation:Required
	CarrierSubnets *metav1.LabelSelector `json:"carrierSubnets"`
	// LoopbackSubnets contains label selector for Subnet object where switch's loopback
	// IP addresses should be reserved
	//+kubebuilder:validation:Required
	LoopbackSubnets *metav1.LabelSelector `json:"loopbackSubnets"`
	// SouthSubnets defines selector for subnets object which will be assigned to switch
	//+kubebuilder:validation:Optional
	SouthSubnets *IPAMSelectionSpec `json:"southSubnets,omitempty"`
	// LoopbackAddresses defines selector for IP objects which should be referenced as switch's loopback addresses
	//+kubebuilder:validation:Optional
	LoopbackAddresses *IPAMSelectionSpec `json:"loopbackAddresses,omitempty"`
}

// SwitchConfigStatus contains observed state of SwitchConfig
//+kubebuilder:object:generate=true
type SwitchConfigStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=sc

// SwitchConfig is the Schema for switch config API
type SwitchConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwitchConfigSpec   `json:"spec,omitempty"`
	Status SwitchConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SwitchConfigList contains a list of SwitchConfig
type SwitchConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SwitchConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SwitchConfig{}, &SwitchConfigList{})
}