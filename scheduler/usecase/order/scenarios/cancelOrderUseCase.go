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

package scenarios

import (
	domain "github.com/onmetal/metal-api/scheduler/domain/order"
	"github.com/onmetal/metal-api/types/common"
)

type CancelOrderUseCase struct {
	orderCancel domain.CancelOrder
}

func NewCancelOrderUseCase(orderCancel domain.CancelOrder) *CancelOrderUseCase {
	return &CancelOrderUseCase{
		orderCancel: orderCancel,
	}
}

func (s *CancelOrderUseCase) Cancel(instanceMetadata common.Metadata) error {
	return s.orderCancel.Cancel(instanceMetadata)
}
