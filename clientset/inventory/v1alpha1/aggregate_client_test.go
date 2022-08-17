package v1alpha1

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/onmetal/metal-api/apis/inventory/v1alpha1"
)

// nolint:forcetypeassert
var _ = Describe("Aggregate client", func() {
	const (
		AggregateName         = "test-aggregate"
		AggregateToDeleteName = "test-aggregate-to-delete"
		DeleteLabel           = "delete-label"
		AggregateNamespace    = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When Aggregate CR is installed", func() {
		It("Should check that Aggregate CR is operational with client", func() {
			By("Creating client")
			finished := make(chan bool)
			ctx := context.Background()

			clientset, err := NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			client := clientset.Aggregates(AggregateNamespace)

			aggregate := &v1alpha1.Aggregate{
				ObjectMeta: v1.ObjectMeta{
					Name:      AggregateName,
					Namespace: AggregateNamespace,
				},
				Spec: v1alpha1.AggregateSpec{
					Aggregates: []v1alpha1.AggregateItem{
						{
							SourcePath: *v1alpha1.JSONPathFromString("spec.cpus"),
							TargetPath: *v1alpha1.JSONPathFromString("status.computed.cpus.cpuCount"),
							Aggregate:  v1alpha1.CCountAggregateType,
						},
					},
				},
			}

			By("Creating watcher")
			watcher, err := client.Watch(ctx, v1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			events := watcher.ResultChan()

			By("Creating Aggregate")
			createdAggregate := &v1alpha1.Aggregate{}
			go func() {
				defer GinkgoRecover()
				createdAggregate, err = client.Create(ctx, aggregate, v1.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(createdAggregate.Spec).Should(Equal(aggregate.Spec))
				finished <- true
			}()

			event := &watch.Event{}
			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Added))
			eventAggregate := event.Object.(*v1alpha1.Aggregate)
			Expect(eventAggregate).NotTo(BeNil())
			Expect(eventAggregate.Spec).Should(Equal(aggregate.Spec))

			<-finished

			By("Updating Aggregate")
			createdAggregate.Spec.Aggregates[0].SourcePath = *v1alpha1.JSONPathFromString("spec.nets")
			go func() {
				defer GinkgoRecover()
				updatedAggregate, err := client.Update(ctx, createdAggregate, v1.UpdateOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedAggregate.Spec).Should(Equal(createdAggregate.Spec))
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventAggregate = event.Object.(*v1alpha1.Aggregate)
			Expect(eventAggregate).NotTo(BeNil())
			Expect(eventAggregate.Spec).Should(Equal(createdAggregate.Spec))

			<-finished

			By("Patching Aggregate")
			patch := []struct {
				Op    string `json:"op"`
				Path  string `json:"path"`
				Value string `json:"value"`
			}{{
				Op:    "replace",
				Path:  "/spec/aggregates/0/targetPath",
				Value: "status.computed.nets.netCount",
			}}

			patchData, err := json.Marshal(patch)
			Expect(err).NotTo(HaveOccurred())

			go func() {
				defer GinkgoRecover()
				patchedAggregate, err := client.Patch(ctx, AggregateName, types.JSONPatchType, patchData, v1.PatchOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(patchedAggregate.Spec.Aggregates[0].TargetPath).Should(BeEquivalentTo(*v1alpha1.JSONPathFromString(patch[0].Value)))
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventAggregate = event.Object.(*v1alpha1.Aggregate)
			Expect(eventAggregate).NotTo(BeNil())
			Expect(eventAggregate.Spec.Aggregates[0].TargetPath).Should(Equal(*v1alpha1.JSONPathFromString(patch[0].Value)))

			<-finished

			// We do not handle status for Aggregate atm,
			// so just a placeholder for now
			By("Updating Aggregate status")
			_, err = client.UpdateStatus(ctx, eventAggregate, v1.UpdateOptions{})
			Expect(err).NotTo(HaveOccurred())
			Eventually(events).Should(Receive())

			aggregateToDelete := &v1alpha1.Aggregate{
				ObjectMeta: v1.ObjectMeta{
					Name:      AggregateToDeleteName,
					Namespace: AggregateNamespace,
					Labels: map[string]string{
						DeleteLabel: "",
					},
				},
				Spec: v1alpha1.AggregateSpec{
					Aggregates: []v1alpha1.AggregateItem{
						{
							SourcePath: *v1alpha1.JSONPathFromString("a.b.c"),
							TargetPath: *v1alpha1.JSONPathFromString("q.w.e"),
							Aggregate:  v1alpha1.CSumAggregateType,
						},
					},
				},
			}

			By("Creating Aggregate collection")
			_, err = client.Create(ctx, aggregateToDelete, v1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
			Eventually(events).Should(Receive())

			By("Listing Aggregates")
			aggregateList, err := client.List(ctx, v1.ListOptions{})
			Expect(aggregateList).NotTo(BeNil())
			Expect(aggregateList.Items).To(HaveLen(2))

			By("Bulk deleting Aggregate")
			Expect(client.DeleteCollection(ctx, v1.DeleteOptions{}, v1.ListOptions{LabelSelector: DeleteLabel})).To(Succeed())

			By("Requesting created Aggregate")
			Eventually(func() bool {
				_, err = client.Get(ctx, AggregateName, v1.GetOptions{})
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Eventually(func() bool {
				_, err = client.Get(ctx, AggregateToDeleteName, v1.GetOptions{})
				return err == nil
			}, timeout, interval).Should(BeFalse())

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Deleted))
			eventAggregate = event.Object.(*v1alpha1.Aggregate)
			Expect(eventAggregate).NotTo(BeNil())
			Expect(eventAggregate.Name).To(Equal(AggregateToDeleteName))

			By("Deleting Aggregate")
			go func() {
				defer GinkgoRecover()
				err := client.Delete(ctx, AggregateName, v1.DeleteOptions{})
				Expect(err).NotTo(HaveOccurred())
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Deleted))
			eventAggregate = event.Object.(*v1alpha1.Aggregate)
			Expect(eventAggregate).NotTo(BeNil())
			Expect(eventAggregate.Name).To(Equal(AggregateName))

			<-finished

			watcher.Stop()
			Eventually(events).Should(BeClosed())
		})
	})
})
