/*
Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved

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

package v1beta1

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	ipamctrl "github.com/onmetal/ipam/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	switchespkg "github.com/ironcore-dev/metal/pkg/switches"
)

const defaultNamespace string = "ironcore"

const (
	timeout        = time.Second * 30
	contextTimeout = time.Second * 35
	interval       = time.Millisecond * 100
)

var (
	samplesPath = filepath.Join("..", "..", "test_samples", "switch")
)

var (
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
	ctx       context.Context
	cancel    context.CancelFunc
)

//nolint:paralleltest,gosec
func TestSwitchController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NetworkSwitch Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	By("bootstrapping test environment")

	metalapiCRDPath := filepath.Join("..", "..", "config", "crd", "bases")
	ipamCRDPath, err := switchespkg.GetCrdPath(ipamv1alpha1.Subnet{})
	Expect(err).ToNot(HaveOccurred())
	ipamWebhookPath, err := switchespkg.GetWebhookPath(ipamv1alpha1.Subnet{})
	Expect(err).ToNot(HaveOccurred())
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			metalapiCRDPath,
			ipamCRDPath,
		},
		ErrorIfCRDPathMissing: true,
		WebhookInstallOptions: envtest.WebhookInstallOptions{
			Paths: []string{
				filepath.Join("..", "..", "config", "webhook"),
				ipamWebhookPath,
			},
		},
	}
	ctx, cancel = context.WithCancel(context.TODO())

	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	ipamv1alpha1.SchemeBuilder.Register(
		&ipamv1alpha1.Network{},
		&ipamv1alpha1.Subnet{},
		&ipamv1alpha1.SubnetList{},
		&ipamv1alpha1.IP{},
		&ipamv1alpha1.IPList{},
	)

	scheme := runtime.NewScheme()
	Expect(corev1.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(admissionv1beta1.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(metalv1alpha4.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(ipamv1alpha1.AddToScheme(scheme)).NotTo(HaveOccurred())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	webhookInstallOptions := &testEnv.WebhookInstallOptions
	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme,
		WebhookServer: webhook.NewServer(webhook.Options{
			Port:    webhookInstallOptions.LocalServingPort,
			Host:    webhookInstallOptions.LocalServingHost,
			CertDir: webhookInstallOptions.LocalServingCertDir,
		}),
		LeaderElection: false,
		Metrics: metricsserver.Options{
			BindAddress: "0",
		},
	})
	Expect(err).ToNot(HaveOccurred())

	Expect((&SwitchReconciler{
		Client:   k8sManager.GetClient(),
		Scheme:   k8sManager.GetScheme(),
		Recorder: k8sManager.GetEventRecorderFor("switch-reconciler"),
		Log:      ctrl.Log.WithName("controllers").WithName("switch-reconciler"),
	}).SetupWithManager(k8sManager)).NotTo(HaveOccurred())
	Expect((&OnboardingReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("onboarding-reconciler"),
	}).SetupWithManager(k8sManager)).NotTo(HaveOccurred())
	Expect((&metalv1alpha4.NetworkSwitch{}).SetupWebhookWithManager(k8sManager)).NotTo(HaveOccurred())
	Expect((&metalv1alpha4.SwitchConfig{}).SetupWebhookWithManager(k8sManager)).NotTo(HaveOccurred())
	Expect((&IPAMReconciler{
		Client:                  k8sManager.GetClient(),
		Scheme:                  k8sManager.GetScheme(),
		Recorder:                k8sManager.GetEventRecorderFor("switch-ipam-reconciler"),
		Log:                     ctrl.Log.WithName("controllers").WithName("switch-ipam-reconciler"),
		SwitchPortsIPAMDisabled: true,
	}).SetupWithManager(k8sManager)).NotTo(HaveOccurred())

	Expect((&ipamctrl.SubnetReconciler{
		Client:        k8sManager.GetClient(),
		Scheme:        k8sManager.GetScheme(),
		EventRecorder: k8sManager.GetEventRecorderFor("subnet-reconciler"),
		Log:           ctrl.Log.WithName("controllers").WithName("subnet-reconciler"),
	}).SetupWithManager(k8sManager)).NotTo(HaveOccurred())
	Expect((&ipamv1alpha1.Subnet{}).SetupWebhookWithManager(k8sManager)).NotTo(HaveOccurred())
	Expect((&ipamctrl.IPReconciler{
		Client:        k8sManager.GetClient(),
		Scheme:        k8sManager.GetScheme(),
		EventRecorder: k8sManager.GetEventRecorderFor("ip-reconciler"),
		Log:           ctrl.Log.WithName("controllers").WithName("ip-reconciler"),
	}).SetupWithManager(k8sManager)).NotTo(HaveOccurred())
	Expect((&ipamv1alpha1.IP{}).SetupWebhookWithManager(k8sManager)).NotTo(HaveOccurred())
	Expect((&ipamctrl.NetworkReconciler{
		Client:        k8sManager.GetClient(),
		Scheme:        k8sManager.GetScheme(),
		EventRecorder: k8sManager.GetEventRecorderFor("network-reconciler"),
		Log:           ctrl.Log.WithName("controllers").WithName("network-reconciler"),
	}).SetupWithManager(k8sManager)).NotTo(HaveOccurred())
	Expect((&ipamv1alpha1.Network{}).SetupWebhookWithManager(k8sManager)).NotTo(HaveOccurred())

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

	seed(ctx, k8sClient)
})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

func seed(ctx context.Context, c client.Client) {
	Expect(seedNamespace(ctx, c)).NotTo(HaveOccurred())
	// Expect(seedInventories(ctx, c)).NotTo(HaveOccurred())
	Expect(seedConfigs(ctx, c)).NotTo(HaveOccurred())
	Expect(seedNetworks(ctx, c)).NotTo(HaveOccurred())
	Expect(seedSubnets(ctx, c)).NotTo(HaveOccurred())
}

func seedNamespace(ctx context.Context, c client.Client) error {
	obj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: defaultNamespace},
		TypeMeta:   metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
	}
	err := c.Create(ctx, obj)
	return err
}

func seedInventories(ctx context.Context, c client.Client) error {
	samplesPath := filepath.Join(samplesPath, "inventories")
	samples, err := switchespkg.GetTestSamples(samplesPath)
	if err != nil {
		return err
	}
	for _, sample := range samples {
		raw, err := os.ReadFile(sample)
		if err != nil {
			return err
		}
		obj := &metalv1alpha4.Inventory{}
		if err := switchespkg.CreateSampleObject(ctx, c, obj, raw); err != nil {
			return err
		}
	}
	return nil
}

func seedConfigs(ctx context.Context, c client.Client) error {
	samplesPath := filepath.Join(samplesPath, "switch_configs")
	samples, err := switchespkg.GetTestSamples(samplesPath)
	if err != nil {
		return err
	}
	for _, sample := range samples {
		raw, err := os.ReadFile(sample)
		if err != nil {
			return err
		}
		obj := &metalv1alpha4.SwitchConfig{}
		if err := switchespkg.CreateSampleObject(ctx, c, obj, raw); err != nil {
			return err
		}
	}
	return nil
}

func seedNetworks(ctx context.Context, c client.Client) error {
	samplesPath := filepath.Join(samplesPath, "networks")
	samples, err := switchespkg.GetTestSamples(samplesPath)
	if err != nil {
		return err
	}
	for _, sample := range samples {
		raw, err := os.ReadFile(sample)
		if err != nil {
			return err
		}
		obj := &ipamv1alpha1.Network{}
		if err := switchespkg.CreateSampleObject(ctx, c, obj, raw); err != nil {
			return err
		}
	}
	networks := &ipamv1alpha1.NetworkList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(ctx, networks)).NotTo(HaveOccurred())
		for _, item := range networks.Items {
			g.Expect(item.Status.State).To(Equal(ipamv1alpha1.CFinishedNetworkState))
		}
	}, timeout, interval).Should(Succeed())
	return nil
}

func seedSubnets(ctx context.Context, c client.Client) error {
	samplesPath := filepath.Join(samplesPath, "subnets")
	samples, err := switchespkg.GetTestSamples(samplesPath)
	if err != nil {
		return err
	}
	for _, sample := range samples {
		raw, err := os.ReadFile(sample)
		if err != nil {
			return err
		}
		obj := &ipamv1alpha1.Subnet{}
		if err := switchespkg.CreateSampleObject(ctx, c, obj, raw); err != nil {
			return err
		}
	}
	subnets := &ipamv1alpha1.SubnetList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(ctx, subnets)).NotTo(HaveOccurred())
		for _, item := range subnets.Items {
			g.Expect(item.Status.State).To(Equal(ipamv1alpha1.CFinishedSubnetState))
		}
	}, timeout, interval).Should(Succeed())
	return nil
}

func seedSwitches(ctx context.Context, c client.Client) error {
	samplesPath := filepath.Join(samplesPath, "switches")
	samples, err := switchespkg.GetTestSamples(samplesPath)
	if err != nil {
		return err
	}
	for _, sample := range samples {
		raw, err := os.ReadFile(sample)
		if err != nil {
			return err
		}
		obj := &metalv1alpha4.NetworkSwitch{}
		if err := switchespkg.CreateSampleObject(ctx, c, obj, raw); err != nil {
			return err
		}
	}
	return nil
}

func seedSwitchesSubnets(ctx context.Context, c client.Client) error {
	samplesPath := filepath.Join(samplesPath, "switch_ipam_objects", "subnets")
	samples, err := switchespkg.GetTestSamples(samplesPath)
	if err != nil {
		return err
	}
	for _, sample := range samples {
		raw, err := os.ReadFile(sample)
		if err != nil {
			return err
		}
		obj := &ipamv1alpha1.Subnet{}
		if err := switchespkg.CreateSampleObject(ctx, c, obj, raw); err != nil {
			return err
		}
	}
	subnets := &ipamv1alpha1.SubnetList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(ctx, subnets)).NotTo(HaveOccurred())
		for _, item := range subnets.Items {
			g.Expect(item.Status.State).To(Equal(ipamv1alpha1.CFinishedSubnetState))
		}
	}, timeout, interval).Should(Succeed())
	return nil
}

func seedSwitchesLoopbacks(ctx context.Context, c client.Client) error {
	samplesPath := filepath.Join(samplesPath, "switch_ipam_objects", "loopbacks")
	samples, err := switchespkg.GetTestSamples(samplesPath)
	if err != nil {
		return err
	}
	for _, sample := range samples {
		raw, err := os.ReadFile(sample)
		if err != nil {
			return err
		}
		obj := &ipamv1alpha1.IP{}
		if err := switchespkg.CreateSampleObject(ctx, c, obj, raw); err != nil {
			return err
		}
	}
	loopbacks := &ipamv1alpha1.IPList{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.List(ctx, loopbacks)).NotTo(HaveOccurred())
		for _, item := range loopbacks.Items {
			g.Expect(item.Status.State).To(Equal(ipamv1alpha1.CFinishedIPState))
		}
	}, timeout, interval).Should(Succeed())
	return nil
}
