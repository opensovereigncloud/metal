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

package order_test

import (
	"testing"

	machinev1alpha2 "github.com/onmetal/metal-api/apis/machine/v1alpha2"
	"github.com/onmetal/metal-api/scheduler/persistence-kubernetes/order"
)

func TestIsTolerated(t *testing.T) {
	t.Parallel()
	tolerations := []machinev1alpha2.Toleration{
		{
			Key:    "uuid",
			Value:  "123",
			Effect: machinev1alpha2.TaintEffectNotAvailable,
		},
	}
	taints := []machinev1alpha2.Taint{
		{
			Key:    "uuid",
			Value:  "123",
			Effect: machinev1alpha2.TaintEffectNotAvailable,
		},
	}

	if !order.IsTolerated(tolerations, taints) {
		t.Error("should be tolerated")
	}
}

func TestIsNotTolerated(t *testing.T) {
	t.Parallel()
	tolerations := []machinev1alpha2.Toleration{
		{
			Key:    "uuid",
			Value:  "123",
			Effect: machinev1alpha2.TaintEffectNotAvailable,
		},
	}
	taints := []machinev1alpha2.Taint{
		{
			Key:    "uuid",
			Value:  "123",
			Effect: machinev1alpha2.TaintEffectNotAvailable,
		},
		{
			Key:    "uuid",
			Value:  "123",
			Effect: machinev1alpha2.TaintEffectError,
		},
	}

	if order.IsTolerated(tolerations, taints) {
		t.Error("should not be tolerated")
	}
}

func TestIsNotToleratedEqual(t *testing.T) {
	t.Parallel()
	tolerations := []machinev1alpha2.Toleration{
		{
			Key:    "uuid",
			Value:  "123",
			Effect: machinev1alpha2.TaintEffectNotAvailable,
		},
		{
			Key:    "uuid",
			Value:  "1234",
			Effect: machinev1alpha2.TaintEffectSuspended,
		},
	}
	taints := []machinev1alpha2.Taint{
		{
			Key:    "uuid",
			Value:  "123",
			Effect: machinev1alpha2.TaintEffectNotAvailable,
		},
		{
			Key:    "uuid",
			Value:  "1234",
			Effect: machinev1alpha2.TaintEffectError,
		},
	}

	if order.IsTolerated(tolerations, taints) {
		t.Error("should not be tolerated")
	}
}
