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

package repository

import (
	"context"
	"errors"

	inventoriesv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/internal/entity"
	machinerr "github.com/onmetal/metal-api/pkg/errors"
	oobv1 "github.com/onmetal/oob-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ServerOnboardingRepo struct {
	client    ctrlclient.Client
	inventory *inventoriesv1alpha1.Inventory
}

func NewServerOnboardingRepo(c ctrlclient.Client) *ServerOnboardingRepo {
	return &ServerOnboardingRepo{
		client: c,
	}
}

func (o *ServerOnboardingRepo) Create(ctx context.Context) error {
	if err := o.client.Create(ctx, o.inventory); err != nil {
		if apierrors.IsAlreadyExists(err) {
			return nil
		}
		return err
	}
	return nil
}

func (o *ServerOnboardingRepo) InitializationStatus(ctx context.Context, e entity.Onboarding) entity.Initialization {
	oobObj := &oobv1.OOB{}
	if err := o.client.Get(ctx, types.NamespacedName{
		Name: e.RequestName, Namespace: e.RequestNamespace}, oobObj); err != nil {
		return entity.Initialization{Require: false, Error: err}
	}

	inventory := &inventoriesv1alpha1.Inventory{}
	if err := o.client.Get(ctx, types.NamespacedName{
		Name: oobObj.Status.UUID, Namespace: e.InitializationObjectNamespace}, inventory); err != nil {
		return entity.Initialization{Require: true, Error: nil}
	}
	return entity.Initialization{Require: false, Error: nil}
}

func (o *ServerOnboardingRepo) Prepare(ctx context.Context, e entity.Onboarding) error {
	oobObj := &oobv1.OOB{}
	if err := o.client.Get(ctx, types.NamespacedName{
		Name: e.RequestName, Namespace: e.RequestNamespace}, oobObj); err != nil {
		return err
	}

	if oobObj.Status.UUID == "" {
		return machinerr.UUIDNotExist(e.RequestName)
	}

	e.InitializationObjectName = oobObj.Status.UUID
	o.inventory = prepareInventory(e)

	return nil
}

func (o *ServerOnboardingRepo) GatherData(ctx context.Context, e entity.Onboarding) error {
	oob := &oobv1.OOB{}
	if err := o.client.Get(ctx, types.NamespacedName{
		Name: e.RequestName, Namespace: e.RequestNamespace}, oob); err != nil {
		return err
	}

	inventory := &inventoriesv1alpha1.Inventory{}
	if err := o.client.Get(ctx, types.NamespacedName{
		Name: oob.Status.UUID, Namespace: e.InitializationObjectNamespace}, inventory); err != nil {
		return err
	}

	if o.IsSizeLabeled(inventory.Labels) {
		inventory.Status.InventoryStatuses.RequestsCount = 0
		return o.client.Update(ctx, inventory)
	}

	if inventory.Status.InventoryStatuses.RequestsCount > 1 {
		return errors.New("machine was booted but inventory not appeared")
	}

	if err := o.enableOOBMachineForInventory(ctx, oob); err != nil {
		return err
	}

	inventory.Status.InventoryStatuses.RequestsCount = 1
	return o.client.Update(ctx, inventory)
}

func (o *ServerOnboardingRepo) IsSizeLabeled(labels map[string]string) bool {
	machine := labels[inventoriesv1alpha1.GetSizeMatchLabel(machineSizeName)]
	switches := labels[inventoriesv1alpha1.GetSizeMatchLabel(switchSizeName)]
	return machine != "" || switches != ""
}

func (o *ServerOnboardingRepo) enableOOBMachineForInventory(ctx context.Context, oobObj *oobv1.OOB) error {
	oobObj.Spec.Power = getPowerState(oobObj.Spec.Power)
	oobObj.Labels = setUpLabels(oobObj)
	return o.client.Update(ctx, oobObj)
}

func prepareInventory(e entity.Onboarding) *inventoriesv1alpha1.Inventory {
	return &inventoriesv1alpha1.Inventory{
		ObjectMeta: metav1.ObjectMeta{
			Name:      e.InitializationObjectName,
			Namespace: e.InitializationObjectNamespace,
		},
		Spec: inventoriesv1alpha1.InventorySpec{
			System: &inventoriesv1alpha1.SystemSpec{
				ID: e.InitializationObjectName,
			},
			Host: &inventoriesv1alpha1.HostSpec{
				Name: "",
			},
		},
	}
}

func getPowerState(state string) string {
	switch state {
	case "On":
		// In case when machine already running Reset is required.
		// Machine should be started from scratch.
		// return "Reset"
		return state
	default:
		return "On"
	}
}

func setUpLabels(oobObj *oobv1.OOB) map[string]string {
	if oobObj.Labels == nil {
		return map[string]string{machinev1alpha2.UUIDLabel: oobObj.Status.UUID}
	}
	if _, ok := oobObj.Labels[machinev1alpha2.UUIDLabel]; !ok {
		oobObj.Labels[machinev1alpha2.UUIDLabel] = oobObj.Status.UUID
	}
	return oobObj.Labels
}
