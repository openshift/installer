// Copyright (C) 2019 oVirt Maintainers
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

const BufferSize = 50 * 1048576 // 50MiB

func resourceOvirtImageTransfer() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtImageTransferCreate,
		Read:   resourceOvirtImageTransferRead,
		Delete: resourceOvirtImageTransferDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// the name of the uploaded image
			"alias": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"sparse": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"disk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOvirtImageTransferCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	alias := d.Get("alias").(string)
	sourceUrl := d.Get("source_url").(string)
	domainId := d.Get("storage_domain_id").(string)
	sparse := d.Get("sparse").(bool)
	correlationID := fmt.Sprintf("image_transfer_%s", alias)

	uploadSize, qcowSize, sourceFile, format, err := PrepareForTransfer(sourceUrl)
	if err != nil {
		return fmt.Errorf("Failed preparing disk for transfer: %s", err)
	}

	diskBuilder := ovirtsdk4.NewDiskBuilder().
		Alias(alias).
		Format(format).
		ProvisionedSize(int64(qcowSize)).
		InitialSize(int64(qcowSize)).
		StorageDomainsOfAny(
			ovirtsdk4.NewStorageDomainBuilder().
				Id(domainId).
				MustBuild())
	diskBuilder.Sparse(sparse)
	disk, err := diskBuilder.Build()
	if err != nil {
		return err
	}

	addDiskRequest := conn.SystemService().DisksService().Add().Disk(disk)
	addDiskRequest.Query("correlation_id", correlationID)
	addResp, err := addDiskRequest.Send()
	if err != nil {
		alias, _ := disk.Alias()
		log.Printf("[DEBUG] Error creating the Disk (%s)", alias)
		return err
	}
	diskID := addResp.MustDisk().MustId()

	// Wait for disk is ready
	log.Printf("[DEBUG] Disk (%s) is created and wait for ready (status is OK)", diskID)

	diskService := conn.SystemService().DisksService().DiskService(diskID)

	for {
		req, _ := diskService.Get().Send()
		if req.MustDisk().MustStatus() == ovirtsdk4.DISKSTATUS_OK {
			break
		}
		log.Print("waiting for disk to be OK")
		time.Sleep(time.Second * 5)
	}

	log.Printf("starting a transfer for disk id: %s", diskID)

	// initialize an image transfer
	if _, err := UploadToDisk(conn, sourceFile, diskID, alias, uploadSize, correlationID); err != nil {
		return err
	}

	jobFinishedConf := &resource.StateChangeConf{
		// An empty list indicates all jobs are completed
		Target:     []string{},
		Refresh:    jobRefreshFunc(conn, correlationID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 15 * time.Second,
	}
	_, err = jobFinishedConf.WaitForState()

	for {
		req, _ := diskService.Get().Send()

		// the system may remove the disk if it find it not compatible
		disk, ok := req.Disk()
		if !ok {
			return fmt.Errorf("the disk was removed, the upload is probably illegal")
		}
		if disk.MustStatus() == ovirtsdk4.DISKSTATUS_OK {
			d.SetId(disk.MustId())
			d.Set("disk_id", disk.MustId())
			break
		}
		log.Printf("waiting for disk to be OK")
		time.Sleep(time.Second * 5)
	}

	return resourceOvirtDiskRead(d, meta)
}

