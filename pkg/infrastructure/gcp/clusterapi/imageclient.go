package clusterapi

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"google.golang.org/api/compute/v1"

	icgcp "github.com/openshift/installer/pkg/asset/installconfig/gcp"
)

//go:generate mockgen -source=./imageclient.go -destination=./mock/imageclient_generated.go -package=mock

// ImageClient is the GCP API surface required to publish an RHCOS image. It
// exists so that the upload flow can be exercised without reaching GCP.
type ImageClient interface {
	// GetImage returns the named compute image, or nil when it does not exist.
	GetImage(ctx context.Context, projectID, imageName string) (*compute.Image, error)
	// CreateStagingBucket creates the bucket used to stage the image tarball.
	CreateStagingBucket(ctx context.Context, bucketName, projectID, region string, labels map[string]string) error
	// UploadObject writes a local file into a bucket.
	UploadObject(ctx context.Context, bucketName, objectName, filePath string) error
	// UniverseDomain returns the universe domain of the active credentials,
	// which determines the storage hostname in sovereign clouds.
	UniverseDomain(ctx context.Context) (string, error)
	// CreateImage creates a compute image from a staged GCS object.
	CreateImage(ctx context.Context, projectID, imageName, gcsSource string, labels map[string]string) error
	// CleanupStaging removes the staged object and its bucket.
	CleanupStaging(ctx context.Context, bucketName, objectName string) error
}

// gcpImageClient implements ImageClient against the real GCP APIs.
type gcpImageClient struct {
	computeSvc    *compute.Service
	storageClient *storage.Client
}

var _ ImageClient = (*gcpImageClient)(nil)

func (c *gcpImageClient) GetImage(ctx context.Context, projectID, imageName string) (*compute.Image, error) {
	return getExistingImage(ctx, c.computeSvc, projectID, imageName)
}

func (c *gcpImageClient) CreateStagingBucket(ctx context.Context, bucketName, projectID, region string, labels map[string]string) error {
	return createStagingBucket(ctx, c.storageClient, bucketName, projectID, region, labels)
}

func (c *gcpImageClient) UploadObject(ctx context.Context, bucketName, objectName, filePath string) error {
	return uploadToGCS(ctx, c.storageClient, bucketName, objectName, filePath)
}

func (c *gcpImageClient) UniverseDomain(ctx context.Context) (string, error) {
	ssn, err := icgcp.GetSession(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get GCP session: %w", err)
	}
	universeDomain, err := ssn.Credentials.GetUniverseDomain()
	if err != nil {
		return "", fmt.Errorf("failed to get universe domain: %w", err)
	}
	return universeDomain, nil
}

func (c *gcpImageClient) CreateImage(ctx context.Context, projectID, imageName, gcsSource string, labels map[string]string) error {
	return createComputeImage(ctx, c.computeSvc, projectID, imageName, gcsSource, labels)
}

func (c *gcpImageClient) CleanupStaging(ctx context.Context, bucketName, objectName string) error {
	return cleanupStaging(ctx, c.storageClient, bucketName, objectName)
}
