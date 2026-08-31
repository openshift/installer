package data

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUnpack(t *testing.T) {
	path := t.TempDir()

	err := Unpack(path, ".")
	if err != nil {
		t.Fatal(err)
	}

	expected := "# Bootstrap Module"
	content, err := os.ReadFile(filepath.Join(path, "libvirt", "bootstrap", "README.md"))
	if err != nil {
		t.Fatal(err)
	}

	firstLine := string(bytes.SplitN(content, []byte("\n"), 2)[0])
	if firstLine != expected {
		t.Fatalf("%q != %q", firstLine, expected)
	}
}

func TestInstallerGatherSSHConnectTimeout(t *testing.T) {
	content := installerGatherContent(t)

	commandCount := 0
	for _, line := range strings.Split(content, "\n") {
		command := strings.TrimSpace(line)
		if !strings.HasPrefix(command, "scp ") && !strings.HasPrefix(command, "ssh ") {
			continue
		}
		commandCount++
		if !strings.Contains(command, "-o ConnectTimeout=30") {
			t.Errorf("SSH command does not set ConnectTimeout=30: %s", command)
		}
	}
	if commandCount != 3 {
		t.Errorf("expected 3 SSH commands, found %d", commandCount)
	}
}

func installerGatherContent(t *testing.T) string {
	t.Helper()
	file, err := Assets.Open("bootstrap/files/usr/local/bin/installer-gather.sh")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	return string(content)
}
