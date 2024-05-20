package ovirtclient

import (
	"fmt"
	"sync"
	"time"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) UpdateDisk(id DiskID, params UpdateDiskParameters, retries ...RetryStrategy) (
	result Disk,
	err error,
) {
	progress, err := o.StartUpdateDisk(id, params, retries...)
	if err != nil {
		return nil, err
	}
	return progress.Wait(retries...)
}

func (o *oVirtClient) StartUpdateDisk(id DiskID, params UpdateDiskParameters, retries ...RetryStrategy) (
	DiskUpdate,
	error,
) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))

	sdkDisk := ovirtsdk.NewDiskBuilder().Id(string(id))
	if alias := params.Alias(); alias != nil {
		sdkDisk.Alias(*alias)
	}
	if provisionedSize := params.ProvisionedSize(); provisionedSize != nil {
		sdkDisk.ProvisionedSize(int64(*provisionedSize))
	}
	correlationID := fmt.Sprintf("disk_update_%s", generateRandomID(5, o.nonSecureRandom))

	var disk Disk

	err := retry(
		fmt.Sprintf("updating disk %s", id),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.
				SystemService().
				DisksService().
				DiskService(string(id)).
				Update().
				Disk(sdkDisk.MustBuild()).
				Query("correlation_id", correlationID).
				Send()
			if err != nil {
				return err
			}
			sdkDisk, ok := response.Disk()
			if !ok {
				return newError(
					EFieldMissing,
					"missing disk object from disk update response",
				)
			}
			disk, err = convertSDKDisk(sdkDisk, o)
			if err != nil {
				return wrap(err, EUnidentified, "failed to convert SDK disk object")
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &diskWait{
		client:        o,
		disk:          disk,
		correlationID: correlationID,
		lock:          &sync.Mutex{},
	}, nil
}

func (m *mockClient) UpdateDisk(id DiskID, params UpdateDiskParameters, retries ...RetryStrategy) (Disk, error) {
	progress, err := m.StartUpdateDisk(id, params, retries...)
	if err != nil {
		return progress.Disk(), err
	}
	return progress.Wait(retries...)
}

func (m *mockClient) StartUpdateDisk(id DiskID, params UpdateDiskParameters, _ ...RetryStrategy) (
	DiskUpdate,
	error,
) {
	m.lock.Lock()
	defer m.lock.Unlock()

	var err error
	disk, ok := m.disks[id]
	if !ok {
		return nil, newError(ENotFound, "disk with ID %s not found", id)
	}
	if err := disk.Lock(); err != nil {
		return nil, err
	}
	if alias := params.Alias(); alias != nil {
		disk = disk.WithAlias(alias)
	}
	if ps := params.ProvisionedSize(); ps != nil {
		disk, err = disk.withProvisionedSize(*ps)
		if err != nil {
			return nil, err
		}
	}
	update := &mockDiskUpdate{
		client: m,
		disk:   disk,
		done:   make(chan struct{}),
	}
	defer update.do()
	return update, nil
}

type mockDiskUpdate struct {
	client *mockClient
	disk   *diskWithData
	done   chan struct{}
}

func (c *mockDiskUpdate) Disk() Disk {
	c.client.lock.Lock()
	defer c.client.lock.Unlock()

	return c.disk
}

func (c *mockDiskUpdate) Wait(_ ...RetryStrategy) (Disk, error) {
	<-c.done

	return c.disk, nil
}

func (c *mockDiskUpdate) do() {
	// Sleep to trigger potential race conditions / improper status handling.
	time.Sleep(time.Second)

	c.client.disks[c.disk.ID()] = c.disk
	c.disk.Unlock()

	close(c.done)
}
