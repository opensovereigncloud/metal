// Copyright 2022 OnMetal authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"

	inventoriesv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	machinev1alpha1 "github.com/onmetal/metal-api/apis/machine/v1alpha1"
	machinerr "github.com/onmetal/metal-api/internal/errors"
	oobonmetal "github.com/onmetal/oob-controller/api/v1"
	switchv1alpha1 "github.com/onmetal/switch-operator/api/v1alpha1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func List(ctx context.Context, c ctrlclient.Client, k kind) (ctrlclient.ObjectList, error) {
	switch k {
	case Machine:
		return getMachineList(ctx, c)
	case Inventory:
		return getInventoryList(ctx, c)
	case Switch:
		return getSwitchList(ctx, c)
	case OOB:
		return getOOBList(ctx, c)
	default:
		return nil, machinerr.NotExist(string(k))
	}
}

func getMachineList(ctx context.Context, c ctrlclient.Client) (*machinev1alpha1.MachineList, error) {
	obj := &machinev1alpha1.MachineList{}
	if err := c.List(ctx, obj); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, machinerr.NotFound(string(Machine))
	}
	return obj, nil
}

func getInventoryList(ctx context.Context, c ctrlclient.Client) (*inventoriesv1alpha1.InventoryList, error) {
	obj := &inventoriesv1alpha1.InventoryList{}
	if err := c.List(ctx, obj); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, machinerr.NotFound(string(Inventory))
	}
	return obj, nil
}

func getSwitchList(ctx context.Context, c ctrlclient.Client) (*switchv1alpha1.SwitchList, error) {
	obj := &switchv1alpha1.SwitchList{}
	if err := c.List(ctx, obj); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, machinerr.NotFound(string(Switch))
	}
	return obj, nil
}

func getOOBList(ctx context.Context, c ctrlclient.Client) (*oobonmetal.MachineList, error) {
	obj := &oobonmetal.MachineList{}
	if err := c.List(ctx, obj); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, machinerr.NotFound(string(Switch))
	}
	return obj, nil
}
