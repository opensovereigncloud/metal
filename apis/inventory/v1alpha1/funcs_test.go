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
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
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

	Context("When tokenized path is inserted to map", func() {
		It("Should be accounted", func() {
			By("Checking for failure when parent path is used to set value")
			acc1 := make(map[string]interface{})
			Expect(accountPath(acc1, []string{"a"})).To(Succeed())
			Expect(accountPath(acc1, []string{"a", "b"})).NotTo(Succeed())

			By("Checking for failure when child path is used to set value")
			acc2 := make(map[string]interface{})
			Expect(accountPath(acc2, []string{"a", "b"})).To(Succeed())
			Expect(accountPath(acc2, []string{"a"})).NotTo(Succeed())

			By("Checking for failure when duplicate paths are used to set value")
			acc3 := make(map[string]interface{})
			Expect(accountPath(acc3, []string{"a", "b"})).To(Succeed())
			Expect(accountPath(acc3, []string{"a", "b"})).NotTo(Succeed())

			By("Checking for success when pats are unique")
			acc4 := make(map[string]interface{})
			Expect(accountPath(acc4, []string{"a", "b"})).To(Succeed())
			Expect(accountPath(acc4, []string{"a", "d"})).To(Succeed())
			Expect(accountPath(acc4, []string{"a", "e", "f"})).To(Succeed())
			Expect(accountPath(acc4, []string{"b", "d"})).To(Succeed())
		})
	})
})
