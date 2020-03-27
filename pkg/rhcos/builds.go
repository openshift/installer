package rhcos

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/openshift/installer/data"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
)

var (
	errInvalidArch = fmt.Errorf("no build metadata for given architecture")
)

type metadata struct {
	AMIs map[string]struct {
		HVM string `json:"hvm"`
	} `json:"amis"`
	Azure struct {
		Image string `json:"image"`
		URL   string `json:"url"`
	}
	GCP struct {
		Image string `json:"image"`
		URL   string `json:"url"`
	}
	BaseURI string `json:"baseURI"`
	Images  struct {
		QEMU struct {
			Path               string `json:"path"`
			SHA256             string `json:"sha256"`
			UncompressedSHA256 string `json:"uncompressed-sha256"`
		} `json:"qemu"`
		OpenStack struct {
			Path               string `json:"path"`
			SHA256             string `json:"sha256"`
			UncompressedSHA256 string `json:"uncompressed-sha256"`
		} `json:"openstack"`
		VMware struct {
			Path   string `json:"path"`
			SHA256 string `json:"sha256"`
		} `json:"vmware"`
	} `json:"images"`
	OSTreeVersion string `json:"ostree-version"`
}

func fetchRHCOSBuild(ctx context.Context, arch types.Architecture) (*metadata, error) {
	file, err := data.Assets.Open(fmt.Sprintf("rhcos-%s.json", arch))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body, err := ioutil.ReadAll(file)
	if os.IsNotExist(err) {
		return nil, errInvalidArch
	} else if err != nil {
		return nil, err
	}

	var meta *metadata
	if err := json.Unmarshal(body, &meta); err != nil {
		return meta, errors.Wrap(err, "failed to parse RHCOS build metadata")
	}

	return meta, nil
}
