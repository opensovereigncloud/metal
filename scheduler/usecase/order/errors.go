// /*
// Copyright (c) 2022 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
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
// */

package usecase

import "errors"

const (
	notFound = "not found"
)

type SchedulerOrderError struct {
	Reason  string
	Message string
}

type SchedulerOrderStatus interface {
	Status() string
}

func ReasonForError(err error) string {
	if reason := SchedulerOrderStatus(nil); errors.As(err, &reason) {
		return reason.Status()
	}
	return ""
}

func (e *SchedulerOrderError) Error() string { return e.Message }

func IsVacantInstanceNotFound(err error) bool {
	return ReasonForError(err) == notFound
}
