// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scenarios_test

import (
	"testing"

	"github.com/ironcore-dev/metal/common/types/base"
	"github.com/ironcore-dev/metal/common/types/events"
	domain "github.com/ironcore-dev/metal/domain/inventory"
	persistence "github.com/ironcore-dev/metal/providers-kubernetes/onboarding"
	"github.com/ironcore-dev/metal/providers-kubernetes/onboarding/fake"
	usecase "github.com/ironcore-dev/metal/usecase/onboarding"
	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
	"github.com/ironcore-dev/metal/usecase/onboarding/invariants"
	"github.com/ironcore-dev/metal/usecase/onboarding/scenarios"
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
