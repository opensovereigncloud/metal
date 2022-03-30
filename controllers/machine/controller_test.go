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

	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	switchv1alpha1 "github.com/onmetal/metal-api/apis/switches/v1alpha1"
	oobv1 "github.com/onmetal/oob-controller/api/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("machine-controller", func() {
	var (
		name      = "a969952c-3475-2b82-a85c-84d8b4f8cd2d"
		namespace = "default"
	)

	Context("Controller Test", func() {
		It("Test machine creation and update", func() {
			testMachine(name, namespace)
		})
		It("Test machine onboarding ", func() {
			testMachineOnboarding()
		})
	})
})

func testMachine(name, namespace string) {
	ctx := context.Background()
	preparedMachine := prepareMachineForTest(name, namespace)
	inventory := prepareInventory(name, namespace)
	switches := prepareSwitch(namespace)

	By("Expect successful switch creation")
	Expect(k8sClient.Create(ctx, switches)).Should(BeNil())

	By("Expect successful switch status update")
	switches.Status = getSwitchesStatus()
	Expect(k8sClient.Status().Update(ctx, switches)).Should(BeNil())

	By("Expect successful machine creation")
	Expect(k8sClient.Create(ctx, preparedMachine)).Should(BeNil())

	By("Expect successful inventory creation")
	Expect(k8sClient.Create(ctx, inventory)).Should(BeNil())

	By("Inspecting machine properties")
	machine := &machinev1alpha2.Machine{}
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
		Expect(machine.Status.OOB.Exist).To(BeFalse())
	}, timeout, interval).Should(Succeed())

	By("Expecting successful inventory deletion")
	Expect(k8sClient.Delete(ctx, inventory)).Should(Succeed())

	Eventually(func() bool {
		Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, machine)).Should(Succeed())
		return machine.Status.Inventory.Exist
	}, timeout, interval).Should(BeFalse())

	By("Expecting successful switch deletion")
	Expect(k8sClient.Delete(ctx, switches)).Should(Succeed())

	By("Trying machine object deletion")
	Expect(k8sClient.Delete(ctx, machine)).Should(Succeed())
}

func testMachineOnboarding() {
	var (
		name      = "a237952c-3475-2b82-a85c-84d8b4f8cd2d"
		namespace = "default"
	)
	ctx := context.Background()
	oob := prepareOOB(name, namespace)

	By("Expect successful oob creation")
	err := k8sClient.Create(ctx, oob)
	Expect(err).Should(BeNil())

	By("Expect oob status to be true")

	m := &machinev1alpha2.Machine{}
	Eventually(func(g Gomega) bool {

		if err := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, m); err != nil {
			return false
		}
		return m.Status.OOB.Exist

	}, timeout, interval).Should(BeTrue())

	By("Expect successful oob deletion")
	Expect(k8sClient.Delete(ctx, oob)).Should(BeNil())

	By("Expect oob status to be false")
	Eventually(func(g Gomega) bool {

		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, m)).Should(BeNil())
		return m.Status.OOB.Exist

	}, timeout, interval).Should(BeFalse())

	By("Expect successful machine deletion")
	Expect(k8sClient.Delete(ctx, m)).Should(BeNil())
}

func prepareMachineForTest(name, namespace string) *machinev1alpha2.Machine {
	return &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"test": "test",
			},
		},
		Spec: machinev1alpha2.MachineSpec{
			InventoryRequested: true,
		},
	}
}

func prepareInventory(name, namespace string) *inventoriesv1alpha1.Inventory {
	return &inventoriesv1alpha1.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"machine.onmetal.de/size-machine": "",
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

func prepareSwitch(namespace string) *switchv1alpha1.Switch {
	return &switchv1alpha1.Switch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "switch",
			Namespace: namespace,
			Labels: map[string]string{
				switchv1alpha1.LabelChassisId: "3c-2c-99-9d-cd-48",
			},
		},
		Spec: switchv1alpha1.SwitchSpec{
			Hostname: "switch1",
			Location: &switchv1alpha1.LocationSpec{
				Room: "room1",
				Row:  1,
				Rack: 2,
				HU:   3,
			},
			Chassis: &switchv1alpha1.ChassisSpec{
				ChassisID: "3c:2c:99:9d:cd:48",
			},
			SoftwarePlatform: &switchv1alpha1.SoftwarePlatformSpec{},
		},
	}
}

func getSwitchesStatus() switchv1alpha1.SwitchStatus {
	return switchv1alpha1.SwitchStatus{
		TotalPorts:      100,
		SwitchPorts:     90,
		Role:            "leaf",
		ConnectionLevel: 2,
		Interfaces: map[string]*switchv1alpha1.InterfaceSpec{
			"Ethernet100": {
				MACAddress: "3c:2c:99:9d:cd:48",
				FEC:        "none",
				IPv4: &switchv1alpha1.IPAddressSpec{
					Address: "100.64.4.70/30",
				},
				IPv6: &switchv1alpha1.IPAddressSpec{
					Address: "64:ff9b:1::220/127",
				},
				MTU:       1500,
				Speed:     10000,
				Lanes:     1,
				State:     "up",
				Direction: "north",
			},
			"Ethernet102": {
				MACAddress: "3c:2c:99:9d:cd:48",
				FEC:        "none",
				IPv4: &switchv1alpha1.IPAddressSpec{
					Address: "100.64.6.70/30",
				},
				IPv6: &switchv1alpha1.IPAddressSpec{
					Address: "64:ff9b:1::220/127",
				},
				MTU:       1500,
				Speed:     10000,
				Lanes:     1,
				State:     "up",
				Direction: "north",
			},
		},
		Configuration: &switchv1alpha1.ConfigurationSpec{
			Managed: true,
			State:   "applied",
		},
		State: "ready",
	}
}

func prepareOOB(name, namespace string) *oobv1.Machine {
	return &oobv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: oobv1.MachineSpec{
			UUID:             name,
			PowerState:       "Off",
			ShutdownDeadline: metav1.Now(),
		},
	}
}
