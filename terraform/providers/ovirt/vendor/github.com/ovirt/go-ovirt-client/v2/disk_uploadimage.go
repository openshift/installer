//
// This file implements the image upload-related functions of the oVirt client.
//

package ovirtclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// Deprecated: use UploadToNewDisk instead.
func (o *oVirtClient) UploadImage(
	alias string,
	storageDomainID StorageDomainID,
	sparse bool,
	size uint64,
	reader io.ReadSeekCloser,
	retries ...RetryStrategy,
) (UploadImageResult, error) {
	o.logger.Warningf("Using UploadImage is deprecated. Please use UploadToNewDisk instead.")
	return o.UploadToNewDisk(
		storageDomainID,
		"",
		size,
		CreateDiskParams().MustWithSparse(sparse).MustWithAlias(alias),
		reader,
		retries...,
	)
}

func (o *oVirtClient) UploadToNewDisk(
	storageDomainID StorageDomainID,
	format ImageFormat,
	size uint64,
	params CreateDiskOptionalParameters,
	reader io.ReadSeekCloser,
	retries ...RetryStrategy,
) (UploadImageResult, error) {
	retries = defaultRetries(retries, defaultLongTimeouts(o))
	progress, err := o.StartUploadToNewDisk(storageDomainID, format, size, params, reader, retries...)
	if err != nil {
		return nil, err
	}
	<-progress.Done()
	if err := progress.Err(); err != nil {
		return nil, err
	}
	return progress, nil
}

// Deprecated: use StartUploadToNewDisk instead.
func (o *oVirtClient) StartImageUpload(
	alias string,
	storageDomainID StorageDomainID,
	sparse bool,
	size uint64,
	reader io.ReadSeekCloser,
	retries ...RetryStrategy,
) (UploadImageProgress, error) {
	o.logger.Warningf("Using StartImageUpload is deprecated. Please use StartUploadToNewDisk instead.")
	return o.StartUploadToNewDisk(
		storageDomainID,
		"",
		size,
		CreateDiskParams().MustWithSparse(sparse).MustWithAlias(alias),
		reader,
		retries...,
	)
}

func (o *oVirtClient) UploadToDisk(
	diskID DiskID,
	size uint64,
	reader io.ReadSeekCloser,
	retries ...RetryStrategy,
) error {
	retries = defaultRetries(retries, defaultLongTimeouts(o))
	progress, err := o.StartUploadToDisk(diskID, size, reader, retries...)
	if err != nil {
		return err
	}
	<-progress.Done()
	return progress.Err()
}

func (o *oVirtClient) StartUploadToDisk(
	diskID DiskID,
	size uint64,
	reader io.ReadSeekCloser,
	retries ...RetryStrategy,
) (UploadImageProgress, error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	o.logger.Infof("Starting disk image upload...")
	disk, err := o.GetDisk(diskID, retries...)
	if err != nil {
		return nil, err
	}

	format, qcowSize, err := extractQCOWParameters(size, reader)
	if err != nil {
		return nil, err
	}

	if qcowSize > disk.ProvisionedSize() {
		return nil, newError(
			EBadArgument,
			"the virtual image size (%d bytes) is larger than the target disk %s (%d bytes)",
			qcowSize,
			diskID,
			disk.ProvisionedSize(),
		)
	}
	ctx, cancel := context.WithCancel(context.Background())
	progress := &uploadToDiskProgress{
		client:        o,
		lock:          &sync.Mutex{},
		done:          make(chan struct{}),
		ctx:           ctx,
		cancel:        cancel,
		correlationID: fmt.Sprintf("image_upload_%s", generateRandomID(5, o.nonSecureRandom)),
		format:        format,
		disk:          disk,
		totalBytes:    size,
		qcowSize:      qcowSize,
		reader:        reader,
		retries:       retries,
	}
	go progress.Do()
	return progress, nil
}

type uploadToDiskProgress struct {
	client           *oVirtClient
	lock             *sync.Mutex
	done             chan struct{}
	ctx              context.Context
	cancel           func()
	disk             Disk
	correlationID    string
	reader           io.ReadSeekCloser
	retries          []RetryStrategy
	transferredBytes uint64
	totalBytes       uint64
	err              error
	format           ImageFormat
	qcowSize         uint64
}

