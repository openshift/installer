package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cavaliercoder/go-cpio"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/diskfs/go-diskfs"
	"github.com/go-openapi/errors"
	"github.com/pkg/diff"
	"github.com/rogpeppe/go-internal/testscript"
	"github.com/stretchr/testify/assert"
	"github.com/vincent-petithory/dataurl"
)

// This file contains a number of functions useful for
// setting up the environment and running the integration
// tests for the agent-based installer

// runAllIntegrationTests runs all the tests found in the (sub)folders
// rooted at rootPath. Folders that do not contain a test file (.txt or .txtar)
// are ignored.
func runAllIntegrationTests(t *testing.T, rootPath string) {
	t.Helper()
	suites := []string{}

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			files, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			for _, f := range files {
				if !f.IsDir() && (strings.HasSuffix(f.Name(), ".txt") || strings.HasSuffix(f.Name(), ".txtar")) {
					for _, s := range suites {
						if s == path {
							return nil
						}
					}
					suites = append(suites, path)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range suites {
		t.Run(generateTestName(s), func(t *testing.T) {
			runIntegrationTest(t, s)
		})
	}
}

func generateTestName(path string) string {
	name := strings.TrimPrefix(path, "testdata/")
	return strings.ReplaceAll(name, "/", "_")
}

func runIntegrationTest(t *testing.T, testFolder string) {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	projectDir, err := os.Getwd()
	assert.NoError(t, err)
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	testscript.Run(t, testscript.Params{
		Dir: testFolder,

		Setup: func(e *testscript.Env) error {
			// This is required to allow proper
			// loading of the embedded resources
			e.Cd = filepath.Join(projectDir, "../../data")

			// For agent commands, let's use the
			// current home dir
			for i, v := range e.Vars {
				if v == "HOME=/no-home" {
					e.Vars[i] = fmt.Sprintf("HOME=%s", homeDir)
					break
				}
			}

			return nil
		},

		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"isocmp": isoCmp,
		},
	})
}

// [!] isoCmp `isoPath` `isoFile` `expectedFile` check that the content of the file
// `isoFile` - extracted from the ISO embedded ignition configuration file referenced
// by `isoPath` - matches the content of the local file `expectedFile`.
func isoCmp(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 3 {
		ts.Fatalf("usage: isocmp isoPath file1 file2")
	}

	workDir := ts.Getenv("WORK")
	isoPath, aFilePath, eFilePath := args[0], args[1], args[2]
	isoPathAbs := filepath.Join(workDir, isoPath)

	aData, err := readFileFromISO(isoPathAbs, aFilePath)
	ts.Check(err)

	eFilePathAbs := filepath.Join(workDir, eFilePath)
	eData, err := os.ReadFile(eFilePathAbs)
	ts.Check(err)

	aText := string(aData)
	eText := string(eData)

	eq := aText == eText
	if neg {
		if eq {
			ts.Fatalf("%s and %s do not differ", aFilePath, eFilePath)
		}
		return
	}
	if eq {
		return
	}

	ts.Logf(aText)

	var sb strings.Builder
	if err := diff.Text(aFilePath, eFilePath, aText, eText, &sb); err != nil {
		ts.Check(err)
	}

	ts.Logf("%s", sb.String())
	ts.Fatalf("%s and %s differ", aFilePath, eFilePath)
}

func readFileFromISO(isoPath string, nodePath string) ([]byte, error) {
	config, err := readIgnitionFromISO(isoPath)
	if err != nil {
		return nil, err
	}

	for _, f := range config.Storage.Files {
		if f.Node.Path == nodePath {
			actualData, err := dataurl.DecodeString(*f.FileEmbedded1.Contents.Source)
			if err != nil {
				return nil, err
			}
			return actualData.Data, nil
		}
	}

	return nil, errors.NotFound(nodePath)
}

func readIgnitionFromISO(isoPath string) (*igntypes.Config, error) {
	disk, err := diskfs.OpenWithMode(isoPath, diskfs.ReadOnly)
	if err != nil {
		return nil, err
	}

	fs, err := disk.GetFilesystem(0)
	if err != nil {
		return nil, err
	}

	ignitionImg, err := fs.OpenFile("/images/ignition.img", os.O_RDONLY)
	if err != nil {
		return nil, err
	}

	gzipReader, err := gzip.NewReader(ignitionImg)
	if err != nil {
		return nil, err
	}

	cpioReader := cpio.NewReader(gzipReader)
	_, err = cpioReader.Next()
	if err != nil {
		return nil, err
	}

	rawContent, err := io.ReadAll(cpioReader)
	if err != nil {
		return nil, err
	}

	var config igntypes.Config
	err = json.Unmarshal(rawContent, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
