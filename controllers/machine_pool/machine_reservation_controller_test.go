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
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	computev1alpha1 "github.com/onmetal/onmetal-api/api/compute/v1alpha1"
	corev1alpha1 "github.com/onmetal/onmetal-api/api/core/v1alpha1"
	"github.com/onmetal/onmetal-api/utils/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("MachineReservation-Controller", func() {
	ctx := testing.SetupContext()
	ns := SetupTest(ctx, machineReservationReconcilers)

	It("Should watch compute machine objects and update metal machine reservation", func() {
		machinePool := &computev1alpha1.MachinePool{}
		metalMachine := &machinev1alpha2.Machine{}
		computeMachine := &computev1alpha1.Machine{}

		u, err := uuid.NewUUID()
		Expect(err).ToNot(HaveOccurred())
		var (
			name      = u.String()
			namespace = ns.Name
		)

		// prepare test data
		By("Create healthy running machine")
		createHealthyRunningMachine(ctx, name, namespace, metalMachine)

		By("Create compute machine with machine class")
		createComputeMachine(ctx, name, namespace, computeMachine)

		By("Create machine pool")
		createMachinePool(ctx, name, namespace, machinePool)

		// testing
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

		By("Expect metal machine has reservation by compute machine")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, metalMachine); err != nil {
				return false
			}

			return metalMachine.Status.Reservation.Reference != nil &&
				metalMachine.Status.Reservation.Reference.Name == computeMachine.Name &&
				metalMachine.Status.Reservation.Reference.Namespace == computeMachine.Namespace
		}).Should(BeTrue())

		By("Compute machine is deleted")
		Expect(k8sClient.Delete(ctx, computeMachine)).Should(Succeed())

		By("Expect metal machine has no reservation")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, metalMachine); err != nil {
				return false
			}

			return metalMachine.Status.Reservation.Reference == nil
		}).Should(BeTrue())
	})
})

func createHealthyRunningMachine(ctx context.Context, name, namespace string, machine *machinev1alpha2.Machine) {
	By("Expect successful machine creation")
	Expect(k8sClient.Create(ctx, prepareTestMachineWithSizeLabels(name, namespace))).Should(Succeed())

	By("Expect machine was created")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
			return false
		}

		return true
	}).Should(BeTrue())

	By("Expect successful machine status update")
	machine.Status = prepareMachineStatus(machinev1alpha2.ReservationStatusRunning)
	Expect(k8sClient.Status().Update(ctx, machine)).To(Succeed())

	By("Expect there is a healthy machine in running reservation status")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
			return false
		}

		return machine.Status.Reservation.Status == "Running" && machine.Status.Health == "Healthy"
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
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machineClass); err != nil {
			return false
		}

		return true
	}).Should(BeTrue())

	By("Expect successful compute machine creation")
	Expect(k8sClient.Create(ctx, prepareTestComputeMachine(name, namespace, machineClass))).Should(Succeed())

	By("Expect compute machine was created")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
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
		Spec: computev1alpha1.MachineSpec{MachineClassRef: corev1.LocalObjectReference{Name: machineClass.Name}},
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
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machinePool); err != nil {
			return false
		}

		return true
	}).Should(BeTrue())
}
