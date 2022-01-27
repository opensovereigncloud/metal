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

	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	benchv1alpha3 "github.com/onmetal/metal-api/apis/benchmark/v1alpha3"
	machinev1lpha1 "github.com/onmetal/metal-api/apis/machine/v1alpha1"
	benchmarkcontroller "github.com/onmetal/metal-api/controllers/benchmark"
	machinecontroller "github.com/onmetal/metal-api/controllers/machine"

	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"

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
	webhookPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		webhookPort = 9443
	}

	var metricsAddr, probeAddr string
	var enableLeaderElection, profiling bool
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
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
	startReconcilers(mgr)
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
	utilruntime.Must(machinev1lpha1.AddToScheme(scheme))
	utilruntime.Must(inventoriesv1alpha1.AddToScheme(scheme))
	utilruntime.Must(switchv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func startReconcilers(mgr ctrl.Manager) {
	var err error
	if err = (&benchmarkcontroller.Reconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Benchmark"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Benchmark")
		os.Exit(1)
	}
	if err = (&machinecontroller.InventoryReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Inventory"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Inventory")
		os.Exit(1)
	}
	if err = (&machinecontroller.MachineReconciler{
		Client:   mgr.GetClient(),
		Log:      ctrl.Log.WithName("controllers").WithName("Machine"),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("machine-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Machine")
		os.Exit(1)
	}
	if err = (&machinecontroller.SwitchReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Switch"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Switch")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder
}

func addHandlers(mgr ctrl.Manager, profiling bool) {
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
