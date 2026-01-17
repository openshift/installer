package azure

import (
	"fmt"

	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/types/azure"
)

// ConvertDataDisks converts a slice of local DataDisk to capz DataDisk types.
func ConvertDataDisks(disks []azure.DataDisk) []capz.DataDisk {
	if disks == nil {
		return nil
	}

	result := make([]capz.DataDisk, len(disks))
	for i, d := range disks {
		result[i] = convertDataDisk(d)
	}
	return result
}

// convertDataDisk converts a local DataDisk to a capz DataDisk.
func convertDataDisk(d azure.DataDisk) capz.DataDisk {
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
