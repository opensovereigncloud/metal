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

package machine

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/onmetal/controller-utils/envtestutils"
	"github.com/onmetal/controller-utils/kustomizeutils"
	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	machinev1lpha1 "github.com/onmetal/metal-api/apis/machine/v1alpha1"
	oobonmetal "github.com/onmetal/oob-controller/api/v1"
	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
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
	timeout  = 3 * time.Second
	interval = 50 * time.Millisecond
)

func TestMachineController(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	By("bootstrapping test environment")

	additionalCRDs := &apiextensionsv1.CustomResourceDefinitionList{}
	Expect(kustomizeutils.RunKustomizeIntoList(".", scheme.Codecs.UniversalDeserializer(), additionalCRDs)).To(Succeed())

	testEnv = &envtest.Environment{
		CRDs: envtestutils.CRDPtrsFromCRDs(additionalCRDs.Items),
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "config", "crd", "bases"),
		},
		ErrorIfCRDPathMissing: true,
	}
	ctx, cancel = context.WithCancel(context.TODO())

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	Expect(inventoriesv1alpha1.AddToScheme(scheme.Scheme)).NotTo(HaveOccurred())
	Expect(switchv1alpha1.AddToScheme(scheme.Scheme)).NotTo(HaveOccurred())
	Expect(machinev1lpha1.AddToScheme(scheme.Scheme)).NotTo(HaveOccurred())
	Expect(oobonmetal.AddToScheme(scheme.Scheme)).NotTo(HaveOccurred())
	Expect(corev1.AddToScheme(scheme.Scheme)).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{})
	Expect(err).ToNot(HaveOccurred())

	err = (&MachineReconciler{
		Client:   k8sManager.GetClient(),
		Recorder: k8sManager.GetEventRecorderFor("metal-api"),
		Log:      ctrl.Log.WithName("controllers").WithName("machine"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&InventoryReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("inventory"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	namespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: switchv1alpha1.CNamespace}}
	Expect(k8sClient.Create(ctx, namespace)).To(BeNil())

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
	By("Remove oob")
	Expect(k8sClient.DeleteAllOf(ctx, &oobonmetal.Machine{}, client.InNamespace(switchv1alpha1.CNamespace))).To(Succeed())
	Eventually(func() bool {
		list := &oobonmetal.MachineList{}
		err := k8sClient.List(ctx, list)
		if err != nil {
			return false
		}
		if len(list.Items) > 0 {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())

	By("Remove switches")
	Expect(k8sClient.DeleteAllOf(ctx, &switchv1alpha1.Switch{}, client.InNamespace(switchv1alpha1.CNamespace))).To(Succeed())
	Eventually(func() bool {
		list := &switchv1alpha1.SwitchList{}
		err := k8sClient.List(ctx, list)
		if err != nil {
			return false
		}
		if len(list.Items) > 0 {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())

	By("Remove inventories")
	Expect(k8sClient.DeleteAllOf(ctx, &inventoriesv1alpha1.Inventory{}, client.InNamespace(switchv1alpha1.CNamespace))).To(Succeed())
	Eventually(func() bool {
		list := &switchv1alpha1.SwitchList{}
		err := k8sClient.List(ctx, list)
		if err != nil {
			return false
		}
		if len(list.Items) > 0 {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())

	By("Remove namespace")
	Expect(k8sClient.Delete(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: switchv1alpha1.CNamespace}})).To(Succeed())
}
