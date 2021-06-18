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

package tests

import (
	"context"
	"crypto/tls"
	"fmt"
	"go/build"
	"io/ioutil"
	"net"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	inventory "github.com/onmetal/k8s-inventory/controllers"
	networkglobalv1alpha1 "github.com/onmetal/k8s-network-global/api/v1alpha1"
	networkglobal "github.com/onmetal/k8s-network-global/controllers"
	subnetv1alpha1 "github.com/onmetal/k8s-subnet/api/v1alpha1"
	subnet "github.com/onmetal/k8s-subnet/controllers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/mod/modfile"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
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
	"github.com/onmetal/switch-operator/controllers"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var ctx context.Context
var testEnv *envtest.Environment
var cancel context.CancelFunc

const (
	DefaultNamespace     = "default"
	OnmetalNamespace     = "onmetal"
	timeout              = time.Second * 30
	interval             = time.Millisecond * 250
	UnderlayNetwork      = "underlay"
	SubnetName           = "switch-subnet"
	SubnetCIDR           = "100.64.0.0/12"
	TestRegion           = "EU-West"
	TestAvailabilityZone = "A"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	ctx, cancel = context.WithCancel(context.TODO())
	var crdPaths []string
	var webhooksPaths []string

	By("Set up network global CRD paths")
	networkGlobalPackagePath := reflect.TypeOf(networkglobalv1alpha1.NetworkGlobal{}).PkgPath()
	goModData, err := ioutil.ReadFile(filepath.Join("..", "go.mod"))
	Expect(err).NotTo(HaveOccurred())
	goModFile, err := modfile.Parse("", goModData, nil)
	Expect(err).NotTo(HaveOccurred())
	networkGlobalModulePath := ""
	for _, req := range goModFile.Require {
		if strings.HasPrefix(networkGlobalPackagePath, req.Mod.Path) {
			networkGlobalModulePath = req.Mod.String()
			break
		}
	}
	Expect(networkGlobalModulePath).NotTo(BeZero())
	networkGlobalCrdPath := filepath.Join(build.Default.GOPATH, "pkg", "mod", networkGlobalModulePath, "config", "crd", "bases")
	crdPaths = append(crdPaths, networkGlobalCrdPath)

	By("Set up subnet CRD paths")
	subnetGlobalPackagePath := reflect.TypeOf(subnetv1alpha1.Subnet{}).PkgPath()
	goModData, err = ioutil.ReadFile(filepath.Join("..", "go.mod"))
	Expect(err).NotTo(HaveOccurred())
	goModFile, err = modfile.Parse("", goModData, nil)
	Expect(err).NotTo(HaveOccurred())
	subnetGlobalModulePath := ""
	for _, req := range goModFile.Require {
		if strings.HasPrefix(subnetGlobalPackagePath, req.Mod.Path) {
			subnetGlobalModulePath = req.Mod.String()
			break
		}
	}
	Expect(networkGlobalModulePath).NotTo(BeZero())
	subnetGlobalCrdPath := filepath.Join(build.Default.GOPATH, "pkg", "mod", subnetGlobalModulePath, "config", "crd", "bases")
	crdPaths = append(crdPaths, subnetGlobalCrdPath)

	By("Set up inventory CRD paths")
	inventoryGlobalPackagePath := reflect.TypeOf(inventoriesv1alpha1.Inventory{}).PkgPath()
	goModData, err = ioutil.ReadFile(filepath.Join("..", "go.mod"))
	Expect(err).NotTo(HaveOccurred())
	goModFile, err = modfile.Parse("", goModData, nil)
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
	crdPaths = append(crdPaths, inventoryGlobalCrdPath)

	By("Set up switch CRD paths")
	switchGlobalCrdPath := filepath.Join("..", "config", "crd", "bases")
	switchGlobalWebhookPath := filepath.Join("..", "config", "webhook")
	crdPaths = append(crdPaths, switchGlobalCrdPath)
	webhooksPaths = append(webhooksPaths, switchGlobalWebhookPath)

	By("Bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     crdPaths,
		ErrorIfCRDPathMissing: true,
		WebhookInstallOptions: envtest.WebhookInstallOptions{
			Paths: webhooksPaths,
		},
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	globalScheme := scheme.Scheme
	err = switchv1alpha1.AddToScheme(globalScheme)
	Expect(err).NotTo(HaveOccurred())
	err = inventoriesv1alpha1.AddToScheme(globalScheme)
	Expect(err).NotTo(HaveOccurred())
	err = networkglobalv1alpha1.AddToScheme(globalScheme)
	Expect(err).NotTo(HaveOccurred())
	err = subnetv1alpha1.AddToScheme(globalScheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: globalScheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	By("Set up k8s-inventory manager")
	inventoryScheme := scheme.Scheme
	err = inventoriesv1alpha1.AddToScheme(inventoryScheme)
	Expect(err).NotTo(HaveOccurred())
	k8sInventoryManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             inventoryScheme,
		LeaderElection:     false,
		MetricsBindAddress: "0",
	})
	Expect(err).ToNot(HaveOccurred())
	err = (&inventory.InventoryReconciler{
		Client: k8sInventoryManager.GetClient(),
		Scheme: k8sInventoryManager.GetScheme(),
		Log:    ctrl.Log.WithName("k8s-inventory").WithName("inventory"),
	}).SetupWithManager(k8sInventoryManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sInventoryManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred())
	}()

	By("Set up k8s-network-global manager")
	networkGlobalScheme := scheme.Scheme
	err = networkglobalv1alpha1.AddToScheme(networkGlobalScheme)
	Expect(err).ToNot(HaveOccurred())
	k8sNetGlobalManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             networkGlobalScheme,
		LeaderElection:     false,
		MetricsBindAddress: "0",
	})
	Expect(err).ToNot(HaveOccurred())
	err = (&networkglobal.NetworkGlobalReconciler{
		Client: k8sNetGlobalManager.GetClient(),
		Scheme: k8sNetGlobalManager.GetScheme(),
		Log:    ctrl.Log.WithName("k8s-network-global").WithName("network-global"),
	}).SetupWithManager(k8sNetGlobalManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sNetGlobalManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred())
	}()

	By("Set up k8s-subnet manager")
	subnetScheme := scheme.Scheme
	err = subnetv1alpha1.AddToScheme(subnetScheme)
	Expect(err).ToNot(HaveOccurred())
	err = networkglobalv1alpha1.AddToScheme(subnetScheme)
	Expect(err).ToNot(HaveOccurred())

	k8sSubnetManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             subnetScheme,
		LeaderElection:     false,
		MetricsBindAddress: "0",
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&subnet.SubnetReconciler{
		Client: k8sSubnetManager.GetClient(),
		Scheme: k8sSubnetManager.GetScheme(),
		Log:    ctrl.Log.WithName("k8s-subnet").WithName("subnet"),
	}).SetupWithManager(k8sSubnetManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sSubnetManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred())
	}()

	By("Set up switch-operator manager")
	switchScheme := scheme.Scheme
	err = admissionv1beta1.AddToScheme(switchScheme)
	Expect(err).NotTo(HaveOccurred())
	err = switchv1alpha1.AddToScheme(switchScheme)
	Expect(err).NotTo(HaveOccurred())
	err = inventoriesv1alpha1.AddToScheme(switchScheme)
	Expect(err).NotTo(HaveOccurred())
	err = subnetv1alpha1.AddToScheme(switchScheme)
	Expect(err).NotTo(HaveOccurred())

	switchWebhookInstallOptions := &testEnv.WebhookInstallOptions
	k8sSwitchManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             switchScheme,
		Host:               switchWebhookInstallOptions.LocalServingHost,
		Port:               switchWebhookInstallOptions.LocalServingPort,
		CertDir:            switchWebhookInstallOptions.LocalServingCertDir,
		LeaderElection:     false,
		MetricsBindAddress: "0",
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&switchv1alpha1.SwitchAssignment{}).SetupWebhookWithManager(k8sSwitchManager)
	Expect(err).NotTo(HaveOccurred())
	err = (&switchv1alpha1.Switch{}).SetupWebhookWithManager(k8sSwitchManager)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:webhook

	err = (&controllers.SwitchReconciler{
		Client: k8sSwitchManager.GetClient(),
		Scheme: k8sSwitchManager.GetScheme(),
		Log:    ctrl.Log.WithName("switch-operator").WithName("switch"),
	}).SetupWithManager(k8sSwitchManager)
	Expect(err).NotTo(HaveOccurred())
	err = (&controllers.SwitchAssignmentReconciler{
		Client: k8sSwitchManager.GetClient(),
		Scheme: k8sSwitchManager.GetScheme(),
		Log:    ctrl.Log.WithName("switch-operator").WithName("switch-assignment"),
	}).SetupWithManager(k8sSwitchManager)
	Expect(err).NotTo(HaveOccurred())
	err = (&controllers.InventoryReconciler{
		Client: k8sSwitchManager.GetClient(),
		Scheme: k8sSwitchManager.GetScheme(),
		Log:    ctrl.Log.WithName("switch-operator").WithName("inventory"),
	}).SetupWithManager(k8sSwitchManager)
	Expect(err).NotTo(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sSwitchManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred())
	}()

	By("Starting webhook servers")
	switchDialer := &net.Dialer{Timeout: time.Second}
	switchAddrPort := fmt.Sprintf("%s:%d", switchWebhookInstallOptions.LocalServingHost, switchWebhookInstallOptions.LocalServingPort)
	Eventually(func() error {
		conn, err := tls.DialWithDialer(switchDialer, "tcp", switchAddrPort, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return err
		}
		conn.Close()
		return nil
	}).Should(Succeed())

	namespace := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: switchv1alpha1.CNamespace}}
	Expect(k8sClient.Create(ctx, namespace)).To(Succeed())

}, 60)

