// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
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

package repository

import (
	"context"

	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/internal/entity"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type AssignmentSynchronizationRepo struct {
	client ctrlclient.Client
}

func NewAssignmentSynchronizationRepo(c ctrlclient.Client) *AssignmentSynchronizationRepo {
	return &AssignmentSynchronizationRepo{
		client: c,
	}
}

func (s *AssignmentSynchronizationRepo) Do(ctx context.Context, e entity.Synchronization) error {
	machine := &machinev1alpha2.Machine{}
	if err := s.client.Get(ctx,
		types.NamespacedName{Name: e.SourceName, Namespace: e.SourceNamespace}, machine); err != nil {
		return err
	}
	metalAssignment := &machinev1alpha2.MachineAssignment{}
	if err := s.client.Get(ctx,
		types.NamespacedName{Name: e.TargetName, Namespace: e.TargetNamespace}, metalAssignment); err != nil {
		return err
	}

	e.SourceStatus = machine.Status.Reservation.Status
	e.TargetStatus = metalAssignment.Status.State

	if !e.IsSyncNeeded() {
		return nil
	}
	metalAssignment.Status.State = e.SourceStatus

	return s.client.Status().Update(ctx, metalAssignment)
}
