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

package persistence

import (
	"context"

	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	oob "github.com/onmetal/oob-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ServerRepository struct {
	client ctrlclient.Client
}

func NewServerRepository(client ctrlclient.Client) *ServerRepository {
	return &ServerRepository{client: client}
}

func (s *ServerRepository) Get(request dto.Request) (dto.Server, error) {
	oobServer, err := s.getOOB(request)
	if err != nil {
		return dto.Server{}, err
	}
	return dto.Server{
		UUID:              oobServer.Status.UUID,
		PowerCapabilities: oobServer.Status.Capabilities,
	}, nil
}

func (s *ServerRepository) getOOB(request dto.Request) (*oob.OOB, error) {
	oobData := &oob.OOB{}
	err := s.
		client.
		Get(
			context.Background(),
			types.NamespacedName{
				Namespace: request.Namespace,
				Name:      request.Name,
			},
			oobData)
	return oobData, err
}
