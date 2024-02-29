// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	domain "github.com/ironcore-dev/metal/domain/machine"
)

type MachinePersister interface {
	Save(machine domain.Machine) error
	Create(machine domain.Machine) error
	Update(machine domain.Machine) error
}
