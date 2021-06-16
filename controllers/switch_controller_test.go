/*
Copyright 2021.

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
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

var _ = Describe("Switch controller", func() {
	const (
		DefaultNamespace = "default"
		SwitchNamespace  = "onmetal"
		timeout          = time.Second * 30
		interval         = time.Millisecond * 250
	)

	AfterEach(func() {
		ctx := context.Background()
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
		Expect(k8sClient.DeleteAllOf(ctx, &switchv1alpha1.Switch{}, client.InNamespace(SwitchNamespace))).To(Succeed())
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
		Expect(k8sClient.DeleteAllOf(ctx, &switchv1alpha1.SwitchAssignment{}, client.InNamespace(SwitchNamespace))).To(Succeed())
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
	})

	Context("Switch CR created", func() {
		It("Should get role, connection level and subnet defined", func() {
			By("SwitchAssignment CR installed")
			ctx := context.Background()

			switchAssignmentSamples := []string{
				filepath.Join("..", "config", "samples", "assignment-1.onmetal.de_v1alpha1_switchassignment.yaml"),
				filepath.Join("..", "config", "samples", "assignment-2.onmetal.de_v1alpha1_switchassignment.yaml"),
				filepath.Join("..", "config", "samples", "assignment-3.onmetal.de_v1alpha1_switchassignment.yaml"),
			}

			for _, sample := range switchAssignmentSamples {
				rawInfo := make(map[string]interface{})
				swa := &switchv1alpha1.SwitchAssignment{}
				sampleBytes, err := ioutil.ReadFile(sample)
				Expect(err).NotTo(HaveOccurred())
				err = yaml.Unmarshal(sampleBytes, rawInfo)
				Expect(err).NotTo(HaveOccurred())

				data, err := json.Marshal(rawInfo)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(data, swa)
				Expect(err).NotTo(HaveOccurred())

				swa.Namespace = SwitchNamespace
				Expect(k8sClient.Create(ctx, swa)).To(Succeed())
				createdSwitchAssignment := &switchv1alpha1.SwitchAssignment{}
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{
						Namespace: swa.Namespace,
						Name:      swa.Name,
					}, createdSwitchAssignment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
			}

			By("Switch CR installed")
			switchesSamples := []string{
				filepath.Join("..", "config", "samples", "spine-0-1.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "spine-0-2.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "spine-0-3.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "spine-1-1.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "spine-1-2.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "spine-1-3.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "spine-1-4.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "spine-1-5.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "spine-1-6.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "leaf-1.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "leaf-2.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "leaf-3.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "leaf-4.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "leaf-5.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "leaf-6.onmetal.de_v1alpha1_inventory.yaml"),
				filepath.Join("..", "config", "samples", "leaf-7.onmetal.de_v1alpha1_inventory.yaml"),
			}

			for _, sample := range switchesSamples {
				rawInfo := make(map[string]interface{})
				inv := &inventoriesv1alpha1.Inventory{}
				sampleBytes, err := ioutil.ReadFile(sample)
				Expect(err).NotTo(HaveOccurred())
				err = yaml.Unmarshal(sampleBytes, rawInfo)

				data, err := json.Marshal(rawInfo)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(data, inv)
				Expect(err).NotTo(HaveOccurred())

				swNamespacedName := types.NamespacedName{
					Namespace: SwitchNamespace,
					Name:      inv.Name,
				}
				inv.Namespace = DefaultNamespace
				Expect(k8sClient.Create(ctx, inv)).To(Succeed())
				createdSwitch := &switchv1alpha1.Switch{}
				Eventually(func() bool {
					err := k8sClient.Get(ctx, swNamespacedName, createdSwitch)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
			}

			list := &switchv1alpha1.SwitchList{}
			Eventually(func() bool {
				err := k8sClient.List(ctx, list)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			for _, sw := range list.Items {
				Expect(sw.Spec.State.Role).Should(Equal(switchv1alpha1.CSpineRole))
				if strings.HasPrefix(sw.Spec.Hostname, "spine-0") {
					Expect(sw.Spec.State.ConnectionLevel).Should(Equal(uint8(0)))
				}
				if strings.HasPrefix(sw.Spec.Hostname, "spine-1") {
					Expect(sw.Spec.State.ConnectionLevel).Should(Equal(uint8(0)))
				}
				if strings.HasPrefix(sw.Spec.Hostname, "leaf") {
					Expect(sw.Spec.State.ConnectionLevel).Should(Equal(uint8(2)))
				}
			}
		})
	})
})
