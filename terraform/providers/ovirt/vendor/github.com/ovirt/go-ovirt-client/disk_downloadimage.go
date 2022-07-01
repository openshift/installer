package ovirtclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// Deprecated: use StartDownloadDisk instead.
func (o *oVirtClient) StartImageDownload(diskID DiskID, format ImageFormat, retries ...RetryStrategy) (ImageDownload, error) {
	o.logger.Warningf("Using StartImageDownload is deprecated, please use StartDownloadDisk instead.")
	return o.StartDownloadDisk(diskID, format, retries...)
}

func (o *oVirtClient) StartDownloadDisk(diskID DiskID, format ImageFormat, retries ...RetryStrategy) (ImageDownload, error) {
	retries = defaultRetries(retries, defaultLongTimeouts(o))

	o.logger.Infof("Starting disk %s image download...", diskID)
	disk, err := o.GetDisk(diskID)
	if err != nil {
		return nil, wrap(err, EUnidentified, "failed to fetch disk for image download")
	}

	realCtx, cancel := context.WithCancel(context.Background())

	dl := &imageDownload{
		disk:       disk,
		lock:       &sync.Mutex{},
		bytesRead:  0,
		size:       disk.TotalSize(),
		lastError:  nil,
		ctx:        realCtx,
		cancel:     cancel,
		conn:       o.conn,
		done:       make(chan struct{}),
		reader:     nil,
		httpClient: o.httpClient,
		createReq:  nil,
		transfer:   nil,
		cli:        o,
		logger:     o.logger,
		retries:    retries,
		format:     format,
	}
	go dl.poll()
	return dl, nil
}

// Deprecated: use DownloadDisk instead.
func (o *oVirtClient) DownloadImage(diskID DiskID, format ImageFormat, retries ...RetryStrategy) (
	ImageDownloadReader,
	error,
) {
	o.logger.Warningf("Using DownloadImage is deprecated, please use DownloadDisk instead.")
	return o.DownloadDisk(diskID, format, retries...)
}

func (o *oVirtClient) DownloadDisk(
	diskID DiskID,
	format ImageFormat,
	retries ...RetryStrategy,
) (ImageDownloadReader, error) {
	download, err := o.StartDownloadDisk(diskID, format, retries...)
	if err != nil {
		return nil, err
	}
	<-download.Initialized()
	if err := download.Err(); err != nil {
		_ = download.Close()
		return nil, err
	}
	return download, nil
}

type imageDownload struct {
	lock *sync.Mutex

	bytesRead uint64
	size      uint64

	lastError error
	ctx       context.Context
	cancel    context.CancelFunc
	conn      *ovirtsdk4.Connection
	done      chan struct{}

	reader     io.ReadCloser
	httpClient http.Client
	createReq  *ovirtsdk4.ImageTransfersServiceAddRequest
	transfer   imageTransfer
	cli        *oVirtClient
	logger     Logger
	retries    []RetryStrategy
	format     ImageFormat
	disk       Disk
}

// poll polls the oVirt API for the status of the transfer and initializes the HTTP request to
// download the image. Once the HTTP transfer is available control is handed back and the
// imageDownload becomes a reader to read the data from. When Close is called on imageDownload
// the transfer is finalized.
func (i *imageDownload) poll() {
	defer close(i.done)
	i.transfer = newImageTransfer(
		i.cli,
		i.logger,
		i.disk.ID(),
		"",
		i.retries,
		ovirtsdk4.IMAGETRANSFERDIRECTION_DOWNLOAD,
		ovirtsdk4.DiskFormat(i.format),
		i.updateDisk,
	)
	transferURL := ""
	var err error
	if transferURL, err = i.transfer.initialize(); err != nil {
		i.lastError = i.transfer.finalize(err)
		return
	}
	var httpResponse *http.Response
	httpResponse, err = i.transferImage(transferURL) //nolint:bodyclose
	if err != nil {
		i.lastError = i.transfer.finalize(err)
		return
	}
	i.reader = httpResponse.Body
}

// updateDisk is a helper function that updates the internally-stored disk object when the transfer updates it.
func (i *imageDownload) updateDisk(disk Disk) {
	i.disk = disk
}

// transferImage will retry a HTTP GET request to download the image and return the HTTP response if successful.
// This call will also set the exact download size in i.size. This function will retry until a valid URL is obtained
// or retries are exhausted.
func (i *imageDownload) transferImage(transferURL string) (httpResponse *http.Response, err error) {
	return httpResponse, retry(
		fmt.Sprintf("transferring image from %s", transferURL),
		i.logger,
		i.retries,
		func() error {
			response, err := i.attemptTransferImage(transferURL) //nolint:bodyclose
			httpResponse = response
			return err
		},
	)
}

// attemptTransferImage will create a single attempt to download an image from the specified transfer URL.
func (i *imageDownload) attemptTransferImage(transferURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", transferURL, nil)
	if err != nil {
		return nil, wrap(err, EBug, "failed to create HTTP request to %s", transferURL)
	}
	httpResponse, err := i.httpClient.Do(req)
	if err != nil {
		return httpResponse, wrap(
			err,
			EConnection,
			"HTTP request to image transfer URL %s failed",
			transferURL,
		)
	}
	if err := i.transfer.checkStatusCode(httpResponse.StatusCode); err != nil {
		_ = httpResponse.Body.Close()
		return nil, wrap(
			err,
			EUnidentified,
			"failed to download image from disk %s",
			i.disk.ID(),
		)
	}
	if err := i.extractDownloadSize(httpResponse); err != nil {
		_ = httpResponse.Body.Close()
		return nil, wrap(
			err,
			EUnidentified,
			"failed to determine download size for disk %s",
			i.disk.ID(),
		)
	}
	return httpResponse, nil
}

