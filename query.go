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
	shallow, err := q.rt.Get(qa...)
	if err != nil {
		return nil, err
	}

	for _, a := range q.associations {
		for _, s := range shallow {
			r, _ := a.with.Get(QueryAttribute{
				Field:     a.otherRef,
				Value:     q.usingValueForAssociation(a, s),
				Condition: Conditions.Equals,
			})
			setValueInFieldWithTag(s, r, a.selfRef)
		}
	}

	return shallow, nil
}

func setValueInFieldWithTag(s interface{}, v []interface{}, ref string) {
	t := reflect.ValueOf(s).Elem().Type()
	// check when v.len == 0
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if tag, ok := f.Tag.Lookup("registry"); ok {
			if tag == ref {
				vType := reflect.ValueOf(v[0]).Elem().Type()
				sl := reflect.MakeSlice(reflect.SliceOf(vType), 0, len(v))
				for _, e := range v {
					sl = reflect.Append(sl, reflect.ValueOf(e).Elem())
				}
				reflect.ValueOf(s).Elem().FieldByName(f.Name).Set(sl)
				return
			}
		}
	}
}

func (q query) usingValueForAssociation(
	a Association, s interface{},
) interface{} {
	t := reflect.ValueOf(s).Elem().Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if tag, ok := f.Tag.Lookup("registry"); ok {
			parts := strings.Split(tag, ",")
			for _, ps := range parts {
				if ps == a.using {
					return reflect.ValueOf(s).Elem().FieldByName(f.Name).Interface()
				}
			}
		}
	}
	return nil
}
