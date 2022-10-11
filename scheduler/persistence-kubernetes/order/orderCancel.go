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
	"context"

	machinev1alpaha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	domain "github.com/onmetal/metal-api/scheduler/domain/order"
	"github.com/onmetal/metal-api/types/common"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type OrderCancelExecutor struct {
	client ctrlclient.Client
}

func NewOrderCancelExecutor(c ctrlclient.Client) *OrderCancelExecutor {
	return &OrderCancelExecutor{
		client: c,
	}
}

func (o *OrderCancelExecutor) Cancel(metadata common.Metadata) error {
	instance := &machinev1alpaha2.Machine{}
	if err := o.client.
		Get(
			context.Background(),
			types.NamespacedName{
				Namespace: metadata.Namespace(),
				Name:      metadata.Name(),
			},
			instance); err != nil {
		return err
	}
	instance.Status.Reservation.Status = domain.OrderStatusAvailable
	instance.Status.Reservation.Reference = common.ResourceReference{}

	return o.client.Update(context.Background(), instance)
}
