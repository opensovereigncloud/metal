/*
Copyright 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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
	"strings"

	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

var newMachineLLDP = inventoriesv1alpha1.LLDPSpec{
	ChassisID:         "1c:34:da:57:3b:44",
	SystemName:        "Fake",
	SystemDescription: "Debian GNU/Linux 10 (buster) Linux 4.19.0-6-2-amd64",
	PortID:            "lan0",
	PortDescription:   "lan0",
	Capabilities:      []inventoriesv1alpha1.LLDPCapabilities{switchv1alpha1.CStationCapability},
}

var newLeafLLDP = inventoriesv1alpha1.LLDPSpec{
	ChassisID:         "68:21:5f:47:0d:5a",
	SystemName:        "spine-1.fra3.infra.onmetal.de",
	SystemDescription: "Debian GNU/Linux 10 (buster) Linux 4.19.0-6-2-amd64 #1 SMP Debian 4.19.67-2+deb10u2 (2019-11-11) x86_64",
	PortID:            "Eth63/1",
	PortDescription:   "Ethernet100",
	Capabilities:      []inventoriesv1alpha1.LLDPCapabilities{switchv1alpha1.CBridgeCapability, switchv1alpha1.CRouterCapability},
}
var newSpineLLDP = inventoriesv1alpha1.LLDPSpec{
	ChassisID:         "68:21:5f:47:11:6e",
	SystemName:        "leaf-1.fra3.infra.onmetal.de",
	SystemDescription: "Debian GNU/Linux 10 (buster) Linux 4.19.0-6-2-amd64 #1 SMP Debian 4.19.67-2+deb10u2 (2019-11-11) x86_64",
	PortID:            "Eth63/1",
	PortDescription:   "Ethernet100",
	Capabilities:      []inventoriesv1alpha1.LLDPCapabilities{switchv1alpha1.CBridgeCapability, switchv1alpha1.CRouterCapability},
}

var _ = Describe("Controllers interaction", func() {
	Context("Processing of switch resources on creation and after", func() {
		It("Should create switches and start to compute configuration", func() {
			By("Start to process switches")
			list := &switchv1alpha1.SwitchList{}
			Eventually(func() bool {
				Expect(k8sClient.List(ctx, list)).Should(Succeed())
				for _, item := range list.Items {
					if string(item.Status.State) == switchv1alpha1.CEmptyString {
						return false
					}
					if strings.HasPrefix(item.Spec.Hostname, "spine") && item.Status.ConnectionLevel != 0 {
						return false
					}
					if strings.HasPrefix(item.Spec.Hostname, "edge-leaf") && item.Status.ConnectionLevel != 1 {
						return false
					}
					if strings.HasPrefix(item.Spec.Hostname, "leaf") && item.Status.ConnectionLevel != 1 {
						return false
					}
					if !item.SubnetsDefined() {
						return false
					}
					if !item.LoopbackAddressesDefined() {
						return false
					}
					if item.Status.State != switchv1alpha1.CSwitchStateReady {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Add LLDP data to emulate switch interconnection")
			spineInv := &inventoriesv1alpha1.Inventory{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: OnmetalNamespace, Name: "a177382d-a3b4-3ecd-97a4-01cc15e749e4"}, spineInv)).To(Succeed())
			for _, nic := range spineInv.Spec.NICs {
				if nic.Name == "Ethernet100" {
					nic.LLDPs = []inventoriesv1alpha1.LLDPSpec{newSpineLLDP}
					break
				}
			}
			Expect(k8sClient.Update(ctx, spineInv)).To(Succeed())
			leafInv := &inventoriesv1alpha1.Inventory{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: OnmetalNamespace, Name: "044ca7d1-c6f8-37d8-83ce-bf6a18318f2d"}, leafInv)).To(Succeed())
			for _, nic := range leafInv.Spec.NICs {
				if nic.Name == "Ethernet100" {
					nic.LLDPs = []inventoriesv1alpha1.LLDPSpec{newLeafLLDP}
					break
				}
			}
			Expect(k8sClient.Update(ctx, leafInv)).To(Succeed())

			spineSw := &switchv1alpha1.Switch{}
			leafSw := &switchv1alpha1.Switch{}
			Eventually(func() bool {
				Expect(k8sClient.Get(ctx, types.NamespacedName{
					Namespace: OnmetalNamespace,
					Name:      "a177382d-a3b4-3ecd-97a4-01cc15e749e4"},
					spineSw)).To(Succeed())
				Expect(k8sClient.Get(ctx, types.NamespacedName{
					Namespace: OnmetalNamespace,
					Name:      "044ca7d1-c6f8-37d8-83ce-bf6a18318f2d"},
					leafSw)).To(Succeed())
				if leafSw.Status.Interfaces["Ethernet100"].Direction != switchv1alpha1.CDirectionNorth {
					return false
				}
				if spineSw.Status.Interfaces["Ethernet100"].Peer == nil {
					return false
				}
				if leafSw.Status.Interfaces["Ethernet100"].Peer == nil {
					return false
				}
				return true
			})

			By("Remove and change LLDP data to emulate peer replacement")
			Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: OnmetalNamespace, Name: "a177382d-a3b4-3ecd-97a4-01cc15e749e4"}, spineInv)).To(Succeed())
			for _, nic := range spineInv.Spec.NICs {
				if nic.Name == "Ethernet100" {
					nic.LLDPs = nil
					break
				}
			}
			Expect(k8sClient.Update(ctx, spineInv)).To(Succeed())
			Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: OnmetalNamespace, Name: "044ca7d1-c6f8-37d8-83ce-bf6a18318f2d"}, leafInv)).To(Succeed())
			for _, nic := range leafInv.Spec.NICs {
				if nic.Name == "Ethernet100" {
					nic.LLDPs = []inventoriesv1alpha1.LLDPSpec{newMachineLLDP}
					break
				}
			}
			Expect(k8sClient.Update(ctx, leafInv)).To(Succeed())

			Eventually(func() bool {
				Expect(k8sClient.Get(ctx, types.NamespacedName{
					Namespace: OnmetalNamespace,
					Name:      "a177382d-a3b4-3ecd-97a4-01cc15e749e4"},
					spineSw)).To(Succeed())
				Expect(k8sClient.Get(ctx, types.NamespacedName{
					Namespace: OnmetalNamespace,
					Name:      "044ca7d1-c6f8-37d8-83ce-bf6a18318f2d"},
					leafSw)).To(Succeed())
				if leafSw.Status.Interfaces["Ethernet100"].Direction != switchv1alpha1.CDirectionSouth {
					return false
				}
				if spineSw.Status.Interfaces["Ethernet100"].Peer != nil {
					return false
				}
				if leafSw.Status.Interfaces["Ethernet100"].Peer == nil {
					return false
				}
				if spineSw.Status.State != switchv1alpha1.CSwitchStateReady {
					return false
				}
				if leafSw.Status.State != switchv1alpha1.CSwitchStateReady {
					return false
				}
				return true
			})
		})
	})
})
