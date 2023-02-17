/*
Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"

	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/internal/constants"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	timeout  = time.Second * 30
	interval = time.Millisecond * 500
)

var (
	loopbacksV4 = map[string]string{
		"b9a234a5-416b-3d49-a4f8-65b6f30c8ee5": "100.64.255.3",
		"044ca7d1-c6f8-37d8-83ce-bf6a18318f2d": "100.64.255.4",
		"a177382d-a3b4-3ecd-97a4-01cc15e749e4": "100.64.255.5",
		"92b9de0f-19f2-3f3b-95d0-fb668b1d3d3b": "100.64.255.6",
	}

	loopbacksV6 = map[string]string{
		"b9a234a5-416b-3d49-a4f8-65b6f30c8ee5": "fd00:afc0:e014:fff::3",
		"044ca7d1-c6f8-37d8-83ce-bf6a18318f2d": "fd00:afc0:e014:fff::4",
		"a177382d-a3b4-3ecd-97a4-01cc15e749e4": "fd00:afc0:e014:fff::5",
		"92b9de0f-19f2-3f3b-95d0-fb668b1d3d3b": "fd00:afc0:e014:fff::6",
	}
)

var _ = Describe("Switch controller", func() {

	BeforeEach(func() {
		preTestContext, preTestCancel := context.WithCancel(ctx)
		defer preTestCancel()
		Expect(seedSwitches(preTestContext, k8sClient)).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		postTestContext, postTestCancel := context.WithCancel(ctx)
		defer postTestCancel()
		deleteSwitches(postTestContext)
		deleteLoopbackIPs(postTestContext)
		deleteSouthSubnets(postTestContext)
	})

	Context("Computing switches' configuration without pre-created IPAM objects", func() {
		It("Should compute configs and create missing IPAM objects", func() {
			By("Expect switches' state 'Pending' due to missing type label")
			checkState(constants.SwitchStatePending)
			setTypeLabel()

			By("Expect successful switches' configuration")
			checkInterfaces()
			checkConfigRef()
			checkLayerAndRole()
			checkLoopbacks()
			checkASN()
			checkSubnets()
			checkIPAddresses()
			checkState(constants.SwitchStateReady)

			By("Expect switches' configuration matches updated global config")
			updateSpinesConfig()
			checkInterfacesUpdated()
			checkState(constants.SwitchStateReady)
		})
	})

	Context("Computing switches' configuration with pre-created IPAM objects", func() {
		JustBeforeEach(func() {
			By("Seeding switches' related IPAM objects")
			preTestContext, preTestCancel := context.WithCancel(ctx)
			defer preTestCancel()
			Expect(seedSwitchesSubnets(preTestContext, k8sClient)).NotTo(HaveOccurred())
			Expect(seedSwitchesLoopbacks(preTestContext, k8sClient)).NotTo(HaveOccurred())
		})

		It("Should compute configs and use existing IPAM objects", func() {
			By("Setting type labels on switches")
			setTypeLabel()

			By("Expect successful switches' configuration")
			checkInterfaces()
			checkConfigRef()
			checkLayerAndRole()
			checkLoopbacks()
			checkASN()
			checkSubnets()
			checkIPAddresses()
			checkState(constants.SwitchStateReady)

			By("Expect pre-created IPAM objects used in switches' configuration")
			checkSeededLoopbacks()
		})
	})
})

func setTypeLabel() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Expect(k8sClient.List(testContext, switches)).NotTo(HaveOccurred())
	for _, item := range switches.Items {
		if item.GetTopSpine() {
			item.Labels[constants.SwitchTypeLabel] = constants.SwitchRoleSpine
		} else {
			item.Labels[constants.SwitchTypeLabel] = constants.SwitchRoleLeaf
		}
		item.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
		Expect(k8sClient.Patch(testContext, &item, client.Apply, patchOpts)).NotTo(HaveOccurred())
	}
}

func checkConfigRef() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switches)).NotTo(HaveOccurred())
		for _, item := range switches.Items {
			g.Expect(item.Status.ConfigRef).NotTo(BeNil())
		}
	}, timeout, interval).Should(Succeed())
}

func checkInterfaces() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switches)).NotTo(HaveOccurred())
		for _, item := range switches.Items {
			g.Expect(item.Status.Interfaces).NotTo(BeZero())
		}
	}, timeout, interval).Should(Succeed())
}

func checkLayerAndRole() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switches)).NotTo(HaveOccurred())
		for _, item := range switches.Items {
			g.Expect(item.Status.Layer).NotTo(BeNil())
			if item.GetTopSpine() {
				g.Expect(item.GetLayer()).To(Equal(uint32(0)))
				g.Expect(item.GetRole()).To(Equal("spine"))
			} else {
				g.Expect(item.GetLayer()).To(Equal(uint32(1)))
				g.Expect(item.Status.Role).NotTo(BeNil())
			}
		}
	}, timeout, interval).Should(Succeed())
}

func checkLoopbacks() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switches)).NotTo(HaveOccurred())
		for _, item := range switches.Items {
			g.Expect(item.Status.LoopbackAddresses).NotTo(BeEmpty())
			g.Expect(len(item.Status.LoopbackAddresses)).To(Equal(2))
		}
	}, timeout, interval).Should(Succeed())
}

func checkASN() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switches)).NotTo(HaveOccurred())
		for _, item := range switches.Items {
			g.Expect(item.Status.ASN).NotTo(BeNil())
		}
	}, timeout, interval).Should(Succeed())
}

func checkSubnets() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switches)).NotTo(HaveOccurred())
		for _, item := range switches.Items {
			g.Expect(item.Status.Subnets).NotTo(BeEmpty())
			g.Expect(len(item.Status.Subnets)).To(Equal(2))
		}
	}, timeout, interval).Should(Succeed())
}

func checkIPAddresses() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switches)).NotTo(HaveOccurred())
		for _, item := range switches.Items {
			for _, nic := range item.Status.Interfaces {
				if nic.GetDirection() == constants.DirectionNorth && nic.Peer == nil {
					continue
				}
				g.Expect(nic.IP).NotTo(BeEmpty())
				g.Expect(len(nic.IP)).To(Equal(2))
			}
		}
	}, timeout, interval).Should(Succeed())
}

func checkState(expected string) {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switches)).NotTo(HaveOccurred())
		for _, item := range switches.Items {
			g.Expect(item.GetState()).To(Equal(expected))
		}
	}, timeout, interval).Should(Succeed())
}

func updateSpinesConfig() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	config := &switchv1beta1.SwitchConfig{}
	Expect(k8sClient.Get(testContext, types.NamespacedName{
		Name:      "spines-config",
		Namespace: "onmetal",
	}, config)).NotTo(HaveOccurred())
	config.Spec.PortsDefaults.SetMTU(9216)
	Expect(k8sClient.Update(testContext, config)).NotTo(HaveOccurred())
}

func checkInterfacesUpdated() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switches)).NotTo(HaveOccurred())
		for _, item := range switches.Items {
			for _, nic := range item.Status.Interfaces {
				if !item.GetTopSpine() && nic.GetDirection() == constants.DirectionSouth {
					continue
				}
				g.Expect(nic.GetMTU()).To(Equal(uint32(9216)))
			}
		}
	}, timeout, interval).Should(Succeed())
}

func checkSeededLoopbacks() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switches)).NotTo(HaveOccurred())
		for _, item := range switches.Items {
			for _, lo := range item.Status.LoopbackAddresses {
				switch lo.GetAddressFamily() {
				case constants.IPv4AF:
					g.Expect(lo.GetAddress()).To(Equal(loopbacksV4[item.Name]))
				case constants.IPv6AF:
					g.Expect(lo.GetAddress()).To(Equal(loopbacksV6[item.Name]))
				}
			}
		}
	}, timeout, interval).Should(Succeed())
}

func deleteSwitches(ctx context.Context) {
	selector := labels.NewSelector()
	req, _ := labels.NewRequirement(constants.InventoriedLabel, selection.Exists, []string{})
	selector = selector.Add(*req)
	opts := client.ListOptions{
		LabelSelector: selector,
		Namespace:     defaultNamespace,
	}
	delOpts := &client.DeleteAllOfOptions{
		ListOptions: opts,
	}
	Expect(k8sClient.DeleteAllOf(ctx, &switchv1beta1.Switch{}, delOpts, client.InNamespace(defaultNamespace))).NotTo(HaveOccurred())
	switches := &switchv1beta1.SwitchList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(ctx, switches, &opts)).NotTo(HaveOccurred())
		g.Expect(switches.Items).To(BeEmpty())
	}, timeout, interval).Should(Succeed())
}

func deleteSouthSubnets(ctx context.Context) {
	selector := labels.NewSelector()
	req, _ := labels.NewRequirement(constants.IPAMObjectPurposeLabel, selection.In, []string{constants.IPAMSouthSubnetPurpose})
	selector = selector.Add(*req)
	opts := client.ListOptions{
		LabelSelector: selector,
		Namespace:     defaultNamespace,
	}
	delOpts := &client.DeleteAllOfOptions{
		ListOptions: opts,
	}
	Expect(k8sClient.DeleteAllOf(ctx, &ipamv1alpha1.Subnet{}, delOpts, client.InNamespace(defaultNamespace))).NotTo(HaveOccurred())
	subnets := &ipamv1alpha1.SubnetList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(ctx, subnets, &opts)).NotTo(HaveOccurred())
		g.Expect(subnets.Items).To(BeEmpty())
	}, timeout, interval).Should(Succeed())
}

func deleteLoopbackIPs(ctx context.Context) {
	selector := labels.NewSelector()
	req, _ := labels.NewRequirement(constants.IPAMObjectPurposeLabel, selection.In, []string{constants.IPAMLoopbackPurpose})
	selector = selector.Add(*req)
	opts := client.ListOptions{
		LabelSelector: selector,
		Namespace:     defaultNamespace,
	}
	delOpts := &client.DeleteAllOfOptions{
		ListOptions: opts,
	}
	Expect(k8sClient.DeleteAllOf(ctx, &ipamv1alpha1.IP{}, delOpts, client.InNamespace(defaultNamespace))).NotTo(HaveOccurred())
	loopbacks := &ipamv1alpha1.IPList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(ctx, loopbacks, &opts)).NotTo(HaveOccurred())
		g.Expect(loopbacks.Items).To(BeEmpty())
	}, timeout, interval).Should(Succeed())
}
