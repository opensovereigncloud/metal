// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers_test

import (
	"testing"

	"github.com/ironcore-dev/metal/common/types/base"
	"github.com/ironcore-dev/metal/common/types/events"
	persistence "github.com/ironcore-dev/metal/providers-kubernetes/onboarding"
	"github.com/ironcore-dev/metal/providers-kubernetes/onboarding/fake"
	"github.com/stretchr/testify/assert"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func newInventoryRepository(
	a *assert.Assertions,
	publisher events.DomainEventPublisher,
	obj ...ctrlclient.Object,
) *persistence.InventoryRepository {
	fakeClient, err := fake.NewFakeWithObjects(obj...)
	a.Nil(err, "must create client with object")
	return persistence.NewInventoryRepository(
		fakeClient,
		publisher,
	)
}

func TestInventoryRepositoryGetSuccess(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	name, namespace := "test", "test"
	inv := fake.InventoryObject(name, namespace)
	publisher := &fakePublisher{}

	repository := newInventoryRepository(a, publisher, inv)

	inventory, err := repository.ByUUID(name)
	a.Nil(err, "must get inventory without error")
	a.Equal(inventory.UUID, name)
	a.Equal(inventory.Namespace, namespace)
}

type fakePublisher struct {
}

func (f *fakePublisher) Publish(events ...base.DomainEvent) {

}
