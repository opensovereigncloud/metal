// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package usecase

import (
	"fmt"

	"github.com/ironcore-dev/metal/common/types/errors"
	domain "github.com/ironcore-dev/metal/domain/machine"
)

type GetMachine interface {
	Execute(machineUUID string) (domain.Machine, error)
}

func MachineNotFound(err error) error {
	return errors.NewBusinessError(notFound, fmt.Sprintf("machine not found: %s", err.Error()))
}
