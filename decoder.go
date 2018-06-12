package ini // import "vimagination.zapto.org/ini"

import (
	"errors"
	"io"
	"reflect"

	"vimagination.zapto.org/parser"
)

type decoder struct {
	parser.Parser
	options
}

func Decode(r io.Reader, v interface{}, options ...Option) error {
	return decode(parser.NewReaderTokeniser(r), v, options...)
}

func DecodeString(s string, v interface{}, options ...Option) error {
	return decode(parser.NewStringTokeniser(s), v, options...)
}

func DecoderBytes(b []byte, v interface{}, options ...Option) error {
	return decode(parser.NewByteTokeniser(b), v, options...)
}

func decode(t parser.Tokeniser, v interface{}, os ...Option) error {
	d := decoder{
		Parser: parser.New(t),
		options: options{
			NameValueDelim:  '=',
			SubSectionDelim: '/',
		},
	}

	for _, o := range os {
		o(&d.options)
	}

	d.TokeniserState(d.name)
	d.PhraserState(d.nameValue)

	var h handler

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return ErrInvalidKey
		}
		switch rv.Type().Elem().Kind() {
		case reflect.String:
			h = newMapString(rv, d.SubSectionDelim)
		case reflect.Struct:
			h = newMapStruct(rv, d.IgnoreTypeErrors)
		case reflect.Map:
			if rv.Type().Elem().Key().Kind() != reflect.String {
				return ErrInvalidKey
			}
			if rv.Type().Elem().Elem().Kind() != reflect.String {
				return ErrInvalidMapType
			}
			h = newMapMapString(rv)
		}
	case reflect.Ptr:
		if rv.Elem().Kind() != reflect.Struct {
			return ErrInvalidType
		}
		if rv.IsNil() {
			return ErrNilPointer
		}
		h = newStruct(rv.Elem(), d.SubSectionDelim, d.IgnoreTypeErrors)
	default:
		return ErrInvalidType
	}

	for {
		p, _ := d.GetPhrase()
		switch p.Type {
		case phraseSection:
			h.Section(p.Data[0].Data)
		case phraseNameValue:
			if err := h.Set(p.Data[0].Data, p.Data[1].Data); err != nil && !d.IgnoreTypeErrors && err != errUnknownType {
				return err
			}
		case parser.PhraseDone:
			h.Close()
			return nil
		default:
			return d.Err
		}
	}
}

// Errors
var (
	ErrInvalidType      = errors.New("needs map or pointer to struct type")
	ErrInvalidSliceType = errors.New("need slice of structs")
	ErrInvalidKey       = errors.New("maps require string keys")
	ErrInvalidMapType   = errors.New("invalid map type")
	ErrNilPointer       = errors.New("nil pointer to struct")
)
