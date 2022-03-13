package ovfdeploy

import (
	"archive/tar"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/network"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

func getTotalBytesRead(totalBytes *int64) int64 {
	return atomic.LoadInt64(totalBytes)
}

func incrementTotalBytesRead(totalBytesRead *int64, n int64) {
	atomic.StoreInt64(totalBytesRead, getTotalBytesRead(totalBytesRead)+n)
}

type ProgressReader struct {
	io.Reader
	Reporter func(r int64)
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	pr.Reporter(int64(n))
	return
}

func DeployOvfAndGetResult(ovfCreateImportSpecResult *types.OvfCreateImportSpecResult, resourcePoolObj *object.ResourcePool,
	folder *object.Folder, host *object.HostSystem, filePath string, deployOva bool, fromLocal bool, allowUnverifiedSSL bool) error {

	var currBytesRead int64 = 0
	var totalBytes int64 = 0

	nfcLease, err := resourcePoolObj.ImportVApp(context.Background(), ovfCreateImportSpecResult.ImportSpec, folder, host)
	if err != nil {
		return err
	}

	leaseInfo, err := nfcLease.Wait(context.Background(), ovfCreateImportSpecResult.FileItem)
	if err != nil {
		return err
	}

	u := nfcLease.StartUpdater(context.Background(), leaseInfo)
	defer u.Done()

	for _, ovfFileItem := range ovfCreateImportSpecResult.FileItem {
		totalBytes += ovfFileItem.Size
	}
	log.Printf("Total size of files to upload is %v bytes", totalBytes)

	statusChannel := make(chan bool)
	// Create a go routine to update progress regularly
	go func() {
		var progress int64 = 0
		for {
			select {
			case <-statusChannel:
				break
			default:
				log.Printf("Uploaded %v of %v Bytes", getTotalBytesRead(&currBytesRead), totalBytes)
				if totalBytes == 0 {
					break
				}
				progress = (getTotalBytesRead(&currBytesRead) / totalBytes) * 100
				nfcLease.Progress(context.Background(), int32(progress))
				time.Sleep(10 * time.Second)
			}
		}
	}()

	for _, ovfFileItem := range ovfCreateImportSpecResult.FileItem {
		for _, deviceObj := range leaseInfo.DeviceUrl {

			if ovfFileItem.DeviceId != deviceObj.ImportKey {
				continue
			}
			if !deployOva {
				if fromLocal {
					err = uploadDisksFromLocal(filePath, ovfFileItem, deviceObj, &currBytesRead)
				} else {
					err = uploadDisksFromUrl(filePath, ovfFileItem, deviceObj, &currBytesRead, allowUnverifiedSSL)
				}
			} else {
				if fromLocal {
					err = uploadOvaDisksFromLocal(filePath, ovfFileItem, deviceObj, &currBytesRead)
				} else {
					err = uploadOvaDisksFromUrl(filePath, ovfFileItem, deviceObj, &currBytesRead, allowUnverifiedSSL)
				}
			}
			if err != nil {
				return fmt.Errorf("error while uploading the disk %s %s", ovfFileItem.Path, err)
			}
			log.Print(" DEBUG : Completed uploading the vmdk file", ovfFileItem.Path)
		}
	}
	statusChannel <- true
	err = nfcLease.Progress(context.Background(), 100)
	if err != nil {
		return err
	}
	if err = nfcLease.Complete(context.Background()); err != nil {
		return err
	}
	return nil
}

func upload(ctx context.Context, item types.OvfFileItem, f io.Reader, url string, size int64, totalBytesRead *int64) error {

	u, err := soap.ParseURL(url)
	if err != nil {
		return err
	}
	c := soap.NewClient(u, true)

	param := soap.Upload{
		ContentLength: size,
	}

	if item.Create {
		param.Method = "PUT"
		param.Headers = map[string]string{
			"Overwrite": "t",
		}
	} else {
		param.Method = "POST"
		param.Type = "application/x-vnd.vmware-streamVmdk"
	}

	pr := &ProgressReader{f, func(r int64) {
		incrementTotalBytesRead(totalBytesRead, r)
	}}
	f = pr

	req, err := http.NewRequest(param.Method, url, f)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	req.ContentLength = param.ContentLength
	req.Header.Set("Content-Type", param.Type)

	for k, v := range param.Headers {
		req.Header.Add(k, v)
	}
	if param.Ticket != nil {
		req.AddCookie(param.Ticket)
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusCreated:
	default:
		err = errors.New(res.Status)
	}
	return err
}

func uploadDisksFromLocal(filePath string, ovfFileItem types.OvfFileItem, deviceObj types.HttpNfcLeaseDeviceUrl, currBytesRead *int64) error {
	absoluteFilePath := ""
	if strings.Contains(filePath, string(os.PathSeparator)) {
		absoluteFilePath = string(filePath[0 : strings.LastIndex(filePath, string(os.PathSeparator))+1])
	}
	vmdkFilePath := absoluteFilePath + ovfFileItem.Path
	log.Print(" [DEBUG] Absolute vmdk path: " + vmdkFilePath)
	file, err := os.Open(vmdkFilePath)
	if err != nil {
		return err
	}
	err = upload(context.Background(), ovfFileItem, file, deviceObj.Url, ovfFileItem.Size, currBytesRead)
	if err != nil {
		return fmt.Errorf("error while uploading the file %s %s", vmdkFilePath, err)
	}
	err = file.Close()
	if err != nil {
		log.Printf("error while closing the file %s", vmdkFilePath)
	}
	return nil
}

func uploadDisksFromUrl(filePath string, ovfFileItem types.OvfFileItem, deviceObj types.HttpNfcLeaseDeviceUrl, currBytesRead *int64,
	allowUnverifiedSSL bool) error {
	absoluteFilePath := ""
	if strings.Contains(filePath, "/") {
		absoluteFilePath = string(filePath[0 : strings.LastIndex(filePath, "/")+1])
	}
	vmdkFilePath := absoluteFilePath + ovfFileItem.Path
	client := getClient(allowUnverifiedSSL)
	resp, err := client.Get(vmdkFilePath)
	log.Print(" [DEBUG] Absolute vmdk path: " + vmdkFilePath)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = upload(context.Background(), ovfFileItem, resp.Body, deviceObj.Url, ovfFileItem.Size, currBytesRead)
	return err
}

func uploadOvaDisksFromLocal(filePath string, ovfFileItem types.OvfFileItem, deviceObj types.HttpNfcLeaseDeviceUrl, currBytesRead *int64) error {
	diskName := ovfFileItem.Path
	ovaFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer ovaFile.Close()

	err = findAndUploadDiskFromOva(ovaFile, diskName, ovfFileItem, deviceObj, currBytesRead)
	return err
}

func uploadOvaDisksFromUrl(filePath string, ovfFileItem types.OvfFileItem, deviceObj types.HttpNfcLeaseDeviceUrl, currBytesRead *int64,
	allowUnverifiedSSL bool) error {
	diskName := ovfFileItem.Path
	client := getClient(allowUnverifiedSSL)
	resp, err := client.Get(filePath)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		err = findAndUploadDiskFromOva(resp.Body, diskName, ovfFileItem, deviceObj, currBytesRead)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("got status %d while getting the file from remote url %s ", resp.StatusCode, filePath)
	}
	return nil
}

