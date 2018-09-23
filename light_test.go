package light

import (
	"testing"
)

func TestReadFile(t *testing.T) {
	bt, err := ReadFile("testdata/ReadFile.txt")
	if err != nil || string(bt) != "abcd" {
		t.Error("expect nil and abcd, got", string(bt), err)
	}
}
