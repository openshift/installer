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
)

// NodeDiskSetup determines the path per disk type, and per platform and role, runs ForDiskSetup.
func NodeDiskSetup(installConfig *installconfig.InstallConfig, role string, diskSetup types.Disk, dataDisk any) (*v1.MachineConfig, error) {
	var path string

	ic := installConfig.Config

	label := string(diskSetup.Type)

	switch diskSetup.Type {
	case types.Etcd:
		path = "/var/lib/etcd"
	case types.Swap:
		path = ""
	case types.UserDefined:
		path = diskSetup.UserDefined.MountPath
		label = diskSetup.UserDefined.PlatformDiskID
	}

	switch ic.Platform.Name() {
	case azuretypes.Name:
		if installConfig.Config.EnabledFeatureGates().Enabled(features.FeatureGateAzureMultiDisk) {
			if azureDataDisk, ok := dataDisk.(v1beta1.DataDisk); ok {
				device := fmt.Sprintf("/dev/disk/azure/scsi1/lun%d", *azureDataDisk.Lun)
				diskSetupIgn, err := machineconfig.ForDiskSetup(role, device, label, path, diskSetup.Type)
				if err != nil {
					return nil, errors.Wrap(err, "failed to create ignition to setup disks for master machines")
				}
				return diskSetupIgn, nil
			}
			return nil, errors.Errorf("unsupported azure data disk type")
		}
	case vspheretypes.Name:
		if installConfig.Config.EnabledFeatureGates().Enabled(features.FeatureGateVSphereMultiDisk) {
			if vsphereDataDisk, ok := dataDisk.(vsphere.DiskInfo); ok {
				// We need to find the index of the datadisk in the array.  Each disk is added in order to the VM so
				// we'll map to that location.  First disk is OS disk so add 1 to index for scsi location
				device := fmt.Sprintf(VsphereScsiByPath, vsphereDataDisk.Index+1)
				diskSetupIgn, err := machineconfig.ForDiskSetup(role, device, label, path, diskSetup.Type)
				if err != nil {
					return nil, errors.Wrap(err, "failed to create ignition to setup disks for master machines")
				}
				return diskSetupIgn, nil
			}
			return nil, errors.Errorf("unsupported vsphere data disk type")
		}
	default:
		return nil, errors.Errorf("unsupported platform %q", ic.Platform.Name())
	}
	return nil, nil
}
