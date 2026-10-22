package machines

import (
	"fmt"

	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/api/features"
	v1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines/machineconfig"
	"github.com/openshift/installer/pkg/asset/machines/vsphere"
	"github.com/openshift/installer/pkg/types"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

const (
	// VsphereScsiByPath defines the path format for vsphere disks being added.
	VsphereScsiByPath = "/dev/disk/by-path/pci-0000:03:00.0-scsi-0:0:%d:0"

	// AzureDataDiskByLun defines the path format for azure data disks being added.
	AzureDataDiskByLun = "/dev/disk/azure/data/by-lun/%d"
)

// DiskName returns the platform disk ID for the given disk setup entry.
func DiskName(diskSetup types.Disk) (string, error) {
	switch diskSetup.Type {
	case types.Etcd:
		return diskSetup.Etcd.PlatformDiskID, nil
	case types.Swap:
		return diskSetup.Swap.PlatformDiskID, nil
	case types.UserDefined:
		return diskSetup.UserDefined.PlatformDiskID, nil
	default:
		return "", errors.Errorf("disk setup type %s is not supported", diskSetup.Type)
	}
}

// NodeDiskSetup determines the path per disk type, and per platform and role, runs ForDiskSetup.
func NodeDiskSetup(installConfig *installconfig.InstallConfig, role string, diskSetup types.Disk, dataDisk any) (*v1.MachineConfig, error) {
	ic := installConfig.Config

	label := string(diskSetup.Type)

	// Etcd disks mount at /var/lib/etcd, user-defined disks at their
	// configured mount path, and swap disks mount nothing.
	path, _ := diskSetup.MountPath()
	if diskSetup.Type == types.UserDefined {
		label = diskSetup.UserDefined.PlatformDiskID
	}

	switch ic.Platform.Name() {
	case azuretypes.Name:
		if installConfig.Config.Enabled(features.FeatureGateAzureMultiDisk) {
			if azureDataDisk, ok := dataDisk.(v1beta1.DataDisk); ok {
				// The azure-vm-utils udev rules link data disks by LUN for both the SCSI and the
				// NVMe disk controller. The older scsi1/lun<N> link exists on SCSI only, so it is
				// never created on NVMe-only instance types and Ignition times out waiting for it.
				device := fmt.Sprintf(AzureDataDiskByLun, *azureDataDisk.Lun)
				diskSetupIgn, err := machineconfig.ForDiskSetup(role, device, label, path, diskSetup.Type)
				if err != nil {
					return nil, errors.Wrap(err, "failed to create ignition to setup disks for master machines")
				}
				return diskSetupIgn, nil
			}
			return nil, errors.Errorf("unsupported azure data disk type")
		}
	case vspheretypes.Name:
		if vsphereDataDisk, ok := dataDisk.(vsphere.DiskInfo); ok {
			device := fmt.Sprintf(VsphereScsiByPath, vsphereDataDisk.Index+1)
			diskSetupIgn, err := machineconfig.ForDiskSetup(role, device, label, path, diskSetup.Type)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create ignition to setup disks for master machines")
			}
			return diskSetupIgn, nil
		}
		return nil, errors.Errorf("unsupported vsphere data disk type")
	default:
		return nil, errors.Errorf("unsupported platform %q", ic.Platform.Name())
	}
	return nil, nil
}
