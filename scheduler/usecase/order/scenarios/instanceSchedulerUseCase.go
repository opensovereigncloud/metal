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
	"github.com/onmetal/metal-api/common/types/base"
	"github.com/onmetal/metal-api/scheduler/usecase/order/access"
	"github.com/onmetal/metal-api/scheduler/usecase/order/dto"
)

type InstanceSchedulerUseCase struct {
	serverExtractor access.ServerExtractor
}

func NewInstanceSchedulerUseCase(serverExtractor access.ServerExtractor) *InstanceSchedulerUseCase {
	return &InstanceSchedulerUseCase{
		serverExtractor: serverExtractor,
	}
}

func (i *InstanceSchedulerUseCase) Execute(instanceMeta base.Metadata) error {
	server, err := i.serverExtractor.Get(instanceMeta)
	if err != nil {
		return err
	}
	if !server.Reserved() {
		return CheckOut(server)
	}
	return CheckIn(server)
}

func CheckIn(server dto.Server) error {
	if err := server.SetPowerState("On"); err != nil {
		return err
	}
	return nil
}

func CheckOut(server dto.Server) error {
	if err := server.SetPowerState("Off"); err != nil {
		return err
	}
	return nil
}
