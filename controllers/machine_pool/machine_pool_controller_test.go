/*
Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	corev1alpha1 "github.com/onmetal/onmetal-api/api/core/v1alpha1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/google/uuid"
	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/controllers/scheduler"
	poolv1alpha1 "github.com/onmetal/onmetal-api/api/compute/v1alpha1"
	"github.com/onmetal/onmetal-api/utils/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var _ = Describe("MachinePool-Controller", func() {
	ctx := testing.SetupContext()
	ns := SetupTest(ctx)

	It("Should watch machine objects and maintain the machinePool", func() {
		machine := &machinev1alpha2.Machine{}
		machinePool := &poolv1alpha1.MachinePool{}

		u, err := uuid.NewUUID()
		Expect(err).ToNot(HaveOccurred())
		var (
			name      = u.String()
			namespace = ns.Name
		)

		// prepare test data
		By("Sizes list created")
		createSizes(ctx, namespace)

		By("Machine classes list created")
		createMachineClasses(ctx, namespace)

		By("Available machine created")
		createAvailableMachine(ctx, name, namespace, machine)

		// testing
		By("MachinePool created")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machinePool); err != nil {
				return false
			}

			return true
		}).Should(BeTrue())

		By("The MachinePool has available machine classes")
		Expect(len(machinePool.Status.AvailableMachineClasses)).Should(Equal(2))

		By("Available machine classes matched with size labels")
		Expect(func() bool {
			var availableSizeLabels = map[string]string{
				"m5-metal-4cpu": "true",
				"m5-metal-2cpu": "true",
			}

			for _, availableMachineClass := range machinePool.Status.AvailableMachineClasses {
				if _, ok := availableSizeLabels[availableMachineClass.Name]; !ok {
					return false
				}
			}

			return true
		}()).Should(BeTrue())

		By("Available machine classes do not contain a size label that is not assigned to a machine")
		Expect(func() bool {
			var notAssignedLabel = map[string]string{
				"m5-metal-6cpu": "true",
			}

			for _, availableMachineClass := range machinePool.Status.AvailableMachineClasses {
				if _, ok := notAssignedLabel[availableMachineClass.Name]; ok {
					return false
				}
			}

			return true
		}()).Should(BeTrue())

		By("Expect successful machine labels update")
		machine.Labels = map[string]string{
			"machine.onmetal.de/size-m5-metal-6cpu": "true",
		}
		Expect(k8sClient.Update(ctx, machine)).To(Succeed())

		By("The available machine classes have been updated following the change in machine labels")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machinePool); err != nil {
				return false
			}

			var availableSizeLabels = map[string]string{
				"m5-metal-6cpu": "true",
			}

			for _, availableMachineClass := range machinePool.Status.AvailableMachineClasses {
				if _, ok := availableSizeLabels[availableMachineClass.Name]; !ok {
					return false
				}
			}

			return true
		}).Should(BeTrue())

		By("Expect successful machine status update to Running")
		machine.Status = prepareMachineStatus(scheduler.ReservationStatusRunning)
		Expect(k8sClient.Status().Update(ctx, machine)).To(Succeed())

		By("Expect there is machine in running reservation status")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
				return false
			}

			return machine.Status.Reservation.Status == "Running"
		}).Should(BeTrue())

		// refresh MachinePool
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machinePool); err != nil {
				return false
			}

			return true
		}).Should(BeTrue())

		By("MachinePool has no available machine classes after machine becomes unavailable")
		Expect(len(machinePool.Status.AvailableMachineClasses)).Should(Equal(0))

		By("Machine deleted")
		Eventually(func() bool {
			if err := k8sClient.Delete(ctx, machine); err != nil {
				return false
			}

			return true
		}).Should(BeTrue())

		By("MachinePool deleted after deleting machine")
		Eventually(func() bool {
			err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machinePool)

			return apierrors.IsNotFound(err)
		}).Should(BeTrue())
	})
})

// nolint reason:temp
func createAvailableMachine(ctx context.Context, name, namespace string, machine *machinev1alpha2.Machine) {
	By("Expect successful machine creation")
	Expect(k8sClient.Create(ctx, prepareTestMachineWithSizeLabels(name, namespace))).Should(Succeed())

	By("Expect machine was created and has finalizer")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
			return false
		}

		return controllerutil.ContainsFinalizer(machine, machineFinalizer)
	}).Should(BeTrue())

	By("Expect successful machine status update")
	machine.Status = prepareMachineStatus(scheduler.ReservationStatusAvailable)
	Expect(k8sClient.Status().Update(ctx, machine)).To(Succeed())

	By("Expect there is machine in available reservation status")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
			return false
		}

		return machine.Status.Reservation.Status == "Available"
	}).Should(BeTrue())
}

// nolint reason:temp
func prepareTestMachineWithSizeLabels(name, namespace string) *machinev1alpha2.Machine {
	return &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"machine.onmetal.de/size-m5-metal-4cpu": "true",
				"machine.onmetal.de/size-m5-metal-2cpu": "true",
			},
		},
		Spec: machinev1alpha2.MachineSpec{
			InventoryRequested: true,
		},
	}
}

// nolint reason:temp
func prepareMachineStatus(status string) machinev1alpha2.MachineStatus {
	return machinev1alpha2.MachineStatus{
		Health: machinev1alpha2.MachineStateHealthy,
		Interfaces: []machinev1alpha2.Interface{
			{Name: "test"},
			{Name: "test2"},
		},
		Reservation: machinev1alpha2.Reservation{Status: status},
	}
}

// nolint reason:temp
func createSizes(ctx context.Context, namespace string) {
	size6cpu := inventoryv1alpha1.Size{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m5-metal-6cpu",
			Namespace: namespace,
		},
	}

	size4cpu := inventoryv1alpha1.Size{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m5-metal-4cpu",
			Namespace: namespace,
		},
	}

	size2cpu := inventoryv1alpha1.Size{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m5-metal-2cpu",
			Namespace: namespace,
		},
	}

	testSizes := []inventoryv1alpha1.Size{
		size6cpu,
		size4cpu,
		size2cpu,
	}

	for _, size := range testSizes {
		Expect(k8sClient.Create(ctx, &size)).Should(Succeed())
	}

	Eventually(func() bool {
		list := &inventoryv1alpha1.SizeList{}

		if err := k8sClient.List(ctx, list); err != nil {
			return false
		}

		return len(list.Items) > 0
	}).Should(BeTrue())
}

// nolint reason:temp
func createMachineClasses(ctx context.Context, namespace string) {
	class6cpu := poolv1alpha1.MachineClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m5-metal-6cpu",
			Namespace: namespace,
		},
		Capabilities: corev1alpha1.ResourceList{
			corev1alpha1.ResourceCPU:    resource.MustParse("6"),
			corev1alpha1.ResourceMemory: resource.MustParse("16Gi"),
		},
	}

	class4cpu := poolv1alpha1.MachineClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m5-metal-4cpu",
			Namespace: namespace,
		},
		Capabilities: corev1alpha1.ResourceList{
			corev1alpha1.ResourceCPU:    resource.MustParse("4"),
			corev1alpha1.ResourceMemory: resource.MustParse("16Gi"),
		},
	}

	class2cpu := poolv1alpha1.MachineClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m5-metal-2cpu",
			Namespace: namespace,
		},
		Capabilities: corev1alpha1.ResourceList{
			corev1alpha1.ResourceCPU:    resource.MustParse("2"),
			corev1alpha1.ResourceMemory: resource.MustParse("16Gi"),
		},
	}

	testMachineClasses := []poolv1alpha1.MachineClass{
		class6cpu,
		class4cpu,
		class2cpu,
	}

	for _, machineClass := range testMachineClasses {
		Expect(k8sClient.Create(ctx, &machineClass)).Should(Succeed())
	}

	Eventually(func() bool {
		list := &poolv1alpha1.MachineClassList{}

		if err := k8sClient.List(ctx, list); err != nil {
			return false
		}

		return len(list.Items) > 0
	}).Should(BeTrue())
}
