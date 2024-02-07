// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import "github.com/ironcore-dev/metal/usecase/onboarding/dto"

type SwitchExtractor interface {
	ByChassisID(chassisID string) (dto.SwitchInfo, error)
}
