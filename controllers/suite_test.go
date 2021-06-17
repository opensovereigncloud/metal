/*
Copyright 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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
	"go/build"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/mod/modfile"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var ctx context.Context
var testEnv *envtest.Environment
var cancel context.CancelFunc

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	inventoryGlobalPackagePath := reflect.TypeOf(inventoriesv1alpha1.Inventory{}).PkgPath()

	goModData, err := ioutil.ReadFile(filepath.Join("..", "go.mod"))
	Expect(err).NotTo(HaveOccurred())

	goModFile, err := modfile.Parse("", goModData, nil)
	Expect(err).NotTo(HaveOccurred())

	inventoryGlobalModulePath := ""
	for _, req := range goModFile.Require {
		if strings.HasPrefix(inventoryGlobalPackagePath, req.Mod.Path) {
			inventoryGlobalModulePath = req.Mod.String()
			break
		}
	}
	Expect(inventoryGlobalModulePath).NotTo(BeZero())

	inventoryGlobalCrdPath := filepath.Join(build.Default.GOPATH, "pkg", "mod", inventoryGlobalModulePath, "config", "crd", "bases")
	switchGlobalCrdPath := filepath.Join("..", "config", "crd", "bases")

	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{switchGlobalCrdPath, inventoryGlobalCrdPath},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = switchv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	err = inventoriesv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme
	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&SwitchReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Switch"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&SwitchAssignmentReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("SwitchAssignment"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&InventoryReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Inventory"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred())
	}()

	k8sClient = k8sManager.GetClient()
	Expect(k8sClient).NotTo(BeNil())

	namespace := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: switchv1alpha1.CNamespace}}
	Expect(k8sClient.Create(ctx, namespace)).To(Succeed())

}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
