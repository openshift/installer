package ovirtclient

import (
	"fmt"
	"sync"
	"time"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) CopyTemplateDiskToStorageDomain(
	diskID DiskID,
	storageDomainID StorageDomainID,
	retries ...RetryStrategy) (result Disk, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	progress, err := o.StartCopyTemplateDiskToStorageDomain(diskID, storageDomainID, retries...)
	if err != nil {
		return nil, err
	}

	return progress.Wait()
}

func (o *oVirtClient) StartCopyTemplateDiskToStorageDomain(
	diskID DiskID,
	storageDomainID StorageDomainID,
	retries ...RetryStrategy) (DiskUpdate, error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	correlationID := fmt.Sprintf("template_disk_copy_%s", generateRandomID(5, o.nonSecureRandom))
	sdkStorageDomain := ovirtsdk.NewStorageDomainBuilder().Id(string(storageDomainID))
	sdkDisk := ovirtsdk.NewDiskBuilder().Id(string(diskID))
	storageDomain, _ := o.GetStorageDomain(storageDomainID)
	disk, _ := o.GetDisk(diskID)

	err := retry(
		fmt.Sprintf("copying disk %s to storage domain %s", diskID, storageDomainID),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.
				SystemService().
				DisksService().
				DiskService(string(diskID)).
				Copy().
				StorageDomain(sdkStorageDomain.MustBuild()).
				Disk(sdkDisk.MustBuild()).
				Query("correlation_id", correlationID).
				Send()

			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &storageDomainDiskWait{
		client:        o,
		disk:          disk,
		storageDomain: storageDomain,
		correlationID: correlationID,
		lock:          &sync.Mutex{},
	}, nil
}

func (m *mockClient) CopyTemplateDiskToStorageDomain(
	diskID DiskID,
	storageDomainID StorageDomainID,
	retries ...RetryStrategy) (result Disk, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	disk, ok := m.disks[diskID]

	if !ok {
		return nil, newError(ENotFound, "disk with ID %s not found", diskID)
	}
	if err := disk.Lock(); err != nil {
		return nil, err
	}
	update := &mockDiskCopy{
		client:          m,
		disk:            disk,
		storageDomainID: storageDomainID,
		done:            make(chan struct{}),
	}
	defer update.do()
	return disk, nil
}

type mockDiskCopy struct {
	client          *mockClient
	disk            *diskWithData
	storageDomainID StorageDomainID
	done            chan struct{}
}

func (c *mockDiskCopy) Disk() Disk {
	c.client.lock.Lock()
	defer c.client.lock.Unlock()

	return c.disk
}

func (c *mockDiskCopy) Wait(_ ...RetryStrategy) (Disk, error) {
	<-c.done

	return c.disk, nil
}

func (c *mockDiskCopy) do() {
	// Sleep to trigger potential race conditions / improper status handling.
	time.Sleep(time.Second)
	c.client.disks[c.disk.ID()] = c.disk
	c.client.disks[c.disk.ID()].storageDomainIDs = append(c.client.disks[c.disk.ID()].storageDomainIDs, c.storageDomainID)
	close(c.done)
}
