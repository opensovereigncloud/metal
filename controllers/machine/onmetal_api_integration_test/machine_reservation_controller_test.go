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

	"github.com/google/uuid"
	machinev1alpha3 "github.com/onmetal/metal-api/apis/machine/v1alpha3"
	domain "github.com/onmetal/metal-api/domain/reservation"
	computev1alpha1 "github.com/onmetal/onmetal-api/api/compute/v1alpha1"
	corev1alpha1 "github.com/onmetal/onmetal-api/api/core/v1alpha1"
	"github.com/onmetal/onmetal-api/utils/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("MachineReservation-Controller", func() {
	ctx := testing.SetupContext()
	ns := SetupTest(ctx, machineReservationReconcilers)

	It("Should watch compute machine objects and update metal machine reservation", func() {
		machinePool := &computev1alpha1.MachinePool{}
		metalMachine := &machinev1alpha3.Machine{}
		computeMachine := &computev1alpha1.Machine{}
		ipxeCM := &corev1.ConfigMap{}

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
		createAvailableMachine(ctx, name, namespace, metalMachine)

		By("Create compute machine with machine class")
		createComputeMachine(ctx, name, namespace, computeMachine)

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

		By("Expect metal machine has no reservation")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, metalMachine); err != nil {
				return false
			}

			return metalMachine.Status.Reservation.Reference == nil
		}).Should(BeTrue())

		By("Expect compute machine has no machine pool ref")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, computeMachine); err != nil {
				return false
			}

			return computeMachine.Spec.MachinePoolRef == nil
		}).Should(BeTrue())

		By("Expect successful computeMachine.Spec.machinePoolRef update")
		computeMachine.Spec.MachinePoolRef = &corev1.LocalObjectReference{Name: machinePool.Name}
		Expect(k8sClient.Update(ctx, computeMachine)).To(Succeed())

		By("Expect compute machine has machine pool ref")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, computeMachine); err != nil {
				return false
			}

			return computeMachine.Spec.MachinePoolRef != nil
		}).Should(BeTrue())

		By("Expect ipxe configmap was created")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx,
				types.NamespacedName{
					Namespace: namespace,
					Name:      "ipxe-" + metalMachine.Name,
				}, ipxeCM); err != nil {
				return false
			}

			return true
		}).Should(BeTrue())

		By("Expect CM data is valid")
		Eventually(func() bool { return ipxeCM.Data["name"] != "" }).Should(BeTrue())

		By("Expect metal machine has reservation by compute machine")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx,
				types.NamespacedName{
					Namespace: namespace,
					Name:      name,
				}, metalMachine); err != nil {
				return false
			}

			return metalMachine.Status.Reservation.Reference != nil &&
				metalMachine.Status.Reservation.Reference.Name == computeMachine.Name &&
				metalMachine.Status.Reservation.Reference.Namespace == computeMachine.Namespace &&
				metalMachine.Status.Reservation.Status == domain.ReservationStatusReserved
		}).Should(BeTrue())

		By("Expect machine pool does not provide any available machine class after machine reservation")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx,
				types.NamespacedName{
					Namespace: namespace,
					Name:      name,
				}, machinePool); err != nil {
				return false
			}

			return len(machinePool.Status.AvailableMachineClasses) == 0
		}).Should(BeTrue())

		By("Compute machine is deleted")
		Expect(k8sClient.Delete(ctx, computeMachine)).Should(Succeed())

		By("Expect metal machine has no reservation")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, metalMachine); err != nil {
				return false
			}

			return metalMachine.Status.Reservation.Reference == nil && metalMachine.Status.Reservation.Status == "Available"
		}).Should(BeTrue())

		By("Expect machine pool provides available machine classes after machine becomes Available")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machinePool); err != nil {
				return false
			}

			return len(machinePool.Status.AvailableMachineClasses) == 2
		}).Should(BeTrue())

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

		By("Expect ipxe was deleted after compute machine deletion")
		Eventually(func() bool {
			return errors.IsNotFound(
				k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: "ipxe-" + metalMachine.Name}, ipxeCM),
			)
		}).Should(BeTrue())
	})
})

func createHealthyRunningMachine(ctx context.Context, name, namespace string, machine *machinev1alpha3.Machine) {
	By("Expect successful machine creation")
	Expect(k8sClient.Create(ctx, prepareTestMachineWithSizeLabels(name, namespace))).Should(Succeed())

	By("Expect machine was created")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx,
			types.NamespacedName{
				Namespace: namespace,
				Name:      name}, machine,
		); err != nil {
			return false
		}

		return true
	}).Should(BeTrue())

	By("Expect successful machine status update")
	machine.Status = prepareMachineStatus(domain.ReservationStatusRunning)
	Expect(k8sClient.Status().Update(ctx, machine)).To(Succeed())

	By("Expect there is a healthy machine in running reservation status")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
			return false
		}

		return machine.Status.Reservation.Status == domain.ReservationStatusRunning &&
			machine.Status.Health == machinev1alpha3.MachineStateHealthy
	}).Should(BeTrue())
}

func createComputeMachine(ctx context.Context, name, namespace string, machine *computev1alpha1.Machine) {
	By("Expect machine class was created")
	machineClass := &computev1alpha1.MachineClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Capabilities: corev1alpha1.ResourceList{
			corev1alpha1.ResourceCPU:    resource.MustParse("8"),
			corev1alpha1.ResourceMemory: resource.MustParse("1Gi"),
		},
	}
	Expect(k8sClient.Create(ctx, machineClass)).To(Succeed(), "failed to create test machine class")

	By("Expect successful machine class creation")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx,
			types.NamespacedName{
				Namespace: namespace,
				Name:      name},
			machineClass); err != nil {
			return false
		}

		return true
	}).Should(BeTrue())

	By("Expect successful compute machine creation")
	Expect(k8sClient.Create(ctx, prepareTestComputeMachine(name, namespace, machineClass))).Should(Succeed())

	By("Expect compute machine was created")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx,
			types.NamespacedName{
				Namespace: namespace,
				Name:      name,
			}, machine); err != nil {
			return false
		}

		return true
	}).Should(BeTrue())
}

func prepareTestComputeMachine(name, namespace string, machineClass *computev1alpha1.MachineClass) *computev1alpha1.Machine {
	return &computev1alpha1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: computev1alpha1.MachineSpec{
			MachineClassRef: corev1.LocalObjectReference{Name: machineClass.Name},
			Image:           "test-url",
		},
	}
}

func createMachinePool(ctx context.Context, name, namespace string, machinePool *computev1alpha1.MachinePool) {
	By("Expect machine pool was created")
	Expect(k8sClient.Create(ctx, &computev1alpha1.MachinePool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	})).To(Succeed())

	By("Expect successful machine pool creation")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx,
			types.NamespacedName{
				Namespace: namespace,
				Name:      name,
			}, machinePool); err != nil {
			return false
		}

		return true
	}).Should(BeTrue())
}
