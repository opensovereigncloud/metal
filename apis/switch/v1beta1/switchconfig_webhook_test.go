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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
						FEC:   MetalAPIString(CFECRS),
						MTU:   MetalAPIUint16(9216),
						State: MetalAPIString(CNICUp),
					},
					IPAM: &GeneralIPAMSpec{
						CarrierSubnets: &v1.LabelSelector{
							MatchLabels: map[string]string{"ipam.onmetal.de/object-purpose": "switch-carrier"},
						},
						LoopbackSubnets: &v1.LabelSelector{
							MatchLabels: map[string]string{"ipam.onmetal.de/object-purpose": "switch-loopbacks"},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, switchConfigObject)).ShouldNot(HaveOccurred())
			Eventually(func(g Gomega) {
				sampleConfig := &SwitchConfig{}
				// check defaulted ports params
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{Name: switchConfigObject.Name, Namespace: defaulNamespace}, sampleConfig)).Should(Succeed())
				g.Expect(GoString(sampleConfig.Spec.PortsDefaults.FEC)).Should(Equal(CFECRS))
				g.Expect(GoString(sampleConfig.Spec.PortsDefaults.State)).Should(Equal(CNICUp))
				g.Expect(sampleConfig.Spec.PortsDefaults.IPv4MaskLength).NotTo(BeNil())
				g.Expect(GoUint8(sampleConfig.Spec.PortsDefaults.IPv4MaskLength)).Should(Equal(uint8(30)))
				g.Expect(sampleConfig.Spec.PortsDefaults.IPv6Prefix).NotTo(BeNil())
				g.Expect(GoUint8(sampleConfig.Spec.PortsDefaults.IPv6Prefix)).Should(Equal(uint8(127)))
				g.Expect(sampleConfig.Spec.PortsDefaults.Lanes).NotTo(BeNil())
				g.Expect(GoUint8(sampleConfig.Spec.PortsDefaults.Lanes)).Should(Equal(uint8(4)))
				g.Expect(GoUint16(sampleConfig.Spec.PortsDefaults.MTU)).Should(Equal(uint16(9216)))
				// check defaulted ipam selectors
				g.Expect(sampleConfig.Spec.IPAM.SouthSubnets).NotTo(BeNil())
				g.Expect(sampleConfig.Spec.IPAM.SouthSubnets.AddressFamilies.IPv4).Should(BeTrue())
				g.Expect(sampleConfig.Spec.IPAM.SouthSubnets.AddressFamilies.IPv6).Should(BeTrue())
				g.Expect(sampleConfig.Spec.IPAM.SouthSubnets.LabelSelector.MatchLabels).Should(Equal(map[string]string{IPAMObjectPurposeLabel: CIPAMPurposeSouthSubnet}))
				g.Expect(sampleConfig.Spec.IPAM.SouthSubnets.FieldSelector.LabelKey).Should(Equal(IPAMObjectOwnerLabel))
				g.Expect(sampleConfig.Spec.IPAM.SouthSubnets.FieldSelector.FieldRef.FieldPath).Should(Equal(CDefaultIPAMFieldRef))
				g.Expect(sampleConfig.Spec.IPAM.LoopbackAddresses).NotTo(BeNil())
				g.Expect(sampleConfig.Spec.IPAM.LoopbackAddresses.AddressFamilies.IPv4).Should(BeTrue())
				g.Expect(sampleConfig.Spec.IPAM.LoopbackAddresses.AddressFamilies.IPv6).Should(BeTrue())
				g.Expect(sampleConfig.Spec.IPAM.LoopbackAddresses.LabelSelector.MatchLabels).Should(Equal(map[string]string{IPAMObjectPurposeLabel: CIPAMPurposeLoopback}))
				g.Expect(sampleConfig.Spec.IPAM.LoopbackAddresses.FieldSelector.LabelKey).Should(Equal(IPAMObjectOwnerLabel))
				g.Expect(sampleConfig.Spec.IPAM.LoopbackAddresses.FieldSelector.FieldRef.FieldPath).Should(Equal(CDefaultIPAMFieldRef))
			}, timeout, interval).Should(Succeed())
		})
	})
})
