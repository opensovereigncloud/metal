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

package onboarding

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Onboarding test", func() {
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

	Context("Onboard switch from inventory", func() {
		var sampleInventory *inventoryv1alpha1.Inventory

		JustBeforeEach(func() {
			samplePath := filepath.Join("..", "samples", "inventories", "spine-1.inventory.yaml")
			sampleBytes, err := ioutil.ReadFile(samplePath)
			Expect(err).NotTo(HaveOccurred())
			sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
			sampleInventory = &inventoryv1alpha1.Inventory{}
			Expect(sampleYAML.Decode(sampleInventory)).NotTo(HaveOccurred())
			Expect(k8sClient.Create(ctx, sampleInventory)).NotTo(HaveOccurred())
		})

		It("Should create switch from inventory", func() {
			onboardingLabels := map[string]string{
				"metalapi.onmetal.de/inventoried":   "true",
				"metalapi.onmetal.de/inventory-ref": sampleInventory.Name,
				"metalapi.onmetal.de/chassis-id":    "68-21-5f-47-17-6e",
			}

			onboardedSwitch := &switchv1beta1.Switch{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: onmetal, Name: sampleInventory.Name}, onboardedSwitch)).NotTo(HaveOccurred())
				g.Expect(onboardedSwitch.Labels).NotTo(BeNil())
			}, timeout, interval).Should(Succeed())
			Expect(reflect.DeepEqual(onboardingLabels, onboardedSwitch.Labels)).Should(BeTrue())
		})
	})

	Context("Onboard existing switch", func() {
		var (
			sampleSwitch    *switchv1beta1.Switch
			sampleInventory *inventoryv1alpha1.Inventory
		)
		JustBeforeEach(func() {
			samplePath := filepath.Join("..", "samples", "switches", "spine-1.switch.yaml")
			sampleBytes, err := ioutil.ReadFile(samplePath)
			Expect(err).NotTo(HaveOccurred())
			sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
			sampleSwitch = &switchv1beta1.Switch{}
			Expect(sampleYAML.Decode(sampleSwitch)).NotTo(HaveOccurred())
			Expect(k8sClient.Create(ctx, sampleSwitch)).NotTo(HaveOccurred())
		})

		It("Should update existing switch with proper labels and annotations", func() {
			Expect(sampleSwitch.Labels["metalapi.onmetal.de/inventoried"]).Should(BeEmpty())

			samplePath := filepath.Join("..", "samples", "inventories", "spine-1.inventory.yaml")
			sampleBytes, err := ioutil.ReadFile(samplePath)
			Expect(err).NotTo(HaveOccurred())
			sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
			sampleInventory = &inventoryv1alpha1.Inventory{}
			Expect(sampleYAML.Decode(sampleInventory)).NotTo(HaveOccurred())
			Expect(k8sClient.Create(ctx, sampleInventory)).NotTo(HaveOccurred())

			onboardingLabels := map[string]string{
				"metalapi.onmetal.de/inventoried":   "true",
				"metalapi.onmetal.de/inventory-ref": sampleInventory.Name,
				"metalapi.onmetal.de/chassis-id":    "68-21-5f-47-17-6e",
			}
			onboardingAnnotations := map[string]string{
				switchv1beta1.CHardwareChassisIdAnnotation: strings.ReplaceAll(
					func() string {
						var chassisID string
						for _, nic := range sampleInventory.Spec.NICs {
							if nic.Name == "eth0" {
								chassisID = nic.MACAddress
							}
						}
						return chassisID
					}(), ":", "",
				),
				switchv1beta1.CHardwareSerialAnnotation:       sampleInventory.Spec.System.SerialNumber,
				switchv1beta1.CHardwareManufacturerAnnotation: sampleInventory.Spec.System.Manufacturer,
				switchv1beta1.CHardwareSkuAnnotation:          sampleInventory.Spec.System.ProductSKU,
				switchv1beta1.CSoftwareOnieAnnotation:         "false",
				switchv1beta1.CSoftwareAsicAnnotation:         sampleInventory.Spec.Distro.AsicType,
				switchv1beta1.CSoftwareVersionAnnotation:      sampleInventory.Spec.Distro.CommitId,
				switchv1beta1.CSoftwareOSAnnotation:           "sonic",
				switchv1beta1.CSoftwareHostnameAnnotation:     sampleInventory.Spec.Host.Name,
			}

			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: onmetal, Name: sampleSwitch.Name}, sampleSwitch)).NotTo(HaveOccurred())
				g.Expect(sampleSwitch.Labels).NotTo(BeNil())
				g.Expect(sampleSwitch.Annotations).NotTo(BeNil())
			}, timeout, interval).Should(Succeed())
			Expect(reflect.DeepEqual(onboardingLabels, sampleSwitch.Labels)).Should(BeTrue())
			Expect(reflect.DeepEqual(onboardingAnnotations, sampleSwitch.Annotations)).Should(BeTrue())
		})
	})

	Context("Onboard switch created after inventory reconciliation finished", func() {
		var (
			sampleSwitch    *switchv1beta1.Switch
			sampleInventory *inventoryv1alpha1.Inventory
		)
		JustBeforeEach(func() {
			samplePath := filepath.Join("..", "samples", "inventories", "spine-1.inventory.yaml")
			sampleBytes, err := ioutil.ReadFile(samplePath)
			Expect(err).NotTo(HaveOccurred())
			sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
			sampleInventory = &inventoryv1alpha1.Inventory{}
			Expect(sampleYAML.Decode(sampleInventory)).NotTo(HaveOccurred())
			Expect(k8sClient.Create(ctx, sampleInventory)).NotTo(HaveOccurred())
		})

		It("Should onboard switch via triggering watches", func() {
			onboardingLabels := map[string]string{
				"metalapi.onmetal.de/inventoried":   "true",
				"metalapi.onmetal.de/inventory-ref": sampleInventory.Name,
				"metalapi.onmetal.de/chassis-id":    "68-21-5f-47-17-6e",
			}
			onboardingAnnotations := map[string]string{
				switchv1beta1.CHardwareChassisIdAnnotation: strings.ReplaceAll(
					func() string {
						var chassisID string
						for _, nic := range sampleInventory.Spec.NICs {
							if nic.Name == "eth0" {
								chassisID = nic.MACAddress
							}
						}
						return chassisID
					}(), ":", "",
				),
				switchv1beta1.CHardwareSerialAnnotation:       sampleInventory.Spec.System.SerialNumber,
				switchv1beta1.CHardwareManufacturerAnnotation: sampleInventory.Spec.System.Manufacturer,
				switchv1beta1.CHardwareSkuAnnotation:          sampleInventory.Spec.System.ProductSKU,
				switchv1beta1.CSoftwareOnieAnnotation:         "false",
				switchv1beta1.CSoftwareAsicAnnotation:         sampleInventory.Spec.Distro.AsicType,
				switchv1beta1.CSoftwareVersionAnnotation:      sampleInventory.Spec.Distro.CommitId,
				switchv1beta1.CSoftwareOSAnnotation:           "sonic",
				switchv1beta1.CSoftwareHostnameAnnotation:     sampleInventory.Spec.Host.Name,
			}

			onboardedSwitch := &switchv1beta1.Switch{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: onmetal, Name: sampleInventory.Name}, onboardedSwitch)).NotTo(HaveOccurred())
				g.Expect(onboardedSwitch.Labels).NotTo(BeNil())
			}, timeout, interval).Should(Succeed())
			Expect(reflect.DeepEqual(onboardingLabels, onboardedSwitch.Labels)).Should(BeTrue())
			Expect(k8sClient.Delete(ctx, onboardedSwitch)).NotTo(HaveOccurred())

			samplePath := filepath.Join("..", "samples", "switches", "spine-1.switch.yaml")
			sampleBytes, err := ioutil.ReadFile(samplePath)
			Expect(err).NotTo(HaveOccurred())
			sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
			sampleSwitch = &switchv1beta1.Switch{}
			Expect(sampleYAML.Decode(sampleSwitch)).NotTo(HaveOccurred())
			Expect(k8sClient.Create(ctx, sampleSwitch)).NotTo(HaveOccurred())

			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: onmetal, Name: sampleSwitch.Name}, sampleSwitch)).NotTo(HaveOccurred())
				g.Expect(sampleSwitch.Labels).NotTo(BeNil())
				g.Expect(sampleSwitch.Annotations).NotTo(BeNil())
			}, timeout, interval).Should(Succeed())
			Expect(reflect.DeepEqual(onboardingLabels, sampleSwitch.Labels)).Should(BeTrue())
			Expect(reflect.DeepEqual(onboardingAnnotations, sampleSwitch.Annotations)).Should(BeTrue())
		})
	})
})
