package main

import (
	"reflect"
)

// MemoryStorage stores types data in memory and exposes methods to
// perform CRUD operations over them
type MemoryStorage struct {
	data map[string][]interface{}
}

// NewMemoryStorage ctor
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: map[string][]interface{}{},
	}
}

// NewType registers a storage for a new type in Registry
func (s *MemoryStorage) NewType(t *Type) error {
	s.data[t.name] = []interface{}{}
	return nil
}

// Create method
func (s *MemoryStorage) Create(t *Type, data interface{}) error {
	s.data[t.name] = append(s.data[t.name], data)
	return nil
}

// Get method
func (s *MemoryStorage) Get(
	t *Type, q ...QueryAttribute,
) ([]interface{}, error) {
	r := []interface{}{}

	for _, x := range s.data[t.name] {
		if s.applyGetFilter(t, x, q) {
			r = append(r, x)
		}
	}
	return r, nil
}

func (s *MemoryStorage) applyGetFilter(
	t *Type, x interface{}, q []QueryAttribute,
) bool {
	for _, qq := range q {
		f := t.tagField[qq.Tag]
		v := reflect.ValueOf(x).FieldByName(f).Interface()
		switch qq.Condition {
		case Conditions.Equals:
			if !(qq.Value == v) {
				return false
			}
		}
	}
	return true
}
