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
	"github.com/google/uuid"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	controllers "github.com/onmetal/metal-api/controllers/machine"
	computev1alpha1 "github.com/onmetal/onmetal-api/api/compute/v1alpha1"
	"github.com/onmetal/onmetal-api/utils/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("machine-power-controller", func() {
	ctx := testing.SetupContext()
	ns := SetupTest(ctx, ipxeReconcilers)

	It("should create ipxe configmap", func() {
		machinePool := &computev1alpha1.MachinePool{}
		metalMachine := &machinev1alpha2.Machine{}
		computeMachine := &computev1alpha1.Machine{}
		ipxeCM := &corev1.ConfigMap{}

		u, err := uuid.NewUUID()
		Expect(err).ToNot(HaveOccurred())
		var (
			name      = u.String()
			namespace = ns.Name
		)

		// prepare test data
		By("Create healthy running machine")
		createHealthyRunningMachine(ctx, name, namespace, metalMachine)

		By("Create compute machine with machine class")
		createComputeMachine(ctx, name, namespace, computeMachine)

		By("Create machine pool")
		createMachinePool(ctx, name, namespace, machinePool)

		// testing
		By("Expect metal machine has no reservation")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, metalMachine); err != nil {
				return false
			}

			return metalMachine.Status.Reservation.Reference == nil
		}).Should(BeTrue())

		By("Expect compute machine has no machine pool ref")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, computeMachine); err != nil {
				return false
			}

			return computeMachine.Spec.MachinePoolRef == nil
		}).Should(BeTrue())

		By("Expect successful computeMachine.Spec.machinePoolRef update")
		computeMachine.Spec.MachinePoolRef = &corev1.LocalObjectReference{Name: machinePool.Name}
		Expect(k8sClient.Update(ctx, computeMachine)).To(Succeed())

		By("Expect compute machine has machine pool ref")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, computeMachine); err != nil {
				return false
			}

			return computeMachine.Spec.MachinePoolRef != nil
		}).Should(BeTrue())

		By("Expect ipxe configmap was created")
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: "ipxe-" + metalMachine.Name}, ipxeCM); err != nil {
				return false
			}

			return true
		}).Should(BeTrue())

		By("Expect CM data is valid")
		Eventually(func() bool { return ipxeCM.Data["name"] != "" }).Should(BeTrue())
	})
})

type OnmetalImageParserFake struct{}

func (f *OnmetalImageParserFake) GetDescription(url string) (controllers.ImageDescription, error) {
	return controllers.ImageDescription{
		KernelDigest:    "5d7ae0f21ba60f208393e18041f9513eb3d2z802d12d4cf54c7049d4a0f3bf99",
		InitRAMFsDigest: "f8b52a8593bs49561c93e7a352d0jdd7da56c27b0072bead3dd6e4fa43430158",
		RootFSDigest:    "87d1c93b8hc23f0c3add9d2e6570x2f9115bf6a1070bdd48d9a8ec75678037d5",
		CommandLine:     "root=LABEL=ROOT ro console=tty0 console=ttyS0,999999 earlyprintk=ttyS0,999999 consoleblank=999 cgroup_enable=cgroup_enable swapaccount=999 ignition.firstboot=999 ignition.platform.id=ignition.platform.id security=security",
	}, nil
}
