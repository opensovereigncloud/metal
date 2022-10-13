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

package domain

import (
	"github.com/onmetal/metal-api/types/common"
)

type Order interface {
	common.Metadata
	SetInstanceType(string)
	InstanceType() string
}

const (
	OrderStatusAvailable = "Available"
	OrderStatusReserved  = "Reserved"
	OrderStatusPending   = "Pending"
	OrderStatusError     = "Error"
	OrderStatusRunning   = "Running"
)

type OrderEntity struct {
	name         string
	namespace    string
	instanceType string
}

func NewOrder(name, namespace string) *OrderEntity {
	return &OrderEntity{name: name, namespace: namespace}
}

func (o *OrderEntity) Name() string {
	return o.name
}

func (o *OrderEntity) Namespace() string {
	return o.namespace
}

func (o *OrderEntity) UID() string {
	return ""
}

func (o *OrderEntity) APIVersion() string {
	return ""
}

func (o *OrderEntity) OwnerReferences() []common.OwnerReference {
	return nil
}

func (o *OrderEntity) SetOwnerReference(_ common.OwnerReference) {}

func (o *OrderEntity) Labels() map[string]string {
	return map[string]string{}
}

func (o *OrderEntity) SetNamespace(namespace string) {
	o.namespace = namespace
}

func (o *OrderEntity) SetInstanceType(instanceType string) {
	o.instanceType = instanceType
}

func (o *OrderEntity) InstanceType() string {
	return o.instanceType
}
