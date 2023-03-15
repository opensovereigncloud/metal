package scenarios_test

import (
	"testing"

	inventories "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	persistence "github.com/onmetal/metal-api/internal/kubernetes/onboarding"
	"github.com/onmetal/metal-api/internal/kubernetes/onboarding/fake"
	usecase "github.com/onmetal/metal-api/internal/usecase/onboarding"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/scenarios"
	"github.com/stretchr/testify/assert"
)

func addMachineUseCase(a *assert.Assertions) usecase.AddMachineUseCase {
	fakeClient, err := fake.NewFakeClient()
	a.Nil(err, "must create client")
	repository := persistence.NewMachineRepository(fakeClient)

	return scenarios.NewAddMachineUseCase(repository)
}

func TestAddMachineUseCaseExecuteSuccess(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	testInventory := inventory("test", "default")

	err := addMachineUseCase(a).Execute(testInventory)
	a.Nil(err, "must create without error")
}

func TestAddMachineUseCaseExecuteFailed(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	testInventory := inventory("", "")

	err := addMachineUseCase(a).Execute(testInventory)
	a.NotNil(err, "must not create")
}

func inventory(uuid, namespace string) dto.Inventory {
	return dto.Inventory{
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
				LLDPs:      nil,
				NDPs:       nil,
			},
		},
	}
}
