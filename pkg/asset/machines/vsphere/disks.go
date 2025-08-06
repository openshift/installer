package vsphere

import "github.com/openshift/installer/pkg/types/vsphere"

// DiskInfo is used to encapsulate additional DiskInfo relating to the overall configuration due to how Disk is passed
// into helper methods.
type DiskInfo struct {
	Index int
	Disk  vsphere.DataDisk
}
