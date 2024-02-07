// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
)

// nolint:forcetypeassert
var _ = Describe("Size controller", func() {
	const (
		SizeName      = "test-size"
		SizeNamespace = "default"

		timeout  = time.Second * 30
		interval = time.Millisecond * 250
	)

	AfterEach(func() {
		ctx := context.Background()
		resources := []struct {
			res   client.Object
			list  client.ObjectList
			count func(client.ObjectList) int
		}{
			{
				res:  &metalv1alpha4.Inventory{},
				list: &metalv1alpha4.InventoryList{},
				count: func(objList client.ObjectList) int {
					list := objList.(*metalv1alpha4.InventoryList)
					return len(list.Items)
				},
			},
			{
				res:  &metalv1alpha4.Aggregate{},
				list: &metalv1alpha4.AggregateList{},
				count: func(objList client.ObjectList) int {
					list := objList.(*metalv1alpha4.AggregateList)
					return len(list.Items)
				},
			},
			{
				res:  &metalv1alpha4.Size{},
				list: &metalv1alpha4.SizeList{},
				count: func(objList client.ObjectList) int {
					list := objList.(*metalv1alpha4.SizeList)
					return len(list.Items)
				},
			},
		}

		for _, r := range resources {
			Expect(k8sClient.DeleteAllOf(ctx, r.res, client.InNamespace(SizeNamespace))).To(Succeed())
			Eventually(func() bool {
				err := k8sClient.List(ctx, r.list)
				if err != nil {
					return false
				}
				if r.count(r.list) > 0 {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
		}
	})

	Context("When size CR is changed", func() {
		It("Should be matched or unmatched to inventory CRs", func() {
			By("Inventories are installed")
			ctx := context.Background()

			testSize := metalv1alpha4.Size{
				ObjectMeta: metav1.ObjectMeta{
					Name:      SizeName,
					Namespace: SizeNamespace,
				},
				Spec: metalv1alpha4.SizeSpec{
					Constraints: []metalv1alpha4.ConstraintSpec{
						{
							Path: *metalv1alpha4.JSONPathFromString("spec.cpus[0].cores"),
							Equal: &metalv1alpha4.ConstraintValSpec{
								Numeric: resource.NewScaledQuantity(16, 0),
							},
						},
					},
				},
			}

			sizeLabel := testSize.GetMatchLabel()

			inventoryShouldMatch := inventoryTemplate()
			inventoryShouldMatch.Name = "should-be-matched-and-updated"
			inventoryShouldMatch.Namespace = SizeNamespace
			inventoryShouldMatch.Spec.CPUs[0].Cores = 16
			inventoryShouldMatch.Spec.CPUs[0].Siblings = 32

			inventoryAlreadyMatched := inventoryTemplate()
			inventoryAlreadyMatched.Name = "already-matched"
			inventoryAlreadyMatched.Namespace = SizeNamespace
			inventoryAlreadyMatched.Labels = map[string]string{
				sizeLabel: "true",
			}
			inventoryAlreadyMatched.Spec.CPUs[0].Cores = 16
			inventoryAlreadyMatched.Spec.CPUs[0].Siblings = 16

			inventoryShouldntMatch := inventoryTemplate()
			inventoryShouldntMatch.Name = "should-not-be-matched"
			inventoryShouldntMatch.Namespace = SizeNamespace
			inventoryShouldntMatch.Spec.CPUs[0].Cores = 32
			inventoryShouldntMatch.Spec.CPUs[0].Siblings = 64

			inventoryShouldUnmatch := inventoryTemplate()
			inventoryShouldUnmatch.Name = "should-be-unmatched"
			inventoryShouldUnmatch.Namespace = SizeNamespace
			inventoryShouldUnmatch.Labels = map[string]string{
				sizeLabel: "true",
			}
			inventoryShouldUnmatch.Spec.CPUs[0].Cores = 32
			inventoryShouldUnmatch.Spec.CPUs[0].Siblings = 32

			testInventories := []metalv1alpha4.Inventory{
				*inventoryShouldMatch,
				*inventoryAlreadyMatched,
				*inventoryShouldntMatch,
				*inventoryShouldUnmatch,
			}

			for _, inventory := range testInventories {
				Expect(k8sClient.Create(ctx, &inventory)).Should(Succeed())
			}

			Eventually(func() bool {
				inventoryList := &metalv1alpha4.InventoryList{}
				err := k8sClient.List(ctx, inventoryList)
				if err != nil {
					return false
				}
				if len(inventoryList.Items) == len(testInventories) {
					return true
				}
				return false
			}, timeout, interval).Should(BeTrue())

			By("Size is installed")
			Expect(k8sClient.Create(ctx, &testSize)).Should(Succeed())

			sizeNamespacedName := types.NamespacedName{
				Namespace: SizeNamespace,
				Name:      SizeName,
			}
			createdSize := metalv1alpha4.Size{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, sizeNamespacedName, &createdSize)
				if err != nil {
					return false
				}
				if len(createdSize.GetFinalizers()) > 0 {
					return true
				}
				return false
			}, timeout, interval).Should(BeTrue())
			Expect(createdSize.Spec).To(Equal(testSize.Spec))

			By(fmt.Sprintf("Size should match %s", inventoryShouldMatch.Name))
			Eventually(func() bool {
				inventory := metalv1alpha4.Inventory{}
				inventoryNamespacedName := types.NamespacedName{
					Namespace: inventoryShouldMatch.Namespace,
					Name:      inventoryShouldMatch.Name,
				}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				for k := range inventory.GetLabels() {
					if k == sizeLabel {
						return true
					}
				}
				return false
			}, timeout, interval).Should(BeTrue())

			By(fmt.Sprintf("Size should match but not get updated %s", inventoryAlreadyMatched.Name))
			Consistently(func() bool {
				inventory := metalv1alpha4.Inventory{}
				inventoryNamespacedName := types.NamespacedName{
					Namespace: inventoryAlreadyMatched.Namespace,
					Name:      inventoryAlreadyMatched.Name,
				}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				for k := range inventory.GetLabels() {
					if k == sizeLabel {
						return true
					}
				}
				return false
			}, timeout, interval).Should(BeTrue())

			By(fmt.Sprintf("Size should unmatch %s", inventoryShouldUnmatch.Name))
			Eventually(func() bool {
				inventory := metalv1alpha4.Inventory{}
				inventoryNamespacedName := types.NamespacedName{
					Namespace: inventoryShouldUnmatch.Namespace,
					Name:      inventoryShouldUnmatch.Name,
				}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				for k := range inventory.GetLabels() {
					if k == sizeLabel {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By(fmt.Sprintf("Size shouldn't match %s", inventoryShouldntMatch.Name))
			Consistently(func() bool {
				inventory := metalv1alpha4.Inventory{}
				inventoryNamespacedName := types.NamespacedName{
					Namespace: inventoryShouldntMatch.Namespace,
					Name:      inventoryShouldntMatch.Name,
				}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				for k := range inventory.GetLabels() {
					if k == sizeLabel {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Size is updated")
			createdSize.Spec.Constraints = []metalv1alpha4.ConstraintSpec{
				{
					Path:        *metalv1alpha4.JSONPathFromString("spec.cpus[0].cores"),
					GreaterThan: resource.NewScaledQuantity(30, 0),
				},
			}

			Expect(k8sClient.Update(ctx, &createdSize)).To(Succeed())

			updatedSize := metalv1alpha4.Size{}
			Eventually(func() metalv1alpha4.SizeSpec {
				err := k8sClient.Get(ctx, sizeNamespacedName, &updatedSize)
				if err != nil {
					return metalv1alpha4.SizeSpec{}
				}
				return updatedSize.Spec
			}, timeout, interval).Should(Equal(createdSize.Spec))

			By("Matched inventories should get unmatched")
			matched := []metalv1alpha4.Inventory{
				*inventoryShouldMatch,
				*inventoryAlreadyMatched,
			}

			for _, i := range matched {
				Eventually(func() bool {
					inventory := metalv1alpha4.Inventory{}
					inventoryNamespacedName := types.NamespacedName{
						Namespace: i.Namespace,
						Name:      i.Name,
					}
					err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
					if err != nil {
						return false
					}
					for k := range inventory.GetLabels() {
						if k == sizeLabel {
							return false
						}
					}
					return true
				}, timeout, interval).Should(BeTrue())
			}

			By("Not matched inventories should get matched")
			unmatched := []metalv1alpha4.Inventory{
				*inventoryShouldntMatch,
				*inventoryShouldUnmatch,
			}

			for _, i := range unmatched {
				Eventually(func() bool {
					inventory := metalv1alpha4.Inventory{}
					inventoryNamespacedName := types.NamespacedName{
						Namespace: i.Namespace,
						Name:      i.Name,
					}
					err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
					if err != nil {
						return false
					}
					for k := range inventory.GetLabels() {
						if k == sizeLabel {
							return true
						}
					}
					return false
				}, timeout, interval).Should(BeTrue())
			}

			By("Size is deleted")
			Expect(k8sClient.Delete(ctx, &updatedSize)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, sizeNamespacedName, &updatedSize)
				return apierrors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			By("All size labels are unset")
			Eventually(func() bool {
				inventoryList := &metalv1alpha4.InventoryList{}
				err := k8sClient.List(ctx, inventoryList)
				if err != nil {
					return false
				}
				for _, inventory := range inventoryList.Items {
					for k := range inventory.GetLabels() {
						if k == sizeLabel {
							return false
						}
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})
	})
})

func inventoryTemplate() *metalv1alpha4.Inventory {
	return &metalv1alpha4.Inventory{
		Spec: metalv1alpha4.InventorySpec{
			System: &metalv1alpha4.SystemSpec{
				ID:           "a967954c-3475-11b2-a85c-84d8b4f8cd2d",
				Manufacturer: "LENOVO",
				ProductSKU:   "LENOVO_MT_20JX_BU_Think_FM_ThinkPad T570 W10DG",
				SerialNumber: "R90QR6J0",
			},
			Blocks: []metalv1alpha4.BlockSpec{
				{
					Name:       "JustDisk",
					Type:       "SCSI",
					Rotational: true,
					Model:      "greatModel",
					Size:       1000,
				},
			},
			Memory: &metalv1alpha4.MemorySpec{
				Total: 1024000,
			},
			CPUs: []metalv1alpha4.CPUSpec{
				{
					PhysicalID: 0,
					LogicalIDs: []uint64{0, 1, 2, 3},
					Cores:      2,
					Siblings:   4,
					VendorID:   "GenuineIntel",
					Model:      "78",
					ModelName:  "Intel(R) Core(TM) i5-6300U CPU @ 2.40GHz",
				},
			},
			NICs: []metalv1alpha4.NICSpec{
				{
					Name:       "enp0s31f6",
					PCIAddress: "0000:00:1f.6",
					MACAddress: "48:2a:e3:02:d9:e8",
					MTU:        1400,
					Speed:      1000,
				},
			},
			Host: &metalv1alpha4.HostSpec{
				Name: "dummy.localdomain",
			},
		},
	}
}
