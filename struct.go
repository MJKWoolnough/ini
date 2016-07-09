package ini

import "reflect"

type vStruct struct {
	Struct           reflect.Value
	SectionStruct    reflect.Value
	IgnoreTypeErrors bool
}

func (d *decoder) NewStruct(s reflect.Value) *vStruct {
	return &vStruct{
		Struct:           s,
		SectionStruct:    s,
		IgnoreTypeErrors: d.IgnoreTypeErrors,
	}
}

func (vs *vStruct) Section(s string) {

}

func (vs *vStruct) Set(k, v string) error {
	if !vs.SectionStruct.IsNil() {
		return nil
	}
	return nil
}

func (vs *vStruct) Close() {

}
