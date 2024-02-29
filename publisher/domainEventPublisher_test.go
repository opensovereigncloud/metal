// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package publisher_test

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/ironcore-dev/metal/common/types/base"
	"github.com/ironcore-dev/metal/common/types/events"
	"github.com/ironcore-dev/metal/publisher"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	log = zap.New(zap.UseFlagOptions(&zap.Options{Development: true}))
)

func TestNewDomainEventPublisher(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	newPublisher := publisher.NewDomainEventPublisher(log)

	listener := newFakeDomainEventListener(log, a)
	anotherListener := newFakeAnotherDomainEventListener(log, a)

	newPublisher.RegisterListeners(listener, anotherListener)

	event := &FakeDomainEvent{id: "test"}
	newPublisher.Publish(event)
}

type FakeDomainEventListener struct {
	test *assert.Assertions
	log  logr.Logger
}

func newFakeDomainEventListener(
	log logr.Logger,
	test *assert.Assertions,
) events.DomainEventListener[base.DomainEvent] {
	return &FakeDomainEventListener{
		log:  log,
		test: test,
	}
}

func (c *FakeDomainEventListener) EventType() base.DomainEvent {
	return &FakeDomainEvent{}
}

func (c *FakeDomainEventListener) Handle(event base.DomainEvent) {
	c.test.NotEmpty(event)
	c.test.Equal("test", event.ID())
	c.log.Info("success")
}

type FakeAnotherDomainEventListener struct {
	test *assert.Assertions
	log  logr.Logger
}

func newFakeAnotherDomainEventListener(
	log logr.Logger,
	test *assert.Assertions,
) events.DomainEventListener[base.DomainEvent] {
	return &FakeAnotherDomainEventListener{
		log:  log,
		test: test,
	}
}

func (c *FakeAnotherDomainEventListener) EventType() base.DomainEvent {
	return &FakeAnotherDomainEvent{}
}

func (c *FakeAnotherDomainEventListener) Handle(event base.DomainEvent) {
	c.test.Empty(event)
}

type FakeDomainEvent struct {
	id string
}

func (m *FakeDomainEvent) ID() string {
	return m.id
}

func (m *FakeDomainEvent) Type() string {
	return "fake_event_created"
}

type FakeAnotherDomainEvent struct {
	id string
}

func (m *FakeAnotherDomainEvent) ID() string {
	return m.id
}

func (m *FakeAnotherDomainEvent) Type() string {
	return "fake_another_event_created"
}
