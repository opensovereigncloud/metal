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

	controllers "github.com/onmetal/metal-api/controllers/machine"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/onmetal/controller-utils/buildutils"
	"github.com/onmetal/controller-utils/modutils"
	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	computev1alpha1 "github.com/onmetal/onmetal-api/api/compute/v1alpha1"
	utilsenvtest "github.com/onmetal/onmetal-api/utils/envtest"
	"github.com/onmetal/onmetal-api/utils/envtest/apiserver"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg            *rest.Config
	k8sClient      client.Client
	reconcilersMap map[string]func(k8sManager manager.Manager, log logr.Logger)
)

const (
	pollingInterval      = 250 * time.Millisecond
	eventuallyTimeout    = 15 * time.Second
	consistentlyDuration = 1 * time.Second
	apiServiceTimeout    = 5 * time.Minute

	machinePoolReconcilers        = "machine-pool-reconcilers"
	machineReservationReconcilers = "machine-reservation-reconcilers"
	ipxeReconcilers               = "ipxe-reconcilers"
)

// nolint
func TestMachinePoolController(t *testing.T) {
	SetDefaultConsistentlyPollingInterval(pollingInterval)
	SetDefaultEventuallyPollingInterval(pollingInterval)
	SetDefaultEventuallyTimeout(eventuallyTimeout)
	SetDefaultConsistentlyDuration(consistentlyDuration)
	RegisterFailHandler(Fail)
	reconcilersMap = make(map[string]func(k8sManager manager.Manager, log logr.Logger))

	reconcilersMap[machinePoolReconcilers] = func(k8sManager manager.Manager, log logr.Logger) {
		err := (&controllers.MachinePoolReconciler{
			Client: k8sManager.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("machine-pool"),
			Scheme: k8sManager.GetScheme(),
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())
	}

	reconcilersMap[machineReservationReconcilers] = func(k8sManager manager.Manager, log logr.Logger) {
		err := (&controllers.MachineReservationReconciler{
			Client: k8sManager.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("machine-reservation"),
			Scheme: k8sManager.GetScheme(),
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())

		err = (&controllers.MachinePoolReconciler{
			Client: k8sManager.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("machine-pool"),
			Scheme: k8sManager.GetScheme(),
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())

		err = (&controllers.IpxeReconciler{
			Client:      k8sManager.GetClient(),
			Log:         ctrl.Log.WithName("controllers").WithName("ipxe"),
			Scheme:      k8sManager.GetScheme(),
			ImageParser: &OnmetalImageParserFake{},
			Templater:   &controllers.IpxeTemplater{},
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())
	}

	reconcilersMap[ipxeReconcilers] = func(k8sManager manager.Manager, log logr.Logger) {
		var err error

		err = (&controllers.MachineReservationReconciler{
			Client: k8sManager.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("machine-reservation"),
			Scheme: k8sManager.GetScheme(),
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())

		err = (&controllers.IpxeReconciler{
			Client:      k8sManager.GetClient(),
			Log:         ctrl.Log.WithName("controllers").WithName("ipxe"),
			Scheme:      k8sManager.GetScheme(),
			ImageParser: &OnmetalImageParserFake{},
			Templater:   &controllers.IpxeTemplater{},
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())
	}

	RunSpecs(t, "Compute machine suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	var err error

	By("bootstrapping test environment")
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
	}
	testEnvExt := &utilsenvtest.EnvironmentExtensions{
		APIServiceDirectoryPaths: []string{
			modutils.Dir("github.com/onmetal/onmetal-api", "config", "apiserver", "apiservice", "bases"),
		},
		ErrorIfAPIServicePathIsMissing: true,
	}

	cfg, err = utilsenvtest.StartWithExtensions(testEnv, testEnvExt)
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())
	DeferCleanup(utilsenvtest.StopWithExtensions, testEnv, testEnvExt)

	Expect(inventoriesv1alpha1.AddToScheme(scheme.Scheme)).NotTo(HaveOccurred())
	Expect(corev1.AddToScheme(scheme.Scheme)).NotTo(HaveOccurred())
	Expect(machinev1alpha2.AddToScheme(scheme.Scheme)).NotTo(HaveOccurred())
	Expect(computev1alpha1.AddToScheme(scheme.Scheme)).To(Succeed())
	Expect(apiregistrationv1.AddToScheme(scheme.Scheme)).To(Succeed())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
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

	Expect(apiSrv.Start()).To(Succeed())
	DeferCleanup(apiSrv.Stop)

	Expect(utilsenvtest.WaitUntilAPIServicesReadyWithTimeout(apiServiceTimeout, testEnvExt, k8sClient, scheme.Scheme)).To(Succeed())
})

func SetupTest(ctx context.Context, reconcilers string) *corev1.Namespace {
	var (
		cancel context.CancelFunc
		ns     = &corev1.Namespace{}
	)

	BeforeEach(func() {
		var mgrCtx context.Context
		mgrCtx, cancel = context.WithCancel(ctx)
		*ns = corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "testns-",
			},
		}
		Expect(k8sClient.Create(ctx, ns)).To(Succeed(), "failed to create test namespace")

		k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
			Scheme:             scheme.Scheme,
			MetricsBindAddress: "0",
		})
		Expect(err).ToNot(HaveOccurred())

		reconcilersMap[reconcilers](k8sManager, ctrl.Log)

		go func() {
			defer GinkgoRecover()
			Expect(k8sManager.Start(mgrCtx)).To(Succeed(), "failed to start manager")
		}()
	})

	AfterEach(func() {
		cancel()
		Expect(k8sClient.Delete(ctx, ns)).To(Succeed(), "failed to delete test namespace")
		Expect(k8sClient.DeleteAllOf(ctx, &computev1alpha1.MachinePool{})).To(Succeed())
		Expect(k8sClient.DeleteAllOf(ctx, &computev1alpha1.MachineClass{})).To(Succeed())
	})

	return ns
}
