package main

import (
	"reflect"
	"strings"
)

type query struct {
	rt           *Type
	attributes   []QueryAttribute
	associations []Association
}

func (q query) Get(qa ...QueryAttribute) ([]interface{}, error) {
	r, err := q.rt.Get(qa...)
	if err != nil {
		return nil, err
	}

	for _, a := range q.associations {
		for i, s := range r {
			nested, _ := a.with.Get(QueryAttribute{
				Tag:       a.otherRef,
				Value:     q.usingValueForAssociation(a, s),
				Condition: Conditions.Equals,
			})
			r[i] = newWithValueInFieldWithTag(q.rt, s, nested, a.selfRef)
		}
	}

	return r, nil
}

func newWithValueInFieldWithTag(
	t *Type, s interface{}, v []interface{}, ref string,
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
	f := t.tagField[ref]
	n.Elem().FieldByName(f).Set(sl)

	return n.Elem().Interface()
}

func (q query) usingValueForAssociation(
	a Association, s interface{},
) interface{} {
	t := reflect.ValueOf(s).Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if tag, ok := f.Tag.Lookup("registry"); ok {
			parts := strings.Split(tag, ",")
			for _, ps := range parts {
				if ps == a.using {
					return reflect.ValueOf(s).FieldByName(f.Name).Interface()
				}
			}
		}
	}
	return nil
}
