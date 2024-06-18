package gcp

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"github.com/openshift/installer/pkg/asset/installconfig"
	gcpic "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	gcpconsts "github.com/openshift/installer/pkg/constants/gcp"
)

const (
	bootstrapIgnitionBucketObjName = "bootstrap.ign"
)

// GetBootstrapStorageName gets the name of the storage bucket for the bootstrap process.
func GetBootstrapStorageName(clusterID string) string {
	return fmt.Sprintf("%s-bootstrap-ignition", clusterID)
}

// NewStorageClient creates a new Google storage client.
func NewStorageClient(ctx context.Context) (*storage.Client, error) {
	ssn, err := gcpic.GetSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get session while creating gcp storage client: %w", err)
	}

	client, err := storage.NewClient(ctx, option.WithCredentials(ssn.Credentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

// CreateBucketHandle will create the bucket handle that can be used as a reference for other storage resources.
func CreateBucketHandle(ctx context.Context, bucketName string) (*storage.BucketHandle, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	client, err := NewStorageClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}
	return client.Bucket(bucketName), nil
}

// CreateStorage creates the gcp bucket/storage. The storage bucket does Not include the bucket object. The
// bucket object is created as a separate process/function, so that the two are not tied together, and
// the data stored inside the object can be set at a later time.
func CreateStorage(ctx context.Context, ic *installconfig.InstallConfig, bucketHandle *storage.BucketHandle, clusterID string) error {
	labels := map[string]string{}
	labels[fmt.Sprintf(gcpconsts.ClusterIDLabelFmt, clusterID)] = "owned"
	for _, label := range ic.Config.GCP.UserLabels {
		labels[label.Key] = label.Value
	}

	bucketAttrs := storage.BucketAttrs{
		UniformBucketLevelAccess: storage.UniformBucketLevelAccess{
			Enabled: true,
		},
		Location: ic.Config.GCP.Region,
		Labels:   labels,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	if err := bucketHandle.Create(ctx, ic.Config.GCP.ProjectID, &bucketAttrs); err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}
	return nil
}

// CreateSignedURL creates a signed url and correlates the signed url with a storage bucket.
func CreateSignedURL(clusterID string) (string, error) {
	bucketName := GetBootstrapStorageName(clusterID)
	handle, err := CreateBucketHandle(context.Background(), bucketName)
	if err != nil {
		return "", fmt.Errorf("creating presigned url, failed to create bucket handle: %w", err)
	}

	// Signing a URL requires credentials authorized to sign a URL. You can pass
	// these in through SignedURLOptions with a Google Access ID with
	// iam.serviceAccounts.signBlob permissions.
	opts := storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(time.Minute * 60),
	}

	// The object has not been created yet. This is ok, it is expected to be created after this call.
	// However, if the object is never created this could cause major issues.
	url, err := handle.SignedURL(bootstrapIgnitionBucketObjName, &opts)
	if err != nil {
		return "", fmt.Errorf("failed to create a signed url: %w", err)
	}

	return url, nil
}

// FillBucket will add the contents to the bootstrap storage bucket object.
func FillBucket(ctx context.Context, bucketHandle *storage.BucketHandle, contents string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	objWriter := bucketHandle.Object(bootstrapIgnitionBucketObjName).NewWriter(ctx)
	if _, err := fmt.Fprint(objWriter, contents); err != nil {
		return fmt.Errorf("failed to store content in bucket object: %w", err)
	}

	if err := objWriter.Close(); err != nil {
		return fmt.Errorf("failed to close bucket object writer: %w", err)
	}

	return nil
}

// DestroyStorage Destroy the bucket and the bucket objects that are associated with the bucket.
func DestroyStorage(ctx context.Context, clusterID string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	client, err := NewStorageClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %w", err)
	}
	bucketName := GetBootstrapStorageName(clusterID)

	if err := client.Bucket(bucketName).Object(bootstrapIgnitionBucketObjName).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete bucket object %s: %w", bootstrapIgnitionBucketObjName, err)
	}

	// Deleting a bucket will delete the managed folders and bucket objects.
	if err := client.Bucket(bucketName).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete bucket %s: %w", bucketName, err)
	}
	return nil
}