var _ = AfterSuite(func() {
	By("Remove test resources")
	Expect(k8sClient.DeleteAllOf(ctx, &inventoriesv1alpha1.Inventory{}, client.InNamespace(DefaultNamespace))).To(Succeed())
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
	Expect(k8sClient.DeleteAllOf(ctx, &switchv1alpha1.Switch{}, client.InNamespace(OnmetalNamespace))).To(Succeed())
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
	Expect(k8sClient.DeleteAllOf(ctx, &switchv1alpha1.SwitchAssignment{}, client.InNamespace(OnmetalNamespace))).To(Succeed())
	Eventually(func() bool {
		list := &switchv1alpha1.SwitchAssignmentList{}
		err := k8sClient.List(ctx, list)
		if err != nil {
			return false
		}
		if len(list.Items) > 0 {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())
	Expect(k8sClient.DeleteAllOf(ctx, &subnetv1alpha1.Subnet{}, client.InNamespace(OnmetalNamespace))).To(Succeed())
	Eventually(func() bool {
		list := &subnetv1alpha1.SubnetList{}
		err := k8sClient.List(ctx, list)
		if err != nil {
			return false
		}
		if len(list.Items) > 0 {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())
	Expect(k8sClient.DeleteAllOf(ctx, &networkglobalv1alpha1.NetworkGlobal{}, client.InNamespace(OnmetalNamespace))).To(Succeed())
	Eventually(func() bool {
		list := &networkglobalv1alpha1.NetworkGlobalList{}
		err := k8sClient.List(ctx, list)
		if err != nil {
			return false
		}
		if len(list.Items) > 0 {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())

	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
