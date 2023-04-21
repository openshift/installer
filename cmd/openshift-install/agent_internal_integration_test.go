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
			"isocmp":              isoCmp,
			"ignitionImgContains": ignitionImgContains,
			"initrdImgContains":   initrdImgContains,
		},
	})
}

// [!] ignitionImgContains `isoPath` `file` check if the specified file `file`
// is stored within /images/ignition.img archive in the ISO `isoPath` image.
func ignitionImgContains(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 2 {
		ts.Fatalf("usage: ignitionImgContains isoPath file")
	}

	workDir := ts.Getenv("WORK")
	isoPath, eFilePath := args[0], args[1]
	isoPathAbs := filepath.Join(workDir, isoPath)

	_, err := extractFileFromIgnitionImg(isoPathAbs, eFilePath)
	ts.Check(err)
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

	aData, err := readFileFromISOIgnitionCfg(isoPathAbs, aFilePath)
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

func readFileFromISOIgnitionCfg(isoPath string, nodePath string) ([]byte, error) {
	config, err := extractIgnitionCfg(isoPath)
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

func extractFileFromIgnitionImg(isoPath string, fileName string) ([]byte, error) {
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

	for {
		header, err := cpioReader.Next()
		if err == io.EOF { //nolint:errorlint
			// end of cpio archive
			break
		}
		if err != nil {
			return nil, err
		}

		if header.Name == fileName {
			rawContent, err := io.ReadAll(cpioReader)
			if err != nil {
				return nil, err
			}
			return rawContent, nil
		}
	}

	return nil, errors.NotFound(fmt.Sprintf("File %s not found within the /images/ignition.img archive", fileName))
}

func extractIgnitionCfg(isoPath string) (*igntypes.Config, error) {
	rawContent, err := extractFileFromIgnitionImg(isoPath, "config.ign")
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

// [!] initrdImgContains `isoPath` `file` check if the specified file `file`
// is stored within a compressed cpio archive by scanning the content of
// /images/ignition.img archive in the ISO `isoPath` image (note: plain cpio
// archives are ignored).
func initrdImgContains(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 2 {
		ts.Fatalf("usage: initrdImgContains isoPath file")
	}

	workDir := ts.Getenv("WORK")
	isoPath, eFilePath := args[0], args[1]
	isoPathAbs := filepath.Join(workDir, isoPath)

	err := checkFileFromInitrdImg(isoPathAbs, eFilePath)
	ts.Check(err)
}

func checkFileFromInitrdImg(isoPath string, fileName string) error {
	disk, err := diskfs.OpenWithMode(isoPath, diskfs.ReadOnly)
	if err != nil {
		return err
	}

	fs, err := disk.GetFilesystem(0)
	if err != nil {
		return err
	}

	initRdImg, err := fs.OpenFile("/images/pxeboot/initrd.img", os.O_RDONLY)
	if err != nil {
		return err
	}
	defer initRdImg.Close()

	const (
		gzipID1     = 0x1f
		gzipID2     = 0x8b
		gzipDeflate = 0x08
	)

	buff := make([]byte, 4096)
	for {
		_, err := initRdImg.Read(buff)
		if err == io.EOF { //nolint:errorlint
			break
		}

		foundAt := -1
		for idx := 0; idx < len(buff)-2; idx++ {
			// scan the buffer for a potential gzip header
			if buff[idx+0] == gzipID1 && buff[idx+1] == gzipID2 && buff[idx+2] == gzipDeflate {
				foundAt = idx
				break
			}
		}

		if foundAt >= 0 {
			// check if it's really a compressed cpio archive
			delta := int64(foundAt - len(buff))
			newPos, err := initRdImg.Seek(delta, io.SeekCurrent)
			if err != nil {
				break
			}

			files, err := lookForCpioFiles(initRdImg)
			if err != nil {
				if _, err := initRdImg.Seek(newPos+2, io.SeekStart); err != nil {
					break
				}
				continue
			}

			// check if the current cpio files match the required ones
			for _, f := range files {
				matched, err := filepath.Match(fileName, f)
				if err != nil {
					return err
				}
				if matched {
					return nil
				}
			}
		}
	}

	return errors.NotFound(fmt.Sprintf("File %s not found within the /images/pxeboot/initrd.img archive", fileName))
}

func lookForCpioFiles(r io.Reader) ([]string, error) {
	var files []string

	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	// skip in case of garbage
	if gr.OS != 255 && gr.OS >= 13 {
		return nil, fmt.Errorf("Unknown OS code: %v", gr.Header.OS)
	}

	cr := cpio.NewReader(gr)
	for {
		h, err := cr.Next()
		if err != nil {
			break
		}

		files = append(files, h.Name)
	}

	return files, nil
}
