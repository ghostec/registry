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
	// TODO: cache structure fields/tags and relations
	registry *Registry
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

func (t *Type) With(selfRefs ...string) query {
	as := []Association{}
	for _, r := range selfRefs {
		a := t.associationFromSelfRef(r)
		if a.selfRef != r {
			continue
		}
		as = append(as, a)
	}
	return query{rt: t, associations: as}
}

func (t *Type) associationFromSelfRef(r string) Association {
	for _, a := range t.associations {
		if a.selfRef == r {
			return a
		}
	}
	return Association{}
}
