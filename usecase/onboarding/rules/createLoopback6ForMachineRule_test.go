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

package rules_test

import (
	"testing"

	ipdomain "github.com/onmetal/metal-api/domain/address"
	domain "github.com/onmetal/metal-api/domain/inventory"
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	"github.com/onmetal/metal-api/usecase/onboarding/providers/mocks"
	"github.com/onmetal/metal-api/usecase/onboarding/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestNewCreateIPv6LoopbackForMachineRuleSuccess(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	inventoryUUID := "test"
	subnetExtractorMock := mocks.NewLoopbackSubnetExtractor(t)
	subnetExtractorMock.
		On("IPv6ByName", inventoryUUID).
		Return(dto.SubnetInfo{Name: "test"}, nil)

	addressPersisterMock := mocks.NewAddressPersister(t)
	addressPersisterMock.
		On("Save", mock.IsType(ipdomain.Address{})).
		Return(nil)
	testInventoryID := domain.NewInventoryID("test")

	inventoryExtractorMock := mocks.NewInventoryExtractor(t)
	inventoryExtractorMock.
		On("ByID", testInventoryID).
		Return(domain.Inventory{
			UUID: inventoryUUID,
			Sizes: map[string]string{
				"machine.onmetal.de/size-machine": "true",
			}}, nil)

	rule := rules.NewCreateLoopback6ForMachineRule(
		subnetExtractorMock,
		addressPersisterMock,
		inventoryExtractorMock,
		zap.New())
	a.IsType(&domain.InventoryFlavorUpdatedDomainEvent{}, rule.EventType())
	rule.Handle(domain.NewInventoryFlavorUpdatedDomainEvent(testInventoryID))
}
