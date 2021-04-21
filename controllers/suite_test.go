/*
Copyright 2021.

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
	"go/build"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/mod/modfile"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/kustomize/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/k8sdeps/transformer"
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/loader"
	"sigs.k8s.io/kustomize/pkg/resid"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/kustomize/pkg/target"

	inventoryv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"

	sizev1alpha1 "github.com/onmetal/k8s-size/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	sizePackagePath := reflect.TypeOf(sizev1alpha1.Size{}).PkgPath()

	goModData, err := ioutil.ReadFile(filepath.Join("..", "go.mod"))
	Expect(err).NotTo(HaveOccurred())

	goModFile, err := modfile.Parse("", goModData, nil)
	Expect(err).NotTo(HaveOccurred())

	sizeModulePath := ""
	for _, req := range goModFile.Require {
		if strings.HasPrefix(sizePackagePath, req.Mod.Path) {
			sizeModulePath = req.Mod.String()
		}
	}
	Expect(sizeModulePath).NotTo(BeZero())

	// https://github.com/kubernetes-sigs/kubebuilder/issues/1999
	sizeCrdPath := filepath.Join(build.Default.GOPATH, "pkg", "mod", sizeModulePath, "config", "crd")

	// Since kubebuilder is not allowing to set complex types for fields with markers
	// we have to build patched configuration with kustomize first
	unstructuredFactory := kunstruct.NewKunstructuredFactoryImpl()
	resourceFactory := resource.NewFactory(unstructuredFactory)
	resmapFactory := resmap.NewFactory(resourceFactory)
	transformerFactory := transformer.NewFactoryImpl()
	kfs := fs.MakeRealFS()

	loader, err := loader.NewLoader(sizeCrdPath, kfs)
	Expect(err).NotTo(HaveOccurred())

	kt, err := target.NewKustTarget(loader, resmapFactory, transformerFactory)
	Expect(err).NotTo(HaveOccurred())

	resMap, err := kt.MakeCustomizedResMap()
	Expect(err).NotTo(HaveOccurred())

	resIds := resMap.GetMatchingIds(func(id resid.ResId) bool {
		return id.Gvk().Kind == "CustomResourceDefinition" &&
			id.Gvk().Group == "apiextensions.k8s.io" &&
			id.Gvk().Version == "v1"
	})

	crds := make([]client.Object, 0)

	for _, resId := range resIds {
		res := resMap[resId]
		resJsonBytes, err := res.MarshalJSON()
		Expect(err).NotTo(HaveOccurred())

		crd := &v1.CustomResourceDefinition{}
		err = json.Unmarshal(resJsonBytes, crd)
		Expect(err).NotTo(HaveOccurred())

		crds = append(crds, crd)
	}

	testEnv = &envtest.Environment{
		CRDs:              crds,
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = sizev1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	err = inventoryv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme
	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&InventoryReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Size"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

	k8sClient = k8sManager.GetClient()
	Expect(k8sClient).NotTo(BeNil())
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
