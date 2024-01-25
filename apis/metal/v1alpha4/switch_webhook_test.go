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

			By("Valid inventoried label but empty spec.inventoryRef")
			switchObject.Labels[inventoried] = validInventoried
			switchObject.Spec.InventoryRef = nil
			_, err = switchObject.ValidateCreate()
			Expect(err).To(HaveOccurred())

			By("Valid inventoried label but empty spec.inventoryRef.name")
			switchObject.Labels[inventoried] = validInventoried
			switchObject.SetInventoryRef("")
			_, err = switchObject.ValidateCreate()
			Expect(err).To(HaveOccurred())

			By("Valid inventoried label but invalid spec.inventoryRef.name")
			switchObject.Labels[inventoried] = validInventoried
			switchObject.SetInventoryRef("123-456")
			_, err = switchObject.ValidateCreate()
			Expect(err).To(HaveOccurred())

			By("Invalid inventoried label but valid spec.inventoryRef.name")
			switchObject.Labels[inventoried] = "123"
			switchObject.SetInventoryRef(validInventoryRef)
			_, err = switchObject.ValidateCreate()
			Expect(err).To(HaveOccurred())
		})

		It("Should not update switch with bad labels", func() {
			switchObject, err := createSwitchFromSampleFile()
			Expect(err).To(Succeed())
			switchObject.Namespace = SwitchNamespace
			baseSwitch := switchObject.DeepCopy()

			By("Valid inventoried label but empty spec.inventoryRef")
			switchObject.Labels[inventoried] = validInventoried
			switchObject.Spec.InventoryRef = nil
			_, err = switchObject.ValidateUpdate(baseSwitch)
			Expect(err).To(HaveOccurred())

			By("Valid inventoried label but empty spec.inventoryRef.name")
			switchObject.Labels[inventoried] = validInventoried
			switchObject.SetInventoryRef("")
			_, err = switchObject.ValidateUpdate(baseSwitch)
			Expect(err).To(HaveOccurred())

			By("Valid inventoried label but invalid spec.inventoryRef.name")
			switchObject.Labels[inventoried] = validInventoried
			switchObject.SetInventoryRef("123-456")
			_, err = switchObject.ValidateUpdate(baseSwitch)
			Expect(err).To(HaveOccurred())

			By("Invalid inventoried label but valid spec.inventoryRef.name")
			switchObject.Labels[inventoried] = "123"
			switchObject.SetInventoryRef(validInventoryRef)
			_, err = switchObject.ValidateUpdate(baseSwitch)
			Expect(err).To(HaveOccurred())
		})

		It("Should create switch with good labels", func() {
			switchObject, err := createSwitchFromSampleFile()
			Expect(err).To(Succeed())
			switchObject.Namespace = SwitchNamespace
			switchObject.Labels[inventoried] = validInventoried
			switchObject.SetInventoryRef(validInventoryRef)
			Expect(k8sClient.Create(ctx, switchObject)).To(Succeed())
		})

		It("Should update switch with good labels", func() {
			switchObject, err := createSwitchFromSampleFile()
			Expect(err).To(Succeed())
			switchObject.Namespace = SwitchNamespace
			switchObject.Labels[inventoried] = validInventoried
			switchObject.SetInventoryRef(validInventoryRef)
			labels := switchObject.Labels
			status := switchObject.Status
			Expect(k8sClient.Create(ctx, switchObject)).To(Succeed())
			switchObject.Labels = labels
			switchObject.Status = status
			Expect(k8sClient.Status().Update(ctx, switchObject)).To(Succeed())
			switchObject.Labels[inventoried] = "false"
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
