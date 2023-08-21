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

package providers

import (
	"context"
	"strings"

	inventories "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	"github.com/onmetal/metal-api/common/types/events"
	domain "github.com/onmetal/metal-api/domain/inventory"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type InventoryRepository struct {
	client               ctrlclient.Client
	domainEventPublisher events.DomainEventPublisher
}

func NewInventoryRepository(
	client ctrlclient.Client,
	domainEventPublisher events.DomainEventPublisher,
) *InventoryRepository {
	return &InventoryRepository{
		client:               client,
		domainEventPublisher: domainEventPublisher,
	}
}

func (r *InventoryRepository) Save(inventory domain.Inventory) error {
	inv := prepareInventory(inventory)
	if err := r.
		client.
		Create(
			context.Background(),
			inv); err != nil {
		return err
	}
	r.domainEventPublisher.Publish(inventory.PopEvents()...)
	return nil
}

func (r *InventoryRepository) ByUUID(uuid string) (domain.Inventory, error) {
	uuidOptions := ctrlclient.MatchingFields{
		"metadata.name": uuid,
	}
	inv, err := r.extractInventoryFromCluster(uuidOptions)
	if err != nil {
		return domain.Inventory{}, err
	}
	return domainInventory(inv), nil
}

func (r *InventoryRepository) ByID(id string) (domain.Inventory, error) {
	idOptions := &ctrlclient.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{"id": id})}
	inv, err := r.extractInventoryFromCluster(idOptions)
	if err != nil {
		return domain.Inventory{}, err
	}
	return domainInventory(inv), nil
}

func (r *InventoryRepository) extractInventoryFromCluster(
	options ctrlclient.ListOption,
) (*inventories.Inventory, error) {
	obj := &inventories.InventoryList{}
	if err := r.
		client.
		List(
			context.Background(),
			obj,
			options,
		); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, errNotFound
	}
	return &obj.Items[0], nil
}

func prepareInventory(inv domain.Inventory) *inventories.Inventory {
	return &inventories.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Name:      inv.UUID,
			Namespace: inv.Namespace,
			Labels: map[string]string{
				"id": inv.ID.String(),
			},
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

func domainInventory(inv *inventories.Inventory) domain.Inventory {
	sizes := sizeLabels(inv.Labels)

	var productSKU, serialNumber string
	if inv.Spec.System != nil {
		productSKU = inv.Spec.System.ProductSKU
		serialNumber = inv.Spec.System.SerialNumber
	}
	return domain.NewInventory(
		domain.NewInventoryID(inv.Labels["id"]),
		inv.Name,
		inv.Namespace,
		productSKU,
		serialNumber,
		sizes,
		inv.Spec.NICs)
}
