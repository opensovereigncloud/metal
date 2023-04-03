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

package scenarios

import (
	"github.com/onmetal/metal-api/internal/usecase/onboarding/access"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/rules"
)

type InventoryOnboardingUseCase struct {
	inventoryRepository access.InventoryRepository
	rule                rules.ServerMustBeEnabledOnFirstTime
}

func NewInventoryOnboardingUseCase(
	inventoryRepository access.InventoryRepository,
	rule rules.ServerMustBeEnabledOnFirstTime) *InventoryOnboardingUseCase {
	return &InventoryOnboardingUseCase{inventoryRepository: inventoryRepository, rule: rule}
}

func (o *InventoryOnboardingUseCase) Execute(request dto.Request) error {
	inv := dto.NewCreateInventory(request.Name, request.Namespace)
	if err := o.inventoryRepository.Create(inv); err != nil {
		return err
	}
	if err := o.rule.Execute(request); err != nil {
		return err
	}
	return nil
}
