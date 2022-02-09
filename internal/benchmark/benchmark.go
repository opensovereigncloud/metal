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

package benchmark

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/onmetal/metal-api/apis/benchmark/v1alpha3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const percentageModifier = 100

const UUIDLabel = "machine.onmetal.de/uuid"

var (
	ErrAlreadyExist     = errors.New("already exist")
	ErrBenchTypeUnknown = errors.New("unknown benchmark type")
)

type Benchmark struct {
	client.Client

	log             logr.Logger
	ctx             context.Context
	name, namespace string
}

func New(ctx context.Context, c client.Client, l logr.Logger, req ctrl.Request) (*Benchmark, bool) {
	if isExist(ctx, c, req) {
		return &Benchmark{}, true
	}
	return &Benchmark{
		Client:    c,
		ctx:       ctx,
		log:       l,
		name:      req.Name,
		namespace: req.Namespace,
	}, false
}

func (b *Benchmark) Create() error {
	obj := prepareMachineBenchmark(b.name, b.namespace)
	return b.Client.Create(b.ctx, obj)
}

func isExist(ctx context.Context, c client.Client, req ctrl.Request) bool {
	m := &v1alpha3.Machine{}
	if err := c.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, m); err != nil {
		if apierrors.IsNotFound(err) {
			return false
		}
		return false
	}
	return true
}

func prepareMachineBenchmark(name, namespace string) *v1alpha3.Machine {
	return &v1alpha3.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{UUIDLabel: name},
		},
	}
}

func CalculateDeviation(oldObj, newObj *v1alpha3.Machine) map[string]v1alpha3.BenchmarkDeviations {
	return calculateMachineDeviation(oldObj, newObj)
}

func calculateMachineDeviation(oldObj, newObj *v1alpha3.Machine) map[string]v1alpha3.BenchmarkDeviations {
	md := make(map[string]v1alpha3.BenchmarkDeviations, len(newObj.Spec.Benchmarks))
	for nn, newValue := range newObj.Spec.Benchmarks {
		oldValue, ok := oldObj.Spec.Benchmarks[nn]
		if !ok {
			continue
		}
		md[nn] = calculateDiffForBenchmark(oldValue, newValue)
	}
	return md
}

func calculateDiffForBenchmark(oldV, newV []v1alpha3.Benchmark) []v1alpha3.BenchmarkDeviation {
	dv := make([]v1alpha3.BenchmarkDeviation, 0, len(newV))
	for n := range newV {
		for o := range oldV {
			if newV[n].Name != oldV[o].Name {
				continue
			}
			dv = append(dv, v1alpha3.BenchmarkDeviation{
				Name:  newV[n].Name,
				Value: percentageChange(oldV[o].Value, newV[n].Value),
			})
		}
	}
	return dv
}

func percentageChange(oldValue, newValue uint64) string {
	if oldValue > newValue {
		decrease := oldValue - newValue
		return fmt.Sprintf("-%.1f%%", float64(decrease)/float64(oldValue)*percentageModifier)
	}
	increase := newValue - oldValue
	return fmt.Sprintf("+%.1f%%", float64(increase)/float64(oldValue)*percentageModifier)
}
