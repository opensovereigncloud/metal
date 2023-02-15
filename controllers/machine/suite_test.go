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
	"github.com/onmetal/controller-utils/buildutils"
	"github.com/onmetal/controller-utils/modutils"
	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	utilsenvtest "github.com/onmetal/onmetal-api/utils/envtest"
	"github.com/onmetal/onmetal-api/utils/envtest/apiserver"
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"path/filepath"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
	"testing"
	"time"

	benchv1alpha3 "github.com/onmetal/metal-api/apis/benchmark/v1alpha3"
	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	computev1alpha1 "github.com/onmetal/onmetal-api/api/compute/v1alpha1"
	networkingv1alpha1 "github.com/onmetal/onmetal-api/api/networking/v1alpha1"
	storagev1alpha1 "github.com/onmetal/onmetal-api/api/storage/v1alpha1"
	oobonmetal "github.com/onmetal/oob-operator/api/v1alpha1"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
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
	timeout  = time.Second * 15
	interval = time.Millisecond * 250
)

var scheme = runtime.NewScheme()

// nolint
func TestMachineController(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Machine Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	var err error
	By("bootstrapping test environment")
	onmetalApiPackagePath := reflect.TypeOf(computev1alpha1.MachinePool{}).PkgPath()

	goModData, err := ioutil.ReadFile(filepath.Join("..", "..", "go.mod"))
	Expect(err).NotTo(HaveOccurred())

	goModFile, err := modfile.Parse("", goModData, nil)
	Expect(err).NotTo(HaveOccurred())

	onmetalApiModulePath := findPackagePath(goModFile, onmetalApiPackagePath)
	Expect(onmetalApiModulePath).NotTo(Equal(""))

	//onmetalApiCrdPath := filepath.Join(build.Default.GOPATH, "pkg", "mod", onmetalApiModulePath, "config", "apiserver", "apiservice", "bases")

	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "config", "crd", "bases")},
	}
	testEnvExt := &utilsenvtest.EnvironmentExtensions{
		APIServiceDirectoryPaths: []string{
			modutils.Dir("github.com/onmetal/onmetal-api", "config", "apiserver", "apiservice", "bases"),
		},
		ErrorIfAPIServicePathIsMissing: true,
	}

	oobonmetal.SchemeBuilder.Register(&oobonmetal.OOB{})
	inventoriesv1alpha1.SchemeBuilder.Register(&inventoriesv1alpha1.Inventory{}, &inventoriesv1alpha1.InventoryList{})
	switchv1beta1.SchemeBuilder.Register(&switchv1beta1.Switch{}, &switchv1beta1.SwitchList{})
	benchv1alpha3.SchemeBuilder.Register(&benchv1alpha3.Machine{}, &benchv1alpha3.MachineList{})
	machinev1alpha2.SchemeBuilder.Register(&machinev1alpha2.Machine{}, &machinev1alpha2.MachineList{})

	Expect(inventoriesv1alpha1.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(benchv1alpha3.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(switchv1beta1.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(oobonmetal.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(corev1.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(machinev1alpha2.AddToScheme(scheme)).NotTo(HaveOccurred())

	Expect(computev1alpha1.AddToScheme(scheme)).To(Succeed())
	Expect(networkingv1alpha1.AddToScheme(scheme)).To(Succeed())
	Expect(ipamv1alpha1.AddToScheme(scheme)).To(Succeed())
	Expect(storagev1alpha1.AddToScheme(scheme)).To(Succeed())
	Expect(apiregistrationv1.AddToScheme(scheme)).To(Succeed())

	cfg, err = utilsenvtest.StartWithExtensions(testEnv, testEnvExt)
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())
	DeferCleanup(utilsenvtest.StopWithExtensions, testEnv, testEnvExt)

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	apiSrv, err := apiserver.New(cfg, apiserver.Options{
		MainPath:     "github.com/onmetal/onmetal-api/cmd/onmetal-apiserver",
		BuildOptions: []buildutils.BuildOption{buildutils.ModModeMod},
		ETCDServers:  []string{testEnv.ControlPlane.Etcd.URL.String()},
		Host:         testEnvExt.APIServiceInstallOptions.LocalServingHost,
		Port:         testEnvExt.APIServiceInstallOptions.LocalServingPort,
		CertDir:      testEnvExt.APIServiceInstallOptions.LocalServingCertDir,
		AttachOutput: true,
	})
	Expect(err).NotTo(HaveOccurred())

	ctx, cancel := context.WithCancel(context.Background())
	DeferCleanup(cancel)
	go func() {
		defer GinkgoRecover()
		err := apiSrv.Start()
		Expect(err).NotTo(HaveOccurred())
	}()

	err = utilsenvtest.WaitUntilAPIServicesReadyWithTimeout(5*time.Minute, testEnvExt, k8sClient, scheme)
	Expect(err).NotTo(HaveOccurred())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: "0",
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&InventoryReconciler{
		Client:    k8sManager.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("machine-inventory"),
		Recorder:  k8sManager.GetEventRecorderFor("inventory-controller"),
		Namespace: "default",
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&OOBReconciler{
		Client:    k8sManager.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("machine-oob"),
		Recorder:  k8sManager.GetEventRecorderFor("Machine-OOB"),
		Namespace: "default",
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&PoolReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("machine-pool"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		Expect(k8sManager.Start(ctx)).To(Succeed())
	}()
})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

func findPackagePath(modFile *modfile.File, packagePath string) string {
	for _, req := range modFile.Require {
		if strings.HasPrefix(packagePath, req.Mod.Path) {
			return req.Mod.String()
		}
	}
	return ""
}
