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
	"github.com/go-logr/logr"
	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/pkg/provider"
	domain "github.com/onmetal/metal-api/scheduler/domain/order"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

type OrderExistExtractor struct {
	client provider.Client
	log    logr.Logger
}

func NewOrderExistExtractor(c provider.Client, l logr.Logger) *OrderExistExtractor {
	return &OrderExistExtractor{
		client: c,
		log:    l,
	}
}

func (e *OrderExistExtractor) Invoke(order domain.Order) bool {
	domainOrder := &machinev1alpha2.MachineAssignment{}
	if err := e.client.Get(domainOrder, order); err != nil {
		if apierrors.IsNotFound(err) {
			return false
		}
		e.log.V(1).Info("orderAlreadyScheduled extractor failed", "error", err)
		return true
	}
	return true
}
