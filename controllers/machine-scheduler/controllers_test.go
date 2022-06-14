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

package controllers

import (
	"context"

	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/internal/entity"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("machine-controller", func() {
	var (
		name      = "b134567c-2475-2b82-a85c-84d8b4f8cb5a"
		namespace = "default"
	)

	Context("Controllers Test", func() {
		It("Test scheduler", func() {
			testScheduler(name, namespace)
		})
	})
})

func testScheduler(name, namespace string) {
	ctx := context.Background()
	requestName := "sample-request"
	preparedRequest := prepareMetalRequest(requestName, namespace)
	preparedMachine := prepareMachineForTest(name, namespace)

	By("Expect successful machine creation")
	Expect(k8sClient.Create(ctx, preparedMachine)).Should(BeNil())

	By("Expect successful machine status update")
	preparedMachine.Status = prepareMachineStatus()
	Expect(k8sClient.Status().Update(ctx, preparedMachine)).To(Succeed())

	By("Expect successful metal assignment creation")
	Expect(k8sClient.Create(ctx, preparedRequest)).Should(BeNil())

	By("Check machine is reserved")

	machine := &machinev1alpha2.Machine{}
	Eventually(func(g Gomega) bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
			return false
		}
		return machine.Status.Reservation.Reference != nil
	}, timeout, interval).Should(BeTrue())

	By("Check request state is reserved")

	request := &machinev1alpha2.MachineAssignment{}
	Eventually(func(g Gomega) bool {
		Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: requestName}, request)).Should(Succeed())
		return request.Status.State == entity.ReservationStatusReserved
	}, timeout, interval).Should(BeTrue())

	By("Expect successful machine status update to running")

	machine.Status.Reservation.Status = entity.ReservationStatusRunning
	Eventually(func(g Gomega) error {
		return k8sClient.Status().Update(ctx, machine)
	}, timeout, interval).Should(BeNil())
}

func prepareMachineForTest(name, namespace string) *machinev1alpha2.Machine {
	return &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"machine.onmetal.de/size-m5.metal": "true",
			},
		},
		Spec: machinev1alpha2.MachineSpec{
			InventoryRequested: true,
		},
	}
}

func prepareMachineStatus() machinev1alpha2.MachineStatus {
	return machinev1alpha2.MachineStatus{
		Health: machinev1alpha2.MachineStateHealthy,
		OOB: machinev1alpha2.ObjectReference{
			Exist: true,
			Reference: &machinev1alpha2.ResourceReference{
				APIVersion: "",
				Kind:       "",
				Name:       "",
				Namespace:  "",
			},
		},
		Inventory: machinev1alpha2.ObjectReference{Exist: true},
		Reservation: machinev1alpha2.Reservation{
			Status:    entity.ReservationStatusAvailable,
			Reference: nil,
		},
		Interfaces: []machinev1alpha2.Interface{
			{Name: "test"},
			{Name: "test2"},
		},
	}
}

func prepareMetalRequest(name, namespace string) *machinev1alpha2.MachineAssignment {
	return &machinev1alpha2.MachineAssignment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: machinev1alpha2.MachineAssignmentSpec{
			MachineClass: v1.LocalObjectReference{
				Name: "m5.metal",
			},
			Image: "myimage_repo_location",
		},
	}
}
