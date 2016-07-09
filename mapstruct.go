package ini

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type mapStruct struct {
	Map              reflect.Value
	Key              reflect.Value
	Value            reflect.Value
	IgnoreTypeErrors bool
}

func (d *decoder) NewMapStruct(m reflect.Value) *mapStruct {
	return &mapStruct{
		Map:              m,
		Key:              reflect.New(m.Type().Key()).Elem(),
		Value:            reflect.New(m.Type().Elem()).Elem(),
		IgnoreTypeErrors: d.IgnoreTypeErrors,
	}
}

func (m *mapStruct) Section(s string) {
	m.Map.SetMapIndex(m.Key, m.Value)
	m.Key.SetString(s)
	m.Value = reflect.New(m.Map.Type().Elem()).Elem()
}

func (m *mapStruct) Set(k, v string) error {
	nf := matchField(m.Value, k)
	if nf < 0 {
		return nil
	}
	f := m.Value.Field(nf)
	switch f.Kind() {
	case reflect.Slice:
		e := reflect.New(f.Type().Elem()).Elem()
		err := setValue(e, v)
		if err == errUnknownType {
			return nil
		} else if err != nil && !m.IgnoreTypeErrors {
			return err
		}
		reflect.Append(f, e)
	case reflect.Map:
		mk := reflect.New(f.Type().Key()).Elem()
		if mk.Kind() == reflect.String {
			mv := reflect.New(f.Type().Elem()).Elem()
			err := setValue(mv, v)
			if err == errUnknownType {
				return nil
			} else if err != nil && !m.IgnoreTypeErrors {
				return err
			}
			f.SetMapIndex(mk, mv)
		}
	default:
		sv := reflect.New(f.Type()).Elem()
		err := setValue(sv, v)
		if err == errUnknownType {
			return nil
		} else if err != nil && !m.IgnoreTypeErrors {
			return err
		}
		f.Set(sv)
	}

	return nil
}

func (m *mapStruct) Close() {
	m.Map.SetMapIndex(m.Key, m.Value)
}

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
