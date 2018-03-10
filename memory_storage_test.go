package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestModel struct {
	Name string `registry:"name"`
}

var _ = Describe("MemoryStorage", func() {
	It("Registering a valid Type should work", func() {
		s := NewMemoryStorage()
		r := New(s)
		TestType, err := r.NewType(&TestModel{}, "test_model")
		Expect(err).NotTo(HaveOccurred())

		TestType.Create(TestModel{
			Name: "some name",
		})
		result, err := TestType.Get(QueryAttribute{
			Field:     "name",
			Value:     "some name",
			Condition: Conditions.Equals,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(len(result)).To(Equal(1))
		Expect(result[0].(TestModel).Name).To(Equal("some name"))
	})
})
