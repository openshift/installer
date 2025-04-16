package azure

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/pageblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/sirupsen/logrus"

	azic "github.com/openshift/installer/pkg/asset/installconfig/azure"
	aztypes "github.com/openshift/installer/pkg/types/azure"
)

var (
	vaultsClient *armkeyvault.VaultsClient
	keysClient   *armkeyvault.KeysClient
)

// CreateStorageAccountInput contains the input parameters for creating a
// storage account.
type CreateStorageAccountInput struct {
	SubscriptionID     string
	ResourceGroupName  string
	StorageAccountName string
	Region             string
	AuthType           azic.AuthenticationType
	Tags               map[string]*string
	CustomerManagedKey *aztypes.CustomerManagedKey
	CloudName          aztypes.CloudEnvironment
	TokenCredential    azcore.TokenCredential
	CloudConfiguration cloud.Configuration
}

// CreateStorageAccountOutput contains the return values after creating a
// storage account.
type CreateStorageAccountOutput struct {
	StorageAccount        *armstorage.Account
	StorageAccountsClient *armstorage.AccountsClient
	StorageClientFactory  *armstorage.ClientFactory
	StorageAccountKeys    []armstorage.AccountKey
}

// CreateStorageAccount creates a new storage account.
func CreateStorageAccount(ctx context.Context, in *CreateStorageAccountInput) (*CreateStorageAccountOutput, error) {
	minimumTLSVersion := armstorage.MinimumTLSVersionTLS10
	cloudConfiguration := in.CloudConfiguration

	/* XXX: Do we support other clouds? */
	switch in.CloudName {
	case aztypes.PublicCloud:
		minimumTLSVersion = armstorage.MinimumTLSVersionTLS12
	case aztypes.USGovernmentCloud:
		minimumTLSVersion = armstorage.MinimumTLSVersionTLS12
	}

	allowSharedKeyAccess := true
	if in.AuthType == azic.ManagedIdentityAuth {
		allowSharedKeyAccess = false
	}

	storageClientFactory, err := armstorage.NewClientFactory(
		in.SubscriptionID,
		in.TokenCredential,
		&arm.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Cloud: cloudConfiguration,
				//Transport: ...,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage account factory %w", err)
	}

	sku := armstorage.SKU{
		Name: to.Ptr(armstorage.SKUNameStandardLRS),
	}
	accountCreateParameters := armstorage.AccountCreateParameters{
		Identity: nil,
		Kind:     to.Ptr(armstorage.KindStorageV2),
		Location: to.Ptr(in.Region),
		SKU:      &sku,
		Properties: &armstorage.AccountPropertiesCreateParameters{
			AllowBlobPublicAccess:       to.Ptr(false),
			AllowSharedKeyAccess:        to.Ptr(allowSharedKeyAccess),
			IsLocalUserEnabled:          to.Ptr(true),
			LargeFileSharesState:        to.Ptr(armstorage.LargeFileSharesStateEnabled),
			PublicNetworkAccess:         to.Ptr(armstorage.PublicNetworkAccessEnabled),
			MinimumTLSVersion:           &minimumTLSVersion,
			AllowCrossTenantReplication: to.Ptr(false), // must remain false to comply with BAFIN and PCI-DSS regulations
		},
		Tags: in.Tags,
	}

	if in.CustomerManagedKey != nil && in.CustomerManagedKey.KeyVault.Name != "" {
		// When encryption is enabled, Ignition is is stored as a page blob
		// (and not a block blob). To support this case, `Kind` can continue to be
		// `StorageV2` and yhe `SKU` needs to be `Premium_LRS`.
		//https://learn.microsoft.com/en-us/azure/storage/common/storage-account-create?tabs=azure-portal
		sku = armstorage.SKU{
			Name: to.Ptr(armstorage.SKUNamePremiumLRS),
		}
		identity := armstorage.Identity{
			Type: to.Ptr(armstorage.IdentityTypeUserAssigned),
			UserAssignedIdentities: map[string]*armstorage.UserAssignedIdentity{
				fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/%s",
					in.SubscriptionID,
					in.CustomerManagedKey.KeyVault.ResourceGroup,
					in.CustomerManagedKey.UserAssignedIdentityKey,
				): {},
			},
		}
		logrus.Debugf("Generating Encrytption for Storage Account using Customer Managed Key")
		encryption, err := GenerateStorageAccountEncryption(
			ctx,
			&CustomerManagedKeyInput{
				SubscriptionID:     in.SubscriptionID,
				ResourceGroupName:  in.ResourceGroupName,
				CustomerManagedKey: in.CustomerManagedKey,
				TokenCredential:    in.TokenCredential,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error generating encryption information for provided customer managed key: %w", err)
		}
		accountCreateParameters.Identity = &identity
		accountCreateParameters.SKU = &sku
		accountCreateParameters.Properties.Encryption = encryption
		accountCreateParameters.Properties.AllowBlobPublicAccess = to.Ptr(true)
	}

	logrus.Debugf("Creating storage account")
	accountsClient := storageClientFactory.NewAccountsClient()
	pollerResponse, err := accountsClient.BeginCreate(
		ctx,
		in.ResourceGroupName,
		in.StorageAccountName,
		accountCreateParameters,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating storage account %s: %w", in.StorageAccountName, err)
	}

	pollDoneResponse, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error waiting for creation of storage account %s: %w", in.StorageAccountName, err)
	}

	logrus.Debugf("Getting storage keys")
	listKeysResponse, err := accountsClient.ListKeys(ctx, in.ResourceGroupName, in.StorageAccountName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve storage account keys for %s: %w", in.StorageAccountName, err)
	}

	out := &CreateStorageAccountOutput{
		StorageAccount:        to.Ptr(pollDoneResponse.Account),
		StorageAccountsClient: accountsClient,
		StorageClientFactory:  storageClientFactory,
	}

	for _, key := range listKeysResponse.Keys {
		out.StorageAccountKeys = append(out.StorageAccountKeys, *key)
	}

	return out, nil
}

