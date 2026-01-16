package azure

import (
	"fmt"

	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/types/azure"
)

// ConvertVMIdentityType converts the local VMIdentityType to the capz VMIdentity type.
func ConvertVMIdentityType(id azure.VMIdentityType) capz.VMIdentity {
	return capz.VMIdentity(id)
}

// ConvertSubnetRole converts the local SubnetRole to the capz SubnetRole type.
func ConvertSubnetRole(role azure.SubnetRole) capz.SubnetRole {
	return capz.SubnetRole(role)
}

// ConvertBootDiagnosticsStorageAccountType converts the local BootDiagnosticsStorageAccountType
// to the capz BootDiagnosticsStorageAccountType.
func ConvertBootDiagnosticsStorageAccountType(t azure.BootDiagnosticsStorageAccountType) capz.BootDiagnosticsStorageAccountType {
	return capz.BootDiagnosticsStorageAccountType(t)
}

// ConvertDataDisks converts a slice of local DataDisk to capz DataDisk types.
func ConvertDataDisks(disks []azure.DataDisk) []capz.DataDisk {
	if disks == nil {
		return nil
	}

	result := make([]capz.DataDisk, len(disks))
	for i, d := range disks {
		result[i] = ConvertDataDisk(d)
	}
	return result
}

// ConvertDataDisk converts a local DataDisk to a capz DataDisk.
func ConvertDataDisk(d azure.DataDisk) capz.DataDisk {
	disk := capz.DataDisk{
		NameSuffix:  d.NameSuffix,
		DiskSizeGB:  d.DiskSizeGB,
		Lun:         d.Lun,
		CachingType: d.CachingType,
	}

	if d.ManagedDisk != nil {
		disk.ManagedDisk = &capz.ManagedDiskParameters{
			StorageAccountType: d.ManagedDisk.StorageAccountType,
		}
		if d.ManagedDisk.DiskEncryptionSet != nil {
			disk.ManagedDisk.DiskEncryptionSet = &capz.DiskEncryptionSetParameters{
				ID: diskEncryptionSetResourceID(d.ManagedDisk.DiskEncryptionSet),
			}
		}
	}

	return disk
}

// diskEncryptionSetResourceID returns the Azure resource ID for a disk encryption set.
func diskEncryptionSetResourceID(des *azure.DiskEncryptionSet) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/diskEncryptionSets/%s",
		des.SubscriptionID, des.ResourceGroup, des.Name)
}
