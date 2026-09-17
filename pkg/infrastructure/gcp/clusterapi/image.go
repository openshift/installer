package clusterapi

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"github.com/thedevsaddam/retry"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"

	icgcp "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	gcpconsts "github.com/openshift/installer/pkg/constants/gcp"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

const (
	rhcosImageObjectName = "rhcos-image.tar.gz"

	// imageStatusReady is the status of a compute image that is usable by machines.
	imageStatusReady = "READY"
)

// Bounds on retrying a failed RHCOS download. The artifact is served by a
// mirror, so a transient failure is worth a few attempts before giving up on
// the install. Overridden in tests.
var (
	downloadRetryCount uint = 3
	downloadRetryDelay      = 5 * time.Second
)

// getExistingImage returns the named compute image, or nil when it does not
// exist. Any other API failure is returned as an error, so that a transient
// lookup failure is not mistaken for a missing image.
func getExistingImage(ctx context.Context, svc *compute.Service, projectID, imageName string) (*compute.Image, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	image, err := svc.Images.Get(projectID, imageName).Context(ctx).Do()
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to check for existing image %s: %w", imageName, err)
	}
	return image, nil
}

// imageDownloader fetches an RHCOS artifact and returns the path of the local
// file it wrote. Injected so that the upload flow can be tested without
// downloading gigabytes over the network.
type imageDownloader func(imageURL, sha256Checksum string) (string, error)

// uploadRHCOSImage builds the GCP clients for the target project and publishes a
// cluster-specific RHCOS image. See publishRHCOSImage for the flow itself.
func uploadRHCOSImage(ctx context.Context, in clusterapi.PreProvisionInput) (string, error) {
	platform := in.InstallConfig.Config.Platform.GCP

	// Create GCP clients with endpoint options for sovereign clouds
	computeOpts := []option.ClientOption{}
	storageOpts := []option.ClientOption{}
	if gcptypes.ShouldUseEndpointForInstaller(platform.Endpoint) {
		computeOpts = append(computeOpts, icgcp.CreateEndpointOption(platform.Endpoint.Name, icgcp.ServiceNameGCPCompute))
		storageOpts = append(storageOpts, icgcp.CreateEndpointOption(platform.Endpoint.Name, icgcp.ServiceNameGCPStorage))
	}

	computeSvc, err := icgcp.GetComputeService(ctx, computeOpts...)
	if err != nil {
		return "", fmt.Errorf("failed to create compute service: %w", err)
	}

	storageClient, err := icgcp.GetStorageService(ctx, storageOpts...)
	if err != nil {
		return "", fmt.Errorf("failed to create storage service: %w", err)
	}

	client := &gcpImageClient{computeSvc: computeSvc, storageClient: storageClient}
	return publishRHCOSImage(ctx, in, client, downloadRHCOSImage)
}

