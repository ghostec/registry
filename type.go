package main

import (
	"fmt"
	"reflect"

	"github.com/mohae/deepcopy"
)

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
	fieldTag     map[string]string
	tagField     map[string]string
}

// Create is a wrapper over t.registry.storage.Create
func (t *Type) Create(data interface{}) error {
	cpy := reflect.Indirect(reflect.ValueOf(deepcopy.Copy(data))).Interface()
	return t.registry.storage.Create(t, cpy)
}

// Get is a wrapper over query.Get with lazy loading
func (t *Type) Get(q ...QueryAttribute) ([]interface{}, error) {
	return query{
		nestingType: queryNestingTypes.lazy,
		rt:          t,
	}.Get(q...)
}

// Eager is a wrapper over query.Get with eager loading
func (t *Type) Eager() query {
	return query{
		nestingType: queryNestingTypes.eager,
		rt:          t,
	}
}

// With is a wrapper over query.Get with custom loading
func (t *Type) With(selfRefs ...string) query {
	return query{
		nestingType: queryNestingTypes.custom,
		rt:          t,
		withNesting: selfRefs,
	}
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
	atype    AssociationType
	with     *Type
	selfRef  string
	otherRef string
	using    string
}

func (t *Type) createAssociation(
	o *Type, selfRef, otherRef, using string, atype AssociationType,
) error {
	t.associations = append(t.associations, Association{
		atype:    atype,
		with:     o,
		selfRef:  selfRef,
		otherRef: otherRef,
		using:    using,
	})
	return nil
}

func (t *Type) HasMany(o *Type, selfRef, otherRef, using string) error {
	return t.createAssociation(o, selfRef, otherRef, using, AssociationTypes.HasMany)
}

func (t *Type) BelongsTo(o *Type, selfRef, otherRef, using string) error {
	return t.createAssociation(o, selfRef, otherRef, using, AssociationTypes.BelongsTo)
}

func (t *Type) associationFromTag(tag string) (Association, error) {
	for _, a := range t.associations {
		if a.selfRef == tag {
			return a, nil
		}
	}
	return Association{}, fmt.Errorf("Association doesn't exist for tag %s", tag)
}
