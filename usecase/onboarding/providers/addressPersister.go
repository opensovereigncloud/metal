// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	domain "github.com/ironcore-dev/metal/domain/address"
)

//go:generate mockery --name AddressPersister
type AddressPersister interface {
	Save(address domain.Address) error
}
