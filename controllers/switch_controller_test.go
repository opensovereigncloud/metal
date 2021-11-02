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

	subnetv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

var chassisIds = []string{"68:21:5f:47:0d:6e", "68:21:5f:47:0b:6e", "68:21:5f:47:0a:6e"}
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
	SystemName:        "spine-1-1.fra3.infra.onmetal.de",
	SystemDescription: "Debian GNU/Linux 10 (buster) Linux 4.19.0-6-2-amd64 #1 SMP Debian 4.19.67-2+deb10u2 (2019-11-11) x86_64",
	PortID:            "Eth63/1",
	PortDescription:   "Ethernet124",
	Capabilities:      []inventoriesv1alpha1.LLDPCapabilities{switchv1alpha1.CBridgeCapability, switchv1alpha1.CRouterCapability},
}
var newSpineLLDP = inventoriesv1alpha1.LLDPSpec{
	ChassisID:         "68:21:5f:47:11:6e",
	SystemName:        "leaf-1.fra3.infra.onmetal.de",
	SystemDescription: "Debian GNU/Linux 10 (buster) Linux 4.19.0-6-2-amd64 #1 SMP Debian 4.19.67-2+deb10u2 (2019-11-11) x86_64",
	PortID:            "Eth63/1",
	PortDescription:   "Ethernet124",
	Capabilities:      []inventoriesv1alpha1.LLDPCapabilities{switchv1alpha1.CBridgeCapability, switchv1alpha1.CRouterCapability},
}

