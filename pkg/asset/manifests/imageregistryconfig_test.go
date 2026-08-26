package manifests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	imageregistryv1 "github.com/openshift/api/imageregistry/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

func TestImageRegistryConfigGenerate(t *testing.T) {
	cases := []struct {
		name             string
		installConfig    *types.InstallConfig
		expectFile       bool
		expectedKeyID    string
		expectedRegion   string
		expectedReplicas int32
	}{
		{
			name: "GCP without KMS encryption produces no file",
			installConfig: icBuild.build(
				icBuild.forGCP(),
				icBuild.withGCPRegion("us-central1"),
				icBuild.withGCPProjectID("test-project"),
			),
			expectFile: false,
		},
		{
			name: "GCP with KMS encryption produces image registry config",
			installConfig: icBuild.build(
				icBuild.forGCP(),
				icBuild.withGCPRegion("us-east1"),
				icBuild.withGCPProjectID("test-project"),
				icBuild.withGCPDefaultMachineKMSKey(&gcptypes.KMSKeyReference{
					Name:     "registry-key",
					KeyRing:  "openshift-keyring",
					Location: "us-east1",
				}),
			),
			expectFile:       true,
			expectedKeyID:    "projects/test-project/locations/us-east1/keyRings/openshift-keyring/cryptoKeys/registry-key",
			expectedRegion:   "us-east1",
			expectedReplicas: 2,
		},
		{
			name: "GCP with KMS encryption and custom KMS project ID",
			installConfig: icBuild.build(
				icBuild.forGCP(),
				icBuild.withGCPRegion("europe-west1"),
				icBuild.withGCPProjectID("default-project"),
				icBuild.withGCPDefaultMachineKMSKey(&gcptypes.KMSKeyReference{
					Name:      "registry-key",
					KeyRing:   "custom-keyring",
					Location:  "europe-west1",
					ProjectID: "kms-project",
				}),
			),
			expectFile:       true,
			expectedKeyID:    "projects/kms-project/locations/europe-west1/keyRings/custom-keyring/cryptoKeys/registry-key",
			expectedRegion:   "europe-west1",
			expectedReplicas: 2,
		},
		{
			name: "Non-GCP platform produces no file",
			installConfig: icBuild.build(
				icBuild.forAWS(),
				icBuild.withAWSRegion("us-east-1"),
			),
			expectFile: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(
				installconfig.MakeAsset(tc.installConfig),
			)

			imageRegistryConfig := &ImageRegistryConfig{}
			err := imageRegistryConfig.Generate(context.Background(), parents)
			if !assert.NoError(t, err, "failed to generate asset") {
				return
			}

			if !tc.expectFile {
				assert.Nil(t, imageRegistryConfig.ConfigFile)
				assert.Empty(t, imageRegistryConfig.Files())
				return
			}

			if !assert.NotNil(t, imageRegistryConfig.ConfigFile) {
				return
			}
			assert.Equal(t, imageRegistryConfigFilename, imageRegistryConfig.ConfigFile.Filename)

			var config imageregistryv1.Config
			err = yaml.Unmarshal(imageRegistryConfig.ConfigFile.Data, &config)
			if !assert.NoError(t, err, "failed to unmarshal image registry config") {
				return
			}

			assert.Equal(t, "cluster", config.Name)
			assert.Equal(t, tc.expectedReplicas, config.Spec.Replicas)

			gcs := config.Spec.Storage.GCS
			if !assert.NotNil(t, gcs) {
				return
			}
			assert.Empty(t, gcs.Bucket, "Bucket should be empty to let the operator generate it")
			assert.Equal(t, tc.expectedRegion, gcs.Region)
			assert.Equal(t, tc.expectedKeyID, gcs.KeyID)
		})
	}
}
