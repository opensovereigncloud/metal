// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scenarios_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/ironcore-dev/metal/common/types/base"
	domain "github.com/ironcore-dev/metal/domain/inventory"
	providers "github.com/ironcore-dev/metal/providers-kubernetes/onboarding"
	"github.com/ironcore-dev/metal/providers-kubernetes/onboarding/fake"
	usecase "github.com/ironcore-dev/metal/usecase/onboarding"
	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
	"github.com/ironcore-dev/metal/usecase/onboarding/invariants"
	"github.com/ironcore-dev/metal/usecase/onboarding/scenarios"
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
			"metal.ironcore.dev/size-m5.metal": "true",
			"metal.ironcore.dev/size-machine":  "true",
		},
		NICs: []metalv1alpha4.NICSpec{
			{
				Name:       "test",
				MACAddress: "123",
				MTU:        1500,
				LLDPs: []metalv1alpha4.LLDPSpec{
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