func (u *uploadToDiskProgress) Close() error {
	return u.reader.Close()
}

func (u *uploadToDiskProgress) updateDisk(disk Disk) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.disk = disk
}

func (u *uploadToDiskProgress) Do() {
	defer func() {
		close(u.done)
		u.cancel()
	}()

	err := u.transfer()

	u.lock.Lock()
	u.err = err
	u.lock.Unlock()
}

func (u *uploadToDiskProgress) transfer() error {
	transfer := newImageTransfer(
		u.client,
		u.client.logger,
		u.disk.ID(),
		u.correlationID,
		u.retries,
		ovirtsdk4.IMAGETRANSFERDIRECTION_UPLOAD,
		ovirtsdk4.DiskFormat(u.format),
		u.updateDisk,
	)
	transferURL := ""
	var err error
	if transferURL, err = transfer.initialize(); err != nil {
		return transfer.finalize(err)
	}
	err = u.transferImage(transfer, transferURL)
	return transfer.finalize(err)
}

// transferImage does an HTTP request to transfer the image to the specified transfer URL.
func (u *uploadToDiskProgress) transferImage(transfer imageTransfer, transferURL string) error {
	return retry(
		fmt.Sprintf(
			"transferring image for disk %s via HTTP request to %s",
			u.disk.ID(),
			transferURL,
		),
		u.client.logger,
		u.retries,
		func() error {
			return u.putRequest(transferURL, transfer)
		},
	)
}

// putRequest performs a single HTTP put request to upload an image. This can be called multiple times to retry an
// upload.
func (u *uploadToDiskProgress) putRequest(transferURL string, transfer imageTransfer) error {
	// We ensure that the reader is at the first byte before attempting a PUT request, otherwise we may upload an
	// incomplete image.
	if _, err := u.reader.Seek(0, io.SeekStart); err != nil {
		return wrap(
			err,
			ELocalIO,
			"could not seek to the first byte of the disk image before upload",
		)
	}

	u.lock.Lock()
	u.transferredBytes = 0
	u.lock.Unlock()

	putRequest, err := http.NewRequest(http.MethodPut, transferURL, u)
	if err != nil {
		return wrap(err, EUnidentified, "failed to create HTTP request")
	}
	putRequest.Header.Add("content-type", "application/octet-stream")
	putRequest.ContentLength = int64(u.totalBytes)
	putRequest.Body = u
	response, err := u.client.httpClient.Do(putRequest)
	if err != nil {
		return wrap(
			err,
			EUnidentified,
			"failed to upload image",
		)
	}
	if err := transfer.checkStatusCode(response.StatusCode); err != nil {
		_ = response.Body.Close()
		return err
	}
	if err := response.Body.Close(); err != nil {
		return wrap(
			err,
			EUnidentified,
			"failed to close response body while uploading image",
		)
	}
	return nil
}

func (u *uploadToDiskProgress) Disk() Disk {
	u.lock.Lock()
	defer u.lock.Unlock()
	return u.disk
}

func (u *uploadToDiskProgress) UploadedBytes() uint64 {
	u.lock.Lock()
	defer u.lock.Unlock()
	return u.transferredBytes
}

func (u *uploadToDiskProgress) TotalBytes() uint64 {
	return u.totalBytes
}

func (u *uploadToDiskProgress) Err() error {
	u.lock.Lock()
	defer u.lock.Unlock()
	return u.err
}

func (u *uploadToDiskProgress) Done() <-chan struct{} {
	return u.done
}

func (u *uploadToDiskProgress) Read(p []byte) (n int, err error) {
	select {
	case <-u.ctx.Done():
		return 0, newError(ETimeout, "timeout while uploading image")
	default:
	}
	n, err = u.reader.Read(p)
	u.transferredBytes += uint64(n)
	return
}

