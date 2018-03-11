package main

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
	atype AssociationType
	with  *Type
	ref   string
}

func (t *Type) createAssociation(
	o *Type, ref string, atype AssociationType,
) error {
	t.associations = append(t.associations, Association{
		atype: atype,
		with:  o,
		ref:   ref,
	})
	return nil
}

func (t *Type) HasMany(o *Type, ref string) error {
	return t.createAssociation(o, ref, AssociationTypes.HasMany)
}

func (t *Type) BelongsTo(o *Type, ref string) error {
	return t.createAssociation(o, ref, AssociationTypes.BelongsTo)
}

func (t *Type) With(refs ...string) query {
	as := []Association{}
	for _, r := range refs {
		a := t.associationFromRef(r)
		if a.ref != r {
			continue
		}
		as = append(as, a)
	}
	return query{rt: t, associations: as}
}

func (t *Type) associationFromRef(r string) Association {
	for _, a := range t.associations {
		if a.ref == r {
			return a
		}
	}
	return Association{}
}
