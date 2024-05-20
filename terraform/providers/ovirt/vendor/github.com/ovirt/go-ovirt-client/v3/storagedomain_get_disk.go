package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) GetDiskFromStorageDomain(id StorageDomainID, diskID DiskID, retries ...RetryStrategy) (result Disk, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting disk %s from storage domain %s", diskID, id),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.SystemService().StorageDomainsService().
				StorageDomainService(string(id)).DisksService().DiskService(string(diskID)).Get().Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.Disk()
			if !ok {
				return newError(
					ENotFound,
					"disk %s not found on storage domain ID %s",
					diskID,
					id,
				)
			}
			result, err = convertSDKDisk(sdkObject, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert disk %s",
					diskID,
				)
			}
			return nil
		})
	return result, err
}

func (m *mockClient) GetDiskFromStorageDomain(id StorageDomainID, diskID DiskID, _ ...RetryStrategy) (Disk, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if disk, ok := m.disks[diskID]; ok {
		for _, domain := range disk.storageDomainIDs {
			if domain == id {
				return disk, nil
			}
		}
		return nil, newError(ENotFound, "disk %s doesnt exist in storage domain %s", diskID, id)
	}
	return nil, newError(ENotFound, "disk %s doesnt exists", diskID)
}
