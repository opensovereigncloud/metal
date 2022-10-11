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

	"github.com/go-logr/logr"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/types/common"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type OrderAlreadyScheduled struct {
	client ctrlclient.Client
	log    logr.Logger
}

func NewOrderAlreadyScheduled(c ctrlclient.Client, l logr.Logger) *OrderAlreadyScheduled {
	return &OrderAlreadyScheduled{
		client: c,
		log:    l,
	}
}

func (e *OrderAlreadyScheduled) Invoke(orderMeta common.Metadata) bool {
	domainOrder := &machinev1alpha2.MachineAssignment{}
	if err := e.client.
		Get(
			context.Background(),
			types.NamespacedName{
				Namespace: orderMeta.Namespace(),
				Name:      orderMeta.Name(),
			}, domainOrder); err != nil {
		e.log.V(1).Info("orderAlreadyScheduled extractor failed", "error", err)
		return true
	}
	return domainOrder.Status.MachineRef.Name != "" && domainOrder.Status.State != ""
}
