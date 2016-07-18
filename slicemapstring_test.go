package ini_test

import "testing"

func TestSliceMapString(t *testing.T) {
	testStruct(t, []structTest{
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
	})
}