var _ = Describe("Controllers interaction", func() {
	Context("Processing of switch resources on creation and after", func() {
		AfterEach(func() {
			var err error
			swl := &switchv1alpha1.SwitchList{}
			snl := &subnetv1alpha1.SubnetList{}
			opts := &client.ListOptions{}

			By("Check there are two subnets - v4 and v6 - for every south NIC")
			Eventually(func() bool {
				Expect(k8sClient.List(ctx, swl)).Should(Succeed())
				for _, sw := range swl.Items {
					southNICs := make([]string, 0)
					for nic := range sw.Status.SouthConnections.Peers {
						southNICs = append(southNICs, strings.ToLower(nic))
					}
					if len(southNICs) == 0 {
						return true
					}
					labels := labelsMap{
						include: map[string][]string{
							switchv1alpha1.LabelSwitchName:    {sw.Name},
							switchv1alpha1.LabelInterfaceName: southNICs,
						},
					}
					opts, err = getListFilter(labels)
					Expect(err).NotTo(HaveOccurred())
					Expect(k8sClient.List(ctx, snl, opts)).Should(Succeed())
					if len(snl.Items) != len(southNICs)*2 {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Check there are no subnets for north NICs")
			Eventually(func() bool {
				Expect(k8sClient.List(ctx, swl)).Should(Succeed())
				for _, sw := range swl.Items {
					northNICs := make([]string, 0)
					for nic := range sw.Status.SouthConnections.Peers {
						northNICs = append(northNICs, strings.ToLower(nic))
					}
					if len(northNICs) == 0 {
						return true
					}
					labels := labelsMap{
						include: map[string][]string{
							switchv1alpha1.LabelSwitchName:    {sw.Name},
							switchv1alpha1.LabelInterfaceName: northNICs,
						},
					}
					opts, err = getListFilter(labels)
					Expect(err).NotTo(HaveOccurred())
					Expect(k8sClient.List(ctx, snl, opts)).Should(Succeed())
					if len(snl.Items) != 0 {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})

		It("Should complete switches processing", func() {
			list := &switchv1alpha1.SwitchList{}
			Eventually(func() bool {
				Expect(k8sClient.List(ctx, list)).Should(Succeed())
				for _, sw := range list.Items {
					// check connection levels
					if strings.HasPrefix(sw.Spec.Hostname, "spine-0") {
						if sw.Status.ConnectionLevel != 0 {
							return false
						}
					}
					if strings.HasPrefix(sw.Spec.Hostname, "spine-1") {
						if sw.Status.ConnectionLevel != 1 {
							return false
						}
					}
					if strings.HasPrefix(sw.Spec.Hostname, "leaf") {
						if sw.Status.ConnectionLevel != 2 {
							return false
						}
					}
					if !(sw.SubnetsOk() && sw.AddressesDefined() && sw.AddressesOk(list)) {
						return false
					}
					// check loopback addresses
					if sw.Spec.IPv4 == switchv1alpha1.CEmptyString {
						return false
					}
					if sw.Spec.IPv6 == switchv1alpha1.CEmptyString {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Update leaf inventory by adding new machine LLDP")
			inv := &inventoriesv1alpha1.Inventory{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Namespace: DefaultNamespace,
				Name:      "7db70ddb-f23d-3d67-8b73-fb0dac5216ab",
			}, inv)).Should(Succeed())
			updateInventory(inv, "Ethernet124", newMachineLLDP)
			Expect(k8sClient.Update(ctx, inv)).Should(Succeed())

			By("Should update south peers, interfaces and switch role")
			sw := &switchv1alpha1.Switch{}
			Eventually(func() bool {
				Expect(k8sClient.Get(ctx, types.NamespacedName{
					Namespace: OnmetalNamespace,
					Name:      "7db70ddb-f23d-3d67-8b73-fb0dac5216ab",
				}, sw)).Should(Succeed())
				if _, ok := sw.Status.SouthConnections.Peers["Ethernet124"]; !ok {
					return false
				}
				if sw.Status.State != switchv1alpha1.CSwitchStateReady {
					return false
				}
				if sw.Status.Role != switchv1alpha1.CLeafRole {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			snl := &subnetv1alpha1.SubnetList{}
			Expect(k8sClient.List(ctx, snl))
			Expect(checkNeededSubnetExist(snl, sw.GetInterfaceSubnetName("Ethernet124", subnetv1alpha1.CIPv4SubnetType))).Should(BeTrue())
			Expect(checkNeededSubnetExist(snl, sw.GetInterfaceSubnetName("Ethernet124", subnetv1alpha1.CIPv6SubnetType))).Should(BeTrue())

			By("Update leaf inventory by changing machine LLDP to switch LLDP")
			leafInv := &inventoriesv1alpha1.Inventory{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Namespace: DefaultNamespace,
				Name:      "7db70ddb-f23d-3d67-8b73-fb0dac5216ab",
			}, leafInv)).Should(Succeed())
			updateInventory(leafInv, "Ethernet124", newLeafLLDP)
			Expect(k8sClient.Update(ctx, leafInv)).Should(Succeed())

			By("Update spine inventory by adding LLDP")
			spineInv := &inventoriesv1alpha1.Inventory{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Namespace: DefaultNamespace,
				Name:      "b9a234a5-416b-3d49-a4f8-65b6f30c8ee5",
			}, spineInv)).Should(Succeed())
			updateInventory(spineInv, "Ethernet124", newSpineLLDP)
			Expect(k8sClient.Update(ctx, spineInv)).Should(Succeed())

			By("Should update peers and interfaces")
			leaf := &switchv1alpha1.Switch{}
			spine := &switchv1alpha1.Switch{}
			Eventually(func() bool {
				Expect(k8sClient.Get(ctx, types.NamespacedName{
					Namespace: OnmetalNamespace,
					Name:      "7db70ddb-f23d-3d67-8b73-fb0dac5216ab",
				}, leaf)).Should(Succeed())
				Expect(k8sClient.Get(ctx, types.NamespacedName{
					Namespace: OnmetalNamespace,
					Name:      "b9a234a5-416b-3d49-a4f8-65b6f30c8ee5",
				}, spine)).Should(Succeed())
				if _, ok := spine.Status.SouthConnections.Peers["Ethernet124"]; !ok {
					return false
				}
				if spine.Status.State != switchv1alpha1.CSwitchStateReady {
					return false
				}
				if _, ok := leaf.Status.NorthConnections.Peers["Ethernet124"]; !ok {
					return false
				}
				if leaf.Status.State != switchv1alpha1.CSwitchStateReady {
					return false
				}
				if leaf.Status.Role != switchv1alpha1.CSpineRole {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(k8sClient.List(ctx, snl))
			Expect(checkNeededSubnetExist(snl, leaf.GetInterfaceSubnetName("Ethernet124", subnetv1alpha1.CIPv4SubnetType))).Should(BeFalse())
			Expect(checkNeededSubnetExist(snl, leaf.GetInterfaceSubnetName("Ethernet124", subnetv1alpha1.CIPv6SubnetType))).Should(BeFalse())
			Expect(checkNeededSubnetExist(snl, spine.GetInterfaceSubnetName("Ethernet124", subnetv1alpha1.CIPv4SubnetType))).Should(BeTrue())
			Expect(checkNeededSubnetExist(snl, spine.GetInterfaceSubnetName("Ethernet124", subnetv1alpha1.CIPv6SubnetType))).Should(BeTrue())
		})
	})
})

func checkNeededSubnetExist(list *subnetv1alpha1.SubnetList, name string) bool {
	result := false
	for _, sn := range list.Items {
		if sn.Name == name {
			return true
		}
	}
	return result
}

func updateInventory(inv *inventoriesv1alpha1.Inventory, nicName string, lldp inventoriesv1alpha1.LLDPSpec) {
	updatedInfIndex := 0
	updatedInf := inventoriesv1alpha1.NICSpec{}
	for i, nic := range inv.Spec.NICs.NICs {
		if nic.Name == nicName {
			updatedInfIndex = i
			updatedInf = *nic.DeepCopy()
		}
	}
	updatedInf.NDPs = []inventoriesv1alpha1.NDPSpec{}
	updatedInf.LLDPs = []inventoriesv1alpha1.LLDPSpec{}
	updatedInf.LLDPs = append(updatedInf.LLDPs, lldp)
	inv.Spec.NICs.NICs[updatedInfIndex] = updatedInf
}
