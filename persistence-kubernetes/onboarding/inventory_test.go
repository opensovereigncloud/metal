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

package persistence_test

import (
	"testing"

	persistence "github.com/onmetal/metal-api/persistence-kubernetes/onboarding"
	"github.com/onmetal/metal-api/persistence-kubernetes/onboarding/fake"
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	"github.com/stretchr/testify/assert"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func newInventoryRepository(a *assert.Assertions, obj ...ctrlclient.Object) *persistence.InventoryRepository {
	fakeClient, err := fake.NewFakeWithObjects(obj...)
	a.Nil(err, "must create client with object")
	return persistence.NewInventoryRepository(fakeClient)
}

func TestInventoryRepositoryGetSuccess(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	name, namespace := "test", "test"
	inv := fake.InventoryObject(name, namespace)
	repository := newInventoryRepository(a, inv)

	inventory, err := repository.Get(dto.Request{Name: name, Namespace: namespace})
	a.Nil(err, "must get inventory without error")
	a.Equal(inventory.UUID, name)
	a.Equal(inventory.Namespace, namespace)
}
