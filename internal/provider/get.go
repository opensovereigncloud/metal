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
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Machine   kind = "machine"
	OOB       kind = "oob"
	Inventory kind = "inventory"
	Switch    kind = "switch"
)

type kind string

func Get(ctx context.Context, c ctrlclient.Client, name, namespace string, k kind) (ctrlclient.Object, error) {
	switch k {
	case Machine:
		return getMachine(ctx, c, name, namespace)
	case Inventory:
		return getInventory(ctx, c, name, namespace)
	case Switch:
		return getSwitch(ctx, c, name, namespace)
	case OOB:
		return getOOB(ctx, c, name, namespace)
	default:
		return nil, machinerr.NotExist(string(k))
	}
}

func GetByLabel(ctx context.Context, c ctrlclient.Client, label map[string]string, k kind) (ctrlclient.Object, error) {
	switch k {
	case Machine:
		return getMachineByLabel(ctx, c, label)
	case Inventory:
		return getInventoryByLabel(ctx, c, label)
	case Switch:
		return getSwitchByLabel(ctx, c, label)
	case OOB:
		return getOOBByLabel(ctx, c, label)
	default:
		return nil, machinerr.NotExist(string(k))
	}
}

func getMachine(ctx context.Context, c ctrlclient.Client, name, namespace string) (*machinev1alpha1.Machine, error) {
	obj := &machinev1alpha1.Machine{}
	if err := c.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func getInventory(ctx context.Context, c ctrlclient.Client, name, namespace string) (*inventoriesv1alpha1.Inventory, error) {
	obj := &inventoriesv1alpha1.Inventory{}
	if err := c.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func getSwitch(ctx context.Context, c ctrlclient.Client, name, namespace string) (*switchv1alpha1.Switch, error) {
	obj := &switchv1alpha1.Switch{}
	if err := c.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func getOOB(ctx context.Context, c ctrlclient.Client, name, namespace string) (*oobonmetal.Machine, error) {
	obj := &oobonmetal.Machine{}
	if err := c.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func getMachineByLabel(ctx context.Context, c ctrlclient.Client, label map[string]string) (*machinev1alpha1.Machine, error) {
	obj := &machinev1alpha1.MachineList{}
	filter := &ctrlclient.ListOptions{
		LabelSelector: ctrlclient.MatchingLabelsSelector{Selector: labels.SelectorFromSet(label)},
	}
	if err := c.List(ctx, obj, filter); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, machinerr.NotFound(string(Machine))
	}
	return &obj.Items[0], nil
}

func getInventoryByLabel(ctx context.Context,
	c ctrlclient.Client, label map[string]string) (*inventoriesv1alpha1.Inventory, error) {
	obj := &inventoriesv1alpha1.InventoryList{}
	filter := &ctrlclient.ListOptions{
		LabelSelector: ctrlclient.MatchingLabelsSelector{Selector: labels.SelectorFromSet(label)},
	}
	if err := c.List(ctx, obj, filter); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, machinerr.NotFound(string(Inventory))
	}
	return &obj.Items[0], nil
}

func getSwitchByLabel(ctx context.Context, c ctrlclient.Client, label map[string]string) (*switchv1alpha1.Switch, error) {
	obj := &switchv1alpha1.SwitchList{}
	filter := &ctrlclient.ListOptions{
		LabelSelector: ctrlclient.MatchingLabelsSelector{Selector: labels.SelectorFromSet(label)},
	}
	if err := c.List(ctx, obj, filter); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, machinerr.NotFound(string(Switch))
	}
	return &obj.Items[0], nil
}

func getOOBByLabel(ctx context.Context, c ctrlclient.Client, label map[string]string) (*oobonmetal.Machine, error) {
	obj := &oobonmetal.MachineList{}
	filter := &ctrlclient.ListOptions{
		LabelSelector: ctrlclient.MatchingLabelsSelector{Selector: labels.SelectorFromSet(label)},
	}
	if err := c.List(ctx, obj, filter); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, machinerr.NotFound(string(Switch))
	}
	return &obj.Items[0], nil
}
