package main

import "reflect"

// Registry makes it easy to perform CRUD operations on any kind of
// struct, by registering custom types and exposing registry.Type
// from where to call Create, Get, Update and Delete
// These are performed by an implementation of the StorageEngine interface
type Registry struct {
	types   map[string]*Type
	storage StorageEngine
}

// New is a Registry ctor
func New(s StorageEngine) *Registry {
	return &Registry{
		storage: s,
		types:   map[string]*Type{},
	}
}

// NewType registers a new type in a Registry
func (r *Registry) NewType(structure interface{}, cue StorageCue) (*Type, error) {
	t := &Type{
		name:       reflect.TypeOf(structure).String(),
		structure:  structure,
		registry:   r,
		storageCue: cue,
		fieldTag:   map[string]string{},
		tagField:   map[string]string{},
	}

	st := reflect.Indirect(reflect.ValueOf(structure)).Type()
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		if tag, ok := f.Tag.Lookup("registry"); ok {
			t.fieldTag[f.Name] = tag
			t.tagField[tag] = f.Name
		}
	}

	if err := r.storage.NewType(t); err != nil {
		return nil, err
	}
	r.types[t.name] = t
	return t, nil
}

// QueryAttribute is used by StorageEngine to query instances of a Type
type QueryAttribute struct {
	Tag       string
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

func main() {}
