// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	"github.com/google/uuid"
	poolv1alpha1 "github.com/ironcore-dev/ironcore/api/compute/v1alpha1"
	corev1alpha1 "github.com/ironcore-dev/ironcore/api/core/v1alpha1"
	"github.com/ironcore-dev/ironcore/utils/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	controllers "github.com/ironcore-dev/metal/controllers/machine"
	domain "github.com/ironcore-dev/metal/domain/reservation"
)

var _ = PDescribe("MachinePool-Controller", func() {
	ctx := testing.SetupContext()
	ns := SetupTest(ctx, machinePoolReconcilers)

	It("Should watch machine objects and maintain the machinePool", func() {
		machine := &metalv1alpha4.Machine{}
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
			if err := k8sClient.Get(ctx,
				types.NamespacedName{
					Namespace: namespace,
					Name:      name,
				}, machinePool); err != nil {
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
			"metal.ironcore.dev/size-m5-metal-6cpu": "true",
		}
		Expect(k8sClient.Update(ctx, machine)).To(Succeed())

		By("The available machine classes have been updated following the change in machine labels")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx,
				types.NamespacedName{
					Namespace: namespace,
					Name:      name,
				}, machinePool); err != nil {
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

		By("Expect successful machine status update to Reserved")
		machine.Status = prepareMachineStatus(domain.ReservationStatusReserved)
		Expect(k8sClient.Status().Update(ctx, machine)).To(Succeed())

		By("Expect machine status is Reserved")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx,
				types.NamespacedName{
					Namespace: namespace,
					Name:      name,
				}, machine); err != nil {
				return false
			}

			return machine.Status.Reservation.Status == domain.ReservationStatusReserved
		}).Should(BeTrue())

		By("MachinePool has no available machine classes after machine becomes unavailable")
		// refresh MachinePool
		Eventually(func() int {
			if err := k8sClient.Get(ctx,
				types.NamespacedName{
					Namespace: namespace,
					Name:      name,
				}, machinePool); err != nil {
				return 1
			}

			return len(machinePool.Status.AvailableMachineClasses)
		}).Should(Equal(0))

		By("Machine deleted")
		Eventually(func() bool {
			if err := k8sClient.Delete(ctx, machine); err != nil {
				return false
			}

			return true
		}).Should(BeTrue())

		By("MachinePool deleted after deleting machine")
		Eventually(func() bool {
			err := k8sClient.Get(ctx,
				types.NamespacedName{
					Namespace: namespace,
					Name:      name,
				}, machinePool)

			return apierrors.IsNotFound(err)
		}).Should(BeTrue())
	})
})

// nolint reason:temp
func createAvailableMachine(ctx context.Context, name, namespace string, machine *metalv1alpha4.Machine) {
	By("Expect successful machine creation")
	Expect(k8sClient.Create(ctx, prepareTestMachineWithSizeLabels(name, namespace))).Should(Succeed())

	By("Expect machine was created and has finalizer")
	Eventually(func() bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
			return false
		}

		return controllerutil.ContainsFinalizer(machine, controllers.MachineFinalizer)
	}).Should(BeTrue())

	By("Expect successful machine status update")
	machine.Status = prepareMachineStatus(domain.ReservationStatusAvailable)
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
func createSizes(ctx context.Context, namespace string) {
	size6cpu := metalv1alpha4.Size{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m5-metal-6cpu",
			Namespace: namespace,
		},
	}

	size4cpu := metalv1alpha4.Size{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m5-metal-4cpu",
			Namespace: namespace,
		},
	}

	size2cpu := metalv1alpha4.Size{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m5-metal-2cpu",
			Namespace: namespace,
		},
	}

	oobNS := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "oob-",
		},
	}
	Expect(k8sClient.Create(ctx, oobNS)).To(Succeed(), "failed to create oob namespace")

	size2cpuDuplicate := metalv1alpha4.Size{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "m5-metal-2cpu",
			Namespace: oobNS.Name,
		},
	}

	testSizes := []metalv1alpha4.Size{
		size6cpu,
		size4cpu,
		size2cpu,
		size2cpuDuplicate,
	}

	for _, size := range testSizes {
		Expect(k8sClient.Create(ctx, &size)).Should(Succeed())
	}

	Eventually(func() bool {
		list := &metalv1alpha4.SizeList{}

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
