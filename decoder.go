package ini

import (
	"errors"
	"io"
	"reflect"

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

	var i interface {
		Section(string)
		Set(string, string) error
		Close()
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return ErrInvalidKey
		}
		switch rv.Type().Elem().Kind() {
		case reflect.String:
			i = d.NewMapString(rv)
		case reflect.Struct:
		case reflect.Map:
			if rv.Type().Elem().Key().Kind() != reflect.String {
				return ErrInvalidKey
			}
			if rv.Type().Elem().Elem().Kind() != reflect.String {
				return ErrInvalidMapType
			}
			i = d.NewMapMapString(rv)
		}
	case reflect.Ptr:
	default:
		return ErrInvalidType
	}

	for {
		p, _ := d.GetPhrase()
		switch p.Type {
		case phraseSection:
			i.Section(p.Data[0].Data)
		case phraseNameValue:
			if err := i.Set(p.Data[0].Data, p.Data[1].Data); err != nil {
				if !d.IgnoreTypeErrors {
					return err
				}
			}
		case parser.PhraseDone:
			i.Close()
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
)
