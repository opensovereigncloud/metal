// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"

	ipamv1alpha1 "github.com/ironcore-dev/ipam/api/ipam/v1alpha1"
	"github.com/ironcore-dev/ironcore-image/oci/remote"
	poolv1alpha1 "github.com/ironcore-dev/ironcore/api/compute/v1alpha1"
	oobv1 "github.com/ironcore-dev/oob/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	benchmarkcontroller "github.com/ironcore-dev/metal/controllers/benchmark"
	inventorycontrollers "github.com/ironcore-dev/metal/controllers/inventory"
	machinecontrollers "github.com/ironcore-dev/metal/controllers/machine"
	onboardingcontroller "github.com/ironcore-dev/metal/controllers/onboarding"
	switchcontroller "github.com/ironcore-dev/metal/controllers/switch"
	onboardingprovider "github.com/ironcore-dev/metal/providers-kubernetes/onboarding"
	"github.com/ironcore-dev/metal/publisher"
	"github.com/ironcore-dev/metal/usecase/onboarding/invariants"
	"github.com/ironcore-dev/metal/usecase/onboarding/rules"
	onboardingscenarios "github.com/ironcore-dev/metal/usecase/onboarding/scenarios"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func addToScheme() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(metalv1alpha4.AddToScheme(scheme))
	utilruntime.Must(oobv1.AddToScheme(scheme))
	utilruntime.Must(ipamv1alpha1.AddToScheme(scheme))
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

	profHandlers := make(map[string]http.Handler)
	if profiling {
		profHandlers["debug/pprof"] = http.HandlerFunc(pprof.Index)
		profHandlers["debug/pprof/profile"] = http.HandlerFunc(pprof.Profile)
		setupLog.Info("profiling activated")
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress:   metricsAddr,
			ExtraHandlers: profHandlers,
		},
		WebhookServer:          webhook.NewServer(webhook.Options{Port: webhookPort}),
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "064f77d7.metal.ironcore.dev",
		Client: client.Options{
			Cache: &client.CacheOptions{
				DisableFor: []client.Object{
					&metalv1alpha4.Machine{},
				},
			},
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
	addHandlers(mgr)
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
		Log:    ctrl.Log.WithName("controllers").WithName("NetworkSwitch-onboarding"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "NetworkSwitch-onboarding")
		os.Exit(1)
	}
	if err = (&switchcontroller.SwitchReconciler{
		Client:   mgr.GetClient(),
		Log:      ctrl.Log.WithName("controllers").WithName("NetworkSwitch"),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("NetworkSwitch"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "NetworkSwitch")
		os.Exit(1)
	}
	if err = (&switchcontroller.IPAMReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("IPAM"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "NetworkSwitch-IPAM")
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
		ImageParser: &machinecontrollers.IroncoreImageParser{
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

func addHandlers(mgr ctrl.Manager) {
	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		if err := (&metalv1alpha4.Size{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Size")
			os.Exit(1)
		}
		if err := (&metalv1alpha4.Aggregate{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Aggregate")
			os.Exit(1)
		}
		if err := (&metalv1alpha4.NetworkSwitch{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "NetworkSwitch")
			os.Exit(1)
		}
		if err := (&metalv1alpha4.SwitchConfig{}).SetupWebhookWithManager(mgr); err != nil {
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
