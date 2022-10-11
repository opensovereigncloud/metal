// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
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
// */

package base

type DomainEvent interface {
	Name() string
	Namespace() string
	EventType() EventType
}

type EventType string

const (
	ServerAccountingCreatedDomainEventType         EventType = "ServerAccountingCreatedDomainEvent"
	ServerAccountingNotAccomplishedDomainEventType EventType = "ServerAccountingNotAccomplishedDomainEventType"
)

type ServerAccountingCreatedDomainEvent struct {
	name, namespace string
}

type ServerAccountingNotAccomplishedDomainEvent struct {
	name, namespace string
}

func NewServerAccountingCreatedDomainEvent(name, namespace string) *ServerAccountingCreatedDomainEvent {
	return &ServerAccountingCreatedDomainEvent{
		name:      name,
		namespace: namespace,
	}
}

func (s *ServerAccountingCreatedDomainEvent) Name() string {
	return s.name
}

func (s *ServerAccountingCreatedDomainEvent) Namespace() string {
	return s.namespace
}

func (s *ServerAccountingCreatedDomainEvent) EventType() EventType {
	return ServerAccountingCreatedDomainEventType
}

func NewServerAccountingNotAccomplishedDomainEvent(name, namespace string) *ServerAccountingCreatedDomainEvent {
	return &ServerAccountingCreatedDomainEvent{
		name:      name,
		namespace: namespace,
	}
}

func (s *ServerAccountingNotAccomplishedDomainEvent) Name() string {
	return s.name
}

func (s *ServerAccountingNotAccomplishedDomainEvent) Namespace() string {
	return s.namespace
}

func (s *ServerAccountingNotAccomplishedDomainEvent) EventType() EventType {
	return ServerAccountingNotAccomplishedDomainEventType
}
