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
	"context"

	domain "github.com/onmetal/metal-api/domain/infrastructure"
	"github.com/onmetal/metal-api/usecase/onboarding/dto"
	oob "github.com/onmetal/oob-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ServerRepository struct {
	client ctrlclient.Client
}

func NewServerRepository(
	client ctrlclient.Client,
) *ServerRepository {
	return &ServerRepository{client: client}
}

func (s *ServerRepository) Get(
	request dto.Request,
) (domain.Server, error) {
	oobServer, err := s.getOOB(request)
	if err != nil {
		return domain.Server{}, err
	}
	return domain.NewServer(
		oobServer.Name,
		oobServer.Namespace,
		oobServer.Status.Capabilities,
	)
}

func (s *ServerRepository) ByUUID(
	uuid string,
) (domain.Server, error) {
	uuidOptions := ctrlclient.MatchingFields{
		"metadata.name": uuid,
	}
	oobServer, err := s.extractOOBFromCluster(uuidOptions)
	if err != nil {
		return domain.Server{}, err
	}
	return domain.NewServer(
		oobServer.Name,
		oobServer.Namespace,
		oobServer.Status.Capabilities,
	)
}

func (s *ServerRepository) getOOB(
	request dto.Request,
) (*oob.OOB, error) {
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

func (s *ServerRepository) extractOOBFromCluster(options ctrlclient.ListOption) (*oob.OOB, error) {
	obj := &oob.OOBList{}
	if err := s.
		client.
		List(
			context.Background(),
			obj,
			options,
		); err != nil {
		return nil, err
	}
	if len(obj.Items) == 0 {
		return nil, errNotFound
	}
	return &obj.Items[0], nil
}
