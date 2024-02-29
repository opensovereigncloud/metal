// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package benchmark

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
)

const percentageModifier = 100

const UUIDLabel = "metal.ironcore.dev/uuid"

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
	m := &metalv1alpha4.Benchmark{}
	if err := c.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, m); err != nil {
		if apierrors.IsNotFound(err) {
			return false
		}
		return false
	}
	return true
}

func prepareMachineBenchmark(name, namespace string) *metalv1alpha4.Benchmark {
	return &metalv1alpha4.Benchmark{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{UUIDLabel: name},
		},
	}
}

func CalculateDeviation(oldObj, newObj *metalv1alpha4.Benchmark) map[string]metalv1alpha4.BenchmarkDeviations {
	return calculateMachineDeviation(oldObj, newObj)
}

func calculateMachineDeviation(oldObj, newObj *metalv1alpha4.Benchmark) map[string]metalv1alpha4.BenchmarkDeviations {
	md := make(map[string]metalv1alpha4.BenchmarkDeviations, len(newObj.Spec.Benchmarks))
	for nn, newValue := range newObj.Spec.Benchmarks {
		oldValue, ok := oldObj.Spec.Benchmarks[nn]
		if !ok {
			continue
		}
		md[nn] = calculateDiffForBenchmark(oldValue, newValue)
	}
	return md
}

func calculateDiffForBenchmark(oldV, newV []metalv1alpha4.BenchmarkResult) []metalv1alpha4.BenchmarkDeviation {
	dv := make([]metalv1alpha4.BenchmarkDeviation, 0, len(newV))
	for n := range newV {
		for o := range oldV {
			if newV[n].Name != oldV[o].Name {
				continue
			}
			dv = append(dv, metalv1alpha4.BenchmarkDeviation{
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