// UploadToDisk reads a file and uploads its content to the ovirt disk.
// Return value is an image transfer object, representing the transfer progress
func UploadToDisk(conn *ovirtsdk4.Connection, sourceFile *os.File, diskID string, alias string, uploadSize int64, correlationID string) (*ovirtsdk4.ImageTransferService, error) {
	defer sourceFile.Close()
	imageTransfersService := conn.SystemService().ImageTransfersService()
	image := ovirtsdk4.NewImageBuilder().Id(diskID).MustBuild()
	log.Printf("the image to transfer: %s", alias)
	transfer := ovirtsdk4.NewImageTransferBuilder().Image(image).MustBuild()
	transferReq := imageTransfersService.Add().ImageTransfer(transfer)
	transferReq.Query("correlation_id", correlationID)
	transferRes, err := transferReq.Send()
	if err != nil {
		log.Printf("failed to initialize an image transfer for image (%v) : %s", transfer, err)
		return nil, err
	}
	transfer = transferRes.MustImageTransfer()
	transferService := imageTransfersService.ImageTransferService(transfer.MustId())
	for {
		req, _ := transferService.Get().Send()
		if req.MustImageTransfer().MustPhase() == ovirtsdk4.IMAGETRANSFERPHASE_TRANSFERRING {
			break
		}
		fmt.Println("waiting for transfer phase to reach transferring")
		time.Sleep(time.Second * 5)
	}
	uploadUrl, err := detectUploadUrl(transfer)
	if err != nil {
		return nil, err
	}
	client := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	urlReader := bufio.NewReaderSize(sourceFile, BufferSize)
	putRequest, err := http.NewRequest(http.MethodPut, uploadUrl, urlReader)
	putRequest.Header.Add("content-type", "application/octet-stream")
	putRequest.ContentLength = uploadSize
	if err != nil {
		log.Printf("failed writing to create a PUT request %s", err)
		return nil, err
	}
	response, err := client.Do(putRequest)
	if response != nil {
		defer func() {
			_ = response.Body.Close()
		}()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Minute)
	defer cancel()
	if err != nil {
		log.Printf("Failed to upload disk image, aborting image transfer... (%v)", err)
		cancelImageTransfer(transferService, correlationID, ctx)
		return transferService, err
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read response body from ImageIO API (%v)", err)
		cancelImageTransfer(transferService, correlationID, ctx)
		return transferService, fmt.Errorf("failed to read response from ImageIO API (%w)", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		cancelImageTransfer(transferService, correlationID, ctx)
		log.Printf("Unexpected HTTP status code from ImageIO API: %d (%s)", response.StatusCode, responseBody)
	}

	client.CloseIdleConnections()

	return transferService, finalizeImageTransfer(conn, transferService, ctx, diskID)
}

func finalizeImageTransfer(
	conn *ovirtsdk4.Connection,
	transferService *ovirtsdk4.ImageTransferService,
	ctx context.Context,
	diskID string,
) error {
	if _, err := transferService.Finalize().Send(); err != nil {
		return fmt.Errorf("failed to finalize image transfer (%w)", err)
	}
	var notFoundError *ovirtsdk4.NotFoundError
	if err := waitForImageTransferState(
		ctx,
		transferService,
		ovirtsdk4.IMAGETRANSFERPHASE_FINISHED_SUCCESS,
		[]ovirtsdk4.ImageTransferPhase{
			ovirtsdk4.IMAGETRANSFERPHASE_FINISHED_FAILURE,
			ovirtsdk4.IMAGETRANSFERPHASE_PAUSED_SYSTEM,
		},
	); err != nil {
		// If it's a not found error we fall back to checking the disk status below. This is used for <4.4.7
		// where the transfer is removed immediately after finalizing.
		if !errors.As(err, &notFoundError) {
			return err
		}
	}
	for {
		log.Printf("Waiting for disk to become OK...")
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout while waiting for disk to become OK")
		case <-time.After(time.Second * 5):
		}
		req, err := conn.SystemService().DisksService().DiskService(diskID).Get().Send()
		if err != nil {
			if errors.As(err, &notFoundError) {
				return fmt.Errorf("upload failed, disk %s removed (%v)", diskID, err)
			}
			log.Printf("failed to fetch disk %s (%v)", diskID, err)
			continue
		}
		switch req.MustDisk().MustStatus() {
		case ovirtsdk4.DISKSTATUS_OK:
			return nil
		case ovirtsdk4.DISKSTATUS_ILLEGAL:
			return fmt.Errorf("upload failed, disk is in %s status", req.MustDisk().MustStatus())
		}
	}
}

func cancelImageTransfer(
	transferService *ovirtsdk4.ImageTransferService,
	correlationID string,
	ctx context.Context,
) {
	if _, e2 := transferService.Cancel().Query("correlation_id", correlationID).Send(); e2 != nil {
		log.Printf("failed to cancel image upload (%v)", e2)
		return
	}
	if e2 := waitForImageTransferState(
		ctx,
		transferService,
		ovirtsdk4.IMAGETRANSFERPHASE_FINISHED_FAILURE,
		[]ovirtsdk4.ImageTransferPhase{
			ovirtsdk4.IMAGETRANSFERPHASE_FINISHED_SUCCESS,
			ovirtsdk4.IMAGETRANSFERPHASE_PAUSED_SYSTEM,
		},
	); e2 != nil {
		log.Printf("failed to wait for canceled image upload to enter failure state (%v)", e2)
	}
}

