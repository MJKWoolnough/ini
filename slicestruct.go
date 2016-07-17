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
		sStruct: newSStruct(reflect.New(s.Type().Elem().Elem()).Elem()),
	}
}

func (ss *sliceStruct) Section(s string) {
	if ss.Changes {
		ss.Close()
		ss.sStruct = newSStruct(reflect.New(ss.Slice.Type().Elem().Elem()).Elem())
	}
	ss.sStruct.Section(s)
}

func (ss *sliceStruct) Set(k, v string) error {
	err := ss.sStruct.Set(k, v)
	if err == nil {
		ss.Changes = true
	}
	return err
}

func (ss *sliceStruct) Close() {
	if ss.Changes {
		ss.Slice.Elem().Set(reflect.Append(ss.Slice.Elem(), ss.Struct))
		ss.Changes = false
	}
}
