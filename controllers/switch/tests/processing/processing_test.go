/*
Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package processing

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Processing test", func() {
	AfterEach(func() {
		By("Remove inventories if exist")
		Expect(k8sClient.DeleteAllOf(ctx, &inventoryv1alpha1.Inventory{}, client.InNamespace(onmetal))).To(Succeed())
		Eventually(func(g Gomega) {
			list := &inventoryv1alpha1.InventoryList{}
			g.Expect(k8sClient.List(ctx, list)).NotTo(HaveOccurred())
			g.Expect(len(list.Items)).To(Equal(0))
		}, timeout, interval).Should(Succeed())

		By("Remove switches if exist")
		Expect(k8sClient.DeleteAllOf(ctx, &switchv1beta1.Switch{}, client.InNamespace(onmetal))).To(Succeed())
		Eventually(func(g Gomega) {
			list := &switchv1beta1.SwitchList{}
			g.Expect(k8sClient.List(ctx, list)).NotTo(HaveOccurred())
			g.Expect(len(list.Items)).To(Equal(0))
		}, timeout, interval).Should(Succeed())
	}, 60)

	Context("Compute switches' configuration without IPAM", func() {
		JustBeforeEach(func() {
			inventoriesSamples := make([]string, 0)
			Expect(filepath.Walk("../samples/inventories", func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					inventoriesSamples = append(inventoriesSamples, path)
				}
				return nil
			})).NotTo(HaveOccurred())
			for _, samplePath := range inventoriesSamples {
				sampleBytes, err := os.ReadFile(samplePath)
				Expect(err).NotTo(HaveOccurred())
				sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
				sampleInventory := &inventoryv1alpha1.Inventory{}
				Expect(sampleYAML.Decode(sampleInventory)).NotTo(HaveOccurred())
				Expect(k8sClient.Create(ctx, sampleInventory)).NotTo(HaveOccurred())
			}

			switchesSamples := make([]string, 0)
			Expect(filepath.Walk("../samples/switches", func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					switchesSamples = append(switchesSamples, path)
				}
				return nil
			})).NotTo(HaveOccurred())
			for _, samplePath := range switchesSamples {
				sampleBytes, err := os.ReadFile(samplePath)
				Expect(err).NotTo(HaveOccurred())
				sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
				sampleSwitch := &switchv1beta1.Switch{}
				Expect(sampleYAML.Decode(sampleSwitch)).NotTo(HaveOccurred())
				sampleSwitch.Labels = map[string]string{
					"metalapi.onmetal.de/inventoried":   "true",
					"metalapi.onmetal.de/inventory-ref": sampleSwitch.Spec.UUID,
				}
				Expect(k8sClient.Create(ctx, sampleSwitch)).NotTo(HaveOccurred())
			}
		})

		It("Should compute switches' configuration", func() {
			var switchesList = &switchv1beta1.SwitchList{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.List(ctx, switchesList)).NotTo(HaveOccurred())
				for _, item := range switchesList.Items {
					g.Expect(item.Status.SwitchState).NotTo(BeNil())
					g.Expect(switchv1beta1.
						GoString(item.Status.SwitchState.State)).
						To(Equal(switchv1beta1.CSwitchStateReady))
					g.Expect(item.ConnectionsOK(switchesList)).Should(BeTrue())
				}
			}, timeout, interval).Should(Succeed())

			var target = &switchv1beta1.Switch{}

			By("Change topSpine flag to false should cause connection level to change from 0")
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "onmetal",
				Name:      "spine-1"}, target)).Should(Succeed())
			Expect(target.Status.ConnectionLevel).To(Equal(uint8(0)))
			target.Spec.TopSpine = false
			Expect(k8sClient.Update(ctx, target)).Should(Succeed())
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{
					Namespace: "onmetal",
					Name:      "spine-1"}, target)).Should(Succeed())
				g.Expect(target.Status.ConnectionLevel).To(Equal(uint8(2)))
			}, timeout, interval).Should(Succeed())

			By("Change topSpine flag to true should cause connection level to change to 0")
			Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: "onmetal", Name: "spine-1"}, target)).Should(Succeed())
			Expect(target.Status.ConnectionLevel).To(Equal(uint8(2)))
			target.Spec.TopSpine = true
			Expect(k8sClient.Update(ctx, target)).Should(Succeed())
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{
					Namespace: "onmetal",
					Name:      "spine-1"}, target)).Should(Succeed())
				g.Expect(target.Status.ConnectionLevel).To(Equal(uint8(0)))
			}, timeout, interval).Should(Succeed())

			By("Recreation of inventory with changed NICs data should cause switch's interfaces update")
			inventory := &inventoryv1alpha1.Inventory{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Namespace: "onmetal",
				Name:      "a177382d-a3b4-3ecd-97a4-01cc15e749e4"}, inventory)).Should(Succeed())
			Expect(k8sClient.Delete(ctx, inventory)).Should(Succeed())

			updatedInventory := func() *inventoryv1alpha1.Inventory {
				samplePath := filepath.Join("..", "samples", "inventories", "spine-1.inventory.yaml")
				sampleBytes, err := os.ReadFile(samplePath)
				Expect(err).NotTo(HaveOccurred())
				sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
				sampleInventory := &inventoryv1alpha1.Inventory{}
				Expect(sampleYAML.Decode(sampleInventory)).NotTo(HaveOccurred())
				updatedNICs := make([]inventoryv1alpha1.NICSpec, 0)

				for idx, nic := range sampleInventory.Spec.NICs {
					if nic.Name != "Ethernet24" {
						continue
					}
					updatedNICs = append(updatedNICs, sampleInventory.Spec.NICs[:idx]...)
					updatedNICs = append(updatedNICs, sampleInventory.Spec.NICs[idx+1:]...)
					break
				}
				for i := 0; i < 4; i++ {
					nicName := fmt.Sprintf("Ethernet2%d", i)
					updatedNICs = append(updatedNICs, inventoryv1alpha1.NICSpec{
						Name:       nicName,
						Lanes:      1,
						Speed:      25000,
						MTU:        9216,
						ActiveFEC:  "none",
						MACAddress: "68:21:5f:47:17:6e",
					})
				}
				sampleInventory.Spec.NICs = updatedNICs
				return sampleInventory
			}()

			Expect(k8sClient.Create(ctx, updatedInventory)).Should(Succeed())
			spine := &switchv1beta1.Switch{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: "onmetal", Name: "spine-1"}, spine)).Should(Succeed())
				for i := 0; i < 4; i++ {
					nicName := fmt.Sprintf("Ethernet2%d", i)
					data, ok := spine.Status.Interfaces[nicName]
					g.Expect(ok).ShouldNot(BeFalse())
					g.Expect(data).ShouldNot(BeNil())
					g.Expect(data.Lanes).Should(Equal(uint8(1)))
				}
			}, timeout, interval).Should(Succeed())
		})
	})
})
