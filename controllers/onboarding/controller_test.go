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

	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/pkg/constants"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
)

var _ = Describe("machine-controller", func() {
	var (
		name      = "b134567c-2475-2b82-a85c-84d8b4f8cb5a"
		namespace = "default"
	)

	Context("Controllers Test", func() {
		It("Test device onboarding", func() {
			testDeviceOnboarding(name, namespace)
		})
	})
})

func testDeviceOnboarding(name, namespace string) {
	ctx := context.Background()
	inventory := prepareInventory(name, namespace)
	switches := prepareSwitch(namespace)

	By("Expect successful switch creation")
	Expect(k8sClient.Create(ctx, switches)).Should(BeNil())

	By("Expect successful switch status update")
	switches.Status = getSwitchesStatus()
	Expect(k8sClient.Status().Update(ctx, switches)).To(Succeed())

	By("Expect successful inventory creation")
	Expect(k8sClient.Create(ctx, inventory)).Should(BeNil())

	By("Expect successful machine creation")
	machine := &machinev1alpha2.Machine{}
	Eventually(func() bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())

	By("Inspecting machine properties")
	Eventually(func(g Gomega) {
		Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine)).Should(Succeed())

		interfaces := 0
		for i := range machine.Status.Interfaces {
			for nic := range inventory.Spec.NICs {
				if inventory.Spec.NICs[nic].Name != machine.Status.Interfaces[i].Name {
					continue
				}
				interfaces++
			}
		}
		Expect(interfaces).To(Equal(len(machine.Status.Interfaces)))
	}, timeout, interval).Should(Succeed())

	By("Expecting successful switch deletion")
	Expect(k8sClient.Delete(ctx, switches)).Should(Succeed())

	By("Trying machine object deletion")
	Expect(k8sClient.Delete(ctx, machine)).Should(Succeed())
}

// nolint reason:temp
func testServerOnboarding(name, namespace string) {
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
		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
			return false
		}

		key, ok = machine.Labels[machinev1alpha2.LeasedLabel]
		if key != "true" || !ok {
			return false
		}
		key, ok = machine.Labels[machinev1alpha2.MetalAssignmentLabel]
		if key != requestName || !ok {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())

	By("Check request state is pending")

	request := &machinev1alpha2.MachineAssignment{}
	Eventually(func(g Gomega) bool {
		Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: requestName}, request)).Should(Succeed())

		return request.Status.State == "Pending"
	}, timeout, interval).Should(BeTrue())

	By("Expect successful machine status update to running")

	machine.Status.Reservation.Status = "Running"
	Eventually(func(g Gomega) error {
		return k8sClient.Status().Update(ctx, machine)
	}, timeout, interval).Should(BeNil())

	By("Check request state is running")

	Eventually(func(g Gomega) bool {
		Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: requestName}, request)).Should(Succeed())

		return request.Status.State == "Running"
	}, timeout, interval).Should(BeTrue())
}

// nolint reason:temp
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

// nolint reason:temp
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

// nolint reason:temp
func prepareMetalRequest(name, namespace string) *machinev1alpha2.MachineAssignment {
	return &machinev1alpha2.MachineAssignment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: machinev1alpha2.MachineAssignmentSpec{
			MachineSize: "m5.metal",
			Image:       "myimage_repo_location",
		},
	}
}

func prepareInventory(name, namespace string) *inventoriesv1alpha1.Inventory {
	return &inventoriesv1alpha1.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"machine.onmetal.de/size-machine": "true",
			},
		},
		Spec: prepareSpecForInventory(name),
	}
}

func prepareSpecForInventory(name string) inventoriesv1alpha1.InventorySpec {
	return inventoriesv1alpha1.InventorySpec{
		System: &inventoriesv1alpha1.SystemSpec{
			ID: name,
		},
		Host: &inventoriesv1alpha1.HostSpec{
			Name: "node1",
		},
		Blocks: []inventoriesv1alpha1.BlockSpec{},
		Memory: &inventoriesv1alpha1.MemorySpec{},
		CPUs:   []inventoriesv1alpha1.CPUSpec{},
		NICs: []inventoriesv1alpha1.NICSpec{
			{
				Name: "enp0s31f6",
				LLDPs: []inventoriesv1alpha1.LLDPSpec{
					{
						ChassisID:         "3c:2c:99:9d:cd:48",
						SystemName:        "EC1817001226",
						SystemDescription: "ECS4100-52T",
						PortID:            "3c:2c:99:9d:cd:77",
						PortDescription:   "Ethernet100",
					},
				},
			},
			{
				Name: "enp1s32f6",
				LLDPs: []inventoriesv1alpha1.LLDPSpec{
					{
						ChassisID:         "3c:2c:99:9d:cd:48",
						SystemName:        "EC1817001226",
						SystemDescription: "ECS4100-52T",
						PortID:            "3c:2c:99:9d:cd:77",
						PortDescription:   "Ethernet102",
					},
				},
			},
		},
	}
}

func prepareSwitch(namespace string) *switchv1beta1.Switch {
	return &switchv1beta1.Switch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "switch",
			Namespace: namespace,
			Labels: map[string]string{
				constants.LabelChassisID: "3c-2c-99-9d-cd-48",
			},
		},
		Spec: switchv1beta1.SwitchSpec{
			InventoryRef: &v1.LocalObjectReference{Name: "40d13165-d984-331c-bca2-5a52f3072e2b"},
			TopSpine:     pointer.Bool(false),
			Managed:      pointer.Bool(true),
			Cordon:       pointer.Bool(false),
			ScanPorts:    pointer.Bool(true),
		},
	}
}

func getSwitchesStatus() switchv1beta1.SwitchStatus {
	return switchv1beta1.SwitchStatus{
		TotalPorts:  pointer.Uint32(100),
		SwitchPorts: pointer.Uint32(90),
		Role:        pointer.String("leaf"),
		Layer:       pointer.Uint32(2),
		Interfaces: map[string]*switchv1beta1.InterfaceSpec{
			"Ethernet100": {
				MACAddress: pointer.String("3c:2c:99:9d:cd:48"),
				PortParametersSpec: &switchv1beta1.PortParametersSpec{
					FEC:   pointer.String("none"),
					MTU:   pointer.Uint32(1500),
					Lanes: pointer.Uint32(1),
					State: pointer.String("up"),
				},
				Direction: pointer.String("south"),
				Speed:     pointer.Uint32(10000),
				IP: []*switchv1beta1.IPAddressSpec{
					{
						Address:      pointer.String("100.64.4.70/30"),
						ExtraAddress: pointer.Bool(false),
					},
					{
						Address:      pointer.String("64:ff9b:1::220/127"),
						ExtraAddress: pointer.Bool(false),
					},
				},
			},
			"Ethernet102": {
				MACAddress: pointer.String("3c:2c:99:9d:cd:48"),
				PortParametersSpec: &switchv1beta1.PortParametersSpec{
					FEC:   pointer.String("none"),
					MTU:   pointer.Uint32(1500),
					Lanes: pointer.Uint32(1),
					State: pointer.String("up"),
				},
				IP: []*switchv1beta1.IPAddressSpec{
					{
						Address:      pointer.String("100.64.6.70/30"),
						ExtraAddress: pointer.Bool(false),
					},
					{
						Address:      pointer.String("64:ff9b:1::220/127"),
						ExtraAddress: pointer.Bool(false),
					},
				},
				Speed:     pointer.Uint32(10000),
				Direction: pointer.String("south"),
			},
		},
		State: pointer.String("Ready"),
	}
}
