// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
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

package base

type Metadata interface {
	Name() string
	Namespace() string
	SetNamespace(string)
}

type InstanceMetadata struct {
	name      string
	namespace string
}

func NewInstanceMetadata(name, namespace string) *InstanceMetadata {
	return &InstanceMetadata{name: name, namespace: namespace}
}

func (o *InstanceMetadata) Name() string {
	return o.name
}

func (o *InstanceMetadata) Namespace() string {
	return o.namespace
}

func (o *InstanceMetadata) SetNamespace(namespace string) {
	o.namespace = namespace
}
