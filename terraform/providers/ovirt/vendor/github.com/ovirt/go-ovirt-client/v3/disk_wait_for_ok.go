package ovirtclient

import (
	"fmt"
	"time"
)

// WaitForDiskOK waits for a disk to be in the OK status, then additionally queries the job that was in progress with
// the correlation ID. This is necessary because the disk returns OK status before the job has actually finished,
// resulting in a "disk locked" error on subsequent operations. It uses checkDiskOk as an underlying function.
func (o *oVirtClient) WaitForDiskOK(diskID DiskID, retries ...RetryStrategy) (disk Disk, err error) {
	err = retry(
		fmt.Sprintf("waiting for disk %s to become OK", diskID),
		o.logger,
		retries,
		func() error {
			disk, err = o.checkDiskOK(diskID)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	return disk, nil
}

// checkDiskOK fetches the disk for the transfer and checks if it is in the OK status. It returns an EPending error if
// it is not.
func (o *oVirtClient) checkDiskOK(diskID DiskID) (Disk, error) {
	disk, err := o.GetDisk(diskID)
	if err != nil {
		return nil, err
	}
	switch disk.Status() {
	case DiskStatusOK:
		return disk, nil
	case DiskStatusLocked:
		return nil, newError(EPending, "disk status is %s, not %s", disk.Status(), DiskStatusOK)
	default:
		return nil, newError(EUnexpectedDiskStatus, "disk status is %s, not %s", disk.Status(), DiskStatusOK)
	}
}

// WaitForDiskOK waits for a disk to be in the OK status, then additionally queries the job that was in progress with
// the correlation ID. This is necessary because the disk returns OK status before the job has actually finished,
// resulting in a "disk locked" error on subsequent operations. It uses checkDiskOk as an underlying function.
func (m *mockClient) WaitForDiskOK(diskID DiskID, retries ...RetryStrategy) (Disk, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	disk, ok := m.disks[diskID]

	if !ok {
		return nil, newError(ENotFound, "Disk with ID %s not found", diskID)
	}
	time.Sleep(2 * time.Second)
	disk.status = DiskStatusOK

	return disk, nil
}
