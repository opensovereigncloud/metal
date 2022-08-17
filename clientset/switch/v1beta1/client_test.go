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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"

	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
)

// nolint:forcetypeassert
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

			res := &switchv1beta1.Switch{
				ObjectMeta: v1.ObjectMeta{
					Name:      SwitchName,
					Namespace: SwitchesNamespace,
				},
				Spec: switchv1beta1.SwitchSpec{
					UUID:     "a177382d-a3b4-3ecd-97a4-01cc15e749e4",
					TopSpine: false,
					Managed:  true,
					Cordon:   false,
				},
			}

			By("Creating watcher")
			watcher, err := client.Watch(ctx, v1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			events := watcher.ResultChan()

			By("Creating Switch")
			createdSwitch := &switchv1beta1.Switch{}
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
			eventSwitch := event.Object.(*switchv1beta1.Switch)
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.Spec).Should(Equal(res.Spec))

			<-finished

			By("Updating Switch")
			createdSwitch, err = client.Get(ctx, SwitchName, v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			createdSwitch.Spec.Cordon = true
			go func() {
				defer GinkgoRecover()
				var updatedSwitch *switchv1beta1.Switch
				updatedSwitch, err = client.Update(ctx, createdSwitch, v1.UpdateOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedSwitch.Spec).Should(Equal(createdSwitch.Spec))
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventSwitch = event.Object.(*switchv1beta1.Switch)
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
				patchedSwitch, err = client.Patch(ctx, SwitchName, types.JSONPatchType, patchData, v1.PatchOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(patchedSwitch.Spec.Managed).Should(BeFalse())
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventSwitch = event.Object.(*switchv1beta1.Switch)
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.Spec.Managed).Should(BeFalse())

			<-finished

			By("Updating Switch status")
			createdSwitch, err = client.Get(ctx, SwitchName, v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			createdSwitch.Status = switchv1beta1.SwitchStatus{
				TotalPorts:      1,
				SwitchPorts:     1,
				Role:            switchv1beta1.MetalAPIString("spine"),
				ConnectionLevel: 0,
				Interfaces: map[string]*switchv1beta1.InterfaceSpec{"Ethernet0": {
					MACAddress: "00:00:00:00:00:01",
					FEC:        switchv1beta1.CFECNone,
					MTU:        9100,
					Speed:      100000,
					Lanes:      4,
					State:      switchv1beta1.CNICUp,
					Direction:  switchv1beta1.CDirectionSouth,
				}},
				SwitchState: &switchv1beta1.SwitchStateSpec{
					State: switchv1beta1.MetalAPIString("initial"),
				},
			}
			go func() {
				defer GinkgoRecover()
				var updatedSwitch *switchv1beta1.Switch
				updatedSwitch, err = client.UpdateStatus(ctx, createdSwitch, v1.UpdateOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedSwitch.Status).Should(Equal(createdSwitch.Status))
				finished <- true
			}()

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Modified))
			eventSwitch = event.Object.(*switchv1beta1.Switch)
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
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Eventually(func() bool {
				_, err = client.Get(ctx, SwitchToDeleteName, v1.GetOptions{})
				return err == nil
			}, timeout, interval).Should(BeFalse())

			Eventually(events).Should(Receive(event))
			Expect(event.Type).To(Equal(watch.Deleted))
			eventSwitch = event.Object.(*switchv1beta1.Switch)
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
			eventSwitch = event.Object.(*switchv1beta1.Switch)
			Expect(eventSwitch).NotTo(BeNil())
			Expect(eventSwitch.Name).To(Equal(SwitchName))

			<-finished

			watcher.Stop()
			Eventually(events).Should(BeClosed())
		})
	})
})