// extractDownloadSize extracts the size of the image download from the content-length header of the
// HTTP response if needed. This is a fallback mechanism for the case when the engine doesn't return
// the correct image size.
func (i *imageDownload) extractDownloadSize(httpResponse *http.Response) error {
	if i.size != 0 {
		return nil
	}
	header := httpResponse.Header
	if header.Get("content-encoding") != "" {
		return newError(
			EBug,
			"the oVirt engine API did not return an image download size and the ImageIO response contained a non-plaintext encoding, so the download size cannot be determined",
		)
	}
	if contentLengthString := header.Get("content-length"); contentLengthString != "" {
		contentLength, err := strconv.ParseUint(contentLengthString, 10, 64)
		if err != nil {
			return wrap(
				err,
				EBug,
				"the oVirt engine API did not return an image download size and the ImageIO response contained an invalid content-length header (%s), so the download size cannot be determined",
			)
		}
		i.size = contentLength
		return nil
	}
	return newError(
		EBug,
		"the oVirt engine API did not return an image download size and the ImageIO response did not contain a content-length header, so the download size cannot be determined",
	)
}

// Err will return the last error from the image download.
func (i *imageDownload) Err() error {
	return i.lastError
}

// Initialized will return a channel that is closed when the image transfer is initialized.
func (i *imageDownload) Initialized() <-chan struct{} {
	return i.done
}

// Read waits for the transfer to be properly initialized and in the transferring state, then
// reads from the HTTP response body. When there are no more bytes left it attempts to automatically
// finalize the transfer by calling the Close function.
func (i *imageDownload) Read(p []byte) (n int, err error) {
	<-i.done
	if i.lastError != nil {
		return 0, i.lastError
	}
	n, err = i.reader.Read(p)
	i.lock.Lock()
	defer i.lock.Unlock()
	i.bytesRead += uint64(n)

	if i.bytesRead == i.size {
		go func() {
			_ = i.Close()
		}()
	}

	return n, err
}

// Close is an implementation of the required function in io.ReadCloser and is responsible for
// closing the HTTP reader body and finalizing the image download. This is important so the disk
// does not stay locked.
func (i *imageDownload) Close() error {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.cancel()
	if i.reader != nil {
		err := i.reader.Close()
		if err == nil {
			i.reader = nil
		}
		return i.transfer.finalize(err)
	}
	return nil
}

// BytesRead returns the number of bytes already read from the download reader.
func (i *imageDownload) BytesRead() uint64 {
	return i.bytesRead
}

// Size returns the size of the download. This does not necessarily match the disk size since the oVirt Engine may not
// return the disk size correctly.
func (i *imageDownload) Size() uint64 {
	return i.size
}

// Deprecated: use StartDownloadDisk instead.
func (m *mockClient) StartImageDownload(diskID DiskID, format ImageFormat, retries ...RetryStrategy) (
	ImageDownload,
	error,
) {
	return m.StartDownloadDisk(diskID, format, retries...)
}

func (m *mockClient) StartDownloadDisk(diskID DiskID, format ImageFormat, _ ...RetryStrategy) (ImageDownload, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	disk, ok := m.disks[diskID]
	if !ok {
		return nil, newError(ENotFound, "disk with ID %s not found", diskID)
	}

	if disk.format != format {
		m.logger.Warningf("the image upload client requested a conversion from from %s to %s; the mock library does not support this and the source image data will be used unmodified which may lead to errors", disk.format, format)
	}

	dl := &mockImageDownload{
		disk:      disk,
		size:      0,
		bytesRead: 0,
		done:      make(chan struct{}),
		lastError: nil,
		lock:      &sync.Mutex{},
		reader:    bytes.NewReader(disk.data),
	}
	go dl.prepare()

	return dl, nil
}

// Deprecated: use DownloadDisk instead.
func (m *mockClient) DownloadImage(diskID DiskID, format ImageFormat, retries ...RetryStrategy) (
	ImageDownloadReader,
	error,
) {
	return m.DownloadDisk(diskID, format, retries...)
}

func (m *mockClient) DownloadDisk(diskID DiskID, format ImageFormat, retries ...RetryStrategy) (
	ImageDownloadReader,
	error,
) {
	download, err := m.StartDownloadDisk(diskID, format, retries...)
	if err != nil {
		return nil, err
	}
	<-download.Initialized()
	if err := download.Err(); err != nil {
		return nil, err
	}
	return download, nil
}

type mockImageDownload struct {
	disk      *diskWithData
	size      uint64
	bytesRead uint64
	done      chan struct{}
	lastError error
	lock      *sync.Mutex
	reader    io.Reader
}

func (m *mockImageDownload) Err() error {
	return m.lastError
}

func (m *mockImageDownload) Initialized() <-chan struct{} {
	return m.done
}

func (m *mockImageDownload) Read(p []byte) (n int, err error) {
	<-m.done
	if m.lastError != nil {
		return 0, m.lastError
	}

	n, err = m.reader.Read(p)

	m.lock.Lock()
	defer m.lock.Unlock()
	if err != nil {
		m.lastError = err
	}
	m.bytesRead += uint64(n)

	if m.bytesRead == m.size {
		go func() {
			_ = m.Close()
		}()
	}

	return n, err
}

func (m *mockImageDownload) Close() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.disk.Unlock()
	return nil
}

func (m *mockImageDownload) BytesRead() uint64 {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.bytesRead
}

func (m *mockImageDownload) Size() uint64 {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.size
}

func (m *mockImageDownload) prepare() {
	// Sleep one second to trigger possible race condition with determining size.
	time.Sleep(time.Second)
	m.lock.Lock()
	defer m.lock.Unlock()
	m.size = uint64(len(m.disk.data))
	close(m.done)
}
