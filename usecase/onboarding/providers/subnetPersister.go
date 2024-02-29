// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import "github.com/ironcore-dev/metal/usecase/onboarding/dto"

//go:generate mockery --name SubnetPersister
type SubnetPersister interface {
	Save(info dto.SubnetInfo) error
}
