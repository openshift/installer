package rhcos

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/openshift/installer/data"
	"github.com/pkg/errors"
)

type metadata struct {
	AMIs map[string]struct {
		HVM string `json:"hvm"`
	} `json:"amis"`
	Azure struct {
		Image string `json:"image"`
		URL string `json:"url"`
	}
	GCP struct {
		Image string `json:"image"`
		URL string `json:"url"`
	}
	BaseURI string `json:"baseURI"`
	Images  struct {
		QEMU struct {
			Path   string `json:"path"`
			SHA256 string `json:"sha256"`
		} `json:"qemu"`
	} `json:"images"`
	OSTreeVersion string `json:"ostree-version"`
}

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
