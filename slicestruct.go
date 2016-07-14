package ini

import "reflect"

type sliceStruct struct {
	Slice reflect.Value
	*sStruct
	Changes bool
}

func newSliceStruct(s reflect.Value) *sliceStruct {
	return &sliceStruct{
		Slice:   s,
		sStruct: newSStruct(reflect.New(s.Type().Elem()).Elem()),
	}
}

func (ss *sliceStruct) Section(s string) {
	if ss.Changes {
		ss.Close()
		ss.sStruct = newSStruct(reflect.New(ss.Slice.Type().Elem()).Elem())
	}
	ss.sStruct.Section(s)
}

func (ss *sliceStruct) Close() {
	if ss.Changes {
		reflect.Append(ss.Slice, ss.Struct)
		ss.Changes = false
	}
}
