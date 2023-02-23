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

package v1beta1

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/utils/pointer"

	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/internal/constants"
)

var _ = PDescribe("Switch client", func() {
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

			res := &switchv1beta1.Switch{
				ObjectMeta: metav1.ObjectMeta{
					Name:      SwitchName,
					Namespace: SwitchesNamespace,
				},
				Spec: switchv1beta1.SwitchSpec{
					InventoryRef: &v1.LocalObjectReference{Name: "a177382d-a3b4-3ecd-97a4-01cc15e749e4"},
					TopSpine:     pointer.Bool(false),
					Managed:      pointer.Bool(true),
					Cordon:       pointer.Bool(false),
					ScanPorts:    pointer.Bool(true),
				},
			}

			By("Creating watcher")
			watcher, err := client.Watch(ctx, metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			events := watcher.ResultChan()

			By("Creating Switch")
			createdSwitch := &switchv1beta1.Switch{}
			go func() {
				defer GinkgoRecover()
				createdSwitch, err = client.Create(ctx, res, metav1.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(createdSwitch.Spec).Should(Equal(res.Spec))
				finished <- true
			}()

			event := &watch.Event{}
			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Added))
			eventSwitch, ok := event.Object.(*switchv1beta1.Switch)
			Expect(ok).To(BeTrue())
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.Spec).Should(Equal(res.Spec))

			<-finished

			By("Updating Switch")
			createdSwitch, err = client.Get(ctx, SwitchName, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			createdSwitch.SetCordon(true)
			go func() {
				defer GinkgoRecover()
				var updatedSwitch *switchv1beta1.Switch
				updatedSwitch, err = client.Update(ctx, createdSwitch, metav1.UpdateOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedSwitch.Spec).Should(Equal(createdSwitch.Spec))
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventSwitch, ok = event.Object.(*switchv1beta1.Switch)
			Expect(ok).To(BeTrue())
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.Spec).Should(Equal(createdSwitch.Spec))

			<-finished

			By("Patching Switch")
			patch := []struct {
				Op    string `json:"op"`
				Path  string `json:"path"`
				Value bool   `json:"value"`
			}{{
				Op:    "replace",
				Path:  "/spec/managed",
				Value: false,
			}}

			patchData, err := json.Marshal(patch)
			Expect(err).NotTo(HaveOccurred())

			go func() {
				defer GinkgoRecover()
				var patchedSwitch *switchv1beta1.Switch
				patchedSwitch, err = client.Patch(ctx, SwitchName, types.JSONPatchType, patchData, metav1.PatchOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(patchedSwitch.GetManaged()).Should(BeFalse())
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventSwitch, ok = event.Object.(*switchv1beta1.Switch)
			Expect(ok).To(BeTrue())
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.GetManaged()).Should(BeFalse())

			<-finished

			By("Updating Switch status")
			createdSwitch, err = client.Get(ctx, SwitchName, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			createdSwitch.Status = switchv1beta1.SwitchStatus{
				TotalPorts:  pointer.Uint32(1),
				SwitchPorts: pointer.Uint32(1),
				Role:        pointer.String("spine"),
				Layer:       pointer.Uint32(0),
				Interfaces: map[string]*switchv1beta1.InterfaceSpec{"Ethernet0": {
					MACAddress: pointer.String("00:00:00:00:00:01"),
					Direction:  pointer.String(constants.DirectionSouth),
					Speed:      pointer.Uint32(100000),
					PortParametersSpec: &switchv1beta1.PortParametersSpec{
						FEC:   pointer.String(constants.FECNone),
						MTU:   pointer.Uint32(9100),
						Lanes: pointer.Uint32(4),
						State: pointer.String(constants.NICUp),
					},
				}},
				State: pointer.String("Initial"),
			}
			go func() {
				defer GinkgoRecover()
				var updatedSwitch *switchv1beta1.Switch
				updatedSwitch, err = client.UpdateStatus(ctx, createdSwitch, metav1.UpdateOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedSwitch.Status).Should(Equal(createdSwitch.Status))
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventSwitch, ok = event.Object.(*switchv1beta1.Switch)
			Expect(ok).To(BeTrue())
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.Status).Should(Equal(createdSwitch.Status))

			switchToDelete := res.DeepCopy()
			switchToDelete.Name = SwitchToDeleteName
			switchToDelete.Labels = map[string]string{
				DeleteLabel: "",
			}

			By("Creating Switch collection")
			_, err = client.Create(ctx, switchToDelete, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
			Eventually(events).Should(Receive())

			By("Listing Switches")
			switchesList, err := client.List(ctx, metav1.ListOptions{})
			Expect(switchesList).NotTo(BeNil())
			Expect(switchesList.Items).To(HaveLen(2))

			By("Bulk deleting Switches")
			Expect(client.DeleteCollection(
				ctx, metav1.DeleteOptions{},
				metav1.ListOptions{LabelSelector: DeleteLabel})).To(Succeed())

			By("Requesting created Switch")
			Eventually(func() bool {
				_, err = client.Get(ctx, SwitchName, metav1.GetOptions{})
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Eventually(func() bool {
				_, err = client.Get(ctx, SwitchToDeleteName, metav1.GetOptions{})
				return err == nil
			}, timeout, interval).Should(BeFalse())

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Deleted))
			eventSwitch, ok = event.Object.(*switchv1beta1.Switch)
			Expect(ok).To(BeTrue())
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.Name).To(Equal(SwitchToDeleteName))

			By("Deleting Switch")
			go func() {
				defer GinkgoRecover()
				err := client.Delete(ctx, SwitchName, metav1.DeleteOptions{})
				Expect(err).NotTo(HaveOccurred())
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Deleted))
			eventSwitch, ok = event.Object.(*switchv1beta1.Switch)
			Expect(ok).To(BeTrue())
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.Name).To(Equal(SwitchName))

			<-finished

			watcher.Stop()
			Eventually(events).Should(BeClosed())
		})
	})
})
