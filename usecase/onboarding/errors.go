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
