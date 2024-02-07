// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import "errors"

var (
	errNotFound         = errors.New("not found")
	errIPNotSet         = errors.New("ip not set")
	errIPNotFound       = errors.New("ip not found")
	errSwitchIsNotReady = errors.New("switch is not ready")
)
