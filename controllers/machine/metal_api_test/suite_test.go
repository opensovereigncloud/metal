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
	"path/filepath"
	"testing"

	machinev1alpha3 "github.com/onmetal/metal-api/apis/machine/v1alpha3"
	controllers "github.com/onmetal/metal-api/controllers/machine"

	benchv1alpha3 "github.com/onmetal/metal-api/apis/benchmark/v1alpha3"
	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	oobonmetal "github.com/onmetal/oob-operator/api/v1alpha1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
	ctx       context.Context
	cancel    context.CancelFunc
)

var scheme = runtime.NewScheme()

// nolint
func TestMachineController(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Machine Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	By("bootstrapping test environment")

	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "config", "crd", "additional"),
			filepath.Join("..", "..", "..", "config", "crd", "bases"),
		},
		ErrorIfCRDPathMissing: true,
	}
	ctx, cancel = context.WithCancel(context.TODO())

	oobonmetal.SchemeBuilder.Register(&oobonmetal.OOB{})
	inventoriesv1alpha1.SchemeBuilder.Register(&inventoriesv1alpha1.Inventory{}, &inventoriesv1alpha1.InventoryList{})
	switchv1beta1.SchemeBuilder.Register(&switchv1beta1.Switch{}, &switchv1beta1.SwitchList{})
	benchv1alpha3.SchemeBuilder.Register(&benchv1alpha3.Machine{}, &benchv1alpha3.MachineList{})
	machinev1alpha3.SchemeBuilder.Register(&machinev1alpha3.Machine{}, &machinev1alpha3.MachineList{})

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	Expect(inventoriesv1alpha1.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(benchv1alpha3.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(switchv1beta1.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(oobonmetal.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(corev1.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(machinev1alpha3.AddToScheme(scheme)).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{Scheme: scheme, MetricsBindAddress: "0"})
	Expect(err).ToNot(HaveOccurred())

	err = (&controllers.MachinePowerReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("machine-power"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()
})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
