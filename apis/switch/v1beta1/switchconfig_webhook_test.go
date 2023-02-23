/*
 * Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1beta1

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/onmetal/metal-api/internal/constants"
)

const (
	defaulNamespace = "default"
	timeout         = time.Second * 30
	interval        = time.Millisecond * 250
)

var _ = Describe("SwitchConfig Webhook", func() {
	AfterEach(func() {
		By("Remove switch configs if exist")
		Expect(k8sClient.DeleteAllOf(ctx, &SwitchConfig{}, client.InNamespace(defaulNamespace))).To(Succeed())
		Eventually(func(g Gomega) {
			list := &SwitchConfigList{}
			g.Expect(k8sClient.List(ctx, list)).NotTo(HaveOccurred())
			g.Expect(len(list.Items)).To(Equal(0))
		}, timeout, interval).Should(Succeed())
	})

	Context("Defaulting switch config", func() {
		It("Should set defaults for switch config", func() {
			switchConfigObject := &SwitchConfig{
				ObjectMeta: v1.ObjectMeta{
					Name:      "sample-config",
					Namespace: defaulNamespace,
				},
				Spec: SwitchConfigSpec{
					Switches: &v1.LabelSelector{
						MatchLabels: map[string]string{"switch.onmetal.de/type": "spine"},
					},
					PortsDefaults: &PortParametersSpec{
						FEC:   pointer.String(constants.FECRS),
						MTU:   pointer.Uint32(9216),
						State: pointer.String(constants.NICUp),
					},
					IPAM: &GeneralIPAMSpec{
						CarrierSubnets: &IPAMSelectionSpec{
							LabelSelector: &v1.LabelSelector{
								MatchLabels: map[string]string{"ipam.onmetal.de/object-purpose": "switch-carrier"},
							},
						},
						LoopbackSubnets: &IPAMSelectionSpec{
							LabelSelector: &v1.LabelSelector{
								MatchLabels: map[string]string{"ipam.onmetal.de/object-purpose": "switch-loopbacks"},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, switchConfigObject)).ShouldNot(HaveOccurred())
			Eventually(func(g Gomega) {
				sampleConfig := &SwitchConfig{}
				// check defaulted ports params
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{Name: switchConfigObject.Name, Namespace: defaulNamespace}, sampleConfig)).Should(Succeed())
				g.Expect(sampleConfig.Spec.PortsDefaults.GetFEC()).Should(Equal(constants.FECRS))
				g.Expect(sampleConfig.Spec.PortsDefaults.GetState()).Should(Equal(constants.NICUp))
				g.Expect(sampleConfig.Spec.PortsDefaults.IPv4MaskLength).NotTo(BeNil())
				g.Expect(sampleConfig.Spec.PortsDefaults.GetIPv4MaskLength()).Should(Equal(uint32(30)))
				g.Expect(sampleConfig.Spec.PortsDefaults.IPv6Prefix).NotTo(BeNil())
				g.Expect(sampleConfig.Spec.PortsDefaults.GetIPv6Prefix()).Should(Equal(uint32(127)))
				g.Expect(sampleConfig.Spec.PortsDefaults.Lanes).NotTo(BeNil())
				g.Expect(sampleConfig.Spec.PortsDefaults.GetLanes()).Should(Equal(uint32(4)))
				g.Expect(sampleConfig.Spec.PortsDefaults.GetMTU()).Should(Equal(uint32(9216)))
				// check defaulted ipam selectors
				g.Expect(sampleConfig.Spec.IPAM.SouthSubnets).NotTo(BeNil())
				g.Expect(sampleConfig.Spec.IPAM.SouthSubnets.LabelSelector.MatchLabels).Should(Equal(map[string]string{constants.IPAMObjectPurposeLabel: constants.IPAMSouthSubnetPurpose}))
				g.Expect(sampleConfig.Spec.IPAM.SouthSubnets.FieldSelector.GetLabelKey()).Should(Equal(constants.IPAMObjectOwnerLabel))
				g.Expect(sampleConfig.Spec.IPAM.SouthSubnets.FieldSelector.FieldRef.FieldPath).Should(Equal(constants.DefaultIPAMFieldRef))
				g.Expect(sampleConfig.Spec.IPAM.LoopbackAddresses).NotTo(BeNil())
				g.Expect(sampleConfig.Spec.IPAM.LoopbackAddresses.LabelSelector.MatchLabels).Should(Equal(map[string]string{constants.IPAMObjectPurposeLabel: constants.IPAMLoopbackPurpose}))
				g.Expect(sampleConfig.Spec.IPAM.LoopbackAddresses.FieldSelector.GetLabelKey()).Should(Equal(constants.IPAMObjectOwnerLabel))
				g.Expect(sampleConfig.Spec.IPAM.LoopbackAddresses.FieldSelector.FieldRef.FieldPath).Should(Equal(constants.DefaultIPAMFieldRef))
				g.Expect(sampleConfig.Spec.IPAM.AddressFamily.GetIPv4()).To(BeTrue())
				g.Expect(sampleConfig.Spec.IPAM.AddressFamily.GetIPv6()).To(BeFalse())
			}, timeout, interval).Should(Succeed())
		})
	})
})
