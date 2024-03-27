package baremetal

import (
	"fmt"
	"net/url"

	libvirt "github.com/digitalocean/go-libvirt"
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
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	o.Logger.Debug("Deleting bare metal resources")

	uri, err := url.Parse(o.LibvirtURI)
	if err != nil {
		return nil, err
	}

	virt, err := libvirt.ConnectToURI(uri)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := virt.Disconnect(); err != nil {
			if o.Logger != nil {
				o.Logger.Warn("failed to disconnect from libvirt", err)
			}
		}
	}()

	err = o.deleteStoragePool(virt)
	if err != nil {
		return nil, fmt.Errorf("failed to clean baremetal bootstrap storage pool: %w", err)
	}

	return nil, nil
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
func (o *ClusterUninstaller) deleteStoragePool(virt *libvirt.Libvirt) error {
	o.Logger.Debug("Deleting baremetal bootstrap volumes")

	pName := o.InfraID + "-bootstrap"
	pool, err := virt.StoragePoolLookupByName(pName)
	if err != nil {
		o.Logger.Warnf("Unable to get storage pool %s: %s", pName, err)
		return nil
	}

	// delete vols
	vols, err := virt.StoragePoolListVolumes(pool, 0)
	if err != nil {
		o.Logger.Warnf("Unable to get volumes in storage pool %s: %s", pName, err)
		return nil
	}

	for _, vName := range vols {
		vol, err := virt.StorageVolLookupByName(pool, vName)
		if err != nil {
			o.Logger.Warnf("Unable to get volume %s in storage pool %s: %s", vName, pName, err)
			return nil
		}
		if err := virt.StorageVolDelete(vol, 0); err != nil {
			o.Logger.Warnf("Unable to delete volume %s in storage pool %s: %s", vName, pName, err)
			return nil
		}
		o.Logger.WithField("volume", vName).Info("Deleted volume")
	}

	if err := virt.StoragePoolDestroy(pool); err != nil {
		o.Logger.Warnf("Unable to destroy storage pool %s: %s", pName, err)
	}

	if err := virt.StoragePoolDelete(pool, 0); err != nil {
		o.Logger.Warnf("Unable to delete storage pool %s: %s", pName, err)
	}

	if err := virt.StoragePoolUndefine(pool); err != nil {
		o.Logger.Warnf("Unable to undefine storage pool %s: %s", pName, err)
	}
	o.Logger.WithField("pool", pName).Info("Deleted pool")

	return nil
}
