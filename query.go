package main

import (
	"reflect"
)

type queryNestingType string

var queryNestingTypes = struct {
	eager  queryNestingType
	lazy   queryNestingType
	custom queryNestingType
}{
	eager:  "eager",
	lazy:   "lazy",
	custom: "custom",
}

type query struct {
	nestingType queryNestingType
	rt          *Type
}

func (q query) Get(qa ...QueryAttribute) ([]interface{}, error) {
	r, err := q.rt.registry.storage.Get(q.rt, qa...)
	if err != nil {
		return nil, err
	}
	if q.nestingType == queryNestingTypes.lazy {
		return r, nil
	}

	var associations []Association
	if q.nestingType == queryNestingTypes.eager {
		associations = q.rt.associations
	}

	for _, a := range associations {
		for i, s := range r {
			nested, _ := a.with.Eager().Get(QueryAttribute{
				Tag:       a.withTag,
				Value:     q.usingValueForAssociation(a, s),
				Condition: Conditions.Equals,
			})
			r[i] = fillValueInFieldWithTag(q.rt, s, nested, a.fromTag)
		}
	}

	return r, nil
}

func fillValueInFieldWithTag(
	t *Type, s interface{}, v []interface{}, fromTag string,
) interface{} {
	if len(v) == 0 {
		return s
	}

	// reflect.Value ptr instance with s' type
	sType := reflect.Indirect(reflect.ValueOf(s)).Type()
	n := reflect.New(sType)
	reflect.Indirect(n).Set(reflect.ValueOf(s))

	// v's slice with original type
	vType := reflect.ValueOf(v[0]).Type()
	sl := reflect.MakeSlice(reflect.SliceOf(vType), 0, len(v))
	for _, e := range v {
		sl = reflect.Append(sl, reflect.ValueOf(e))
	}

	// search field in n with ref tag in registry and set it
	f := t.tagField[fromTag]
	n.Elem().FieldByName(f).Set(sl)

	return n.Elem().Interface()
}

func (q query) usingValueForAssociation(
	a Association, s interface{},
) interface{} {
	f := a.from.tagField[a.using]
	return reflect.ValueOf(s).FieldByName(f).Interface()
}
