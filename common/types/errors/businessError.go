// Copyright 2023 OnMetal authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
