// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package usecase

import (
	domain "github.com/ironcore-dev/metal/domain/machine"
	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
)

type CreateMachine interface {
	Execute(machineInfo dto.MachineInfo) (domain.MachineID, error)
}
