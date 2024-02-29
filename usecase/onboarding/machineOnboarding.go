// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package usecase

import (
	"fmt"

	"github.com/ironcore-dev/metal/common/types/errors"
	invdomain "github.com/ironcore-dev/metal/domain/inventory"
	domain "github.com/ironcore-dev/metal/domain/machine"
)

type MachineOnboarding interface {
	Execute(machine domain.Machine, inventory invdomain.Inventory) error
}

func MachineAlreadCreated(name string) error {
	return errors.NewBusinessError(
		alreadyCreated,
		fmt.Sprintf("Machine Already Onboarded: %s", name),
	)
}
