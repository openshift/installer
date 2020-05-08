// Copyright (C) 2019 oVirt Maintainers
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"bufio"
	"crypto/tls"
	"encoding/binary"
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

	addResp, err := conn.SystemService().DisksService().Add().Disk(disk).Send()
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
	transferService, err := UploadToDisk(conn, sourceFile, diskID, alias, uploadSize)
	if err != nil {
		return err
	}

	log.Printf("finalizing...")
	_, err = transferService.Finalize().Send()

	if err != nil {
		return err
	}

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
func UploadToDisk(conn *ovirtsdk4.Connection, sourceFile *os.File, diskID string, alias string, uploadSize int64) (*ovirtsdk4.ImageTransferService, error) {
	defer sourceFile.Close()
	imageTransfersService := conn.SystemService().ImageTransfersService()
	image := ovirtsdk4.NewImageBuilder().Id(diskID).MustBuild()
	log.Printf("the image to transfer: %s", alias)
	transfer := ovirtsdk4.NewImageTransferBuilder().Image(image).MustBuild()
	transferRes, err := imageTransfersService.Add().ImageTransfer(transfer).Send()
	if err != nil {
		log.Printf("failed to initialize an image transfer for image (%v) : %s", transfer, err)
		return nil, err
	}
	log.Printf("transfer response: %v", transferRes)
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
	_, err = client.Do(putRequest)
	if err != nil {
		return nil, err
	}
	return transferService, nil
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
