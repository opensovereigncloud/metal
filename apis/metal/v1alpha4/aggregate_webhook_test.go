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

package v1alpha4

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Aggregate webhook", func() {
	const (
		AggregateNamespace = "default"
	)

	Context("When Aggregate is not created", func() {
		It("Should check that invalid CR will be rejected", func() {
			ctx := context.Background()
			crs := invalidAggregates(AggregateNamespace)

			for _, cr := range crs {
				By(fmt.Sprintf("Attempting to create Aggregate %s with invalid configuration", cr.Name))
				Expect(k8sClient.Create(ctx, &cr)).ShouldNot(Succeed())
			}
		})
	})

	Context("When Aggregate is not created", func() {
		It("Should check that valid CR will be accepted", func() {
			ctx := context.Background()
			crs := validAggregates(AggregateNamespace)

			for _, cr := range crs {
				By(fmt.Sprintf("Attempting to create Aggregate %s with valid configuration", cr.Name))
				Expect(k8sClient.Create(ctx, &cr)).Should(Succeed())
			}
		})
	})
})

func validAggregates(namespace string) []Aggregate {
	return []Aggregate{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "without-aggregate-type",
				Namespace: namespace,
			},
			Spec: AggregateSpec{
				Aggregates: []AggregateItem{
					{
						SourcePath: *JSONPathFromString("spec.cpus[*].siblings"),
						TargetPath: *JSONPathFromString("cpus.threads"),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-one-aggregate",
				Namespace: namespace,
			},
			Spec: AggregateSpec{
				Aggregates: []AggregateItem{
					{
						SourcePath: *JSONPathFromString("spec.cpus[*].siblings"),
						TargetPath: *JSONPathFromString("cpus.threads"),
						Aggregate:  CSumAggregateType,
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-multiple-aggregates",
				Namespace: namespace,
			},
			Spec: AggregateSpec{
				Aggregates: []AggregateItem{
					{
						SourcePath: *JSONPathFromString("spec.cpus[*].siblings"),
						TargetPath: *JSONPathFromString("cpus.threads"),
						Aggregate:  CSumAggregateType,
					},
					{
						SourcePath: *JSONPathFromString("spec.cpus[*].mhz"),
						TargetPath: *JSONPathFromString("cpus.freq"),
						Aggregate:  CAverageAggregateType,
					},
				},
			},
		},
	}
}

func invalidAggregates(namespace string) []Aggregate {
	return []Aggregate{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-non-existing-source-path",
				Namespace: namespace,
			},
			Spec: AggregateSpec{
				Aggregates: []AggregateItem{
					{
						SourcePath: *JSONPathFromString("spec.cpus[*].siblings123"),
						TargetPath: *JSONPathFromString("cpus.threads"),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-malformed-source-path",
				Namespace: namespace,
			},
			Spec: AggregateSpec{
				Aggregates: []AggregateItem{
					{
						SourcePath: *JSONPathFromString(":::"),
						TargetPath: *JSONPathFromString("cpus.threads"),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-malformed-target-path",
				Namespace: namespace,
			},
			Spec: AggregateSpec{
				Aggregates: []AggregateItem{
					{
						SourcePath: *JSONPathFromString("spec.cpus[*].siblings123"),
						TargetPath: *JSONPathFromString(":::"),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-clashing-target-paths",
				Namespace: namespace,
			},
			Spec: AggregateSpec{
				Aggregates: []AggregateItem{
					{
						SourcePath: *JSONPathFromString("spec.cpus[*].siblings"),
						TargetPath: *JSONPathFromString("cpus.threads"),
						Aggregate:  CSumAggregateType,
					},
					{
						SourcePath: *JSONPathFromString("spec.cpus[*].mhz"),
						TargetPath: *JSONPathFromString("cpus.threads"),
						Aggregate:  CAverageAggregateType,
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-parent-target-paths",
				Namespace: namespace,
			},
			Spec: AggregateSpec{
				Aggregates: []AggregateItem{
					{
						SourcePath: *JSONPathFromString("spec.cpus[*].siblings"),
						TargetPath: *JSONPathFromString("cpus"),
						Aggregate:  CSumAggregateType,
					},
					{
						SourcePath: *JSONPathFromString("spec.cpus[*].mhz"),
						TargetPath: *JSONPathFromString("cpus.threads"),
						Aggregate:  CAverageAggregateType,
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-child-target-paths",
				Namespace: namespace,
			},
			Spec: AggregateSpec{
				Aggregates: []AggregateItem{
					{
						SourcePath: *JSONPathFromString("spec.cpus[*].siblings"),
						TargetPath: *JSONPathFromString("cpus.threads"),
						Aggregate:  CSumAggregateType,
					},
					{
						SourcePath: *JSONPathFromString("spec.cpus[*].mhz"),
						TargetPath: *JSONPathFromString("cpus.threads.freq"),
						Aggregate:  CAverageAggregateType,
					},
				},
			},
		},
	}
}
