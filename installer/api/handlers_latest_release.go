package api

import (
	"io"
	"net/http"

	"github.com/coreos/tectonic-installer/installer/version"
)

// latestReleaseHandler gets the tectonic's latest release from coreos.com.
func latestReleaseHandler(w http.ResponseWriter, req *http.Request, _ *Context) error {
	res, err := http.Get("https://coreos.com/tectonic/api/releases/latest")
	if err != nil {
		return newBadRequestError("Failed to get a response from coreos.com: %s", err)
	}
	io.Copy(w, res.Body)
	res.Body.Close()
	return nil
}

// Fetch tectonic's Build version and return in JSON format
func currentReleaseHandler(w http.ResponseWriter, req *http.Request, _ *Context) error {
	version := struct {
		TectonicVersion string `json:"tectonicVersion"`
		BuildTime       string `json:"buildTime"`
	}{
		TectonicVersion: version.TectonicVersion,
		BuildTime:       version.BuildTime,
	}

	return writeJSONResponse(w, req, http.StatusOK, version)
}
