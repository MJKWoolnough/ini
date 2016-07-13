package ini

import "reflect"

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
