package rhcos

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	// DefaultChannel is the default RHCOS channel for the cluster.
	DefaultChannel = "maipo"

	baseURL = "https://releases-rhcos.svc.ci.openshift.org/storage/releases"
)

type metadata struct {
	AMIs []struct {
		HVM  string `json:"hvm"`
		Name string `json:"name"`
	} `json:"amis"`
	Images struct {
		QEMU struct {
			Path   string `json:"path"`
			SHA256 string `json:"sha256"`
		} `json:"qemu"`
	} `json:"images"`
	OSTreeVersion string `json:"ostree-version"`
}

func fetchMetadata(ctx context.Context, channel string, build string) (metadata, error) {
	url := fmt.Sprintf("%s/%s/%s/meta.json", baseURL, channel, build)
	logrus.Debugf("Fetching RHCOS metadata from %q", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return metadata{}, errors.Wrap(err, "failed to build request")
	}

	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return metadata{}, errors.Wrapf(err, "failed to fetch metadata for build %s", build)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return metadata{}, errors.Errorf("incorrect HTTP response (%s)", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return metadata{}, errors.Wrap(err, "failed to read HTTP response")
	}

	var meta metadata
	if err := json.Unmarshal(body, &meta); err != nil {
		return meta, errors.Wrap(err, "failed to parse HTTP response")
	}

	return meta, nil
}
