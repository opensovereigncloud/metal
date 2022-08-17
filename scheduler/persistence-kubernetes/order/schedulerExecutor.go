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

package order

import (
	machinev1alpaha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/pkg/provider"
	domain "github.com/onmetal/metal-api/scheduler/domain/order"
)

type SchedulerExecutor struct {
	client      provider.Client
	instance    *machinev1alpaha2.Machine
	domainOrder *machinev1alpaha2.MachineAssignment
}

func newSchedulerExecutor(client provider.Client,
	instance *machinev1alpaha2.Machine,
	domainOrder *machinev1alpaha2.MachineAssignment) *SchedulerExecutor {
	return &SchedulerExecutor{
		client:      client,
		instance:    instance,
		domainOrder: domainOrder,
	}
}

func (s *SchedulerExecutor) Execute() error {
	s.domainOrder.Status.MachineRef = orderReference(s.instance.Name, s.instance.Namespace)
	s.domainOrder.Status.State = domain.OrderStatusPending
	if err := s.client.Update(s.domainOrder); err != nil {
		return err
	}

	s.instance.Status.Reservation.Status = domain.OrderStatusPending
	s.instance.Status.Reservation.Reference = orderReference(s.domainOrder.Name, s.domainOrder.Namespace)

	return s.client.Update(s.instance)
}

func orderReference(requestName, namespace string) *machinev1alpaha2.ResourceReference {
	return &machinev1alpaha2.ResourceReference{
		Name: requestName, Namespace: namespace,
	}
}
