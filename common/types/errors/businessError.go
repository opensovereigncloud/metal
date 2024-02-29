// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package errors

import "github.com/pkg/errors"

const unknown = "unknown"

type BusinessError interface {
	Error() string
	Status() string
}

type BusinessErrorImpl struct {
	Reason  string
	Message string
}

func NewBusinessError(
	reason string,
	message string,
) *BusinessErrorImpl {
	return &BusinessErrorImpl{
		Reason:  reason,
		Message: message,
	}
}

type BusinessErrorStatus interface {
	Status() string
}

func ReasonForError(err error) string {
	if reason := BusinessErrorStatus(nil); errors.As(err, &reason) {
		return reason.Status()
	}
	return unknown
}

func (e *BusinessErrorImpl) Error() string { return e.Message }

func (e *BusinessErrorImpl) Status() string { return e.Reason }
