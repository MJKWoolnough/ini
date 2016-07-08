package ini_test

import (
	"reflect"
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

func TestMapMapString(t *testing.T) {
	tests := []struct {
		Input  string
		Output map[string]map[string]string
	}{
		{
			"",
			map[string]map[string]string{},
		},
		{
			";Comment",
			map[string]map[string]string{},
		},
		{
			"#Comment",
			map[string]map[string]string{},
		},
		{
			"a=b",
			map[string]map[string]string{
				"": map[string]string{"a": "b"},
			},
		},
		{
			"a=b\n",
			map[string]map[string]string{
				"": map[string]string{"a": "b"},
			},
		},
		{
			"a=b\nc=d",
			map[string]map[string]string{
				"": map[string]string{"a": "b", "c": "d"},
			},
		},
		{
			";Comment1\na=b\n#Comment2",
			map[string]map[string]string{
				"": map[string]string{"a": "b"},
			},
		},
		{
			"a=b\n\nc=d",
			map[string]map[string]string{
				"": map[string]string{"a": "b", "c": "d"},
			},
		},
		{
			"[EmptySection]",
			map[string]map[string]string{},
		},
		{
			"[Section1]\na=b",
			map[string]map[string]string{
				"Section1": map[string]string{"a": "b"},
			},
		},
		{
			"\n\n[Section1]\na=b",
			map[string]map[string]string{
				"Section1": map[string]string{"a": "b"},
			},
		},
		{
			"[Section1]\n\na=b",
			map[string]map[string]string{
				"Section1": map[string]string{"a": "b"},
			},
		},
		{
			"[Section1]\na=b",
			map[string]map[string]string{
				"Section1": map[string]string{"a": "b"},
			},
		},
		{
			"[Section1]\na=b\nc=d",
			map[string]map[string]string{
				"Section1": map[string]string{"a": "b", "c": "d"},
			},
		},
		{
			"[Section1]\na=b\nc=d\n\n[Section2]\ne=ff\ngggghdfg=4465746",
			map[string]map[string]string{
				"Section1": map[string]string{"a": "b", "c": "d"},
				"Section2": map[string]string{"e": "ff", "gggghdfg": "4465746"},
			},
		},
		{
			"#\n#Section 1 Follows\n#\n\n[Section1]\n;A=B/C\na=b\nc=d\n\n;\n#Section 2 Follow\n;\n[Section2]\ne=ff\ngggghdfg=4465746\n;Ending on a Comment",
			map[string]map[string]string{
				"Section1": map[string]string{"a": "b", "c": "d"},
				"Section2": map[string]string{"e": "ff", "gggghdfg": "4465746"},
			},
		},
	}

	for n, test := range tests {
		m := make(map[string]map[string]string)
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
			if len(v) != len(w) {
				t.Errorf("Test %d: expecting %d elements in key %q, got %d", n+1, len(w), k, len(v))
			}
			for i, x := range v {
				y, ok := w[i]
				if !ok {
					t.Errorf("Test %d: key missing in %q, %q", n+1, k, i)
				}
				if x != y {
					t.Errorf("Test %d: key %q in %q not correct, expecting %q, got %q", n+1, i, k, x, y)
				}
			}
		}
	}
}

func TestMapMapStringCustomTypes(t *testing.T) {
	type A string
	type B string
	type C string
	m := make(map[A]map[B]C)
	err := ini.DecodeString("1=2\n[SectionA]\nA=B\nC=D\n[SectionB]\nE=F", m)
	if err != nil {
		t.Errorf("Test: unexpected error: %s", err)
	}
	if len(m) != 3 {
		t.Errorf("Test: expecting 3 elements, got %d", len(m))
	}
	if n, ok := m[""]; !ok {
		t.Errorf("Test: missing default section")
	} else if len(n) != 1 {
		t.Errorf("Test: expecting 1 element in default section, got %d", len(n))
	} else {
		if n["1"] != "2" {
			t.Errorf("Test: expecting \"1\" == \"2\", got %q", n["1"])
		}
	}
	if a, ok := m["SectionA"]; !ok {
		t.Errorf("Test: missing SectionA")
	} else if len(a) != 2 {
		t.Errorf("Test: expecting 2 elements in SectionA, got %d", len(a))
	} else {
		if a["A"] != "B" {
			t.Errorf("Test: expecting \"SectionA/A\" == \"B\", got %q", a["A"])
		}
		if a["C"] != "D" {
			t.Errorf("Test: expecting \"SectionA/C\" == \"D\", got %q", a["C"])
		}
	}
	if b, ok := m["SectionB"]; !ok {
		t.Errorf("Test: missing SectionB")
	} else if len(b) != 1 {
		t.Errorf("Test: expecting 1 element in SectionB, got %d", len(b))
	} else {
		if b["E"] != "F" {
			t.Errorf("Test: expecting \"SectionB/E\" == \"F\", got %q", b["E"])
		}
	}
}

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
