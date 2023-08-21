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

	inventories "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	"github.com/onmetal/metal-api/common/types/base"
	domain "github.com/onmetal/metal-api/domain/inventory"
	providers "github.com/onmetal/metal-api/providers-kubernetes/onboarding"
	"github.com/onmetal/metal-api/providers-kubernetes/onboarding/fake"
	usecase "github.com/onmetal/metal-api/usecase/onboarding"
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	"github.com/onmetal/metal-api/usecase/onboarding/invariants"
	"github.com/onmetal/metal-api/usecase/onboarding/scenarios"
	"github.com/stretchr/testify/assert"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func addMachineUseCase(
	a *assert.Assertions,
	objects ...ctrlclient.Object,
) usecase.CreateMachine {
	fakeClient, err := fake.NewFakeWithObjects(objects...)
	a.Nil(err, "must create client")
	publisher := &fakeEventPublisher{}
	machineRepository := providers.NewMachineRepository(fakeClient, publisher)
	machineIDGenerator := providers.NewKubernetesMachineIDGenerator()
	machineAlreadyExist := invariants.NewMachineAlreadyExist(machineRepository)
	return scenarios.NewCreateMachineUseCase(
		machineRepository,
		machineIDGenerator,
		machineAlreadyExist,
	)
}

func TestAddMachineUseCaseExecuteSuccess(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	testInventory := inventory("test", "default")
	machineInfo := dto.NewMachineInfoFromInventory(testInventory)

	id, err := addMachineUseCase(a).Execute(machineInfo)
	a.NotEmpty(id)
	a.Nil(err)
}

func TestAddMachineUseCaseExecuteFailed(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	uuid, namespace := "test", "test"
	testInventory := inventory(uuid, namespace)
	machineInfo := dto.NewMachineInfoFromInventory(testInventory)
	machineObj := fake.FakeMachineObject(uuid, namespace)

	_, err := addMachineUseCase(a, machineObj).Execute(machineInfo)

	a.True(usecase.IsAlreadyCreated(err))
	a.NotNil(err, "must not create")
}

func inventory(uuid, namespace string) domain.Inventory {
	return domain.Inventory{
		UUID:         uuid,
		Namespace:    namespace,
		ProductSKU:   "1",
		SerialNumber: "1",
		Sizes: map[string]string{
			"machine.onmetal.de/size-m5.metal": "true",
			"machine.onmetal.de/size-machine":  "true",
		},
		NICs: []inventories.NICSpec{
			{
				Name:       "test",
				MACAddress: "123",
				MTU:        1500,
				LLDPs: []inventories.LLDPSpec{
					{
						ChassisID:         "test",
						SystemName:        "test",
						SystemDescription: "test",
						PortID:            "test",
						PortDescription:   "test",
					},
				},
				NDPs: nil,
			},
		},
	}
}

type fakeEventPublisher struct {
}

func (f fakeEventPublisher) Publish(_ ...base.DomainEvent) {

}
