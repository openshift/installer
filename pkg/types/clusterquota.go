package types

import (
	"github.com/openshift/installer/pkg/types/gcp"
)

// ClusterQuota contains the size, in cloud quota, of
// the cluster that was created by installer.
type ClusterQuota struct {
	GCP *gcp.Quota `json:"gcp,omitempty"`
}
