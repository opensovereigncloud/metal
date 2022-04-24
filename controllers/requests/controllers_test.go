package controllers

import (
	"context"

	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	requestv1alpha1 "github.com/onmetal/metal-api/apis/request/v1alpha1"
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

	preparedRequest := prepareMetalRequest(namespace)
	preparedMachine := prepareMachineForTest(name, namespace)
	preparedMachineStatus := prepareMachineStatus()

	By("Expect successful machine creation")
	Expect(k8sClient.Create(ctx, preparedMachine)).Should(BeNil())

	preparedMachine.Status = preparedMachineStatus

	By("Expect successful machine status update")
	Expect(k8sClient.Status().Update(ctx, preparedMachine)).Should(BeNil())

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
		if key != "sample-request" && !ok {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())

	By("Check status is running")

	request := &requestv1alpha1.Request{}
	Eventually(func(g Gomega) bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, request); err != nil {
			return false
		}

		if request.Status.State != machinev1alpha2.RequestStateRunning {
			return false
		}
		return true
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

func prepareMetalRequest(namespace string) *requestv1alpha1.Request {
	return &requestv1alpha1.Request{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-request",
			Namespace: namespace,
		},
		Spec: requestv1alpha1.RequestSpec{
			Hostname: "myhost",
			Kind:     "Machine",
			MachineClass: v1.LocalObjectReference{
				Name: "m5.metal",
			},
			MachinePool: v1.LocalObjectReference{
				Name: "metal",
			},
			Image: "myimage_repo_location",
		},
	}
}
