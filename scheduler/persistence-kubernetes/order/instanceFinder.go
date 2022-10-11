// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package order

import (
	"context"

	inventoryv1alpaha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	domain "github.com/onmetal/metal-api/scheduler/domain/order"
	usecase "github.com/onmetal/metal-api/scheduler/usecase/order"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type InstanceFinderExtractor struct {
	client ctrlclient.Client
}

func NewInstanceFinderExtractor(c ctrlclient.Client) *InstanceFinderExtractor {
	return &InstanceFinderExtractor{
		client: c,
	}
}

func (e *InstanceFinderExtractor) FindVacantInstanceForOrder(o domain.Order) (domain.OrderScheduler, error) {
	domainOrder := &machinev1alpha2.MachineAssignment{}
	if err := e.client.
		Get(
			context.Background(),
			types.NamespacedName{
				Namespace: o.Namespace(),
				Name:      o.Name(),
			},
			domainOrder); err != nil {
		return nil, err
	}
	o.SetInstanceType(domainOrder.Spec.MachineClass.Name)

	instance, err := e.findInstanceForType(o.InstanceType(), domainOrder.Spec.Tolerations)
	if err != nil {
		return nil, err
	}
	return newSchedulerExecutor(e.client, instance, domainOrder), nil
}

func (e *InstanceFinderExtractor) findInstanceForType(instanceType string,
	orderTolerations []machinev1alpha2.Toleration) (*machinev1alpha2.Machine, error) {
	instances, err := e.findInstancesForType(instanceType)
	if err != nil {
		return nil, err
	}
	instance := findToleratedMachine(orderTolerations, instances)
	if instance == nil {
		return nil, usecase.VacantInstanceNotFound(instanceType)
	}
	return instance, nil
}

func (e *InstanceFinderExtractor) findInstancesForType(instanceType string) (*machinev1alpha2.MachineList, error) {
	instanceList := &machinev1alpha2.MachineList{}
	instanceSelector := getLabelSelectorForAvailableMachine(instanceType)
	listOpts := &ctrlclient.ListOptions{
		LabelSelector: ctrlclient.MatchingLabelsSelector{Selector: labels.SelectorFromSet(instanceSelector)},
	}
	if err := e.client.
		List(
			context.Background(), instanceList, listOpts); err != nil {
		return nil, err
	}
	return instanceList, nil
}

func findToleratedMachine(
	orderTolerations []machinev1alpha2.Toleration,
	instanceList *machinev1alpha2.MachineList) *machinev1alpha2.Machine {
	for m := range instanceList.Items {
		if !IsTolerated(orderTolerations, instanceList.Items[m].Spec.Taints) {
			continue
		}
		if instanceList.Items[m].Status.Reservation.Reference.Name != "" {
			continue
		}
		if instanceList.Items[m].Status.Health != machinev1alpha2.MachineStateHealthy {
			continue
		}
		return &instanceList.Items[m]
	}
	return nil
}

func getLabelSelectorForAvailableMachine(instanceType string) map[string]string {
	sizeLabel := inventoryv1alpaha1.CLabelPrefix + instanceType
	return map[string]string{sizeLabel: "true"}
}

func IsTolerated(tolerations []machinev1alpha2.Toleration, taints []machinev1alpha2.Taint) bool {
	if len(taints) == 0 {
		return true
	}
	if len(tolerations) != len(taints) {
		return false
	}
	tolerated := 0
	for t := range tolerations {
		for taint := range taints {
			if !toleratesTaint(tolerations[t], taints[taint]) {
				continue
			}
			tolerated++
		}
	}
	return tolerated == len(taints)
}

func toleratesTaint(toleration machinev1alpha2.Toleration, taint machinev1alpha2.Taint) bool {
	if toleration.Effect != taint.Effect {
		return false
	}

	if toleration.Key != taint.Key {
		return false
	}

	switch toleration.Operator {
	case "", machinev1alpha2.TolerationOpEqual: // empty operator means Equal
		return toleration.Value == taint.Value
	case machinev1alpha2.TolerationOpExists:
		return true
	default:
		return false
	}
}
