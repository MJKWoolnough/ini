package ini

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

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

func (d *decoder) readStruct(s reflect.Value) {
	t := s.Type()
	nf := t.NumField()
Loop:
	for d.Peek().Type == tokenName {
		p, _ := d.GetPhrase()
		if p.Type != phraseNameValue {
			break
		}
		pn := p.Data[0].Data
		pv := p.Data[1].Data
		match := -1
		score := -1
		for i := 0; i < nf; i++ {
			f := t.Field(i)
			if f.PkgPath != "" {
				continue
			}
			tag := f.Tag.Get("ini")
			if tag == "" {
				tag = f.Name
			} else if tag[0] == ',' {
				tag = f.Name + tag
			}
			n, o := parseTag(tag)
			if n == pn {
				match = i
				break
			}
			if l := len(n); l > score && l >= len(pn) && o.Contains("prefix") && pn[:l] == n {
				score = l
				match = i
			}
		}
		if match < 0 {
			continue
		}
		f := s.Field(match)
		switch f.Kind() {
		case reflect.Slice:
			v := reflect.New(f.Type().Elem()).Elem()
			err := setValue(v, pv)
			if err == errUnknownType {
				continue Loop
			} else if err != nil && !d.IgnoreTypeErrors {
				d.Err = err
				return
			}
			reflect.Append(f, v)
		case reflect.Map:
			k := reflect.New(f.Type().Key()).Elem()
			if k.Kind() == reflect.String {
				v := reflect.New(f.Type().Elem()).Elem()
				err := setValue(v, pv)
				if err == errUnknownType {
					continue Loop
				} else if err != nil && !d.IgnoreTypeErrors {
					d.Err = err
					return
				}
				f.SetMapIndex(k, v)
			}
		default:
			v := reflect.New(f.Type()).Elem()
			err := setValue(v, pv)
			if err == errUnknownType {
				continue Loop
			} else if err != nil && !d.IgnoreTypeErrors {
				d.Err = err
				return
			}
			f.Set(v)
		}
	}
}

func setValue(v reflect.Value, pv string) (err error) {
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
		err = errUnknownType
	}
	return err
}

// Errors
var (
	ErrInvalidBool = errors.New("invalid boolean value")
	errUnknownType = errors.New("unknown type")
)
