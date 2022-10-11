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

	"github.com/go-logr/logr"
	machinev1alpaha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/pkg/provider/kubernetes-provider/fake"
	"github.com/onmetal/metal-api/scheduler/controllers"
	"github.com/onmetal/metal-api/scheduler/persistence-kubernetes/order"
	"github.com/onmetal/metal-api/scheduler/usecase/order/scenarios"
	"github.com/onmetal/metal-api/types/common"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func newLogger() logr.Logger {
	opts := zap.Options{
		Development: true,
		DestWriter:  nil,
	}
	return zap.New(zap.UseFlagOptions(&opts))
}

func orderScheduler(a *assert.Assertions) *controllers.Scheduler {
	l := newLogger()
	fakeClient, err := fake.NewFakeClient()
	a.Nil(err, "must create client")
	orderAlreadyScheduled := order.NewOrderAlreadyScheduled(fakeClient, l)
	instanceExtractor := order.NewInstanceFinderExtractor(fakeClient)
	instanceForOrderUseCase := scenarios.NewFindVacantInstanceUseCase(instanceExtractor)
	cancelOrder := order.NewOrderCancelExecutor(fakeClient)

	return controllers.NewSchedulerController(
		l,
		orderAlreadyScheduled,
		cancelOrder,
		instanceForOrderUseCase,
	)
}

func BenchmarkSchedulerReconcile(t *testing.B) {
	t.ReportAllocs()
	a := assert.New(t)

	scheduler := orderScheduler(a)
	a.NotNil(scheduler, "Scheduler: must not be nil")
	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      fake.ExistingOrderName,
			Namespace: "default",
		},
	}
	ctx := context.Background()
	for i := 0; i < t.N; i++ {
		_, err := scheduler.Reconcile(ctx, req)
		a.Nil(err, "Scheduler: must reconcile without error")
	}
}

func TestSchedulerReconcile(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	scheduler := orderScheduler(a)
	a.NotNil(scheduler, "Scheduler: must not be nil")
	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      fake.ExistingOrderName,
			Namespace: "default",
		},
	}
	_, err := scheduler.Reconcile(context.Background(), req)
	a.Nil(err, "Scheduler: must reconcile without error")
}

func TestSchedulerIsAlreadyScheduled(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	scheduler := orderScheduler(a)
	a.NotNil(scheduler, "Scheduler: must not be nil")
	updateEvent := event.UpdateEvent{
		ObjectNew: orderForController(),
	}
	a.False(scheduler.AlreadyOrdered(updateEvent), "Scheduler: must be already proceeded")
}

func TestSchedulerCancelOrder(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	scheduler := orderScheduler(a)
	a.NotNil(scheduler, "Scheduler: must not be nil")
	deleteEvent := event.DeleteEvent{
		Object:             orderForController(),
		DeleteStateUnknown: false,
	}
	a.False(scheduler.CancelOrder(deleteEvent), "Scheduler: must delete successfully")
}

func orderForController() *machinev1alpaha2.MachineAssignment {
	return &machinev1alpaha2.MachineAssignment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fake.ExistingOrderName,
			Namespace: "default",
		},
		Spec: machinev1alpaha2.MachineAssignmentSpec{
			MachineClass: corev1.LocalObjectReference{
				Name: "m5-metal",
			},
			Image: "myimage_repo_location",
		},
		Status: machinev1alpaha2.MachineAssignmentStatus{
			State: "Running",
			MachineRef: common.ResourceReference{
				Name:      fake.ExistingServerUUID,
				Namespace: "default",
			},
		},
	}
}
