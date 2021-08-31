package types

import (
	"github.com/openshift/installer/pkg/types/gcp"
)

// ClusterQuota contains the size, in cloud quota, of
// the cluster that was created by installer.
type ClusterQuota struct {
	GCP *gcp.Quota `json:"gcp,omitempty"`
}

// Platform returns a string representation of the platform
// (e.g. "aws" if AWS is non-nil).  It returns an empty string if no
// platform is configured.
func (cpm *ClusterQuota) Platform() string {
	if cpm == nil {
		return ""
	}
	if cpm.GCP != nil {
		return gcp.Name
	}
	return ""
}
