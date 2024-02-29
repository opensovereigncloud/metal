// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package domain_test

import (
	"testing"

	domain "github.com/ironcore-dev/metal/domain/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestNewServerHasPowerCapsSuccess(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	server, err := domain.NewServer(
		"test",
		"test",
		[]string{domain.PowerCapabilities},
	)
	a.Nil(err)
	a.True(server.HasPowerCapabilities())
}

func TestNewServerHasPowerCapsFailed(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	server, err := domain.NewServer(
		"test",
		"test",
		[]string{"cmd"},
	)
	a.Nil(err)
	a.False(server.HasPowerCapabilities())
}

func TestNewServerEmptyUUID(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	server, err := domain.NewServer(
		"",
		"test",
		[]string{domain.PowerCapabilities},
	)
	a.NotNil(err)
	a.Empty(server)
}
