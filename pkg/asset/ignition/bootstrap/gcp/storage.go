package gcp

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"

	"github.com/openshift/installer/pkg/asset/installconfig"
	gcpconsts "github.com/openshift/installer/pkg/constants/gcp"
)

const (
	bootstrapIgnitionBucketObjName = "bootstrap.ign"
)

// GetBootstrapStorageName gets the name of the storage bucket for the bootstrap process.
func GetBootstrapStorageName(clusterID string) string {
	return fmt.Sprintf("%s-bootstrap-ignition", clusterID)
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
func CreateSignedURL(client *storage.Client, clusterID string) (string, error) {
	bucketName := GetBootstrapStorageName(clusterID)

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
	url, err := client.Bucket(bucketName).SignedURL(bootstrapIgnitionBucketObjName, &opts)
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
func DestroyStorage(ctx context.Context, client *storage.Client, clusterID string) error {
	bucketName := GetBootstrapStorageName(clusterID)

	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()
	if err := client.Bucket(bucketName).Object(bootstrapIgnitionBucketObjName).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete bucket object %s: %w", bootstrapIgnitionBucketObjName, err)
	}

	// Deleting a bucket will delete the managed folders and bucket objects.
	if err := client.Bucket(bucketName).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete bucket %s: %w", bucketName, err)
	}
	return nil
}
