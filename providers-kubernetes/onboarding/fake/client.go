// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// nolint
package fake

import (
	"net/netip"

	ipam "github.com/onmetal/ipam/api/v1alpha1"
	oob "github.com/onmetal/oob-operator/api/v1alpha1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	domain "github.com/ironcore-dev/metal/domain/infrastructure"
	"github.com/ironcore-dev/metal/pkg/constants"
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
			&metalv1alpha4.Machine{}, "metadata.name", machineIndex).
		WithIndex(
			&metalv1alpha4.Inventory{}, "metadata.name", inventoryIndex).
		WithIndex(
			&ipam.IP{}, "metadata.name", ipIndex)

	return fakeWithObjects.Build(), nil
}

func machineIndex(rawObj ctrlclient.Object) []string {
	obj := rawObj.(*metalv1alpha4.Machine)
	return []string{obj.ObjectMeta.Name}
}

func inventoryIndex(rawObj ctrlclient.Object) []string {
	obj := rawObj.(*metalv1alpha4.Inventory)
	return []string{obj.ObjectMeta.Name}
}

func ipIndex(rawObj ctrlclient.Object) []string {
	obj := rawObj.(*ipam.IP)
	return []string{obj.ObjectMeta.Name}
}

func getScheme() (*k8sRuntime.Scheme, error) {
	scheme := k8sRuntime.NewScheme()
	if err := metalv1alpha4.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := ipam.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := oob.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := core.AddToScheme(scheme); err != nil {
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

func InventoryObject(name, namespace string) *metalv1alpha4.Inventory {
	return &metalv1alpha4.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"id": name,
			},
		},
		Spec: metalv1alpha4.InventorySpec{
			System: &metalv1alpha4.SystemSpec{
				ID: "",
			},
			Host: &metalv1alpha4.HostSpec{
				Name: "",
			},
		},
		Status: metalv1alpha4.InventoryStatus{
			InventoryStatuses: metalv1alpha4.InventoryStatuses{
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

func FakeMachineObject(name, namespace string) *metalv1alpha4.Machine {
	return &metalv1alpha4.Machine{
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
	prefix string) *metalv1alpha4.NetworkSwitch {
	ips := make([]*metalv1alpha4.IPAddressSpec, 0)
	ips = append(ips, &metalv1alpha4.IPAddressSpec{Address: &prefix})
	interfaces := make(map[string]*metalv1alpha4.InterfaceSpec)
	interfaces[portDescription] = &metalv1alpha4.InterfaceSpec{IP: ips}
	ready := constants.SwitchStateReady
	return &metalv1alpha4.NetworkSwitch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				constants.LabelChassisID: chassisID,
			},
		},
		Status: metalv1alpha4.NetworkSwitchStatus{
			Interfaces: interfaces,
			State:      &ready,
		},
	}
}
