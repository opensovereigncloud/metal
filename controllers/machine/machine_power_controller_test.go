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
	"github.com/onmetal/metal-api/controllers/scheduler"
	oobv1 "github.com/onmetal/oob-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("machine-power-controller", func() {
	It("Should watch machine objects and turn it on if reservation exists", func() {
		machine := &machinev1alpha2.Machine{}
		oob := &oobv1.OOB{}

		u, err := uuid.NewUUID()
		Expect(err).ToNot(HaveOccurred())
		var (
			name      = u.String()
			namespace = "default"
		)

		// prepare test data
		createOOB(ctx, name, namespace, oob)
		createHealthyRunningMachine(ctx, name, namespace, machine)

		// testing
		By("Expect machine has no reservation")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
				return false
			}

			return machine.Status.Reservation.Reference == nil
		}).Should(BeTrue())

		By("Expect OOB is turned off")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, oob); err != nil {
				return false
			}

			return oob.Spec.Power == "Off"
		}).Should(BeTrue())

		By("Expect machine status reservation updated successfully")
		machine.Status.Reservation = machinev1alpha2.Reservation{
			Status: machine.Status.Reservation.Status,
			Reference: &machinev1alpha2.ResourceReference{
				Name:      name,
				Namespace: namespace,
			},
		}
		Expect(k8sClient.Status().Update(ctx, machine)).To(Succeed())

		By("Expect machine has reservation")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
				return false
			}

			return machine.Status.Reservation.Reference != nil &&
				machine.Status.Reservation.Reference.Name == name &&
				machine.Status.Reservation.Reference.Namespace == namespace
		}).Should(BeTrue())

		By("Expect OOB is turned on")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, oob); err != nil {
				return false
			}

			return oob.Spec.Power == "On"
		}).Should(BeTrue())
	})
})

func createOOB(ctx context.Context, name, namespace string, oob *oobv1.OOB) {
	By("Create OOB")
	Expect(k8sClient.Create(ctx, prepareOOB(name, namespace))).To(Succeed())

	By("Expect OOB created")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, oob); err != nil {
			return false
		}

		return true
	}).Should(BeTrue())
}

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
	machine.Status = prepareMachineStatus(scheduler.ReservationStatusRunning)
	Expect(k8sClient.Status().Update(ctx, machine)).To(Succeed())

	By("Expect there is a healthy machine in running reservation status")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
			return false
		}

		return machine.Status.Reservation.Status == "Running" && machine.Status.Health == "Healthy"
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
