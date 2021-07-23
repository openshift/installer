package azurestackprivate

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/Microsoft/azure-vhd-utils/vhdcore/common"
	"github.com/Microsoft/azure-vhd-utils/vhdcore/diskstream"
	"github.com/Microsoft/azure-vhd-utils/vhdcore/validator"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/openshift/azure-sdk-for-go/storage"

	"github.com/openshift/installer/pkg/terraform/exec/plugins/azurestackprivate/upload"
	"github.com/openshift/installer/pkg/terraform/exec/plugins/azurestackprivate/upload/metadata"
)

const (
	// PageBlobPageSize is the size of the page blob in bytes.
	PageBlobPageSize  int64 = 2 * 1024 * 1024
	apiProfileVersion       = "2019-07-07"
)

func resourceArmVHDBlob() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVHDBlobCreate,
		Read:   resourceArmStorageBlobRead,
		Exists: resourceArmStorageBlobExists,
		Delete: resourceArmStorageBlobDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_container_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmVHDBlobCreate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	cont := d.Get("storage_container_name").(string)
	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	source := d.Get("source").(string)
	storageAccountName := d.Get("storage_account_name").(string)

	ensureVHDSanity(source)
	diskStream, err := diskstream.CreateNewDiskStream(source)
	if err != nil {
		return err
	}
	defer diskStream.Close()

	storageAccountKey, accountExists, err := armClient.getKeyForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		return fmt.Errorf("storage account %q not found in resource group %q while uploading RHCOS image", storageAccountName, resourceGroupName)
	}

	blobServiceBaseURL := strings.TrimPrefix(armClient.environment.ResourceManagerEndpoint, "https://management.")
	storageClient, err := storage.NewClient(storageAccountName, storageAccountKey, blobServiceBaseURL, apiProfileVersion, true)
	if err != nil {
		return err
	}
	blobClient := storageClient.GetBlobService()

	localMetaData, err := getMetaDataFromLocalVHD(source)
	if err != nil {
		return err
	}
	var rangesToSkip []*common.IndexRange
	createBlob(blobClient, cont, name, diskStream.GetSize(), localMetaData)

	uploadableRanges, err := upload.LocateUploadableRanges(diskStream, rangesToSkip, PageBlobPageSize)
	if err != nil {
		return err
	}

	uploadableRanges, err = upload.DetectEmptyRanges(diskStream, uploadableRanges)
	if err != nil {
		return err
	}

	cxt := &upload.DiskUploadContext{
		VhdStream:             diskStream,
		UploadableRanges:      uploadableRanges,
		AlreadyProcessedBytes: common.TotalRangeLength(rangesToSkip),
		BlobServiceClient:     blobClient,
		ContainerName:         cont,
		BlobName:              name,
		Parallelism:           8 * runtime.NumCPU(),
		Resume:                false,
	}
	log.Printf("[INFO] Creating blob %q in storage account %q", name, storageAccountName)

	err = upload.Upload(cxt)
	if err != nil {
		return err
	}

	// Delete the VHD because we just unzipped a 16GB file in a dir the user is
	// probably unaware of.
	err = os.Remove(source)
	if err != nil {
		log.Printf("[WARN] Unable to delete temp file %s after upload, please delete manually: %s", source, err.Error())
	}

	d.SetId(name)
	return resourceArmStorageBlobRead(d, meta)
}

func getMetaDataFromLocalVHD(vhdPath string) (*metadata.MetaData, error) {
	fileStat, err := getFileStat(vhdPath)
	if err != nil {
		return nil, err
	}

	fileMetaData := &metadata.FileMetaData{
		FileName:         fileStat.Name(),
		FileSize:         fileStat.Size(),
		LastModifiedTime: fileStat.ModTime(),
	}

	diskStream, err := diskstream.CreateNewDiskStream(vhdPath)
	if err != nil {
		return nil, err
	}
	defer diskStream.Close()
	return &metadata.MetaData{
		FileMetaData: fileMetaData,
	}, nil
}

// getFileStat returns os.FileInfo of a file.
//
func getFileStat(filePath string) (os.FileInfo, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("fileMetaData.getFileStat: %v", err)
	}
	defer fd.Close()
	return fd.Stat()
}

// createBlob creates a page blob of specific size and sets custom metadata
// The parameter client is the Azure blob service client, parameter containerName is the name of an existing container
// in which the page blob needs to be created, parameter blobName is name for the new page blob, size is the size of
// the new page blob in bytes and parameter vhdMetaData is the custom metadata to be associacted with the page blob
//
func createBlob(client storage.BlobStorageClient, containerName, blobName string, size int64, vhdMetaData *metadata.MetaData) {
	if err := client.PutPageBlob(containerName, blobName, size, nil); err != nil {
		log.Fatal(err)
	}
	m, _ := vhdMetaData.ToMap()
	if err := client.SetBlobMetadata(containerName, blobName, m, make(map[string]string)); err != nil {
		log.Fatal(err)
	}
}

// ensureVHDSanity ensure is VHD is valid for Azure.
//
func ensureVHDSanity(localVHDPath string) {
	if err := validator.ValidateVhd(localVHDPath); err != nil {
		log.Fatal(err)
	}

	if err := validator.ValidateVhdSize(localVHDPath); err != nil {
		log.Fatal(err)
	}
}

func resourceArmStorageBlobRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	resourceGroupName := d.Get("resource_group_name").(string)
	storageAccountName := d.Get("storage_account_name").(string)

	blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		log.Printf("[DEBUG] Storage account %q not found, removing blob %q from state", storageAccountName, d.Id())
		d.SetId("")
		return nil
	}

	exists, err := resourceArmStorageBlobExists(d, meta)
	if err != nil {
		return err
	}

	if !exists {
		// Exists already removed this from state
		return nil
	}

	name := d.Get("name").(string)
	storageContainerName := d.Get("storage_container_name").(string)

	container := blobClient.GetContainerReference(storageContainerName)
	blob := container.GetBlobReference(name)
	url := blob.GetURL()
	if url == "" {
		log.Printf("[INFO] URL for %q is empty", name)
	}
	d.Set("url", url)

	return nil
}

func resourceArmStorageBlobExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	armClient := meta.(*ArmClient)
	ctx := armClient.StopContext

	resourceGroupName := d.Get("resource_group_name").(string)
	storageAccountName := d.Get("storage_account_name").(string)

	blobClient, accountExists, err := armClient.getBlobStorageClientForStorageAccount(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return false, err
	}
	if !accountExists {
		log.Printf("[DEBUG] Storage account %q not found, removing blob %q from state", storageAccountName, d.Id())
		d.SetId("")
		return false, nil
	}

	name := d.Get("name").(string)
	storageContainerName := d.Get("storage_container_name").(string)

	log.Printf("[INFO] Checking for existence of storage blob %q.", name)
	container := blobClient.GetContainerReference(storageContainerName)
	blob := container.GetBlobReference(name)
	exists, err := blob.Exists()
	if err != nil {
		return false, fmt.Errorf("error testing existence of storage blob %q: %s", name, err)
	}

	if !exists {
		log.Printf("[INFO] Storage blob %q no longer exists, removing from state...", name)
		d.SetId("")
	}

	return exists, nil
}

// We do not use this Terraform provider for deletion but this implementation is required.
func resourceArmStorageBlobDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
