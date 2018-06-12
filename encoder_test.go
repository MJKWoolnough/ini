package ini_test

import (
	"bytes"
	"testing"

	"vimagination.zapto.org/ini"
)

type encodeTest struct {
	Input  interface{}
	Output []byte
}

func testEncode(t *testing.T, tests []encodeTest) {
	var b bytes.Buffer
	for n, test := range tests {
		err := ini.Encode(&b, test.Input)
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else if !bytes.Equal(b.Bytes(), test.Output) {
			t.Errorf("test %d: expecting:\n%s\ngot:\n%s", n+1, test.Output, b.Bytes())
		}
		b.Reset()
	}
}
