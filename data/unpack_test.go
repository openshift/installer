package data

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestUnpack(t *testing.T) {
	path, err := ioutil.TempDir("", "installer-data-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	err = Unpack(path, ".")
	if err != nil {
		t.Fatal(err)
	}

	expected := "# Bootstrap Module"
	content, err := ioutil.ReadFile(filepath.Join(path, "aws", "bootstrap", "README.md"))
	if err != nil {
		t.Fatal(err)
	}

	firstLine := string(bytes.SplitN(content, []byte("\n"), 2)[0])
	if firstLine != expected {
		t.Fatalf("%q != %q", firstLine, expected)
	}
}
