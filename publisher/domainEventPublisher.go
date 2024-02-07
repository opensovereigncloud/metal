// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package publisher

import (
	"github.com/go-logr/logr"
	"github.com/ironcore-dev/metal/common/types/base"
	"github.com/ironcore-dev/metal/common/types/events"
)

type DomainEventPublisher struct {
	listeners map[string][]events.DomainEventListener[base.DomainEvent]
	log       logr.Logger
}

func NewDomainEventPublisher(
	log logr.Logger,
) *DomainEventPublisher {
	listeners := make(map[string][]events.DomainEventListener[base.DomainEvent])
	return &DomainEventPublisher{
		listeners: listeners,
		log:       log,
	}
}

func (d *DomainEventPublisher) RegisterListeners(
	domainEventListeners ...events.DomainEventListener[base.DomainEvent],
) {
	for _, domainEventListener := range domainEventListeners {
		domainEvent := domainEventListener.EventType()
		d.listeners[domainEvent.Type()] = append(d.listeners[domainEvent.Type()], domainEventListener)
	}
}

func (d *DomainEventPublisher) Publish(events ...base.DomainEvent) {
	for _, event := range events {
		listener, ok := d.listeners[event.Type()]
		if !ok {
			d.log.Info("listener for event not found", "event", event.Type())
			continue
		}
		d.log.Info("event published", "id", event.ID(), "event", event.Type())
		d.sendEvent(listener, event)
	}
}

func (d *DomainEventPublisher) sendEvent(
	listener []events.DomainEventListener[base.DomainEvent],
	event base.DomainEvent,
) {
	for l := range listener {
		listener[l].Handle(event)
	}
}