func (o *oVirtClient) StartUploadToNewDisk(
	storageDomainID StorageDomainID,
	format ImageFormat,
	size uint64,
	params CreateDiskOptionalParameters,
	reader io.ReadSeekCloser,
	retries ...RetryStrategy,
) (UploadImageProgress, error) {
	retries = defaultRetries(retries, defaultLongTimeouts(o))

	o.logger.Infof("Starting disk image upload...")

	imageFormat, qcowSize, err := extractQCOWParameters(size, reader)
	if err != nil {
		return nil, err
	}

	if format == "" {
		format = imageFormat
	} else if err := format.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	progress := &uploadToNewDiskProgress{
		uploadToDiskProgress: uploadToDiskProgress{
			client:        o,
			lock:          &sync.Mutex{},
			done:          make(chan struct{}),
			ctx:           ctx,
			cancel:        cancel,
			correlationID: fmt.Sprintf("image_upload_%s", generateRandomID(5, o.nonSecureRandom)),
			format:        imageFormat,
			disk:          nil,
			totalBytes:    size,
			qcowSize:      qcowSize,
			reader:        reader,
			retries:       retries,
		},

		storageDomainID: storageDomainID,
		diskFormat:      format,
		diskParams:      params,
	}

	go progress.Do()
	return progress, nil
}

type uploadToNewDiskProgress struct {
	uploadToDiskProgress

	storageDomainID StorageDomainID
	diskFormat      ImageFormat
	diskParams      CreateDiskOptionalParameters
}

func (u *uploadToNewDiskProgress) Do() {
	defer func() {
		close(u.done)
		u.cancel()
	}()

	size := u.qcowSize
	if size < 1024*1024 {
		size = 1024 * 1024
	}

	disk, err := u.client.CreateDisk(
		u.storageDomainID,
		u.diskFormat,
		size,
		u.diskParams,
		u.retries...,
	)
	if err != nil {
		u.lock.Lock()
		u.err = err
		u.lock.Unlock()
		return
	}

	u.updateDisk(disk)

	err = u.uploadToDiskProgress.transfer()
	u.lock.Lock()
	u.err = err
	u.lock.Unlock()

	if err != nil {
		u.client.logger.Infof("Image upload to new disk failed, removing created disk (%v)", err)
		if err := disk.Remove(u.retries...); err != nil && !HasErrorCode(err, ENotFound) {
			u.client.logger.Warningf(
				"Failed to remove newly created disk %s after failed image upload, please remove manually. (%v)",
				disk.ID(),
				err,
			)
		}
	}
}

func (m *mockClient) StartImageUpload(
	alias string,
	storageDomainID StorageDomainID,
	sparse bool,
	size uint64,
	reader io.ReadSeekCloser,
	retries ...RetryStrategy,
) (UploadImageProgress, error) {
	return m.StartUploadToNewDisk(
		storageDomainID,
		"",
		size,
		CreateDiskParams().MustWithSparse(sparse).MustWithAlias(alias),
		reader,
		retries...,
	)
}

func (m *mockClient) UploadImage(
	alias string,
	storageDomainID StorageDomainID,
	sparse bool,
	size uint64,
	reader io.ReadSeekCloser,
	retries ...RetryStrategy,
) (UploadImageResult, error) {
	return m.UploadToNewDisk(
		storageDomainID,
		"",
		size,
		CreateDiskParams().MustWithSparse(sparse).MustWithAlias(alias),
		reader,
		retries...,
	)
}

func (m *mockClient) StartUploadToDisk(
	diskID DiskID,
	size uint64,
	reader io.ReadSeekCloser,
	retries ...RetryStrategy,
) (UploadImageProgress, error) {
	disk, err := m.getDisk(diskID, retries...)
	if err != nil {
		return nil, err
	}

	imageFormat, qcowSize, err := extractQCOWParameters(size, reader)
	if err != nil {
		return nil, err
	}

	if qcowSize > disk.TotalSize() {
		return nil, newError(
			EBadArgument,
			"the specified size (%d bytes) is larger than the target disk %s (%d bytes)",
			size,
			diskID,
			disk.TotalSize(),
		)
	}

	if imageFormat != disk.Format() {
		return nil, newError(
			EBadArgument,
			"the mock facility doesn't support uploading %s images to %s disks,"+
				" please upload in the disk format in your tests.",
			imageFormat,
			disk.Format(),
		)
	}

	progress := &mockImageUploadProgress{
		err:    nil,
		disk:   disk,
		client: m,
		reader: reader,
		size:   size,
		done:   make(chan struct{}),
	}

	// Lock the disk to simulate the upload being initialized.
	if err := progress.disk.Lock(); err != nil {
		return nil, newError(EDiskLocked, "disk locked after creation")
	}

	go progress.do()

	return progress, nil
}

