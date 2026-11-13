package data

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestUnpack(t *testing.T) {
	const expected = "test asset\n"
	originalAssets := Assets
	t.Cleanup(func() {
		Assets = originalAssets
	})
	Assets = http.FS(fstest.MapFS{
		"nested/asset.txt": {Data: []byte(expected)},
	})

	path := t.TempDir()

	err := Unpack(path, ".")
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(path, "nested", "asset.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != expected {
		t.Fatalf("%q != %q", content, expected)
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
