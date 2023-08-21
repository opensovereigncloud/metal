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

package publisher_test

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/onmetal/metal-api/common/types/base"
	"github.com/onmetal/metal-api/common/types/events"
	"github.com/onmetal/metal-api/publisher"
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
