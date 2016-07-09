package ini

import "reflect"

type mapString struct {
	Map     reflect.Value
	Key     reflect.Value
	Value   reflect.Value
	Delim   rune
	Section string
}

func (d *decoder) NewMapString(m reflect.Value) *mapString {
	return &mapString{
		Map:     m,
		Key:     reflect.New(m.Type().Key()).Elem(),
		Value:   reflect.New(m.Type().Elem()).Elem(),
		Delim:   d.SubSectionDelim,
		Section: "",
	}
}

func (m *mapString) Section(s string) {
	m.Section = s + string(m.delim)
}

func (m *mapString) Set(k, v string) error {
	m.Key.SetString(m.Section + k)
	m.Value.SetString(v)
	m.Map.SetMapIndex(m.Key, m.Value)
	return nil
}

func (m *mapString) Close() {
}
