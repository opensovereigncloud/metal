// Copyright 2023 OnMetal authors
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

package scenarios_test

import (
	"testing"

	domain "github.com/onmetal/metal-api/domain/machine"
	persistence "github.com/onmetal/metal-api/persistence-kubernetes/onboarding"
	"github.com/onmetal/metal-api/persistence-kubernetes/onboarding/fake"
	usecase "github.com/onmetal/metal-api/usecase/onboarding"
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	"github.com/onmetal/metal-api/usecase/onboarding/scenarios"
	"github.com/stretchr/testify/assert"
)

func newMachineOnboardingUseCase(a *assert.Assertions,
	fakeInventory dto.Inventory) usecase.MachineOnboardingUseCase {
	fakeClient, err := fake.NewFakeClient()
	a.Nil(err, "must create client")

	machineNetwork := persistence.NewMachineInterfaces(fakeClient)
	machineRepository := &fakeMachineRepository{
		test:      a,
		inventory: fakeInventory,
	}
	return scenarios.NewMachineOnboardingUseCase(
		machineNetwork,
		machineRepository)
}

func TestMachineOnboardingUseCaseExecuteSuccess(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	testInventory := inventory("test", "default")
	err := newMachineOnboardingUseCase(a, testInventory).Execute(machine(), testInventory)
	a.Nil(err, "must onboard machine without error")
}

func machine() domain.Machine {
	return domain.Machine{
		UUID:         "test",
		Namespace:    "default",
		SKU:          "",
		SerialNumber: "",
		Interfaces:   nil,
		Size:         nil,
	}
}

type fakeMachineRepository struct {
	test      *assert.Assertions
	inventory dto.Inventory
}

func (f *fakeMachineRepository) Create(inventory dto.Inventory) error {
	return nil
}

func (f *fakeMachineRepository) Update(machine domain.Machine) error {
	f.test.Equal(len(machine.Interfaces), len(f.inventory.NICs))
	f.test.Equal(machine.Size, f.inventory.Sizes)
	f.test.Equal(machine.SKU, f.inventory.ProductSKU)

	return nil
}

func (f *fakeMachineRepository) Get(request dto.Request) (domain.Machine, error) {
	return domain.Machine{}, nil
}
