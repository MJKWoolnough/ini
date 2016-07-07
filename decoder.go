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
	d := decoder{
		Tokeniser:       parser.NewReaderTokeniser(r),
		NameValueDelim:  '=',
		SubSectionDelim: '/',
	}
	for _, o := range options {
		o(&d)
	}
	return d.decode(v)
}

func DecodeString(s string, v interface{}, options ...Option) error {
	d := decoder{
		Tokeniser:       parser.NewStringTokeniser(s),
		NameValueDelim:  '=',
		SubSectionDelim: '/',
	}
	for _, o := range options {
		o(&d)
	}
	return d.decode(v)
}

func DecoderBytes(b []byte, v interface{}, options ...Option) error {
	d := decoder{
		Tokeniser:       parser.NewByteTokeniser(b),
		NameValueDelim:  '=',
		SubSectionDelim: '/',
	}
	for _, o := range options {
		o(&d)
	}
	return d.decode(v)
}

func (d *decoder) decode(v interface{}) error {
	d.TokeniserState(d.name)
	rv := reflect.ValueOf(v)

	if k := rv.Kind(); k != reflect.Map {
		switch rv.Type().Elem().Kind() {
		case reflect.Map:
		case reflect.String:
		case reflect.Struct:
		}
	} else if k != reflect.Ptr {
		return ErrNotPointer
	}

	rv = rv.Elem()

	switch rv.Kind() {
	case reflect.Slice:
		if rv.Type().Elem().Kind() != reflect.Struct {
			return ErrInvalidSliceType
		}
	case reflect.Struct:
	default:
		return ErrInvalidType
	}

	return nil

	// if slice of struct type, read until a section and for each add to the slice

	// if struct type, read global section into struct
	// read sections into relevant structs
}

// Errors
var (
	ErrNotPointer       = errors.New("need pointer to type")
	ErrInvalidType      = errors.New("needs map, slice or struct type")
	ErrInvalidSliceType = errors.New("need slice of structs")
)