func GetOvfDescriptor(filePath string, deployOva bool, fromLocal bool, allowUnverifiedSSL bool) (string, error) {

	ovfDescriptor := ""
	if !deployOva {
		if fromLocal {
			fileBuffer, err := ioutil.ReadFile(filePath)
			if err != nil {
				return "", err
			}
			ovfDescriptor = string(fileBuffer)
		} else {
			client := getClient(allowUnverifiedSSL)
			resp, err := client.Get(filePath)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return "", err
				}
				ovfDescriptor = string(bodyBytes)
			}
		}
	} else {
		if fromLocal {
			ovaFile, err := os.Open(filePath)
			if err != nil {
				return "", err
			}
			defer ovaFile.Close()
			ovfDescriptor, err = getOvfDescriptorFromOva(ovaFile)
			if err != nil {
				return "", err
			}

		} else {
			client := getClient(allowUnverifiedSSL)
			resp, err := client.Get(filePath)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {

				ovfDescriptor, err = getOvfDescriptorFromOva(resp.Body)
				if err != nil {
					return "", err
				}
			} else {
				return "", fmt.Errorf("got status %d while getting the file from remote url %s ", resp.StatusCode, filePath)
			}
		}
	}
	return ovfDescriptor, nil
}

func getOvfDescriptorFromOva(ovaFile io.Reader) (string, error) {
	ovaReader := tar.NewReader(ovaFile)
	for {
		fileHdr, err := ovaReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if strings.HasSuffix(fileHdr.Name, ".ovf") {
			content, _ := ioutil.ReadAll(ovaReader)
			ovfDescriptor := string(content)
			return ovfDescriptor, nil
		}
	}
	return "", fmt.Errorf("ovf file not found inside the ova")
}

func findAndUploadDiskFromOva(ovaFile io.Reader, diskName string, ovfFileItem types.OvfFileItem, deviceObj types.HttpNfcLeaseDeviceUrl, currBytesRead *int64) error {
	ovaReader := tar.NewReader(ovaFile)
	for {
		fileHdr, err := ovaReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if fileHdr.Name == diskName {
			err = upload(context.Background(), ovfFileItem, ovaReader, deviceObj.Url, ovfFileItem.Size, currBytesRead)
			if err != nil {
				return fmt.Errorf("error while uploading the file %s %s", diskName, err)
			}
			return nil
		}
	}
	return fmt.Errorf("disk %s not found inside ova", diskName)
}

func GetNetworkMapping(client *govmomi.Client, d *schema.ResourceData) ([]types.OvfNetworkMapping, error) {
	var ovfNetworkMappings []types.OvfNetworkMapping
	m := d.Get("ovf_deploy.0.ovf_network_map").(map[string]interface{})
	if m != nil {
		for key, val := range m {

			networkObj, err := network.FromID(client, fmt.Sprint(val))
			if err != nil {
				return nil, err
			}
			networkMapping := types.OvfNetworkMapping{Name: key, Network: networkObj.Reference()}
			ovfNetworkMappings = append(ovfNetworkMappings, networkMapping)
		}
	}
	return ovfNetworkMappings, nil
}

func getClient(allowUnverifiedSSL bool) *http.Client {
	if allowUnverifiedSSL {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		return &http.Client{Transport: tr}
	} else {
		return &http.Client{}
	}
}
