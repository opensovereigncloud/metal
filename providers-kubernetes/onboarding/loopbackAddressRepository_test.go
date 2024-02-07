// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers_test

import (
	"fmt"
	"testing"

	providers "github.com/ironcore-dev/metal/providers-kubernetes/onboarding"
	"github.com/ironcore-dev/metal/providers-kubernetes/onboarding/fake"
	"github.com/stretchr/testify/assert"
)

func TestLoopbackRepositoryIPv4ByMachineUUIDSuccess(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	uuid := "123"
	ip := fake.IPIPAMObject(fmt.Sprintf("%s-lo-ipv4", uuid), "test")
	fakeClient, _ := fake.NewFakeWithObjects(ip)

	repo := providers.NewLoopbackAddressRepository(fakeClient)

	address, err := repo.IPv4ByMachineUUID(uuid)
	a.Nil(err)
	a.NotEmpty(address.Prefix)
}
