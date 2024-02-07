// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	ipdomain "github.com/ironcore-dev/metal/domain/address"
)

type LoopbackAddressExtractor interface {
	Try(times int) LoopbackAddressExtractor
	IPv4ByMachineUUID(uuid string) (ipdomain.Address, error)
	IPv6ByMachineUUID(uuid string) (ipdomain.Address, error)
}
