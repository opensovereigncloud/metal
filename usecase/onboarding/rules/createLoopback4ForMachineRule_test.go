// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package rules_test

import (
	"testing"

	ipdomain "github.com/ironcore-dev/metal/domain/address"
	domain "github.com/ironcore-dev/metal/domain/inventory"
	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
	"github.com/ironcore-dev/metal/usecase/onboarding/providers"
	"github.com/ironcore-dev/metal/usecase/onboarding/providers/mocks"
	"github.com/ironcore-dev/metal/usecase/onboarding/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestNewCreateIPv4LoopbackForMachineRuleSuccess(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	subnetExtractorMock := mocks.NewLoopbackSubnetExtractor(t)
	subnetExtractorMock.
		On("ByType", providers.IPv4).
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
			UUID: "test",
			Sizes: map[string]string{
				"metal.ironcore.dev/size-machine": "true",
			}}, nil)

	rule := rules.NewCreateLoopback4ForMachineRule(
		subnetExtractorMock,
		addressPersisterMock,
		inventoryExtractorMock,
		zap.New())
	a.IsType(&domain.InventoryFlavorUpdatedDomainEvent{}, rule.EventType())
	rule.Handle(domain.NewInventoryFlavorUpdatedDomainEvent(testInventoryID))
}
