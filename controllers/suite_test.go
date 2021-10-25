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

	"github.com/google/uuid"
	subnetv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	subnet "github.com/onmetal/ipam/controllers"
	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	inventory "github.com/onmetal/k8s-inventory/controllers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/mod/modfile"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
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

const (
	DefaultNamespace     = "default"
	OnmetalNamespace     = "onmetal"
	timeout              = time.Second * 90
	interval             = time.Millisecond * 250
	UnderlayNetwork      = "underlay"
	SubnetNameV4         = "switch-ranges-v4"
	LoopbackSubnetV4     = "switches-v4"
	SubnetNameV6         = "switch-ranges-v6"
	LoopbackSubnetV6     = "switches-v6"
	SubnetIPv4CIDR       = "100.64.0.0/16"
	LoopbackIPv4CIDR     = "100.64.0.0/24"
	SubnetIPv6CIDR       = "64:ff9b:1::/112"
	LoopbackIPv6CIDR     = "64:ff9b:1::/120"
	TestRegion           = "eu-west"
	TestAvailabilityZone = "A"
)

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
	inventoryGlobalCrdPath := getCrdPath(inventoriesv1alpha1.Inventory{})
	subnetGlobalCrdPath := getCrdPath(subnetv1alpha1.Subnet{})
	switchGlobalCrdPath := filepath.Join("..", "config", "crd", "bases")

	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{switchGlobalCrdPath, inventoryGlobalCrdPath, subnetGlobalCrdPath},
		ErrorIfCRDPathMissing: true,
		WebhookInstallOptions: envtest.WebhookInstallOptions{
			Paths: []string{filepath.Join("..", "config", "webhook")},
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

	By("Set up k8s-subnet manager")
	subnetScheme := scheme.Scheme
	err = subnetv1alpha1.AddToScheme(subnetScheme)
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
	err = (&subnet.NetworkReconciler{
		Client: k8sSubnetManager.GetClient(),
		Scheme: k8sSubnetManager.GetScheme(),
		Log:    ctrl.Log.WithName("k8s-subnet").WithName("network"),
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

	network := &subnetv1alpha1.Network{
		ObjectMeta: metav1.ObjectMeta{
			Name:      UnderlayNetwork,
			Namespace: OnmetalNamespace,
		},
		Spec: subnetv1alpha1.NetworkSpec{
			Description: "test underlay network",
		},
	}
	Expect(k8sClient.Create(ctx, network)).To(Succeed())

	By("Prepare subnets")
	cidrV4, _ := subnetv1alpha1.CIDRFromString(SubnetIPv4CIDR)
	subnetV4 := &subnetv1alpha1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      SubnetNameV4,
			Namespace: OnmetalNamespace,
		},
		Spec: subnetv1alpha1.SubnetSpec{
			CIDR:        cidrV4,
			NetworkName: UnderlayNetwork,
			Regions: []subnetv1alpha1.Region{
				{
					Name:              TestRegion,
					AvailabilityZones: []string{TestAvailabilityZone},
				},
			},
		},
	}
	Expect(k8sClient.Create(ctx, subnetV4)).To(Succeed())
	Eventually(func() bool {
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: SubnetNameV4, Namespace: OnmetalNamespace}, subnetV4)).Should(Succeed())
		return subnetV4.Status.State == subnetv1alpha1.CFinishedSubnetState
	}, timeout, interval).Should(BeTrue())

	loopbackCidrV4, _ := subnetv1alpha1.CIDRFromString(LoopbackIPv4CIDR)
	loopbackSubnetV4 := &subnetv1alpha1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      LoopbackSubnetV4,
			Namespace: OnmetalNamespace,
		},
		Spec: subnetv1alpha1.SubnetSpec{
			CIDR:             loopbackCidrV4,
			ParentSubnetName: SubnetNameV4,
			NetworkName:      UnderlayNetwork,
			Regions: []subnetv1alpha1.Region{
				{
					Name:              TestRegion,
					AvailabilityZones: []string{TestAvailabilityZone},
				},
			},
		},
	}
	Expect(k8sClient.Create(ctx, loopbackSubnetV4)).To(Succeed())
	Eventually(func() bool {
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: LoopbackSubnetV4, Namespace: OnmetalNamespace}, loopbackSubnetV4)).Should(Succeed())
		return loopbackSubnetV4.Status.State == subnetv1alpha1.CFinishedSubnetState
	}, timeout, interval).Should(BeTrue())

	cidrV6, _ := subnetv1alpha1.CIDRFromString(SubnetIPv6CIDR)
	subnetV6 := &subnetv1alpha1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      SubnetNameV6,
			Namespace: OnmetalNamespace,
		},
		Spec: subnetv1alpha1.SubnetSpec{
			CIDR:        cidrV6,
			NetworkName: UnderlayNetwork,
			Regions: []subnetv1alpha1.Region{
				{
					Name:              TestRegion,
					AvailabilityZones: []string{TestAvailabilityZone},
				},
			},
		},
	}
	Expect(k8sClient.Create(ctx, subnetV6)).To(Succeed())
	Eventually(func() bool {
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: SubnetNameV6, Namespace: OnmetalNamespace}, subnetV6)).Should(Succeed())
		return subnetV6.Status.State == subnetv1alpha1.CFinishedSubnetState
	}, timeout, interval).Should(BeTrue())

	loopbackCidrV6, _ := subnetv1alpha1.CIDRFromString(LoopbackIPv6CIDR)
	loopbackSubnetV6 := &subnetv1alpha1.Subnet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      LoopbackSubnetV6,
			Namespace: OnmetalNamespace,
		},
		Spec: subnetv1alpha1.SubnetSpec{
			CIDR:             loopbackCidrV6,
			ParentSubnetName: SubnetNameV6,
			NetworkName:      UnderlayNetwork,
			Regions: []subnetv1alpha1.Region{
				{
					Name:              TestRegion,
					AvailabilityZones: []string{TestAvailabilityZone},
				},
			},
		},
	}
	Expect(k8sClient.Create(ctx, loopbackSubnetV6)).To(Succeed())
	Eventually(func() bool {
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: LoopbackSubnetV6, Namespace: OnmetalNamespace}, loopbackSubnetV6)).Should(Succeed())
		return loopbackSubnetV6.Status.State == subnetv1alpha1.CFinishedSubnetState
	}, timeout, interval).Should(BeTrue())

}, 60)

var _ = AfterSuite(func() {
	By("Remove subnets")
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
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

func getCrdPath(crdPackageScheme interface{}) string {
	globalPackagePath := reflect.TypeOf(crdPackageScheme).PkgPath()
	goModData, err := ioutil.ReadFile(filepath.Join("..", "go.mod"))
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

func getUUID(identifier string) string {
	namespaceUUID := uuid.NewMD5(uuid.UUID{}, []byte(OnmetalNamespace))
	newUUID := uuid.NewMD5(namespaceUUID, []byte(identifier))
	return newUUID.String()
}
