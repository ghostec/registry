package main

import "reflect"

// NewType registers a new type in a Registry
func (r *Registry) NewType(structure interface{}, cue StorageCue) (*Type, error) {
	t := &Type{
		structure:  structure,
		registry:   r,
		storageCue: cue,
	}
	name := reflect.TypeOf(structure).String()
	if err := r.storage.NewType(t); err != nil {
		return nil, err
	}
	r.types[name] = t
	return t, nil
}

// Create is a wrapper over t.registry.storage.Create
func (t *Type) Create(data interface{}) error {
	return t.registry.storage.Create(t, data)
}

// Get is a wrapper over t.registry.storage.Get
func (t *Type) Get(q ...QueryAttribute) ([]interface{}, error) {
	return t.registry.storage.Get(t, q...)
}

// AssociationType is a type alias
type AssociationType string

// AssociationTypes are all the associations supported by registry
var AssociationTypes = struct {
	HasMany   AssociationType
	BelongsTo AssociationType
}{
	HasMany:   "HasMany",
	BelongsTo: "BelongsTo",
}

// Association describes an association a *Type has with another *Type
type Association struct {
	Type AssociationType
	With *Type
}
