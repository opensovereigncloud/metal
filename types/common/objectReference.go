// /*
// Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package common

func (m *ObjectMetadata) Reference() ResourceReference {
	return ResourceReference{Name: m.Name(), Namespace: m.Namespace()}
}

// ResourceReference defines related resource info.
type ResourceReference struct {
	// Name refers to the resource name
	// +optional
	Name string `json:"name,omitempty"`
	// Namespace refers to the resource namespace
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceReference) DeepCopyInto(out *ResourceReference) {
	*out = *in
	if in != nil {
		in, out := &in, &out
		*out = new(ResourceReference)
		**out = **in
	}
}
