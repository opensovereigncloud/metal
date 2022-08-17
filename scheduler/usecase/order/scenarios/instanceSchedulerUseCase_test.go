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

package scenarios

import (
	"testing"

	kubernetestpersictence "github.com/onmetal/metal-api/scheduler/persistence-kubernetes/order"
	usecase "github.com/onmetal/metal-api/scheduler/usecase/order"

	"github.com/onmetal/metal-api/common/types/base"
	"github.com/onmetal/metal-api/pkg/provider/kubernetes-provider/fake"
	"github.com/stretchr/testify/assert"
)

func instanceSchedule(a *assert.Assertions) usecase.InstanceSchedulerUseCase {
	k8sFakeClient, err := fake.NewFakeClient()
	a.Nil(err, "InstanceSchedulerUseCase: must create client")
	serverExtractor := kubernetestpersictence.NewServerExtractor(k8sFakeClient)
	return NewInstanceSchedulerUseCase(serverExtractor)
}

func TestNewInstanceSchedulerUseCaseCheckIn(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	schedule := instanceSchedule(a)
	instance := base.NewInstanceMetadata(
		fake.ExistingServerUUID,
		"default")
	a.Nil(schedule.Execute(instance), "InstanceSchedulerUseCase: muse schedule instance without error")
}

func TestNewInstanceSchedulerUseCaseCheckOut(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	schedule := instanceSchedule(a)
	instance := base.NewInstanceMetadata(
		fake.AvailiableServerUUID,
		"default")
	a.Nil(schedule.Execute(instance), "InstanceSchedulerUseCase: muse schedule instance without error")
}
