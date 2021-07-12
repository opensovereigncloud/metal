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
	"reflect"
	"strings"

	"github.com/google/uuid"
	subnetv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
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

var _ = Describe("Integration between operators", func() {
	Context("Operators interaction", func() {
		It("Prepare assignments", func() {
			for _, id := range chassisIds[1:] {
				swa := &switchv1alpha1.SwitchAssignment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      getUUID(id),
						Namespace: OnmetalNamespace,
					},
					Spec: switchv1alpha1.SwitchAssignmentSpec{
						ChassisID:        id,
						Region:           TestRegion,
						AvailabilityZone: TestAvailabilityZone,
					},
				}
				Expect(k8sClient.Create(ctx, swa)).To(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, swa.NamespacedName(), swa)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				Expect(swa.Labels).Should(Equal(map[string]string{switchv1alpha1.LabelChassisId: switchv1alpha1.MacToLabel(id)}))
				list := &switchv1alpha1.SwitchAssignmentList{}
				Eventually(func() bool {
					Expect(k8sClient.List(ctx, list)).Should(Succeed())
					for _, item := range list.Items {
						if item.Status.State != switchv1alpha1.StatePending {
							return false
						}
					}
					return true
				}, timeout, interval).Should(BeTrue())
			}
		})

		It("Prepare switches", func() {
			By("Create inventories")
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
				Expect(sw.Labels).Should(Equal(map[string]string{switchv1alpha1.LabelChassisId: switchv1alpha1.MacToLabel(sw.Spec.Chassis.ChassisID)}))
			}
		})

		It("Should update connection levels", func() {
			By("Create new assignment")
			swa := &switchv1alpha1.SwitchAssignment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      getUUID(chassisIds[0]),
					Namespace: OnmetalNamespace,
				},
				Spec: switchv1alpha1.SwitchAssignmentSpec{
					ChassisID:        chassisIds[0],
					Region:           TestRegion,
					AvailabilityZone: TestAvailabilityZone,
				},
			}
			Expect(k8sClient.Create(ctx, swa)).To(Succeed())
			Eventually(func() bool {
				Expect(k8sClient.Get(ctx, swa.NamespacedName(), swa)).Should(Succeed())
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(swa.Labels).Should(Equal(map[string]string{switchv1alpha1.LabelChassisId: switchv1alpha1.MacToLabel(chassisIds[0])}))
			swaList := &switchv1alpha1.SwitchAssignmentList{}
			Expect(k8sClient.List(ctx, swaList)).To(Succeed())

			By("Switches reconciliation in progress")
			list := &switchv1alpha1.SwitchList{}
			Eventually(func() bool {
				Expect(k8sClient.List(ctx, list)).Should(Succeed())
				for _, sw := range list.Items {
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
				}
				return true
			}, timeout, interval).Should(BeTrue())

			list = &switchv1alpha1.SwitchList{}
			Eventually(func() bool {
				Expect(k8sClient.List(ctx, list)).Should(Succeed())
				if !list.AllConnectionsOk() {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			list = &switchv1alpha1.SwitchList{}
			Eventually(func() bool {
				Expect(k8sClient.List(ctx, list)).Should(Succeed())
				for _, sw := range list.Items {
					if sw.Spec.SouthSubnetV4 == nil {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})

		It("Should change status for assignments", func() {
			list := &switchv1alpha1.SwitchAssignmentList{}
			Eventually(func() bool {
				Expect(k8sClient.List(ctx, list)).Should(Succeed())
				for _, item := range list.Items {
					if item.Status.State != switchv1alpha1.StateFinished {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})

		It("Should update south v6 subnet", func() {
			By("Subnets defining")
			cidr, _ := subnetv1alpha1.CIDRFromString(SubnetIPv6CIDR)
			subnetV6 := &subnetv1alpha1.Subnet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      SubnetNameV6,
					Namespace: OnmetalNamespace,
				},
				Spec: subnetv1alpha1.SubnetSpec{
					CIDR:              cidr,
					NetworkName:       UnderlayNetwork,
					Regions:           []string{TestRegion},
					AvailabilityZones: []string{TestAvailabilityZone},
				},
			}
			Expect(k8sClient.Create(ctx, subnetV6)).To(Succeed())

			list := &switchv1alpha1.SwitchList{}
			Eventually(func() bool {
				Expect(k8sClient.List(ctx, list)).Should(Succeed())
				for _, sw := range list.Items {
					if sw.Spec.SouthSubnetV4 == nil || sw.Spec.SouthSubnetV6 == nil {
						return false
					}
					if !sw.AddressesDefined() {
						return false
					}
					if sw.Status.State != switchv1alpha1.StateFinished {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})

		It("Should update related resources on switch delete", func() {
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

			subnets := &subnetv1alpha1.SubnetList{}
			Eventually(func() bool {
				Expect(k8sClient.List(ctx, subnets)).Should(Succeed())
				for _, subnet := range subnets.Items {
					cidr := *subnet.Spec.CIDR
					if !reflect.DeepEqual(subnet.Status.Vacant, []subnetv1alpha1.CIDR{cidr}) {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())

			list := &switchv1alpha1.SwitchAssignmentList{}
			Eventually(func() bool {
				Expect(k8sClient.List(ctx, list)).Should(Succeed())
				for _, item := range list.Items {
					if item.Status.State != switchv1alpha1.StatePending {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Cleanup environment")
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
	})

	Context("Testing on real data", func() {
		It("Should create switches and define peers", func() {
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

			switchesSamples := []string{
				filepath.Join("..", "config", "samples", "realdata", "sp2.yaml"),
				filepath.Join("..", "config", "samples", "realdata", "lf1.yaml"),
				filepath.Join("..", "config", "samples", "realdata", "lf2.yaml"),
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
				Expect(sw.Labels).Should(Equal(map[string]string{switchv1alpha1.LabelChassisId: switchv1alpha1.MacToLabel(sw.Spec.Chassis.ChassisID)}))
			}

			swa := &switchv1alpha1.SwitchAssignment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      getUUID(chassisIds[0]),
					Namespace: OnmetalNamespace,
				},
				Spec: switchv1alpha1.SwitchAssignmentSpec{
					ChassisID:        chassisIds[0],
					Region:           TestRegion,
					AvailabilityZone: TestAvailabilityZone,
				},
			}
			Expect(k8sClient.Create(ctx, swa)).To(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, swa.NamespacedName(), swa)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(swa.Labels).Should(Equal(map[string]string{switchv1alpha1.LabelChassisId: switchv1alpha1.MacToLabel(chassisIds[0])}))

			list := &switchv1alpha1.SwitchList{}
			Eventually(func() bool {
				Expect(k8sClient.List(ctx, list)).Should(Succeed())
				for _, sw := range list.Items {
					if strings.HasPrefix(sw.Spec.Hostname, "sp") {
						if sw.Status.ConnectionLevel != 0 {
							return false
						}
					}
					if strings.HasPrefix(sw.Spec.Hostname, "lf") {
						if sw.Status.ConnectionLevel != 1 {
							return false
						}
					}
					if sw.Status.State != switchv1alpha1.StateFinished {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})
	})
})

func getUUID(identifier string) string {
	namespaceUUID := uuid.NewMD5(uuid.UUID{}, []byte(OnmetalNamespace))
	newUUID := uuid.NewMD5(namespaceUUID, []byte(identifier))
	return newUUID.String()
}
