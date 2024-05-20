package ovirtclient

import (
	"sync"
	"time"
)

func (m *mockClient) StartCreateDisk(
	storageDomainID StorageDomainID,
	format ImageFormat,
	size uint64,
	params CreateDiskOptionalParameters,
	_ ...RetryStrategy,
) (DiskCreation, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	disk, err := m.createDisk(storageDomainID, format, size, params)
	if err != nil {
		return nil, err
	}

	creation := &mockDiskCreation{
		client: m,
		disk:   disk,
		done:   make(chan struct{}),
	}
	creation.do()
	return creation, nil
}

func (m *mockClient) createDisk(
	storageDomainID StorageDomainID,
	format ImageFormat,
	size uint64,
	params CreateDiskOptionalParameters,
) (*diskWithData, error) {
	if err := validateDiskCreationParameters(format, size); err != nil {
		return nil, err
	}

	if _, ok := m.storageDomains[storageDomainID]; !ok {
		return nil, newError(ENotFound, "storage domain with ID %s not found", storageDomainID)
	}

	disk := &diskWithData{
		disk: disk{
			client:           m,
			id:               DiskID(m.GenerateUUID()),
			format:           format,
			provisionedSize:  size,
			totalSize:        size,
			storageDomainIDs: []StorageDomainID{storageDomainID},
			status:           DiskStatusLocked,
		},
		lock: &sync.Mutex{},
		data: nil,
	}

	if params != nil {
		if alias := params.Alias(); alias != "" {
			disk.disk.alias = alias
		}
		if sparse := params.Sparse(); sparse != nil {
			disk.disk.sparse = *sparse
		}
	}

	m.disks[disk.id] = disk

	return disk, nil
}

func (m *mockClient) CreateDisk(
	storageDomainID StorageDomainID,
	format ImageFormat,
	size uint64,
	params CreateDiskOptionalParameters,
	retries ...RetryStrategy,
) (Disk, error) {
	result, err := m.StartCreateDisk(storageDomainID, format, size, params, retries...)
	if err != nil {
		return nil, err
	}

	return result.Wait()
}

type mockDiskCreation struct {
	client *mockClient
	disk   *diskWithData
	done   chan struct{}
}

func (c *mockDiskCreation) Disk() Disk {
	c.client.lock.Lock()
	defer c.client.lock.Unlock()

	return c.disk
}

func (c *mockDiskCreation) Wait(_ ...RetryStrategy) (Disk, error) {
	<-c.done

	return c.disk, nil
}

func (c *mockDiskCreation) do() {
	// Sleep to trigger potential race conditions / improper status handling.
	time.Sleep(time.Second)

	c.disk.Unlock()

	close(c.done)
}
