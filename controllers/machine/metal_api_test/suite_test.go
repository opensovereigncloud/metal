// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"path/filepath"
	"testing"

	ipamv1alpha1 "github.com/ironcore-dev/ipam/api/ipam/v1alpha1"
	switchespkg "github.com/ironcore-dev/metal/pkg/switches"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	controllers "github.com/ironcore-dev/metal/controllers/machine"

	oobv1alpha1 "github.com/ironcore-dev/oob/api/v1alpha1"

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
	var (
		metalCRDPath, oobCRDPath, ipamCRDPath string
		err                                   error
	)

	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	By("bootstrapping test environment")

	metalCRDPath = filepath.Join("..", "..", "..", "config", "crd", "bases")
	ipamCRDPath, err = switchespkg.GetCrdPath(ipamv1alpha1.Subnet{}, filepath.Join("..", "..", "..", "go.mod"))
	Expect(err).ToNot(HaveOccurred())
	oobCRDPath, err = switchespkg.GetCrdPath(oobv1alpha1.OOB{}, filepath.Join("..", "..", "..", "go.mod"))
	Expect(err).ToNot(HaveOccurred())
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			metalCRDPath,
			oobCRDPath,
			ipamCRDPath,
		},
		ErrorIfCRDPathMissing: true,
	}
	ctx, cancel = context.WithCancel(context.TODO())

	oobv1alpha1.SchemeBuilder.Register(&oobv1alpha1.OOB{})

	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	Expect(oobv1alpha1.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(corev1.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(metalv1alpha4.AddToScheme(scheme)).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: "0",
		},
	})
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
