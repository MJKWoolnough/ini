package ini_test

import "testing"

func TestSliceStruct(t *testing.T) {
	type tStruct struct {
		Key string
		ZB  string
	}
	testStruct(t, []structTest{
		{
			Input: "[Section1]\nKey=Value\nZB=42\n[Section2]\nKey=Master\nZB=YA",
			Output: struct {
				Sections []tStruct `ini:"Section,prefix"`
			}{
				Sections: []tStruct{
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
	})
}
