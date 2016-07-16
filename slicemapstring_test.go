package ini_test

import (
	"reflect"
	"testing"

	"github.com/MJKWoolnough/ini"
)

func TestSliceMapString(t *testing.T) {
	tests := []struct {
		Input  string
		Output interface{}
	}{
		{
			"[Section]\nTest=1",
			struct {
				Sections []map[string]string `ini:",prefix"`
			}{
				Sections: []map[string]string{
					map[string]string{
						"Section/Test": "1",
					},
				},
			},
		},
		{
			"[Section1]\nTest=1\nHELLO=WORLD\n[Section2]\nP=NP",
			struct {
				Sections []map[string]string `ini:",prefix"`
			}{
				Sections: []map[string]string{
					map[string]string{
						"Section1/Test":  "1",
						"Section1/HELLO": "WORLD",
					},
					map[string]string{
						"Section2/P": "NP",
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
