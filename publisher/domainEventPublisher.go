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
