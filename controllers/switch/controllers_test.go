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
	"bytes"
	"context"
	"os"
	"path/filepath"
	"time"

	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/utils/pointer"

	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/pkg/constants"
	switchespkg "github.com/onmetal/metal-api/pkg/switches"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	timeout  = time.Second * 60
	interval = time.Millisecond * 1000
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

// The following onboarding cases are covered:
// - Inventory object is created. Onboarding-controller should reconcile object and create corresponding Switch object;
//
// - Onboarded Switch object was updated: onboarding metadata was deleted. Onboarding-controller should restore label
//   and annotations related to onboarding process;
//
// - Inventory object exists, for some reason automatically created Switch object was deleted, new Switch object was
//   created manually without onboarding metadata (labels/annotations/inventory reference). Onboarding-controller
//   should handle "Create" event and update existing Switch object with proper metadata and Inventory reference;
//   Constraints:
//   - Switch object should either have the same name as Inventory object OR contain .spec.inventoryRef.name field
//     filled with proper Inventory object name;
//
// - Switch object exists, it was created without required labels/annotations/inventory reference. After creation of
//   Inventory object, onboarding-controller should update existing Switch object with proper metadata and Inventory
//   reference;
//   Constraints:
//   - Switch object should either have the same name as Inventory object OR contain .spec.inventoryRef.name field
//     filled with proper Inventory object name;
//
// The following cases for configuration processing are covered:
// - IPAM objects are pre-created, switch-controller consumes existing subnets and IPs;
// - IPAM objects are created during switch reconciliation;
// - Interfaces' parameters are changed by changing of parameters defined in SwitchConfig object spec;

