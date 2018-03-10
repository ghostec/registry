package main

// Registry makes it easy to perform CRUD operations on any kind of
// struct, by registering custom types and exposing registry.Type
// from where to call Create, Get, Update and Delete
// These are performed by an implementation of the StorageEngine interface
type Registry struct {
	types   map[string]*Type
	storage StorageEngine
}

// Type is the registered entity that holds the instructureion needed about
// a custom struct in order to perform CRUD operation via registry.storage
type Type struct {
	name string
	// structure is an instance of the underlying struct
	structure interface{}
	registry  *Registry
	// storageCue is the information used by a StorageEngine implementation
	// to access the Type's instances (eg. table name in Postgres)
	storageCue   StorageCue
	associations []Association
}

// QueryAttribute is used by StorageEngine to query instances of a Type
type QueryAttribute struct {
	Field     string
	Value     interface{}
	Condition Condition
}

// Condition is a type alias
type Condition string

// StorageCue is a type alias
type StorageCue string

// StorageEngine is an interface that any storage engine needs to implement
// to be supported by a Registry
type StorageEngine interface {
	NewType(*Type) error
	Create(*Type, interface{}) error
	Get(*Type, ...QueryAttribute) ([]interface{}, error)
	// Update(*Type)
	// Delete(*Type)
}

// Conditions are all the supported conditions by StorageEngine queries
var Conditions = struct {
	Equals Condition
}{
	Equals: "equal",
}

// New is a Registry ctor
func New(s StorageEngine) *Registry {
	return &Registry{
		storage: s,
		types:   map[string]*Type{},
	}
}
