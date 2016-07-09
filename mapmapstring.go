package ini

import "reflect"

type mapMapString struct {
	Map, MapMap       reflect.Value
	KeyA, KeyB, Value reflect.Value
}

func (d *decoder) NewMapMapString(m reflect.Value) *mapMapString {
	ka := reflect.New(m.Type().Key()).Elem()
	kb := reflect.New(m.Type().Elem().Key()).Elem()
	v := reflect.New(m.Type().Elem().Elem()).Elem()
	mm := reflect.MakeMap(m.Type().Elem())
	return &mapMapString{
		Map:    m,
		MapMap: mm,
		KeyA:   ka,
		KeyB:   kb,
		Value:  v,
	}
}

func (m *mapMapString) Section(s string) {
	m.Close()
	m.MapMap = reflect.MakeMap(m.Map.Type().Elem())
	m.KeyA.SetString(s)
}

func (m *mapMapString) Set(k, v string) error {
	m.KeyB.SetString(k)
	m.Value.SetString(v)
	m.MapMap.SetMapIndex(m.KeyB, m.Value)
	return nil
}

func (m *mapMapString) Close() {
	if m.MapMap.Len() > 0 {
		m.Map.SetMapIndex(m.KeyA, m.MapMap)
	}
}