var _ = Describe("Switch controller", func() {

	AfterEach(func() {
		postTestContext, postTestCancel := context.WithCancel(ctx)
		defer postTestCancel()
		deleteInventories(postTestContext)
		deleteSwitches(postTestContext)
		deleteLoopbackIPs(postTestContext)
		deleteSouthSubnets(postTestContext)
	})

	Context("Creating switches from inventories", func() {
		It("Switches should be created from inventories by onboarding-controller", func() {
			testContext, testCancel := context.WithCancel(ctx)
			defer testCancel()

			By("Seed Inventory objects")
			Expect(seedInventories(testContext, k8sClient)).NotTo(HaveOccurred())

			By("Expect switches exist")
			checkSwitches()
		})
	})

	Context("Onboarding metadata should be persistent", func() {
		It("Onboarding-controller should restore labels and annotations if deleted", func() {
			testContext, testCancel := context.WithCancel(ctx)
			defer testCancel()

			By("Seed Inventory objects")
			Expect(seedInventories(testContext, k8sClient)).NotTo(HaveOccurred())

			By("Expect switches exist")
			checkSwitches()

			By("Remove onboarding metadata")
			switches := &switchv1beta1.SwitchList{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.List(testContext, switches)).To(Succeed())
				for _, item := range switches.Items {
					delete(item.Labels, constants.InventoriedLabel)
					delete(item.Annotations, constants.HardwareChassisIDAnnotation)
					g.Expect(k8sClient.Update(testContext, &item)).To(Succeed())
				}
			}, timeout, interval).Should(Succeed())

			By("Expect switches' onboarding metadata was restored")
			checkSwitches()
		})
	})

	Context("Onboarding-controller handles CREATE events from switches", func() {
		JustBeforeEach(func() {
			preTestContext, preTestCancel := context.WithCancel(ctx)
			defer preTestCancel()
			Expect(seedInventories(preTestContext, k8sClient)).NotTo(HaveOccurred())
			switches := &switchv1beta1.SwitchList{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.List(preTestContext, switches)).To(Succeed())
				g.Expect(len(switches.Items)).To(Equal(4))
			}, timeout, interval).Should(Succeed())
			deleteSwitches(preTestContext)
		})

		It("Created switches should be updated with onboarding metadata after creation", func() {
			testContext, testCancel := context.WithCancel(ctx)
			defer testCancel()

			By("Seed Switch objects without onboarding metadata - labels, annotations and inventory reference")
			names := []string{
				"b9a234a5-416b-3d49-a4f8-65b6f30c8ee5",
				"044ca7d1-c6f8-37d8-83ce-bf6a18318f2d",
				"a177382d-a3b4-3ecd-97a4-01cc15e749e4",
				"92b9de0f-19f2-3f3b-95d0-fb668b1d3d3b",
			}
			for _, name := range names {
				topSpine := false
				if name == "a177382d-a3b4-3ecd-97a4-01cc15e749e4" || name == "92b9de0f-19f2-3f3b-95d0-fb668b1d3d3b" {
					topSpine = true
				}
				obj := &switchv1beta1.Switch{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: defaultNamespace,
					},
					Spec: switchv1beta1.SwitchSpec{
						TopSpine: pointer.Bool(topSpine),
					},
				}
				Expect(k8sClient.Create(testContext, obj)).To(Succeed())
			}

			By("Expect Switch objects are updated with onboarding metadata")
			switches := &switchv1beta1.SwitchList{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.List(testContext, switches)).To(Succeed())
				for _, item := range switches.Items {
					g.Expect(item.Labels[constants.InventoriedLabel]).NotTo(BeEmpty())
					g.Expect(item.Annotations[constants.HardwareChassisIDAnnotation]).NotTo(BeEmpty())
					g.Expect(item.GetInventoryRef()).NotTo(BeEmpty())
				}
			}, timeout, interval).Should(Succeed())
		})
	})

	Context("Existing switches are updated by onboarding-controller on inventory creation", func() {
		It("Existing switches should be updated with onboarding metadata after inventories are created", func() {
			testContext, testCancel := context.WithCancel(ctx)
			defer testCancel()

			By("Seed Switch objects without onboarding metadata - labels, annotations and inventory reference")
			names := []string{
				"b9a234a5-416b-3d49-a4f8-65b6f30c8ee5",
				"044ca7d1-c6f8-37d8-83ce-bf6a18318f2d",
				"a177382d-a3b4-3ecd-97a4-01cc15e749e4",
				"92b9de0f-19f2-3f3b-95d0-fb668b1d3d3b",
			}
			for _, name := range names {
				topSpine := false
				if name == "a177382d-a3b4-3ecd-97a4-01cc15e749e4" || name == "92b9de0f-19f2-3f3b-95d0-fb668b1d3d3b" {
					topSpine = true
				}
				obj := &switchv1beta1.Switch{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: defaultNamespace,
					},
					Spec: switchv1beta1.SwitchSpec{
						TopSpine: pointer.Bool(topSpine),
					},
				}
				Expect(k8sClient.Create(testContext, obj)).To(Succeed())
			}

			By("Seed Inventory objects")
			Expect(seedInventories(testContext, k8sClient)).NotTo(HaveOccurred())

			By("Expect Switch objects are updated with onboarding metadata")
			switches := &switchv1beta1.SwitchList{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.List(testContext, switches)).To(Succeed())
				for _, item := range switches.Items {
					g.Expect(item.Labels[constants.InventoriedLabel]).NotTo(BeEmpty())
					g.Expect(item.Annotations[constants.HardwareChassisIDAnnotation]).NotTo(BeEmpty())
					g.Expect(item.GetInventoryRef()).NotTo(BeEmpty())
				}
			}, timeout, interval).Should(Succeed())
		})
	})

	Context("Onboarding-controller updates condition of the switch on inventory update", func() {
		JustBeforeEach(func() {
			preTestContext, preTestCancel := context.WithCancel(ctx)
			defer preTestCancel()
			Expect(seedSwitches(preTestContext, k8sClient)).NotTo(HaveOccurred())
			Expect(seedInventories(preTestContext, k8sClient)).NotTo(HaveOccurred())
		})

		It("Should update interfaces so they are aligned with corresponding inventory", func() {
			testContext, testCancel := context.WithCancel(ctx)
			defer testCancel()

			By("Expect switches' state 'Pending' due to missing type label")
			checkState(constants.SwitchStatePending)
			setConfigSelector()

			By("Expect switch spec matches initial state")
			checkInterfaces()
			target := &switchv1beta1.Switch{}
			Consistently(func(g Gomega) {
				g.Expect(k8sClient.Get(testContext, types.NamespacedName{
					Namespace: "onmetal",
					Name:      "b9a234a5-416b-3d49-a4f8-65b6f30c8ee5",
				}, target)).NotTo(HaveOccurred())
				targetInterface := target.Status.Interfaces["Ethernet100"]
				g.Expect(targetInterface.Peer).To(BeNil())
			}, timeout, interval).Should(Succeed())

			By("Expect switches' configuration matches updated inventory")
			updateInventory()
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(testContext, types.NamespacedName{
					Namespace: "onmetal",
					Name:      "b9a234a5-416b-3d49-a4f8-65b6f30c8ee5",
				}, target)).NotTo(HaveOccurred())
				targetInterface := target.Status.Interfaces["Ethernet100"]
				g.Expect(targetInterface.Peer).NotTo(BeNil())
				g.Expect(targetInterface.Peer.GetChassisID()).To(Equal("2a30fd70-008e-4975-ba77-8f5683505e37"))
			}, timeout*2, interval).Should(Succeed())
		})
	})

	Context("Computing switches' configuration without pre-created IPAM objects", func() {
		JustBeforeEach(func() {
			preTestContext, preTestCancel := context.WithCancel(ctx)
			defer preTestCancel()
			Expect(seedSwitches(preTestContext, k8sClient)).NotTo(HaveOccurred())
			Expect(seedInventories(preTestContext, k8sClient)).NotTo(HaveOccurred())
		})

		It("Should compute configs and create missing IPAM objects", func() {
			By("Expect successful switches' configuration")
			checkInterfaces()
			checkLayerAndRole()

			By("Expect switches' state 'Pending'")
			checkState(constants.SwitchStatePending)
			setConfigSelector()

			checkConfigRef()
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
			Expect(seedSwitches(preTestContext, k8sClient)).NotTo(HaveOccurred())
			Expect(seedInventories(preTestContext, k8sClient)).NotTo(HaveOccurred())
			Expect(seedSwitchesSubnets(preTestContext, k8sClient)).NotTo(HaveOccurred())
			Expect(seedSwitchesLoopbacks(preTestContext, k8sClient)).NotTo(HaveOccurred())
		})

		It("Should compute configs and use existing IPAM objects", func() {
			By("Expect successful switches' configuration")
			checkInterfaces()
			checkLayerAndRole()

			By("Expect switches' state 'Pending'")
			checkState(constants.SwitchStatePending)
			setConfigSelector()

			checkConfigRef()
			checkLoopbacks()
			checkASN()
			checkSubnets()
			checkIPAddresses()
			checkState(constants.SwitchStateReady)

			By("Expect pre-created IPAM objects used in switches' configuration")
			checkSeededLoopbacks()
		})
	})

	Context("Updating mapping between switches and switch configs", func() {
		JustBeforeEach(func() {
			preTestContext, preTestCancel := context.WithCancel(ctx)
			defer preTestCancel()
			Expect(seedInventories(preTestContext, k8sClient)).NotTo(HaveOccurred())
			checkSwitches()
			setTopSpines()
		})

		It("Should update mapping between switches and switch configs", func() {
			By("Expect switches' state 'Pending' due to inability to discover related switch config")
			checkState(constants.SwitchStatePending)
			checkLayerAndRole()
			checkConfigSelectorPopulated()
			setConfigSelector()

			By("Configuration process proceeds further")
			checkConfigRef()
			checkState(constants.SwitchStateReady)

			By("Removing of config selector causes 'Pending' state for switches")
			flushConfigSelector()
			checkConfigSelectorPopulated()
			checkState(constants.SwitchStatePending)

			By("Updating of switch configs labels should lead to proper switch configuration")
			updateSwitchConfigLabels()
			checkState(constants.SwitchStateReady)
		})
	})
})

