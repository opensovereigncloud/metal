package v1alpha1

import (
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Size match to Inventory", func() {
	Context("When value is collected with JSONPath", func() {
		It("Should be converted to string", func() {
			By("It is a string or pointer to string")
			expected := "ishouldmatch"
			values := []reflect.Value{
				reflect.ValueOf(expected),
				reflect.ValueOf(&expected),
			}

			for _, value := range values {
				actual, err := valueToString(&value)
				Expect(err).NotTo(HaveOccurred())
				Expect(actual).To(Equal(expected))
			}
		})

		It("Should not be converted to string and return an error", func() {
			By("It is other type than string")
			var testNil *string = nil
			testInt := 1
			testFloat := -1.43
			testStruct := resource.NewScaledQuantity(1, 5)
			values := []reflect.Value{
				reflect.ValueOf(testNil),
				reflect.ValueOf(&testNil),
				reflect.ValueOf(testInt),
				reflect.ValueOf(&testInt),
				reflect.ValueOf(testFloat),
				reflect.ValueOf(&testFloat),
				reflect.ValueOf(testStruct),
				reflect.ValueOf(&testStruct),
			}

			for _, value := range values {
				actual, err := valueToString(&value)
				Expect(err).To(HaveOccurred())
				Expect(actual).To(Equal(""))
			}
		})

		It("Should be converted to Quantity", func() {
			By("It is a numeric type or Quantity or numeric string")
			testInt8 := int8(-12)
			testInt16 := int16(-12)
			testInt32 := int32(-12)
			testInt64 := int64(-12)
			testInt := int(-12)
			testByte := byte(129)
			testUint8 := uint8(5)
			testUint16 := uint16(5)
			testUint32 := uint32(5)
			testUint64 := uint64(5)
			testUint := uint(5)
			testFloat32 := float32(5)
			testFloat64 := float64(5)
			testQuantityPtr := resource.NewScaledQuantity(100, 2)
			testQuantity := *testQuantityPtr
			testString := "1.487"

			values := []struct {
				value reflect.Value
				qty   resource.Quantity
			}{
				{
					value: reflect.ValueOf(testInt8),
					qty:   *resource.NewScaledQuantity(int64(testInt8), 0),
				},
				{
					value: reflect.ValueOf(&testInt8),
					qty:   *resource.NewScaledQuantity(int64(testInt8), 0),
				},
				{
					value: reflect.ValueOf(testInt16),
					qty:   *resource.NewScaledQuantity(int64(testInt16), 0),
				},
				{
					value: reflect.ValueOf(&testInt16),
					qty:   *resource.NewScaledQuantity(int64(testInt16), 0),
				},
				{
					value: reflect.ValueOf(testInt32),
					qty:   *resource.NewScaledQuantity(int64(testInt32), 0),
				},
				{
					value: reflect.ValueOf(&testInt32),
					qty:   *resource.NewScaledQuantity(int64(testInt32), 0),
				},
				{
					value: reflect.ValueOf(testInt64),
					qty:   *resource.NewScaledQuantity(testInt64, 0),
				},
				{
					value: reflect.ValueOf(&testInt64),
					qty:   *resource.NewScaledQuantity(testInt64, 0),
				},
				{
					value: reflect.ValueOf(testInt),
					qty:   *resource.NewScaledQuantity(int64(testInt), 0),
				},
				{
					value: reflect.ValueOf(&testInt),
					qty:   *resource.NewScaledQuantity(int64(testInt), 0),
				},
				{
					value: reflect.ValueOf(testByte),
					qty:   *resource.NewScaledQuantity(int64(testByte), 0),
				},
				{
					value: reflect.ValueOf(&testByte),
					qty:   *resource.NewScaledQuantity(int64(testByte), 0),
				},
				{
					value: reflect.ValueOf(testUint8),
					qty:   *resource.NewScaledQuantity(int64(testUint8), 0),
				},
				{
					value: reflect.ValueOf(&testUint8),
					qty:   *resource.NewScaledQuantity(int64(testUint8), 0),
				},
				{
					value: reflect.ValueOf(testUint16),
					qty:   *resource.NewScaledQuantity(int64(testUint16), 0),
				},
				{
					value: reflect.ValueOf(&testUint16),
					qty:   *resource.NewScaledQuantity(int64(testUint16), 0),
				},
				{
					value: reflect.ValueOf(testUint32),
					qty:   *resource.NewScaledQuantity(int64(testUint32), 0),
				},
				{
					value: reflect.ValueOf(&testUint32),
					qty:   *resource.NewScaledQuantity(int64(testUint32), 0),
				},
				{
					value: reflect.ValueOf(testUint64),
					qty:   resource.MustParse(fmt.Sprintf("%d", testUint64)),
				},
				{
					value: reflect.ValueOf(&testUint64),
					qty:   resource.MustParse(fmt.Sprintf("%d", testUint64)),
				},
				{
					value: reflect.ValueOf(testUint),
					qty:   resource.MustParse(fmt.Sprintf("%d", testUint)),
				},
				{
					value: reflect.ValueOf(&testUint),
					qty:   resource.MustParse(fmt.Sprintf("%d", testUint)),
				},
				{
					value: reflect.ValueOf(testFloat32),
					qty:   resource.MustParse(fmt.Sprintf("%f", testFloat32)),
				},
				{
					value: reflect.ValueOf(&testFloat32),
					qty:   resource.MustParse(fmt.Sprintf("%f", testFloat32)),
				},
				{
					value: reflect.ValueOf(testFloat64),
					qty:   resource.MustParse(fmt.Sprintf("%f", testFloat64)),
				},
				{
					value: reflect.ValueOf(&testFloat64),
					qty:   resource.MustParse(fmt.Sprintf("%f", testFloat64)),
				},
				{
					value: reflect.ValueOf(testQuantity),
					qty:   testQuantity,
				},
				{
					value: reflect.ValueOf(&testQuantity),
					qty:   testQuantity,
				},
				{
					value: reflect.ValueOf(&testQuantityPtr),
					qty:   testQuantity,
				},
				{
					value: reflect.ValueOf(testString),
					qty:   resource.MustParse(testString),
				},
				{
					value: reflect.ValueOf(&testString),
					qty:   resource.MustParse(testString),
				},
			}

			for _, value := range values {
				actual, err := valueToQuantity(&value.value)
				Expect(err).NotTo(HaveOccurred())
				Expect(actual.Cmp(value.qty)).To(Equal(0))
			}
		})

		It("Should not be converted to quantity and return an error", func() {
			By("It is non-numeric type")
			testString := "nonNumericString"
			testOtherStruct := Size{}
			var testNil *resource.Quantity = nil

			values := []reflect.Value{
				reflect.ValueOf(testString),
				reflect.ValueOf(&testString),
				reflect.ValueOf(testOtherStruct),
				reflect.ValueOf(&testOtherStruct),
				reflect.ValueOf(testNil),
				reflect.ValueOf(&testNil),
			}

			for _, value := range values {
				actual, err := valueToQuantity(&value)
				Expect(err).To(HaveOccurred())
				Expect(actual).To(BeNil())
			}
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
			CPUs: &CPUTotalSpec{
				Sockets: 2,
				Cores:   16,
				Threads: 32,
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
						Path: "system.id",
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
						Path: "system.serialNumber",
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
						Path: "system.id",
						Equal: &ConstraintValSpec{
							Literal: stringPtr("myInventoryId"),
						},
					},
					{
						Path: "system.manufacturer",
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
						Path: "system.id",
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
						Path: "cpus.cores",
						Equal: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(16, 0),
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
						Path: "cpus.cpus[0].model",
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
						Path: "cpus.cores",
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
						Path: "cpus.cpus[*].cores",
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
						Path: "cpus.cpus[*].cores",
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
						Path:      "cpus.cpus[*].cores",
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
						Path:      "cpus.cpus[*].cores",
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
						Path:      "cpus.cpus[*].cores",
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
						Path:      "cpus.cpus[*].cores",
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
						Path:        "cpus.cores",
						GreaterThan: resource.NewScaledQuantity(11, 0),
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
						Path:               "cpus.cores",
						GreaterThanOrEqual: resource.NewScaledQuantity(16, 0),
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
						Path:     "cpus.cores",
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
						Path:            "cpus.cores",
						LessThanOrEqual: resource.NewScaledQuantity(16, 0),
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
						Path:        "cpus.cores",
						GreaterThan: resource.NewScaledQuantity(8, 0),
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
						Path:               "cpus.cores",
						GreaterThanOrEqual: resource.NewScaledQuantity(16, 0),
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
						Path:            "cpus.cores",
						GreaterThan:     resource.NewScaledQuantity(12, 0),
						LessThanOrEqual: resource.NewScaledQuantity(24, 0),
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
						Path:               "cpus.cores",
						GreaterThanOrEqual: resource.NewScaledQuantity(12, 0),
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
						Path:        "cpus.cpus[*].cores",
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
						Path:               "cpus.cpus[*].cores",
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
						Path:      "cpus.cpus[*].cores",
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
						Path:            "cpus.cpus[*].cores",
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
						Path:        "cpus.cpus[*].cores",
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
						Path:               "cpus.cpus[*].cores",
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
						Path:            "cpus.cpus[*].cores",
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
						Path:               "cpus.cpus[*].cores",
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
						Path:        "cpus.cpus[*].cores",
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
						Path:               "cpus.cpus[*].cores",
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
						Path:      "cpus.cpus[*].cores",
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
						Path:            "cpus.cpus[*].cores",
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
						Path:        "cpus.cpus[*].cores",
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
						Path:               "cpus.cpus[*].cores",
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
						Path:            "cpus.cpus[*].cores",
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
						Path:               "cpus.cpus[*].cores",
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
						Path: "system.id",
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
						Path: "system.id",
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
						Path: "system.id",
						Equal: &ConstraintValSpec{
							Literal: stringPtr("myInventoryId"),
						},
					},
					{
						Path: "system.manufacturer",
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
						Path: "ipmis.ipAddress",
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
						Path: "cpus.cores",
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
						Path: "cpus.cores",
						NotEqual: &ConstraintValSpec{
							Numeric: resource.NewScaledQuantity(16, 0),
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
						Path: "cpus.cpus[0].model",
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
						Path: "cpus.memory.total",
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
						Path: "cpus.cpus[*].cores",
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
						Path: "cpus.cpus[*].cores",
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
						Path:      "cpus.cpus[*].cores",
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
						Path:      "cpus.cpus[*].cores",
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
						Path:      "cpus.cpus[*].cores",
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
						Path:      "cpus.cpus[*].cores",
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
						Path:        "cpus.cores",
						GreaterThan: resource.NewScaledQuantity(16, 0),
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
						Path:               "cpus.cores",
						GreaterThanOrEqual: resource.NewScaledQuantity(17, 0),
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
						Path:     "cpus.cores",
						LessThan: resource.NewScaledQuantity(16, 0),
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
						Path:            "cpus.cores",
						LessThanOrEqual: resource.NewScaledQuantity(10, 0),
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
						Path:        "cpus.cores",
						GreaterThan: resource.NewScaledQuantity(1, 0),
						LessThan:    resource.NewScaledQuantity(16, 0),
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
						Path:               "cpus.cores",
						GreaterThanOrEqual: resource.NewScaledQuantity(1, 0),
						LessThan:           resource.NewScaledQuantity(16, 0),
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
						Path:            "cpus.cores",
						GreaterThan:     resource.NewScaledQuantity(1, 0),
						LessThanOrEqual: resource.NewScaledQuantity(15, 0),
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
						Path:               "cpus.cores",
						GreaterThanOrEqual: resource.NewScaledQuantity(17, 0),
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
						Path:        "cpus.cpus[*].cores",
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
						Path:               "cpus.cpus[*].cores",
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
						Path:      "cpus.cpus[*].cores",
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
						Path:            "cpus.cpus[*].cores",
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
						Path:        "cpus.cpus[*].cores",
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
						Path:               "cpus.cpus[*].cores",
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
						Path:            "cpus.cpus[*].cores",
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
						Path:               "cpus.cpus[*].cores",
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
						Path:        "cpus.cpus[*].cores",
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
						Path:               "cpus.cpus[*].cores",
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
						Path:      "cpus.cpus[*].cores",
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
						Path:            "cpus.cpus[*].cores",
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
						Path:        "cpus.cpus[*].cores",
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
						Path:               "cpus.cpus[*].cores",
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
						Path:            "cpus.cpus[*].cores",
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
						Path:               "cpus.cpus[*].cores",
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
						Path: "system.id",
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
						Path: "cpus.sockets",
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
						Path: "cpus",
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
						Path:  "cpus",
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
