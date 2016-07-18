package ini_test

import (
	"reflect"
	"testing"

	"github.com/MJKWoolnough/ini"
)

type structTest struct {
	Input  string
	Output interface{}
}

func testStruct(t *testing.T, tests []structTest) {
	for n, test := range tests {
		result := reflect.New(reflect.TypeOf(test.Output))
		err := ini.DecodeString(test.Input, result.Interface())
		if err != nil {
			t.Errorf("%s: test %d: unexpected error, %s", n+1, err)
		}
		if !reflect.DeepEqual(result.Elem().Interface(), test.Output) {
			t.Errorf("test %d: expecting %s, got %s", n+1, test.Output, result.Elem().Interface())
		}
	}
}

func TestStruct(t *testing.T) {
	type sliceValues struct {
		Vals []uint8 `ini:"Val,prefix"`
	}
	testStruct(t, []structTest{
		{
			Input:  "",
			Output: struct{}{},
		},
		{
			Input: "Test=123",
			Output: struct {
				Test uint
			}{
				Test: 123,
			},
		},
		{
			Input: "Test=-5\nBool=true",
			Output: struct {
				Test int
				Bool bool
			}{
				Test: -5,
				Bool: true,
			},
		},
		{
			Input: "Test=3.142\n[Section1]\nTest=A\n[Section2]\nTest=B\n[Section3]\nTest=C\n[Other]\nA=1\nB=2\n[Misc]\nZB=42\n[Other2]\nAB=CD",
			Output: struct {
				Test     float32
				Sections []struct {
					Test string
				} `ini:"Section,prefix"`
				M map[string]map[string]string `ini:",prefix"`
			}{
				Test: 3.142,
				Sections: []struct{ Test string }{
					{Test: "A"},
					{Test: "B"},
					{Test: "C"},
				},
				M: map[string]map[string]string{
					"Other":  map[string]string{"A": "1", "B": "2"},
					"Misc":   map[string]string{"ZB": "42"},
					"Other2": map[string]string{"AB": "CD"},
				},
			},
		},
		{
			Input: "[SectionA]\nVal1=1\nVal2=2\nVal3=3\n",
			Output: struct {
				SectionA sliceValues
			}{
				SectionA: sliceValues{
					Vals: []uint8{1, 2, 3},
				},
			},
		},
	})
}