func waitForImageTransferState(
	ctx context.Context,
	transferService *ovirtsdk4.ImageTransferService,
	waitForPhase ovirtsdk4.ImageTransferPhase,
	disallowedPhases []ovirtsdk4.ImageTransferPhase,
) error {
	log.Printf("Waiting for image transfer to enter %s state...", waitForPhase)
	var notFoundError *ovirtsdk4.NotFoundError
	for {
		log.Printf("Waiting for image transfer to enter finished state...")
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout while trying to finalize transfer")
		case <-time.After(time.Second * 5):
		}
		req, err := transferService.Get().Send()
		if err != nil {
			if errors.As(err, &notFoundError) {
				// Return the error directly and let the calling party deal with it.
				// This is the case on <4.4.7, where the transfer is removed immediately.
				return err
			} else {
				log.Printf("Error while fetching image transfer (%v)", err)
				continue
			}
		}
		imageTransfer, ok := req.ImageTransfer()
		if !ok {
			return fmt.Errorf("image transfer no longer exists (possibly deleted by engine?) (%v)", err)
		}

		log.Printf("Image transfer is in %s state.", req.MustImageTransfer().MustPhase())
		if imageTransfer.MustPhase() == waitForPhase {
			log.Printf("Image transfer finished.")
			break
		}
		for _, phase := range disallowedPhases {
			if req.MustImageTransfer().MustPhase() == phase {
				return fmt.Errorf(
					"image transfer is in an incorrect state: %s",
					req.MustImageTransfer().MustPhase(),
				)
			}
		}
	}
	return nil
}

func detectUploadUrl(transfer *ovirtsdk4.ImageTransfer) (string, error) {
	// hostUrl means a direct upload to an oVirt node
	insecureClient := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	hostUrl, err := url.Parse(transfer.MustTransferUrl())
	if err == nil {
		optionsReq, err := http.NewRequest(http.MethodOptions, hostUrl.String(), strings.NewReader(""))
		res, err := insecureClient.Do(optionsReq)
		if err == nil && res.StatusCode == 200 {
			return hostUrl.String(), nil
		}
		log.Printf("OPTIONS call to %s failed with %v. Trying the proxy URL", hostUrl.String(), err)
		// can't reach the host url, try the proxy.
	}

	proxyUrl, err := url.Parse(transfer.MustProxyUrl())
	if err != nil {
		log.Printf("failed to parse the proxy url (%s) : %s", transfer.MustProxyUrl(), err)
		return "", err
	}
	optionsReq, err := http.NewRequest(http.MethodOptions, proxyUrl.String(), strings.NewReader(""))
	res, err := insecureClient.Do(optionsReq)
	if err == nil && res.StatusCode == 200 {
		return proxyUrl.String(), nil
	}
	log.Printf("OPTIONS call to %s failed with  %v", proxyUrl.String(), err)
	return "", err
}

func resourceOvirtImageTransferRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	getDiskResp, err := conn.SystemService().DisksService().
		DiskService(d.Id()).Get().Send()
	if err != nil {
		return err
	}

	disk, ok := getDiskResp.Disk()
	if !ok {
		d.SetId("")
		return nil
	}
	d.Set("name", disk.MustAlias())
	d.Set("size", disk.MustProvisionedSize()/int64(math.Pow(2, 30)))
	d.Set("format", disk.MustFormat())
	d.Set("disk_id", disk.MustId())

	if sds, ok := disk.StorageDomains(); ok {
		if len(sds.Slice()) > 0 {
			d.Set("storage_domain_id", sds.Slice()[0].MustId())
		}
	}
	if alias, ok := disk.Alias(); ok {
		d.Set("alias", alias)
	}
	if shareable, ok := disk.Shareable(); ok {
		d.Set("shareable", shareable)
	}
	if sparse, ok := disk.Sparse(); ok {
		d.Set("sparse", sparse)
	}

	return nil
}

func resourceOvirtImageTransferDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	diskService := conn.SystemService().
		DisksService().
		DiskService(d.Id())

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] Now to remove Disk (%s)", d.Id())
		_, e := diskService.Remove().Send()
		if e != nil {
			if _, ok := e.(*ovirtsdk4.NotFoundError); ok {
				log.Printf("[DEBUG] Disk (%s) has been removed", d.Id())
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Error removing Disk (%s): %s", d.Id(), e))
		}
		return resource.RetryableError(fmt.Errorf("Disk (%s) is still being removed", d.Id()))
	})
}

// ImageTransferStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an oVirt image transfer.
func ImageTransferStateRefreshFunc(conn *ovirtsdk4.Connection, transferId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := conn.SystemService().
			ImageTransfersService().
			ImageTransferService(transferId).
			Get().
			Send()
		if err != nil {
			if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
				// should occur only if the transfer was deleted
				return nil, "", nil
			}
			return nil, "", err
		}

		return r.MustImageTransfer(), string(r.MustImageTransfer().MustPhase()), nil
	}
}

func parseQcowSize(header []byte) (uint64, error) {
	isQCOW := string(header[0:4]) == "QFI\xfb"
	if !isQCOW {
		return 0, fmt.Errorf("not a qcow header")
	}
	size := binary.BigEndian.Uint64(header[24:32])
	return size, nil

}

// PrepareForTransfer examine the source url, downloads the file locally if needed
// and return the intended upload size, and format of the image, errors otherwise
func PrepareForTransfer(sourceUrl string) (uploadSize int64, qcowSize uint64, sourceFile *os.File, diskFormat ovirtsdk4.DiskFormat, err error) {
	var sFile *os.File
	if strings.HasPrefix(sourceUrl, "file://") || strings.HasPrefix(sourceUrl, "/") {
		// skip url download, its a local file
		local, err := os.Open(sourceUrl)
		if err != nil {
			return 0, 0, nil, "", err
		}
		sFile = local
	} else {
		resp, err := http.Get(sourceUrl)
		if err != nil {
			return 0, 0, nil, "", err
		}
		defer resp.Body.Close()

		sFile, err = ioutil.TempFile("/tmp", "*-ovirt-image.downloaded")
		io.Copy(sFile, resp.Body)

		// reset cursor to prep for reads
		_, err = sFile.Seek(0, 0)
		if err != nil {
			return 0, 0, nil, "", err
		}

		// remove it when done or cache? maybe --cache
		// defer os.Remove(sFile.Name())
	}
	getFileInfo, err := sFile.Stat()
	if err != nil {
		log.Printf("the sourceUrl is unreachable %s", sourceUrl)
		return 0, 0, nil, "", err
	}

	// gather details about the file
	uploadSize = getFileInfo.Size()
	header := make([]byte, 32)
	_, err = sFile.Read(header)
	if err != nil {
		return 0, 0, nil, "", err
	}
	// we already read 32 bytes - go back to the start of the stream.
	_, err = sFile.Seek(0, 0)
	if err != nil {
		return 0, 0, nil, "", err
	}

	format := ovirtsdk4.DISKFORMAT_COW
	qcowSize, err = parseQcowSize(header)
	if err != nil {
		format = ovirtsdk4.DISKFORMAT_RAW
	}
	log.Printf("upload size is %v", qcowSize)

	// provisioned size must be the virtual size of the QCOW image
	// so must parse the virtual size from the disk. see the UI handling for that
	// for block format, the initial size of disk == uploadSize
	return uploadSize, qcowSize, sFile, format, nil
}

func jobRefreshFunc(conn *ovirtsdk4.Connection, correlationID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		jobResp, err := conn.SystemService().JobsService().List().Search(fmt.Sprintf("correlation_id=%s", correlationID)).Send()
		if err != nil {
			return nil, "", err
		}
		if jobSlice, ok := jobResp.Jobs(); ok && len(jobSlice.Slice()) > 0 {
			jobs := jobSlice.Slice()
			for _, job := range jobs {
				if status, _ := job.Status(); status == ovirtsdk4.JOBSTATUS_STARTED {
					return job, string(job.MustId()), nil
				}
			}
		}

		return nil, "", nil
	}
}
