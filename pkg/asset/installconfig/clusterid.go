package installconfig

import (
	"github.com/pborman/uuid"

	"github.com/openshift/installer/pkg/asset"
)

type clusterID struct {
	ClusterID string
}

var _ asset.Asset = (*clusterID)(nil)

// Dependencies returns no dependencies.
func (a *clusterID) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates a new UUID
func (a *clusterID) Generate(asset.Parents) error {
	a.ClusterID = uuid.New()
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *clusterID) Name() string {
	return "Cluster ID"
}
