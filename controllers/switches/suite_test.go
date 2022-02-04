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
	"bytes"
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

	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	ipamctrl "github.com/onmetal/ipam/controllers"
	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/mod/modfile"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	switchv1alpha1 "github.com/onmetal/metal-api/apis/switches/v1alpha1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

const (
	OnmetalNamespace = "onmetal"
	timeout          = time.Second * 45
	interval         = time.Millisecond * 250
)

var k8sClient client.Client
var ctx context.Context
var testEnv *envtest.Environment
var cancel context.CancelFunc
var ipv4Used, ipv6Used = false, false

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
	subnetGlobalCrdPath := getCrdPath(ipamv1alpha1.Subnet{})
	switchGlobalCrdPath := filepath.Join("..", "..", "config", "crd", "bases")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{switchGlobalCrdPath, subnetGlobalCrdPath},
		ErrorIfCRDPathMissing: true,
		WebhookInstallOptions: envtest.WebhookInstallOptions{
			Paths: []string{filepath.Join("..", "..", "config", "webhook")},
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
	err = ipamv1alpha1.AddToScheme(globalScheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: globalScheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	By("Set up k8s-subnet manager")
	ipamScheme := scheme.Scheme
	err = ipamv1alpha1.AddToScheme(ipamScheme)
	Expect(err).ToNot(HaveOccurred())
	k8sSubnetManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             ipamScheme,
		LeaderElection:     false,
		MetricsBindAddress: "0",
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&ipamctrl.SubnetReconciler{
		Client: k8sSubnetManager.GetClient(),
		Scheme: k8sSubnetManager.GetScheme(),
		Log:    ctrl.Log.WithName("k8s-subnet").WithName("subnet"),
	}).SetupWithManager(k8sSubnetManager)
	Expect(err).ToNot(HaveOccurred())
	err = (&ipamctrl.NetworkReconciler{
		Client: k8sSubnetManager.GetClient(),
		Scheme: k8sSubnetManager.GetScheme(),
		Log:    ctrl.Log.WithName("k8s-subnet").WithName("network"),
	}).SetupWithManager(k8sSubnetManager)
	Expect(err).ToNot(HaveOccurred())
	err = (&ipamctrl.IPReconciler{
		Client: k8sSubnetManager.GetClient(),
		Scheme: k8sSubnetManager.GetScheme(),
		Log:    ctrl.Log.WithName("k8s-subnet").WithName("ip"),
	}).SetupWithManager(k8sSubnetManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sSubnetManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred())
	}()

	webhookInstallOptions := &testEnv.WebhookInstallOptions
	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             globalScheme,
		Host:               webhookInstallOptions.LocalServingHost,
		Port:               webhookInstallOptions.LocalServingPort,
		CertDir:            webhookInstallOptions.LocalServingCertDir,
		LeaderElection:     false,
		MetricsBindAddress: "0",
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&switchv1alpha1.SwitchAssignment{}).SetupWebhookWithManager(k8sManager)
	Expect(err).NotTo(HaveOccurred())
	err = (&switchv1alpha1.Switch{}).SetupWebhookWithManager(k8sManager)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:webhook

	err = (&SwitchReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("Switch"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())
	err = (&SwitchAssignmentReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("SwitchAssignment"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())
	err = (&InventoryReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("Inventory"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred())
	}()

	dialer := &net.Dialer{Timeout: time.Second}
	addrPort := fmt.Sprintf("%s:%d", webhookInstallOptions.LocalServingHost, webhookInstallOptions.LocalServingPort)
	Eventually(func() error {
		conn, err := tls.DialWithDialer(dialer, "tcp", addrPort, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return err
		}
		_ = conn.Close()
		return nil
	}).Should(Succeed())

	namespace := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: switchv1alpha1.CNamespace}}
	Expect(k8sClient.Create(ctx, namespace)).To(Succeed())
}, 60)

