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

package domain

import (
	machine "github.com/onmetal/metal-api/apis/machine/v1alpha2"
)

type Machine struct {
	UUID         string
	Namespace    string
	SKU          string
	SerialNumber string
	Interfaces   []machine.Interface
	Size         map[string]string
}

func NewMachine(
	UUID string,
	namespace string,
	SKU string,
	serialNumber string,
	interfaces []machine.Interface,
	size map[string]string) Machine {
	return Machine{
		UUID:         UUID,
		Namespace:    namespace,
		SKU:          SKU,
		SerialNumber: serialNumber,
		Interfaces:   interfaces,
		Size:         size}
}

func (m *Machine) MachineSizes(sizes map[string]string) {
	m.Size = sizes
}
