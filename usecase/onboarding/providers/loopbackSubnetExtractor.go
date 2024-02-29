// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import "github.com/ironcore-dev/metal/usecase/onboarding/dto"

const (
	IPv4 = "IPv4"
	IPv6 = "IPv6"
)

//go:generate mockery --name LoopbackSubnetExtractor
type LoopbackSubnetExtractor interface {
	ByType(ipType string) (dto.SubnetInfo, error)
	IPv6ByName(name string) (dto.SubnetInfo, error)
}