// CreateBlobContainerInput contains the input parameters used for creating a
// blob storage container.
type CreateBlobContainerInput struct {
	SubscriptionID       string
	ResourceGroupName    string
	StorageAccountName   string
	ContainerName        string
	PublicAccess         *armstorage.PublicAccess
	StorageClientFactory *armstorage.ClientFactory
}

// CreateBlobContainerOutput contains the return values after creating a blob
// storage container.
type CreateBlobContainerOutput struct {
	BlobContainer *armstorage.BlobContainer
}

// CreateBlobContainer creates a blob container in a storage account.
func CreateBlobContainer(ctx context.Context, in *CreateBlobContainerInput) (*CreateBlobContainerOutput, error) {
	blobContainersClient := in.StorageClientFactory.NewBlobContainersClient()

	logrus.Debugf("Creating blob container")
	blobContainerResponse, err := blobContainersClient.Create(
		ctx,
		in.ResourceGroupName,
		in.StorageAccountName,
		in.ContainerName,
		armstorage.BlobContainer{
			ContainerProperties: &armstorage.ContainerProperties{
				PublicAccess: in.PublicAccess,
			},
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob container %s: %w", in.ContainerName, err)
	}

	return &CreateBlobContainerOutput{
		BlobContainer: to.Ptr(blobContainerResponse.BlobContainer),
	}, nil
}

// CreatePageBlobInput containers the input parameters used for creating a page
// blob.
type CreatePageBlobInput struct {
	StorageURL         string
	BlobURL            string
	ImageURL           string
	StorageAccountName string
	BootstrapIgnData   []byte
	ImageLength        int64
	StorageAccountKeys []armstorage.AccountKey
	CloudConfiguration cloud.Configuration
}

// CreatePageBlobOutput contains the return values after creating a page blob.
type CreatePageBlobOutput struct {
	PageBlobClient      *pageblob.Client
	SharedKeyCredential *azblob.SharedKeyCredential
}

// CreatePageBlob creates a blob and uploads a file from a URL to it.
func CreatePageBlob(ctx context.Context, in *CreatePageBlobInput) (string, error) {
	logrus.Debugf("Getting page blob credentials")

	// XXX: Should try all of them until one is successful
	sharedKeyCredential, err := azblob.NewSharedKeyCredential(in.StorageAccountName, *in.StorageAccountKeys[0].Value)
	if err != nil {
		return "", fmt.Errorf("failed to get shared credentials for storage account: %w", err)
	}

	logrus.Debugf("Getting page blob client")
	pageBlobClient, err := pageblob.NewClientWithSharedKeyCredential(
		in.BlobURL,
		sharedKeyCredential,
		&pageblob.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Cloud: in.CloudConfiguration,
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to get page blob client: %w", err)
	}

	logrus.Debugf("Creating Page blob and uploading image to it")
	if in.ImageURL == "" {
		_, err = pageBlobClient.Create(ctx, in.ImageLength, nil)
		if err != nil {
			return "", fmt.Errorf("failed to create page blob with image contents: %w", err)
		}
		// This image (example: ignition shim) needs to be uploaded from a local file.
		err = doUploadPages(ctx, pageBlobClient, in.BootstrapIgnData, in.ImageLength)
		if err != nil {
			return "", fmt.Errorf("failed to upload page blob image contents: %w", err)
		}
	} else {
		// This is used in terraform, not sure if it matters
		metadata := map[string]*string{
			"source_uri": to.Ptr(in.ImageURL),
		}

		_, err = pageBlobClient.Create(ctx, in.ImageLength, &pageblob.CreateOptions{
			Metadata: metadata,
		})
		if err != nil {
			return "", fmt.Errorf("failed to create page blob with image URL: %w", err)
		}

		err = doUploadPagesFromURL(ctx, pageBlobClient, in.ImageURL, in.ImageLength)
		if err != nil {
			return "", fmt.Errorf("failed to upload page blob image from URL %s: %w", in.ImageURL, err)
		}
	}

	// Is this addition OK for when CreatePageBlob() is called from InfraReady()
	sasURL, err := pageBlobClient.GetSASURL(sas.BlobPermissions{Read: true}, time.Now().Add(time.Minute*60), &blob.GetSASURLOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get Page Blob SAS URL: %w", err)
	}
	return sasURL, nil
}

func doUploadPages(ctx context.Context, pageBlobClient *pageblob.Client, imageData []byte, imageLength int64) error {
	logrus.Debugf("Uploading to Page Blob with Image of length :%d", imageLength)

	// Page blobs file size must be a multiple of 512, hence a little padding is needed to push the file.
	// imageLength has already been adjusted to the next highest size divisible by 512.
	// So, here we are padding the image to match this size.
	// Bootstrap Ignition is a json file. For parsing of this file to succeed with the padding, the
	// file needs to end with a }.
	logrus.Debugf("Original Image length: %d", int64(len(imageData)))
	padding := imageLength - int64(len(imageData))
	paddingString := strings.Repeat(" ", int(padding)) + string(imageData[len(imageData)-1])
	imageData = append(imageData[0:len(imageData)-1], paddingString...)
	logrus.Debugf("New Image length (after padding): %d", int64(len(imageData)))

	pageSize := int64(1024 * 1024 * 4)
	newOffset := int64(0)
	remainingImageLength := imageLength

	for remainingImageLength > 0 {
		if remainingImageLength < pageSize {
			pageSize = remainingImageLength
		}

		logrus.Debugf("Uploading pages with Offset :%d and Count :%d", newOffset, pageSize)

		_, err := pageBlobClient.UploadPages(
			ctx,
			streaming.NopCloser(bytes.NewReader(imageData)),
			blob.HTTPRange{
				Offset: newOffset,
				Count:  pageSize,
			},
			nil)
		if err != nil {
			return fmt.Errorf("failed uploading Image to page blob: %w", err)
		}
		newOffset += pageSize
		remainingImageLength -= pageSize
		logrus.Debugf("newOffset :%d and remainingImageLength :%d", newOffset, remainingImageLength)
	}
	return nil
}

func doUploadPagesFromURL(ctx context.Context, pageBlobClient *pageblob.Client, imageURL string, imageLength int64) error {
	// Azure only allows 4MB chunks, See
	// https://docs.microsoft.com/rest/api/storageservices/put-page-from-url
	pageSize := int64(1024 * 1024 * 4)
	leftOverBytes := imageLength % pageSize
	offset := int64(0)
	pages := int64(0)

	if imageLength > pageSize {
		pages = imageLength / pageSize
		if imageLength%pageSize > 0 {
			pages++
		}
	} else {
		pageSize = imageLength
		pages = 1
	}

	threadsPerGroup := int64(64)
	if pages < threadsPerGroup {
		threadsPerGroup = pages
	}

	threadGroups := pages / threadsPerGroup
	if pages%threadsPerGroup > 0 {
		threadGroups++
	}

	var wg sync.WaitGroup
	var threadError error
	var res error

	pagesLeft := pages
	for threadGroup := int64(0); threadGroup < threadGroups; threadGroup++ {
		if pagesLeft < threadsPerGroup {
			threadsPerGroup = pagesLeft
		}

		errors := make(chan error, 1)
		defer close(errors)

		results := make(chan int64, threadsPerGroup)
		defer close(results)

		for thread := int64(0); thread < threadsPerGroup; thread++ {
			if offset+pageSize >= imageLength && leftOverBytes > 0 {
				pageSize = leftOverBytes
				leftOverBytes = 0
			} else if offset > imageLength {
				break
			}

			wg.Add(1)
			go func(ctx context.Context, source string, thread, sourceOffset, destOffset, count int64, wg *sync.WaitGroup) {
				defer wg.Done()
				var err error
				nretries := 3
				for i := 0; i < nretries; i++ {
					_, err = pageBlobClient.UploadPagesFromURL(ctx, imageURL, sourceOffset, destOffset, count, nil)
					if err == nil {
						break
					}
				}
				errors <- err
				results <- thread
			}(ctx, imageURL, thread, offset, offset, pageSize, &wg)

			offset += pageSize
		}
		pagesLeft -= threadsPerGroup
		for thread := int64(0); thread < threadsPerGroup; thread++ {
			threadError = <-errors

			// XXX: Save first error only. Should we care about the
			// rest?
			if threadError != nil && res == nil {
				res = threadError
			}
			<-results
		}
		wg.Wait()
		if res != nil {
			logrus.Debug("Failed to upload rhcos image")
			break
		}

		logrus.Debugf("%d out of %d pages uploaded", pages-pagesLeft, pages)
	}

	logrus.Debugf("Done uploading")
	return res
}

// CreateBlockBlobInput containers the input parameters used for creating a
// block blob.
type CreateBlockBlobInput struct {
	StorageURL         string
	BlobURL            string
	StorageAccountName string
	BootstrapIgnData   []byte
	StorageAccountKeys []armstorage.AccountKey
	CloudConfiguration cloud.Configuration
}

// CreateBlockBlobOutput contains the return values after creating a block
// blob.
type CreateBlockBlobOutput struct {
	PageBlobClient      *pageblob.Client
	SharedKeyCredential *azblob.SharedKeyCredential
}

// CreateBlockBlob creates a block blob and uploads a file from a URL to it.
func CreateBlockBlob(ctx context.Context, in *CreateBlockBlobInput) (string, error) {
	logrus.Debugf("Getting block blob credentials")

	// XXX: Should try all of them until one is successful
	sharedKeyCredential, err := azblob.NewSharedKeyCredential(in.StorageAccountName, *in.StorageAccountKeys[0].Value)
	if err != nil {
		return "", fmt.Errorf("failed to get shared crdentials for storage account: %w", err)
	}

	logrus.Debugf("Getting block blob client")
	blockBlobClient, err := blockblob.NewClientWithSharedKeyCredential(
		in.BlobURL,
		sharedKeyCredential,
		&blockblob.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Cloud: in.CloudConfiguration,
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to get page blob client: %w", err)
	}

	logrus.Debugf("Creating block blob")

	accessTier := blob.AccessTierHot
	_, err = blockBlobClient.Upload(ctx, streaming.NopCloser(bytes.NewReader(in.BootstrapIgnData)), &blockblob.UploadOptions{
		Tier: &accessTier,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create block blob: %w", err)
	}

	sasURL, err := blockBlobClient.GetSASURL(sas.BlobPermissions{Read: true}, time.Now().Add(time.Minute*60), &blob.GetSASURLOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get SAS URL: %w", err)
	}

	return sasURL, nil
}

// CustomerManagedKeyInput contains the input parameters for creating the
// customer managed key and identity.
type CustomerManagedKeyInput struct {
	SubscriptionID     string
	ResourceGroupName  string
	CustomerManagedKey *aztypes.CustomerManagedKey
	TokenCredential    azcore.TokenCredential
}

// GenerateStorageAccountEncryption generates all the Encryption information for the Storage Account
// using the Customer Managed Key.
func GenerateStorageAccountEncryption(ctx context.Context, in *CustomerManagedKeyInput) (*armstorage.Encryption, error) {
	logrus.Debugf("Generating Encryption for Storage Account")

	if in.CustomerManagedKey == nil {
		logrus.Debugf("No Customer Managed Key provided. So, Encryption not enabled on storage account.")
		return &armstorage.Encryption{}, nil
	}

	keyvaultClientFactory, err := armkeyvault.NewClientFactory(
		in.SubscriptionID,
		in.TokenCredential,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get key vault client factory %w", err)
	}

	keysClient = keyvaultClientFactory.NewKeysClient()

	_, err = keysClient.Get(
		ctx,
		in.CustomerManagedKey.KeyVault.ResourceGroup,
		in.CustomerManagedKey.KeyVault.Name,
		in.CustomerManagedKey.KeyVault.KeyName,
		&armkeyvault.KeysClientGetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get customer managed key %s from key vault %s: %w", in.CustomerManagedKey.KeyVault.KeyName, in.CustomerManagedKey.KeyVault.Name, err)
	}

	vaultsClient = keyvaultClientFactory.NewVaultsClient()

	keyVault, err := vaultsClient.Get(
		ctx,
		in.CustomerManagedKey.KeyVault.ResourceGroup,
		in.CustomerManagedKey.KeyVault.Name,
		&armkeyvault.VaultsClientGetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get key vault %s which contains customer managed key: %w", in.CustomerManagedKey.KeyVault.Name, err)
	}

	encryption := &armstorage.Encryption{
		Services: &armstorage.EncryptionServices{
			Blob: &armstorage.EncryptionService{
				Enabled: to.Ptr(true),
				KeyType: to.Ptr(armstorage.KeyTypeAccount),
			},
			File: &armstorage.EncryptionService{
				Enabled: to.Ptr(true),
				KeyType: to.Ptr(armstorage.KeyTypeAccount),
			},
		},
		EncryptionIdentity: &armstorage.EncryptionIdentity{
			EncryptionUserAssignedIdentity: to.Ptr(fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/%s",
				in.SubscriptionID,
				in.CustomerManagedKey.KeyVault.ResourceGroup,
				in.CustomerManagedKey.UserAssignedIdentityKey,
			)),
		},
		KeySource: to.Ptr(armstorage.KeySourceMicrosoftKeyvault),
		KeyVaultProperties: &armstorage.KeyVaultProperties{
			KeyName:     to.Ptr(in.CustomerManagedKey.KeyVault.KeyName),
			KeyVersion:  to.Ptr(""),
			KeyVaultURI: keyVault.Properties.VaultURI,
		},
	}

	return encryption, nil
}
