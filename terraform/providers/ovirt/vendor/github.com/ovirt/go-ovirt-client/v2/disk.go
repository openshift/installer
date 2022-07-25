package ovirtclient

import (
	"io"
	"strings"
	"sync"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

//go:generate go run scripts/rest/rest.go -i "Disk" -n "disk" -T DiskID

// DiskID is the identifier for disks.
type DiskID string

// DiskClient is the client interface part that deals with disks.
type DiskClient interface {
	// StartImageUpload uploads an image file into a disk. The actual upload takes place in the
	// background and can be tracked using the returned UploadImageProgress object.
	//
	// Parameters are as follows:
	//
	// - alias: this is the name used for the uploaded image.
	// - storageDomainID: this is the UUID of the storage domain that the image should be uploaded to.
	// - sparse: use sparse provisioning
	// - size: this is the file size of the image. This must match the bytes read.
	// - reader: this is the source of the image data.
	// - retries: a set of optional retry options.
	//
	// You can wait for the upload to complete using the Done() method:
	//
	//     progress, err := cli.StartImageUpload(...)
	//     if err != nil {
	//         //...
	//     }
	//     <-progress.Done()
	//
	// After the upload is complete you can check the Err() method if it completed successfully:
	//
	//     if err := progress.Err(); err != nil {
	//         //...
	//     }
	//
	// Deprecated: Use StartUploadToNewDisk instead.
	StartImageUpload(
		alias string,
		storageDomainID StorageDomainID,
		sparse bool,
		size uint64,
		reader io.ReadSeekCloser,
		retries ...RetryStrategy,
	) (UploadImageProgress, error)

	// StartUploadToNewDisk uploads an image file into a disk. The actual upload takes place in the
	// background and can be tracked using the returned UploadImageProgress object. If the process fails a removal
	// of the created disk is attempted.
	//
	// Parameters are as follows:
	//
	// - storageDomainID: this is the UUID of the storage domain that the image should be uploaded to.
	// - format: format of the created disk. This does not necessarily have to be identical to the format of the image
	//   being uploaded as the oVirt engine converts images on upload.
	// - size: file size of the uploaded image on the disk.
	// - reader: this is the source of the image data. It is a reader that must support seek and close operations.
	// - retries: a set of optional retry options.
	//
	// You can wait for the upload to complete using the Done() method:
	//
	//     progress, err := cli.StartUploadToNewDisk(...)
	//     if err != nil {
	//         //...
	//     }
	//     <-progress.Done()
	//
	// After the upload is complete you can check the Err() method if it completed successfully:
	//
	//     if err := progress.Err(); err != nil {
	//         //...
	//     }
	StartUploadToNewDisk(
		storageDomainID StorageDomainID,
		format ImageFormat,
		size uint64,
		params CreateDiskOptionalParameters,
		reader io.ReadSeekCloser,
		retries ...RetryStrategy,
	) (UploadImageProgress, error)

	// UploadImage is identical to StartImageUpload, but waits until the upload is complete. It returns the disk ID
	// as a result, or the error if one happened.
	//
	// Deprecated: Use UploadToNewDisk instead.
	UploadImage(
		alias string,
		storageDomainID StorageDomainID,
		sparse bool,
		size uint64,
		reader io.ReadSeekCloser,
		retries ...RetryStrategy,
	) (UploadImageResult, error)

	// UploadToNewDisk is identical to StartUploadToNewDisk, but waits until the upload is complete. It
	// returns the disk ID as a result, or the error if one happened.
	UploadToNewDisk(
		storageDomainID StorageDomainID,
		format ImageFormat,
		size uint64,
		params CreateDiskOptionalParameters,
		reader io.ReadSeekCloser,
		retries ...RetryStrategy,
	) (UploadImageResult, error)

	// StartUploadToDisk uploads a disk image to an existing disk. The actual upload takes place in the background
	// and can be tracked using the returned UploadImageProgress object. Parameters are as follows:
	//
	// - diskID: ID of the disk to upload to.
	// - reader This is the source of the image data.
	// - retries: A set of optional retry options.
	//
	// You can wait for the upload to complete using the Done() method:
	//
	//     progress, err := cli.StartUploadToDisk(...)
	//     if err != nil {
	//         //...
	//     }
	//     <-progress.Done()
	//
	// After the upload is complete you can check the Err() method if it completed successfully:
	//
	//     if err := progress.Err(); err != nil {
	//         //...
	//     }
	StartUploadToDisk(
		diskID DiskID,
		size uint64,
		reader io.ReadSeekCloser,
		retries ...RetryStrategy,
	) (UploadImageProgress, error)

	// UploadToDisk runs StartUploadDisk and then waits for the upload to complete. It returns an error if the upload
	// failed despite retries.
	//
	// Parameters are as follows:
	//
	// - diskID: ID of the disk to upload to.
	// - size: size of the file on disk.
	// - reader: this is the source of the image data. The format is automatically determined from the file being
	//   uploaded. The reader must support seeking and close.
	// - retries: a set of optional retry options.
	UploadToDisk(
		diskID DiskID,
		size uint64,
		reader io.ReadSeekCloser,
		retries ...RetryStrategy,
	) error

	// StartImageDownload starts the download of the image file of a specific disk.
	// The caller can then wait for the initialization using the Initialized() call:
	//
	//     <-download.Initialized()
	//
	// Alternatively, the downloader can use the Read() function to wait for the download to become available
	// and then read immediately.
	//
	// The caller MUST close the returned reader, otherwise the disk will remain locked in the oVirt engine.
	//
	// Deprecated: please use StartDownloadDisk instead.
	StartImageDownload(
		diskID DiskID,
		format ImageFormat,
		retries ...RetryStrategy,
	) (ImageDownload, error)

	// StartDownloadDisk starts the download of the image file of a specific disk.
	// The caller can then wait for the initialization using the Initialized() call:
	//
	//     <-download.Initialized()
	//
	// Alternatively, the downloader can use the Read() function to wait for the download to become available
	// and then read immediately.
	//
	// The caller MUST close the returned reader, otherwise the disk will remain locked in the oVirt engine.
	StartDownloadDisk(
		diskID DiskID,
		format ImageFormat,
		retries ...RetryStrategy,
	) (ImageDownload, error)

	// DownloadImage runs StartDownloadDisk, then waits for the download to be ready before returning the reader.
	// The caller MUST close the ImageDownloadReader in order to properly unlock the disk in the oVirt engine.
	//
	// Deprecated: please use DownloadDisk instead.
	DownloadImage(
		diskID DiskID,
		format ImageFormat,
		retries ...RetryStrategy,
	) (ImageDownloadReader, error)

	// DownloadDisk runs StartDownloadDisk, then waits for the download to be ready before returning the reader.
	// The caller MUST close the ImageDownloadReader in order to properly unlock the disk in the oVirt engine.
	DownloadDisk(
		diskID DiskID,
		format ImageFormat,
		retries ...RetryStrategy,
	) (ImageDownloadReader, error)

	// StartCreateDisk starts creating an empty disk with the specified parameters and returns a DiskCreation object,
	// which can be queried for completion. Optional parameters can be created using CreateDiskParams().
	StartCreateDisk(
		storageDomainID StorageDomainID,
		format ImageFormat,
		size uint64,
		params CreateDiskOptionalParameters,
		retries ...RetryStrategy,
	) (DiskCreation, error)

	// CreateDisk is a shorthand for calling StartCreateDisk, and then waiting for the disk creation to complete.
	// Optional parameters can be created using CreateDiskParams().
	//
	// Caution! The CreateDisk method may return both an error and a disk that has been created, but has not reached
	// the ready state. Since the disk is almost certainly in a locked state, this may mean that there is a disk left
	// behind.
	CreateDisk(
		storageDomainID StorageDomainID,
		format ImageFormat,
		size uint64,
		params CreateDiskOptionalParameters,
		retries ...RetryStrategy,
	) (Disk, error)

	// StartUpdateDisk sends the disk update request to the oVirt API and returns a DiskUpdate
	// object, which can be used to wait for the update to complete. Use UpdateDiskParams to
	// obtain a builder for the parameters structure.
	StartUpdateDisk(
		id DiskID,
		params UpdateDiskParameters,
		retries ...RetryStrategy,
	) (DiskUpdate, error)

	// UpdateDisk updates the specified disk with the specified parameters. Use UpdateDiskParams to
	// obtain a builder for the parameters structure.
	UpdateDisk(
		id DiskID,
		params UpdateDiskParameters,
		retries ...RetryStrategy,
	) (Disk, error)

	// ListDisks lists all disks.
	ListDisks(retries ...RetryStrategy) ([]Disk, error)
	// GetDisk fetches a disk with a specific ID from the oVirt Engine.
	GetDisk(diskID DiskID, retries ...RetryStrategy) (Disk, error)
	// ListDisksByAlias fetches a disks with a specific name from the oVirt Engine.
	ListDisksByAlias(alias string, retries ...RetryStrategy) ([]Disk, error)
	// RemoveDisk removes a disk with a specific ID.
	RemoveDisk(diskID DiskID, retries ...RetryStrategy) error
	// WaitForDiskOK waits for a disk to be in OK status
	WaitForDiskOK(diskID DiskID, retries ...RetryStrategy) (Disk, error)
}

// UpdateDiskParams creates a builder for the params for updating a disk.
func UpdateDiskParams() BuildableUpdateDiskParameters {
	return &updateDiskParams{}
}

// UpdateDiskParameters describes the possible parameters for updating a disk.
type UpdateDiskParameters interface {
	// Alias returns the disk alias to set. It can return nil to leave the alias unchanged.
	Alias() *string
	// ProvisionedSize returns the disk provisioned size to set.
	// It can return nil to leave the provisioned size unchanged.
	ProvisionedSize() *uint64
}

// BuildableUpdateDiskParameters is a buildable version of UpdateDiskParameters.
type BuildableUpdateDiskParameters interface {
	UpdateDiskParameters

	// WithAlias changes the params structure to set the alias to the specified value. It returns an error
	// if the alias is invalid.
	WithAlias(alias string) (BuildableUpdateDiskParameters, error)
	// MustWithAlias is identical to WithAlias, but panics instead of returning an error.
	MustWithAlias(alias string) BuildableUpdateDiskParameters

	// WithProvisionedSize changes the params structure to set the provisioned size to the specified value.
	// It returns an error if the provisioned size is invalid.
	WithProvisionedSize(size uint64) (BuildableUpdateDiskParameters, error)
	// MustWithProvisionedSize is identical to WithProvisionedSize, but panics instead of returning an error.
	MustWithProvisionedSize(size uint64) BuildableUpdateDiskParameters
}

type updateDiskParams struct {
	alias           *string
	provisionedSize *uint64
}

func (u *updateDiskParams) Alias() *string {
	return u.alias
}

func (u *updateDiskParams) WithAlias(alias string) (BuildableUpdateDiskParameters, error) {
	u.alias = &alias
	return u, nil
}

func (u *updateDiskParams) MustWithAlias(alias string) BuildableUpdateDiskParameters {
	builder, err := u.WithAlias(alias)
	if err != nil {
		panic(err)
	}
	return builder
}

func (u *updateDiskParams) ProvisionedSize() *uint64 {
	return u.provisionedSize
}

func (u *updateDiskParams) WithProvisionedSize(size uint64) (BuildableUpdateDiskParameters, error) {
	err := validateDiskSize(size)
	if err != nil {
		return u, err
	}

	u.provisionedSize = &size
	return u, nil
}

func (u *updateDiskParams) MustWithProvisionedSize(size uint64) BuildableUpdateDiskParameters {
	builder, err := u.WithProvisionedSize(size)
	if err != nil {
		panic(err)
	}
	return builder
}

// CreateDiskOptionalParameters is a structure that serves to hold the optional parameters for DiskClient.CreateDisk.
type CreateDiskOptionalParameters interface {
	// Alias is a secondary name for the disk.
	Alias() string

	// Sparse indicates that the disk should be sparse-provisioned.If it returns nil, the default will be used.
	Sparse() *bool
}

// BuildableCreateDiskParameters is a buildable version of CreateDiskOptionalParameters.
type BuildableCreateDiskParameters interface {
	CreateDiskOptionalParameters

	// WithAlias sets the alias of the disk.
	WithAlias(alias string) (BuildableCreateDiskParameters, error)
	// MustWithAlias is the same as WithAlias, but panics instead of returning an error.
	MustWithAlias(alias string) BuildableCreateDiskParameters

	// WithSparse sets sparse provisioning for the disk.
	WithSparse(sparse bool) (BuildableCreateDiskParameters, error)
	// MustWithSparse is the same as WithSparse, but panics instead of returning an error.
	MustWithSparse(sparse bool) BuildableCreateDiskParameters
}

// CreateDiskParams creates a buildable set of CreateDiskOptionalParameters for use with
// Client.CreateDisk.
func CreateDiskParams() BuildableCreateDiskParameters {
	return &createDiskParams{}
}

type createDiskParams struct {
	alias  string
	sparse *bool
}

func (c *createDiskParams) Alias() string {
	return c.alias
}

func (c *createDiskParams) WithAlias(alias string) (BuildableCreateDiskParameters, error) {
	c.alias = alias
	return c, nil
}

func (c *createDiskParams) MustWithAlias(alias string) BuildableCreateDiskParameters {
	builder, err := c.WithAlias(alias)
	if err != nil {
		panic(err)
	}
	return builder
}

func (c *createDiskParams) Sparse() *bool {
	return c.sparse
}

func (c *createDiskParams) WithSparse(sparse bool) (BuildableCreateDiskParameters, error) {
	c.sparse = &sparse
	return c, nil
}

func (c *createDiskParams) MustWithSparse(sparse bool) BuildableCreateDiskParameters {
	builder, err := c.WithSparse(sparse)
	if err != nil {
		panic(err)
	}
	return builder
}

// DiskCreation is a process object that lets you query the status of the disk creation.
type DiskCreation interface {
	// Disk returns the disk that has been created, even if it is not yet ready.
	Disk() Disk
	// Wait waits until the disk creation is complete and returns when it is done. It returns the created disk and
	// an error if one happened.
	Wait(retries ...RetryStrategy) (Disk, error)
}

// DiskUpdate is an object to monitor the progress of an update.
type DiskUpdate interface {
	// Disk returns the disk as it was during the last update call.
	Disk() Disk
	// Wait waits until the disk update is complete and returns when it is done. It returns the created disk and
	// an error if one happened.
	Wait(retries ...RetryStrategy) (Disk, error)
}

// ImageDownloadReader is a special reader for reading image downloads. On the first Read call
// it waits until the image download is ready and then returns the desired bytes. It also
// tracks how many bytes are read for an async display of a progress bar.
type ImageDownloadReader interface {
	io.Reader
	// Read reads the specified bytes from the disk image. This call will block if
	// the image is not yet ready for download.
	Read(p []byte) (n int, err error)
	// Close closes the image download and unlocks the disk.
	Close() error
	// BytesRead returns the number of bytes read so far. This can be used to
	// provide a progress bar.
	BytesRead() uint64
	// Size returns the size of the disk image in bytes. This is ONLY available after the initialization is complete and
	// MAY return 0 before.
	Size() uint64
}

// ImageDownload represents an image download in progress. The caller MUST
// close the image download when it is finished otherwise the disk will not be unlocked.
type ImageDownload interface {
	ImageDownloadReader

	// Err returns the error that happened during initializing the download, or the last error reading from the
	// image server.
	Err() error
	// Initialized returns a channel that will be closed when the initialization is complete. This can be either
	// in an errored state (check Err()) or when the image is ready.
	Initialized() <-chan struct{}
}

// UploadImageResult represents the completed image upload.
type UploadImageResult interface {
	// Disk returns the disk that has been created as the result of the image upload.
	Disk() Disk
}

// DiskData is the core of a Disk, only exposing data functions, but not the client functions.
// This can be used for cases where not a full Disk is required, but only the data functionality.
type DiskData interface {
	// ID is the unique ID for this disk.
	ID() DiskID
	// Alias is the name for this disk set by the user.
	Alias() string
	// ProvisionedSize is the size visible to the virtual machine.
	ProvisionedSize() uint64
	// TotalSize is the size of the image file.
	// This value can be zero in some cases, for example when the disk upload wasn't properly finalized.
	TotalSize() uint64
	// Format is the format of the image.
	Format() ImageFormat
	// StorageDomainIDs returns a list of storage domains this disk is present on. This will typically be a single
	// disk, but may have multiple disk when the disk has been copied over to other storage domains. The disk is always
	// present on at least one disk, so this list will never be empty.
	StorageDomainIDs() []StorageDomainID
	// Status returns the status the disk is in.
	Status() DiskStatus
	// Sparse indicates sparse provisioning on the disk.
	Sparse() bool
}

// Disk is a disk in oVirt.
type Disk interface {
	DiskData

	// StartDownload starts the download of the image file the current disk.
	// The caller can then wait for the initialization using the Initialized() call:
	//
	//     <-download.Initialized()
	//
	// Alternatively, the downloader can use the Read() function to wait for the download to become available
	// and then read immediately.
	//
	// The caller MUST close the returned reader, otherwise the disk will remain locked in the oVirt engine.
	// The passed context is observed only for the initialization phase.
	StartDownload(
		format ImageFormat,
		retries ...RetryStrategy,
	) (ImageDownload, error)

	// Download runs StartDownload, then waits for the download to be ready before returning the reader.
	// The caller MUST close the ImageDownloadReader in order to properly unlock the disk in the oVirt engine.
	Download(
		format ImageFormat,
		retries ...RetryStrategy,
	) (ImageDownloadReader, error)

	// Remove removes the current disk in the oVirt engine.
	Remove(retries ...RetryStrategy) error

	// AttachToVM attaches a disk to this VM.
	AttachToVM(
		vmID VMID,
		diskInterface DiskInterface,
		params CreateDiskAttachmentOptionalParams,
		retries ...RetryStrategy,
	) (DiskAttachment, error)

	// StartUpdate starts an update to the disk. The returned DiskUpdate can be used to wait
	// for the update to complete. Use UpdateDiskParams() to obtain a buildable structure.
	StartUpdate(
		params UpdateDiskParameters,
		retries ...RetryStrategy,
	) (DiskUpdate, error)

	// Update updates the current disk with the specified parameters.
	// Use UpdateDiskParams() to obtain a buildable structure.
	Update(
		params UpdateDiskParameters,
		retries ...RetryStrategy,
	) (Disk, error)

	// StorageDomains will fetch and return the storage domains associated with this disk.
	StorageDomains(retries ...RetryStrategy) ([]StorageDomain, error)

	// WaitForOK waits for the disk status to return to OK.
	WaitForOK(retries ...RetryStrategy) (Disk, error)
}

// DiskStatus shows the status of a disk. Certain operations lock a disk, which is important because the disk can then
// not be changed.
type DiskStatus string

const (
	// DiskStatusOK represents a disk status that operations can be performed on.
	DiskStatusOK DiskStatus = "ok"
	// DiskStatusLocked represents a disk status where no operations can be performed on the disk.
	DiskStatusLocked DiskStatus = "locked"
	// DiskStatusIllegal indicates that the disk cannot be accessed by the virtual machine, and the user needs
	// to take action to resolve the issue.
	DiskStatusIllegal DiskStatus = "illegal"
)

// DiskStatusList is a list of DiskStatus values.
type DiskStatusList []DiskStatus

// DiskStatusValues returns all possible values for DiskStatus.
func DiskStatusValues() DiskStatusList {
	return []DiskStatus{
		DiskStatusOK,
		DiskStatusLocked,
		DiskStatusIllegal,
	}
}

// Strings returns a list of strings.
func (l DiskStatusList) Strings() []string {
	result := make([]string, len(l))
	for i, status := range l {
		result[i] = string(status)
	}
	return result
}

// UploadImageProgress is a tracker for the upload progress happening in the background.
type UploadImageProgress interface {
	// Disk returns the disk created as part of the upload process once the upload is complete. Before the upload
	// is complete it will return nil.
	Disk() Disk
	// UploadedBytes returns the number of bytes already uploaded.
	//
	// Caution! This number may decrease or reset to 0 if the upload has to be retried.
	UploadedBytes() uint64
	// TotalBytes returns the total number of bytes to be uploaded.
	TotalBytes() uint64
	// Err returns the error of the upload once the upload is complete or errored.
	Err() error
	// Done returns a channel that will be closed when the upload is complete.
	Done() <-chan struct{}
}

// ImageFormat is a constant for representing the format that images can be in. This is relevant
// for both image uploads and image downloads, as the oVirt engine has the capability of converting
// between these formats.
//
// Note: the mocking facility cannot convert between the formats due to the complexity of the QCOW2
// format. It is recommended to write tests only using the raw format as comparing QCOW2 files
// is complex.
type ImageFormat string

// Validate returns an error if the image format doesn't have a valid value.
func (f ImageFormat) Validate() error {
	for _, format := range ImageFormatValues() {
		if format == f {
			return nil
		}
	}
	return newError(
		EBadArgument,
		"invalid image format: %s must be one of: %s",
		f,
		strings.Join(ImageFormatValues().Strings(), ", "),
	)
}

const (
	// ImageFormatCow is an image conforming to the QCOW2 image format. This image format can use
	// compression, supports snapshots, and other features.
	// See https://github.com/qemu/qemu/blob/master/docs/interop/qcow2.txt for details.
	ImageFormatCow ImageFormat = "cow"
	// ImageFormatRaw is not actually a format, it only contains the raw bytes on the block device.
	ImageFormatRaw ImageFormat = "raw"
)

// ImageFormatList is a list of ImageFormat values.
type ImageFormatList []ImageFormat

// ImageFormatValues returns all possible ImageFormat values.
func ImageFormatValues() ImageFormatList {
	return []ImageFormat{
		ImageFormatCow,
		ImageFormatRaw,
	}
}

// Strings creates a string list of the values.
func (l ImageFormatList) Strings() []string {
	result := make([]string, len(l))
	for i, status := range l {
		result[i] = string(status)
	}
	return result
}

func convertSDKDisk(sdkDisk *ovirtsdk4.Disk, client Client) (Disk, error) {
	id, ok := sdkDisk.Id()
	if !ok {
		return nil, newError(EFieldMissing, "disk does not contain an ID")
	}
	var storageDomainIDs []StorageDomainID
	if sdkStorageDomain, ok := sdkDisk.StorageDomain(); ok {
		storageDomainID, _ := sdkStorageDomain.Id()
		storageDomainIDs = append(storageDomainIDs, StorageDomainID(storageDomainID))
	}
	if sdkStorageDomains, ok := sdkDisk.StorageDomains(); ok {
		for _, sd := range sdkStorageDomains.Slice() {
			storageDomainID, _ := sd.Id()
			storageDomainIDs = append(storageDomainIDs, StorageDomainID(storageDomainID))
		}
	}
	if len(storageDomainIDs) == 0 {
		return nil, newError(EFieldMissing, "failed to find a valid storage domain for disk %s", id)
	}
	alias, ok := sdkDisk.Alias()
	if !ok {
		return nil, newError(EFieldMissing, "disk %s does not contain an alias", id)
	}
	provisionedSize, ok := sdkDisk.ProvisionedSize()
	if !ok {
		return nil, newError(EFieldMissing, "disk %s does not contain a provisioned size", id)
	}
	totalSize, ok := sdkDisk.TotalSize()
	if !ok {
		return nil, newError(EFieldMissing, "disk %s does not contain a total size", id)
	}
	format, ok := sdkDisk.Format()
	if !ok {
		return nil, newError(EFieldMissing, "disk %s has no format field", id)
	}
	status, ok := sdkDisk.Status()
	if !ok {
		return nil, newError(EFieldMissing, "disk %s has no status field", id)
	}
	sparse, ok := sdkDisk.Sparse()
	if !ok {
		return nil, newError(EFieldMissing, "disk %s has no sparse field", id)
	}
	return &disk{
		client: client,

		id:               DiskID(id),
		alias:            alias,
		provisionedSize:  uint64(provisionedSize),
		totalSize:        uint64(totalSize),
		format:           ImageFormat(format),
		storageDomainIDs: storageDomainIDs,
		status:           DiskStatus(status),
		sparse:           sparse,
	}, nil
}

type disk struct {
	client Client

	id               DiskID
	alias            string
	provisionedSize  uint64
	format           ImageFormat
	storageDomainIDs []StorageDomainID
	status           DiskStatus
	totalSize        uint64
	sparse           bool
}

func (d *disk) WaitForOK(retries ...RetryStrategy) (Disk, error) {
	return d.client.WaitForDiskOK(d.id, retries...)
}

func (d *disk) StorageDomainIDs() []StorageDomainID {
	return d.storageDomainIDs
}

func (d *disk) StorageDomains(retries ...RetryStrategy) ([]StorageDomain, error) {
	storageDomains := make([]StorageDomain, len(d.storageDomainIDs))
	for i, id := range d.storageDomainIDs {
		storageDomain, err := d.client.GetStorageDomain(id, retries...)
		if err != nil {
			return nil, err
		}
		storageDomains[i] = storageDomain
	}
	return storageDomains, nil
}

func (d *disk) Update(params UpdateDiskParameters, retries ...RetryStrategy) (Disk, error) {
	return d.client.UpdateDisk(d.id, params, retries...)
}

func (d *disk) StartUpdate(params UpdateDiskParameters, retries ...RetryStrategy) (DiskUpdate, error) {
	return d.client.StartUpdateDisk(d.id, params, retries...)
}

func (d *disk) Sparse() bool {
	return d.sparse
}

func (d *disk) AttachToVM(
	vmID VMID,
	diskInterface DiskInterface,
	params CreateDiskAttachmentOptionalParams,
	retries ...RetryStrategy,
) (DiskAttachment, error) {
	return d.client.CreateDiskAttachment(vmID, d.id, diskInterface, params, retries...)
}

func (d *disk) Remove(retries ...RetryStrategy) error {
	return d.client.RemoveDisk(d.id, retries...)
}

func (d *disk) TotalSize() uint64 {
	return d.totalSize
}

func (d disk) Status() DiskStatus {
	return d.status
}

func (d disk) ID() DiskID {
	return d.id
}

func (d disk) Alias() string {
	return d.alias
}

func (d disk) ProvisionedSize() uint64 {
	return d.provisionedSize
}

func (d disk) Format() ImageFormat {
	return d.format
}

func (d disk) StartDownload(format ImageFormat, retries ...RetryStrategy) (ImageDownload, error) {
	return d.client.StartDownloadDisk(d.id, format, retries...)
}

func (d disk) Download(format ImageFormat, retries ...RetryStrategy) (ImageDownloadReader, error) {
	return d.client.DownloadDisk(d.id, format, retries...)
}

type diskWait struct {
	client        *oVirtClient
	disk          Disk
	correlationID string
	lock          *sync.Mutex
}

func (d *diskWait) Disk() Disk {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.disk
}

func (d *diskWait) Wait(retries ...RetryStrategy) (Disk, error) {
	retries = defaultRetries(retries, defaultLongTimeouts(d.client))
	if err := d.client.waitForJobFinished(d.correlationID, retries); err != nil {
		return d.disk, err
	}

	d.lock.Lock()
	diskID := d.disk.ID()
	d.lock.Unlock()

	disk, err := d.client.GetDisk(diskID, retries...)

	d.lock.Lock()
	defer d.lock.Unlock()
	if disk != nil {
		d.disk = disk
	}
	return disk, err
}