var _ = AfterSuite(func() {
	cleanUp()
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

func getCrdPath(crdPackageScheme interface{}) string {
	globalPackagePath := reflect.TypeOf(crdPackageScheme).PkgPath()
	goModData, err := ioutil.ReadFile(filepath.Join("..", "..", "go.mod"))
	Expect(err).NotTo(HaveOccurred())
	goModFile, err := modfile.Parse("", goModData, nil)
	Expect(err).NotTo(HaveOccurred())
	globalModulePath := ""
	for _, req := range goModFile.Require {
		if strings.HasPrefix(globalPackagePath, req.Mod.Path) {
			globalModulePath = req.Mod.String()
			break
		}
	}
	Expect(globalModulePath).NotTo(BeZero())
	return filepath.Join(build.Default.GOPATH, "pkg", "mod", globalModulePath, "config", "crd", "bases")
}

func prepareEnv() {
	By("Prepare inventories")
	switchesSamples := []string{
		filepath.Join("..", "..", "config", "samples", "testdata", "edge-leaf-1.fra3.infra.onmetal.de.yaml"),
		filepath.Join("..", "..", "config", "samples", "testdata", "edge-leaf-2.fra3.infra.onmetal.de.yaml"),
		filepath.Join("..", "..", "config", "samples", "testdata", "leaf-1.fra3.infra.onmetal.de.yaml"),
		filepath.Join("..", "..", "config", "samples", "testdata", "leaf-2.fra3.infra.onmetal.de.yaml"),
		filepath.Join("..", "..", "config", "samples", "testdata", "spine-1.fra3.infra.onmetal.de.yaml"),
		filepath.Join("..", "..", "config", "samples", "testdata", "spine-2.fra3.infra.onmetal.de.yaml"),
	}
	for _, sample := range switchesSamples {
		inv := &inventoriesv1alpha1.Inventory{}
		sampleBytes, err := ioutil.ReadFile(sample)
		Expect(err).NotTo(HaveOccurred())
		dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
		Expect(dec.Decode(inv)).NotTo(HaveOccurred())
		Expect(k8sClient.Create(ctx, inv)).To(Succeed())
	}

	By("Prepare assignments")
	assignmentsSamples := []string{
		filepath.Join("..", "..", "config", "samples", "testdata", "assignment-sp1.yaml"),
		filepath.Join("..", "..", "config", "samples", "testdata", "assignment-sp2.yaml"),
	}
	for _, sample := range assignmentsSamples {
		swa := &switchv1alpha1.SwitchAssignment{}
		sampleBytes, err := ioutil.ReadFile(sample)
		Expect(err).NotTo(HaveOccurred())
		dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
		Expect(dec.Decode(swa)).NotTo(HaveOccurred())
		Expect(k8sClient.Create(ctx, swa)).To(Succeed())
		Eventually(func() bool {
			err := k8sClient.Get(ctx, swa.NamespacedName(), swa)
			return err == nil
		}, timeout, interval).Should(BeTrue())
		Expect(swa.Labels).Should(Equal(map[string]string{switchv1alpha1.LabelChassisId: switchv1alpha1.MacToLabel(swa.Spec.ChassisID)}))
	}

	By("Switches created from inventories")
	list := &switchv1alpha1.SwitchList{}
	Expect(k8sClient.List(ctx, list)).Should(Succeed())
	Expect(len(list.Items)).To(Equal(len(switchesSamples)))
	for _, sw := range list.Items {
		Expect(sw.Labels).Should(Equal(map[string]string{switchv1alpha1.LabelChassisId: switchv1alpha1.MacToLabel(sw.Spec.Chassis.ChassisID)}))
	}
}

func prepareNetwork(subnetsSamples []string) {
	By("Prepare network")
	networkSample := filepath.Join("..", "..", "config", "samples", "testdata", "underlay-network.yaml")
	network := &ipamv1alpha1.Network{}
	sampleBytes, err := ioutil.ReadFile(networkSample)
	Expect(err).NotTo(HaveOccurred())
	dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
	Expect(dec.Decode(network)).NotTo(HaveOccurred())
	Expect(k8sClient.Create(ctx, network)).To(Succeed())

	By("Prepare subnets")
	for _, sample := range subnetsSamples {
		subnet := &ipamv1alpha1.Subnet{}
		sampleBytes, err := ioutil.ReadFile(sample)
		Expect(err).NotTo(HaveOccurred())
		dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(sampleBytes), len(sampleBytes))
		Expect(dec.Decode(subnet)).NotTo(HaveOccurred())
		Expect(k8sClient.Create(ctx, subnet)).To(Succeed())
		Eventually(func() bool {
			Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: subnet.Namespace, Name: subnet.Name}, subnet)).To(Succeed())
			return subnet.Status.State == ipamv1alpha1.CFinishedSubnetState
		}, timeout, interval).Should(BeTrue())
		if subnet.Status.Type == ipamv1alpha1.CIPv4SubnetType {
			ipv4Used = true
		}
		if subnet.Status.Type == ipamv1alpha1.CIPv6SubnetType {
			ipv6Used = true
		}
	}
}

func cleanUp() {
	By("Remove switches")
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

	By("Check assignments in pending state")
	list := &switchv1alpha1.SwitchAssignmentList{}
	Eventually(func() bool {
		Expect(k8sClient.List(ctx, list)).Should(Succeed())
		for _, item := range list.Items {
			if item.Status.State != switchv1alpha1.CAssignmentStatePending {
				return false
			}
		}
		return true
	}, timeout, interval).Should(BeTrue())

	By("Remove assignments")
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

	By("Remove inventories")
	Expect(k8sClient.DeleteAllOf(ctx, &inventoriesv1alpha1.Inventory{}, client.InNamespace(OnmetalNamespace))).To(Succeed())
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

	By("Remove IPs")
	Expect(k8sClient.DeleteAllOf(ctx, &ipamv1alpha1.IP{}, client.InNamespace(OnmetalNamespace))).To(Succeed())
	Eventually(func() bool {
		list := &ipamv1alpha1.IPList{}
		err := k8sClient.List(ctx, list)
		if err != nil {
			return false
		}
		if len(list.Items) > 0 {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())

	By("Remove subnets")
	Expect(k8sClient.DeleteAllOf(ctx, &ipamv1alpha1.Subnet{}, client.InNamespace(OnmetalNamespace))).To(Succeed())
	Eventually(func() bool {
		list := &ipamv1alpha1.SubnetList{}
		err := k8sClient.List(ctx, list)
		if err != nil {
			return false
		}
		if len(list.Items) > 0 {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())

	By("Remove networks")
	Expect(k8sClient.DeleteAllOf(ctx, &ipamv1alpha1.Network{}, client.InNamespace(OnmetalNamespace))).To(Succeed())
	Eventually(func() bool {
		list := &ipamv1alpha1.NetworkList{}
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
