package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/storage"

	"github.com/openshift/installer/pkg/asset/installconfig"
	gcpic "github.com/openshift/installer/pkg/asset/installconfig/gcp"
)

const (
	BootstrapIgnitionBucket = "bootstrap.ign"
)

// GetBootstrapStorageName gets the name of the storage bucket for the bootstrap process.
func GetBootstrapStorageName(clusterID string) string {
	return fmt.Sprintf("%s-bootstrap-ignition", clusterID)
}

// NewClient creates a new Google storage client.
func NewClient(ctx context.Context) (*storage.Client, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

// CreateBucketHandle will create the bucket handle that can be used as a reference for other storage resources.
func CreateBucketHandle(ctx context.Context, bucketName string) (*storage.BucketHandle, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	client, err := NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.Bucket(bucketName), nil
}

// CreateStorage creates the gcp bucket/storage. The storage bucket does Not include the bucket object. The
// bucket object is created as a separate process/function, so that the two are not tied together, and
// the data stored inside the object can be set at a later time.
func CreateStorage(ctx context.Context, ic *installconfig.InstallConfig, clusterID, bucketName string) error {
	bucketHandle, err := CreateBucketHandle(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to create bucket handle: %w", err)
	}

	labels := map[string]string{}
	labels[fmt.Sprintf("kubernetes-io-cluster-%s", clusterID)] = "owned"
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

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	if err := bucketHandle.Create(ctx, ic.Config.GCP.ProjectID, &bucketAttrs); err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}
	return nil
}

// CreateSignedURL creates a signed url and correlates the signed url with a storage bucket.
func CreateSignedURL(handle *storage.BucketHandle, objectName string) (string, error) {
	opts := storage.SignedURLOptions{
		Method:  "PUT",
		Expires: time.Now().Add(time.Minute * 60),
	}

	ctx := context.Background()
	session, err := gcpic.GetSession(ctx)
	if err != nil {
		return "", err
	}

	// TODO: make sure all cases are handled including the cases required by https://github.com/openshift/installer/pull/7697
	if session.Credentials.JSON != nil {
		var credsMap map[string]interface{}
		if err := json.Unmarshal(session.Credentials.JSON, &credsMap); err != nil {
			return "", err
		}
		opts.GoogleAccessID = credsMap["client_email"].(string)
		opts.PrivateKey = []byte(credsMap["private_key"].(string))
	}

	// The object has not been created yet. This is ok, it is expected to be created after this call.
	// However, if the object is never created this could cause major issues.
	url, err := handle.SignedURL(objectName, &opts)
	if err != nil {
		return "", fmt.Errorf("failed to create a signed url: %w", err)
	}

	return url, nil
}

// ProvisionBootstrapStorage will provision the required storage bucket and signed url for the bootstrap process.
func ProvisionBootstrapStorage(ic *installconfig.InstallConfig, clusterID string) (string, error) {
	ctx := context.Background()

	if err := CreateStorage(ctx, ic, clusterID, BootstrapIgnitionBucket); err != nil {
		return "", nil
	}

	bucketHandle, err := CreateBucketHandle(ctx, GetBootstrapStorageName(clusterID))
	if err != nil {
		return "", err
	}

	url, err := CreateSignedURL(bucketHandle, BootstrapIgnitionBucket)
	if err != nil {
		return "", err
	}

	return url, nil
}

// FillBucket will add the contents to the bootstrap storage bucket object. The bucketName is the
// name of the bucket, and the object name refers to the object that should be filled within the bucket.
func FillBucket(ctx context.Context, bucketName, objectName, contents string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	bucketHandle, err := CreateBucketHandle(ctx, bucketName)
	if err != nil {
		return err
	}

	objWriter := bucketHandle.Object(objectName).NewWriter(ctx)
	if _, err := fmt.Fprintf(objWriter, contents); err != nil {
		return fmt.Errorf("failed to store content in bucket object: %w", err)
	}

	if err := objWriter.Close(); err != nil {
		return fmt.Errorf("failed to close bucket object writer: %w", err)
	}

	return nil
}
