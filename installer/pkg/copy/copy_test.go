package copy

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "workflow-test-")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)

	sourcePath := filepath.Join(dir, "source")
	sourceContent := []byte("Hello, World!\n")
	err = ioutil.WriteFile(sourcePath, sourceContent, 0600)
	if err != nil {
		t.Error(err)
	}

	targetPath := filepath.Join(dir, "target")
	err = Copy(sourcePath, targetPath)
	if err != nil {
		t.Error(err)
	}

	targetContent, err := ioutil.ReadFile(targetPath)
	if err != nil {
		t.Error(err)
	}

	if string(targetContent) != string(sourceContent) {
		t.Errorf("target %q != source %q", string(targetContent), string(sourceContent))
	}
}
