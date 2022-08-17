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

package kubernetesmanager

import (
	"context"
	"log"

	ipam "github.com/onmetal/ipam/api/v1alpha1"
	benchv1alpha3 "github.com/onmetal/metal-api/apis/benchmark/v1alpha3"
	inventoryv1alpaha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpaha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	authv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ControllerManager struct {
	Manager ctrl.Manager
}

type ManagerConfig struct {
	MetricsAddr          string
	EnableLeaderElection bool
	ProbeAddr            string
}

func NewManager(config ManagerConfig) (*ControllerManager, error) {
	scheme, err := getScheme()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	controllerManager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     config.MetricsAddr,
		HealthProbeBindAddress: config.ProbeAddr,
		LeaderElection:         config.EnableLeaderElection,
		LeaderElectionID:       "38b1eb41.onmetal.de",
	})
	if err != nil {
		return nil, err
	}
	return &ControllerManager{Manager: controllerManager}, nil
}

func getScheme() (*k8sRuntime.Scheme, error) {
	scheme := k8sRuntime.NewScheme()
	if err := benchv1alpha3.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := inventoryv1alpaha1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := machinev1alpaha2.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := authv1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := ipam.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return scheme, nil
}

func (c *ControllerManager) Start() error {
	log.Println("controller manager starting")
	return c.Manager.Start(context.Background())
}
