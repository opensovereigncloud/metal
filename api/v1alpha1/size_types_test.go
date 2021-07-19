package v1alpha1

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
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
})