// publishRHCOSImage downloads the RHCOS image from the stream, uploads it to a GCS
// staging bucket, creates a GCP compute image from it, then cleans up the staging
// resources. Returns the full image reference for use in machine manifests.
//
// A single image is created and shared by the control plane and compute machines.
// This relies on both pools having the same architecture: sovereign clouds do not
// offer arm64 instances, so a heterogeneous cluster cannot be provisioned there.
// If arm64 becomes available in a sovereign region, this must upload one image per
// architecture, since in.RhcosImage.ControlPlane and in.RhcosImage.Compute are
// resolved per pool and would then point at different artifacts.
func publishRHCOSImage(ctx context.Context, in clusterapi.PreProvisionInput, client ImageClient, download imageDownloader) (string, error) {
	imageURL := in.RhcosImage.ControlPlane
	platform := in.InstallConfig.Config.Platform.GCP
	projectID := platform.ProjectID
	region := platform.Region
	imageName := fmt.Sprintf("%s-rhcos", in.InfraID)
	bucketName := fmt.Sprintf("%s-rhcos-image", in.InfraID)

	// Guard the shared-image assumption described above. Nothing in install-config
	// validation rejects heterogeneous architectures on GCP when a multi-arch
	// payload is used, so fail loudly here rather than booting compute machines
	// from a control plane image.
	if in.RhcosImage.Compute != imageURL {
		return "", fmt.Errorf("heterogeneous architectures are not supported on sovereign clouds: "+
			"control plane and compute resolve to different RHCOS artifacts (%s, %s)",
			imageURL, in.RhcosImage.Compute)
	}

	logrus.Infof("Uploading RHCOS image for cluster %s", in.InfraID)

	// Parse download URL and extract sha256 checksum
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse RHCOS image URL: %w", err)
	}
	sha256Checksum := parsedURL.Query().Get("sha256")
	parsedURL.RawQuery = ""
	downloadURL := parsedURL.String()

	// Reuse an image left by an earlier attempt. PreProvision can run again
	// after a partial failure, and without this the retry re-downloads
	// gigabytes of image data only for Images.Insert to fail with alreadyExists.
	existing, err := client.GetImage(ctx, projectID, imageName)
	if err != nil {
		return "", err
	}
	if existing != nil {
		if existing.Status != imageStatusReady {
			return "", fmt.Errorf("image %s already exists in project %s but is not usable (status %s); "+
				"delete it and retry the install", imageName, projectID, existing.Status)
		}
		imageRef := gcptypes.RHCOSImageRef(projectID, in.InfraID)
		logrus.Infof("Reusing existing RHCOS image: %s", imageRef)
		return imageRef, nil
	}

	// Download the RHCOS tar.gz image
	localPath, err := download(downloadURL, sha256Checksum)
	if err != nil {
		return "", fmt.Errorf("failed to download RHCOS image: %w", err)
	}
	// downloadRHCOSImage owns the enclosing temp directory, so remove the
	// directory rather than just the file it holds.
	defer os.RemoveAll(filepath.Dir(localPath))

	// Create staging bucket
	labels := buildImageLabels(in.InstallConfig.Config.GCP, in.InfraID)
	if err := client.CreateStagingBucket(ctx, bucketName, projectID, region, labels); err != nil {
		return "", fmt.Errorf("failed to create staging bucket: %w", err)
	}

	// Tear the staging resources down on every exit path, not just success:
	// they are pure scratch space and are useless once this function returns.
	// The context is detached because cleanup must still run when the install
	// context is already cancelled, which is exactly when we are bailing out.
	defer func() {
		if err := client.CleanupStaging(context.WithoutCancel(ctx), bucketName, rhcosImageObjectName); err != nil {
			logrus.Warnf("Failed to clean up staging bucket %s: %v (will be cleaned up during destroy)", bucketName, err)
		}
	}()

	// Upload image to GCS
	if err := client.UploadObject(ctx, bucketName, rhcosImageObjectName, localPath); err != nil {
		return "", fmt.Errorf("failed to upload RHCOS image to GCS: %w", err)
	}

	// Build GCS URL using the universe domain from credentials so that
	// sovereign clouds get the correct storage hostname.
	universeDomain, err := client.UniverseDomain(ctx)
	if err != nil {
		return "", err
	}
	gcsURI := fmt.Sprintf("https://storage.%s/%s/%s", universeDomain, bucketName, rhcosImageObjectName)
	if err := client.CreateImage(ctx, projectID, imageName, gcsURI, labels); err != nil {
		return "", fmt.Errorf("failed to create compute image: %w", err)
	}

	imageRef := gcptypes.RHCOSImageRef(projectID, in.InfraID)
	logrus.Infof("RHCOS image created: %s", imageRef)
	return imageRef, nil
}

// downloadRHCOSImage downloads the RHCOS tar.gz file without decompressing it.
// GCP compute image creation requires the original tar.gz format.
func downloadRHCOSImage(imageURL string, sha256Checksum string) (string, error) {
	logrus.Infof("Downloading RHCOS image from %s", imageURL)

	tmpDir, err := os.MkdirTemp("", "rhcos-gcp-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	filePath := filepath.Join(tmpDir, "rhcos-gcp.tar.gz")

	httpClient := &http.Client{Timeout: 30 * time.Minute}
	err = retry.DoFunc(downloadRetryCount, downloadRetryDelay, func() error {
		resp, err := httpClient.Get(imageURL) //nolint:gosec
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("bad status downloading RHCOS image: %s", resp.Status)
		}

		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		var reader io.Reader = resp.Body
		hasher := sha256.New()
		if sha256Checksum != "" {
			reader = io.TeeReader(resp.Body, hasher)
		}

		written, err := io.Copy(f, reader)
		if err != nil {
			return fmt.Errorf("failed to write RHCOS image: %w", err)
		}
		logrus.Debugf("Downloaded RHCOS image: %d bytes", written)

		if sha256Checksum != "" {
			foundChecksum := fmt.Sprintf("%x", hasher.Sum(nil))
			if sha256Checksum != foundChecksum {
				return fmt.Errorf("checksum mismatch for RHCOS image: expected=%s found=%s", sha256Checksum, foundChecksum)
			}
			logrus.Debug("RHCOS image checksum verification passed")
		}

		return nil
	})
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	return filePath, nil
}

