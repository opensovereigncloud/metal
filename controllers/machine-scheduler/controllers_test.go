package controllers

import (
	"context"

	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
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

	By("Expect successful metal request creation")
	Expect(k8sClient.Create(ctx, preparedRequest)).Should(BeNil())

	By("Check machine is reserved")

	var key string
	var ok bool
	machine := &machinev1alpha2.Machine{}
	Eventually(func(g Gomega) bool {
		Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine)).Should(Succeed())

		key, ok = machine.Labels[machinev1alpha2.LeasedLabel]
		if key != "true" && !ok {
			return false
		}
		key, ok = machine.Labels[machinev1alpha2.MetalRequestLabel]
		if key != requestName && !ok {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())

	By("Check request state is pending")

	request := &machinev1alpha2.MachineAssignment{}
	Eventually(func(g Gomega) bool {
		Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: requestName}, request)).Should(Succeed())

		return request.Status.State == machinev1alpha2.RequestStatePending
	}, timeout, interval).Should(BeTrue())

	By("Expect successful machine status update to running")

	machine.Status.Reservation.RequestState = machinev1alpha2.RequestStateRunning
	Eventually(func(g Gomega) error {
		return k8sClient.Status().Update(ctx, machine)
	}, timeout, interval).Should(BeNil())

	By("Check request state is running")

	Eventually(func(g Gomega) bool {
		Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: requestName}, request)).Should(Succeed())

		return request.Status.State == machinev1alpha2.RequestStateRunning
	}, timeout, interval).Should(BeTrue())
}

func prepareMachineForTest(name, namespace string) *machinev1alpha2.Machine {
	return &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"machine.onmetal.de/leased-size":  "m5.metal",
				"machine.onmetal.de/leasing-pool": "metal",
			},
		},
		Spec: machinev1alpha2.MachineSpec{
			InventoryRequested: true,
		},
	}
}

func prepareMachineStatus() machinev1alpha2.MachineStatus {
	return machinev1alpha2.MachineStatus{
		Health:    machinev1alpha2.MachineStateHealthy,
		OOB:       machinev1alpha2.ObjectReference{Exist: true},
		Inventory: machinev1alpha2.ObjectReference{Exist: true},
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
