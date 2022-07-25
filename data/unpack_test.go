package data

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestUnpack(t *testing.T) {
	path := t.TempDir()

	err := Unpack(path, ".")
	if err != nil {
		t.Fatal(err)
	}

	expected := "# Bootstrap Module"
	content, err := ioutil.ReadFile(filepath.Join(path, "baremetal", "bootstrap", "README.md"))
	if err != nil {
		t.Fatal(err)
	}

	firstLine := string(bytes.SplitN(content, []byte("\n"), 2)[0])
	if firstLine != expected {
		t.Fatalf("%q != %q", firstLine, expected)
	}
}
