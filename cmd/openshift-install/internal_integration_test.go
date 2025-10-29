package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	"github.com/openshift/installer/internal/tshelpers"
	"github.com/openshift/installer/pkg/asset/releaseimage"
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
		// Uncomment below line to help debug the testcases
		// TestWork: true,

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

			var pullspec string
			// If set, let's use $OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE to replace test data
			// and also pass it to the testscript environment. It will be usually set in a CI job
			// to reference the ephemeral payload release.
			if releaseImageOverride, ok := os.LookupEnv("OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE"); ok && releaseImageOverride != "" {
				pullspec = releaseImageOverride
				e.Vars = append(e.Vars, fmt.Sprintf("OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE=%s", pullspec))
			} else {
				// Let's get the current release version, so that
				// it could be used within the tests
				pullspec, err = releaseimage.Default()
				if err != nil {
					return err
				}
			}
			e.Vars = append(e.Vars, fmt.Sprintf("RELEASE_IMAGE=%s", pullspec))
			if xdgCacheHome, ok := os.LookupEnv("XDG_CACHE_HOME"); ok && xdgCacheHome != "" {
				e.Vars = append(e.Vars, fmt.Sprintf("XDG_CACHE_HOME=%s", xdgCacheHome))
			}
			// When AUTH_FILE is set in the CI integration-tests job
			if authFilePath, ok := os.LookupEnv("AUTH_FILE"); ok && authFilePath != "" {
				workDir := e.Getenv("WORK")
				err := updatePullSecret(workDir, authFilePath)
				if err != nil {
					return err
				}
				t.Log("PullSecret updated successfully")
			}

			return nil
		},

		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"isocmp":                  tshelpers.IsoCmp,
			"ignitionImgContains":     tshelpers.IgnitionImgContains,
			"configImgContains":       tshelpers.ConfigImgContains,
			"initrdImgContains":       tshelpers.InitrdImgContains,
			"unconfiguredIgnContains": tshelpers.UnconfiguredIgnContains,
			"unconfiguredIgnCmp":      tshelpers.UnconfiguredIgnCmp,
			"expandFile":              tshelpers.ExpandFile,
			"isoContains":             tshelpers.IsoContains,
			"isoIgnitionContains":     tshelpers.IsoIgnitionContains,
			"isoSizeMin":              tshelpers.IsoSizeMin,
			"isoSizeMax":              tshelpers.IsoSizeMax,
		},
	})
}

func updatePullSecret(workDir, authFilePath string) error {
	authFile, err := os.ReadFile(authFilePath)
	if err != nil {
		return err
	}
	ciPullSecret := string(authFile)
	expectedInstallConfigPathAbs := filepath.Join(workDir, "install-config.yaml")
	_, err = os.Stat(expectedInstallConfigPathAbs)
	if err == nil {
		installConfigFile, err := os.ReadFile(expectedInstallConfigPathAbs)
		if err != nil {
			return err
		}
		var config map[string]interface{}
		if err := yaml.Unmarshal(installConfigFile, &config); err != nil {
			return err
		}

		config["pullSecret"] = ciPullSecret

		updatedConfig, err := yaml.Marshal(&config)
		if err != nil {
			return err
		}
		if err := os.WriteFile(expectedInstallConfigPathAbs, updatedConfig, 0600); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return nil
}
