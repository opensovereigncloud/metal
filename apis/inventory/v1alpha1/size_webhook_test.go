// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */
// nolint
package v1alpha1

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Size webhook", func() {
	const (
		SizeNamespace = "default"
	)

	Context("When Size is not created", func() {
		It("Should check that invalid CR will be rejected", func() {
			By("Attempting to create Size with invalid configuration")

			ctx := context.Background()
			crs := invalidSizes(SizeNamespace)

			for _, cr := range crs {
				Expect(k8sClient.Create(ctx, &cr)).ShouldNot(Succeed())
			}
		})
	})

	Context("When Size is not created", func() {
		It("Should check that valid CR will be accepted", func() {
			By("Attempting to create Size with valid configuration")

			ctx := context.Background()
			crs := validSizes(SizeNamespace)

			for _, cr := range crs {
				By(fmt.Sprintf("Creating size %s", cr.Name))
				Expect(k8sClient.Create(ctx, &cr)).Should(Succeed())
			}
		})
	})
})

func validSizes(namespace string) []Size {
	sampleLiteralValue := "abc"

	return []Size{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-empty-spec",
				Namespace: namespace,
			},
			Spec: SizeSpec{},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "without-description",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.system.id"),
						Equal: &ConstraintValSpec{
							Literal: &sampleLiteralValue,
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-full-json-path",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("{.spec.system.id}"),
						Equal: &ConstraintValSpec{
							Literal: &sampleLiteralValue,
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-leading-dot-path",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString(".spec.system.id"),
						Equal: &ConstraintValSpec{
							Literal: &sampleLiteralValue,
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-simple-path",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.system.id"),
						Equal: &ConstraintValSpec{
							Literal: &sampleLiteralValue,
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-multiple-constraints",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("{.spec.system.id}"),
						Equal: &ConstraintValSpec{
							Literal: &sampleLiteralValue,
						},
					},
					{
						Path: *JSONPathFromString("{.spec.system.manufacturer}"),
						Equal: &ConstraintValSpec{
							Literal: &sampleLiteralValue,
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-equal-to-literal-constraint",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.system.id"),
						Equal: &ConstraintValSpec{
							Literal: &sampleLiteralValue,
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-not-equal-to-literal-constraint",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.system.id"),
						NotEqual: &ConstraintValSpec{
							Literal: &sampleLiteralValue,
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-equal-to-numeric-constraint",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.blocks[0].size"),
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(0, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-not-equal-to-numeric-constraint",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.blocks[0].size"),
						NotEqual: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(0, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-lt-constraint",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:     *JSONPathFromString("spec.blocks[0].size"),
						LessThan: resource.NewScaledQuantity(5, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-lte-constraint",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.blocks[0].size"),
						LessThanOrEqual: resource.NewScaledQuantity(5, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-gt-constraint",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.blocks[0].size"),
						GreaterThan: resource.NewScaledQuantity(1, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-gte-constraint",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.blocks[0].size"),
						GreaterThanOrEqual: resource.NewScaledQuantity(1, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-lt-gt-constraint",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.blocks[0].size"),
						LessThan:    resource.NewScaledQuantity(5, 0),
						GreaterThan: resource.NewScaledQuantity(1, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-lte-gt-constraint",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.blocks[0].size"),
						LessThanOrEqual: resource.NewScaledQuantity(5, 0),
						GreaterThan:     resource.NewScaledQuantity(1, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-lte-gte-constraint",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.blocks[0].size"),
						LessThanOrEqual:    resource.NewScaledQuantity(5, 0),
						GreaterThanOrEqual: resource.NewScaledQuantity(1, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-lt-gte-constraint",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.blocks[0].size"),
						LessThan:           resource.NewScaledQuantity(5, 0),
						GreaterThanOrEqual: resource.NewScaledQuantity(1, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-constraint-on-computed-aggregates",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("status.computed.name.key"),
						LessThan:           resource.NewScaledQuantity(5, 0),
						GreaterThanOrEqual: resource.NewScaledQuantity(1, 0),
					},
				},
			},
		},
	}
}

func invalidSizes(namespace string) []Size {
	sampleLiteralValue := "abc"

	return []Size{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-malformed-path",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("{.spec.blocks.count"),
						GreaterThanOrEqual: resource.NewScaledQuantity(1, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-invalid-path",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("path.that.doesnt.exists"),
						GreaterThanOrEqual: resource.NewScaledQuantity(1, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-two-constraints-for-one-path",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.blocks[0].size"),
						GreaterThanOrEqual: resource.NewScaledQuantity(1, 0),
					},
					{
						Path:            *JSONPathFromString("spec.blocks[0].size"),
						LessThanOrEqual: resource.NewScaledQuantity(5, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-constraint-without-conditions",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.blocks[0].size"),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-aggregate-on-literals",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:      *JSONPathFromString("spec.blocks[*].size"),
						Aggregate: CSumAggregateType,
						Equal: &ConstraintValSpec{
							Literal: &sampleLiteralValue,
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-eq-and-neq",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.blocks[0].size"),
						Equal: &ConstraintValSpec{
							Literal: &sampleLiteralValue,
						},
						NotEqual: &ConstraintValSpec{
							Literal: &sampleLiteralValue,
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-lt-and-lte",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.blocks[0].size"),
						LessThan:        resource.NewScaledQuantity(5, 0),
						LessThanOrEqual: resource.NewScaledQuantity(5, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-gt-and-gte",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.blocks[0].size"),
						GreaterThan:        resource.NewScaledQuantity(5, 0),
						GreaterThanOrEqual: resource.NewScaledQuantity(5, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-lt-and-eq",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:     *JSONPathFromString("spec.blocks[0].size"),
						LessThan: resource.NewScaledQuantity(5, 0),
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(5, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-lte-and-eq",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.blocks[0].size"),
						LessThanOrEqual: resource.NewScaledQuantity(5, 0),
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(5, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-gt-and-eq",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.blocks[0].size"),
						GreaterThan: resource.NewScaledQuantity(5, 0),
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(5, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-gte-and-eq",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.blocks[0].size"),
						GreaterThanOrEqual: resource.NewScaledQuantity(5, 0),
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(5, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "with-wrong-interval",
				Namespace: namespace,
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.blocks[0].size"),
						GreaterThan: resource.NewScaledQuantity(5, 0),
						LessThan:    resource.NewScaledQuantity(1, 0),
					},
				},
			},
		},
	}
}
