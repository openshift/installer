package ovirtclient

import (
	"sync"

	"github.com/google/uuid"
)

// diskWithData adds the ability to store the data directly in the disk for mocking purposes.
type diskWithData struct {
	disk
	lock *sync.Mutex
	data []byte
}

func (d *diskWithData) Lock() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.status != DiskStatusOK {
		return newError(EDiskLocked, "disk %s is %s", d.id, d.status)
	}
	d.status = DiskStatusLocked
	return nil
}

func (d *diskWithData) Unlock() {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.status = DiskStatusOK
}

func (d *diskWithData) WithAlias(alias *string) *diskWithData {
	return &diskWithData{
		disk{
			client:           d.client,
			id:               d.id,
			alias:            *alias,
			provisionedSize:  d.provisionedSize,
			format:           d.format,
			storageDomainIDs: d.storageDomainIDs,
			status:           d.status,
			totalSize:        d.totalSize,
			sparse:           d.sparse,
		},
		d.lock,
		d.data,
	}
}

func (d *diskWithData) withProvisionedSize(ps uint64) (*diskWithData, error) {
	if d.provisionedSize > ps {
		return nil, newError(
			EBadArgument,
			"Cannot edit Virtual Disk. New disk size must be larger than the current disk size",
		)
	}
	return &diskWithData{
		disk{
			client:           d.client,
			id:               d.id,
			alias:            d.alias,
			provisionedSize:  ps,
			format:           d.format,
			storageDomainIDs: d.storageDomainIDs,
			status:           d.status,
			totalSize:        ps,
			sparse:           d.sparse,
		},
		d.lock,
		d.data,
	}, nil
}

// clone is an internal function that makes a copy of the disk object with a new UUID.
func (d *diskWithData) clone(sparse *bool) *diskWithData {
	if sparse == nil {
		sparse = &d.sparse
	}
	return &diskWithData{
		disk{
			d.client,
			DiskID(uuid.NewString()),
			d.alias,
			d.provisionedSize,
			d.format,
			d.storageDomainIDs,
			d.status,
			d.totalSize,
			*sparse,
		},
		&sync.Mutex{},
		d.data,
	}
}
