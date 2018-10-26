package types

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/openstack"
)

// ClusterMetadata contains information
// regarding the cluster that was created by installer.
type ClusterMetadata struct {
	ClusterName             string `json:"clusterName"`
	ClusterPlatformMetadata `json:",inline"`
}

// ClusterPlatformMetadata contains metadata for platfrom.
type ClusterPlatformMetadata struct {
	AWS       *aws.Metadata       `json:"aws,omitempty"`
	OpenStack *openstack.Metadata `json:"openstack,omitempty"`
	Libvirt   *libvirt.Metadata   `json:"libvirt,omitempty"`
}

// Platform returns a string representation of the platform
// (e.g. "aws" if AWS is non-nil).  It returns an empty string if no
// platform is configured.
func (cpm *ClusterPlatformMetadata) Platform() string {
	if cpm == nil {
		return ""
	}
	if cpm.AWS != nil {
		return "aws"
	}
	if cpm.Libvirt != nil {
		return "libvirt"
	}
	if cpm.OpenStack != nil {
		return "openstack"
	}
	return ""
}

// LoadClusterMetadata loads the cluster metadata from an asset directory.
func LoadClusterMetadata(dir string) (cmetadata *ClusterMetadata, err error) {
	raw, err := ioutil.ReadFile(filepath.Join(dir, "metadata.json"))
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(raw, &cmetadata); err != nil {
		return nil, err
	}

	return cmetadata, err
}
