package ini

import "reflect"

type mapString struct {
	Map     reflect.Value
	Key     reflect.Value
	Value   reflect.Value
	Delim   rune
	section string
}

func (d *decoder) NewMapString(m reflect.Value) *mapString {
	return &mapString{
		Map:   m,
		Key:   reflect.New(m.Type().Key()).Elem(),
		Value: reflect.New(m.Type().Elem()).Elem(),
		Delim: d.SubSectionDelim,
	}
}

func (m *mapString) Section(s string) {
	m.section = s + string(m.Delim)
}

func (m *mapString) Set(k, v string) error {
	m.Key.SetString(m.section + k)
	m.Value.SetString(v)
	m.Map.SetMapIndex(m.Key, m.Value)
	return nil
}

func (m *mapString) Close() {
}
