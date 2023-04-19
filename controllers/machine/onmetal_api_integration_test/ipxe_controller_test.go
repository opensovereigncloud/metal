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
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("machine-power-controller", func() {
	//ctx := testing.SetupContext()
	//ns := SetupTest(ctx, ipxeReconcilers)
	//
	//It("should create ipxe configmap", func() {
	//	machinePool := &computev1alpha1.MachinePool{}
	//	metalMachine := &machinev1alpha2.Machine{}
	//	computeMachine := &computev1alpha1.Machine{}
	//	defaultIpxeCM := &corev1.ConfigMap{}
	//
	//	u, err := uuid.NewUUID()
	//	Expect(err).ToNot(HaveOccurred())
	//	var (
	//		name      = u.String()
	//		namespace = ns.Name
	//	)
	//
	//	// prepare test data
	//	By("Create healthy running machine")
	//	createHealthyRunningMachine(ctx, name, namespace, metalMachine)
	//
	//	By("Create compute machine with machine class")
	//	createComputeMachine(ctx, name, namespace, computeMachine)
	//
	//	By("Create machine pool")
	//	createMachinePool(ctx, name, namespace, machinePool)
	//
	//	By("adding template config")
	//	defaultMC := &corev1.ConfigMap{
	//		ObjectMeta: metav1.ObjectMeta{
	//			Name:      controllers.IpxeDefaultTemplateName,
	//			Namespace: namespace,
	//		},
	//		Data: map[string]string{
	//			"name": "|\n    #!ipxe\n\n    chain --replace --autofree http://host/path/boot.ipxe",
	//		},
	//	}
	//	Expect(k8sClient.Create(ctx, defaultMC)).To(Succeed())
	//
	//	By("Expect default CM was created")
	//	Eventually(func() bool {
	//		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: controllers.IpxeDefaultTemplateName}, defaultMC); err != nil {
	//			return false
	//		}
	//
	//		return true
	//	}).Should(BeTrue())
	//
	//	// testing
	//	By("Expect metal machine has no reservation")
	//	Eventually(func() bool {
	//		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, metalMachine); err != nil {
	//			return false
	//		}
	//
	//		return metalMachine.Status.Reservation.Reference == nil
	//	}).Should(BeTrue())
	//
	//	By("Expect compute machine has no machine pool ref")
	//	Eventually(func() bool {
	//		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, computeMachine); err != nil {
	//			return false
	//		}
	//
	//		return computeMachine.Spec.MachinePoolRef == nil
	//	}).Should(BeTrue())
	//
	//	By("Expect successful computeMachine.Spec.machinePoolRef update")
	//	computeMachine.Spec.MachinePoolRef = &corev1.LocalObjectReference{Name: machinePool.Name}
	//	Expect(k8sClient.Update(ctx, computeMachine)).To(Succeed())
	//
	//	By("Expect compute machine has machine pool ref")
	//	Eventually(func() bool {
	//		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, computeMachine); err != nil {
	//			return false
	//		}
	//
	//		return computeMachine.Spec.MachinePoolRef != nil
	//	}).Should(BeTrue())
	//
	//	By("Waiting for default ipxe configmap to be created")
	//	Eventually(func() bool {
	//		if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: "ipxe-" + metalMachine.Name}, defaultIpxeCM); err != nil {
	//			return false
	//		}
	//
	//		return true
	//	}).Should(BeTrue())
	//})
})