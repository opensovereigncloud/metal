// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/ironcore-dev/metal/common/types/events"
	domain "github.com/ironcore-dev/metal/domain/inventory"
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

func (r *InventoryRepository) ByID(id domain.InventoryID) (domain.Inventory, error) {
	idOptions := &ctrlclient.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{"id": id.String()})}
	inv, err := r.extractInventoryFromCluster(idOptions)
	if err != nil {
		return domain.Inventory{}, err
	}
	return domainInventory(inv), nil
}

func (r *InventoryRepository) extractInventoryFromCluster(
	options ctrlclient.ListOption,
) (*metalv1alpha4.Inventory, error) {
	obj := &metalv1alpha4.InventoryList{}
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

func prepareInventory(inv domain.Inventory) *metalv1alpha4.Inventory {
	return &metalv1alpha4.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Name:      inv.UUID,
			Namespace: inv.Namespace,
			Labels: map[string]string{
				"id": inv.ID.String(),
			},
		},
		Spec: metalv1alpha4.InventorySpec{
			System: &metalv1alpha4.SystemSpec{
				ID: "",
			},
			Host: &metalv1alpha4.HostSpec{
				Name: "",
			},
		},
		Status: metalv1alpha4.InventoryStatus{
			InventoryStatuses: metalv1alpha4.InventoryStatuses{
				RequestsCount: 1,
			},
		},
	}
}

func sizeLabels(labels map[string]string) map[string]string {
	machineLabels := make(map[string]string, len(labels))
	for key, value := range labels {
		if !strings.Contains(key, metalv1alpha4.CLabelPrefix) {
			continue
		}
		machineLabels[key] = value
	}
	return machineLabels
}

func domainInventory(inv *metalv1alpha4.Inventory) domain.Inventory {
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
