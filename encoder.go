package ini

import (
	"io"
	"reflect"
)

type mapKeys []reflect.Value

func (m mapKeys) Len() int {
	return len(m)
}

func (m mapKeys) Less(i, j int) bool {
	return m[i].String() < m[j].String()
}

func (m mapKeys) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type encoder struct {
	io.Writer
	written bool
	options
}

func Encode(w io.Writer, v interface{}, options ...Option) error {
	var e encoder
	for _, o := range options {
		o(&e.options)
	}
	e.Writer = w
	i := reflect.ValueOf(v)
	switch i.Kind() {
	case reflect.Map:
		if i.Type().Key().Kind() == reflect.String && i.IsValid() {
			switch i.Type().Elem().Kind() {
			case reflect.Map:
				if i.Type().Elem().Key().Kind() == reflect.String && i.Type().Elem().Elem().Kind() == reflect.String {
					return d.encodeMapMap(i)
				}
			case reflect.Slice:
				return d.encodeMapSlice(i)
			case reflect.String:
				return d.encodeMapString(i)
			case reflect.Struct:
				return d.encodeMapStruct(i)
			}
		}
	case reflect.Ptr:
		if i.Type().Elem().Kind() == reflect.Struct {
			if i.IsNil() {
				return ErrNilPointer
			}
			return d.encodeStruct(i.Elem())
		}
	case reflect.Struct:
		return d.encodeStruct(i)
	}
	return ErrInvalidType
}

func (e *encoder) WriteSection(s string) error {
	if e.written {
		if _, err := e.Writer.Write([]byte{'\n', '\n'}); err != nil {
			return err
		}
	}
	if _, err := e.Writer.Write([]byte{'['}); err != nil {
		return err
	}
	if _, err := e.Writer.Write([]byte(s)); err != nil {
		return err
	}
	if _, err := e.Writer.Write([]byte{']'}); err != nil {
		return err
	}
	e.written = true
	return nil
}

func (e *encoder) WriteKeyValue(k, v string) error {
	if e.written {
		if _, err := e.Writer.Write([]byte{'\n'}); err != nil {
			return err
		}
	}
	if _, err := e.Writer.Write([]byte(k)); err != nil {
		return err
	}
	if _, err := e.Writer.Write([]byte{e.NameValueDelim}); err != nil {
		return err
	}
	if _, err := e.Writer.Write([]byte(v)); err != nil {
		return err
	}
	e.written = true
}
