// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package usecase

import "github.com/ironcore-dev/metal/common/types/errors"

const (
	alreadyCreated = "already onboarded"
	notFound       = "not found"
)

func IsAlreadyCreated(err error) bool {
	return errors.ReasonForError(err) == alreadyCreated
}

func IsNotFound(err error) bool {
	return errors.ReasonForError(err) == notFound
}
