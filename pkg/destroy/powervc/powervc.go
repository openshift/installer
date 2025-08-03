package powervc

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	od "github.com/openshift/installer/pkg/destroy/openstack"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	Metadata *types.ClusterMetadata
	Logger   logrus.FieldLogger
}

// New returns an PowerVC destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		Metadata: metadata,
		Logger:   logger,
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	openstackMetadata := &openstack.Metadata{
		Cloud:      o.Metadata.ClusterPlatformMetadata.PowerVC.Cloud,
		Identifier: o.Metadata.ClusterPlatformMetadata.PowerVC.Identifier,
	}
	o.Metadata.ClusterPlatformMetadata.OpenStack = openstackMetadata

	openstackDestroyer, err := od.New(o.Logger, o.Metadata)
	if err != nil {
		return nil, errors.New("destroy PowerVC cannot call New OpenStack")
	}

	return openstackDestroyer.Run()
}
