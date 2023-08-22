// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

// nolint
package fake

import (
	"net/netip"

	ipam "github.com/onmetal/ipam/api/v1alpha1"
	inventories "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machine "github.com/onmetal/metal-api/apis/machine/v1alpha3"
	switches "github.com/onmetal/metal-api/apis/switch/v1beta1"
	domain "github.com/onmetal/metal-api/domain/infrastructure"
	"github.com/onmetal/metal-api/pkg/constants"
	oob "github.com/onmetal/oob-operator/api/v1alpha1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func NewFakeWithObjects(objects ...ctrlclient.Object) (ctrlclient.Client, error) {
	fakeBuilder := fake.NewClientBuilder()
	scheme, err := getScheme()
	if err != nil {
		return nil, err
	}

	fakeWithObjects := fakeBuilder.WithScheme(scheme).WithObjects(objects...)
	fakeWithObjects = fakeWithObjects.
		WithIndex(
			&machine.Machine{}, "metadata.name", machineIndex).
		WithIndex(
			&inventories.Inventory{}, "metadata.name", inventoryIndex).
		WithIndex(
			&ipam.IP{}, "metadata.name", ipIndex)

	return fakeWithObjects.Build(), nil
}

func machineIndex(rawObj ctrlclient.Object) []string {
	obj := rawObj.(*machine.Machine)
	return []string{obj.ObjectMeta.Name}
}

func inventoryIndex(rawObj ctrlclient.Object) []string {
	obj := rawObj.(*inventories.Inventory)
	return []string{obj.ObjectMeta.Name}
}

func ipIndex(rawObj ctrlclient.Object) []string {
	obj := rawObj.(*ipam.IP)
	return []string{obj.ObjectMeta.Name}
}

func getScheme() (*k8sRuntime.Scheme, error) {
	scheme := k8sRuntime.NewScheme()
	if err := machine.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := ipam.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := inventories.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := oob.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := core.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := switches.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return scheme, nil
}

func getInitObjectList(objects ...ctrlclient.Object) []ctrlclient.Object {
	objectList := make([]ctrlclient.Object, len(objects))
	for o := range objects {
		objectList = append(objectList, objects[o])
	}
	return objectList
}

func InventoryObject(name, namespace string) *inventories.Inventory {
	return &inventories.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"id": name,
			},
		},
		Spec: inventories.InventorySpec{
			System: &inventories.SystemSpec{
				ID: "",
			},
			Host: &inventories.HostSpec{
				Name: "",
			},
		},
		Status: inventories.InventoryStatus{
			InventoryStatuses: inventories.InventoryStatuses{
				RequestsCount: 1,
			},
		},
	}
}

func OOBObject(name, namespace string) *oob.OOB {
	return &oob.OOB{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Status: oob.OOBStatus{
			UUID:         name,
			Capabilities: []string{domain.PowerCapabilities},
		},
	}
}

func FakeMachineObject(name, namespace string) *machine.Machine {
	return &machine.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"id": name,
			},
		},
	}
}

func OOBObjectWithoutPowerCaps(name, namespace string) *oob.OOB {
	return &oob.OOB{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Status: oob.OOBStatus{
			UUID: name,
		},
	}
}

func IPIPAMObject(name, namespace string) *ipam.IP {
	return &ipam.IP{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Status: ipam.IPStatus{
			Reserved: &ipam.IPAddr{
				Net: netip.MustParseAddr("192.168.1.1"),
			},
		},
	}
}

func FakeSwitchObject(
	name,
	namespace,
	chassisID string,
	portDescription string,
	prefix string) *switches.Switch {
	ips := make([]*switches.IPAddressSpec, 0)
	ips = append(ips, &switches.IPAddressSpec{Address: &prefix})
	interfaces := make(map[string]*switches.InterfaceSpec)
	interfaces[portDescription] = &switches.InterfaceSpec{IP: ips}
	ready := constants.SwitchStateReady
	return &switches.Switch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				constants.LabelChassisID: chassisID,
			},
		},
		Status: switches.SwitchStatus{
			Interfaces: interfaces,
			State:      &ready,
		},
	}
}
