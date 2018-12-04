package installconfig

import (
	"github.com/pborman/uuid"

	"github.com/openshift/installer/pkg/asset"
)

// ClusterID is the unique ID of the cluster, immutable during the cluster's life
type ClusterID struct {
	ClusterID string
}

var _ asset.Asset = (*ClusterID)(nil)

// Dependencies returns no dependencies.
func (a *ClusterID) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates a new UUID
func (a *ClusterID) Generate(asset.Parents) error {
	a.ClusterID = uuid.New()
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *ClusterID) Name() string {
	return "Cluster ID"
}
