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

	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	oobv1 "github.com/onmetal/oob-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("machine-controller", func() {
	Context("Controller Test", func() {
		It("Test machine onboarding ", func() {
			testMachineOOB()
		})
	})
})

func testMachineOOB() {
	var (
		name      = "a237952c-3475-2b82-a85c-84d8b4f8cd2d"
		namespace = "default"
	)
	ctx := context.Background()
	oob := prepareOOB(name, namespace)

	preparedMachine := prepareTestMachine(name, namespace)

	By("Expect successful test machine for oob creation")
	Expect(k8sClient.Create(ctx, preparedMachine)).To(Succeed())

	By("Expect successful oob creation")

	Expect(k8sClient.Create(ctx, oob)).To(Succeed())

	By("Expect successful oob status update")
	Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, oob)).To(Succeed())
	oob.Status = prepareOOBStatus(name)

	Expect(k8sClient.Status().Update(ctx, oob)).To(Succeed())

	By("Expect oob status to be true")

	m := &machinev1alpha2.Machine{}
	Eventually(func(g Gomega) bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, m); err != nil {
			return false
		}
		return m.Status.OOB.Exist
	}, timeout, interval).Should(BeTrue())

	By("Expect successful oob deletion")
	Expect(k8sClient.Delete(ctx, oob)).To(Succeed())

	By("Expect oob status to be false")
	Eventually(func(g Gomega) bool {
		if err := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, m); err != nil {
			return false
		}
		return m.Status.OOB.Exist
	}, timeout, interval).Should(BeFalse())
}

func prepareOOB(name, namespace string) *oobv1.OOB {
	return &oobv1.OOB{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: oobv1.OOBSpec{
			Power: "Off",
		},
	}
}

func prepareOOBStatus(name string) oobv1.OOBStatus {
	return oobv1.OOBStatus{
		UUID: name,
	}
}

func prepareTestMachine(name, namespace string) *machinev1alpha2.Machine {
	return &machinev1alpha2.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: machinev1alpha2.MachineSpec{},
	}
}
