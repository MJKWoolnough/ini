package ini

import (
	"errors"
	"reflect"
	"strconv"
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
		Delim:            delim,
		IgnoreTypeErrors: ignoreTypeErrors,
	}
}

func (vs *vStruct) Section(s string) {
	vs.Close()
	vs.handler = null{}
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
			if field.Type().Key().Kind() != reflect.String || field.Type().Elem().Kind() != reflect.String {
				return
			}
			vs.handler = newSliceMapString(field, vs.Delim)
		case reflect.Struct: // []struct
			vs.handler = newSliceStruct(field)
		}
	case reflect.Struct:
		vs.handler = newSStruct(field)
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

type sStruct struct {
	Struct reflect.Value
}

func newSStruct(s reflect.Value) *sStruct {
	return &sStruct{s}
}

func (sStruct) Section(_ string) {}

func (s *sStruct) Set(k, v string) error {
	nf := matchField(s.Struct, k)
	if nf < 0 {
		return errUnknownType
	}
	f := s.Struct.Field(nf)
	switch f.Kind() {
	case reflect.Slice:
		e := reflect.New(f.Type().Elem()).Elem()
		if err := setValue(e, v); err != nil {
			return err
		}
		reflect.Append(f, e)
	case reflect.Map:
		mk := reflect.New(f.Type().Key()).Elem()
		if mk.Kind() == reflect.String {
			mv := reflect.New(f.Type().Elem()).Elem()
			if err := setValue(mv, v); err != nil {
				return err
			}
			f.SetMapIndex(mk, mv)
		}
	default:
		sv := reflect.New(f.Type()).Elem()
		if err := setValue(sv, v); err != nil {
			return err
		}
		f.Set(sv)
	}
	return nil
}

func (sStruct) Close() {}

func matchField(v reflect.Value, name string) int {
	match := -1
	score := -1
	nf := v.NumField()
	for i := 0; i < nf; i++ {
		f := v.Type().Field(i)
		if f.PkgPath != "" {
			continue
		}
		tag := f.Tag.Get("ini")
		n, o := parseTag(tag)
		if n == "" {
			n = f.Name
		}
		if n == name {
			match = i
			break
		}
		if l := len(n); l > score && l >= len(name) && o.Contains("prefix") && name[:l] == n {
			score = l
			match = i
		}
	}
	return match
}

func setValue(v reflect.Value, pv string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(pv)
	case reflect.Bool:
		switch strings.ToUpper(pv) {
		case "TRUE", "T", "ON", "1", "YES", "Y":
			v.SetBool(true)
		case "FALSE", "F", "OFF", "0", "NO", "N":
			v.SetBool(false)
		default:
			return ErrInvalidBool
		}
	case reflect.Uint:
		n, err := strconv.ParseUint(pv, 0, 0)
		if err != nil {
			return err
		}
		v.SetUint(n)
	case reflect.Uint8:
		n, err := strconv.ParseUint(pv, 0, 8)
		if err != nil {
			return err
		}
		v.SetUint(n)
	case reflect.Uint16:
		n, err := strconv.ParseUint(pv, 0, 16)
		if err != nil {
			return err
		}
		v.SetUint(n)
	case reflect.Uint32:
		n, err := strconv.ParseUint(pv, 0, 32)
		if err != nil {
			return err
		}
		v.SetUint(n)
	case reflect.Uint64:
		n, err := strconv.ParseUint(pv, 0, 64)
		if err != nil {
			return err
		}
		v.SetUint(n)
	case reflect.Int:
		n, err := strconv.ParseInt(pv, 0, 0)
		if err != nil {
			return err
		}
		v.SetInt(n)
	case reflect.Int8:
		n, err := strconv.ParseInt(pv, 0, 8)
		if err != nil {
			return err
		}
		v.SetInt(n)
	case reflect.Int16:
		n, err := strconv.ParseInt(pv, 0, 16)
		if err != nil {
			return err
		}
		v.SetInt(n)
	case reflect.Int32:
		n, err := strconv.ParseInt(pv, 0, 32)
		if err != nil {
			return err
		}
		v.SetInt(n)
	case reflect.Int64:
		n, err := strconv.ParseInt(pv, 0, 64)
		if err != nil {
			return err
		}
		v.SetInt(n)
	case reflect.Float32:
		n, err := strconv.ParseFloat(pv, 32)
		if err != nil {
			return err
		}
		v.SetFloat(n)
	case reflect.Float64:
		n, err := strconv.ParseFloat(pv, 64)
		if err != nil {
			return err
		}
		v.SetFloat(n)
	default:
		return errUnknownType
	}
	return nil
}

// Errors
var (
	ErrInvalidBool = errors.New("invalid boolean value")
	errUnknownType = errors.New("unknown type")
)
