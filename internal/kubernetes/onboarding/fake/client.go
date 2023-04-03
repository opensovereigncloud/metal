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
	inventories "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machine "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	oob "github.com/onmetal/oob-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func NewFakeClient() (ctrlclient.Client, error) {
	fakeK8sClient, err := newFakeClient()
	if err != nil {
		return nil, err
	}
	return fakeK8sClient, nil
}

func NewFakeWithObjects(objects ...ctrlclient.Object) (ctrlclient.Client, error) {
	fakeBuilder := fake.NewClientBuilder()
	scheme, err := getScheme()
	if err != nil {
		return nil, err
	}

	fakeWithObjects := fakeBuilder.WithScheme(scheme).WithObjects(objects...)
	return fakeWithObjects.Build(), nil
}

func newFakeClient(objects ...ctrlclient.Object) (ctrlclient.Client, error) {
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
	if err := machine.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := inventories.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := oob.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := corev1.AddToScheme(scheme); err != nil {
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
