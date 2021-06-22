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

package tests

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	networkglobalv1alpha1 "github.com/onmetal/k8s-network-global/api/v1alpha1"
	subnetv1alpha1 "github.com/onmetal/k8s-subnet/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
)

var _ = Describe("Integration between operators", func() {
	Context("Operators interaction", func() {
		It("Prepare assignments", func() {
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

				swa.Namespace = OnmetalNamespace
				Expect(k8sClient.Create(ctx, swa)).To(Succeed())
				assignment := &switchv1alpha1.SwitchAssignment{}
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{
						Namespace: swa.Namespace,
						Name:      swa.Name,
					}, assignment)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(assignment.Labels).Should(Equal(map[string]string{switchv1alpha1.LabelChassisId: strings.ReplaceAll(assignment.Spec.ChassisID, ":", "-")}))
			}
		})

		It("Prepare networks", func() {
			ng := &networkglobalv1alpha1.NetworkGlobal{
				ObjectMeta: v1.ObjectMeta{
					Name:      UnderlayNetwork,
					Namespace: OnmetalNamespace,
				},
				Spec: networkglobalv1alpha1.NetworkGlobalSpec{
					Description: "test network global",
				},
			}
			Expect(k8sClient.Create(ctx, ng)).To(Succeed())
			subnetCidrV4, err := networkglobalv1alpha1.CIDRFromString(SubnetIPv4CIDR)
			Expect(err).NotTo(HaveOccurred())
			subnetV4 := &subnetv1alpha1.Subnet{
				ObjectMeta: v1.ObjectMeta{
					Name:      SubnetNameV4,
					Namespace: OnmetalNamespace,
				},
				Spec: subnetv1alpha1.SubnetSpec{
					CIDR:              *subnetCidrV4,
					ParentSubnetName:  "",
					NetworkGlobalName: UnderlayNetwork,
					Regions:           []string{TestRegion},
					AvailabilityZones: []string{TestAvailabilityZone},
				},
			}
			Expect(k8sClient.Create(ctx, subnetV4)).To(Succeed())
			createdSubnetV4 := &subnetv1alpha1.Subnet{}
			subnetNamespacedName := types.NamespacedName{
				Namespace: OnmetalNamespace,
				Name:      SubnetNameV4,
			}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, subnetNamespacedName, createdSubnetV4)
				if err != nil {
					return false
				}
				if createdSubnetV4.Status.State != subnetv1alpha1.CFinishedSubnetState {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			subnetCidrV6, err := networkglobalv1alpha1.CIDRFromString(SubnetIPv6CIDR)
			Expect(err).NotTo(HaveOccurred())
			subnetV6 := &subnetv1alpha1.Subnet{
				ObjectMeta: v1.ObjectMeta{
					Name:      SubnetNameV6,
					Namespace: OnmetalNamespace,
				},
				Spec: subnetv1alpha1.SubnetSpec{
					CIDR:              *subnetCidrV6,
					ParentSubnetName:  "",
					NetworkGlobalName: UnderlayNetwork,
					Regions:           []string{TestRegion},
					AvailabilityZones: []string{TestAvailabilityZone},
				},
			}
			Expect(k8sClient.Create(ctx, subnetV6)).To(Succeed())
			createdSubnetV6 := &subnetv1alpha1.Subnet{}
			subnetNamespacedName = types.NamespacedName{
				Namespace: OnmetalNamespace,
				Name:      SubnetNameV6,
			}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, subnetNamespacedName, createdSubnetV6)
				if err != nil {
					return false
				}
				if createdSubnetV6.Status.State != subnetv1alpha1.CFinishedSubnetState {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})

		It("Prepare switches", func() {
			By("Create inventories")
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
					Namespace: OnmetalNamespace,
					Name:      inv.Name,
				}
				inv.Namespace = DefaultNamespace
				Expect(k8sClient.Create(ctx, inv)).To(Succeed())
				sw := &switchv1alpha1.Switch{}
				Eventually(func() bool {
					err := k8sClient.Get(ctx, swNamespacedName, sw)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(sw.Labels).Should(Equal(map[string]string{switchv1alpha1.LabelChassisId: strings.ReplaceAll(sw.Spec.SwitchChassis.ChassisID, ":", "-")}))
			}

			By("Switches reconciliation running")
			Eventually(func() bool {
				list := &switchv1alpha1.SwitchList{}
				Expect(k8sClient.List(ctx, list)).Should(Succeed())
				for _, sw := range list.Items {
					if sw.Spec.State.ConnectionLevel == 255 {
						return false
					}
					if strings.HasPrefix(sw.Spec.Hostname, "spine-0") {
						Expect(sw.Spec.State.ConnectionLevel).ShouldNot(Equal(0))
					}
					if strings.HasPrefix(sw.Spec.Hostname, "spine-1") {
						Expect(sw.Spec.State.ConnectionLevel).ShouldNot(Equal(1))
					}
					if strings.HasPrefix(sw.Spec.Hostname, "leaf") {
						Expect(sw.Spec.State.ConnectionLevel).ShouldNot(Equal(2))
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Subnets defining")
			Eventually(func() bool {
				list := &switchv1alpha1.SwitchList{}
				Expect(k8sClient.List(ctx, list)).Should(Succeed())
				for _, sw := range list.Items {
					if sw.Spec.SouthSubnetV4 == nil || sw.Spec.SouthSubnetV6 == nil {
						return false
					}
					if sw.Spec.SouthSubnetV4.CIDR == "" || sw.Spec.SouthSubnetV6.CIDR == "" {
						return false
					}
					for _, iface := range sw.GetSwitchPorts() {
						if iface.LLDPChassisID != "" {
							if iface.IPv4 == "" || iface.IPv6 == "" {
								return false
							}
						}
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})
	})
})
