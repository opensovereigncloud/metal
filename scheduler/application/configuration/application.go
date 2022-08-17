// /*
// Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package configuration

import (
	kubernetestprovider "github.com/onmetal/metal-api/pkg/provider/kubernetes-provider"
	"github.com/onmetal/metal-api/pkg/server"
	kubernetesmanager "github.com/onmetal/metal-api/pkg/server/kubernetes-manager"
	"github.com/onmetal/metal-api/scheduler/controllers"
	"github.com/onmetal/metal-api/scheduler/persistence-kubernetes/order"
	"github.com/onmetal/metal-api/scheduler/usecase/order/scenarios"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

func Initialize() (server.Server, error) {
	logger := NewLogger()
	config := newConfiguration()
	managerConfig := kubernetesmanager.ManagerConfig{
		MetricsAddr:          config.metricsAddr,
		EnableLeaderElection: config.enableLeaderElection,
		ProbeAddr:            config.probeAddr,
	}
	k8sManager, err := kubernetesmanager.NewManager(managerConfig)
	if err != nil {
		return nil, err
	}
	providerClient, err := kubernetestprovider.NewClient(k8sManager.Manager.GetClient())
	if err != nil {
		return nil, err
	}

	orderAlreadyScheduledExecutor := order.NewOrderAlreadyScheduled(providerClient, logger)
	orderAlreadyScheduledUseCase := scenarios.NewOrderAlreadyScheduledUseCase(orderAlreadyScheduledExecutor)
	instanceExtractor := order.NewInstanceFinderExtractor(providerClient)
	instanceForOrderUseCase := scenarios.NewFindVacantInstanceUseCase(instanceExtractor)
	cancelOrderExecutor := order.NewOrderCancelExecutor(providerClient)
	cancelOrderUseCase := scenarios.NewCancelOrderUseCase(cancelOrderExecutor)

	if err = controllers.NewSchedulerController(
		logger.WithName("controllers").WithName("Scheduler"),
		orderAlreadyScheduledUseCase,
		cancelOrderUseCase,
		instanceForOrderUseCase).SetupWithManager(k8sManager.Manager); err != nil {
		return nil, err
	}

	serverExtractor := order.NewServerExtractor(providerClient)
	instanceSchedulerUseCase := scenarios.NewInstanceSchedulerUseCase(serverExtractor)

	if err = controllers.NewInstanceScheduler(
		ctrl.Log.WithName("controllers").WithName("Machine-scheduler"),
		instanceSchedulerUseCase,
	).SetupWithManager(k8sManager.Manager); err != nil {
		return nil, err
	}

	if healthErr := k8sManager.Manager.AddHealthzCheck("healthz", healthz.Ping); healthErr != nil {
		return nil, healthErr
	}
	if readyErr := k8sManager.Manager.AddReadyzCheck("readyz", healthz.Ping); readyErr != nil {
		return nil, readyErr
	}
	return k8sManager, nil
}
