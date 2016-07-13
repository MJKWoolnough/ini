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
	Delim            rune
	IgnoreTypeErrors bool
}

func newStruct(s reflect.Value, delim rune, ignoreTypeErrors bool) *vStruct {
	return &vStruct{
		Struct:           s,
		handler:          null{},
		Delim:            d.SubSectionDelim,
		IgnoreTypeErrors: ignoreTypeErros,
	}
}

func (vs *vStruct) Section(s string) {
	vs.Value.Close()
	vs.Value = null{}
	section := getSection(vs.Struct, s, string(vs.Delim))
	if section == nil {
		return
	}
	field := vs.Struct.FieldByIndex(section)
	sect := ""
	switch field.Kind() {
	case reflect.Map:
		if field.Type().Key().Kind() != reflect.String {
			return
		}
		switch sm := field.Type().Elem(); sm.Kind() {
		case reflect.Map: // map[string]map[string]string
			if sm.Key().Kind() != reflect.String || sm.Elem().Kind() != reflect.String {
				return
			}
			vs.handler = newMapMapString(field)
		case reflect.String: // map[string]string
			vs.handler = newMapString(field, vs.Delim)
		case reflect.Struct: //map[string]struct
			vs.handler = newMapStruct(field, vs.IgnoreTypeErrors)
		}
	case reflect.Slice:
		switch field.Type().Elem().Kind() {
		case reflect.Map: // []map[string]string
			//vs.handler = newSliceMapString(field, vs.IgnoreTypeErrors)
		case reflect.Struct: // []struct
			//vs.handler = newSliceStruct(field, vs.IgnoreTypeErrors)
		}
	case reflect.Struct:
		//vs.handler = newSStruct(field)
	}
	vs.handler.Section(sect)
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
