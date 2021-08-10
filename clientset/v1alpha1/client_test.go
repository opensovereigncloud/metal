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

	"github.com/onmetal/switch-operator/api/v1alpha1"
)

var _ = Describe("Switch client", func() {
	const (
		SwitchName         = "test-switch"
		SwitchToDeleteName = "test-switch-to-delete"
		DeleteLabel        = "delete-label"
		SwitchesNamespace  = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When Switch CR is installed", func() {
		It("Should check that Switch CR is operational with client", func() {
			By("Creating client")
			finished := make(chan bool)
			ctx := context.Background()

			clientset, err := NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			client := clientset.Switches(SwitchesNamespace)

			res := &v1alpha1.Switch{
				ObjectMeta: v1.ObjectMeta{
					Name:      SwitchName,
					Namespace: SwitchesNamespace,
				},
				Spec: v1alpha1.SwitchSpec{
					Hostname:    SwitchName,
					Location:    &v1alpha1.LocationSpec{},
					TotalPorts:  5,
					SwitchPorts: 4,
					Distro: &v1alpha1.SwitchDistroSpec{
						OS:      "SONiC",
						Version: "1.0.0",
						ASIC:    "broadcom",
					},
					Chassis: &v1alpha1.SwitchChassisSpec{
						Manufacturer: "Edgecore",
						SKU:          "1",
						Serial:       "00000X00001",
						ChassisID:    "68:21:5f:47:0d:6e",
					},
					Interfaces: map[string]*v1alpha1.InterfaceSpec{
						"eth0": {
							Speed:               1000,
							MTU:                 1500,
							Lanes:               1,
							FEC:                 "none",
							MacAddress:          "68:21:5f:47:0d:6e",
							IPv4:                "",
							IPv6:                "",
							PeerType:            "Switch",
							PeerSystemName:      "mgmt-sw-1",
							PeerChassisID:       "64:9d:99:15:56:1c",
							PeerPortID:          "Eth 1/44",
							PeerPortDescription: "Ethernet Port on unit 1, port 44",
						},
						"Ethernet0": {
							Speed:      100000,
							MTU:        9100,
							Lanes:      4,
							FEC:        "none",
							MacAddress: "68:21:5f:47:0d:6e",
						},
						"Ethernet4": {
							Speed:      100000,
							MTU:        9100,
							Lanes:      4,
							FEC:        "none",
							MacAddress: "68:21:5f:47:0d:6e",
						},
						"Ethernet8": {
							Speed:      100000,
							MTU:        9100,
							Lanes:      4,
							FEC:        "none",
							MacAddress: "68:21:5f:47:0d:6e",
						},
						"Ethernet12": {
							Speed:      100000,
							MTU:        9100,
							Lanes:      4,
							FEC:        "none",
							MacAddress: "68:21:5f:47:0d:6e",
						},
					},
				},
			}

			By("Creating watcher")
			watcher, err := client.Watch(ctx, v1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			events := watcher.ResultChan()

			By("Creating Switch")
			createdSwitch := &v1alpha1.Switch{}
			go func() {
				defer GinkgoRecover()
				createdSwitch, err = client.Create(ctx, res, v1.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(createdSwitch.Spec).Should(Equal(res.Spec))
				finished <- true
			}()

			event := &watch.Event{}
			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Added))
			eventSwitch := event.Object.(*v1alpha1.Switch)
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.Spec).Should(Equal(res.Spec))

			<-finished

			By("Updating Switch")
			createdSwitch, err = client.Get(ctx, SwitchName, v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			createdSwitch.Spec.Location.Room = "dummy-room"
			createdSwitch.Spec.Location.HU = 5
			createdSwitch.Spec.Location.Rack = 2
			createdSwitch.Spec.Location.Row = 9
			go func() {
				defer GinkgoRecover()
				updatedSwitch, err := client.Update(ctx, createdSwitch, v1.UpdateOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedSwitch.Spec).Should(Equal(createdSwitch.Spec))
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventSwitch = event.Object.(*v1alpha1.Switch)
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.Spec).Should(Equal(createdSwitch.Spec))

			<-finished

			By("Patching Switch")
			patch := []struct {
				Op    string `json:"op"`
				Path  string `json:"path"`
				Value string `json:"value"`
			}{{
				Op:    "replace",
				Path:  "/spec/location/room",
				Value: "patched-dummy-room",
			}}

			patchData, err := json.Marshal(patch)
			Expect(err).NotTo(HaveOccurred())

			go func() {
				defer GinkgoRecover()
				patchedSwitch, err := client.Patch(ctx, SwitchName, types.JSONPatchType, patchData, v1.PatchOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(patchedSwitch.Spec.Location.Room).Should(Equal(patch[0].Value))
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventSwitch = event.Object.(*v1alpha1.Switch)
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.Spec.Location.Room).Should(Equal(patch[0].Value))

			<-finished

			By("Updating Switch status")
			createdSwitch, err = client.Get(ctx, SwitchName, v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			createdSwitch.FillStatusOnCreate()
			go func() {
				defer GinkgoRecover()
				updatedSwitch, err := client.UpdateStatus(ctx, createdSwitch, v1.UpdateOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedSwitch.Status).Should(Equal(createdSwitch.Status))
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventSwitch = event.Object.(*v1alpha1.Switch)
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.Status).Should(Equal(createdSwitch.Status))

			switchToDelete := res.DeepCopy()
			switchToDelete.Name = SwitchToDeleteName
			switchToDelete.Labels = map[string]string{
				DeleteLabel: "",
			}

			By("Creating Switch collection")
			_, err = client.Create(ctx, switchToDelete, v1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
			Eventually(events).Should(Receive())

			By("Listing Switches")
			switchesList, err := client.List(ctx, v1.ListOptions{})
			Expect(switchesList).NotTo(BeNil())
			Expect(switchesList.Items).To(HaveLen(2))

			By("Bulk deleting Switches")
			Expect(client.DeleteCollection(ctx, v1.DeleteOptions{}, v1.ListOptions{LabelSelector: DeleteLabel})).To(Succeed())

			By("Requesting created Switch")
			Eventually(func() bool {
				_, err = client.Get(ctx, SwitchName, v1.GetOptions{})
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Eventually(func() bool {
				_, err = client.Get(ctx, SwitchToDeleteName, v1.GetOptions{})
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeFalse())

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Deleted))
			eventSwitch = event.Object.(*v1alpha1.Switch)
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.Name).To(Equal(SwitchToDeleteName))

			By("Deleting Switch")
			go func() {
				defer GinkgoRecover()
				err := client.Delete(ctx, SwitchName, v1.DeleteOptions{})
				Expect(err).NotTo(HaveOccurred())
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Deleted))
			eventSwitch = event.Object.(*v1alpha1.Switch)
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.Name).To(Equal(SwitchName))

			<-finished

			watcher.Stop()
			Eventually(events).Should(BeClosed())
		})
	})
})
