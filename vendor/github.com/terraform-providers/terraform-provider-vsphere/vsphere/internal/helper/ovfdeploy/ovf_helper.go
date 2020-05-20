package ovfdeploy

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/network"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"io"
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

func DeployOVFAndGetResult(ovfCreateImportSpecResult *types.OvfCreateImportSpecResult, resourcePoolObj *object.ResourcePool,
	folder *object.Folder, host *object.HostSystem, OvfPath string, ovfFromLocal bool) error {

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
				progress = getTotalBytesRead(&currBytesRead) * 100 / totalBytes
				nfcLease.Progress(context.Background(), int32(progress))
				time.Sleep(10 * time.Second)
			}
		}
	}()

	for _, ovfFileItem := range ovfCreateImportSpecResult.FileItem {
		for _, deviceObj := range leaseInfo.DeviceUrl {

			if ovfFileItem.DeviceId == deviceObj.ImportKey {
				if ovfFromLocal {
					absoluteFilePath := ""
					if strings.Contains(OvfPath, string(os.PathSeparator)) {
						absoluteFilePath = string(OvfPath[0 : strings.LastIndex(OvfPath, string(os.PathSeparator))+1])
					}
					vmdkFilePath := absoluteFilePath + ovfFileItem.Path
					log.Print(" [DEBUG] Absolute VMDK path: " + vmdkFilePath)
					file, err := os.Open(vmdkFilePath)
					if err != nil {
						return err
					}
					err = Upload(context.Background(), ovfFileItem, file, deviceObj.Url, ovfFileItem.Size, &currBytesRead)
					if err != nil {
						return fmt.Errorf("error while uploading the file %s %s", vmdkFilePath, err)
					}
					err = file.Close()
					if err != nil {
						log.Printf("error while closing the file %s", vmdkFilePath)
					}

				} else {
					absoluteFilePath := ""
					if strings.Contains(OvfPath, "/") {
						absoluteFilePath = string(OvfPath[0 : strings.LastIndex(OvfPath, "/")+1])
					}
					vmdkFilePath := absoluteFilePath + ovfFileItem.Path
					resp, err := http.Get(vmdkFilePath)
					log.Print("DEBUG Absolute VMDK path: " + vmdkFilePath)
					if err != nil {
						return err
					}
					defer resp.Body.Close()
					err = Upload(context.Background(), ovfFileItem, resp.Body, deviceObj.Url, ovfFileItem.Size, &currBytesRead)
					if err != nil {
						return err
					}
				}
				log.Print("DEBUG : Completed uploading the VMDK file")
			}
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

func Upload(ctx context.Context, item types.OvfFileItem, f io.Reader, url string, size int64, totalBytesRead *int64) error {

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
