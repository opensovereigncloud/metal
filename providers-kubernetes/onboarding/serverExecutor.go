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

package providers

import (
	"github.com/go-logr/logr"
	domain "github.com/onmetal/metal-api/domain/infrastructure"
)

type FakeServerExecutor struct {
	log logr.Logger
}

func NewFakeServerExecutor(log logr.Logger) *FakeServerExecutor {
	return &FakeServerExecutor{log: log}
}

func (f *FakeServerExecutor) Enable(serverInfo domain.Server) error {
	f.log.Info("server turned on after inventory onboarding", "server", serverInfo)
	return nil
}
