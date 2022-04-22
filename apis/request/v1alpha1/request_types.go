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

package v1alpha1

import (
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"inet.af/netaddr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// A toleration operator is the set of operators that can be used in a toleration.
type TolerationOperator string

const (
	TolerationOpEqual  TolerationOperator = "Equal"
	TolerationOpExists TolerationOperator = "Exists"
)

type RequestState string

const (
	RequestStateReserved RequestState = "Reserved"
	RequestStatePending  RequestState = "Pending"
	RequestStateError    RequestState = "Error"
	RequestStateRunning  RequestState = "Running"
)

type RequestKind string

const (
	Machine RequestKind = "Machine"
	Switch  RequestKind = "Switch"
	Router  RequestKind = "Router"
)

// RequestSpec defines the desired state of Request
type RequestSpec struct {
	// Hostname is the hostname of the machine
	Hostname string `json:"hostname"`
	// Kind defines request server type. Machine or Switch or Router.
	//+optional
	//+kubebuilder:default:=Machine
	Kind RequestKind `json:"kind,omitempty"`
	// MachineClass is a reference to the machine class/flavor of the machine.
	MachineClass corev1.LocalObjectReference `json:"machineClass"`
	// MachinePoolSelector selects a suitable MachinePool by the given labels.
	MachinePoolSelector map[string]string `json:"machinePoolSelector,omitempty"`
	// MachinePool defines machine pool to run the machine in.
	// If empty, a scheduler will figure out an appropriate pool to run the machine in.
	MachinePool corev1.LocalObjectReference `json:"machinePool,omitempty"`
	// Image is the URL providing the operating system image of the machine.
	Image string `json:"image"`
	// SSHPublicKeys is a list of SSH public key secret references of a machine.
	SSHPublicKeys []ObjectSelector `json:"sshPublicKeys,omitempty"`
	// Interfaces define a list of network interfaces present on the machine
	Interfaces []Interface `json:"interfaces,omitempty"`
	// SecurityGroups is a list of security groups of a machine
	SecurityGroups []corev1.LocalObjectReference `json:"securityGroups,omitempty"`
	// VolumeAttachments are volumes attached to this machine.
	VolumeAttachments []VolumeAttachment `json:"volumeAttachments,omitempty"`
	// Ignition is a reference to a config map containing the ignition YAML for the machine to boot up.
	// If key is empty, DefaultIgnitionKey will be used as fallback.
	Ignition *ObjectSelector `json:"ignition,omitempty"`
	// EFIVars are variables to pass to EFI while booting up.
	EFIVars []EFIVar `json:"efiVars,omitempty"`
	// Tolerations define tolerations the Machine has. Only MachinePools whose taints
	// covered by Tolerations will be considered to run the Machine.
	Tolerations []Toleration `json:"tolerations,omitempty"`
}

// VolumeAttachment defines a volume attachment of a machine
type VolumeAttachment struct {
	// Name is the name of the VolumeAttachment
	Name string `json:"name"`
	// Priority is the OS priority of the volume.
	Priority int32 `json:"priority,omitempty"`
	// VolumeAttachmentSource is the source where the storage for the VolumeAttachment resides at.
	VolumeAttachmentSource `json:",inline"`
}

// VolumeAttachmentSource specifies the source to use for a VolumeAttachment.
type VolumeAttachmentSource struct {
	// VolumeClaim instructs the VolumeAttachment to use a VolumeClaim as source for the attachment.
	VolumeClaim *VolumeClaimAttachmentSource `json:"volumeClaim,omitempty"`
}

// VolumeClaimAttachmentSource references a VolumeClaim as VolumeAttachment source.
type VolumeClaimAttachmentSource struct {
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
type Interface struct {
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

// The resource this Toleration is attached to tolerates any taint that matches
// the triple <key,value,effect> using the matching operator <operator>.
type Toleration struct {
	// Key is the taint key that the toleration applies to. Empty means match all taint keys.
	// If the key is empty, operator must be Exists; this combination means to match all values and all keys.
	Key string `json:"key,omitempty"`
	// Operator represents a key's relationship to the value.
	// Valid operators are Exists and Equal. Defaults to Equal.
	// Exists is equivalent to wildcard for value, so that a resource can
	// tolerate all taints of a particular category.
	Operator TolerationOperator `json:"operator,omitempty"`
	// Value is the taint value the toleration matches to.
	// If the operator is Exists, the value should be empty, otherwise just a regular string.
	Value string `json:"value,omitempty"`
	// Effect indicates the taint effect to match. Empty means match all taint effects.
	// When specified, allowed values are NoSchedule.
	Effect machinev1alpha2.TaintEffect `json:"effect,omitempty"`
}

// RequestStatus defines the observed state of Request
type RequestStatus struct {
	State     RequestState       `json:"state,omitempty"`
	Reference *ResourceReference `json:"reference,omitempty"`
}

// ResourceReference defines related resource info
type ResourceReference struct {
	// APIVersion refers to the resource API version
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind refers to the resource kind
	// +optional
	Kind string `json:"kind,omitempty"`
	// Name refers to the resource name
	// +optional
	Name string `json:"name,omitempty"`
	// Namespace refers to the resource namespace
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Request is the Schema for the requests API
type Request struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RequestSpec   `json:"spec,omitempty"`
	Status RequestStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RequestList contains a list of Request
type RequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Request `json:"items"`
}

// DeepCopyInto is an deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPAddress) DeepCopyInto(out *IPAddress) {
	*out = *in
}
