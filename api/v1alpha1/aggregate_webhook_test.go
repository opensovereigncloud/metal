package v1alpha1

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aggregate webhook", func() {
	const (
		AggregateNamespace = "default"
	)

	Context("When ", func() {
		It("Should ", func() {
			By("")
			acc1 := make(map[string]interface{})
			Expect(accountPath(acc1, []string{"a"})).To(Succeed())
			Expect(accountPath(acc1, []string{"a", "b"})).NotTo(Succeed())

			acc2 := make(map[string]interface{})
			Expect(accountPath(acc2, []string{"a", "b"})).To(Succeed())
			Expect(accountPath(acc2, []string{"a"})).NotTo(Succeed())

			acc3 := make(map[string]interface{})
			Expect(accountPath(acc3, []string{"a", "b"})).To(Succeed())
			Expect(accountPath(acc3, []string{"a", "b"})).NotTo(Succeed())

			acc4 := make(map[string]interface{})
			Expect(accountPath(acc4, []string{"a", "b"})).To(Succeed())
			Expect(accountPath(acc4, []string{"a", "d"})).To(Succeed())
			Expect(accountPath(acc4, []string{"a", "e", "f"})).To(Succeed())
			Expect(accountPath(acc4, []string{"b", "d"})).To(Succeed())
		})
	})
})