func (m *mockClient) UploadToDisk(diskID DiskID, size uint64, reader io.ReadSeekCloser, retries ...RetryStrategy) error {
	progress, err := m.StartUploadToDisk(diskID, size, reader, retries...)
	if err != nil {
		return err
	}
	<-progress.Done()
	return progress.Err()
}

func (m *mockClient) StartUploadToNewDisk(
	storageDomainID StorageDomainID,
	format ImageFormat,
	size uint64,
	params CreateDiskOptionalParameters,
	reader io.ReadSeekCloser,
	_ ...RetryStrategy,
) (UploadImageProgress, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.storageDomains[storageDomainID]; !ok {
		return nil, newError(ENotFound, "storage domain with ID %s not found", storageDomainID)
	}

	imageFormat, qcowSize, err := extractQCOWParameters(size, reader)
	if err != nil {
		return nil, err
	}

	if imageFormat != format {
		return nil, newError(
			EBadArgument,
			"the mock facility doesn't support uploading %s images to %s disks,"+
				" please upload in the disk format in your tests.",
			imageFormat,
			format,
		)
	}

	if qcowSize < 1024*1024 {
		qcowSize = 1024 * 1024
	}

	disk, err := m.createDisk(storageDomainID, format, qcowSize, params)
	if err != nil {
		return nil, err
	}
	// Unlock the disk to simulate disk creation being complete.
	disk.Unlock()

	progress := &mockImageUploadProgress{
		err:    nil,
		disk:   disk,
		client: m,
		reader: reader,
		size:   size,
		done:   make(chan struct{}),
	}

	// Lock the disk to simulate the upload being initialized.
	if err := progress.disk.Lock(); err != nil {
		return nil, newError(EDiskLocked, "disk locked after creation")
	}

	go progress.do()

	return progress, nil
}

func (m *mockClient) UploadToNewDisk(
	storageDomainID StorageDomainID,
	format ImageFormat,
	size uint64,
	params CreateDiskOptionalParameters,
	reader io.ReadSeekCloser,
	retries ...RetryStrategy,
) (UploadImageResult, error) {
	progress, err := m.StartUploadToNewDisk(storageDomainID, format, size, params, reader, retries...)
	if err != nil {
		return nil, err
	}
	<-progress.Done()
	if err := progress.Err(); err != nil {
		return nil, err
	}
	return progress, nil
}

type mockImageUploadProgress struct {
	err           error
	disk          *diskWithData
	client        *mockClient
	reader        io.ReadSeekCloser
	size          uint64
	uploadedBytes uint64
	done          chan struct{}
}

func (m *mockImageUploadProgress) Disk() Disk {
	disk := m.disk
	if disk.id == "" {
		return nil
	}
	return disk
}

func (m *mockImageUploadProgress) UploadedBytes() uint64 {
	return m.uploadedBytes
}

func (m *mockImageUploadProgress) TotalBytes() uint64 {
	return m.size
}

func (m *mockImageUploadProgress) Err() error {
	return m.err
}

func (m *mockImageUploadProgress) Done() <-chan struct{} {
	return m.done
}

func (m *mockImageUploadProgress) do() {
	defer func() {
		m.disk.Unlock()
		close(m.done)
	}()

	var err error
	if _, err = m.reader.Seek(0, io.SeekStart); err != nil {
		m.err = fmt.Errorf("failed to seek to start of image file (%w)", err)
		return
	}
	m.disk.data, err = io.ReadAll(m.reader)
	m.err = err
	if err != nil {
		m.uploadedBytes = m.size
	}
}
