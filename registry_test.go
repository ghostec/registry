package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestModel struct {
	Name string `registry:"name"`
}

var _ = Describe("Registry tests", func() {
	Describe("MemoryStorage", func() {
		It("Registering a valid Type should work", func() {
			s := NewMemoryStorage()
			r := New(s)
			TestType, err := r.NewType(&TestModel{}, "test_models")
			Expect(err).NotTo(HaveOccurred())

			TestType.Create(TestModel{
				Name: "some name",
			})
			result, err := TestType.Get(QueryAttribute{
				Tag:       "name",
				Value:     "some name",
				Condition: Conditions.Equals,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(result)).To(Equal(1))
			Expect(result[0].(TestModel).Name).To(Equal("some name"))
		})

		It("Type.Create should deep copy original content", func() {
			s := NewMemoryStorage()
			r := New(s)
			TestType, err := r.NewType(&TestModel{}, "test_models")
			Expect(err).NotTo(HaveOccurred())

			original := &TestModel{
				Name: "some name",
			}
			TestType.Create(original)
			original.Name = "changed name"

			result, err := TestType.Get(QueryAttribute{
				Tag:       "name",
				Value:     "some name",
				Condition: Conditions.Equals,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(result)).To(Equal(1))
			Expect(result[0].(TestModel).Name).To(Equal("some name"))
		})
	})

	Describe("query without associations", func() {
		It("query Get should return the same as Type.Get", func() {
			s := NewMemoryStorage()
			r := New(s)
			TestType, err := r.NewType(&TestModel{}, "test_models")
			Expect(err).NotTo(HaveOccurred())

			TestType.Create(&TestModel{
				Name: "some name",
			})
			resultT, err := TestType.Get(QueryAttribute{
				Tag:       "name",
				Value:     "some name",
				Condition: Conditions.Equals,
			})
			Expect(err).NotTo(HaveOccurred())

			query := &query{rt: TestType}
			resultQ, err := query.Get(QueryAttribute{
				Tag:       "name",
				Value:     "some name",
				Condition: Conditions.Equals,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(resultQ)).To(Equal(len(resultT)))
			Expect(resultQ[0].(TestModel).Name).
				To(Equal(resultT[0].(TestModel).Name))
		})
	})

	Describe("Associations", func() {
		It("Get() should get nested types ", func() {
			type Y struct {
				ID   string `registry:"id"`
				Attr string `registry:"attr"`
				XID  string `registry:"x_id"`
			}
			type X struct {
				ID   string `registry:"id"`
				Name string `registry:"name"`
				Ys   []Y    `registry:"ys"`
			}

			// registering types
			s := NewMemoryStorage()
			r := New(s)
			XType, err := r.NewType(&X{}, "xs")
			Expect(err).NotTo(HaveOccurred())
			YType, err := r.NewType(&Y{}, "ys")
			Expect(err).NotTo(HaveOccurred())

			// maybe one can trigger the other?
			// or maybe only needs one, and scrap the other
			XType.HasMany(YType, "ys", "x_id", "id")
			YType.BelongsTo(XType, "x_id", "ys", "id")

			// creating instances
			x := X{
				ID:   "xID",
				Name: "some name",
			}
			XType.Create(x)
			YType.Create(Y{
				Attr: "some attr",
				XID:  x.ID,
			})

			// checking if y instance was created
			resultY, err := YType.Get(QueryAttribute{
				Tag:       "x_id",
				Value:     x.ID,
				Condition: Conditions.Equals,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(resultY)).To(Equal(1))
			Expect(resultY[0].(Y).Attr).To(Equal("some attr"))

			// checking if x instance was created
			result, err := XType.With("ys").Get(QueryAttribute{
				Tag:       "name",
				Value:     "some name",
				Condition: Conditions.Equals,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(result)).To(Equal(1))
			Expect(result[0].(X).Name).To(Equal("some name"))

			// checking if X.Ys[] was filled (Get() eagerly)
			Expect(len(result[0].(X).Ys)).To(Equal(1))
			Expect(result[0].(X).Ys[0].Attr).To(Equal("some attr"))
		})
	})
})
