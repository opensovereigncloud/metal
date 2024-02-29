// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package usecase

import (
	"fmt"

	"github.com/ironcore-dev/metal/common/types/errors"
	domain "github.com/ironcore-dev/metal/domain/infrastructure"
)

type GetServer interface {
	Execute(uuid string) (domain.Server, error)
}

func ServerNotFound(err error) error {
	return errors.NewBusinessError(notFound, fmt.Sprintf("server not found: %s", err.Error()))
}
