package ini

import (
	"errors"
	"io"
	"reflect"
	"strings"

	"github.com/MJKWoolnough/parser"
)

type decoder struct {
	parser.Parser

	NameValueDelim                                                  rune
	SubSections                                                     bool
	SubSectionDelim                                                 rune
	DuplicateSecond, IgnoreQuotes, IgnoreTypeErrors, AllowMultiline bool
}

type Option func(*decoder)

func Decode(r io.Reader, v interface{}, options ...Option) error {
	return decode(parser.NewReaderTokeniser(r), v, options...)
}

func DecodeString(s string, v interface{}, options ...Option) error {
	return decode(parser.NewStringTokeniser(s), v, options...)
}

func DecoderBytes(b []byte, v interface{}, options ...Option) error {
	return decode(parser.NewByteTokeniser(b), v, options...)
}

func decode(t parser.Tokeniser, v interface{}, options ...Option) error {
	d := decoder{
		Parser:          parser.New(t),
		NameValueDelim:  '=',
		SubSectionDelim: '/',
	}
	for _, o := range options {
		o(&d)
	}
	d.TokeniserState(d.name)
	d.PhraserState(d.nameValue)
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return ErrInvalidKey
		}
		switch rv.Type().Elem().Kind() {
		case reflect.Map: // map[string]map[string]???
			switch rv.Type().Elem().Elem().Kind() {
			case reflect.String: //map[string]map[string]string
				var section string
				s := reflect.New(rv.Type().Key()).Elem()
				for {
					s.SetString(section)
					m := rv.MapIndex(s)
					if !m.IsValid() {
						m = reflect.MakeMap(rv.Type().Elem())
					}
					d.readMap(m, "")
					if m.Len() != 0 {
						rv.SetMapIndex(s, m)
					}
					p, _ := d.GetPhrase()
					if p.Type != phraseSection {
						if d.Err == io.EOF {
							return nil
						}
						return d.Err
					}
					section = p.Data[0].Data
				}
			default:
				return ErrInvalidMapType
			}
		case reflect.String: //map[string]string
			var prefix string
			for {
				d.readMap(rv, prefix)
				p, _ := d.GetPhrase()
				if p.Type != phraseSection {
					if d.Err == io.EOF {
						return nil
					}
					return d.Err
				}
				prefix = p.Data[0].Data + string(d.SubSectionDelim)
			}
		case reflect.Struct: //map[string]struct
		default:
			return ErrInvalidMapType
		}
	case reflect.Ptr:
		rv = rv.Elem()

		if rv.Kind() != reflect.Struct {
			return ErrInvalidType
		}

	default:
		return ErrInvalidType
	}

	return nil
}

func (d *decoder) readMap(m reflect.Value, prefix string) {
	k := reflect.New(m.Type().Key()).Elem()
	v := reflect.New(m.Type().Elem()).Elem()
	for d.Peek().Type == tokenName {
		p, _ := d.GetPhrase()
		if p.Type != phraseNameValue {
			break
		}
		k.SetString(prefix + p.Data[0].Data)
		v.SetString(p.Data[1].Data)
		m.SetMapIndex(k, v)
	}
}

func (d *decoder) readStruct(s reflect.Value) {
	t := s.Type()
	nf := t.NumField()
	for d.Peek().Type == tokenName {
		p, _ := d.GetPhrase()
		pn := p.Data[0].Data
		pv := p.Data[1].Data
		if p.Type != phraseNameValue {
			break
		}
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
			setValue(v, pv)
			reflect.Append(f, v)
		case reflect.Map:
			k := reflect.New(f.Type().Key()).Elem()
			if k.Kind() == reflect.String {
				v := reflect.New(f.Type().Elem()).Elem()
				setValue(v, pv)
				f.SetMapIndex(k, v)
			}
		default:
			v := reflect.New(f.Type()).Elem()
			err := setValue(v, pv)
			if err != nil && !d.IgnoreTypeErrors {
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
		case "TRUE", "T", "ON", "1":
			v.SetBool(true)
		default:
			v.SetBool(false)
		}
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.Complex64:
	case reflect.Complex128:
	}
	return err
}

// Errors
var (
	ErrInvalidType      = errors.New("needs map or pointer to struct type")
	ErrInvalidSliceType = errors.New("need slice of structs")
	ErrInvalidKey       = errors.New("maps require string keys")
	ErrInvalidMapType   = errors.New("invalid map type")
)
