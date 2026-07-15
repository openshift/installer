package clusterapi

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"github.com/thedevsaddam/retry"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"k8s.io/utils/ptr"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8syaml "sigs.k8s.io/yaml"

	machineapi "github.com/openshift/api/machine/v1beta1"
	icgcp "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/asset/machines"
	gcpconsts "github.com/openshift/installer/pkg/constants/gcp"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

const (
	rhcosImageObjectName = "rhcos-image.tar.gz"
)

// uploadRHCOSImage downloads the RHCOS image from the stream, uploads it to a GCS
// staging bucket, creates a GCP compute image from it, then cleans up the staging
// resources. Returns the full image reference for use in machine manifests.
func uploadRHCOSImage(ctx context.Context, in clusterapi.PreProvisionInput) (string, error) {
	imageURL := in.RhcosImage.ControlPlane
	platform := in.InstallConfig.Config.Platform.GCP
	projectID := platform.ProjectID
	region := platform.Region
	imageName := fmt.Sprintf("%s-rhcos", in.InfraID)
	bucketName := fmt.Sprintf("%s-rhcos-image", in.InfraID)

	logrus.Infof("Uploading RHCOS image for cluster %s", in.InfraID)

	// Parse download URL and extract sha256 checksum
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse RHCOS image URL: %w", err)
	}
	sha256Checksum := parsedURL.Query().Get("sha256")
	parsedURL.RawQuery = ""
	downloadURL := parsedURL.String()

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

	// Download the RHCOS tar.gz image
	localPath, err := downloadRHCOSImage(downloadURL, sha256Checksum)
	if err != nil {
		return "", fmt.Errorf("failed to download RHCOS image: %w", err)
	}
	defer os.Remove(localPath)

	// Create staging bucket
	labels := buildImageLabels(in.InstallConfig.Config.GCP, in.InfraID)
	if err := createStagingBucket(ctx, storageClient, bucketName, projectID, region, labels); err != nil {
		return "", fmt.Errorf("failed to create staging bucket: %w", err)
	}

	// Upload image to GCS
	if err := uploadToGCS(ctx, storageClient, bucketName, rhcosImageObjectName, localPath); err != nil {
		return "", fmt.Errorf("failed to upload RHCOS image to GCS: %w", err)
	}

	// Build GCS URL using the universe domain from credentials so that
	// sovereign clouds get the correct storage hostname.
	ssn, err := icgcp.GetSession(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get GCP session: %w", err)
	}
	universeDomain, err := ssn.Credentials.GetUniverseDomain()
	if err != nil {
		return "", fmt.Errorf("failed to get universe domain: %w", err)
	}
	gcsURI := fmt.Sprintf("https://storage.%s/%s/%s", universeDomain, bucketName, rhcosImageObjectName)
	if err := createComputeImage(ctx, computeSvc, projectID, imageName, gcsURI, labels); err != nil {
		return "", fmt.Errorf("failed to create compute image: %w", err)
	}

	// Clean up staging resources
	if err := cleanupStaging(ctx, storageClient, bucketName, rhcosImageObjectName); err != nil {
		logrus.Warnf("Failed to clean up staging bucket %s: %v (will be cleaned up during destroy)", bucketName, err)
	}

	imageRef := fmt.Sprintf("projects/%s/global/images/%s", projectID, imageName)
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
	err = retry.DoFunc(3, 5*time.Second, func() error {
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

	if err := client.Bucket(bucketName).Object(objectName).Delete(ctx); err != nil {
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

// needsImageUpload checks whether any GCPMachine manifest has an empty image,
// indicating that the RHCOS image needs to be uploaded to the cluster's project.
func needsImageUpload(machineManifests []client.Object) bool {
	for i := range machineManifests {
		if gcpMachine, ok := machineManifests[i].(*capg.GCPMachine); ok {
			if ptr.Deref(gcpMachine.Spec.Image, "") == "" {
				return true
			}
		}
	}
	return false
}

// updateMachineManifestImages sets the image reference on all GCPMachine manifests.
func updateMachineManifestImages(machineManifests []client.Object, imageRef string) {
	for i := range machineManifests {
		if gcpMachine, ok := machineManifests[i].(*capg.GCPMachine); ok {
			gcpMachine.Spec.Image = ptr.To(imageRef)
		}
	}
}

// updateWorkerMachineSetImages updates the RHCOS image reference in worker MachineSet
// provider specs. Worker MachineSets use the MAPI GCPMachineProviderSpec format,
// stored as serialized YAML in the Worker asset's MachineSetFiles.
func updateWorkerMachineSetImages(workersAsset *machines.Worker, imageRef string) error {
	machineSets, err := workersAsset.MachineSets()
	if err != nil {
		return fmt.Errorf("failed to get worker machinesets: %w", err)
	}

	for i := range machineSets {
		providerSpec, ok := machineSets[i].Spec.Template.Spec.ProviderSpec.Value.Object.(*machineapi.GCPMachineProviderSpec)
		if !ok {
			continue
		}

		updated := false
		for d := range providerSpec.Disks {
			if providerSpec.Disks[d].Image == "" || strings.HasPrefix(providerSpec.Disks[d].Image, "http") {
				providerSpec.Disks[d].Image = imageRef
				updated = true
			}
		}

		if updated {
			rawProviderSpec, err := json.Marshal(providerSpec)
			if err != nil {
				return fmt.Errorf("failed to marshal updated provider spec for machineset %d: %w", i, err)
			}
			machineSets[i].Spec.Template.Spec.ProviderSpec.Value.Raw = rawProviderSpec
			machineSets[i].Spec.Template.Spec.ProviderSpec.Value.Object = nil

			data, err := k8syaml.Marshal(&machineSets[i])
			if err != nil {
				return fmt.Errorf("failed to marshal updated machineset %d: %w", i, err)
			}
			workersAsset.MachineSetFiles[i].Data = data
		}
	}

	return nil
}
