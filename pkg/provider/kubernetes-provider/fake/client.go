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
	"time"

	ipam "github.com/onmetal/ipam/api/v1alpha1"
	benchv1alpha3 "github.com/onmetal/metal-api/apis/benchmark/v1alpha3"
	inventoryv1alpaha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpaha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/types/common"
	oobv1 "github.com/onmetal/oob-operator/api/v1alpha1"
	authv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	ExistingMacAddress                     = "903cb33381fb"
	ExistingInventoryUUID                  = "b127954c-3475-22b2-b91c-62d8b4f8cd3f"
	DeleteInventoryUUID                    = "c235954c-3475-22b2-b91c-23d8b4f8cd1a"
	NotAccomplishedInventoryUUID           = "a137954c-3475-22b2-b91c-62d8b4f8cd4f"
	ExistingServerUUID                     = "b127954c-3475-22b2-b91c-62d8b4f8cd3f"
	AvailableInstanceUUID                  = "a132954c-3475-22b2-v20c-74d8b4f8cd4f"
	AvailableSwitchInstanceUUID            = "s245654f-3475-22b2-v20c-12s8b4f8cd4f"
	AvailableInstanceUUIDWithoutServerInfo = "z142934c-3475-22b2-v20c-91d8b4f8cd4a"
	ExistingOrderName                      = "test-order"
)

const (
	benchmarkConfigMap = `benchmarks:
- name: fio-test
  app: fio
  jsonpathInputSelector: "spec.blocks.*.name"
  resources:
    cpu: "8"
  args:
    - '--rw=randread'
    - '--filename={{ inputSelector }}'`
)

func NewFakeClient() (ctrlclient.Client, error) {
	fakeK8sClient, err := newFakeClient()
	if err != nil {
		return nil, err
	}
	return fakeK8sClient, nil
}

func newFakeClient() (ctrlclient.Client, error) {
	fakeBuilder := fake.NewClientBuilder()
	scheme, err := getScheme()
	if err != nil {
		return nil, err
	}
	objects := getInitObjectList("default")

	fakeWithObjects := fakeBuilder.WithScheme(scheme).WithObjects(objects...)
	return fakeWithObjects.Build(), nil
}

func getScheme() (*k8sRuntime.Scheme, error) {
	scheme := k8sRuntime.NewScheme()
	if err := benchv1alpha3.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := inventoryv1alpaha1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := machinev1alpaha2.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := authv1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := ipam.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := oobv1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := switchv1beta1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return scheme, nil
}

