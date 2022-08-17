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
	inventoryv1alpaha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/pkg/provider"
	domain "github.com/onmetal/metal-api/scheduler/domain/order"
	usecase "github.com/onmetal/metal-api/scheduler/usecase/order"
)

type InstanceFinderExtractor struct {
	client provider.Client
}

func NewInstanceFinderExtractor(c provider.Client) *InstanceFinderExtractor {
	return &InstanceFinderExtractor{
		client: c,
	}
}

func (e *InstanceFinderExtractor) FindVacantInstanceForOrder(o domain.Order) (domain.OrderScheduler, error) {
	domainOrder := &machinev1alpha2.MachineAssignment{}
	if err := e.client.Get(domainOrder, o); err != nil {
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
	metalList := &machinev1alpha2.MachineList{}
	instanceSelector := getLabelSelectorForAvailableMachine(instanceType)
	listOpts := &provider.ListOptions{
		Filter: instanceSelector,
	}
	if err := e.client.List(metalList, listOpts); err != nil {
		return nil, err
	}
	for m := range metalList.Items {
		if !IsTolerated(orderTolerations, metalList.Items[m].Spec.Taints) {
			continue
		}
		if metalList.Items[m].Status.Reservation.Reference != nil {
			continue
		}
		if metalList.Items[m].Status.Health != machinev1alpha2.MachineStateHealthy {
			continue
		}
		return &metalList.Items[m], nil
	}
	return nil, usecase.VacantInstanceNotFound(instanceType)
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