func setTopSpines() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	spineOne := &switchv1beta1.Switch{}
	Expect(k8sClient.Get(testContext, types.NamespacedName{
		Namespace: "onmetal",
		Name:      "a177382d-a3b4-3ecd-97a4-01cc15e749e4",
	}, spineOne)).To(Succeed())
	spineOne.SetTopSpine(true)
	Expect(k8sClient.Update(testContext, spineOne)).To(Succeed())
	spineTwo := &switchv1beta1.Switch{}
	Expect(k8sClient.Get(testContext, types.NamespacedName{
		Namespace: "onmetal",
		Name:      "92b9de0f-19f2-3f3b-95d0-fb668b1d3d3b",
	}, spineTwo)).To(Succeed())
	spineTwo.SetTopSpine(true)
	Expect(k8sClient.Update(testContext, spineTwo)).To(Succeed())
}

func setConfigSelector() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switches)).NotTo(HaveOccurred())
		for _, item := range switches.Items {
			if item.GetTopSpine() {
				item.Spec.ConfigSelector = &metav1.LabelSelector{
					MatchLabels: map[string]string{constants.SwitchTypeLabel: constants.SwitchRoleSpine},
				}
			} else {
				item.Spec.ConfigSelector = &metav1.LabelSelector{
					MatchLabels: map[string]string{constants.SwitchTypeLabel: constants.SwitchRoleLeaf},
				}
			}
			item.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
			g.Expect(k8sClient.Patch(testContext, &item, client.Apply, switchespkg.PatchOpts)).NotTo(HaveOccurred())
		}
	}, timeout, interval).Should(Succeed())
}

