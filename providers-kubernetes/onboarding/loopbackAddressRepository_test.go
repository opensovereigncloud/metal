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
