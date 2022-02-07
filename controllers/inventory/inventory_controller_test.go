package controllers

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
)

var _ = Describe("Inventory controller", func() {
	const (
		InventoryName      = "test-inventory"
		InventoryNamespace = "default"

		timeout  = time.Second * 30
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
				res:  &inventoryv1alpha1.Inventory{},
				list: &inventoryv1alpha1.InventoryList{},
				count: func(objList client.ObjectList) int {
					list := objList.(*inventoryv1alpha1.InventoryList)
					return len(list.Items)
				},
			},
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

	nestedMapLookup := func(theMap map[string]interface{}, keys ...string) (interface{}, error) {
		nextMap := theMap
		lastKeyIdx := len(keys) - 1
		if nextMap == nil {
			return nil, errors.New("map is nil")
		}
		if lastKeyIdx < 0 {
			return nil, errors.New("keys are not provided")
		}
		var result interface{}
		var err error
		for i, key := range keys {
			mapIface, ok := nextMap[key]
			if !ok {
				result, err = nil, errors.Errorf("key %d, %s not found", i, key)
				break
			}
			if lastKeyIdx == i {
				result, err = mapIface, nil
				break
			}
			if mapIface == nil {
				result, err = nil, errors.Errorf("key %d, %s returns nil instead of map", i, key)
				break
			}
			nextMap, ok = mapIface.(map[string]interface{})
			if !ok {
				result, err = nil, errors.Errorf("cant cast value for key %d, %s to the map", i, key)
				break
			}
		}
		return result, err
	}

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
							Path: *inventoryv1alpha1.JSONPathFromString("spec.cpus[0].cores"),
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
							Path: *inventoryv1alpha1.JSONPathFromString("spec.cpus[0].siblings"),
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
							Path: *inventoryv1alpha1.JSONPathFromString("spec.cpus[0].cores"),
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
							Path: *inventoryv1alpha1.JSONPathFromString("spec.cpus[0].siblings"),
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
							SourcePath: *inventoryv1alpha1.JSONPathFromString("spec.cpus[*].logicalIds[*]"),
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
				for k := range inventory.GetLabels() {
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
				for k := range inventory.GetLabels() {
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
				for k := range inventory.GetLabels() {
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
				for k := range inventory.GetLabels() {
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
				iface, err := nestedMapLookup(inventory.Status.Computed.Object, testAggregate.Name, "cpus", "maxLogicalId")
				if err != nil {
					return false
				}
				maxLogicalId := iface.(string)
				return maxLogicalId == "3"
			}, timeout, interval).Should(BeTrue())

			By("Inventory is updated")
			Expect(k8sClient.Get(ctx, inventoryNamespacedName, &createdInventory)).To(Succeed())

			createdInventory.Spec.CPUs[0].LogicalIDs = append(createdInventory.Spec.CPUs[0].LogicalIDs, 5)
			createdInventory.Spec.CPUs[0].Cores = 8
			createdInventory.Spec.CPUs[0].Siblings = 16

			Expect(k8sClient.Update(ctx, &createdInventory)).To(Succeed())

			updatedInventory := inventoryv1alpha1.Inventory{}
			Eventually(func() inventoryv1alpha1.InventorySpec {
				err := k8sClient.Get(ctx, inventoryNamespacedName, &updatedInventory)
				if err != nil {
					return inventoryv1alpha1.InventorySpec{}
				}
				return updatedInventory.Spec
			}, timeout, interval).Should(Equal(createdInventory.Spec))

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
					for k := range inventory.GetLabels() {
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
					for k := range inventory.GetLabels() {
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
				iface, err := nestedMapLookup(inventory.Status.Computed.Object, testAggregate.Name, "cpus", "maxLogicalId")
				if err != nil {
					return false
				}
				maxLogicalId := iface.(string)
				return maxLogicalId == "5"
			}, timeout, interval).Should(BeTrue())

			By("Inventory is deleted")
			Expect(k8sClient.Get(ctx, inventoryNamespacedName, &updatedInventory)).To(Succeed())

			Expect(k8sClient.Delete(ctx, &updatedInventory)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, inventoryNamespacedName, &updatedInventory)
				return apierrors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
		})
	})
})
