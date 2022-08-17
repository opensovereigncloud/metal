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
	"github.com/onmetal/metal-api/scheduler/usecase/order/access"
)

type FindVacantInstanceUseCase struct {
	instanceExtractor access.InstanceFinder
}

func NewFindVacantInstanceUseCase(instanceExtractor access.InstanceFinder) *FindVacantInstanceUseCase {
	return &FindVacantInstanceUseCase{
		instanceExtractor: instanceExtractor,
	}
}

func (f *FindVacantInstanceUseCase) Execute(order domain.Order) (domain.OrderScheduler, error) {
	orderScheduler, err := f.instanceExtractor.FindVacantInstanceForOrder(order)
	if err != nil {
		return nil, err
	}
	return orderScheduler, nil
}
