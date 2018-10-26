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

const (
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

func fetchLatestMetadata(ctx context.Context, channel string) (metadata, error) {
	build, err := fetchLatestBuild(ctx, channel)
	if err != nil {
		return metadata{}, errors.Wrap(err, "failed to fetch latest build")
	}

	url := fmt.Sprintf("%s/%s/%s/meta.json", baseURL, channel, build)
	logrus.Debugf("Fetching RHCOS metadata from %q", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return metadata{}, errors.Wrap(err, "failed to build request")
	}

	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return metadata{}, errors.Wrap(err, "failed to fetch metadata")
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

func fetchLatestBuild(ctx context.Context, channel string) (string, error) {
	url := fmt.Sprintf("%s/%s/builds.json", baseURL, channel)
	logrus.Debugf("Fetching RHCOS builds from %q", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to build request")
	}

	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch builds")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("incorrect HTTP response (%s)", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read HTTP response")
	}

	var builds struct {
		Builds []string `json:"builds"`
	}
	if err := json.Unmarshal(body, &builds); err != nil {
		return "", errors.Wrap(err, "failed to parse HTTP response")
	}

	if len(builds.Builds) == 0 {
		return "", errors.Errorf("no builds found")
	}

	return builds.Builds[0], nil
}
