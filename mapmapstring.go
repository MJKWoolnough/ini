package ini

import (
	"reflect"
	"sort"
)

type mapMapString struct {
	Map, MapMap       reflect.Value
	KeyA, KeyB, Value reflect.Value
}

func newMapMapString(m reflect.Value) *mapMapString {
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

func (e *encoder) encodeMapMap(m reflect.Value) error {
	keys := mapKeys(m.MapKeys())
	sort.Sort(keys)
	for _, key := range keys {
		k := key.String()
		if k != "" {
			if err := e.WriteSection(k); err != nil {
				return err
			}
		}
		mv := m.MapIndex(key)
		mvk := mapKeys(mv.MapKeys())
		sort.Sort(mvk)
		for _, vk := range mvk {
			if err := e.WriteKeyValue(vk.String(), mv.MapIndex(vk).String()); err != nil {
				return err
			}
		}
	}
	return nil
}
