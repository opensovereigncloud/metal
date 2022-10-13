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

package common

import "time"

type Metadata interface {
	Name() string
	Namespace() string
	UID() string
	APIVersion() string
	OwnerReferences() []OwnerReference
	SetOwnerReference(reference OwnerReference)
	Labels() map[string]string
	SetNamespace(namespace string)
}

type OwnerReference struct {
	Name       string
	Kind       string
	APIVersion string
	UniqueID   string
}

type ObjectMetadata struct {
	name         string
	namespace    string
	creationTime time.Time
}

func NewObjectMetadata(name, namespace string) *ObjectMetadata {
	return &ObjectMetadata{name: name, namespace: namespace, creationTime: time.Now()}
}

func (m *ObjectMetadata) Name() string {
	return m.name
}

func (m *ObjectMetadata) UID() string {
	return ""
}

func (m *ObjectMetadata) APIVersion() string {
	return ""
}

func (m *ObjectMetadata) OwnerReferences() []OwnerReference {
	return nil
}

func (m *ObjectMetadata) SetOwnerReference(_ OwnerReference) {}

func (m *ObjectMetadata) Namespace() string {
	return m.namespace
}

func (m *ObjectMetadata) Labels() map[string]string {
	return map[string]string{}
}

func (m *ObjectMetadata) SetNamespace(namespace string) {
	m.namespace = namespace
}
