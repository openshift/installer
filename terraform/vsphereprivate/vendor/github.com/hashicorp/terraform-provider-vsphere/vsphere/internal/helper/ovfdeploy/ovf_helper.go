package ovfdeploy

import (
	"archive/tar"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/datastore"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/folder"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/resourcepool"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/network"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
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
		for {
			select {
			case <-statusChannel:
				break
			default:
				if totalBytes == 0 {
					_ = nfcLease.Progress(context.Background(), 100)
					return
				}
				log.Printf("Uploaded %v of %v Bytes", getTotalBytesRead(&currBytesRead), totalBytes)
				progress := (getTotalBytesRead(&currBytesRead) / totalBytes) * 100
				_ = nfcLease.Progress(context.Background(), int32(progress))
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

func GetNetworkMapping(client *govmomi.Client, m map[string]interface{}) ([]types.OvfNetworkMapping, error) {
	var ovfNetworkMappings []types.OvfNetworkMapping
	for key, val := range m {
		networkObj, err := network.FromID(client, fmt.Sprint(val))
		if err != nil {
			return nil, err
		}
		networkMapping := types.OvfNetworkMapping{Name: key, Network: networkObj.Reference()}
		ovfNetworkMappings = append(ovfNetworkMappings, networkMapping)
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

func CheckDeploymentOption(client *govmomi.Client, deploymentOption, ovfDescriptor string) error {
	ovfManager := ovf.NewManager(client.Client)

	ovfParseDescriptorParams := types.OvfParseDescriptorParams{}
	ovfParsedDescriptor, err := ovfManager.ParseDescriptor(context.Background(), ovfDescriptor, ovfParseDescriptorParams)
	if err != nil {
		return fmt.Errorf("error while parsing the ovf descriptor file %s", err)
	}
	var validDeployments []string
	for _, option := range ovfParsedDescriptor.DeploymentOption {
		validDeployments = append(validDeployments, option.Key)
		if deploymentOption == option.Key {
			return nil
		}
	}
	// If we get to this point it means that no matches were found.
	return fmt.Errorf("invalid ovf deployment %s specified, valid deployments are: %s", deploymentOption, strings.Join(validDeployments, ", "))
}

type OvfHelper struct {
	AllowUnverifiedSSL bool
	Datastore          *object.Datastore
	DeploymentOption   string
	DeployOva          bool
	DiskProvisioning   string
	FilePath           string
	Folder             *object.Folder
	IsLocal            bool
	Name               string
	HostSystem         *object.HostSystem
	IpAllocationPolicy string
	IpProtocol         string
	NetworkMapping     []types.OvfNetworkMapping
	ResourcePool       *object.ResourcePool
}

type OvfHelperParams struct {
	AllowUnverifiedSSL bool
	DatastoreId        string
	DeploymentOption   string
	DiskProvisioning   string
	FilePath           string
	Folder             string
	HostId             string
	IpAllocationPolicy string
	IpProtocol         string
	Name               string
	NetworkMappings    map[string]interface{}
	OvfUrl             string
	PoolId             string
}

func NewOvfHelper(client *govmomi.Client, o *OvfHelperParams) (*OvfHelper, error) {
	ovfParams := &OvfHelper{
		AllowUnverifiedSSL: o.AllowUnverifiedSSL,
		DeploymentOption:   o.DeploymentOption,
		DiskProvisioning:   o.DiskProvisioning,
		IpAllocationPolicy: o.IpAllocationPolicy,
		IpProtocol:         o.IpProtocol,
		Name:               o.Name,
	}

	ovfParams.DeployOva = false
	ovfParams.IsLocal = true
	ovfParams.FilePath = o.FilePath

	ovfUrl := o.OvfUrl
	if ovfUrl != "" {
		ovfParams.IsLocal = false
		ovfParams.FilePath = ovfUrl
	}

	if strings.HasSuffix(ovfParams.FilePath, ".ova") {
		ovfParams.DeployOva = true
	}

	//Resource pool
	poolID := o.PoolId
	poolObj, err := resourcepool.FromID(client, poolID)
	if err != nil {
		return nil, fmt.Errorf("could not find resource pool ID %q: %s", poolID, err)
	}
	ovfParams.ResourcePool = poolObj

	// Folder
	folderObj, err := folder.VirtualMachineFolderFromObject(client, poolObj, o.Folder)
	if err != nil {
		return nil, err
	}
	ovfParams.Folder = folderObj

	//Host
	hostId := o.HostId
	if hostId == "" {
		return nil, fmt.Errorf("host system ID is required for ovf deployment")
	}
	hostObj, err := hostsystem.FromID(client, hostId)
	if err != nil {
		return nil, fmt.Errorf("could not find host with ID %q: %s", hostId, err)
	}
	ovfParams.HostSystem = hostObj

	//Datastore
	dsId := o.DatastoreId
	if dsId == "" {
		return nil, fmt.Errorf("data store ID is required for ovf deployment")
	}
	dsObj, err := datastore.FromID(client, dsId)
	if err != nil {
		return nil, fmt.Errorf("could not find datastore with ID %q: %s", dsId, err)
	}
	ovfParams.Datastore = dsObj

	//Network Mapping
	networkMapping, err := GetNetworkMapping(client, o.NetworkMappings)
	if err != nil {
		return nil, fmt.Errorf("while getting OVF network mapping: %s", err)
	}
	ovfParams.NetworkMapping = networkMapping

	return ovfParams, nil
}

func (o *OvfHelper) GetImportSpec(client *govmomi.Client) (*types.OvfCreateImportSpecResult, error) {

	hsRef := o.HostSystem.Reference()
	importSpecParam := types.OvfCreateImportSpecParams{
		EntityName:         o.Name,
		HostSystem:         &hsRef,
		NetworkMapping:     o.NetworkMapping,
		IpAllocationPolicy: o.IpAllocationPolicy,
		IpProtocol:         o.IpProtocol,
		DiskProvisioning:   o.DiskProvisioning,
	}

	ovfDescriptor, err := GetOvfDescriptor(o.FilePath, o.DeployOva, o.IsLocal, o.AllowUnverifiedSSL)
	if err != nil {
		return nil, fmt.Errorf("error while reading the ovf file %s, %s ", o.FilePath, err)
	}

	if ovfDescriptor == "" {
		return nil, fmt.Errorf("the given ovf file %s is empty", o.FilePath)
	}

	ovfManager := ovf.NewManager(client.Client)
	deploymentOption := o.DeploymentOption
	if deploymentOption != "" {
		err := CheckDeploymentOption(client, deploymentOption, ovfDescriptor)
		if err != nil {
			return nil, fmt.Errorf("while checking deployment option: %s", err)
		}
		importSpecParam.DeploymentOption = deploymentOption
	}

	is, err := ovfManager.CreateImportSpec(context.Background(), ovfDescriptor,
		o.ResourcePool.Reference(), o.Datastore.Reference(), importSpecParam)
	if len(is.Error) > 0 {
		out := "while getting ovf import spec: \n"
		for _, e := range is.Error {
			out = fmt.Sprintf("%s\n- %s", out, e.LocalizedMessage)
		}
		return nil, fmt.Errorf(out)
	}
	return is, nil
}

func (o *OvfHelper) DeployOvf(spec *types.OvfCreateImportSpecResult) error {
	return DeployOvfAndGetResult(spec, o.ResourcePool, o.Folder, o.HostSystem,
		o.FilePath, o.DeployOva, o.IsLocal, o.AllowUnverifiedSSL)
}
