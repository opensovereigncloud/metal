package v1alpha1

import (
	"context"

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
						Path: "system.id",
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
						Path: "{.system.id}",
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
						Path: ".system.id",
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
						Path: "system.id",
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
						Path: "{.system.id}",
						Equal: &ConstraintValSpec{
							Literal: &sampleLiteralValue,
						},
					},
					{
						Path: "{.system.manufacturer}",
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
						Path: "system.id",
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
						Path: "system.id",
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
						Path: "blocks.count",
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
						Path: "blocks.count",
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
						Path:     "blocks.count",
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
						Path:            "blocks.count",
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
						Path:        "blocks.count",
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
						Path:               "blocks.count",
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
						Path:        "blocks.count",
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
						Path:            "blocks.count",
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
						Path:               "blocks.count",
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
						Path:               "blocks.count",
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
						Path:               "{.blocks.count",
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
						Path:               "path.that.doesnt.eists",
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
						Path:               "blocks.count",
						GreaterThanOrEqual: resource.NewScaledQuantity(1, 0),
					},
					{
						Path:            "blocks.count",
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
						Path: "blocks.count",
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
						Path:      "blocks.blocks[*].size",
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
						Path: "blocks.count",
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
						Path:            "blocks.count",
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
						Path:               "blocks.count",
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
						Path:     "blocks.count",
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
						Path:            "blocks.count",
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
						Path:        "blocks.count",
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
						Path:               "blocks.count",
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
						Path:        "blocks.count",
						GreaterThan: resource.NewScaledQuantity(5, 0),
						LessThan:    resource.NewScaledQuantity(1, 0),
					},
				},
			},
		},
	}
}
