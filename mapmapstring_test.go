package ini_test

import (
	"testing"

	"vimagination.zapto.org/ini"
)

func TestDecodeMapMapString(t *testing.T) {
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

func TestDecodeMapMapStringCustomTypes(t *testing.T) {
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

func TestEncodeMapMapString(t *testing.T) {
	testEncode(t, []encodeTest{
		{
			map[string]map[string]string{
				"": map[string]string{
					"A": "1",
					"C": "2",
					"B": "3",
				},
			},
			[]byte("A=1\nB=3\nC=2"),
		},
		{
			map[string]map[string]string{
				"ABC": map[string]string{
					"A": "1",
					"C": "2",
					"B": "3",
				},
			},
			[]byte("[ABC]\nA=1\nB=3\nC=2"),
		},
		{
			map[string]map[string]string{
				"": map[string]string{
					"A": "1",
					"C": "2",
					"B": "3",
				},
				"DEF": map[string]string{
					"D": "4",
					"F": "5",
					"E": "6",
				},
			},
			[]byte("A=1\nB=3\nC=2\n\n[DEF]\nD=4\nE=6\nF=5"),
		},
	})
}