// createStagingBucket creates a GCS bucket for staging the RHCOS image.
func createStagingBucket(ctx context.Context, client *storage.Client, bucketName, projectID, region string, labels map[string]string) error {
	logrus.Infof("Creating staging bucket %s", bucketName)
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	bucketAttrs := storage.BucketAttrs{
		UniformBucketLevelAccess: storage.UniformBucketLevelAccess{
			Enabled: true,
		},
		Location: region,
		Labels:   labels,
	}

	if err := client.Bucket(bucketName).Create(ctx, projectID, &bucketAttrs); err != nil {
		// An earlier attempt that failed before cleanup leaves the bucket
		// behind. The object write below overwrites whatever it holds, so
		// reusing it is safe and keeps a retried PreProvision working.
		if isConflictError(err) {
			logrus.Debugf("Reusing existing staging bucket %s", bucketName)
			return nil
		}
		return fmt.Errorf("failed to create staging bucket: %w", err)
	}
	return nil
}

// uploadToGCS uploads a local file to a GCS bucket.
func uploadToGCS(ctx context.Context, client *storage.Client, bucketName, objectName, filePath string) error {
	logrus.Infof("Uploading %s to gs://%s/%s", filePath, bucketName, objectName)

	f, err := os.Open(filepath.Clean(filePath)) //nolint:gosec // filePath is from our own downloadRHCOSImage, not user input
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	writer := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err := io.Copy(writer, f); err != nil {
		return fmt.Errorf("failed to upload to GCS: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to finalize GCS upload: %w", err)
	}

	logrus.Debug("RHCOS image uploaded to GCS")
	return nil
}

// createComputeImage creates a GCP compute image from a GCS object.
func createComputeImage(ctx context.Context, svc *compute.Service, projectID, imageName, gcsSource string, labels map[string]string) error {
	logrus.Infof("Creating GCP compute image %s", imageName)

	// Building an image from a multi-gigabyte tarball routinely takes several
	// minutes, so this is far more generous than the other staging steps. It
	// still needs a bound: without one a stalled operation blocks the install
	// indefinitely.
	ctx, cancel := context.WithTimeout(ctx, time.Minute*15)
	defer cancel()

	image := &compute.Image{
		Name: imageName,
		RawDisk: &compute.ImageRawDisk{
			Source: gcsSource,
		},
		GuestOsFeatures: []*compute.GuestOsFeature{
			{Type: "GVNIC"},
			{Type: "UEFI_COMPATIBLE"},
			{Type: "VIRTIO_SCSI_MULTIQUEUE"},
			{Type: "SEV_CAPABLE"},
			{Type: "SEV_SNP_CAPABLE"},
		},
		Labels: labels,
	}

	op, err := svc.Images.Insert(projectID, image).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to insert image: %w", err)
	}

	if err := WaitForOperationGlobal(ctx, svc, projectID, op); err != nil {
		return fmt.Errorf("failed waiting for image creation: %w", err)
	}

	logrus.Infof("GCP compute image %s created successfully", imageName)
	return nil
}

// cleanupStaging removes the staging GCS object and bucket.
func cleanupStaging(ctx context.Context, client *storage.Client, bucketName, objectName string) error {
	logrus.Infof("Cleaning up staging bucket %s", bucketName)
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	// Cleanup now also runs after a failed upload, where the object was never
	// written. Treating that as an error would skip the bucket delete below and
	// leak the bucket.
	if err := client.Bucket(bucketName).Object(objectName).Delete(ctx); err != nil && !errors.Is(err, storage.ErrObjectNotExist) {
		return fmt.Errorf("failed to delete staging object: %w", err)
	}
	if err := client.Bucket(bucketName).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete staging bucket: %w", err)
	}
	return nil
}

// buildImageLabels creates labels for the RHCOS compute image and staging bucket.
func buildImageLabels(platform *gcptypes.Platform, infraID string) map[string]string {
	labels := map[string]string{
		fmt.Sprintf(gcpconsts.ClusterIDLabelFmt, infraID): "owned",
	}
	for _, label := range platform.UserLabels {
		labels[label.Key] = label.Value
	}
	return labels
}
