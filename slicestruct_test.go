package ini_test

import (
	"reflect"
	"testing"

	"github.com/MJKWoolnough/ini"
)

func TestSliceStruct(t *testing.T) {
	type testStruct struct {
		Key string
		ZB  string
	}
	tests := []struct {
		Input  string
		Output interface{}
	}{
		{
			Input: "[Section1]\nKey=Value\nZB=42\n[Section2]\nKey=Master\nZB=YA",
			Output: struct {
				Sections []testStruct `ini:"Section,prefix"`
			}{
				Sections: []testStruct{
					{
						Key: "Value",
						ZB:  "42",
					},
					{
						Key: "Master",
						ZB:  "YA",
					},
				},
			},
		},
	}

	for n, test := range tests {
		result := reflect.New(reflect.TypeOf(test.Output))
		err := ini.DecodeString(test.Input, result.Interface())
		if err != nil {
			t.Errorf("test %d: unexpected error, %s", n+1, err)
		}
		if !reflect.DeepEqual(result.Elem().Interface(), test.Output) {
			t.Errorf("test %d: expecting %s, got %s", n+1, test.Output, result.Elem().Interface())
		}
	}
}
