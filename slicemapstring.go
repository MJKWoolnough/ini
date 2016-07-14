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
		mapString: newMapString(reflect.MakeMap(s.Type().Elem()), delim),
	}
}

func (sm *sliceMapString) Section(s string) {
	if sm.Changes {
		sm.Close()
		sm.mapString = newMapString(reflect.MakeMap(sm.Slice.Type().Elem()), sm.Delim)
	}
	sm.mapString.Section(s)
}

func (sm *sliceMapString) Close() {
	if sm.Changes {
		if sm.Map.Len() > 0 {
			reflect.Append(sm.Slice, sm.Map)
		}
		sm.Changes = false
	}
}
