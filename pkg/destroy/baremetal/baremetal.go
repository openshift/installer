// +build baremetal

package baremetal

import (
	"github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	InfraID                 string
	LibvirtURI              string
	BootstrapProvisioningIP string
	Logger                  logrus.FieldLogger
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() error {
	o.Logger.Debug("Deleting bare metal resources")

	// FIXME: close the connection
	conn, err := libvirt.NewConnect(o.LibvirtURI)
	if err != nil {
		return errors.Wrap(err, "failed to connect to Libvirt daemon")
	}
	err = o.deleteStoragePool(conn)
	if err != nil {
		return errors.Wrap(err, "failed to clean baremetal bootstrap storage pool")
	}

	o.Logger.Debug("FIXME: delete resources!")

	return nil
}

// New returns bare metal Uninstaller from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		InfraID:                 metadata.InfraID,
		LibvirtURI:              metadata.ClusterPlatformMetadata.BareMetal.LibvirtURI,
		BootstrapProvisioningIP: metadata.ClusterPlatformMetadata.BareMetal.BootstrapProvisioningIP,
		Logger:                  logger,
	}, nil
}

// deleteStoragePool destroys, deletes and undefines any storagePool left behind during the creation
// of the bootstrap VM
func (o *ClusterUninstaller) deleteStoragePool(conn *libvirt.Connect) error {
	o.Logger.Debug("Deleting baremetal bootstrap volumes")

	pName := o.InfraID + "-bootstrap"
	pool, err := conn.LookupStoragePoolByName(pName)
	if err != nil {
		o.Logger.Warnf("Unable to get storage pool %s: %s", pName, err)
		return nil
	}
	defer pool.Free()

	// delete vols
	vols, err := pool.ListAllStorageVolumes(0)
	if err != nil {
		o.Logger.Warnf("Unable to get volumes in storage pool %s: %s", pName, err)
		return nil
	}

	for _, vol := range vols {
		defer vol.Free()
		vName, err := vol.GetName()
		if err != nil {
			o.Logger.Warnf("Unable to get volume %s in storage pool %s: %s", vName, pName, err)
			return nil
		}
		if err := vol.Delete(0); err != nil {
			o.Logger.Warnf("Unable to delete volume %s in storage pool %s: %s", vName, pName, err)
			return nil
		}
		o.Logger.WithField("volume", vName).Info("Deleted volume")
	}

	if err := pool.Destroy(); err != nil {
		o.Logger.Warnf("Unable to destroy storage pool %s: %s", pName, err)
	}

	if err := pool.Delete(0); err != nil {
		o.Logger.Warnf("Unable to delete storage pool %s: %s", pName, err)
	}

	if err := pool.Undefine(); err != nil {
		o.Logger.Warnf("Unable to undefine storage pool %s: %s", pName, err)
	}
	o.Logger.WithField("pool", pName).Info("Deleted pool")

	return nil
}
