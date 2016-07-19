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

func Encode(w io.Writer, v interface{}, options ...Option) error {
	var d decoder
	for _, o := range options {
		o(&d)
	}
	i := reflect.ValueOf(v)
	switch i.Kind() {
	case reflect.Map:
		if i.Type().Key().Kind() != reflect.String {

		}
		switch i.Type().Elem().Kind() {
		case reflect.Map:
			if i.Type().Elem().Key().Kind() == reflect.String && i.Type().Elem().Elem().Kind() == reflect.String {
				return d.encodeMapMap(i)
			}
		case reflect.Slice:
			return d.encodeMapSlice(i)
		case reflect.String:
			return d.encodeMapString
		case reflect.Struct:
			return d.encodeMapStruct(i)
		}
	case reflect.Ptr:
		if i.Type().Elem().Kind() != reflect.Struct {

		}
		i = i.Elem()
		return d.encodeStruct(i)
	case reflect.Struct:
		return d.encodeStruct(i)
	}
	return ErrInvalidType
}
