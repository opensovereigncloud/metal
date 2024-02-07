// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	domain "github.com/ironcore-dev/metal/domain/machine"
)

type MachineExtractor interface {
	ByID(machineID domain.MachineID) (domain.Machine, error)
	ByUUID(uuid string) (domain.Machine, error)
}
