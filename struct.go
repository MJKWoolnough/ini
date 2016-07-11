package ini

import (
	"reflect"
	"strings"
)

type handler interface {
	Section(string)
	Set(string, string) error
	Close()
}

type vStruct struct {
	Struct reflect.Value
	handler
	Delim string
}

func (d *decoder) NewStruct(s reflect.Value) *vStruct {
	return &vStruct{
		Struct:  s,
		handler: null{},
		Delim:   string(d.SubSectionDelim),
	}
}

func (vs *vStruct) Section(s string) {
	section := getSection(vs.Struct, s, vs.Delim)
	if section == nil {
		return
	}
}

func getSection(s reflect.Value, section, delim string) []int {
	if s.Kind() != reflect.Struct {
		return nil
	}
	parts := strings.SplitN(section, delim, 2)
	pos := matchField(s, parts[0])
	if pos >= 0 {
		if len(parts) == 1 {
			return []int{pos}
		}
		v := getSection(s.Field(pos), parts[1], delim)
		if v != nil {
			w := make([]int, 1, len(v)+1)
			w[0] = pos
			return append(w, v...)
		}
	}
	if len(parts) == 1 {
		return nil
	}
	pos = matchField(s, section)
	if pos < 0 {
		return nil
	}
	return []int{pos}
}

type null struct{}

func (null) Section(string) {}

func (null) Set(string, string) error {
	return nil
}

func (null) Close() {}
