/*
 * Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha4

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ironcore-dev/metal/pkg/constants"
)

func createSwitchFromSampleFile() (switchObject *NetworkSwitch, err error) {
	samplePath := filepath.Join("..", "..", "..", "config", "samples", "switch_v1beta1_switch.yaml")
	sampleBytes, err := os.ReadFile(samplePath)
	if err != nil {
		return
	}
	yamlRepresentation := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
	switchObject = &NetworkSwitch{}
	err = yamlRepresentation.Decode(switchObject)

	return
}

var _ = Describe("NetworkSwitch Webhook", func() {
	const (
		SwitchNamespace   = "default"
		timeout           = time.Second * 30
		interval          = time.Millisecond * 250
		validInventoryRef = "cdb75442-a44d-4918-b925-b195a3ae0f09"
		validInventoried  = "true"
		inventoryRef      = "metal.ironcore.dev/inventory-ref"
		inventoried       = "metal.ironcore.dev/inventoried"
	)

	AfterEach(func() {
		Expect(k8sClient.DeleteAllOf(ctx, &NetworkSwitch{}, client.InNamespace(SwitchNamespace))).To(Succeed())
		Eventually(
			func() bool {
				list := &NetworkSwitchList{}
				err := k8sClient.List(ctx, list)
				if err != nil || len(list.Items) > 0 {
					return false
				}
				return true
			},
			timeout, interval).Should(BeTrue())
	})

	Context("On NetworkSwitch overrides update", func() {
		It("Should not update overrides on north interfaces", func() {
			switchObject, err := createSwitchFromSampleFile()
			Expect(err).To(Succeed())
			switchObject.Namespace = SwitchNamespace
			switchStatus := switchObject.Status

			Expect(k8sClient.Create(ctx, switchObject)).To(Succeed())

			By("On empty Status")
			Expect(k8sClient.Update(ctx, switchObject)).To(Succeed())

			switchObject.Status = switchStatus
			interfaceStatus := switchObject.Status.Interfaces["Ethernet0"]
			interfaceStatus.SetDirection("north")
			switchObject.Status.Interfaces["Ethernet0"] = interfaceStatus
			Expect(k8sClient.Status().Update(ctx, switchObject)).To(Succeed())

			currentMTU := switchObject.Spec.Interfaces.Overrides[0].GetMTU()
			currentFEC := switchObject.Spec.Interfaces.Overrides[0].GetFEC()
			currentLanes := switchObject.Spec.Interfaces.Overrides[0].GetLanes()
			var newMTU uint32 = 576
			if currentMTU == newMTU {
				newMTU = 577
			}
			var newFEC = "rs"
			if currentFEC == newFEC {
				newFEC = "none"
			}
			var newLanes uint32 = 2
			if currentLanes == newLanes {
				newLanes = 1
			}

			By("On updating MTU")
			switchObject.Spec.Interfaces.Overrides[0].SetMTU(newMTU)
			Expect(k8sClient.Update(ctx, switchObject)).To(HaveOccurred())

			By("On updating FEC")
			switchObject.Spec.Interfaces.Overrides[0].SetMTU(currentMTU)
			switchObject.Spec.Interfaces.Overrides[0].SetFEC(newFEC)
			Expect(k8sClient.Update(ctx, switchObject)).To(HaveOccurred())

			By("On updating Lanes")
			switchObject.Spec.Interfaces.Overrides[0].SetFEC(currentFEC)
			switchObject.Spec.Interfaces.Overrides[0].SetLanes(newLanes)
			Expect(k8sClient.Update(ctx, switchObject)).To(HaveOccurred())

			By("But passing changes on south")
			interfaceStatus = switchObject.Status.Interfaces["Ethernet0"]
			interfaceStatus.SetDirection("south")
			switchObject.Status.Interfaces["Ethernet0"] = interfaceStatus
			Expect(k8sClient.Status().Update(ctx, switchObject)).To(Succeed())
			switchObject.Spec.Interfaces.Overrides[0].SetLanes(newLanes)
			switchObject.Spec.Interfaces.Overrides[0].SetMTU(newMTU)
			switchObject.Spec.Interfaces.Overrides[0].SetFEC(newFEC)
			Expect(k8sClient.Update(ctx, switchObject)).To(Succeed())
		})
	})

	Context("On certain NetworkSwitch labels", func() {
		It("Should not create switch with bad labels", func() {
			switchObject, err := createSwitchFromSampleFile()
			Expect(err).To(Succeed())
			switchObject.Namespace = SwitchNamespace

			By("Empty inventoried label but valid inventory-ref label")
			switchObject.Labels[inventoried] = ""
			switchObject.Labels[inventoryRef] = validInventoryRef
			Expect(k8sClient.Create(ctx, switchObject)).To(HaveOccurred())

			By("Bad value inventoried label but valid inventory-ref label")
			switchObject.Labels[inventoried] = "123"
			switchObject.Labels[inventoryRef] = validInventoryRef
			Expect(k8sClient.Create(ctx, switchObject)).To(HaveOccurred())

			By("Empty inventory-ref label but valid inventoried label")
			switchObject.Labels[inventoried] = validInventoried
			switchObject.Labels[inventoryRef] = ""
			Expect(k8sClient.Create(ctx, switchObject)).To(HaveOccurred())

			By("Bad value inventory-ref label but valid inventoried label")
			switchObject.Labels[inventoried] = validInventoried
			switchObject.Labels[inventoryRef] = "123-456"
			Expect(k8sClient.Create(ctx, switchObject)).To(HaveOccurred())

			By("No inventory-ref label but valid inventoried label")
			switchObject.Labels[inventoried] = validInventoried
			delete(switchObject.Labels, inventoryRef)
			Expect(k8sClient.Create(ctx, switchObject)).To(HaveOccurred())

			By("No inventoried label but valid inventory-ref label")
			switchObject.Labels[inventoryRef] = validInventoryRef
			delete(switchObject.Labels, inventoried)
			Expect(k8sClient.Create(ctx, switchObject)).To(HaveOccurred())
		})

		It("Should not update switch with bad labels", func() {
			switchObject, err := createSwitchFromSampleFile()
			Expect(err).To(Succeed())
			switchObject.Namespace = SwitchNamespace
			switchObject.Labels[inventoried] = validInventoried
			switchObject.Labels[inventoryRef] = validInventoryRef
			labels := switchObject.Labels
			status := switchObject.Status
			Expect(k8sClient.Create(ctx, switchObject)).To(Succeed())
			switchObject.Labels = labels
			switchObject.Status = status
			Expect(k8sClient.Status().Update(ctx, switchObject)).To(Succeed())

			By("Empty inventoried label but valid inventory-ref label")
			switchObject.Labels[inventoried] = ""
			switchObject.Labels[inventoryRef] = validInventoryRef
			Expect(k8sClient.Update(ctx, switchObject)).To(HaveOccurred())

			By("Bad value inventoried label but valid inventory-ref label")
			switchObject.Labels[inventoried] = "123"
			switchObject.Labels[inventoryRef] = validInventoryRef
			Expect(k8sClient.Update(ctx, switchObject)).To(HaveOccurred())

			By("Empty inventory-ref label but valid inventoried label")
			switchObject.Labels[inventoried] = validInventoried
			switchObject.Labels[inventoryRef] = ""
			Expect(k8sClient.Update(ctx, switchObject)).To(HaveOccurred())

			By("Bad value inventory-ref label but valid inventoried label")
			switchObject.Labels[inventoried] = validInventoried
			switchObject.Labels[inventoryRef] = "123-456"
			Expect(k8sClient.Update(ctx, switchObject)).To(HaveOccurred())

			By("No inventory-ref label but valid inventoried label")
			switchObject.Labels[inventoried] = validInventoried
			delete(switchObject.Labels, inventoryRef)
			Expect(k8sClient.Update(ctx, switchObject)).To(HaveOccurred())

			By("No inventoried label but valid inventory-ref label")
			switchObject.Labels[inventoryRef] = validInventoryRef
			delete(switchObject.Labels, inventoried)
			Expect(k8sClient.Update(ctx, switchObject)).To(HaveOccurred())
		})

		It("Should create switch with good labels", func() {
			switchObject, err := createSwitchFromSampleFile()
			Expect(err).To(Succeed())
			switchObject.Namespace = SwitchNamespace
			switchObject.Labels[inventoried] = validInventoried
			switchObject.Labels[inventoryRef] = validInventoryRef
			Expect(k8sClient.Create(ctx, switchObject)).To(Succeed())
		})

		It("Should update switch with good labels", func() {
			switchObject, err := createSwitchFromSampleFile()
			Expect(err).To(Succeed())
			switchObject.Namespace = SwitchNamespace
			switchObject.Labels[inventoried] = validInventoried
			switchObject.Labels[inventoryRef] = validInventoryRef
			labels := switchObject.Labels
			status := switchObject.Status
			Expect(k8sClient.Create(ctx, switchObject)).To(Succeed())
			switchObject.Labels = labels
			switchObject.Status = status
			Expect(k8sClient.Status().Update(ctx, switchObject)).To(Succeed())
			switchObject.Labels[inventoried] = "false"
			switchObject.Labels[inventoryRef] = "e0e223f5-032a-48d7-8481-a828c3cd868a"
			Expect(k8sClient.Update(ctx, switchObject)).To(Succeed())
		})

		It("Should mutate config selector if it is not set", func() {
			switchObject, err := createSwitchFromSampleFile()
			Expect(err).To(Succeed())
			switchObject.Namespace = SwitchNamespace
			Expect(k8sClient.Create(ctx, switchObject)).To(Succeed())
			Expect(switchObject.GetConfigSelector()).NotTo(BeNil())
			Expect(switchObject.GetConfigSelector().MatchLabels[constants.SwitchConfigLayerLabel]).To(Equal("0"))
			switchObject.Spec.ConfigSelector.MatchLabels[constants.SwitchTypeLabel] = "spine"
			Expect(k8sClient.Update(ctx, switchObject)).To(Succeed())
			Expect(switchObject.GetConfigSelector()).NotTo(BeNil())
			Expect(switchObject.GetConfigSelector().MatchLabels[constants.SwitchConfigLayerLabel]).To(BeEmpty())
			Expect(len(switchObject.GetConfigSelector().MatchLabels)).To(Equal(1))
		})

		It("Should bypass config selector mutating in case it is populated", func() {
			switchObject, err := createSwitchFromSampleFile()
			Expect(err).To(Succeed())
			switchObject.Namespace = SwitchNamespace
			switchObject.Spec.ConfigSelector = &metav1.LabelSelector{
				MatchLabels: map[string]string{constants.SwitchTypeLabel: "spine"},
			}
			Expect(k8sClient.Create(ctx, switchObject)).To(Succeed())
			Expect(switchObject.GetConfigSelector()).NotTo(BeNil())
			Expect(len(switchObject.GetConfigSelector().MatchLabels)).To(Equal(1))
			Expect(switchObject.GetConfigSelector().MatchLabels[constants.SwitchConfigLayerLabel]).To(BeEmpty())
		})
	})
})
