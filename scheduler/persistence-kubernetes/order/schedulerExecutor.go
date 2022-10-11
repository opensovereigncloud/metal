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
	"context"

	machinev1alpaha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	domain "github.com/onmetal/metal-api/scheduler/domain/order"
	"github.com/onmetal/metal-api/types/common"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type SchedulerExecutor struct {
	client      ctrlclient.Client
	instance    *machinev1alpaha2.Machine
	domainOrder *machinev1alpaha2.MachineAssignment
}

func newSchedulerExecutor(client ctrlclient.Client,
	instance *machinev1alpaha2.Machine,
	domainOrder *machinev1alpaha2.MachineAssignment) *SchedulerExecutor {
	return &SchedulerExecutor{
		client:      client,
		instance:    instance,
		domainOrder: domainOrder,
	}
}

func (s *SchedulerExecutor) Schedule() error {
	s.domainOrder.Status.MachineRef = common.NewObjectMetadata(
		s.instance.Name,
		s.instance.Namespace).
		Reference()
	s.domainOrder.Status.State = domain.OrderStatusPending
	if err := s.client.
		Update(
			context.Background(),
			s.domainOrder); err != nil {
		return err
	}

	s.instance.Status.Reservation.Status = domain.OrderStatusPending
	s.instance.Status.Reservation.Reference = common.NewObjectMetadata(
		s.domainOrder.Name,
		s.domainOrder.Namespace).
		Reference()

	return s.client.
		Update(
			context.Background(),
			s.instance)
}
