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

var _ = Describe("pool-controller", func() {
	ctx := testing.SetupContext()
	ns := SetupTest(ctx)

	Context("Controller Test", func() {
		It("Should watch machine objects and maintain the pool", func() {
			machine := &machinev1alpha2.Machine{}
			pool := &poolv1alpha1.MachinePool{}

			u, err := uuid.NewUUID()
			Expect(err).ToNot(HaveOccurred())
			var (
				name      = u.String()
				namespace = ns.Namespace
			)

			// prepare test data
			createSizes(namespace)
			createAvailableMachine(name, namespace, machine)

			// testing
			By("Pool created")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, pool); err != nil {
					return false
				}

				return true
			}, timeout, interval).Should(BeTrue())

			By("The pool has available machine classes")
			Expect(len(pool.Status.AvailableMachineClasses)).Should(Equal(2))

			By("Available machine classes matched with size labels")
			Expect(func() bool {
				var availableSizeLabels = map[string]string{
					"m5.metal.4cpu": "true",
					"m5.metal.2cpu": "true",
				}

				for _, availableMachineClass := range pool.Status.AvailableMachineClasses {
					if _, ok := availableSizeLabels[availableMachineClass.Name]; !ok {
						return false
					}
				}

				return true
			}()).Should(BeTrue())

			By("Available machine classes do not contain a size label that is not assigned to a machine")
			Expect(func() bool {
				var notAssignedLabel = map[string]string{
					"m5.metal.6cpu": "true",
				}

				for _, availableMachineClass := range pool.Status.AvailableMachineClasses {
					if _, ok := notAssignedLabel[availableMachineClass.Name]; ok {
						return false
					}
				}

				return true
			}()).Should(BeTrue())

			By("Expect successful machine labels update")
			machine.Labels = map[string]string{
				"machine.onmetal.de/size-m5.metal.6cpu": "true",
			}
			Expect(k8sClient.Update(ctx, machine)).To(Succeed())

			By("The available machine classes have been updated following the change in machine labels")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, pool); err != nil {
					return false
				}

				var availableSizeLabels = map[string]string{
					"m5.metal.6cpu": "true",
				}

				for _, availableMachineClass := range pool.Status.AvailableMachineClasses {
					if _, ok := availableSizeLabels[availableMachineClass.Name]; !ok {
						return false
					}
				}

				return true
			}, timeout, interval).Should(BeTrue())

			By("Expect successful machine status update to Running")
			machine.Status = prepareMachineStatus(scheduler.ReservationStatusRunning)
			Expect(k8sClient.Status().Update(ctx, machine)).To(Succeed())

			By("Expect there is machine in running reservation status")
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
					return false
				}

				return machine.Status.Reservation.Status == "Running"
			}, timeout, interval).Should(BeTrue())

			// refresh pool obj
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, pool); err != nil {
					return false
				}

				return true
			}, timeout, interval).Should(BeTrue())

			By("Pool has no available machine classes after machine becomes unavailable")
			Expect(len(pool.Status.AvailableMachineClasses)).Should(Equal(0))

			By("Machine deleted")
			Eventually(func() bool {
				if err := k8sClient.Delete(ctx, machine); err != nil {
					return false
				}

				return true
			}, timeout, interval).Should(BeTrue())

			By("Pool deleted after deleting machine")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, pool)

				return apierrors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
		})
	})
})

// nolint reason:temp
func createAvailableMachine(name, namespace string, machine *machinev1alpha2.Machine) {
	By("Expect successful machine creation")
	Expect(k8sClient.Create(ctx, prepareTestMachineWithSizeLabels(name, namespace))).Should(Succeed())

	By("Expect machine was created and has finalizer")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
			return false
		}

		return controllerutil.ContainsFinalizer(machine, machineFinalizer)
	}, timeout, interval).Should(BeTrue())

	By("Expect successful machine status update")
	machine.Status = prepareMachineStatus(scheduler.ReservationStatusAvailable)
	Expect(k8sClient.Status().Update(ctx, machine)).To(Succeed())

	By("Expect there is machine in available reservation status")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
			return false
		}

		return machine.Status.Reservation.Status == "Available"
	}, timeout, interval).Should(BeTrue())
}

// nolint reason:temp
func prepareTestMachineWithSizeLabels(name, namespace string) *machinev1alpha2.Machine {
	return &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"machine.onmetal.de/size-m5.metal.4cpu": "true",
				"machine.onmetal.de/size-m5.metal.2cpu": "true",
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
		Health:    machinev1alpha2.MachineStateHealthy,
		OOB:       machinev1alpha2.ObjectReference{Exist: true},
		Inventory: machinev1alpha2.ObjectReference{Exist: true},
		Interfaces: []machinev1alpha2.Interface{
			{Name: "test"},
			{Name: "test2"},
		},
		Reservation: machinev1alpha2.Reservation{Status: status},
	}
}

// nolint reason:temp
func createSizes(namespace string) {
	size6cpu := inventoryv1alpha1.Size{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m5.metal.6cpu",
			Namespace: namespace,
		},
	}

	size4cpu := inventoryv1alpha1.Size{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m5.metal.4cpu",
			Namespace: namespace,
		},
	}

	size2cpu := inventoryv1alpha1.Size{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m5.metal.2cpu",
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
}
