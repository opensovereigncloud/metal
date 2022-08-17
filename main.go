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

package main

import (
	"flag"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"

	ipamv1alpha1 "github.com/onmetal/ipam/api/v1alpha1"
	oobv1 "github.com/onmetal/oob-controller/api/v1"

	benchv1alpha3 "github.com/onmetal/metal-api/apis/benchmark/v1alpha3"
	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1lpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	benchmarkcontroller "github.com/onmetal/metal-api/controllers/benchmark"
	inventorycontrollers "github.com/onmetal/metal-api/controllers/inventory"
	machinecontroller "github.com/onmetal/metal-api/controllers/machine"
	schedulercontrollers "github.com/onmetal/metal-api/controllers/machine-scheduler"
	onboardingcontroller "github.com/onmetal/metal-api/controllers/onboarding"
	switchcontroller "github.com/onmetal/metal-api/controllers/switch/v1beta1"

	"github.com/onmetal/metal-api/internal/repository"
	"github.com/onmetal/metal-api/internal/usecase"

	// to ensure that exec-entrypoint and run can make use of them.
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func main() {
	addToScheme()
	webhookPort, err := strconv.Atoi(os.Getenv("WEBHOOK_PORT"))
	if err != nil {
		setupLog.Info("unable to read `WEBHOOK_PORT` env", "error", err)
		webhookPort = 9443
	}

	var metricsAddr, probeAddr, namespace string
	var enableLeaderElection, profiling bool
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.StringVar(&namespace, "namespace", "default", "Namespace name for object creation")
	flag.IntVar(&webhookPort, "webhook-bind-address", webhookPort, "The address the webhook endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&profiling, "profiling", false, "Enabling this will activate profiling that will be listen on :8080")
	opts := zap.Options{
		Development: false,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   webhookPort,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "064f77d7.machine.onmetal.de",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	envNS := os.Getenv("NAMESPACE")
	if namespace == "default" && envNS != "" {
		namespace = os.Getenv("NAMESPACE")
	}

	startReconcilers(mgr, namespace)
	addHandlers(mgr, profiling)
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func addToScheme() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(benchv1alpha3.AddToScheme(scheme))
	utilruntime.Must(machinev1lpha2.AddToScheme(scheme))
	utilruntime.Must(inventoriesv1alpha1.AddToScheme(scheme))
	utilruntime.Must(oobv1.AddToScheme(scheme))
	utilruntime.Must(ipamv1alpha1.AddToScheme(scheme))
	utilruntime.Must(switchv1beta1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func startReconcilers(mgr ctrl.Manager, namespace string) {
	var err error

	deviceOnboardingRepo := repository.NewOnboardingRepo(mgr.GetClient())
	deviceOnboardingUseCase := usecase.NewDeviceOnboarding(deviceOnboardingRepo)

	serverOnboardingRepo := repository.NewServerOnboardingRepo(mgr.GetClient())
	serverOnboardingUseCase := usecase.NewServerOnboarding(serverOnboardingRepo)

	if err = (&benchmarkcontroller.Reconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Benchmark"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Benchmark")
		os.Exit(1)
	}
	if err = (&benchmarkcontroller.OnboardingReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Benchmark-onboarding"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Benchmark-onboarding")
		os.Exit(1)
	}
	if err = (&machinecontroller.InventoryReconciler{
		Client:    mgr.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("Machine-inventory"),
		Scheme:    mgr.GetScheme(),
		Recorder:  mgr.GetEventRecorderFor("Machine-inventory"),
		Namespace: namespace,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Machine-inventory")
		os.Exit(1)
	}
	if err = (&machinecontroller.OOBReconciler{
		Client:    mgr.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("Machine-OOB"),
		Scheme:    mgr.GetScheme(),
		Recorder:  mgr.GetEventRecorderFor("Machine-OOB"),
		Namespace: namespace,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Machine-OOB")
		os.Exit(1)
	}
	if err = (&switchcontroller.OnboardingReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Switch-onboarding"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Switch-onboarding")
		os.Exit(1)
	}
	if err = (&switchcontroller.SwitchReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Switch"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Switch")
		os.Exit(1)
	}
	if err = (&inventorycontrollers.InventoryReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Inventory"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Inventory")
		os.Exit(1)
	}
	if err = (&inventorycontrollers.SizeReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Size"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Size")
		os.Exit(1)
	}
	if err = (&inventorycontrollers.AggregateReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Aggregate"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Aggregate")
		os.Exit(1)
	}
	if err = (&onboardingcontroller.OnboardingReconciler{
		Client:               mgr.GetClient(),
		Log:                  ctrl.Log.WithName("controllers").WithName("Device-onboarding"),
		Scheme:               mgr.GetScheme(),
		OnboardingRepo:       deviceOnboardingUseCase,
		DestinationNamespace: namespace,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Device-onboarding")
		os.Exit(1)
	}
	if err = (&onboardingcontroller.InventoryOnboardingReconciler{
		Client:               mgr.GetClient(),
		Log:                  ctrl.Log.WithName("controllers").WithName("Server-onboarding"),
		Scheme:               mgr.GetScheme(),
		OnboardingRepo:       serverOnboardingUseCase,
		DestinationNamespace: namespace,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Server-onboarding")
		os.Exit(1)
	}
	if err = (&schedulercontrollers.IgnitionReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("Ignition"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Ignition")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder
}

func addHandlers(mgr ctrl.Manager, profiling bool) {
	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		if err := (&inventoriesv1alpha1.Size{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Size")
			os.Exit(1)
		}
		if err := (&inventoriesv1alpha1.Aggregate{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Aggregate")
			os.Exit(1)
		}
		if err := (&switchv1beta1.Switch{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Switch")
			os.Exit(1)
		}
		if err := (&switchv1beta1.SwitchConfig{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "SwitchConfig")
			os.Exit(1)
		}
	}
	if healthErr := mgr.AddHealthzCheck("healthz", healthz.Ping); healthErr != nil {
		setupLog.Error(healthErr, "unable to set up health check")
		os.Exit(1)
	}
	if readyErr := mgr.AddReadyzCheck("readyz", healthz.Ping); readyErr != nil {
		setupLog.Error(readyErr, "unable to set up ready check")
		os.Exit(1)
	}
	if profiling {
		err := mgr.AddMetricsExtraHandler("/debug/pprof/", http.HandlerFunc(pprof.Index))
		if err != nil {
			setupLog.Error(err, "unable to attach pprof to webserver")
			os.Exit(1)
		}
		err = mgr.AddMetricsExtraHandler("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		if err != nil {
			setupLog.Error(err, "unable to attach cpu pprof to webserver")
			os.Exit(1)
		}
		setupLog.Info("profiling activated")
	}
}
