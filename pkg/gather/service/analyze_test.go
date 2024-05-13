package service

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func generateSuccessOutput(stage string) string {
	return `[
{"phase":"service start"},
{"phase":"stage start", "stage":"` + stage + `"},
{"phase":"stage end", "stage":"` + stage + `", "result":"success"},
{"phase":"service end", "result":"success"}
]`
}

func generateFailureOutput(stage string) string {
	return `[
{"phase":"service start"},
{"phase":"stage start", "stage":"` + stage + `"},
{"phase":"stage end", "stage":"` + stage + `", "result":"failure", "errorMessage":"Line 1\nLine 2\nLine 3"}
]`
}

func failedReleaseImage() []logrus.Entry {
	return []logrus.Entry{
		{Level: logrus.ErrorLevel, Message: "The bootstrap machine failed to download the release image"},
		{Level: logrus.InfoLevel, Message: "Line 1"},
		{Level: logrus.InfoLevel, Message: "Line 2"},
		{Level: logrus.InfoLevel, Message: "Line 3"},
	}
}

func failedURLChecks() []logrus.Entry {
	return []logrus.Entry{
		{Level: logrus.InfoLevel, Message: "Line 1"},
		{Level: logrus.InfoLevel, Message: "Line 2"},
		{Level: logrus.InfoLevel, Message: "Line 3"},
	}
}

func TestAnalyzeGatherBundle(t *testing.T) {
	cases := []struct {
		name           string
		files          map[string]string
		expectedOutput []logrus.Entry
	}{
		{
			name: "no files",
			expectedOutput: []logrus.Entry{
				{Level: logrus.ErrorLevel, Message: "The bootstrap machine did not execute the release-image.service systemd unit"},
			},
		},
		{
			name: "release-image not started",
			files: map[string]string{
				"log-bundle/bootstrap/services/release-image.json": "[]",
			},
			expectedOutput: []logrus.Entry{
				{Level: logrus.ErrorLevel, Message: "The bootstrap machine did not execute the release-image.service systemd unit"},
			},
		},
		{
			name: "bootkube not started",
			files: map[string]string{
				"log-bundle/bootstrap/services/release-image.json": generateSuccessOutput("pull-release-image"),
				"log-bundle/bootstrap/services/bootkube.json":      "[]",
			},
			expectedOutput: []logrus.Entry{
				{Level: logrus.ErrorLevel, Message: "The bootstrap machine did not execute the bootkube.service systemd unit"},
			},
		},
		{
			name: "release-image and API Server URL successful",
			files: map[string]string{
				"log-bundle/bootstrap/services/release-image.json": generateSuccessOutput("pull-release-image"),
				"log-bundle/bootstrap/services/bootkube.json":      generateSuccessOutput("check-api-url"),
			},
		},
		{
			name: "release-image and API Server URL successful bootstrap-in-place",
			files: map[string]string{
				"log-bundle/log-bundle-bootstrap/bootstrap/services/release-image.json": generateSuccessOutput("pull-release-image"),
				"log-bundle/bootstrap/services/bootkube.json":                           generateSuccessOutput("check-api-url"),
			},
		},
		{
			name: "only release-image failed",
			files: map[string]string{
				"log-bundle/bootstrap/services/release-image.json": generateFailureOutput("pull-release-image"),
				"log-bundle/bootstrap/services/bootkube.json":      generateSuccessOutput("check-api-url"),
			},
			expectedOutput: failedReleaseImage(),
		},
		{
			name: "API Server URL failed",
			files: map[string]string{
				"log-bundle/log-bundle-bootstrap/bootstrap/services/release-image.json": generateSuccessOutput("pull-release-image"),
				"log-bundle/bootstrap/services/bootkube.json":                           generateFailureOutput("check-api-url"),
			},
			expectedOutput: failedURLChecks(),
		},
		{
			name: "API-INT Server URL failed",
			files: map[string]string{
				"log-bundle/log-bundle-bootstrap/bootstrap/services/release-image.json": generateSuccessOutput("pull-release-image"),
				"log-bundle/bootstrap/services/bootkube.json":                           generateFailureOutput("check-api-int-url"),
			},
			expectedOutput: failedURLChecks(),
		},
		{
			name: "both release-image and API Server URLs failed",
			files: map[string]string{
				"log-bundle/log-bundle-bootstrap/bootstrap/services/release-image.json": generateFailureOutput("pull-release-image"),
				"log-bundle/bootstrap/services/bootkube.json":                           generateFailureOutput("check-api-url"),
			},
			expectedOutput: failedReleaseImage(),
		},
		{
			name: "empty release-image.json",
			files: map[string]string{
				"log-bundle/bootstrap/services/release-image.json": "",
			},
			expectedOutput: []logrus.Entry{
				{Level: logrus.InfoLevel, Message: "Could not analyze the release-image.service: service entries file does not begin with a token: EOF"},
				{Level: logrus.ErrorLevel, Message: "The bootstrap machine did not execute the release-image.service systemd unit"},
			},
		},
		{
			name: "empty bootkube.json",
			files: map[string]string{
				"log-bundle/bootstrap/services/release-image.json": generateSuccessOutput("pull-release-image"),
				"log-bundle/bootstrap/services/bootkube.json":      "",
			},
			expectedOutput: []logrus.Entry{
				{Level: logrus.InfoLevel, Message: "Could not analyze the bootkube.service: service entries file does not begin with a token: EOF"},
				{Level: logrus.ErrorLevel, Message: "The bootstrap machine did not execute the bootkube.service systemd unit"},
			},
		},
		{
			name: "malformed release-image.json",
			files: map[string]string{
				"log-bundle/bootstrap/services/release-image.json": "{}",
			},
			expectedOutput: []logrus.Entry{
				{Level: logrus.InfoLevel, Message: "Could not analyze the release-image.service: service entries file does not begin with an array"},
				{Level: logrus.ErrorLevel, Message: "The bootstrap machine did not execute the release-image.service systemd unit"},
			},
		},
		{
			name: "malformed bootkube.json",
			files: map[string]string{
				"log-bundle/bootstrap/services/release-image.json": generateSuccessOutput("pull-release-image"),
				"log-bundle/bootstrap/services/bootkube.json":      "{}",
			},
			expectedOutput: []logrus.Entry{
				{Level: logrus.InfoLevel, Message: "Could not analyze the bootkube.service: service entries file does not begin with an array"},
				{Level: logrus.ErrorLevel, Message: "The bootstrap machine did not execute the bootkube.service systemd unit"},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gatherBuilder bytes.Buffer
			gzipWriter := gzip.NewWriter(&gatherBuilder)
			defer gzipWriter.Close()
			tarWriter := tar.NewWriter(gzipWriter)
			defer tarWriter.Close()
			for filename, contents := range tc.files {
				contentsAsBytes := []byte(contents)
				if err := tarWriter.WriteHeader(&tar.Header{
					Typeflag: tar.TypeReg,
					Name:     filename,
					Size:     int64(len(contentsAsBytes)),
				}); err != nil {
					t.Fatal(err)
				}
				if _, err := tarWriter.Write(contentsAsBytes); err != nil {
					t.Fatal(err)
				}
			}
			gzipWriter.Close()
			hook := test.NewLocal(logrus.StandardLogger())
			defer hook.Reset()
			err := analyzeGatherBundle(&gatherBuilder)
			assert.NoError(t, err, "unexpected error from analysis")
			for i, e := range hook.Entries {
				hook.Entries[i] = logrus.Entry{
					Level:   e.Level,
					Message: e.Message,
				}
			}
			assert.Equal(t, tc.expectedOutput, hook.Entries)
		})
	}
}
