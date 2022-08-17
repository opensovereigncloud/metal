package v1alpha1

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/onmetal/metal-api/apis/inventory/v1alpha1"
)

// nolint:forcetypeassert
var _ = Describe("Inventory client", func() {
	const (
		InventoryName         = "test-inventory"
		InventoryToDeleteName = "test-inventory-to-delete"
		DeleteLabel           = "delete-label"
		InventoryNamespace    = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When Inventory CR is installed", func() {
		It("Should check that Inventory CR is operational with client", func() {
			By("Creating client")
			finished := make(chan bool)
			ctx := context.Background()

			clientset, err := NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			client := clientset.Inventories(InventoryNamespace)

			inventory := &v1alpha1.Inventory{
				ObjectMeta: v1.ObjectMeta{
					Name:      InventoryName,
					Namespace: InventoryNamespace,
				},
				Spec: v1alpha1.InventorySpec{
					System: &v1alpha1.SystemSpec{
						ID:           "a967954c-3475-11b2-a85c-84d8b4f8cd2d",
						Manufacturer: "LENOVO",
						ProductSKU:   "LENOVO_MT_20JX_BU_Think_FM_ThinkPad T570 W10DG",
						SerialNumber: "R90QR6J0",
					},
					Blocks: []v1alpha1.BlockSpec{
						{
							Name:       "JustDisk",
							Type:       "SCSI",
							Rotational: true,
							Model:      "greatModel",
							Size:       1000,
						},
					},
					Memory: &v1alpha1.MemorySpec{
						Total: 1024000,
					},
					CPUs: []v1alpha1.CPUSpec{
						{
							PhysicalID: 0,
							LogicalIDs: []uint64{0, 1, 2, 3},
							Cores:      2,
							Siblings:   4,
							VendorID:   "GenuineIntel",
							Model:      "78",
							ModelName:  "Intel(R) Core(TM) i5-6300U CPU @ 2.40GHz",
							// These two values should be set through parse method
							// in order to avoid failure on deep equality.
							// On deserialization, Quantity sets internal string representation
							// to the string value coming from JSON.
							// If left empty or constructed from int64, string representation
							// will not be set, and DeepEqual method will return false in result.
							MHz:      resource.MustParse("0"),
							BogoMIPS: resource.MustParse("0"),
						},
					},
					NICs: []v1alpha1.NICSpec{
						{
							Name:       "enp0s31f6",
							PCIAddress: "0000:00:1f.6",
							MACAddress: "48:2a:e3:02:d9:e8",
							MTU:        1400,
							Speed:      1000,
						},
					},
					Host: &v1alpha1.HostSpec{
						Name: "dummy.localdomain",
					},
				},
			}

			By("Creating watcher")
			watcher, err := client.Watch(ctx, v1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			events := watcher.ResultChan()

			By("Creating Inventory")
			createdInventory := &v1alpha1.Inventory{}
			go func() {
				defer GinkgoRecover()
				createdInventory, err = client.Create(ctx, inventory, v1.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(createdInventory.Spec).Should(Equal(inventory.Spec))
				finished <- true
			}()

			event := &watch.Event{}
			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Added))
			eventInventory := event.Object.(*v1alpha1.Inventory)
			Expect(eventInventory).NotTo(BeNil())
			Expect(eventInventory.Spec).Should(Equal(inventory.Spec))

			<-finished

			By("Updating Inventory")
			createdInventory.Spec.Host.Name = "updateddummy.localdomain"
			go func() {
				defer GinkgoRecover()
				updatedInventory, err := client.Update(ctx, createdInventory, v1.UpdateOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedInventory.Spec).Should(Equal(createdInventory.Spec))
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventInventory = event.Object.(*v1alpha1.Inventory)
			Expect(eventInventory).NotTo(BeNil())
			Expect(eventInventory.Spec).Should(Equal(createdInventory.Spec))

			<-finished

			By("Patching Inventory")
			patch := []struct {
				Op    string `json:"op"`
				Path  string `json:"path"`
				Value string `json:"value"`
			}{{
				Op:    "replace",
				Path:  "/spec/host/name",
				Value: "patcheddummy.localdomain",
			}}

			patchData, err := json.Marshal(patch)
			Expect(err).NotTo(HaveOccurred())

			go func() {
				defer GinkgoRecover()
				patchedInventory, err := client.Patch(ctx, InventoryName, types.JSONPatchType, patchData, v1.PatchOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(patchedInventory.Spec.Host.Name).Should(Equal(patch[0].Value))
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventInventory = event.Object.(*v1alpha1.Inventory)
			Expect(eventInventory).NotTo(BeNil())
			Expect(eventInventory.Spec.Host.Name).Should(Equal(patch[0].Value))

			<-finished

			// We do not handle status for Inventory atm,
			// so just a placeholder for now
			By("Updating Inventory status")
			_, err = client.UpdateStatus(ctx, eventInventory, v1.UpdateOptions{})
			Expect(err).NotTo(HaveOccurred())
			Eventually(events).Should(Receive())

			inventoryToDelete := inventory.DeepCopy()
			inventoryToDelete.Name = InventoryToDeleteName
			inventoryToDelete.Labels = map[string]string{
				DeleteLabel: "",
			}

			By("Creating Inventory collection")
			_, err = client.Create(ctx, inventoryToDelete, v1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
			Eventually(events).Should(Receive())

			By("Listing Inventories")
			inventoryList, err := client.List(ctx, v1.ListOptions{})
			Expect(inventoryList).NotTo(BeNil())
			Expect(inventoryList.Items).To(HaveLen(2))

			By("Bulk deleting Inventory")
			Expect(client.DeleteCollection(ctx, v1.DeleteOptions{}, v1.ListOptions{LabelSelector: DeleteLabel})).To(Succeed())

			By("Requesting created Inventory")
			Eventually(func() bool {
				_, err = client.Get(ctx, InventoryName, v1.GetOptions{})
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Eventually(func() bool {
				_, err = client.Get(ctx, InventoryToDeleteName, v1.GetOptions{})
				return err == nil
			}, timeout, interval).Should(BeFalse())

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Deleted))
			eventInventory = event.Object.(*v1alpha1.Inventory)
			Expect(eventInventory).NotTo(BeNil())
			Expect(eventInventory.Name).To(Equal(InventoryToDeleteName))

			By("Deleting Inventory")
			go func() {
				defer GinkgoRecover()
				err := client.Delete(ctx, InventoryName, v1.DeleteOptions{})
				Expect(err).NotTo(HaveOccurred())
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Deleted))
			eventInventory = event.Object.(*v1alpha1.Inventory)
			Expect(eventInventory).NotTo(BeNil())
			Expect(eventInventory.Name).To(Equal(InventoryName))

			<-finished

			watcher.Stop()
			Eventually(events).Should(BeClosed())
		})
	})
})
