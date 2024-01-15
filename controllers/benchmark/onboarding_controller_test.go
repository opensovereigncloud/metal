// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package controllers

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
)

var _ = Describe("machine-controller", func() {
	var (
		name      = "a969952c-3475-2b82-a85c-94d8b4f8cd2d"
		namespace = "default"
	)

	Context("Controller Test", func() {
		It("Test benchmark onboarding ", func() {
			testBenchmarkOnboarding(name, namespace)
		})
	})
})

func testBenchmarkOnboarding(name, namespace string) {
	ctx := context.Background()

	inventory := prepareInventory(name, namespace)

	By("Expect successful inventory creation")
	Expect(k8sClient.Create(ctx, inventory)).Should(BeNil())

	By("Expect successful benchmark creation")
	b := &metalv1alpha4.Benchmark{}
	Eventually(func(g Gomega) error {
		return k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, b)
	}, timeout, interval).Should(Succeed())
}

func prepareInventory(name, namespace string) *metalv1alpha4.Inventory {
	return &metalv1alpha4.Inventory{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec:       prepareSpecForInventory(name),
	}
}

func prepareSpecForInventory(name string) metalv1alpha4.InventorySpec {
	return metalv1alpha4.InventorySpec{
		System: &metalv1alpha4.SystemSpec{
			ID: name,
		},
		Host: &metalv1alpha4.HostSpec{
			Name: "node1",
		},
		Blocks: []metalv1alpha4.BlockSpec{},
		Memory: &metalv1alpha4.MemorySpec{},
		CPUs:   []metalv1alpha4.CPUSpec{},
		NICs: []metalv1alpha4.NICSpec{
			{
				Name: "enp0s31f6",
				LLDPs: []metalv1alpha4.LLDPSpec{
					{
						ChassisID:         "3c:2c:99:9d:cd:48",
						SystemName:        "EC1817001226",
						SystemDescription: "ECS4100-52T",
						PortID:            "3c:2c:99:9d:cd:77",
						PortDescription:   "Ethernet100",
					},
				},
			},
			{
				Name: "enp1s32f6",
				LLDPs: []metalv1alpha4.LLDPSpec{
					{
						ChassisID:         "3c:2c:99:9d:cd:48",
						SystemName:        "EC1817001226",
						SystemDescription: "ECS4100-52T",
						PortID:            "3c:2c:99:9d:cd:77",
						PortDescription:   "Ethernet102",
					},
				},
			},
		},
	}
}
