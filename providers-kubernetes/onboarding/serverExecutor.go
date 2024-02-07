// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	"github.com/go-logr/logr"
	domain "github.com/ironcore-dev/metal/domain/infrastructure"
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
