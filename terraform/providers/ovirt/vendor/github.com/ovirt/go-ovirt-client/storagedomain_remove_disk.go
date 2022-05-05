package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) RemoveDiskFromStorageDomain(id StorageDomainID, diskID DiskID, retries ...RetryStrategy) (err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("removing disk %s from storage domain %s", diskID, id),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.SystemService().StorageDomainsService().
				StorageDomainService(string(id)).DisksService().DiskService(string(diskID)).Remove().Send()
			if err != nil {
				o.logger.Infof("error removing disk..")
				return err
			}

			return nil
		})
	return
}

func (m *mockClient) RemoveDiskFromStorageDomain(id StorageDomainID, diskID DiskID, _ ...RetryStrategy) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.disks[diskID]; !ok {
		return newError(ENotFound, "disk with ID %s not found", diskID)
	}

	domains := m.disks[diskID].storageDomainIDs

	// if there is only 1 domain just delete the disk
	if len(domains) == 1 {
		delete(m.disks, diskID)
		return nil
	}

	// remove the storagedomain from the disk slice
	for i, sdomain := range domains {
		if sdomain == id {
			// gocritic will complain on the following line due to appendAssign, but that's legit here
			m.disks[diskID].storageDomainIDs = append(domains[:i], domains[i+1:]...) //nolint:gocritic
			return nil
		}
	}
	return newError(ENotFound, "disk %s is not found in StorageDomain %s", diskID, id)
}
