package azure

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/pageblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/sirupsen/logrus"

	aztypes "github.com/openshift/installer/pkg/types/azure"
)

// CreateStorageAccountInput contains the input parameters for creating a
// storage account.
type CreateStorageAccountInput struct {
	SubscriptionID     string
	ResourceGroupName  string
	StorageAccountName string
	Region             string
	Tags               map[string]*string
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

	logrus.Debugf("Creating storage account")
	accountsClient := storageClientFactory.NewAccountsClient()
	pollerResponse, err := accountsClient.BeginCreate(
		ctx,
		in.ResourceGroupName,
		in.StorageAccountName,
		armstorage.AccountCreateParameters{
			Kind:     to.Ptr(armstorage.KindStorageV2),
			Location: to.Ptr(in.Region),
			SKU: &armstorage.SKU{
				Name: to.Ptr(armstorage.SKUNameStandardLRS), // XXX Premium_LRS if disk encryption if used
			},
			Properties: &armstorage.AccountPropertiesCreateParameters{
				AllowBlobPublicAccess: to.Ptr(true), // XXX true if using disk encryption
				AllowSharedKeyAccess:  to.Ptr(true),
				IsLocalUserEnabled:    to.Ptr(true),
				LargeFileSharesState:  to.Ptr(armstorage.LargeFileSharesStateEnabled),
				PublicNetworkAccess:   to.Ptr(armstorage.PublicNetworkAccessEnabled),
				MinimumTLSVersion:     &minimumTLSVersion,
			},
			Tags: in.Tags,
		},
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
				PublicAccess: to.Ptr(armstorage.PublicAccessContainer),
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
func CreatePageBlob(ctx context.Context, in *CreatePageBlobInput) (*CreatePageBlobOutput, error) {
	logrus.Debugf("Getting page blob credentials")

	// XXX: Should try all of them until one is successful
	sharedKeyCredential, err := azblob.NewSharedKeyCredential(in.StorageAccountName, *in.StorageAccountKeys[0].Value)
	if err != nil {
		return nil, fmt.Errorf("failed to get shared credentials for storage account: %w", err)
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
		return nil, fmt.Errorf("failed to get page blob client: %w", err)
	}

	// This is used in terraform, not sure if it matters
	metadata := make(map[string]*string, 1)
	metadata["source_uri"] = to.Ptr(in.ImageURL)

	logrus.Debugf("Creating blob")
	_, err = pageBlobClient.Create(ctx, in.ImageLength, &pageblob.CreateOptions{
		Metadata: metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create blob: %w", err)
	}

	logrus.Debugf("Uploading to blob")
	err = doUploadPagesFromURL(ctx, pageBlobClient, in.ImageURL, in.ImageLength)
	if err != nil {
		return nil, fmt.Errorf("failed to upload blob image %s: %w", in.ImageURL, err)
	}
	return &CreatePageBlobOutput{
		PageBlobClient:      pageBlobClient,
		SharedKeyCredential: sharedKeyCredential,
	}, nil
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
