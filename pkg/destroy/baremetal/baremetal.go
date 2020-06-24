// +build baremetal

package baremetal

import (
	"context"

	"github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	LibvirtURI              string
	BootstrapProvisioningIP string
	Logger                  logrus.FieldLogger
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run(context.Context) error {
	o.Logger.Debug("Deleting bare metal resources")

	// FIXME: close the connection
	_, err := libvirt.NewConnect(o.LibvirtURI)
	if err != nil {
		return errors.Wrap(err, "failed to connect to Libvirt daemon")
	}

	o.Logger.Debug("FIXME: delete resources!")

	return nil
}

// New returns bare metal Uninstaller from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		LibvirtURI:              metadata.ClusterPlatformMetadata.BareMetal.LibvirtURI,
		BootstrapProvisioningIP: metadata.ClusterPlatformMetadata.BareMetal.BootstrapProvisioningIP,
		Logger:                  logger,
	}, nil
}
