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

package fake

import (
	"context"
	"errors"
	"fmt"

	ipam "github.com/onmetal/ipam/api/v1alpha1"
	benchv1alpha3 "github.com/onmetal/metal-api/apis/benchmark/v1alpha3"
	inventoryv1alpaha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpaha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/common/types/base"
	"github.com/onmetal/metal-api/pkg/provider"
	oobv1 "github.com/onmetal/oob-controller/api/v1"
	authv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	typesv1 "k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	ExistingMacAddress    = "903cb33381fb"
	ExistingInventoryUUID = "b127954c-3475-22b2-b91c-62d8b4f8cd3f"
	ExistingServerUUID    = "b127954c-3475-22b2-b91c-62d8b4f8cd3f"
	AvailiableServerUUID  = "a132954c-3475-22b2-v20c-74d8b4f8cd4f"
	ExistingOrderName     = "test-order"
)

var (
	errMustNotBeNil                        = errors.New("should not be nil")
	errWrongTypeDefined                    = errors.New("wrong type defined")
	errInventorySpecMuseBeUpdated          = errors.New("inventory spec must be updated")
	errOrderInstanceReferenceMustBeUpdated = errors.New("order instance ref is nil")
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

type fakeClient struct {
	client ctrlclient.Client
	err    error
}

func NewFakeClient() (provider.Client, error) {
	fakeK8sClient, err := newFakeClient()
	if err != nil {
		return nil, err
	}
	return &fakeClient{
		client: fakeK8sClient,
	}, nil
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

func (f *fakeClient) Create(obj any) error {
	switch k8sObject := obj.(type) {
	case ctrlclient.Object:
		return f.client.Create(context.Background(), k8sObject)
	default:
		return fmt.Errorf("%w, %T", errWrongTypeDefined, obj)
	}
}

func (f *fakeClient) Update(obj any) error {
	if f.err != nil {
		return f.err
	}
	if obj == nil {
		return fmt.Errorf("%w", errMustNotBeNil)
	}
	switch v := obj.(type) {
	case *inventoryv1alpaha1.Inventory:
		if v.Spec.System.ID != "newid" {
			return fmt.Errorf("%w", errInventorySpecMuseBeUpdated)
		}
	case *machinev1alpaha2.MachineAssignment:
		if v.Status.MachineRef == nil {
			return fmt.Errorf("%w", errOrderInstanceReferenceMustBeUpdated)
		}
	default:
		return nil
	}
	return nil
}

func (f *fakeClient) Patch(obj any, patch []byte) error {
	switch k8sObject := obj.(type) {
	case ctrlclient.Object:
		return f.client.Patch(context.Background(), k8sObject, ctrlclient.RawPatch(typesv1.MergePatchType, patch))
	default:
		return fmt.Errorf("%w, %T", errWrongTypeDefined, obj)
	}
}

func (f *fakeClient) Get(obj any, sa base.Metadata) error {
	switch k8sObject := obj.(type) {
	case ctrlclient.Object:
		return f.client.Get(
			context.Background(),
			typesv1.NamespacedName{
				Name:      sa.Name(),
				Namespace: sa.Namespace(),
			}, k8sObject)
	default:
		return fmt.Errorf("%w, %T", errWrongTypeDefined, obj)
	}
}

func (f *fakeClient) List(objList any, listOptions *provider.ListOptions) error {
	switch k8sObjectList := objList.(type) {
	case ctrlclient.ObjectList:
		listOptionFilter := &ctrlclient.ListOptions{
			LabelSelector: ctrlclient.MatchingLabelsSelector{Selector: labels.SelectorFromSet(listOptions.Filter)},
		}
		return f.client.List(
			context.Background(),
			k8sObjectList,
			listOptionFilter)
	default:
		return fmt.Errorf("%w, %T", errWrongTypeDefined, objList)
	}
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
				Name:      ExistingInventoryUUID,
				Namespace: namespace,
			},
			Spec: prepareSpecForInventory(ExistingInventoryUUID),
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
					Reference: &machinev1alpaha2.ResourceReference{
						Name:      "non-exist-order",
						Namespace: "default",
					},
				},
				OOB: machinev1alpaha2.ObjectReference{
					Exist: true,
					Reference: &machinev1alpaha2.ResourceReference{
						Name:      ExistingServerUUID,
						Namespace: namespace,
					},
				},
			},
		},
		&machinev1alpaha2.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      AvailiableServerUUID,
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
				OOB: machinev1alpaha2.ObjectReference{
					Exist: true,
					Reference: &machinev1alpaha2.ResourceReference{
						Name:      ExistingServerUUID,
						Namespace: namespace,
					},
				},
			},
		},
		&oobv1.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ExistingServerUUID,
				Namespace: namespace,
			},
			Spec: oobv1.MachineSpec{
				PowerState: "Off",
			},
		},
		&machinev1alpaha2.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "list-test-object",
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
