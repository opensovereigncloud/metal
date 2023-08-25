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

package fake

import (
	inventoryv1alpaha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	authv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
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
	return fakeWithObjects.Build(), nil
}

func getScheme() (*k8sRuntime.Scheme, error) {
	scheme := k8sRuntime.NewScheme()
	if err := inventoryv1alpaha1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := authv1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := rbacv1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return scheme, nil
}

func IPObjectEndpoint(address string, port int32) *corev1.Endpoints {
	return &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubernetes",
			Namespace: "default",
		},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{{
					IP: address,
				}},
				Ports: []corev1.EndpointPort{{
					Name:     "kubernetes",
					Port:     port,
					Protocol: "TCP",
				}},
			},
		},
	}
}

func InventoryObject(name, namespace string) *inventoryv1alpaha1.Inventory {
	return &inventoryv1alpaha1.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				inventoryv1alpaha1.GetSizeMatchLabel("machine"): "true",
			},
		},
		Spec: prepareSpecForInventory(name),
	}
}
func RawEndpoints(address string, port int32) corev1.Endpoints {
	return corev1.Endpoints{
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{{
					IP: address,
				}},
				Ports: []corev1.EndpointPort{{
					Name:     "kubernetes",
					Port:     port,
					Protocol: "TCP",
				}},
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
