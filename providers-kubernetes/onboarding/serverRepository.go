// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	"context"

	domain "github.com/ironcore-dev/metal/domain/infrastructure"
	"github.com/ironcore-dev/metal/usecase/onboarding/dto"
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
