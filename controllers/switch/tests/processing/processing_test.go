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
	"io/fs"
	"io/ioutil"
	"path/filepath"

	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
				sampleBytes, err := ioutil.ReadFile(samplePath)
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
				sampleBytes, err := ioutil.ReadFile(samplePath)
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
					g.Expect(switchv1beta1.GoString(item.Status.SwitchState.State)).To(Equal(switchv1beta1.CSwitchStateReady))
					g.Expect(item.ConnectionsOK(switchesList)).Should(BeTrue())
				}
			}, timeout, interval).Should(Succeed())
		})
	})
})
