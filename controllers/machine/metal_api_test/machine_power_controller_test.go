// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/ironcore-dev/metal/domain/reservation"

	oobv1 "github.com/onmetal/oob-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
)

var _ = Describe("machine-power-controller", func() {
	It("Should watch machine objects and turn it on if reservation exists", func() {
		machine := &metalv1alpha4.Machine{}
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
		machine.Status.Reservation = metalv1alpha4.Reservation{
			Status: machine.Status.Reservation.Status,
			Reference: &metalv1alpha4.ResourceReference{
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

// nolint reason:temp
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

// nolint reason:temp
func createHealthyRunningMachine(ctx context.Context, name, namespace string, machine *metalv1alpha4.Machine) {
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
	machine.Status = prepareMachineStatus(domain.ReservationStatusRunning)
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
func prepareTestMachineWithSizeLabels(name, namespace string) *metalv1alpha4.Machine {
	return &metalv1alpha4.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"metal.ironcore.dev/size-m5-metal-4cpu": "true",
				"metal.ironcore.dev/size-m5-metal-2cpu": "true",
			},
		},
		Spec: metalv1alpha4.MachineSpec{
			InventoryRequested: true,
		},
	}
}

// nolint reason:temp
func prepareMachineStatus(status string) metalv1alpha4.MachineStatus {
	return metalv1alpha4.MachineStatus{
		Health: metalv1alpha4.MachineStateHealthy,
		Network: metalv1alpha4.Network{
			Interfaces: []metalv1alpha4.Interface{
				{Name: "test"},
				{Name: "test2"},
			},
		},
		Reservation: metalv1alpha4.Reservation{Status: status},
	}
}

// nolint reason:temp
func prepareOOB(name, namespace string) *oobv1.OOB {
	return &oobv1.OOB{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: oobv1.OOBSpec{
			Power: "Off",
		},
	}
}