func flushConfigSelector() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switches)).NotTo(HaveOccurred())
		for _, item := range switches.Items {
			item.Spec.ConfigSelector = nil
			item.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
			g.Expect(k8sClient.Patch(testContext, &item, client.Apply, switchespkg.PatchOpts)).NotTo(HaveOccurred())
		}
	}, timeout, interval).Should(Succeed())
}

func checkConfigSelectorPopulated() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Consistently(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switches))
		for _, item := range switches.Items {
			if item.GetTopSpine() {
				g.Expect(item.GetConfigSelector().MatchLabels).NotTo(BeNil())
				g.Expect(item.GetConfigSelector().MatchLabels[constants.SwitchConfigLayerLabel]).To(Equal("0"))
			} else {
				g.Expect(item.GetConfigSelector().MatchLabels).NotTo(BeNil())
				g.Expect(item.GetConfigSelector().MatchLabels[constants.SwitchConfigLayerLabel]).To(Equal("1"))
			}
		}
	}, timeout, interval).Should(Succeed())
}

func updateSwitchConfigLabels() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switchConfigs := &switchv1beta1.SwitchConfigList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switchConfigs)).NotTo(HaveOccurred())
		for _, item := range switchConfigs.Items {
			typeValue := item.Labels[constants.SwitchTypeLabel]
			if typeValue == "spine" {
				item.Labels[constants.SwitchConfigLayerLabel] = "0"
			}
			if typeValue == "leaf" {
				item.Labels[constants.SwitchConfigLayerLabel] = "1"
			}
			delete(item.Labels, constants.SwitchTypeLabel)
			g.Expect(k8sClient.Update(ctx, &item)).NotTo(HaveOccurred())
		}
	}, timeout, interval).Should(Succeed())
}

func checkSwitches() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	switches := &switchv1beta1.SwitchList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(testContext, switches)).NotTo(HaveOccurred())
		g.Expect(switches.Items).NotTo(BeEmpty())
		for _, item := range switches.Items {
			g.Expect(item.GetInventoryRef()).NotTo(BeEmpty())
			g.Expect(item.Labels[constants.InventoriedLabel]).NotTo(BeEmpty())
			g.Expect(item.Annotations[constants.HardwareChassisIDAnnotation]).NotTo(BeEmpty())
		}
	}, timeout, interval).Should(Succeed())
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
	}, timeout*2, interval).Should(Succeed())
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

func updateInventory() {
	testContext, testCancel := context.WithTimeout(ctx, timeout*2)
	defer testCancel()
	samplePath := filepath.Join(samplesPath, "updatedInventory", "leaf-1.inventory.yaml")
	raw, err := os.ReadFile(samplePath)
	Expect(err).To(BeNil())
	updatedInventory := &inventoryv1alpha1.Inventory{}
	sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(raw), len(raw))
	Expect(sampleYAML.Decode(updatedInventory)).To(Succeed())
	updatedInventory.TypeMeta = metav1.TypeMeta{
		Kind:       "Inventory",
		APIVersion: "machine.onmetal.de/v1alpha1",
	}
	Expect(k8sClient.Patch(testContext, updatedInventory, client.Apply, switchespkg.PatchOpts)).NotTo(HaveOccurred())
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

func deleteInventories(ctx context.Context) {
	selector := labels.NewSelector()
	req, _ := labels.NewRequirement(constants.SizeLabel, selection.Exists, []string{})
	selector = selector.Add(*req)
	opts := client.ListOptions{
		LabelSelector: selector,
		Namespace:     defaultNamespace,
	}
	delOpts := &client.DeleteAllOfOptions{
		ListOptions: opts,
	}
	Expect(k8sClient.DeleteAllOf(ctx, &inventoryv1alpha1.Inventory{}, delOpts, client.InNamespace(defaultNamespace))).NotTo(HaveOccurred())
	inventories := &inventoryv1alpha1.InventoryList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(ctx, inventories, &opts)).NotTo(HaveOccurred())
		g.Expect(inventories.Items).To(BeEmpty())
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
