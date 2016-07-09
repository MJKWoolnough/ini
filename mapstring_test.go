package ini_test

import (
	"testing"

	"github.com/MJKWoolnough/ini"
)

func TestMapString(t *testing.T) {
	tests := []struct {
		Input  string
		Output map[string]string
	}{
		{
			"",
			map[string]string{},
		},
		{
			";Comment",
			map[string]string{},
		},
		{
			"#Comment",
			map[string]string{},
		},
		{
			"a=b",
			map[string]string{"a": "b"},
		},
		{
			"a=b\n",
			map[string]string{"a": "b"},
		},
		{
			"a=b\nc=d",
			map[string]string{"a": "b", "c": "d"},
		},
		{
			";Comment1\na=b\n#Comment2",
			map[string]string{"a": "b"},
		},
		{
			"a=b\n\nc=d",
			map[string]string{"a": "b", "c": "d"},
		},
		{
			"[EmptySection]",
			map[string]string{},
		},
		{
			"[Section1]\na=b",
			map[string]string{"Section1/a": "b"},
		},
		{
			"\n\n[Section1]\na=b",
			map[string]string{"Section1/a": "b"},
		},
		{
			"[Section1]\n\na=b",
			map[string]string{"Section1/a": "b"},
		},
		{
			"[Section1]\na=b",
			map[string]string{"Section1/a": "b"},
		},
		{
			"[Section1]\na=b\nc=d",
			map[string]string{"Section1/a": "b", "Section1/c": "d"},
		},
		{
			"[Section1]\na=b\nc=d\n\n[Section2]\ne=ff\ngggghdfg=4465746",
			map[string]string{"Section1/a": "b", "Section1/c": "d", "Section2/e": "ff", "Section2/gggghdfg": "4465746"},
		},
		{
			"#\n#Section 1 Follows\n#\n\n[Section1]\n;A=B/C\na=b\nc=d\n\n;\n#Section 2 Follow\n;\n[Section2]\ne=ff\ngggghdfg=4465746\n;Ending on a Comment",
			map[string]string{"Section1/a": "b", "Section1/c": "d", "Section2/e": "ff", "Section2/gggghdfg": "4465746"},
		},
	}

	for n, test := range tests {
		m := make(map[string]string)
		err := ini.DecodeString(test.Input, m)
		if err != nil {
			t.Errorf("Test %d: unexpected error: %s", n+1, err)
			continue
		}
		if len(m) != len(test.Output) {
			t.Errorf("Test %d: expecting %d elements, got %d", n+1, len(test.Output), len(m))
			continue
		}
		for k, v := range m {
			w, ok := test.Output[k]
			if !ok {
				t.Errorf("Test %d: key missing %q", n+1, k)
				continue
			}
			if v != w {
				t.Errorf("Test %d: key %q not correct, got %q, expecting %q", n+1, k, w, v)
			}
		}
	}
}

func TestMapStringCustomTypes(t *testing.T) {
	type A string
	type B string
	m := make(map[A]B)
	err := ini.DecodeString("1=2\n[SectionA]\nA=B\nC=D\n[SectionB]\nE=F", m)
	if err != nil {
		t.Errorf("Test: unexpected error: %s", err)
	}
	if len(m) != 4 {
		t.Errorf("Test: expecting 3 elements, got %d", len(m))
	}
	if m["1"] != "2" {
		t.Errorf("Test: expecting \"1\" == \"2\", got %q", m["1"])
	}
	if m["SectionA/A"] != "B" {
		t.Errorf("Test: expecting \"SectionA/A\" == \"B\", got %q", m["SectionA/A"])
	}
	if m["SectionA/C"] != "D" {
		t.Errorf("Test: expecting \"SectionA/C\" == \"D\", got %q", m["SectionA/C"])
	}
	if m["SectionB/E"] != "F" {
		t.Errorf("Test: expecting \"SectionB/E\" == \"F\", got %q", m["SectionB/E"])
	}
}
