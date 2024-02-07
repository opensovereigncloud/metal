// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha4

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

var _ = Describe("ConstraintValSpec marshalling and unmarshalling", func() {
	Context("When JSON is deserialized to ConstraintValSpec", func() {
		It("Should accept integers, numeric strings or literals", func() {
			By("Deserializing integers to Quantity")
			integerJsons := []string{
				`12345`,
			}
			for _, j := range integerJsons {
				cvs := &ConstraintValSpec{}

				Expect(json.Unmarshal([]byte(j), cvs)).Should(Succeed())
				Expect(cvs.Numeric).NotTo(BeNil())
				Expect(cvs.Literal).To(BeNil())

				q := resource.MustParse(j)

				Expect(q.Equal(*cvs.Numeric)).To(BeTrue())
			}

			By("Deserializing numeric strings to Quantity")
			numericStringJsons := []string{
				`"12"`,
				`"3.14"`,
				`"1e-5"`,
				`"1e7"`,
				`"5Gi"`,
				`"1.5M"`,
			}
			for _, j := range numericStringJsons {
				cvs := &ConstraintValSpec{}

				Expect(json.Unmarshal([]byte(j), cvs)).Should(Succeed())
				Expect(cvs.Numeric).NotTo(BeNil())
				Expect(cvs.Literal).To(BeNil())

				q := resource.MustParse(strings.Trim(j, "\""))

				Expect(q.Equal(*cvs.Numeric)).To(BeTrue())
			}

			By("Deserializing literals to string")
			nonNumericStringJsons := []string{
				`"123justalphas"`,
				`"justalphas"`,
				`"just1alphas"`,
				`"justalphas123"`,
				`"null"`,
				`""`,
			}
			for _, j := range nonNumericStringJsons {
				cvs := &ConstraintValSpec{}

				Expect(json.Unmarshal([]byte(j), cvs)).Should(Succeed())
				Expect(cvs.Numeric).To(BeNil())
				Expect(cvs.Literal).NotTo(BeNil())

				Expect(*cvs.Literal).To(BeEquivalentTo(strings.Trim(j, "\"")))
			}

			By("Deserializing null to nil struct values")
			nullStringJsons := `null`
			cvs := &ConstraintValSpec{}

			Expect(json.Unmarshal([]byte(nullStringJsons), cvs)).Should(Succeed())
			Expect(cvs.Numeric).To(BeNil())
			Expect(cvs.Literal).To(BeNil())
		})
	})

	Context("When ConstraintValSpec is serialized to Json", func() {
		It("Should be transformed to single JSON value", func() {
			By("Serializing Quantity to numeric string")
			qSpec := &ConstraintValSpec{
				Numeric: resource.NewScaledQuantity(1, 0),
			}
			b, err := json.Marshal(qSpec)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(b)).To(BeEquivalentTo(fmt.Sprintf(`"%s"`, qSpec.Numeric.String())))

			By("Serializing string to string")
			s := "abc"
			sSpec := &ConstraintValSpec{
				Literal: &s,
			}
			b, err = json.Marshal(sSpec)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(b)).To(BeEquivalentTo(fmt.Sprintf(`"%s"`, s)))

			By("Failing to serialize when both values are set")
			fSpec := &ConstraintValSpec{
				Literal: &s,
				Numeric: resource.NewScaledQuantity(1, 0),
			}
			_, err = json.Marshal(fSpec)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("When inventory is getting matched to size", func() {
		It("Should succeed", func() {
			sizes := sizesShouldMatch()
			inventory := inventory()

			for _, size := range sizes {
				By(fmt.Sprintf("Matching to size %s", size.Name))
				matches, err := size.Matches(inventory)
				Expect(err).NotTo(HaveOccurred())
				Expect(matches).To(BeTrue())
			}
		})

		It("Should fail", func() {
			sizes := sizesShouldNotMatch()
			inventory := inventory()

			for _, size := range sizes {
				By(fmt.Sprintf("Matching to size %s", size.Name))
				matches, err := size.Matches(inventory)
				Expect(err).NotTo(HaveOccurred())
				Expect(matches).To(BeFalse())
			}
		})

		It("Should produce error", func() {
			sizes := sizesShouldReturnErr()
			inventory := inventory()

			for _, size := range sizes {
				By(fmt.Sprintf("Matching to size %s", size.Name))
				matches, err := size.Matches(inventory)
				Expect(err).To(HaveOccurred())
				Expect(matches).To(BeFalse())
			}
		})
	})
})

func inventory() *Inventory {
	return &Inventory{
		Spec: InventorySpec{
			System: &SystemSpec{
				ID:           "myInventoryId",
				Manufacturer: "myManufacturer",
			},
			CPUs: []CPUSpec{
				{
					Model: "78",
					Cores: 8,
				},
				{
					Cores: 8,
				},
			},
		},
	}
}

func sizesShouldMatch() []Size {
	return []Size{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "empty-spec",
			},
			Spec: SizeSpec{},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "string-equal",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.system.id"),
						Equal: &ConstraintValSpec{
							Literal: stringPtr("myInventoryId"),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "empty-string-equal",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.system.serialNumber"),
						Equal: &ConstraintValSpec{
							Literal: stringPtr(""),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "multiple-string-equal",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.system.id"),
						Equal: &ConstraintValSpec{
							Literal: stringPtr("myInventoryId"),
						},
					},
					{
						Path: *JSONPathFromString("spec.system.manufacturer"),
						Equal: &ConstraintValSpec{
							Literal: stringPtr("myManufacturer"),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "string-not-equal",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.system.id"),
						NotEqual: &ConstraintValSpec{
							Literal: stringPtr("blacklistedSystemId"),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-equal",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.cpus[0].cores"),
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(8, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-equal-to-numeric-string",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.cpus[0].model"),
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(78, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-not-equal",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.cpus[1].cores"),
						NotEqual: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(64, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-equal-foreach",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.cpus[*].cores"),
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(8, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-not-equal-foreach",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.cpus[*].cores"),
						NotEqual: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(16, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-equal-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:      *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate: CSumAggregateType,
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(16, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-equal-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:      *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate: CAverageAggregateType,
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(8, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-not-equal-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:      *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate: CSumAggregateType,
						NotEqual: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(24, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-not-equal-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:      *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate: CAverageAggregateType,
						NotEqual: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(24, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.cpus[0].cores"),
						GreaterThan: resource.NewScaledQuantity(7, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[0].cores"),
						GreaterThanOrEqual: resource.NewScaledQuantity(8, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-lt",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:     *JSONPathFromString("spec.cpus[0].cores"),
						LessThan: resource.NewScaledQuantity(17, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-lte",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.cpus[0].cores"),
						LessThanOrEqual: resource.NewScaledQuantity(8, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-lt",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.cpus[0].cores"),
						GreaterThan: resource.NewScaledQuantity(4, 0),
						LessThan:    resource.NewScaledQuantity(24, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-lt",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[0].cores"),
						GreaterThanOrEqual: resource.NewScaledQuantity(8, 0),
						LessThan:           resource.NewScaledQuantity(24, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-lte",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.cpus[0].cores"),
						GreaterThan:     resource.NewScaledQuantity(3, 0),
						LessThanOrEqual: resource.NewScaledQuantity(8, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-lte",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[0].cores"),
						GreaterThanOrEqual: resource.NewScaledQuantity(4, 0),
						LessThanOrEqual:    resource.NewScaledQuantity(12, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:   CAverageAggregateType,
						GreaterThan: resource.NewScaledQuantity(6, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:          CAverageAggregateType,
						GreaterThanOrEqual: resource.NewScaledQuantity(8, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-lt-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:      *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate: CAverageAggregateType,
						LessThan:  resource.NewScaledQuantity(11, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-lte-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:       CAverageAggregateType,
						LessThanOrEqual: resource.NewScaledQuantity(8, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-lt-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:   CAverageAggregateType,
						GreaterThan: resource.NewScaledQuantity(6, 0),
						LessThan:    resource.NewScaledQuantity(10, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-lt-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:          CAverageAggregateType,
						GreaterThanOrEqual: resource.NewScaledQuantity(8, 0),
						LessThan:           resource.NewScaledQuantity(11, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-lte-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:       CAverageAggregateType,
						GreaterThan:     resource.NewScaledQuantity(6, 0),
						LessThanOrEqual: resource.NewScaledQuantity(8, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-lte-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:          CAverageAggregateType,
						GreaterThanOrEqual: resource.NewScaledQuantity(4, 0),
						LessThanOrEqual:    resource.NewScaledQuantity(24, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:   CSumAggregateType,
						GreaterThan: resource.NewScaledQuantity(6, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:          CSumAggregateType,
						GreaterThanOrEqual: resource.NewScaledQuantity(16, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-lt-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:      *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate: CSumAggregateType,
						LessThan:  resource.NewScaledQuantity(20, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-lte-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:       CSumAggregateType,
						LessThanOrEqual: resource.NewScaledQuantity(16, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-lt-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:   CSumAggregateType,
						GreaterThan: resource.NewScaledQuantity(6, 0),
						LessThan:    resource.NewScaledQuantity(18, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-lt-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:          CSumAggregateType,
						GreaterThanOrEqual: resource.NewScaledQuantity(16, 0),
						LessThan:           resource.NewScaledQuantity(20, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-lte-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:       CSumAggregateType,
						GreaterThan:     resource.NewScaledQuantity(6, 0),
						LessThanOrEqual: resource.NewScaledQuantity(16, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-lte-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:          CSumAggregateType,
						GreaterThanOrEqual: resource.NewScaledQuantity(10, 0),
						LessThanOrEqual:    resource.NewScaledQuantity(20, 0),
					},
				},
			},
		},
	}
}

func sizesShouldNotMatch() []Size {
	return []Size{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "string-equal",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.system.id"),
						Equal: &ConstraintValSpec{
							Literal: stringPtr("anotherInventoryId"),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "string-not-equal",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.system.id"),
						NotEqual: &ConstraintValSpec{
							Literal: stringPtr("myInventoryId"),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "string-multiple-equal-with-one-match",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.system.id"),
						Equal: &ConstraintValSpec{
							Literal: stringPtr("myInventoryId"),
						},
					},
					{
						Path: *JSONPathFromString("spec.system.manufacturer"),
						Equal: &ConstraintValSpec{
							Literal: stringPtr("nonExistentManufacturer"),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "string-not-equal-on-absent",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.ipmis.ipAddress"),
						NotEqual: &ConstraintValSpec{
							Literal: stringPtr("192.168.1.1"),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-equal",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.cpus[0].cores"),
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(10, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-not-equal",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.cpus[0].cores"),
						NotEqual: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(8, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-equal-to-numeric-string",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.cpus[0].model"),
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(77, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-not-equal-on-absent",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.cpus.memory.total"),
						NotEqual: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(8, resource.Giga),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-equal-foreach",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.cpus[*].cores"),
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(9, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-not-equal-foreach",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.cpus[*].cores"),
						NotEqual: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(8, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-equal-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:      *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate: CSumAggregateType,
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(20, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-equal-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:      *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate: CAverageAggregateType,
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(6, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-not-equal-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:      *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate: CSumAggregateType,
						NotEqual: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(16, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-not-equal-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:      *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate: CAverageAggregateType,
						NotEqual: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(8, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.cpus[0].cores"),
						GreaterThan: resource.NewScaledQuantity(8, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[0].cores"),
						GreaterThanOrEqual: resource.NewScaledQuantity(9, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-lt",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:     *JSONPathFromString("spec.cpus[0].cores"),
						LessThan: resource.NewScaledQuantity(8, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-lte",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.cpus[0].cores"),
						LessThanOrEqual: resource.NewScaledQuantity(7, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-lt",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.cpus[0].cores"),
						GreaterThan: resource.NewScaledQuantity(1, 0),
						LessThan:    resource.NewScaledQuantity(8, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-lt",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[0].cores"),
						GreaterThanOrEqual: resource.NewScaledQuantity(1, 0),
						LessThan:           resource.NewScaledQuantity(8, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-lte",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.cpus[0].cores"),
						GreaterThan:     resource.NewScaledQuantity(1, 0),
						LessThanOrEqual: resource.NewScaledQuantity(7, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-lte",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[0].cores"),
						GreaterThanOrEqual: resource.NewScaledQuantity(9, 0),
						LessThanOrEqual:    resource.NewScaledQuantity(24, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:   CAverageAggregateType,
						GreaterThan: resource.NewScaledQuantity(8, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:          CAverageAggregateType,
						GreaterThanOrEqual: resource.NewScaledQuantity(10, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-lt-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:      *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate: CAverageAggregateType,
						LessThan:  resource.NewScaledQuantity(8, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-lte-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:       CAverageAggregateType,
						LessThanOrEqual: resource.NewScaledQuantity(7, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-lt-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:   CAverageAggregateType,
						GreaterThan: resource.NewScaledQuantity(8, 0),
						LessThan:    resource.NewScaledQuantity(10, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-lt-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:          CAverageAggregateType,
						GreaterThanOrEqual: resource.NewScaledQuantity(9, 0),
						LessThan:           resource.NewScaledQuantity(11, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-lte-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:       CAverageAggregateType,
						GreaterThan:     resource.NewScaledQuantity(6, 0),
						LessThanOrEqual: resource.NewScaledQuantity(7, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-lte-avg-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:          CAverageAggregateType,
						GreaterThanOrEqual: resource.NewScaledQuantity(4, 0),
						LessThanOrEqual:    resource.NewScaledQuantity(4, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:   CSumAggregateType,
						GreaterThan: resource.NewScaledQuantity(16, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:          CSumAggregateType,
						GreaterThanOrEqual: resource.NewScaledQuantity(17, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-lt-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:      *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate: CSumAggregateType,
						LessThan:  resource.NewScaledQuantity(16, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-lte-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:       CSumAggregateType,
						LessThanOrEqual: resource.NewScaledQuantity(15, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-lt-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:        *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:   CSumAggregateType,
						GreaterThan: resource.NewScaledQuantity(16, 0),
						LessThan:    resource.NewScaledQuantity(18, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-lt-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:          CSumAggregateType,
						GreaterThanOrEqual: resource.NewScaledQuantity(17, 0),
						LessThan:           resource.NewScaledQuantity(20, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gt-lte-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:            *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:       CSumAggregateType,
						GreaterThan:     resource.NewScaledQuantity(6, 0),
						LessThanOrEqual: resource.NewScaledQuantity(15, 0),
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-gte-lte-sum-aggregate",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:               *JSONPathFromString("spec.cpus[*].cores"),
						Aggregate:          CSumAggregateType,
						GreaterThanOrEqual: resource.NewScaledQuantity(2, 0),
						LessThanOrEqual:    resource.NewScaledQuantity(2, 0),
					},
				},
			},
		},
	}
}

func sizesShouldReturnErr() []Size {
	return []Size{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "numeric-equal-to-literal",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.system.id"),
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(3, 0),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "literal-equal-to-numeric",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.cpus[0].cores"),
						Equal: &ConstraintValSpec{
							Literal: stringPtr("nonNumericString"),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "equal-to-unsupported-struct",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path: *JSONPathFromString("spec.cpus"),
						Equal: &ConstraintValSpec{
							Literal: stringPtr("nonNumericString"),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "equal-to-nil",
			},
			Spec: SizeSpec{
				Constraints: []ConstraintSpec{
					{
						Path:  *JSONPathFromString("spec.cpus"),
						Equal: &ConstraintValSpec{},
					},
				},
			},
		},
	}
}

func stringPtr(s string) *string {
	return &s
}
