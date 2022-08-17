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

package controllers_test

import (
	"context"
	"testing"

	"github.com/onmetal/metal-api/pkg/provider/kubernetes-provider/fake"
	"github.com/onmetal/metal-api/scheduler/controllers"
	"github.com/onmetal/metal-api/scheduler/persistence-kubernetes/order"
	"github.com/onmetal/metal-api/scheduler/usecase/order/scenarios"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func orderCleaner(a *assert.Assertions) *controllers.InstanceCleaner {
	l := newLogger()
	fakeClient, err := fake.NewFakeClient()
	a.Nil(err, "must create client")
	instanceExtractor := order.NewInstanceExtractor(fakeClient)
	orderExistExtractor := order.NewOrderExistExtractor(fakeClient, l)
	instanceCleanerUseCase := scenarios.NewOrderCleanerUseCase(instanceExtractor, orderExistExtractor)

	return controllers.NewInstanceCleaner(
		l,
		instanceCleanerUseCase,
	)
}

func TestCleanerReconcile(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	scheduler := orderCleaner(a)
	a.NotNil(scheduler, "Cleaner: must not be nil")
	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      fake.ExistingServerUUID,
			Namespace: "default",
		},
	}
	_, err := scheduler.Reconcile(context.Background(), req)
	a.Nil(err, "Cleaner: must reconcile without error")
}
