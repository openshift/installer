package ovirtclient

import (
	"strings"
	"sync"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// StorageDomainID is a specialized type for storage domain IDs.
type StorageDomainID string

// StorageDomainClient contains the portion of the goVirt API that deals with storage domains.
type StorageDomainClient interface {
	// ListStorageDomains lists all storage domains.
	ListStorageDomains(retries ...RetryStrategy) (StorageDomainList, error)
	// GetStorageDomain returns a single storage domain, or an error if the storage domain could not be found.
	GetStorageDomain(id StorageDomainID, retries ...RetryStrategy) (StorageDomain, error)
	// GetDiskFromStorageDomain returns a single disk from a specific storage domain, or an error if no disk can be found.
	GetDiskFromStorageDomain(id StorageDomainID, diskID DiskID, retries ...RetryStrategy) (result Disk, err error)
	// RemoveDiskFromStorageDomain removes a disk from a specific storage domain, but leaves the disk on other storage
	// domains if any. If the disk is not present on any more storage domains, the entire disk will be removed.
	RemoveDiskFromStorageDomain(id StorageDomainID, diskID DiskID, retries ...RetryStrategy) error
}

// StorageDomainData is the core of StorageDomain, providing only data access functions.
type StorageDomainData interface {
	// ID is the unique identified for the storage system connected to oVirt.
	ID() StorageDomainID
	// Name is the user-given name for the storage domain.
	Name() string
	// Available returns the number of available bytes on the storage domain
	Available() uint64
	// StorageType returns the type of the storage domain
	StorageType() StorageDomainType
	// Status returns the status of the storage domain. This status may be unknown if the storage domain is external.
	// Check ExternalStatus as well.
	Status() StorageDomainStatus
	// ExternalStatus returns the external status of a storage domain.
	ExternalStatus() StorageDomainExternalStatus
}

// StorageDomain represents a storage domain returned from the oVirt Engine API.
type StorageDomain interface {
	StorageDomainData
}

// StorageDomainList represents a list of storage domains.
type StorageDomainList []StorageDomain

// Filter applies the passed in filter function to all items in the current list and returns a new list with items that returned true.
func (list StorageDomainList) Filter(apply func(StorageDomain) bool) StorageDomainList {
	var filtered StorageDomainList
	for _, sd := range list {
		if apply(sd) {
			filtered = append(filtered, sd)
		}
	}
	return filtered
}

// StorageDomainType represents the type of the storage domain.
type StorageDomainType string

const (
	// StorageDomainTypeCinder represents a cinder host storage type.
	StorageDomainTypeCinder StorageDomainType = "cinder"
	// StorageDomainTypeFCP represents a fcp host storage type.
	StorageDomainTypeFCP StorageDomainType = "fcp"
	// StorageDomainTypeGlance represents a glance host storage type.
	StorageDomainTypeGlance StorageDomainType = "glance"
	// StorageDomainTypeGlusterFS represents a glusterfs host storage type.
	StorageDomainTypeGlusterFS StorageDomainType = "glusterfs"
	// StorageDomainTypeISCSI represents a iscsi host storage type.
	StorageDomainTypeISCSI StorageDomainType = "iscsi"
	// StorageDomainTypeLocalFS represents a localfs host storage type.
	StorageDomainTypeLocalFS StorageDomainType = "localfs"
	// StorageDomainTypeManagedBlockStorage represents a managed block storage host storage type.
	StorageDomainTypeManagedBlockStorage StorageDomainType = "managed_block_storage"
	// StorageDomainTypeNFS represents a nfs host storage type.
	StorageDomainTypeNFS StorageDomainType = "nfs"
	// StorageDomainTypePosixFS represents a posixfs host storage type.
	StorageDomainTypePosixFS StorageDomainType = "posixfs"
)

// StorageDomainTypeList is a list of possible StorageDomainTypes.
type StorageDomainTypeList []StorageDomainType

// FileStorageDomainTypeList is a list of possible StorageDomainTypes which are considered file storage.
type FileStorageDomainTypeList []StorageDomainType

// FileStorageDomainTypeValues returns all the StorageDomainTypes values which are considered file storage.
func FileStorageDomainTypeValues() FileStorageDomainTypeList {
	return []StorageDomainType{
		StorageDomainTypeGlusterFS,
		StorageDomainTypeLocalFS,
		StorageDomainTypeNFS,
		StorageDomainTypePosixFS,
	}
}

// StorageDomainTypeValues returns all possible StorageDomainTypeValues values.
func StorageDomainTypeValues() StorageDomainTypeList {
	return []StorageDomainType{
		StorageDomainTypeCinder,
		StorageDomainTypeFCP,
		StorageDomainTypeGlance,
		StorageDomainTypeGlusterFS,
		StorageDomainTypeISCSI,
		StorageDomainTypeLocalFS,
		StorageDomainTypeManagedBlockStorage,
		StorageDomainTypeNFS,
		StorageDomainTypePosixFS,
	}
}

// StorageDomainStatus represents the status a domain can be in. Either this status field, or the
// StorageDomainExternalStatus must be set.
//
// Note: this is not well documented due to missing source documentation. If you know something about these statuses
// please contribute here:
// https://github.com/oVirt/ovirt-engine-api-model/blob/master/src/main/java/types/StorageDomainStatus.java
type StorageDomainStatus string

// Validate returns an error if the storage domain status doesn't have a valid value.
func (s StorageDomainStatus) Validate() error {
	for _, format := range StorageDomainStatusValues() {
		if format == s {
			return nil
		}
	}
	return newError(
		EBadArgument,
		"invalid storage domain status: %s must be one of: %s",
		s,
		strings.Join(ImageFormatValues().Strings(), ", "),
	)
}

const (
	// StorageDomainStatusActivating indicates that the storage domain is currently activating and will soon be active.
	StorageDomainStatusActivating StorageDomainStatus = "activating"
	// StorageDomainStatusActive is the normal status for a storage domain when it's working.
	StorageDomainStatusActive StorageDomainStatus = "active"
	// StorageDomainStatusDetaching is the status when it is being disconnected.
	StorageDomainStatusDetaching StorageDomainStatus = "detaching"
	// StorageDomainStatusInactive is an undocumented status of the storage domain.
	StorageDomainStatusInactive StorageDomainStatus = "inactive"
	// StorageDomainStatusLocked is an undocumented status of the storage domain.
	StorageDomainStatusLocked StorageDomainStatus = "locked"
	// StorageDomainStatusMaintenance is an undocumented status of the storage domain.
	StorageDomainStatusMaintenance StorageDomainStatus = "maintenance"
	// StorageDomainStatusMixed is an undocumented status of the storage domain.
	StorageDomainStatusMixed StorageDomainStatus = "mixed"
	// StorageDomainStatusPreparingForMaintenance is an undocumented status of the storage domain.
	StorageDomainStatusPreparingForMaintenance StorageDomainStatus = "preparing_for_maintenance"
	// StorageDomainStatusUnattached is an undocumented status of the storage domain.
	StorageDomainStatusUnattached StorageDomainStatus = "unattached"
	// StorageDomainStatusUnknown is an undocumented status of the storage domain.
	StorageDomainStatusUnknown StorageDomainStatus = "unknown"
	// StorageDomainStatusNA indicates that the storage domain does not have a status. Please check the external status
	// instead.
	StorageDomainStatusNA StorageDomainStatus = ""
)

// StorageDomainStatusList is a list of StorageDomainStatus.
type StorageDomainStatusList []StorageDomainStatus

// StorageDomainStatusValues returns all possible StorageDomainStatus values.
func StorageDomainStatusValues() StorageDomainStatusList {
	return []StorageDomainStatus{
		StorageDomainStatusActivating,
		StorageDomainStatusActive,
		StorageDomainStatusDetaching,
		StorageDomainStatusInactive,
		StorageDomainStatusLocked,
		StorageDomainStatusMaintenance,
		StorageDomainStatusMixed,
		StorageDomainStatusPreparingForMaintenance,
		StorageDomainStatusUnattached,
		StorageDomainStatusUnknown,
		StorageDomainStatusNA,
	}
}

// Strings creates a string list of the values.
func (l StorageDomainStatusList) Strings() []string {
	result := make([]string, len(l))
	for i, status := range l {
		result[i] = string(status)
	}
	return result
}

// StorageDomainExternalStatus represents the status of an external storage domain. This status is updated externally.
//
// Note: this is not well-defined as the oVirt model has only a very generic description. See
// https://github.com/oVirt/ovirt-engine-api-model/blob/9869596c298925538d510de5019195b488970738/src/main/java/types/ExternalStatus.java
// for details.
type StorageDomainExternalStatus string

const (
	// StorageDomainExternalStatusNA represents an external status that is not applicable.
	// Most likely, the status should be obtained from StorageDomainStatus, since the
	// storage domain in question is not an external storage.
	StorageDomainExternalStatusNA StorageDomainExternalStatus = ""
	// StorageDomainExternalStatusError indicates an error state.
	StorageDomainExternalStatusError StorageDomainExternalStatus = "error"
	// StorageDomainExternalStatusFailure indicates a failure state.
	StorageDomainExternalStatusFailure StorageDomainExternalStatus = "failure"
	// StorageDomainExternalStatusInfo indicates an OK status, but there is information available for the administrator
	// that might be relevant.
	StorageDomainExternalStatusInfo StorageDomainExternalStatus = "info"
	// StorageDomainExternalStatusOk indicates a working status.
	StorageDomainExternalStatusOk StorageDomainExternalStatus = "ok"
	// StorageDomainExternalStatusWarning indicates that the storage domain has warnings that may be relevant for the
	// administrator.
	StorageDomainExternalStatusWarning StorageDomainExternalStatus = "warning"
)

// StorageDomainExternalStatusList is a list of StorageDomainStatus.
type StorageDomainExternalStatusList []StorageDomainExternalStatus

// StorageDomainExternalStatusValues returns all possible StorageDomainExternalStatus values.
func StorageDomainExternalStatusValues() StorageDomainExternalStatusList {
	return []StorageDomainExternalStatus{
		StorageDomainExternalStatusNA,
		StorageDomainExternalStatusError,
		StorageDomainExternalStatusFailure,
		StorageDomainExternalStatusInfo,
		StorageDomainExternalStatusOk,
		StorageDomainExternalStatusWarning,
	}
}

// Strings creates a string list of the values.
func (l StorageDomainExternalStatusList) Strings() []string {
	result := make([]string, len(l))
	for i, status := range l {
		result[i] = string(status)
	}
	return result
}

func convertSDKStorageDomain(sdkStorageDomain *ovirtsdk4.StorageDomain, client Client) (StorageDomain, error) {
	id, ok := sdkStorageDomain.Id()
	if !ok {
		return nil, newError(EFieldMissing, "failed to fetch ID of storage domain")
	}
	name, ok := sdkStorageDomain.Name()
	if !ok {
		return nil, newError(EFieldMissing, "failed to fetch name of storage domain")
	}
	available, ok := sdkStorageDomain.Available()
	if !ok {
		// If this is not OK the status probably doesn't allow for reading disk space (e.g. unattached), so we return 0.
		available = 0
	}
	if available < 0 {
		return nil, newError(EBug, "invalid available bytes returned from storage domain: %d", available)
	}
	storage, ok := sdkStorageDomain.Storage()
	if !ok {
		return nil, newError(EFieldMissing, "failed to fetch hostStorage of storage domain")
	}
	storageType, ok := storage.Type()
	if !ok {
		return nil, newError(EFieldMissing, "failed to fetch storage type of storage domain")
	}
	// It is OK for the storage domain status to not be present if the external status is present.
	status, _ := sdkStorageDomain.Status()
	// It is OK for the storage domain external status to not be present if the status is present.
	externalStatus, _ := sdkStorageDomain.ExternalStatus()
	if status == "" && externalStatus == "" {
		return nil, newError(EFieldMissing, "neither the status nor the external status is set for storage domain %s", id)
	}

	return &storageDomain{
		client: client,

		id:             StorageDomainID(id),
		name:           name,
		available:      uint64(available),
		storageType:    StorageDomainType(storageType),
		status:         StorageDomainStatus(status),
		externalStatus: StorageDomainExternalStatus(externalStatus),
	}, nil
}

type storageDomain struct {
	client Client

	id             StorageDomainID
	name           string
	available      uint64
	storageType    StorageDomainType
	status         StorageDomainStatus
	externalStatus StorageDomainExternalStatus
}

func (s storageDomain) ID() StorageDomainID {
	return s.id
}

func (s storageDomain) Name() string {
	return s.name
}

func (s storageDomain) Available() uint64 {
	return s.available
}

func (s storageDomain) StorageType() StorageDomainType {
	return s.storageType
}

func (s storageDomain) Status() StorageDomainStatus {
	return s.status
}

func (s storageDomain) ExternalStatus() StorageDomainExternalStatus {
	return s.externalStatus
}

type storageDomainDiskWait struct {
	client        *oVirtClient
	disk          Disk
	storageDomain StorageDomain
	correlationID string
	lock          *sync.Mutex
}

func (d *storageDomainDiskWait) Disk() Disk {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.disk
}

func (d *storageDomainDiskWait) Wait(retries ...RetryStrategy) (Disk, error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(d.client))
	d.lock.Lock()
	diskID := d.disk.ID()
	storageDomainID := d.storageDomain.ID()
	d.lock.Unlock()

	if _, err := d.client.WaitForDiskOK(diskID, retries...); err != nil {
		return nil, err
	}

	if err := d.client.waitForJobFinished(d.correlationID, retries); err != nil {
		return nil, err
	}

	disk, err := d.client.GetDiskFromStorageDomain(storageDomainID, diskID)

	if disk != nil {
		d.disk = disk
	}
	return disk, err
}
