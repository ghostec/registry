package main

import (
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
	atype   AssociationType
	from    *Type
	with    *Type
	fromTag string
	withTag string
	using   string
}

func (t *Type) createAssociation(
	w *Type, fromTag, withTag, using string, atype AssociationType,
) error {
	t.associations = append(t.associations, Association{
		atype:   atype,
		from:    t,
		with:    w,
		fromTag: fromTag,
		withTag: withTag,
		using:   using,
	})
	return nil
}

// HasMany creates an association (t has_many w)
func (t *Type) HasMany(w *Type, fromTag, withTag, using string) error {
	return t.createAssociation(w, fromTag, withTag, using, AssociationTypes.HasMany)
}

// BelongsTo creates an association (t belongs_to w)
func (t *Type) BelongsTo(w *Type, fromTag, withTag, using string) error {
	return t.createAssociation(w, fromTag, withTag, using, AssociationTypes.BelongsTo)
}
