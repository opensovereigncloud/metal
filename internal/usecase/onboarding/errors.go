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

package usecase

import "errors"

const (
	alreadyOnboarded = "already onboarded"
	notValidated     = "not validated"
	notFound         = "not found"
	unknown          = "unknown"
)

type OnboardingError struct {
	Reason  string
	Message string
}

type OnboardingStatus interface {
	Status() string
}

func ReasonForError(err error) string {
	if reason := OnboardingStatus(nil); errors.As(err, &reason) {
		return reason.Status()
	}
	return unknown
}

func (e *OnboardingError) Error() string { return e.Message }

func (e *OnboardingError) Status() string { return e.Reason }

func IsAlreadyOnboarded(err error) bool {
	return ReasonForError(err) == alreadyOnboarded
}

func IsNotFound(err error) bool {
	return ReasonForError(err) == notFound
}
