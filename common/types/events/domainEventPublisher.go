// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package events

import "github.com/ironcore-dev/metal/common/types/base"

type DomainEventPublisher interface {
	Publish(events ...base.DomainEvent)
}
