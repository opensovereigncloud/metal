// /*
// Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
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
	"github.com/onmetal/metal-api/common/types/base"
	"github.com/onmetal/metal-api/pkg/provider"
	domain "github.com/onmetal/metal-api/scheduler/domain/order"
)

type OrderCancelExecutor struct {
	client provider.Client
}

func NewOrderCancelExecutor(c provider.Client) *OrderCancelExecutor {
	return &OrderCancelExecutor{
		client: c,
	}
}

func (o *OrderCancelExecutor) Cancel(metadata base.Metadata) error {
	instance := &machinev1alpaha2.Machine{}
	if err := o.client.Get(instance, metadata); err != nil {
		return err
	}

	instance.Status.Reservation.Status = domain.OrderStatusAvailable
	instance.Status.Reservation.Reference = nil

	return o.client.Update(instance)
}
