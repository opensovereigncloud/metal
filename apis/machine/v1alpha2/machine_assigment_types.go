/*
Copyright 2022.

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

package v1alpha2

import (
	"inet.af/netaddr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MachineAssignmentSpec defines the desired state of Request
type MachineAssignmentSpec struct {
	// Tolerations define tolerations the Machine has. Only MachinePools whose taints
	// covered by Tolerations will be considered to run the Machine.
	Tolerations []Toleration `json:"tolerations,omitempty"`
	// MachineClass is a reference to the machine class/flavor of the machine.
	MachineClass corev1.LocalObjectReference `json:"machineClass"`
	// Image is the URL providing the operating system image of the machine.
	Image string `json:"image"`
	// Interfaces define a list of network interfaces present on the machine
	NetworkInterfaces []NetworkInterfaces `json:"networkInterfaces,omitempty"`
	// Volumes are volumes attached to this machine.
	Volumes []Volume `json:"volumes,omitempty"`
	// Ignition is a reference to a config map containing the ignition YAML for the machine to boot up.
	// If key is empty, DefaultIgnitionKey will be used as fallback.
	Ignition *ObjectSelector `json:"ignition,omitempty"`
	// EFIVars are variables to pass to EFI while booting up.
	EFIVars []EFIVar `json:"efiVars,omitempty"`
}

// Volume defines a volume attachment of a machine
type Volume struct {
	// Name is the name of the VolumeAttachment
	Name string `json:"name"`
	// Priority is the OS priority of the volume.
	Priority int32 `json:"priority,omitempty"`
	// VolumeAttachmentSource is the source where the storage for the VolumeAttachment resides at.
	VolumeSource `json:",inline"`
}

// VolumeSource specifies the source to use for a VolumeAttachment.
type VolumeSource struct {
	// VolumeClaim instructs the VolumeAttachment to use a VolumeClaim as source for the attachment.
	VolumeClaim *VolumeClaimSource `json:"volumeClaim,omitempty"`
}

// VolumeClaimSource references a VolumeClaim as VolumeAttachment source.
type VolumeClaimSource struct {
	// Ref is a reference to the VolumeClaim.
	Ref corev1.LocalObjectReference `json:"ref"`
}

// ObjectSelector is a reference to a specific 'key' within a ConfigMap resource.
// In some instances, `key` is a required field.
type ObjectSelector struct {
	// The name of the ConfigMap resource being referred to.
	corev1.LocalObjectReference `json:",inline"`
	// The key of the entry in the ConfigMap resource's `data` field to be used.
	// Some instances of this field may be defaulted, in others it may be
	// required.
	// +optional
	Key string `json:"key,omitempty"`
}

// EFIVar is a variable to pass to EFI while booting up.
type EFIVar struct {
	Name  string `json:"name,omitempty"`
	UUID  string `json:"uuid,omitempty"`
	Value string `json:"value"`
}

// Interface is the definition of a single interface
type NetworkInterfaces struct {
	// Name is the name of the interface
	Name string `json:"name"`
	// Target is the referenced resource of this interface.
	Target corev1.LocalObjectReference `json:"target"`
	// Priority is the priority level of this interface
	Priority int32 `json:"priority,omitempty"`
	// IP specifies a concrete IP address which should be allocated from a Subnet
	IP *IPAddress `json:"ip,omitempty"`
}

// IP is an IP address.
// +kubebuilder:validation:Type=string
type IPAddress struct {
	netaddr.IP `json:"-"`
}

// MachineAssignmentStatus defines the observed state of Request
type MachineAssignmentStatus struct {
	State     RequestState       `json:"state,omitempty"`
	Reference *ResourceReference `json:"reference,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MachineAssignment is the Schema for the requests API
type MachineAssignment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineAssignmentSpec   `json:"spec,omitempty"`
	Status MachineAssignmentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MachineAssignmentList contains a list of Request
type MachineAssignmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MachineAssignment `json:"items"`
}

// DeepCopyInto is an deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPAddress) DeepCopyInto(out *IPAddress) {
	*out = *in
}
