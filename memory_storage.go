package main

import "reflect"

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
		if s.applyGetFilter(x, q) {
			r = append(r, x)
		}
	}
	return r, nil
}

func (s *MemoryStorage) applyGetFilter(
	x interface{}, q []QueryAttribute,
) bool {
	t := reflect.Indirect(reflect.ValueOf(x)).Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		for _, qq := range q {
			if tag, ok := f.Tag.Lookup("registry"); ok && tag == qq.Field {
				v := reflect.ValueOf(x).Elem().FieldByName(f.Name).Interface()
				switch qq.Condition {
				case Conditions.Equals:
					return qq.Value == v
				}
			}
		}
	}
	return false
}
