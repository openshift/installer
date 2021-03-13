package service

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// AnalyzeGatherBundle will analyze the bootstrap gather bundle at the specified path.
// Analysis will be logged.
// Returns an error if there was a problem reading the bundle.
func AnalyzeGatherBundle(bundlePath string) error {
	// open the bundle file for reading
	bundleFile, err := os.Open(bundlePath)
	if err != nil {
		return errors.Wrap(err, "could not open the gather bundle")
	}
	defer bundleFile.Close()
	return analyzeGatherBundle(bundleFile)
}

func analyzeGatherBundle(bundleFile io.Reader) error {
	// decompress the bundle
	uncompressedStream, err := gzip.NewReader(bundleFile)
	if err != nil {
		return errors.Wrap(err, "could not decompress the gather bundle")
	}
	defer uncompressedStream.Close()

	// read through the tar for relevant files
	tarReader := tar.NewReader(uncompressedStream)
	var releaseImageAnalysis *analysis
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "encountered an error reading from the gather bundle")
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		filenameParts := strings.SplitN(header.Name, "/", 2)
		if len(filenameParts) != 2 {
			continue
		}
		// we only care about the release-image.service for now. in the future, we will look at other services, too.
		if filenameParts[1] == "bootstrap/services/release-image.json" {
			var err error
			releaseImageAnalysis, err = analyzeService(tarReader)
			if err != nil {
				logrus.Infof("Could not analyze the release-image.service: %v", err)
			}
			break
		}
	}

	// log details about the release-image.service.
	if releaseImageAnalysis != nil && releaseImageAnalysis.starts > 0 {
		if !releaseImageAnalysis.successful {
			logrus.Error("The bootstrap machine failed to download the release image")
			for _, l := range strings.Split(releaseImageAnalysis.lastError, "\n") {
				logrus.Info(l)
			}
		}
	} else {
		logrus.Error("The bootstrap machine did not execute the release-image.service systemd unit")
	}

	return nil
}

type analysis struct {
	// starts is the number of times that the service started
	starts int
	// successful is true if the last invocation of the service ended in success
	successful bool
	// failingStage is the stage that failed in the last unsuccessful invocation of the service
	failingStage string
	// lastError is the last error recorded in the last failure of the service
	lastError string
}

func analyzeService(r io.Reader) (*analysis, error) {
	a := &analysis{}
	decoder := json.NewDecoder(r)
	t, err := decoder.Token()
	if err != nil {
		return a, errors.Wrap(err, "service entries file does not begin with a token")
	}
	delim, isDelim := t.(json.Delim)
	if !isDelim {
		return a, errors.New("service entries file does not begin with a delimiter")
	}
	if delim != '[' {
		return a, errors.New("service entries file does not begin with an array")
	}
	var lastEntry *Entry
	for decoder.More() {
		entry := &Entry{}
		if err := decoder.Decode(entry); err != nil {
			return nil, errors.Wrap(err, "could not decode an entry in the service entries file")
		}

		// record a new start of the service
		if entry.Phase == ServiceStart {
			a.starts++
		}

		// the service is only considered considered successful if the very last entry is either the service ending
		// successfully or a post-command ending successfully.
		a.successful = entry.Result == Success && (entry.Phase == ServiceEnd || entry.Phase == PostCommandEnd)

		// save the last error
		if entry.Result == Failure {
			// if a stage failure causes a service (or pre- or post-command) failure, we want to preserve the failing
			// stage from the stage end entry.
			if lastEntry == nil || lastEntry.Phase != StageEnd || lastEntry.Result != Failure {
				a.failingStage = entry.Stage
			}
			a.lastError = entry.ErrorMessage
		}
		lastEntry = entry
	}
	return a, nil
}
