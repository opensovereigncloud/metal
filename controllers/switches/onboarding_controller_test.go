/*
Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
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

package controllers

import (
	"context"

	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1alpha1 "github.com/onmetal/metal-api/apis/switches/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("machine-controller", func() {
	Context("Controller Test", func() {
		It("Should create switch", func() {
			testSwitch()
		})
	})
})

func testSwitch() {
	ctx := context.Background()
	inventory := prepareInventory()

	By("Expect successful inventory creation")
	err := k8sClient.Create(ctx, inventory)
	Expect(err).Should(BeNil())

	By("Expect successful switch creation")
	Eventually(func(g Gomega) error {
		m := &switchv1alpha1.Switch{}
		return k8sClient.Get(ctx, types.NamespacedName{Name: inventory.Name, Namespace: inventory.Namespace}, m)
	}, timeout, interval).Should(BeNil())

}

func prepareInventory() *inventoriesv1alpha1.Inventory {
	var (
		name      = "a967954c-3475-2b82-a85c-84d8b4f8cd2d"
		namespace = switchv1alpha1.CNamespace
	)

	return &inventoriesv1alpha1.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: prepareSpecForInventory(name),
	}
}

func prepareSpecForInventory(name string) inventoriesv1alpha1.InventorySpec {
	return inventoriesv1alpha1.InventorySpec{
		System: &inventoriesv1alpha1.SystemSpec{
			ID: name,
		},
		Host: &inventoriesv1alpha1.HostSpec{
			Type: "Switch",
			Name: "node1",
		},
		Distro: &inventoriesv1alpha1.DistroSpec{
			AsicType: "test",
			CommitId: "123",
		},
		Benchmark: &inventoriesv1alpha1.BenchmarkSpec{
			Blocks:   []inventoriesv1alpha1.BlockBenchmarkCollection{},
			Networks: []inventoriesv1alpha1.NetworkBenchmarkResult{},
		},
		Blocks: []inventoriesv1alpha1.BlockSpec{},
		Memory: &inventoriesv1alpha1.MemorySpec{},
		CPUs:   []inventoriesv1alpha1.CPUSpec{},
		NICs: []inventoriesv1alpha1.NICSpec{
			{
				Name: "enp0s31f6",
				LLDPs: []inventoriesv1alpha1.LLDPSpec{
					{
						ChassisID:         "3c:2c:99:9d:cd:48",
						SystemName:        "EC1817001226",
						SystemDescription: "ECS4100-52T",
						PortID:            "3c:2c:99:9d:cd:77",
						PortDescription:   "Ethernet MachinePort on unit 1, port 47",
					},
				},
			},
		},
	}
}
