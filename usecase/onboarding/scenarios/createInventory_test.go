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

package scenarios_test

import (
	"testing"

	"github.com/onmetal/metal-api/common/types/base"
	"github.com/onmetal/metal-api/common/types/events"
	domain "github.com/onmetal/metal-api/domain/inventory"
	persistence "github.com/onmetal/metal-api/providers-kubernetes/onboarding"
	"github.com/onmetal/metal-api/providers-kubernetes/onboarding/fake"
	usecase "github.com/onmetal/metal-api/usecase/onboarding"
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	"github.com/onmetal/metal-api/usecase/onboarding/invariants"
	"github.com/onmetal/metal-api/usecase/onboarding/scenarios"
	"github.com/stretchr/testify/assert"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func newInventoryOnboardingUseCase(
	a *assert.Assertions,
	eventPublisher events.DomainEventPublisher,
	objs ...ctrlclient.Object) usecase.CreateInventory {
	fakeClient, err := fake.NewFakeWithObjects(objs...)
	a.Nil(err, "must create client")

	inventoryRepository := persistence.NewInventoryRepository(fakeClient, eventPublisher)
	inventoryAlreadyExist := invariants.NewInventoryAlreadyExist(inventoryRepository)
	inventoryIDGenerator := persistence.NewKubernetesInventoryIDGenerator()
	return scenarios.NewCreateInventoryUseCase(
		inventoryAlreadyExist,
		inventoryIDGenerator,
		inventoryRepository,
	)
}

func TestInventoryOnboardingUseCaseExecuteSuccess(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	request := dto.InventoryInfo{
		UUID:      "newtest",
		Namespace: "default",
	}
	eventPublihser := &fakePublisher{
		test:  a,
		event: &domain.InventoryCreatedDomainEvent{},
	}

	err := newInventoryOnboardingUseCase(a, eventPublihser).Execute(request)
	a.Nil(err, "must onboard inventory without error")
}

func TestInventoryOnboardingUseCaseExecuteAlreadyOnboarded(t *testing.T) {
	t.Parallel()

	name, namespace := "exist", "default"

	inv := fake.InventoryObject(name, namespace)
	a := assert.New(t)

	request := dto.InventoryInfo{
		UUID:      name,
		Namespace: namespace,
	}
	eventPublihser := &fakePublisher{
		test:  a,
		event: nil,
	}
	err := newInventoryOnboardingUseCase(a, eventPublihser, inv).Execute(request)
	a.True(usecase.IsAlreadyCreated(err))
}

type fakePublisher struct {
	test  *assert.Assertions
	event base.DomainEvent
}

func (f *fakePublisher) Publish(events ...base.DomainEvent) {
	f.test.IsType(f.event, events[0])
}
