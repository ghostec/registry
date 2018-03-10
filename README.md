Type Registry library that accepts StorageEngine implementations to perform CRUD operations over custom structs

## Usage:

```go
import "github.com/ghostec/registry"

s := registry.NewMemoryStorage()
r := registry.New(s)

type Book struct {
	Name string `registry:"name"`
}

// register a custom type using an empty instance and a storage cue
BookType := r.NewType(Book{}, "books")
BookType.Create(Book{Name: "Alice in the rabbit's hole"})
results, _ := BookType.Get(&QueryAttribute{
	Field: "name",
	Value: "Alice in the rabbit's hole",
	Condition: registry.Conditions.Equals,
})
// results = []interface{}{Book{Name: "Alice in the rabbit's hole"}}
book = results[0].(Book)
```
