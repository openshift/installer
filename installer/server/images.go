package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/go-semver/semver"
	"golang.org/x/net/context"

	"github.com/coreos/tectonic-installer/installer/server/ctxh"
)

const (
	imagesTimeout = 5 * time.Second
)

var (
	coreosMinVersion = semver.New("1010.1.0")
)

// Image represents a versioned set of OS image assets.
type Image struct {
	Version string `json:"version"`
}

// Images represents the list images response format.
type Images struct {
	CoreOS []Image `json:"coreos"`
}

// listImagesHandler returns the list of available CoreOS images at the
// matchbox assets endpoint.
func listImagesHandler() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		endpoint, err := endpointFromRequest(req)
		if err != nil {
			return ctxh.NewAppError(err, "endpoint argument required", http.StatusBadRequest)
		}
		log.Debugf("querying endpoint for images: %v", endpoint)

		coreosImages, err := listCoreOSImages(ctx, endpoint)
		if err != nil {
			return ctxh.NewAppError(err, fmt.Sprintf("Failed to query available images: %v", err.Error()), http.StatusBadRequest)
		}

		b, err := json.Marshal(&Images{
			CoreOS: filterVersions(coreosImages, coreosMinVersion),
		})
		if err != nil {
			return ctxh.NewAppError(err, fmt.Sprintf("Failed to marshal images listing: %v", err.Error()), http.StatusInternalServerError)
		}

		writeJSON(w, b)
		return nil
	}
	return ctxh.ContextHandlerFuncWithError(fn)
}

func endpointFromRequest(req *http.Request) (string, error) {
	endpoint := req.URL.Query().Get("endpoint")
	if endpoint == "" {
		return "", errors.New("installer: No endpoint provided")
	}
	return endpoint, nil
}

// listCoreOSImages fetches and parses the CoreOS images available at an
// endpoint (e.g. http://matchbox.foo:8080/assets/coreos/).
//
// The caller must provide a suitable image endpoint, CoreOS image indexes
// don't have a consistent identifier, just a list of image anchors.
func listCoreOSImages(ctx context.Context, endpoint string) ([]Image, error) {
	client := &http.Client{
		Timeout:   time.Duration(imagesTimeout),
		Transport: http.DefaultTransport,
	}

	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var images []Image

	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch {
		case tt == html.StartTagToken:
			token := z.Token()
			if token.Data == "a" {
				images = append(images, getCoreOSImage(token))
			}
		case tt == html.ErrorToken:
			return images, nil
		}
	}
}

func getCoreOSImage(token html.Token) Image {
	for _, attr := range token.Attr {
		if attr.Key == "href" {
			return Image{
				Version: strings.Replace(attr.Val, "/", "", -1),
			}
		}
	}
	return Image{}
}

// filterVersions returns images with versions greater than the minVersion.
func filterVersions(images []Image, minVersion *semver.Version) (filtered []Image) {
	for _, image := range images {
		v, err := semver.NewVersion(image.Version)
		// skip images with non-semver versions
		if err == nil && v != nil && coreosMinVersion.LessThan(*v) {
			filtered = append(filtered, image)
		}
	}

	return filtered
}
