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

package usecase

import (
	"context"

	"github.com/onmetal/metal-api/internal/entity"
	"github.com/onmetal/metal-api/internal/repository"
)

type SyncUseCase struct {
	repo repository.Synchronization
}

func NewSyncUseCase(r repository.Synchronization) *SyncUseCase {
	return &SyncUseCase{repo: r}
}

// Do - synchronize state of two objects.
func (s *SyncUseCase) Do(ctx context.Context, e entity.Reservation) error {
	sync := entity.Synchronization{
		SourceName:      e.RequestName,
		SourceNamespace: e.RequestNamespace,
		TargetName:      e.OrderName,
		TargetNamespace: e.OrderNamespace,
	}
	return s.repo.Do(ctx, sync)
}