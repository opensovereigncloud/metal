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

var _ = Describe("Inventory controller", func() {
	const (
		InventoryName      = "test-inventory"
		InventoryNamespace = "default"

		timeout  = time.Second * 20
		interval = time.Millisecond * 500
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
			Expect(k8sClient.DeleteAllOf(ctx, r.res, client.InNamespace(InventoryNamespace))).To(Succeed())
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

	Context("When inventory CR is created and updated", func() {
		It("Should be matched or unmatched to size CRs", func() {
			By("Sizes are installed")
			ctx := context.Background()

			sizeShouldMatch := inventoryv1alpha1.Size{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "should-match",
					Namespace: InventoryNamespace,
				},
				Spec: inventoryv1alpha1.SizeSpec{
					Constraints: []inventoryv1alpha1.ConstraintSpec{
						{
							Path: "cpus.cores",
							Equal: &inventoryv1alpha1.ConstraintValSpec{
								Numeric: resource.NewScaledQuantity(2, 0),
							},
						},
					},
				},
			}

			sizeAlreadyMatched := inventoryv1alpha1.Size{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "already-matched",
					Namespace: InventoryNamespace,
				},
				Spec: inventoryv1alpha1.SizeSpec{
					Constraints: []inventoryv1alpha1.ConstraintSpec{
						{
							Path: "cpus.threads",
							Equal: &inventoryv1alpha1.ConstraintValSpec{
								Numeric: resource.NewScaledQuantity(4, 0),
							},
						},
					},
				},
			}

			sizeShouldNotMatch := inventoryv1alpha1.Size{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "should-not-match",
					Namespace: InventoryNamespace,
				},
				Spec: inventoryv1alpha1.SizeSpec{
					Constraints: []inventoryv1alpha1.ConstraintSpec{
						{
							Path: "cpus.cores",
							Equal: &inventoryv1alpha1.ConstraintValSpec{
								Numeric: resource.NewScaledQuantity(8, 0),
							},
						},
					},
				},
			}

			sizeShouldUnmatch := inventoryv1alpha1.Size{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "should-unmatch",
					Namespace: InventoryNamespace,
				},
				Spec: inventoryv1alpha1.SizeSpec{
					Constraints: []inventoryv1alpha1.ConstraintSpec{
						{
							Path: "cpus.threads",
							Equal: &inventoryv1alpha1.ConstraintValSpec{
								Numeric: resource.NewScaledQuantity(16, 0),
							},
						},
					},
				},
			}

			testSizes := []inventoryv1alpha1.Size{
				sizeShouldMatch,
				sizeAlreadyMatched,
				sizeShouldNotMatch,
				sizeShouldUnmatch,
			}

			for _, size := range testSizes {
				Expect(k8sClient.Create(ctx, &size)).Should(Succeed())
			}

			Eventually(func() bool {
				sizeList := &inventoryv1alpha1.SizeList{}
				err := k8sClient.List(ctx, sizeList)
				if err != nil {
					return false
				}
				if len(sizeList.Items) == len(testSizes) {
					return true
				}
				return false
			}, timeout, interval).Should(BeTrue())

			By("Aggregate is installed")
			testAggregate := inventoryv1alpha1.Aggregate{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-aggregate",
					Namespace: InventoryNamespace,
				},
				Spec: inventoryv1alpha1.AggregateSpec{
					Aggregates: []inventoryv1alpha1.AggregateItem{
						{
							SourcePath: *inventoryv1alpha1.JSONPathFromString("cpus.cpus[*].logicalIds[*]"),
							TargetPath: *inventoryv1alpha1.JSONPathFromString("cpus.maxLogicalId"),
							Aggregate:  inventoryv1alpha1.CMaxAggregateType,
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, &testAggregate)).Should(Succeed())

			Eventually(func() bool {
				createdAggregateNamespacedName := types.NamespacedName{
					Namespace: testAggregate.Namespace,
					Name:      testAggregate.Name,
				}
				createdAggregate := inventoryv1alpha1.Aggregate{}
				err := k8sClient.Get(ctx, createdAggregateNamespacedName, &createdAggregate)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Inventory is installed")
			testInventory := inventoryv1alpha1.Inventory{
				ObjectMeta: metav1.ObjectMeta{
					Name:      InventoryName,
					Namespace: InventoryNamespace,
					Labels: map[string]string{
						sizeAlreadyMatched.GetMatchLabel(): "true",
						sizeShouldUnmatch.GetMatchLabel():  "true",
					},
				},
				Spec: inventoryv1alpha1.InventorySpec{
					System: &inventoryv1alpha1.SystemSpec{
						ID:           "a967954c-3475-11b2-a85c-84d8b4f8cd2d",
						Manufacturer: "LENOVO",
						ProductSKU:   "LENOVO_MT_20JX_BU_Think_FM_ThinkPad T570 W10DG",
						SerialNumber: "R90QR6J0",
					},
					Blocks: &inventoryv1alpha1.BlockTotalSpec{
						Count:    1,
						Capacity: 1,
						Blocks: []inventoryv1alpha1.BlockSpec{
							{
								Name:       "JustDisk",
								Type:       "SCSI",
								Rotational: true,
								Model:      "greatModel",
								Size:       1000,
							},
						},
					},
					Memory: &inventoryv1alpha1.MemorySpec{
						Total: 1024000,
					},
					CPUs: &inventoryv1alpha1.CPUTotalSpec{
						Sockets: 1,
						Cores:   2,
						Threads: 4,
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
					},
					NICs: &inventoryv1alpha1.NICTotalSpec{
						Count: 1,
						NICs: []inventoryv1alpha1.NICSpec{
							{
								Name:       "enp0s31f6",
								PCIAddress: "0000:00:1f.6",
								MACAddress: "48:2a:e3:02:d9:e8",
								MTU:        1400,
								Speed:      1000,
							},
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

			Expect(k8sClient.Create(ctx, &testInventory)).Should(Succeed())
			inventoryNamespacedName := types.NamespacedName{
				Namespace: InventoryNamespace,
				Name:      InventoryName,
			}
			createdInventory := inventoryv1alpha1.Inventory{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, inventoryNamespacedName, &createdInventory)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(createdInventory.Spec).To(Equal(testInventory.Spec))

			By(fmt.Sprintf("Inventory should match %s", sizeShouldMatch.Name))
			Eventually(func() bool {
				inventory := inventoryv1alpha1.Inventory{}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				for k, _ := range inventory.GetLabels() {
					if k == sizeShouldMatch.GetMatchLabel() {
						return true
					}
				}
				return false
			}, timeout, interval).Should(BeTrue())

			By(fmt.Sprintf("Inventory should match but not get updated %s", sizeAlreadyMatched.Name))
			Consistently(func() bool {
				inventory := inventoryv1alpha1.Inventory{}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				for k, _ := range inventory.GetLabels() {
					if k == sizeAlreadyMatched.GetMatchLabel() {
						return true
					}
				}
				return false
			}, timeout, interval).Should(BeTrue())

			By(fmt.Sprintf("Inventory should unmatch %s", sizeShouldUnmatch.Name))
			Eventually(func() bool {
				inventory := inventoryv1alpha1.Inventory{}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				for k, _ := range inventory.GetLabels() {
					if k == sizeShouldUnmatch.GetMatchLabel() {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By(fmt.Sprintf("Inventory shouldn't match %s", sizeShouldNotMatch.Name))
			Consistently(func() bool {
				inventory := inventoryv1alpha1.Inventory{}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				for k, _ := range inventory.GetLabels() {
					if k == sizeShouldNotMatch.GetMatchLabel() {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Inventory has aggregate computed")
			Eventually(func() bool {
				inventory := inventoryv1alpha1.Inventory{}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				if inventory.Status.Computed.Object == nil {
					return false
				}
				aggregateMap, ok := inventory.Status.Computed.Object[testAggregate.Name]
				if !ok {
					return false
				}
				maxLogicalId := aggregateMap.(map[string]interface{})["cpus"].(map[string]interface{})["maxLogicalId"].(string)
				if maxLogicalId != "3" {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Inventory is updated")
			Expect(k8sClient.Get(ctx, inventoryNamespacedName, &createdInventory)).To(Succeed())

			createdInventory.Spec.CPUs.Cores = 8
			createdInventory.Spec.CPUs.Threads = 16
			createdInventory.Spec.CPUs.CPUs[0].LogicalIDs = append(createdInventory.Spec.CPUs.CPUs[0].LogicalIDs, 5)
			createdInventory.Spec.Benchmark = &inventoryv1alpha1.BenchmarkSpec{
				Blocks: []inventoryv1alpha1.BlockBenchmarkResult{
					{
						BlockName:           "/dev/hda",
						SmallBlockReadIOPS:  10000,
						SmallBlockWriteIOPS: 9000,
						BandwidthReadIOPS:   8000,
						BandwidthWriteIOPS:  7000,
					},
				},
				Network: &inventoryv1alpha1.NetworkBenchmarkResult{AverageNetworkThroughputBPS: 15000000000},
			}

			Expect(k8sClient.Update(ctx, &createdInventory)).To(Succeed())

			updatedInventory := inventoryv1alpha1.Inventory{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, inventoryNamespacedName, &updatedInventory)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(updatedInventory.Spec).To(Equal(createdInventory.Spec))

			By("Matched size should get unmatched")
			matchedLabels := []string{
				sizeShouldMatch.GetMatchLabel(),
				sizeAlreadyMatched.GetMatchLabel(),
			}

			for _, label := range matchedLabels {
				Eventually(func() bool {
					inventory := inventoryv1alpha1.Inventory{}
					err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
					if err != nil {
						return false
					}
					for k, _ := range inventory.GetLabels() {
						if k == label {
							return false
						}
					}
					return true
				}, timeout, interval).Should(BeTrue())
			}

			By("Not matched sizes should get matched")
			unmatchedLabels := []string{
				sizeShouldNotMatch.GetMatchLabel(),
				sizeShouldUnmatch.GetMatchLabel(),
			}

			for _, label := range unmatchedLabels {
				Eventually(func() bool {
					inventory := inventoryv1alpha1.Inventory{}
					err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
					if err != nil {
						return false
					}
					for k, _ := range inventory.GetLabels() {
						if k == label {
							return true
						}
					}
					return false
				}, timeout, interval).Should(BeTrue())
			}

			By("Inventory has aggregate recalculated")
			Eventually(func() bool {
				inventory := inventoryv1alpha1.Inventory{}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				if inventory.Status.Computed.Object == nil {
					return false
				}
				aggregateMap, ok := inventory.Status.Computed.Object[testAggregate.Name]
				if !ok {
					return false
				}
				maxLogicalId := aggregateMap.(map[string]interface{})["cpus"].(map[string]interface{})["maxLogicalId"].(string)
				if maxLogicalId != "5" {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Inventory is deleted")
			Expect(k8sClient.Get(ctx, inventoryNamespacedName, &updatedInventory)).To(Succeed())

			Expect(k8sClient.Delete(ctx, &updatedInventory)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, inventoryNamespacedName, &updatedInventory)
				if apierrors.IsNotFound(err) {
					return true
				}
				return false
			}, timeout, interval).Should(BeTrue())
		})
	})

	Context("When inventory CR is created and updated", func() {
		It("Labels should be updated accordingly", func() {
			ctx := context.Background()
			By("Inventory is installed")
			testInventory := inventoryv1alpha1.Inventory{
				ObjectMeta: metav1.ObjectMeta{
					Name:      InventoryName,
					Namespace: InventoryNamespace,
				},
				Spec: inventoryv1alpha1.InventorySpec{
					System: &inventoryv1alpha1.SystemSpec{
						ID:           "a967954c-3475-11b2-a85c-84d8b4f8cd2d",
						Manufacturer: "LENOVO",
						ProductSKU:   "LENOVO_MT_20JX_BU_Think_FM_ThinkPad T570 W10DG",
						SerialNumber: "R90QR6J0",
					},
					Blocks: &inventoryv1alpha1.BlockTotalSpec{
						Count:    1,
						Capacity: 1,
						Blocks: []inventoryv1alpha1.BlockSpec{
							{
								Name:       "JustDisk",
								Type:       "SCSI",
								Rotational: true,
								Model:      "greatModel",
								Size:       1000,
							},
						},
					},
					Memory: &inventoryv1alpha1.MemorySpec{
						Total: 1024000,
					},
					CPUs: &inventoryv1alpha1.CPUTotalSpec{
						Sockets: 1,
						Cores:   2,
						Threads: 4,
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
					},
					NICs: &inventoryv1alpha1.NICTotalSpec{
						Count: 1,
						NICs: []inventoryv1alpha1.NICSpec{
							{
								Name:       "enp0s31f6",
								PCIAddress: "0000:00:1f.6",
								MACAddress: "48:2a:e3:02:d9:e8",
								MTU:        1400,
								Speed:      1000,
							},
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

			Expect(k8sClient.Create(ctx, &testInventory)).Should(Succeed())
			inventoryNamespacedName := types.NamespacedName{
				Namespace: InventoryNamespace,
				Name:      InventoryName,
			}
			createdInventory := inventoryv1alpha1.Inventory{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, inventoryNamespacedName, &createdInventory)
				if err != nil {
					return false
				}
				return createdInventory.Labels[inventoryv1alpha1.CInventoryTypeLabel] == "Machine" &&
					createdInventory.Labels[inventoryv1alpha1.CInventoryHostnameLabel] == "dummy.localdomain"
			}, timeout, interval).Should(BeTrue())

			By("When hostname and type changed")
			createdInventory.Spec.Host.Type = "Switch"
			createdInventory.Spec.Host.Name = "coolname.localdomain"
			Expect(k8sClient.Update(ctx, &createdInventory)).Should(Succeed())

			By("Then hostname and type labels should also change eventually")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, inventoryNamespacedName, &createdInventory)
				if err != nil {
					return false
				}
				return createdInventory.Labels[inventoryv1alpha1.CInventoryTypeLabel] == "Switch" &&
					createdInventory.Labels[inventoryv1alpha1.CInventoryHostnameLabel] == "coolname.localdomain"
			}, timeout, interval).Should(BeTrue())
		})
	})
})
