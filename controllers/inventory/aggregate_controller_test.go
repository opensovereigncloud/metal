package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
)

var _ = Describe("Aggregate controller", func() {
	const (
		InventoryName      = "test-inventory"
		AggregateName      = "test-aggregate"
		AggregateNamespace = "default"

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
			Expect(k8sClient.DeleteAllOf(ctx, r.res, client.InNamespace(AggregateNamespace))).To(Succeed())
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

	Context("When aggregate CR is changed", func() {
		It("Should recalculate aggregates for inventory CRs", func() {
			By("Inventories are installed")
			ctx := context.Background()

			testInventory := &inventoryv1alpha1.Inventory{
				ObjectMeta: controllerruntime.ObjectMeta{
					Name:      InventoryName,
					Namespace: AggregateNamespace,
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
					Benchmark: &inventoryv1alpha1.BenchmarkSpec{
						Blocks:   []inventoryv1alpha1.BlockBenchmarkCollection{},
						Networks: []inventoryv1alpha1.NetworkBenchmarkResult{},
					},
				},
			}

			Expect(k8sClient.Create(ctx, testInventory)).Should(Succeed())

			inventoryNamespacedName := types.NamespacedName{
				Namespace: testInventory.Namespace,
				Name:      testInventory.Name,
			}

			Eventually(func() bool {
				createdInventory := inventoryv1alpha1.Inventory{}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &createdInventory)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Aggregate is installed")
			testAggregate := inventoryv1alpha1.Aggregate{
				ObjectMeta: metav1.ObjectMeta{
					Name:      AggregateName,
					Namespace: AggregateNamespace,
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
				err := k8sClient.Get(ctx, createdAggregateNamespacedName, &testAggregate)
				if err != nil {
					return false
				}
				if !controllerutil.ContainsFinalizer(&testAggregate, CAggregateFinalizer) {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Inventory has aggregate calculated")
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
				if iface == nil {
					return false
				}
				maxLogicalId := iface.(string)
				if maxLogicalId != "3" {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Aggregate is updated")
			testAggregate.Spec.Aggregates = []inventoryv1alpha1.AggregateItem{
				{
					SourcePath: *inventoryv1alpha1.JSONPathFromString("spec.cpus[*].cores"),
					TargetPath: *inventoryv1alpha1.JSONPathFromString("cpus.coreCount"),
					Aggregate:  inventoryv1alpha1.CSumAggregateType,
				},
			}

			Expect(k8sClient.Update(ctx, &testAggregate)).Should(Succeed())

			Eventually(func() bool {
				aggregateNamespacedName := types.NamespacedName{
					Namespace: testAggregate.Namespace,
					Name:      testAggregate.Name,
				}
				aggregate := inventoryv1alpha1.Aggregate{}
				err := k8sClient.Get(ctx, aggregateNamespacedName, &aggregate)
				if testAggregate.Spec.Aggregates[0].SourcePath != aggregate.Spec.Aggregates[0].SourcePath {
					return false
				}
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Inventory has aggregate recalculated")
			Eventually(func() bool {
				inventory := inventoryv1alpha1.Inventory{}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				iface, err := nestedMapLookup(inventory.Status.Computed.Object, testAggregate.Name, "cpus", "coreCount")
				if err != nil {
					return false
				}
				if iface == nil {
					return false
				}
				coreCount := iface.(string)
				if coreCount != "2" {
					return false
				}
				_, err = nestedMapLookup(inventory.Status.Computed.Object, testAggregate.Name, "cpus", "maxLogicalId")
				if err == nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Aggregate is deleted")
			Expect(k8sClient.Delete(ctx, &testAggregate)).Should(Succeed())

			By("Inventory has no aggregate")
			Eventually(func() bool {
				inventory := inventoryv1alpha1.Inventory{}
				err := k8sClient.Get(ctx, inventoryNamespacedName, &inventory)
				if err != nil {
					return false
				}
				if _, ok := inventory.Status.Computed.Object[testAggregate.Name]; ok {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})
	})
})
