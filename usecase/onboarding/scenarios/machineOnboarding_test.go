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
	"net/netip"
	"testing"

	switches "github.com/onmetal/metal-api/apis/switch/v1beta1"
	ipdomain "github.com/onmetal/metal-api/domain/address"
	invdomain "github.com/onmetal/metal-api/domain/inventory"
	domain "github.com/onmetal/metal-api/domain/machine"
	usecase "github.com/onmetal/metal-api/usecase/onboarding"
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	"github.com/onmetal/metal-api/usecase/onboarding/providers"
	"github.com/onmetal/metal-api/usecase/onboarding/scenarios"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	log = zap.New()
)

func newMachineOnboardingUseCase(
	a *assert.Assertions,
	fakeInventory invdomain.Inventory,
) usecase.MachineOnboarding {
	loopbackRepository := &fakeLoopbackRepository{
		address: netip.MustParsePrefix("0.0.0.1/32"),
		err:     nil,
	}
	machineRepository := &fakeMachineRepository{
		test:      a,
		inventory: fakeInventory,
	}
	swutchExtractor := &fakeSwitchRepository{}
	return scenarios.NewMachineOnboardingUseCase(
		machineRepository,
		machineRepository,
		swutchExtractor,
		loopbackRepository,
		log)
}

func TestMachineOnboardingUseCaseExecuteSuccess(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	testInventory := inventory("test", "default")
	err := newMachineOnboardingUseCase(a, testInventory).Execute(newMachine(), testInventory)
	a.Nil(err, "must onboard machine without error")
}

func newMachine() domain.Machine {
	return domain.NewMachine(
		domain.NewMachineID("test"),
		"test",
		"",
		65535,
		"sju",
		"serialNumber",
		nil,
		domain.Loopbacks{},
		map[string]string{"machine": "true"},
	)
}

type fakeMachineRepository struct {
	test      *assert.Assertions
	inventory invdomain.Inventory
}

func (f *fakeMachineRepository) Save(machine domain.Machine) error {
	return f.Update(machine)
}

func (f *fakeMachineRepository) Create(_ domain.Machine) error {
	return nil
}

func (f *fakeMachineRepository) Update(machine domain.Machine) error {
	f.test.Equal(len(machine.Interfaces), len(f.inventory.NICs))
	f.test.Equal(machine.Size, f.inventory.Sizes)
	f.test.Equal(machine.SKU, f.inventory.ProductSKU)

	return nil
}

func (f *fakeMachineRepository) ByUUID(_ string) (domain.Machine, error) {
	return domain.Machine{}, nil
}

func (f *fakeMachineRepository) ByID(_ domain.MachineID) (domain.Machine, error) {
	return domain.Machine{}, nil
}

type fakeLoopbackRepository struct {
	address netip.Prefix
	err     error
}

func (f *fakeLoopbackRepository) Try(_ int) providers.LoopbackAddressExtractor {
	return f
}
func (f *fakeLoopbackRepository) IPv4ByMachineUUID(uuid string) (ipdomain.Address, error) {
	return ipdomain.Address{Prefix: f.address}, f.err
}

func (f *fakeLoopbackRepository) IPv6ByMachineUUID(uuid string) (ipdomain.Address, error) {
	return ipdomain.Address{Prefix: f.address}, f.err
}

type fakeSwitchRepository struct {
}

func (f *fakeSwitchRepository) ByChassisID(id string) (dto.SwitchInfo, error) {
	return dto.SwitchInfo{
		Name:  "test",
		Lanes: 1,
		Interfaces: &switches.InterfaceSpec{
			IP: nil,
		},
	}, nil
}