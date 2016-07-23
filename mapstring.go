package ini

import (
	"reflect"
	"sort"
	"strings"
)

type mapString struct {
	Map     reflect.Value
	Key     reflect.Value
	Value   reflect.Value
	Delim   rune
	section string
}

func newMapString(m reflect.Value, delim rune) *mapString {
	return &mapString{
		Map:   m,
		Key:   reflect.New(m.Type().Key()).Elem(),
		Value: reflect.New(m.Type().Elem()).Elem(),
		Delim: delim,
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

func (e *encoder) encodeMapString(m reflect.Value) error {
	keys := mapKeys(m.MapKeys())
	sort.Sort(keys)
	var last string
	for _, key := range keys {
		k := key.String()
		v := m.MapIndex(key)
		var section string
		p := strings.LastIndexAny(key, string(e.SubSectionDelim))
		if p >= 0 {
			section, k = k[:p], k[p+1:]
		}
		if section != last {
			e.WriteSection(section)
			last = section
		}
		e.WriteKeyValue(k, v)
	}
	return nil
}
