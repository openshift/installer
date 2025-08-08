package service

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

// regex matching the path of a service entries file. The captured group is the name of the service.
// For example, if the filename is "log-bundle-20210329190553/bootstrap/services/release-image.json",
// then the name of the service is "release-image".
// In case the log-bundle is from bootstrap-in-place installation the file name is:
// "log-bundle-20210329190553/log-bundle-bootstrap/bootstrap/services/release-image.json"
var serviceEntriesFilePathRegex = regexp.MustCompile(`\/bootstrap\/services\/([^.]+)\.json$`)

// AnalyzeGatherBundle will analyze the bootstrap gather bundle at the specified path.
// Analysis will be logged.
// Returns an error if there was a problem reading the bundle.
func AnalyzeGatherBundle(bundlePath string) error {
	// open the bundle file for reading
	bundleFile, err := os.Open(bundlePath)
	if err != nil {
		return fmt.Errorf("could not open the gather bundle: %w", err)
	}
	defer bundleFile.Close()
	return analyzeGatherBundle(bundleFile)
}

func analyzeGatherBundle(bundleFile io.Reader) error {
	// decompress the bundle
	uncompressedStream, err := gzip.NewReader(bundleFile)
	if err != nil {
		return fmt.Errorf("could not decompress the gather bundle: %w", err)
	}
	defer uncompressedStream.Close()

	// read through the tar for relevant files
	tarReader := tar.NewReader(uncompressedStream)
	serviceAnalyses := make(map[string]analysis)
	servicesFound := make([]string, 0)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("encountered an error reading from the gather bundle: %w", err)
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}

		serviceEntriesFileSubmatch := serviceEntriesFilePathRegex.FindStringSubmatch(header.Name)
		if serviceEntriesFileSubmatch == nil {
			continue
		}
		serviceName := serviceEntriesFileSubmatch[1]
		servicesFound = append(servicesFound, serviceName)

		serviceAnalysis, err := analyzeService(tarReader)
		if err != nil {
			logrus.Infof("Could not analyze the %s.service: %v", serviceName, err)
			continue
		}

		serviceAnalyses[serviceName] = serviceAnalysis
	}

	if len(servicesFound) == 0 {
		logrus.Error("Invalid log bundle or the bootstrap machine could not be reached and bootstrap logs were not collected")
		return nil
	}

	analysisChecks := []struct {
		name     string
		check    func(analysis) bool
		optional bool
	}{
		{name: "node-image-pull", check: checkReleaseImageDownload, optional: false},
		{name: "release-image", check: checkReleaseImageDownload, optional: false},
		{name: "bootkube", check: checkBootkubeService, optional: false},
	}
	for _, check := range analysisChecks {
		a := serviceAnalyses[check.name]
		if a.starts == 0 {
			if check.optional {
				logrus.Infof("The bootstrap machine did not execute the %s.service systemd unit", check.name)
				break
			}
			logrus.Errorf("The bootstrap machine did not execute the %s.service systemd unit", check.name)
			break
		}
		if !check.check(a) {
			break
		}
	}

	return nil
}

func checkReleaseImageDownload(a analysis) bool {
	if a.successful {
		return true
	}
	logrus.Error("The bootstrap machine failed to download the release image")
	a.logLastError()
	return false
}

// bootstrap-verify-api-servel-urls.sh is currently running as part of the bootkube service.
// And the verification of the API and API-Int URLs are the only stage where a failure is
// currently reported. So, here we are able to conclude that a failure corresponds to a
// failure to resolve either the API URL or API-Int URL or both. If that changes and if
// any other stage in the bootkube service starts reporting a failure, we need to revisit
// this. At that point verification of the URLs could be moved to its own service.
func checkBootkubeService(a analysis) bool {
	if a.successful {
		return true
	}
	// Note: Even when there is a stage failure, we are not returning false here. That is
	// intentional because we donot want to report this as an error in the "analyze" output.
	a.logLastError()
	return true
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

func analyzeService(r io.Reader) (analysis, error) {
	a := analysis{}
	decoder := json.NewDecoder(r)
	t, err := decoder.Token()
	if err != nil {
		return a, fmt.Errorf("service entries file does not begin with a token: %w", err)
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
			return a, fmt.Errorf("could not decode an entry in the service entries file: %w", err)
		}

		// record a new start of the service
		if entry.Phase == ServiceStart {
			a.starts++
		}

		// the service is only considered successful if the last entry is either the service ending successfully or a
		// post-command ending successfully.
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

func (a analysis) logLastError() {
	for _, l := range strings.Split(a.lastError, "\n") {
		logrus.Info(l)
	}
}
