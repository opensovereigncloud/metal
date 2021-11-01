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
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

var chassisIds = []string{"68:21:5f:47:0d:6e", "68:21:5f:47:0b:6e", "68:21:5f:47:0a:6e"}

var _ = Describe("Controllers interaction", func() {
	Context("Processing of switch resources on creation and after", func() {
		BeforeEach(func() {
			By("Prepare inventories")
			switchesSamples := []string{
				filepath.Join("..", "config", "samples", "inventories", "spine-0-1.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "spine-0-2.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "spine-0-3.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "spine-1-1.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "spine-1-2.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "spine-1-3.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "spine-1-4.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "spine-1-5.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "spine-1-6.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "leaf-1.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "leaf-2.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "leaf-3.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "leaf-4.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "leaf-5.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "leaf-6.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "inventories", "leaf-7.onmetal.de_v1alpha1_inventory.yaml"),
			}
			for _, sample := range switchesSamples {
				rawInfo := make(map[string]interface{})
				inv := &inventoriesv1alpha1.Inventory{}
				sampleBytes, err := ioutil.ReadFile(sample)
				Expect(err).NotTo(HaveOccurred())
				err = yaml.Unmarshal(sampleBytes, rawInfo)
				Expect(err).NotTo(HaveOccurred())
				data, err := json.Marshal(rawInfo)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(data, inv)
				Expect(err).NotTo(HaveOccurred())

				swNamespacedName := types.NamespacedName{
					Namespace: OnmetalNamespace,
					Name:      inv.Name,
				}
				inv.Namespace = DefaultNamespace
				Expect(k8sClient.Create(ctx, inv)).To(Succeed())
				sw := &switchv1alpha1.Switch{}
				Eventually(func() bool {
					err := k8sClient.Get(ctx, swNamespacedName, sw)
					return err == nil
				}, timeout, interval).Should(BeTrue())
				Expect(sw.Labels).Should(Equal(map[string]string{switchv1alpha1.LabelChassisId: switchv1alpha1.MacToLabel(sw.Spec.Chassis.ChassisID)}))
			}

			By("Prepare assignments")
			for _, id := range chassisIds {
				swa := &switchv1alpha1.SwitchAssignment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      getUUID(id),
						Namespace: OnmetalNamespace,
					},
					Spec: switchv1alpha1.SwitchAssignmentSpec{
						ChassisID: id,
						Region: &switchv1alpha1.RegionSpec{
							Name:             TestRegion,
							AvailabilityZone: TestAvailabilityZone,
						},
					},
				}
				Expect(k8sClient.Create(ctx, swa)).To(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, swa.NamespacedName(), swa)
					return err == nil
				}, timeout, interval).Should(BeTrue())
				Expect(swa.Labels).Should(Equal(map[string]string{switchv1alpha1.LabelChassisId: switchv1alpha1.MacToLabel(id)}))
			}

			By("Processing finished")
			list := &switchv1alpha1.SwitchList{}
			Eventually(func() bool {
				Expect(k8sClient.List(ctx, list)).Should(Succeed())
				for _, sw := range list.Items {
					if sw.Status.State != switchv1alpha1.CSwitchStateReady {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})

		AfterEach(func() {
			By("Remove switches")
			Expect(k8sClient.DeleteAllOf(ctx, &switchv1alpha1.Switch{}, client.InNamespace(OnmetalNamespace))).To(Succeed())
			Eventually(func() bool {
				list := &switchv1alpha1.SwitchList{}
				err := k8sClient.List(ctx, list)
				if err != nil {
					return false
				}
				if len(list.Items) > 0 {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Check assignments in pending state")
			list := &switchv1alpha1.SwitchAssignmentList{}
			Eventually(func() bool {
				Expect(k8sClient.List(ctx, list)).Should(Succeed())
				for _, item := range list.Items {
					if item.Status.State != switchv1alpha1.CAssignmentStatePending {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Remove assignments")
			Expect(k8sClient.DeleteAllOf(ctx, &switchv1alpha1.SwitchAssignment{}, client.InNamespace(OnmetalNamespace))).To(Succeed())
			Eventually(func() bool {
				list := &switchv1alpha1.SwitchAssignmentList{}
				err := k8sClient.List(ctx, list)
				if err != nil {
					return false
				}
				if len(list.Items) > 0 {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Remove inventories")
			Expect(k8sClient.DeleteAllOf(ctx, &inventoriesv1alpha1.Inventory{}, client.InNamespace(DefaultNamespace))).To(Succeed())
			Eventually(func() bool {
				list := &inventoriesv1alpha1.InventoryList{}
				err := k8sClient.List(ctx, list)
				if err != nil {
					return false
				}
				if len(list.Items) > 0 {
					return false
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

			By("Update inventory by adding new LLDPs")
			inv := &inventoriesv1alpha1.Inventory{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Namespace: DefaultNamespace,
				Name:      "7db70ddb-f23d-3d67-8b73-fb0dac5216ab",
			}, inv)).Should(Succeed())
			updatedInfIndex := 0
			updatedInf := inventoriesv1alpha1.NICSpec{}
			for i, nic := range inv.Spec.NICs.NICs {
				if nic.Name == "Ethernet124" {
					updatedInfIndex = i
					updatedInf = *nic.DeepCopy()
				}
			}
			updatedInf.NDPs = []inventoriesv1alpha1.NDPSpec{}
			updatedInf.LLDPs = []inventoriesv1alpha1.LLDPSpec{}
			updatedInf.LLDPs = append(updatedInf.LLDPs, inventoriesv1alpha1.LLDPSpec{
				ChassisID:         "1c:34:da:57:3b:44",
				SystemName:        "Fake",
				SystemDescription: "Debian GNU/Linux 10 (buster) Linux 4.19.0-6-2-amd64",
				PortID:            "lan0",
				PortDescription:   "lan0",
				Capabilities:      []inventoriesv1alpha1.LLDPCapabilities{switchv1alpha1.CStationCapability},
			})
			inv.Spec.NICs.NICs[updatedInfIndex] = updatedInf
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
		})
	})
})
