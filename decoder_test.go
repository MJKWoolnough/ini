package ini_test

import (
	"reflect"
	"testing"

	"github.com/MJKWoolnough/ini"
)

func TestMapStruct(t *testing.T) {
	type Test1 struct {
		A string
		C uint8
	}
	tests := []struct {
		Input  string
		Output interface{}
		Result interface{}
	}{
		{"A=B\nC=4", map[string]Test1{"": {"B", 4}}, map[string]Test1{}},
		{"[Section1]\nA=FG\nC=123", map[string]Test1{"": {"", 0}, "Section1": {"FG", 123}}, map[string]Test1{}},
		{"A=B\nC=4\n[Section1]\nA=FG\nC=123", map[string]Test1{"": {"B", 4}, "Section1": {"FG", 123}}, map[string]Test1{}},
	}

	for n, test := range tests {
		err := ini.DecodeString(test.Input, test.Result)
		if err != nil {
			t.Errorf("Test %d: unexpected error, %s", n+1, err)
			continue
		}
		if !reflect.DeepEqual(test.Result, test.Output) {
			t.Errorf("Test %d: result does not match expected", n+1)
		}
	}
}
