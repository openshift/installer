package gcp

import (
	imageregistryv1 "github.com/openshift/api/imageregistry/v1"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

// BuildGCSConfig constructs the GCS storage configuration for the ImageRegistry CR.
// The bucket name is intentionally left empty so that the cluster-image-registry-operator
// generates one with proper length truncation and padding.
func BuildGCSConfig(platform *gcptypes.Platform) *imageregistryv1.ImageRegistryConfigStorageGCS {
	gcsConfig := &imageregistryv1.ImageRegistryConfigStorageGCS{
		Region:    platform.Region,
		ProjectID: platform.ProjectID,
	}

	if kmsKey := gcptypes.GetDefaultEncryptionKey(platform); kmsKey != nil {
		gcsConfig.KeyID = gcptypes.FormatKMSKeyResourcePath(kmsKey, platform.ProjectID)
	}

	return gcsConfig
}
