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

package order_test

import (
	"testing"

	"github.com/onmetal/metal-api/pkg/provider/kubernetes-provider/fake"
	domain "github.com/onmetal/metal-api/scheduler/domain/order"
	"github.com/onmetal/metal-api/scheduler/persistence-kubernetes/order"
	"github.com/stretchr/testify/assert"
)

func instanceExtractor(a *assert.Assertions) *order.InstanceFinderExtractor {
	k8sFakeClient, err := fake.NewFakeClient()
	a.Nil(err, "InstanceFinderExtractor: must crete fake client")

	return order.NewInstanceFinderExtractor(k8sFakeClient)
}

func TestInstanceExtractorFindVacantInstanceForOrder(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	extractor := instanceExtractor(a)
	a.NotNil(extractor, "InstanceFinderExtractor: must create instance extractor")

	o := domain.NewOrder(fake.ExistingOrderName, "default")
	orderExecutor, err := extractor.FindVacantInstanceForOrder(o)
	a.Nil(err, "InstanceFinderExtractor: must have not error")
	a.NotNil(orderExecutor, "InstanceFinderExtractor: order scheduler must not be nil")
}

func TestInstanceExtractorFindVacantInstanceForOrderNotFound(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	extractor := instanceExtractor(a)
	a.NotNil(extractor, "InstanceFinderExtractor: must create instance extractor")

	o := domain.NewOrder("fake", "default")
	orderExecutor, err := extractor.FindVacantInstanceForOrder(o)
	a.NotNil(err, "InstanceFinderExtractor: must have an error")
	a.Nil(orderExecutor, "InstanceFinderExtractor: order scheduler must not be created")
}
