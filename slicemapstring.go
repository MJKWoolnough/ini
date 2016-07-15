package ini

import "reflect"

type sliceMapString struct {
	Slice reflect.Value
	*mapString
	Delim   rune
	Changes bool
}

func newSliceMapString(s reflect.Value, delim rune) *sliceMapString {
	return &sliceMapString{
		Slice:     s,
		mapString: newMapString(reflect.MakeMap(s.Type().Elem().Elem()), delim),
	}
}

func (sm *sliceMapString) Section(s string) {
	if sm.Changes {
		sm.Close()
		sm.mapString = newMapString(reflect.MakeMap(sm.Slice.Type().Elem().Elem()), sm.Delim)
	}
	sm.mapString.Section(s)
}

func (sm *sliceMapString) Set(k, v string) error {
	sm.Changes = true
	return sm.mapString.Set(k, v)
}

func (sm *sliceMapString) Close() {
	if sm.Changes {
		if sm.Map.Len() > 0 {
			sm.Slice.Elem().Set(reflect.Append(sm.Slice.Elem(), sm.Map))
		}
		sm.Changes = false
	}
}
