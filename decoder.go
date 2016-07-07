package ini

import (
	"errors"
	"io"
	"reflect"

	"github.com/MJKWoolnough/parser"
)

type decoder struct {
	parser.Tokeniser

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
		Tokeniser:       t,
		NameValueDelim:  '=',
		SubSectionDelim: '/',
	}
	for _, o := range options {
		o(&d)
	}
	d.TokeniserState(d.name)
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
			default:
				return ErrInvalidMapType
			}
		case reflect.String: //map[string]string
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

// Errors
var (
	ErrInvalidType      = errors.New("needs map or pointer to struct type")
	ErrInvalidSliceType = errors.New("need slice of structs")
	ErrInvalidKey       = errors.New("maps require string keys")
	ErrInvalidMapType   = errors.New("invalid map type")
)
