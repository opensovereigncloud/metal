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

package v1beta1

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func createSwitchFromSampleFile() (switchObject *Switch, err error) {
	samplePath := filepath.Join("..", "..", "..", "config", "samples", "switch_v1beta1_switch.yaml")
	sampleBytes, err := ioutil.ReadFile(samplePath)
	if err != nil {
		return
	}
	yamlRepresentation := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
	switchObject = &Switch{}
	err = yamlRepresentation.Decode(switchObject)

	return
}

var _ = Describe("Switch Webhook", func() {
	const (
		SwitchNamespace   = "default"
		timeout           = time.Second * 30
		interval          = time.Millisecond * 250
		validInventoryRef = "cdb75442-a44d-4918-b925-b195a3ae0f09"
		validInventoried  = "true"
		inventoryRef      = "machine.onmetal.de/inventory-ref"
		inventoried       = "machine.onmetal.de/inventoried"
	)

	AfterEach(func() {
		Expect(k8sClient.DeleteAllOf(ctx, &Switch{}, client.InNamespace(SwitchNamespace))).To(Succeed())
		Eventually(
			func() bool {
				list := &SwitchList{}
				err := k8sClient.List(ctx, list)
				if err != nil || len(list.Items) > 0 {
					return false
				}
				return true
			},
			timeout, interval).Should(BeTrue())
	})

	Context("On Switch overrides update", func() {
		It("Should not update overrides on north interfaces", func() {
			switchObject, err := createSwitchFromSampleFile()
			Expect(err).To(Succeed())
			switchObject.Namespace = SwitchNamespace
			switchStatus := switchObject.Status

			Expect(k8sClient.Create(ctx, switchObject)).To(Succeed())

			By("On empty Status")
			Expect(k8sClient.Update(ctx, switchObject)).To(HaveOccurred())

			switchObject.Status = switchStatus
			interfaceStatus := switchObject.Status.Interfaces["Ethernet0"]
			interfaceStatus.Direction = "north"
			switchObject.Status.Interfaces["Ethernet0"] = interfaceStatus
			Expect(k8sClient.Status().Update(ctx, switchObject)).To(Succeed())

			currentMTU := GoUint16(switchObject.Spec.Interfaces.Overrides[0].MTU)
			currentFEC := GoString(switchObject.Spec.Interfaces.Overrides[0].FEC)
			currentLanes := GoUint8(switchObject.Spec.Interfaces.Overrides[0].Lanes)
			var newMTU uint16 = 576
			if currentMTU == newMTU {
				newMTU = 577
			}
			var newFEC = "rs"
			if currentFEC == newFEC {
				newFEC = "none"
			}
			var newLanes uint8 = 2
			if currentLanes == newLanes {
				newLanes = 1
			}

			By("On updating MTU")
			switchObject.Spec.Interfaces.Overrides[0].MTU = &newMTU
			Expect(k8sClient.Update(ctx, switchObject)).To(HaveOccurred())

			By("On updating FEC")
			switchObject.Spec.Interfaces.Overrides[0].MTU = &currentMTU
			switchObject.Spec.Interfaces.Overrides[0].FEC = &newFEC
			Expect(k8sClient.Update(ctx, switchObject)).To(HaveOccurred())

			By("On updating Lanes")
			switchObject.Spec.Interfaces.Overrides[0].FEC = &currentFEC
			switchObject.Spec.Interfaces.Overrides[0].Lanes = &newLanes
			Expect(k8sClient.Update(ctx, switchObject)).To(HaveOccurred())

			By("But passing changes on south")
			interfaceStatus = switchObject.Status.Interfaces["Ethernet0"]
			interfaceStatus.Direction = "south"
			switchObject.Status.Interfaces["Ethernet0"] = interfaceStatus
			Expect(k8sClient.Status().Update(ctx, switchObject)).To(Succeed())
			switchObject.Spec.Interfaces.Overrides[0].Lanes = &newLanes
			switchObject.Spec.Interfaces.Overrides[0].MTU = &newMTU
			switchObject.Spec.Interfaces.Overrides[0].FEC = &newFEC
			Expect(k8sClient.Update(ctx, switchObject)).To(Succeed())
		})
	})

	Context("On certain Switch labels", func() {
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
	})
})
