package rules_test

import (
	"errors"
	"testing"

	domain "github.com/onmetal/metal-api/domain/inventory"
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	"github.com/onmetal/metal-api/usecase/onboarding/providers"
	"github.com/onmetal/metal-api/usecase/onboarding/providers/mocks"
	"github.com/onmetal/metal-api/usecase/onboarding/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestNewCreateIPv6SubnetFromParentForInventoryRule(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	subnetExtractorMock := mocks.NewLoopbackSubnetExtractor(t)
	subnetExtractorMock.
		On("ByType", providers.IPv6).
		Return(dto.SubnetInfo{Name: "test"}, nil)

	subnetPersisterMock := mocks.NewSubnetPersister(t)
	subnetPersisterMock.
		On("Save", mock.AnythingOfType("dto.SubnetInfo")).
		Return(nil)
	testInventoryID := domain.NewInventoryID("test")

	inventoryExtractorMock := mocks.NewInventoryExtractor(t)
	inventoryExtractorMock.
		On("ByID", testInventoryID).
		Return(domain.Inventory{UUID: "test"}, nil)

	rule := rules.NewCreateIPv6SubnetFromParentForInventoryRule(
		subnetExtractorMock,
		subnetPersisterMock,
		inventoryExtractorMock,
		zap.New())
	a.IsType(&domain.InventoryCreatedDomainEvent{}, rule.EventType())

	rule.Handle(domain.NewInventoryCreatedDomainEvent(testInventoryID))
}

func TestNewCreateIPv6SubnetFromParentForInventoryRuleNoParentSubnet(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	subnetExtractorMock := mocks.NewLoopbackSubnetExtractor(t)
	subnetExtractorMock.
		On("ByType", providers.IPv6).
		Return(dto.SubnetInfo{Name: "test"}, errors.New("subnet not found"))

	testInventoryID := domain.NewInventoryID("test")

	rule := rules.NewCreateIPv6SubnetFromParentForInventoryRule(
		subnetExtractorMock,
		nil,
		nil,
		zap.New())
	a.IsType(&domain.InventoryCreatedDomainEvent{}, rule.EventType())
	rule.Handle(domain.NewInventoryCreatedDomainEvent(testInventoryID))
}

func TestNewCreateIPv6SubnetFromParentForInventoryRuleNoInventoryFound(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	subnetExtractorMock := mocks.NewLoopbackSubnetExtractor(t)
	subnetExtractorMock.
		On("ByType", providers.IPv6).
		Return(dto.SubnetInfo{Name: "test"}, nil)

	testInventoryID := domain.NewInventoryID("test")

	inventoryExtractorMock := mocks.NewInventoryExtractor(t)
	inventoryExtractorMock.
		On("ByID", testInventoryID).
		Return(domain.Inventory{UUID: "test"}, errors.New("not found"))

	rule := rules.NewCreateIPv6SubnetFromParentForInventoryRule(
		subnetExtractorMock,
		nil,
		inventoryExtractorMock,
		zap.New())
	a.IsType(&domain.InventoryCreatedDomainEvent{}, rule.EventType())

	rule.Handle(domain.NewInventoryCreatedDomainEvent(testInventoryID))
}

func TestNewCreateIPv6SubnetFromParentForInventoryRuleSaveBroken(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	subnetExtractorMock := mocks.NewLoopbackSubnetExtractor(t)
	subnetExtractorMock.
		On("ByType", providers.IPv6).
		Return(dto.SubnetInfo{Name: "test"}, nil)

	testInventoryID := domain.NewInventoryID("test")

	inventoryExtractorMock := mocks.NewInventoryExtractor(t)
	inventoryExtractorMock.
		On("ByID", testInventoryID).
		Return(domain.Inventory{UUID: "test"}, nil)

	subnetPersisterMock := mocks.NewSubnetPersister(t)
	subnetPersisterMock.
		On("Save", mock.AnythingOfType("dto.SubnetInfo")).
		Return(errors.New("subnet save failed"))

	rule := rules.NewCreateIPv6SubnetFromParentForInventoryRule(
		subnetExtractorMock,
		subnetPersisterMock,
		inventoryExtractorMock,
		zap.New())
	a.IsType(&domain.InventoryCreatedDomainEvent{}, rule.EventType())

	rule.Handle(domain.NewInventoryCreatedDomainEvent(testInventoryID))
}