func getInitObjectList(namespace string) []ctrlclient.Object {
	return []ctrlclient.Object{
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "benchmark-config",
				Namespace: namespace,
			},
			Data: map[string]string{"config": benchmarkConfigMap},
		},
		&benchv1alpha3.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ExistingServerUUID,
				Namespace: namespace,
			},
			Spec: benchv1alpha3.MachineSpec{
				Benchmarks: map[string]benchv1alpha3.Benchmarks{"nvme1": {{Name: "test", Value: 1}}},
			},
		},
		&inventoryv1alpaha1.Inventory{
			ObjectMeta: metav1.ObjectMeta{
				Name:      NotAccomplishedInventoryUUID,
				Namespace: namespace,
				Labels: map[string]string{
					inventoryv1alpaha1.GetSizeMatchLabel("dummy"): "true",
				},
				CreationTimestamp: metav1.Time{
					Time: time.Now().Add(-30 * time.Minute)},
			},
			Spec: prepareSpecForInventory(NotAccomplishedInventoryUUID),
		},
		&inventoryv1alpaha1.Inventory{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ExistingInventoryUUID,
				Namespace: namespace,
				Labels: map[string]string{
					inventoryv1alpaha1.GetSizeMatchLabel("machine"): "true",
				},
				CreationTimestamp: metav1.Time{
					Time: time.Now().Add(-30 * time.Minute)},
			},
			Spec: prepareSpecForInventory(ExistingInventoryUUID),
		},
		&switchv1beta1.Switch{
			ObjectMeta: metav1.ObjectMeta{
				Name:      AvailableSwitchInstanceUUID,
				Namespace: namespace,
			},
		},
		&inventoryv1alpaha1.Inventory{
			ObjectMeta: metav1.ObjectMeta{
				Name:      AvailableSwitchInstanceUUID,
				Namespace: namespace,
				Labels: map[string]string{
					inventoryv1alpaha1.GetSizeMatchLabel("switch"): "true",
				},
				CreationTimestamp: metav1.Time{
					Time: time.Now().Add(-30 * time.Minute)},
			},
			Spec: prepareSpecForInventory(AvailableSwitchInstanceUUID),
		},
		&inventoryv1alpaha1.Inventory{
			ObjectMeta: metav1.ObjectMeta{
				Name:      DeleteInventoryUUID,
				Namespace: namespace,
				Labels: map[string]string{
					inventoryv1alpaha1.GetSizeMatchLabel("machine"): "true",
				},
				CreationTimestamp: metav1.Time{
					Time: time.Now().Add(-30 * time.Minute)},
			},
		},
		&ipam.IP{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testip1",
				Namespace: namespace,
				Labels: map[string]string{
					"mac": ExistingMacAddress,
				},
			},
			Spec: ipam.IPSpec{},
		},
		&machinev1alpaha2.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ExistingServerUUID,
				Namespace: namespace,
				Labels: map[string]string{
					"machine.onmetal.de/size-m5-metal": "true",
				},
			},
			Spec: machinev1alpaha2.MachineSpec{},
			Status: machinev1alpaha2.MachineStatus{
				Health: machinev1alpaha2.MachineStateHealthy,
				Reservation: machinev1alpaha2.Reservation{
					Status: "Running",
					Reference: common.ResourceReference{
						Name:      "non-exist-order",
						Namespace: "default",
					},
				},
				OOB: common.ResourceReference{
					Name:      ExistingServerUUID,
					Namespace: namespace,
				},
			},
		},
		&machinev1alpaha2.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      AvailableInstanceUUID,
				Namespace: namespace,
				Labels: map[string]string{
					"machine.onmetal.de/size-m5-metal": "true",
				},
			},
			Spec: machinev1alpaha2.MachineSpec{},
			Status: machinev1alpaha2.MachineStatus{
				Health: machinev1alpaha2.MachineStateHealthy,
				Reservation: machinev1alpaha2.Reservation{
					Status: "Available",
				},
				OOB: common.ResourceReference{
					Name:      ExistingServerUUID,
					Namespace: namespace,
				},
			},
		},
		&oobv1.OOB{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ExistingServerUUID,
				Namespace: namespace,
				Labels: map[string]string{
					"onmetal.de/oob-ignore": "true",
				},
			},
			Spec: oobv1.OOBSpec{
				Power: "Off",
			},
			Status: oobv1.OOBStatus{
				OSMessage:    "Ok",
				Capabilities: []string{"power"},
			},
		},
		&machinev1alpaha2.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      AvailableInstanceUUIDWithoutServerInfo,
				Namespace: namespace,
				Labels: map[string]string{
					"machine.onmetal.de/size-m5-metal": "true",
				},
			},
			Spec: machinev1alpaha2.MachineSpec{},
			Status: machinev1alpaha2.MachineStatus{
				Health: machinev1alpaha2.MachineStateUnhealthy,
			},
		},
		&machinev1alpaha2.MachineAssignment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ExistingOrderName,
				Namespace: namespace,
			},
			Spec: machinev1alpaha2.MachineAssignmentSpec{
				MachineClass: corev1.LocalObjectReference{
					Name: "m5-metal",
				},
				Image: "myimage_repo_location",
			},
		},
	}
}

func prepareSpecForInventory(name string) inventoryv1alpaha1.InventorySpec {
	return inventoryv1alpaha1.InventorySpec{
		System: &inventoryv1alpaha1.SystemSpec{
			ID: name,
		},
		Host: &inventoryv1alpaha1.HostSpec{
			Name: "node1",
		},
		Blocks: []inventoryv1alpaha1.BlockSpec{
			{
				Name:       "nvme1",
				Type:       "",
				Rotational: true,
			},
		},
		Memory: &inventoryv1alpaha1.MemorySpec{},
		CPUs:   []inventoryv1alpaha1.CPUSpec{},
		NICs: []inventoryv1alpaha1.NICSpec{
			{
				Name: "enp0s31f6",
				LLDPs: []inventoryv1alpaha1.LLDPSpec{
					{
						ChassisID:         "3c:2c:99:9d:cd:48",
						SystemName:        "EC1817001226",
						SystemDescription: "ECS4100-52T",
						PortID:            "3c:2c:99:9d:cd:77",
						PortDescription:   "Ethernet100",
					},
				},
			},
			{
				Name: "enp0s32f6",
				LLDPs: []inventoryv1alpaha1.LLDPSpec{
					{
						ChassisID:         "3c:2c:99:9d:cd:48",
						SystemName:        "EC1817001226",
						SystemDescription: "ECS4100-52T",
						PortID:            "3c:2c:99:9d:cd:77",
						PortDescription:   "Ethernet102",
					},
				},
			},
		},
	}
}
