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
	"time"

	benchv1alpha3 "github.com/onmetal/metal-api/apis/benchmark/v1alpha3"
	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"

	. "github.com/onsi/ginkgo"
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

const (
	timeout  = time.Second * 60
	interval = time.Millisecond * 250
)

var scheme = runtime.NewScheme()

func TestMachineController(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Benchmark Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	By("bootstrapping test environment")

	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "config", "crd", "bases"),
		},
		ErrorIfCRDPathMissing: true,
	}
	ctx, cancel = context.WithCancel(context.TODO())

	inventoriesv1alpha1.SchemeBuilder.Register(&inventoriesv1alpha1.Inventory{}, &inventoriesv1alpha1.InventoryList{})
	benchv1alpha3.SchemeBuilder.Register(&benchv1alpha3.Machine{}, &benchv1alpha3.MachineList{})

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	Expect(inventoriesv1alpha1.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(benchv1alpha3.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(corev1.AddToScheme(scheme)).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{Scheme: scheme, Host: "127.0.0.1", MetricsBindAddress: "0"})
	Expect(err).ToNot(HaveOccurred())

	err = (&OnboardingReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("benchmark-onboarding"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()
}, 60)

var _ = AfterSuite(func() {
	cleanUp()
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

func cleanUp() {
	By("Remove benchmarks")
	Expect(k8sClient.DeleteAllOf(ctx, &benchv1alpha3.Machine{}, client.InNamespace("default"))).To(Succeed())
	Eventually(func() bool {
		list := &benchv1alpha3.MachineList{}
		err := k8sClient.List(ctx, list)
		if err != nil {
			return false
		}
		if len(list.Items) > 0 {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())

	By("Remove invetories")
	Expect(k8sClient.DeleteAllOf(ctx, &inventoriesv1alpha1.Inventory{}, client.InNamespace("default"))).To(Succeed())
	Eventually(func() bool {
		list := &inventoriesv1alpha1.InventoryList{}
		err := k8sClient.List(ctx, list)
		if err != nil {
			return false
		}
		if len(list.Items) > 0 {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())
}
