package persistence_test

import (
	"testing"

	persistence "github.com/onmetal/metal-api/internal/kubernetes/onboarding"
	"github.com/onmetal/metal-api/internal/kubernetes/onboarding/fake"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
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
