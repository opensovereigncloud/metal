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

package networking

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	inventoryv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("IPAM interaction test", func() {
	BeforeEach(func() {
		subnetsList := &ipamv1alpha1.SubnetList{}

		By("Create switches configs")
		configsSamples := make([]string, 0)
		Expect(filepath.Walk("../samples/switchconfigs", func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				configsSamples = append(configsSamples, path)
			}
			return nil
		})).NotTo(HaveOccurred())
		for _, samplePath := range configsSamples {
			sampleBytes, err := ioutil.ReadFile(samplePath)
			Expect(err).NotTo(HaveOccurred())
			sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
			sampleConfig := &switchv1beta1.SwitchConfig{}
			Expect(sampleYAML.Decode(sampleConfig)).NotTo(HaveOccurred())
			Expect(k8sClient.Create(ctx, sampleConfig)).NotTo(HaveOccurred())
		}

		By("Create ipam objects: networks")
		networkSample := filepath.Join("..", "samples", "networks", "overlay-network.yaml")
		sampleBytes, err := ioutil.ReadFile(networkSample)
		Expect(err).NotTo(HaveOccurred())
		sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
		sampleNetwork := &ipamv1alpha1.Network{}
		Expect(sampleYAML.Decode(sampleNetwork)).NotTo(HaveOccurred())
		Expect(k8sClient.Create(ctx, sampleNetwork)).NotTo(HaveOccurred())

		By("Create ipam objects: parent subnets")
		parentSubnets := make([]string, 0)
		Expect(filepath.Walk("../samples/subnets/parent", func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				parentSubnets = append(parentSubnets, path)
			}
			return nil
		})).NotTo(HaveOccurred())
		for _, samplePath := range parentSubnets {
			sampleBytes, err := ioutil.ReadFile(samplePath)
			Expect(err).NotTo(HaveOccurred())
			sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
			sampleSubnet := &ipamv1alpha1.Subnet{}
			Expect(sampleYAML.Decode(sampleSubnet)).NotTo(HaveOccurred())
			Expect(k8sClient.Create(ctx, sampleSubnet)).NotTo(HaveOccurred())
		}

		Eventually(func(g Gomega) {
			g.Expect(k8sClient.List(ctx, subnetsList)).NotTo(HaveOccurred())
			for _, item := range subnetsList.Items {
				g.Expect(item.Status.State).To(Equal(ipamv1alpha1.CFinishedSubnetState))
			}
		}, timeout, interval).Should(Succeed())

		By("Create ipam objects: child subnets")
		childSubnets := make([]string, 0)
		Expect(filepath.Walk("../samples/subnets/child", func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				childSubnets = append(childSubnets, path)
			}
			return nil
		})).NotTo(HaveOccurred())
		for _, samplePath := range childSubnets {
			sampleBytes, err := ioutil.ReadFile(samplePath)
			Expect(err).NotTo(HaveOccurred())
			sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
			sampleSubnet := &ipamv1alpha1.Subnet{}
			Expect(sampleYAML.Decode(sampleSubnet)).NotTo(HaveOccurred())
			Expect(k8sClient.Create(ctx, sampleSubnet)).NotTo(HaveOccurred())
		}

		Eventually(func(g Gomega) {
			g.Expect(k8sClient.List(ctx, subnetsList)).NotTo(HaveOccurred())
			for _, item := range subnetsList.Items {
				g.Expect(item.Status.State).To(Equal(ipamv1alpha1.CFinishedSubnetState))
			}
		}, timeout, interval).Should(Succeed())

		By("Create ipam objects: loopback addresses")
		ipList := &ipamv1alpha1.IPList{}
		loopbacks := make([]string, 0)
		Expect(filepath.Walk("../samples/loopbacks", func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				loopbacks = append(loopbacks, path)
			}
			return nil
		})).NotTo(HaveOccurred())
		for _, samplePath := range loopbacks {
			sampleBytes, err := ioutil.ReadFile(samplePath)
			Expect(err).NotTo(HaveOccurred())
			sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
			sampleLoopback := &ipamv1alpha1.IP{}
			Expect(sampleYAML.Decode(sampleLoopback)).NotTo(HaveOccurred())
			Expect(k8sClient.Create(ctx, sampleLoopback)).NotTo(HaveOccurred())
		}
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.List(ctx, ipList)).NotTo(HaveOccurred())
			for _, item := range ipList.Items {
				g.Expect(item.Status.State).To(Equal(ipamv1alpha1.CFinishedIPState))
			}
		}, timeout, interval).Should(Succeed())
	})

	AfterEach(func() {
		By("Remove inventories if exist")
		Expect(k8sClient.DeleteAllOf(ctx, &inventoryv1alpha1.Inventory{}, client.InNamespace(onmetal))).To(Succeed())
		Eventually(func(g Gomega) {
			list := &inventoryv1alpha1.InventoryList{}
			g.Expect(k8sClient.List(ctx, list)).NotTo(HaveOccurred())
			g.Expect(len(list.Items)).To(Equal(0))
		}, timeout, interval).Should(Succeed())

		By("Remove switches if exist")
		Expect(k8sClient.DeleteAllOf(ctx, &switchv1beta1.Switch{}, client.InNamespace(onmetal))).To(Succeed())
		Eventually(func(g Gomega) {
			list := &switchv1beta1.SwitchList{}
			g.Expect(k8sClient.List(ctx, list)).NotTo(HaveOccurred())
			g.Expect(len(list.Items)).To(Equal(0))
		}, timeout, interval).Should(Succeed())

		By("Remove IP addresses if exist")
		Expect(k8sClient.DeleteAllOf(ctx, &ipamv1alpha1.IP{}, client.InNamespace(onmetal))).To(Succeed())
		Eventually(func(g Gomega) {
			list := &ipamv1alpha1.IPList{}
			g.Expect(k8sClient.List(ctx, list)).NotTo(HaveOccurred())
			g.Expect(len(list.Items)).To(Equal(0))
		}, timeout, interval).Should(Succeed())

		By("Remove child subnets if exist")
		subnets := &ipamv1alpha1.SubnetList{}
		req, _ := labels.NewRequirement("ipam.onmetal.de/object-purpose", selection.In, []string{"south-subnet", "switch-loopbacks"})
		sel := labels.NewSelector().Add(*req)
		opts := &client.ListOptions{
			LabelSelector: sel,
			Namespace:     onmetal,
			Limit:         100,
		}
		Expect(k8sClient.List(ctx, subnets, opts)).To(Succeed())
		for _, item := range subnets.Items {
			Expect(k8sClient.Delete(ctx, &item)).To(Succeed())
		}
		req, _ = labels.NewRequirement("ipam.onmetal.de/object-purpose", selection.In, []string{"switch-carrier"})
		sel = labels.NewSelector().Add(*req)
		opts = &client.ListOptions{
			LabelSelector: sel,
			Namespace:     onmetal,
			Limit:         100,
		}
		Expect(k8sClient.List(ctx, subnets, opts)).To(Succeed())
		for _, item := range subnets.Items {
			Expect(k8sClient.Delete(ctx, &item)).To(Succeed())
		}

		Eventually(func(g Gomega) {
			list := &ipamv1alpha1.SubnetList{}
			g.Expect(k8sClient.List(ctx, list)).NotTo(HaveOccurred())
			g.Expect(len(list.Items)).To(Equal(0))
		}, timeout, interval).Should(Succeed())

		By("Remove networks if exist")
		Expect(k8sClient.DeleteAllOf(ctx, &ipamv1alpha1.Network{}, client.InNamespace(onmetal))).To(Succeed())
		Eventually(func(g Gomega) {
			list := &ipamv1alpha1.NetworkList{}
			g.Expect(k8sClient.List(ctx, list)).NotTo(HaveOccurred())
			g.Expect(len(list.Items)).To(Equal(0))
		}, timeout, interval).Should(Succeed())

		By("Remove switches configs")
		Expect(k8sClient.DeleteAllOf(ctx, &switchv1beta1.SwitchConfig{}, client.InNamespace(onmetal))).To(Succeed())
		Eventually(func(g Gomega) {
			list := &switchv1beta1.SwitchConfigList{}
			g.Expect(k8sClient.List(ctx, list)).NotTo(HaveOccurred())
			g.Expect(len(list.Items)).To(Equal(0))
		}, timeout, interval).Should(Succeed())
	}, 60)

	Context("Set switches' south subnets, loopbacks addresses and switch ports' ips", func() {
		JustBeforeEach(func() {
			inventoriesSamples := make([]string, 0)
			Expect(filepath.Walk("../samples/inventories", func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					inventoriesSamples = append(inventoriesSamples, path)
				}
				return nil
			})).NotTo(HaveOccurred())
			for _, samplePath := range inventoriesSamples {
				sampleBytes, err := ioutil.ReadFile(samplePath)
				Expect(err).NotTo(HaveOccurred())
				sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
				sampleInventory := &inventoryv1alpha1.Inventory{}
				Expect(sampleYAML.Decode(sampleInventory)).NotTo(HaveOccurred())
				Expect(k8sClient.Create(ctx, sampleInventory)).NotTo(HaveOccurred())
			}

			switchesSamples := make([]string, 0)
			Expect(filepath.Walk("../samples/switches", func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					switchesSamples = append(switchesSamples, path)
				}
				return nil
			})).NotTo(HaveOccurred())
			for _, samplePath := range switchesSamples {
				stype := "leaf"
				sampleBytes, err := ioutil.ReadFile(samplePath)
				Expect(err).NotTo(HaveOccurred())
				sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
				sampleSwitch := &switchv1beta1.Switch{}
				Expect(sampleYAML.Decode(sampleSwitch)).NotTo(HaveOccurred())
				if strings.Contains(sampleSwitch.Name, "spine") {
					stype = "spine"
				}
				sampleSwitch.Labels = map[string]string{
					"metalapi.onmetal.de/inventoried":   "true",
					"metalapi.onmetal.de/inventory-ref": sampleSwitch.Spec.UUID,
					"metalapi.onmetal.de/type":          stype,
				}
				Expect(k8sClient.Create(ctx, sampleSwitch)).NotTo(HaveOccurred())
			}
		})

		It("Should set switches' subnets, loopback addresses and switch ports' ips", func() {
			By("Base configuration processing")
			var switchesList = &switchv1beta1.SwitchList{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.List(ctx, switchesList)).NotTo(HaveOccurred())
				for _, item := range switchesList.Items {
					g.Expect(item.Status.SwitchState).NotTo(BeNil())
					g.Expect(switchv1beta1.GoString(item.Status.SwitchState.State)).To(Equal(switchv1beta1.CSwitchStateReady))
					g.Expect(item.ConnectionsOK(switchesList)).Should(BeTrue())
				}
			}, timeout, interval).Should(Succeed())

			By("Networking configuration processing")
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.List(ctx, switchesList)).Should(Succeed())
				for _, item := range switchesList.Items {
					g.Expect(len(item.Status.LoopbackAddresses)).ShouldNot(BeZero())
					g.Expect(len(item.Status.LoopbackAddresses)).Should(Equal(2))
					g.Expect(len(item.Status.Subnets)).ShouldNot(BeZero())
					g.Expect(len(item.Status.Subnets)).Should(Equal(2))
					g.Expect(item.IPaddressesOK(switchesList)).Should(BeTrue())
					for _, nic := range item.Status.Interfaces {
						g.Expect(len(nic.IP)).ShouldNot(BeZero())
						g.Expect(len(nic.IP)).Should(Equal(2))
					}
				}
			}, timeout, interval).Should(Succeed())
		})
	})

	Context("Update switches' loopbacks on changes", func() {
		JustBeforeEach(func() {
			inventoriesSamples := make([]string, 0)
			Expect(filepath.Walk("../samples/inventories", func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					inventoriesSamples = append(inventoriesSamples, path)
				}
				return nil
			})).NotTo(HaveOccurred())
			for _, samplePath := range inventoriesSamples {
				sampleBytes, err := ioutil.ReadFile(samplePath)
				Expect(err).NotTo(HaveOccurred())
				sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
				sampleInventory := &inventoryv1alpha1.Inventory{}
				Expect(sampleYAML.Decode(sampleInventory)).NotTo(HaveOccurred())
				Expect(k8sClient.Create(ctx, sampleInventory)).NotTo(HaveOccurred())
			}

			switchesSamples := make([]string, 0)
			Expect(filepath.Walk("../samples/switches", func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					switchesSamples = append(switchesSamples, path)
				}
				return nil
			})).NotTo(HaveOccurred())
			for _, samplePath := range switchesSamples {
				stype := "leaf"
				sampleBytes, err := ioutil.ReadFile(samplePath)
				Expect(err).NotTo(HaveOccurred())
				sampleYAML := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
				sampleSwitch := &switchv1beta1.Switch{}
				Expect(sampleYAML.Decode(sampleSwitch)).NotTo(HaveOccurred())
				if strings.Contains(sampleSwitch.Name, "spine") {
					stype = "spine"
				}
				sampleSwitch.Labels = map[string]string{
					"metalapi.onmetal.de/inventoried":   "true",
					"metalapi.onmetal.de/inventory-ref": sampleSwitch.Spec.UUID,
					"metalapi.onmetal.de/type":          stype,
				}
				Expect(k8sClient.Create(ctx, sampleSwitch)).NotTo(HaveOccurred())
			}
		})

		It("Should set switches' loopback and change them when loopback IP objects change", func() {
			By("Base configuration processing")
			var switchesList = &switchv1beta1.SwitchList{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.List(ctx, switchesList)).NotTo(HaveOccurred())
				for _, item := range switchesList.Items {
					g.Expect(item.Status.SwitchState).NotTo(BeNil())
					g.Expect(switchv1beta1.GoString(item.Status.SwitchState.State)).To(Equal(switchv1beta1.CSwitchStateReady))
					g.Expect(item.ConnectionsOK(switchesList)).Should(BeTrue())
				}
			}, timeout, interval).Should(Succeed())

			By("Networking configuration done")
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.List(ctx, switchesList)).Should(Succeed())
				for _, item := range switchesList.Items {
					g.Expect(len(item.Status.LoopbackAddresses)).ShouldNot(BeZero())
					g.Expect(len(item.Status.LoopbackAddresses)).Should(Equal(2))
					g.Expect(len(item.Status.Subnets)).ShouldNot(BeZero())
					g.Expect(len(item.Status.Subnets)).Should(Equal(2))
				}
			}, timeout, interval).Should(Succeed())

			By("Changing loopbacks' IP objects should cause changes in switch configuration")
			loopbacks := &ipamv1alpha1.IPList{}
			Expect(k8sClient.List(ctx, loopbacks)).Should(Succeed())
			for _, lo := range loopbacks.Items {
				if lo.Spec.Subnet.Name == "loopbacks-v6" {
					delete(lo.Labels, switchv1beta1.IPAMObjectOwnerLabel)
					Expect(k8sClient.Update(ctx, &lo))
				}
			}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.List(ctx, switchesList)).Should(Succeed())
				for _, item := range switchesList.Items {
					g.Expect(len(item.Status.LoopbackAddresses)).ShouldNot(BeZero())
					g.Expect(len(item.Status.LoopbackAddresses)).Should(Equal(1))
				}
			}, timeout, interval).Should(Succeed())

			By("Changing subnets objects should cause changes in switch configuration")
			subnets := &ipamv1alpha1.SubnetList{}
			Expect(k8sClient.List(ctx, subnets)).Should(Succeed())
			for _, sn := range subnets.Items {
				if sn.Spec.ParentSubnet.Name == "switch-ranges-v6" {
					delete(sn.Labels, switchv1beta1.IPAMObjectOwnerLabel)
					Expect(k8sClient.Update(ctx, &sn))
				}
			}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.List(ctx, switchesList)).Should(Succeed())
				for _, item := range switchesList.Items {
					g.Expect(len(item.Status.Subnets)).ShouldNot(BeZero())
					g.Expect(len(item.Status.Subnets)).Should(Equal(1))
				}
			}, timeout, interval).Should(Succeed())
		})
	})
})
