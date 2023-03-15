// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
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

package persistence

import (
	"context"
	"strings"

	inventories "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	usecase "github.com/onmetal/metal-api/internal/usecase/onboarding"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type InventoryRepository struct {
	client ctrlclient.Client
}

func NewInventoryRepository(client ctrlclient.Client) *InventoryRepository {
	return &InventoryRepository{client: client}
}

func (r *InventoryRepository) Create(inventory dto.CreateInventory) error {
	inv := prepareInventory(inventory)
	err := r.
		client.
		Create(
			context.Background(),
			inv)
	if apierrors.IsAlreadyExists(err) {
		return usecase.InventoryAlreadyOnboarded(inventory.Name)
	}
	return err
}

func (r *InventoryRepository) Get(request dto.Request) (dto.Inventory, error) {
	inv := &inventories.Inventory{}
	if err := r.
		client.
		Get(
			context.Background(),
			types.NamespacedName{
				Namespace: request.Namespace,
				Name:      request.Name,
			},
			inv); err != nil {
		return dto.Inventory{}, err
	}
	sizes := sizeLabels(inv.Labels)
	inventory := dto.NewInventory(
		inv.Name,
		inv.Namespace,
		inv.Spec.System.ProductSKU,
		inv.Spec.System.SerialNumber,
		sizes,
		inv.Spec.NICs)
	return inventory, nil
}

func prepareInventory(inv dto.CreateInventory) *inventories.Inventory {
	return &inventories.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Name:      inv.Name,
			Namespace: inv.Namespace,
		},
		Spec: inventories.InventorySpec{
			System: &inventories.SystemSpec{
				ID: "",
			},
			Host: &inventories.HostSpec{
				Name: "",
			},
		},
		Status: inventories.InventoryStatus{
			InventoryStatuses: inventories.InventoryStatuses{
				RequestsCount: 1,
			},
		},
	}
}

func sizeLabels(labels map[string]string) map[string]string {
	machineLabels := make(map[string]string, len(labels))
	for key, value := range labels {
		if !strings.Contains(key, inventories.CLabelPrefix) {
			continue
		}
		machineLabels[key] = value
	}
	return machineLabels
}
