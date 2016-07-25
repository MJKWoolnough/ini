package ini

import (
	"reflect"
	"sort"
)

type mapStruct struct {
	Map reflect.Value
	sStruct
	Key              reflect.Value
	Changes          bool
	IgnoreTypeErrors bool
}

func newMapStruct(m reflect.Value, ignoreTypeErrors bool) *mapStruct {
	return &mapStruct{
		Map: m,
		Key: reflect.New(m.Type().Key()).Elem(),
		sStruct: sStruct{
			Struct: reflect.New(m.Type().Elem()).Elem(),
		},
		IgnoreTypeErrors: ignoreTypeErrors,
	}
}

func (m *mapStruct) Section(s string) {
	m.Close()
	m.Key.SetString(s)
	m.Struct = reflect.New(m.Map.Type().Elem()).Elem()
}

func (m *mapStruct) Set(k, v string) error {
	if err := m.sStruct.Set(k, v); err != nil {
		if !m.IgnoreTypeErrors {
			return err
		}
	} else {
		m.Changes = true
	}
	return nil
}

func (m *mapStruct) Close() {
	if m.Changes {
		m.Map.SetMapIndex(m.Key, m.Struct)
		m.Changes = false
	}
}

func (e *encoder) encodeMapStruct(m reflect.Value) error {
	keys := mapKeys(m.MapKeys())
	sort.Sort(keys)
	for _, key := range keys {
		k := key.String()
		if k != "" {
			if err := e.WriteSection(k); err != nil {
				return err
			}
		}
		if err := e.encodeSStruct(m.MapIndex(key)); err != nil {
			return err
		}
	}
	return nil
}
