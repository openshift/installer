package rhcos

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/openshift/installer/data"
	"github.com/pkg/errors"
)

// metadata is a subset of the `meta.json` generated
// by https://github.com/coreos/coreos-assembler
type metadata struct {
	AMIs map[string]struct {
		HVM string `json:"hvm"`
	} `json:"amis"`
	BaseURI string `json:"baseURI"`
	Images  struct {
		QEMU struct {
			Path   string `json:"path"`
			SHA256 string `json:"sha256"`
		} `json:"qemu"`
	} `json:"images"`
	OSTreeVersion string `json:"ostree-version"`
}

// fetchRHCOSBuild retrieves the pinned RHEL CoreOS metadata.
func fetchRHCOSBuild(ctx context.Context) (*metadata, error) {
	file, err := data.Assets.Open("rhcos.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var meta *metadata
	if err := json.Unmarshal(body, &meta); err != nil {
		return meta, errors.Wrap(err, "failed to parse RHCOS build metadata")
	}

	return meta, nil
}

// FetchVersion retrives the pinned RHCOS version.
func FetchVersion() (string, error) {
	meta, err := fetchRHCOSBuild(context.TODO())
	if err != nil {
		return "", err
	}
	return meta.OSTreeVersion, nil
}
