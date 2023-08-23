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

package providers_test

import (
	"testing"

	"github.com/onmetal/metal-api/common/types/base"
	"github.com/onmetal/metal-api/common/types/events"
	persistence "github.com/onmetal/metal-api/providers-kubernetes/onboarding"
	"github.com/onmetal/metal-api/providers-kubernetes/onboarding/fake"
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
