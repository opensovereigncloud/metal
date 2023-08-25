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
	benchv1alpha3 "github.com/onmetal/metal-api/apis/benchmark/v1alpha3"
	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1lpha3 "github.com/onmetal/metal-api/apis/machine/v1alpha3"
	switchv1beta1 "github.com/onmetal/metal-api/apis/switch/v1beta1"
	benchmarkcontroller "github.com/onmetal/metal-api/controllers/benchmark"
	inventorycontrollers "github.com/onmetal/metal-api/controllers/inventory"
	machinecontrollers "github.com/onmetal/metal-api/controllers/machine"
	onboardingcontroller "github.com/onmetal/metal-api/controllers/onboarding"
	switchcontroller "github.com/onmetal/metal-api/controllers/switch"
	onboardingprovider "github.com/onmetal/metal-api/providers-kubernetes/onboarding"
	"github.com/onmetal/metal-api/publisher"
	"github.com/onmetal/metal-api/usecase/onboarding/invariants"
	"github.com/onmetal/metal-api/usecase/onboarding/rules"
	onboardingscenarios "github.com/onmetal/metal-api/usecase/onboarding/scenarios"
	poolv1alpha1 "github.com/onmetal/onmetal-api/api/compute/v1alpha1"
	"github.com/onmetal/onmetal-image/oci/remote"
	oobv1 "github.com/onmetal/oob-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func addToScheme() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(benchv1alpha3.AddToScheme(scheme))
	utilruntime.Must(machinev1lpha3.AddToScheme(scheme))
	utilruntime.Must(inventoriesv1alpha1.AddToScheme(scheme))
	utilruntime.Must(oobv1.AddToScheme(scheme))
	utilruntime.Must(ipamv1alpha1.AddToScheme(scheme))
	utilruntime.Must(switchv1beta1.AddToScheme(scheme))
	utilruntime.Must(poolv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	addToScheme()
	webhookPort, err := strconv.Atoi(os.Getenv("WEBHOOK_PORT"))
	if err != nil {
		setupLog.Info("unable to read `WEBHOOK_PORT` env", "error", err)
		webhookPort = 9443
	}

	var metricsAddr, probeAddr, namespace, bootstrapAPIServer, loopbackSubnetLabelValue string
	var ipV6PrefixBits int
	var enableLeaderElection, profiling bool
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.StringVar(&namespace, "namespace", "default", "Namespace name for object creation")
	flag.StringVar(&bootstrapAPIServer, "bootstrap-api-server", "", "Endpoint of the the k8s api server to join to like https://1.2.3.4:6443")
	flag.StringVar(
		&loopbackSubnetLabelValue,
		"loopback_subnet_value_name",
		"loopback", "Loopback subnet label value name")
	flag.IntVar(
		&ipV6PrefixBits,
		"ip_v6_prefix_bits",
		64,
		"Subnet prefix bit length")
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
		WebhookServer:          webhook.NewServer(webhook.Options{Port: webhookPort}),
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "064f77d7.machine.onmetal.de",
		ClientDisableCacheFor: []client.Object{
			&machinev1lpha3.Machine{},
		},
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	envNS := os.Getenv("NAMESPACE")
	if namespace == "default" && envNS != "" {
		namespace = os.Getenv("NAMESPACE")
	}

	startReconcilers(mgr, bootstrapAPIServer, loopbackSubnetLabelValue)
	addHandlers(mgr, profiling)
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func startReconcilers(
	mgr ctrl.Manager,
	bootstrapAPIServer string,
	loopbackSubnetLabelValue string,
) {
	var err error
	eventPublisher := publisher.NewDomainEventPublisher(mgr.GetLogger())

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
	if err = (&switchcontroller.SwConfigReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("SwitchConfig"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "SwitchConfig")
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
		Client:         mgr.GetClient(),
		Log:            ctrl.Log.WithName("controllers").WithName("Size"),
		Scheme:         mgr.GetScheme(),
		EventPublisher: eventPublisher,
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
	if err = (&inventorycontrollers.AccessReconciler{
		Client:             mgr.GetClient(),
		Log:                ctrl.Log.WithName("controllers").WithName("Access"),
		Scheme:             mgr.GetScheme(),
		BootstrapAPIServer: bootstrapAPIServer,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Access")
		os.Exit(1)
	}

	if err = inventoryOnboardingReconciler(mgr, eventPublisher).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Inventory-onboarding")
		os.Exit(1)
	}

	if err := machineOnboardingReconciler(
		mgr,
		eventPublisher,
		loopbackSubnetLabelValue,
	).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Machine-onboarding")
		os.Exit(1)
	}

	if err = (&machinecontrollers.MachinePoolReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("Machine-Pool"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Machine-Pool")
		os.Exit(1)
	}
	if err = (&machinecontrollers.MachineReservationReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("Machine-Reservation"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Machine-Reservation")
		os.Exit(1)
	}
	if err = (&machinecontrollers.MachinePowerReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("Machine-Power"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Machine-Power")
		os.Exit(1)
	}

	registry, err := remote.DockerRegistry(nil)
	if err != nil {
		setupLog.Error(err, "unable to create registry")
		os.Exit(1)
	}
	if err = (&machinecontrollers.IpxeReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("Ipxe"),
		ImageParser: &machinecontrollers.OnmetalImageParser{
			Registry: registry,
			Log:      ctrl.Log.WithName("controllers").WithName("Ipxe").WithName("Image-Parser"),
		},
		Templater: &machinecontrollers.IpxeTemplater{},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Ipxe")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder
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

func inventoryOnboardingReconciler(
	mgr ctrl.Manager,
	eventPublisher *publisher.DomainEventPublisher,
) *onboardingcontroller.InventoryOnboardingReconciler {
	inventoryRepository := onboardingprovider.NewInventoryRepository(mgr.GetClient(), eventPublisher)
	inventoryIDGenerator := onboardingprovider.NewKubernetesInventoryIDGenerator()
	inventoryAlreadyExist := invariants.NewInventoryAlreadyExist(inventoryRepository)
	serverRepository := onboardingprovider.NewServerRepository(mgr.GetClient())
	serverExecutor := onboardingprovider.NewFakeServerExecutor(mgr.GetLogger())

	enableServerAfterInventoryCreationRule := rules.NewServerMustBeEnabledOnFirstTimeRule(
		serverExecutor,
		inventoryRepository,
		mgr.GetLogger(),
	)
	eventPublisher.RegisterListeners(enableServerAfterInventoryCreationRule)

	inventoryOnboardingUseCase := onboardingscenarios.NewCreateInventoryUseCase(
		inventoryAlreadyExist,
		inventoryIDGenerator,
		inventoryRepository,
	)
	getServerUseCase := onboardingscenarios.NewGetServerUseCase(serverRepository)

	return onboardingcontroller.NewInventoryOnboardingReconciler(
		ctrl.Log.WithName("controllers").WithName("Inventory-onboarding"),
		inventoryOnboardingUseCase,
		getServerUseCase,
	)
}

func machineOnboardingReconciler(
	mgr ctrl.Manager,
	eventPublisher *publisher.DomainEventPublisher,
	loopbackSubnetLabelValue string,
) *onboardingcontroller.OnboardingMachineReconciler {
	machineRepository := onboardingprovider.NewMachineRepository(mgr.GetClient(), eventPublisher)
	switchExtractor := onboardingprovider.NewSwitchRepository(mgr.GetClient())
	subnetRepository := onboardingprovider.NewLoopbackSubnetRepository(mgr.GetClient(), loopbackSubnetLabelValue)
	loopbackRepository := onboardingprovider.NewLoopbackAddressRepository(mgr.GetClient())
	machineAlreadyExist := invariants.NewMachineAlreadyExist(machineRepository)
	machineIDGenerator := onboardingprovider.NewKubernetesMachineIDGenerator()
	inventoryRepository := onboardingprovider.NewInventoryRepository(mgr.GetClient(), eventPublisher)

	createLoopback4ForMachineRule := rules.NewCreateLoopback4ForMachineRule(
		subnetRepository,
		loopbackRepository,
		inventoryRepository,
		mgr.GetLogger(),
	)
	createLoopback6ForMachineRule := rules.NewCreateLoopback6ForMachineRule(
		subnetRepository,
		loopbackRepository,
		inventoryRepository,
		mgr.GetLogger(),
	)

	createIPv6SubnetFromParentForInventoryRule := rules.NewCreateIPv6SubnetFromParentForInventoryRule(
		subnetRepository,
		subnetRepository,
		inventoryRepository,
		mgr.GetLogger(),
	)
	eventPublisher.RegisterListeners(
		createLoopback4ForMachineRule,
		createLoopback6ForMachineRule,
		createIPv6SubnetFromParentForInventoryRule)

	machineUseCase := onboardingscenarios.NewGetMachineUseCase(machineRepository)

	getInventoryUseCase := onboardingscenarios.NewGetInventoryUseCase(inventoryRepository)

	createMachineUseCase := onboardingscenarios.NewCreateMachineUseCase(
		machineRepository,
		machineIDGenerator,
		machineAlreadyExist,
	)

	machineOnboardUseCase := onboardingscenarios.NewMachineOnboardingUseCase(
		machineRepository,
		machineRepository,
		switchExtractor,
		loopbackRepository,
		mgr.GetLogger())

	return onboardingcontroller.NewOnboardingMachineReconciler(
		ctrl.Log.WithName("controllers").WithName("Machine-onboarding"),
		machineUseCase,
		getInventoryUseCase,
		createMachineUseCase,
		machineOnboardUseCase,
	)
}
