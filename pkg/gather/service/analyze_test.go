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
			name: "release-image successful",
			files: map[string]string{
				"log-bundle/bootstrap/services/release-image.json": `[
{"phase":"service start"},
{"phase":"stage start", "stage":"pull-release-image"},
{"phase":"stage end", "stage":"pull-release-image", "result":"success"},
{"phase":"service end", "result":"success"}
]`,
			},
		},
		{
			name: "release-image successful bootstrap-in-place",
			files: map[string]string{
				"log-bundle/log-bundle-bootstrap/bootstrap/services/release-image.json": `[
{"phase":"service start"},
{"phase":"stage start", "stage":"pull-release-image"},
{"phase":"stage end", "stage":"pull-release-image", "result":"success"},
{"phase":"service end", "result":"success"}
]`,
			},
		},
		{
			name: "release-image failed",
			files: map[string]string{
				"log-bundle/bootstrap/services/release-image.json": `[
{"phase":"service start"},
{"phase":"stage start", "stage":"pull-release-image"},
{"phase":"stage end", "stage":"pull-release-image", "result":"failure", "errorMessage":"Line 1\nLine 2\nLine 3"}
]`,
			},
			expectedOutput: []logrus.Entry{
				{Level: logrus.ErrorLevel, Message: "The bootstrap machine failed to download the release image"},
				{Level: logrus.InfoLevel, Message: "Line 1"},
				{Level: logrus.InfoLevel, Message: "Line 2"},
				{Level: logrus.InfoLevel, Message: "Line 3"},
			},
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
			name: "malformed release-image.json",
			files: map[string]string{
				"log-bundle/bootstrap/services/release-image.json": "{}",
			},
			expectedOutput: []logrus.Entry{
				{Level: logrus.InfoLevel, Message: "Could not analyze the release-image.service: service entries file does not begin with an array"},
				{Level: logrus.ErrorLevel, Message: "The bootstrap machine did not execute the release-image.service systemd unit"},
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
