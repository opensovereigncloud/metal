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

package scenarios_test

import (
	"errors"
	"testing"

	persistence "github.com/onmetal/metal-api/internal/kubernetes/onboarding"
	"github.com/onmetal/metal-api/internal/kubernetes/onboarding/fake"
	usecase "github.com/onmetal/metal-api/internal/usecase/onboarding"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/dto"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/rules"
	"github.com/onmetal/metal-api/internal/usecase/onboarding/scenarios"
	"github.com/stretchr/testify/assert"
)

func newInventoryOnboardingUseCase(a *assert.Assertions,
	rule rules.ServerMustBeEnabledOnFirstTime) usecase.OnboardingUseCase {
	fakeClient, err := fake.NewFakeClient()
	a.Nil(err, "must create client")

	inventoryRepository := persistence.NewInventoryRepository(fakeClient)

	return scenarios.NewInventoryOnboardingUseCase(
		inventoryRepository,
		rule)
}

func TestInventoryOnboardingUseCaseExecuteSuccess(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	rule := &fakeRule{err: nil}
	request := dto.Request{
		Name:      "newtest",
		Namespace: "default",
	}
	err := newInventoryOnboardingUseCase(a, rule).Execute(request)
	a.Nil(err, "must onboard inventory without error")
}

func TestInventoryOnboardingUseCaseExecuteRuleFailed(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	rule := &fakeRule{err: errors.New("forbidden")}
	request := dto.Request{
		Name:      "newtest",
		Namespace: "default",
	}
	err := newInventoryOnboardingUseCase(a, rule).Execute(request)
	a.NotNil(err, "server must not be enabled")
	a.Contains(err.Error(), "forbidden", "must contain forbidden as an error")
}

type fakeRule struct {
	err error
}

func (f *fakeRule) Execute(_ dto.Request) error {
	return f.err
}
