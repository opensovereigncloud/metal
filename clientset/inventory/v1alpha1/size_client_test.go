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
var _ = Describe("Size client", func() {
	const (
		SizeName         = "test-size"
		SizeToDeleteName = "test-size-to-delete"
		DeleteLabel      = "delete-label"
		SizeNamespace    = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When Size CR is installed", func() {
		It("Should check that Size CR is operational with client", func() {
			By("Creating client")
			finished := make(chan bool)
			ctx := context.Background()

			clientset, err := NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			client := clientset.Sizes(SizeNamespace)

			qty := resource.MustParse("12Gi")
			size := &v1alpha1.Size{
				ObjectMeta: v1.ObjectMeta{
					Name:      SizeName,
					Namespace: SizeNamespace,
				},
				Spec: v1alpha1.SizeSpec{
					Constraints: []v1alpha1.ConstraintSpec{
						{
							Path:        *v1alpha1.JSONPathFromString("a.b.c"),
							GreaterThan: &qty,
						},
					},
				},
			}

			By("Creating watcher")
			watcher, err := client.Watch(ctx, v1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			events := watcher.ResultChan()

			By("Creating Size")
			createdSize := &v1alpha1.Size{}
			go func() {
				defer GinkgoRecover()
				createdSize, err = client.Create(ctx, size, v1.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(createdSize.Spec).Should(Equal(size.Spec))
				finished <- true
			}()

			event := &watch.Event{}
			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Added))
			eventSize := event.Object.(*v1alpha1.Size)
			Expect(eventSize).NotTo(BeNil())
			Expect(eventSize.Spec).Should(Equal(size.Spec))

			<-finished

			By("Updating Size")
			createdSize.Spec.Constraints[0].Path = *v1alpha1.JSONPathFromString("d.e.f")
			go func() {
				defer GinkgoRecover()
				updatedSize, err := client.Update(ctx, createdSize, v1.UpdateOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedSize.Spec).Should(Equal(createdSize.Spec))
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventSize = event.Object.(*v1alpha1.Size)
			Expect(eventSize).NotTo(BeNil())
			Expect(eventSize.Spec).Should(Equal(createdSize.Spec))

			<-finished

			By("Patching Size")
			patch := []struct {
				Op    string `json:"op"`
				Path  string `json:"path"`
				Value string `json:"value"`
			}{{
				Op:    "replace",
				Path:  "/spec/constraints/0/path",
				Value: "g.h.i",
			}}

			patchData, err := json.Marshal(patch)
			Expect(err).NotTo(HaveOccurred())

			go func() {
				defer GinkgoRecover()
				patchedSize, err := client.Patch(ctx, SizeName, types.JSONPatchType, patchData, v1.PatchOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(patchedSize.Spec.Constraints[0].Path.String()).Should(Equal(v1alpha1.JSONPathFromString(patch[0].Value).String()))
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventSize = event.Object.(*v1alpha1.Size)
			Expect(eventSize).NotTo(BeNil())
			Expect(eventSize.Spec.Constraints[0].Path.String()).Should(Equal(v1alpha1.JSONPathFromString(patch[0].Value).String()))

			<-finished

			// We do not handle status for Size atm,
			// so just a placeholder for now
			By("Updating Size status")
			_, err = client.UpdateStatus(ctx, eventSize, v1.UpdateOptions{})
			Expect(err).NotTo(HaveOccurred())
			Eventually(events).Should(Receive())

			sizeToDelete := &v1alpha1.Size{
				ObjectMeta: v1.ObjectMeta{
					Name:      SizeToDeleteName,
					Namespace: SizeNamespace,
					Labels: map[string]string{
						DeleteLabel: "",
					},
				},
				Spec: v1alpha1.SizeSpec{},
			}

			By("Creating Size collection")
			_, err = client.Create(ctx, sizeToDelete, v1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
			Eventually(events).Should(Receive())

			By("Listing Sizes")
			sizeList, err := client.List(ctx, v1.ListOptions{})
			Expect(sizeList).NotTo(BeNil())
			Expect(sizeList.Items).To(HaveLen(2))

			By("Bulk deleting Size")
			Expect(client.DeleteCollection(ctx, v1.DeleteOptions{}, v1.ListOptions{LabelSelector: DeleteLabel})).To(Succeed())

			By("Requesting created Size")
			Eventually(func() bool {
				_, err = client.Get(ctx, SizeName, v1.GetOptions{})
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Eventually(func() bool {
				_, err = client.Get(ctx, SizeToDeleteName, v1.GetOptions{})
				return err == nil
			}, timeout, interval).Should(BeFalse())

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Deleted))
			eventSize = event.Object.(*v1alpha1.Size)
			Expect(eventSize).NotTo(BeNil())
			Expect(eventSize.Name).To(Equal(SizeToDeleteName))

			By("Deleting Size")
			go func() {
				defer GinkgoRecover()
				err := client.Delete(ctx, SizeName, v1.DeleteOptions{})
				Expect(err).NotTo(HaveOccurred())
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Deleted))
			eventSize = event.Object.(*v1alpha1.Size)
			Expect(eventSize).NotTo(BeNil())
			Expect(eventSize.Name).To(Equal(SizeName))

			<-finished

			watcher.Stop()
			Eventually(events).Should(BeClosed())
		})
	})
})
