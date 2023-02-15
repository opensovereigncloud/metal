/*
Copyright 2022.

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

	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	"github.com/onmetal/metal-api/internal/repository"
	"github.com/onmetal/metal-api/internal/usecase"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//+kubebuilder:scaffold:imports
)

var (
	k8sClient client.Client
	testEnv   *envtest.Environment
	ctx       context.Context
	cancel    context.CancelFunc
)

const (
	timeout  = time.Second * 75
	interval = time.Millisecond * 250
)

var scheme = runtime.NewScheme()

// nolint
func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Onboarding Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}
	ctx, cancel = context.WithCancel(context.TODO())

	machinev1alpha2.SchemeBuilder.Register(&machinev1alpha2.Machine{}, &machinev1alpha2.MachineList{})
	machinev1alpha2.SchemeBuilder.Register(&machinev1alpha2.MachineAssignment{}, &machinev1alpha2.MachineAssignmentList{})
	inventoriesv1alpha1.SchemeBuilder.Register(&inventoriesv1alpha1.Inventory{}, &inventoriesv1alpha1.InventoryList{})
	switchv1beta1.SchemeBuilder.Register(&switchv1beta1.Switch{}, &switchv1beta1.SwitchList{})
	Expect(machinev1alpha2.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(inventoriesv1alpha1.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(switchv1beta1.AddToScheme(scheme)).NotTo(HaveOccurred())

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{Scheme: scheme, MetricsBindAddress: "0"})
	Expect(err).ToNot(HaveOccurred())

	deviceOnboardingRepo := repository.NewOnboardingRepo(k8sManager.GetClient())
	deviceOnboardingUseCase := usecase.NewDeviceOnboarding(deviceOnboardingRepo)

	err = (&OnboardingReconciler{
		Client:               k8sManager.GetClient(),
		Log:                  ctrl.Log.WithName("controllers").WithName("device-onboarding"),
		OnboardingRepo:       deviceOnboardingUseCase,
		DestinationNamespace: "default",
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
