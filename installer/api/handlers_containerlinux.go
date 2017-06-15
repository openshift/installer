package api

import (
	"net/http"
	"time"

	"github.com/coreos/go-semver/semver"

	"github.com/coreos/tectonic-installer/installer/pkg/containerlinux"
)

const (
	containerLinuxListTimeout = 10 * time.Second
)

var containerLinuxMinVersion = semver.New("1010.1.0")

// listMatchboxImagesHandler returns the list of available Container Linux
// images at the matchbox assets endpoint.
func listMatchboxImagesHandler(w http.ResponseWriter, req *http.Request, _ *Context) error {
	endpoint := req.URL.Query().Get("endpoint")
	if endpoint == "" {
		return newBadRequestError("No endpoint provided")
	}

	images, err := containerlinux.ListMatchboxImages(endpoint, *containerLinuxMinVersion, containerLinuxListTimeout)
	if err != nil {
		return newInternalServerError("Failed to query available images: %s", err)
	}

	type Image struct {
		Version string `json:"version"`
	}
	response := struct {
		CoreOS []Image `json:"coreos"`
	}{}
	for _, image := range images {
		response.CoreOS = append(response.CoreOS, Image{Version: image.String()})
	}

	return writeJSONResponse(w, req, http.StatusOK, response)
}

// listAMIImagesHandler returns the list of available Container Linux AMIs.
func listAMIImagesHandler(w http.ResponseWriter, req *http.Request, _ *Context) error {
	amis, err := containerlinux.ListAMIImages(containerLinuxListTimeout)
	if err != nil {
		return newInternalServerError("Failed to query available images: %s", err)
	}

	return writeJSONResponse(w, req, http.StatusOK, amis)
}
