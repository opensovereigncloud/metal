package controllers

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	inventoryv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
)

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
				res:  &inventoryv1alpha1.Aggregate{},
				list: &inventoryv1alpha1.AggregateList{},
				count: func(objList client.ObjectList) int {
					list := objList.(*inventoryv1alpha1.AggregateList)
					return len(list.Items)
				},
			},
			{
				res:  &inventoryv1alpha1.Size{},
				list: &inventoryv1alpha1.SizeList{},
				count: func(objList client.ObjectList) int {
					list := objList.(*inventoryv1alpha1.SizeList)
					return len(list.Items)
				},
			},
			{
				res:  &inventoryv1alpha1.Inventory{},
				list: &inventoryv1alpha1.InventoryList{},
				count: func(objList client.ObjectList) int {
					list := objList.(*inventoryv1alpha1.InventoryList)
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

			testSize := inventoryv1alpha1.Size{
				ObjectMeta: metav1.ObjectMeta{
					Name:      SizeName,
					Namespace: SizeNamespace,
				},
				Spec: inventoryv1alpha1.SizeSpec{
					Constraints: []inventoryv1alpha1.ConstraintSpec{
						{
							Path: "cpus[0].cores",
							Equal: &inventoryv1alpha1.ConstraintValSpec{
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

			testInventories := []inventoryv1alpha1.Inventory{
				*inventoryShouldMatch,
				*inventoryAlreadyMatched,
				*inventoryShouldntMatch,
				*inventoryShouldUnmatch,
			}

			for _, inventory := range testInventories {
				Expect(k8sClient.Create(ctx, &inventory)).Should(Succeed())
			}

			Eventually(func() bool {
				inventoryList := &inventoryv1alpha1.InventoryList{}
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
			createdSize := inventoryv1alpha1.Size{}
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
				inventory := inventoryv1alpha1.Inventory{}
				inventoryNamespacedName := types.NamespacedName{
					Namespace: inventoryShouldMatch.Namespace,
					Name:      inventoryShouldMatch.Name,
				}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				for k, _ := range inventory.GetLabels() {
					if k == sizeLabel {
						return true
					}
				}
				return false
			}, timeout, interval).Should(BeTrue())

			By(fmt.Sprintf("Size should match but not get updated %s", inventoryAlreadyMatched.Name))
			Consistently(func() bool {
				inventory := inventoryv1alpha1.Inventory{}
				inventoryNamespacedName := types.NamespacedName{
					Namespace: inventoryAlreadyMatched.Namespace,
					Name:      inventoryAlreadyMatched.Name,
				}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				for k, _ := range inventory.GetLabels() {
					if k == sizeLabel {
						return true
					}
				}
				return false
			}, timeout, interval).Should(BeTrue())

			By(fmt.Sprintf("Size should unmatch %s", inventoryShouldUnmatch.Name))
			Eventually(func() bool {
				inventory := inventoryv1alpha1.Inventory{}
				inventoryNamespacedName := types.NamespacedName{
					Namespace: inventoryShouldUnmatch.Namespace,
					Name:      inventoryShouldUnmatch.Name,
				}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				for k, _ := range inventory.GetLabels() {
					if k == sizeLabel {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By(fmt.Sprintf("Size shouldn't match %s", inventoryShouldntMatch.Name))
			Consistently(func() bool {
				inventory := inventoryv1alpha1.Inventory{}
				inventoryNamespacedName := types.NamespacedName{
					Namespace: inventoryShouldntMatch.Namespace,
					Name:      inventoryShouldntMatch.Name,
				}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				for k, _ := range inventory.GetLabels() {
					if k == sizeLabel {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Size is updated")
			createdSize.Spec.Constraints = []inventoryv1alpha1.ConstraintSpec{
				{
					Path:        "cpus[0].cores",
					GreaterThan: resource.NewScaledQuantity(30, 0),
				},
			}

			Expect(k8sClient.Update(ctx, &createdSize)).To(Succeed())

			updatedSize := inventoryv1alpha1.Size{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, sizeNamespacedName, &updatedSize)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(updatedSize.Spec).To(Equal(createdSize.Spec))

			By("Matched inventories should get unmatched")
			matched := []inventoryv1alpha1.Inventory{
				*inventoryShouldMatch,
				*inventoryAlreadyMatched,
			}

			for _, i := range matched {
				Eventually(func() bool {
					inventory := inventoryv1alpha1.Inventory{}
					inventoryNamespacedName := types.NamespacedName{
						Namespace: i.Namespace,
						Name:      i.Name,
					}
					err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
					if err != nil {
						return false
					}
					for k, _ := range inventory.GetLabels() {
						if k == sizeLabel {
							return false
						}
					}
					return true
				}, timeout, interval).Should(BeTrue())
			}

			By("Not matched inventories should get matched")
			unmatched := []inventoryv1alpha1.Inventory{
				*inventoryShouldntMatch,
				*inventoryShouldUnmatch,
			}

			for _, i := range unmatched {
				Eventually(func() bool {
					inventory := inventoryv1alpha1.Inventory{}
					inventoryNamespacedName := types.NamespacedName{
						Namespace: i.Namespace,
						Name:      i.Name,
					}
					err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
					if err != nil {
						return false
					}
					for k, _ := range inventory.GetLabels() {
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
				if apierrors.IsNotFound(err) {
					return true
				}
				return false
			}, timeout, interval).Should(BeTrue())

			By("All size labels are unset")
			Eventually(func() bool {
				inventoryList := &inventoryv1alpha1.InventoryList{}
				err := k8sClient.List(ctx, inventoryList)
				if err != nil {
					return false
				}
				for _, inventory := range inventoryList.Items {
					for k, _ := range inventory.GetLabels() {
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

func inventoryTemplate() *inventoryv1alpha1.Inventory {
	return &inventoryv1alpha1.Inventory{
		Spec: inventoryv1alpha1.InventorySpec{
			System: &inventoryv1alpha1.SystemSpec{
				ID:           "a967954c-3475-11b2-a85c-84d8b4f8cd2d",
				Manufacturer: "LENOVO",
				ProductSKU:   "LENOVO_MT_20JX_BU_Think_FM_ThinkPad T570 W10DG",
				SerialNumber: "R90QR6J0",
			},
			Blocks: []inventoryv1alpha1.BlockSpec{
				{
					Name:       "JustDisk",
					Type:       "SCSI",
					Rotational: true,
					Model:      "greatModel",
					Size:       1000,
				},
			},
			Memory: &inventoryv1alpha1.MemorySpec{
				Total: 1024000,
			},
			CPUs: []inventoryv1alpha1.CPUSpec{
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
			NICs: []inventoryv1alpha1.NICSpec{
				{
					Name:       "enp0s31f6",
					PCIAddress: "0000:00:1f.6",
					MACAddress: "48:2a:e3:02:d9:e8",
					MTU:        1400,
					Speed:      1000,
				},
			},
			Host: &inventoryv1alpha1.HostSpec{
				Type: "Machine",
				Name: "dummy.localdomain",
			},
			Benchmark: &inventoryv1alpha1.BenchmarkSpec{
				Blocks:  []inventoryv1alpha1.BlockBenchmarkResult{},
				Network: &inventoryv1alpha1.NetworkBenchmarkResult{},
			},
		},
	}
}
