// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package base

type DomainEntity interface {
	AddEvent(event DomainEvent)
	PopEvents() []DomainEvent
}

type DomainEntityImpl struct {
	events []DomainEvent
}

func NewDomainEntity() *DomainEntityImpl {
	events := make([]DomainEvent, 0)
	return &DomainEntityImpl{events: events}
}

func (d *DomainEntityImpl) AddEvent(event DomainEvent) {
	d.events = append(d.events, event)
}

func (d *DomainEntityImpl) PopEvents() []DomainEvent {
	return d.events
}
