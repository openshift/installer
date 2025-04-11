// Package bootstrap uses Terraform to remove bootstrap resources.
package bootstrap

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/asset/cluster/metadata"
	osp "github.com/openshift/installer/pkg/destroy/openstack"
	"github.com/openshift/installer/pkg/infrastructure/openstack/preprovision"
	infra "github.com/openshift/installer/pkg/infrastructure/platform"
	ibmcloudtfvars "github.com/openshift/installer/pkg/tfvars/ibmcloud"
	"github.com/openshift/installer/pkg/types"
	typesazure "github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/featuregates"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/openstack"
)

// Destroy uses Terraform to remove bootstrap resources.
func Destroy(ctx context.Context, dir string) (err error) {
	metadata, err := metadata.Load(dir)
	if err != nil {
		return err
	}

	platform := metadata.Platform()
	if platform == "" {
		return errors.New("no platform configured in metadata")
	}

	if platform == openstack.Name {
		if err := preprovision.SetTerraformEnvironment(); err != nil {
			return errors.Wrapf(err, "Failed to  initialize infrastructure")
		}

		imageName := metadata.InfraID + "-ignition"
		if err := osp.DeleteGlanceImage(ctx, imageName, metadata.OpenStack.Cloud); err != nil {
			return errors.Wrapf(err, "Failed to delete glance image %s", imageName)
		}
	}

	// Azure Stack uses the Azure platform but has its own Terraform configuration.
	if platform == typesazure.Name && metadata.Azure.CloudName == typesazure.StackCloud {
		platform = typesazure.StackTerraformName
	}

	// IBM Cloud allows override of service endpoints, which would be required during bootstrap destroy.
	// Create a JSON file with overrides, if these endpoints are present
	if platform == ibmcloudtypes.Name && metadata.IBMCloud != nil && len(metadata.IBMCloud.ServiceEndpoints) > 0 {
		// Build the JSON containing the endpoint overrides for IBM Cloud Services.
		jsonData, err := ibmcloudtfvars.CreateEndpointJSON(metadata.IBMCloud.ServiceEndpoints, metadata.IBMCloud.Region)
		if err != nil {
			return fmt.Errorf("failed generating endpoint override JSON data for bootstrap destroy: %w", err)
		}

		// If JSON data was generated, create the JSON file for IBM Cloud Terraform provider to use during destroy.
		if jsonData != nil {
			endpointsFilePath := filepath.Join(dir, ibmcloudtfvars.IBMCloudEndpointJSONFileName)
			if err := os.WriteFile(endpointsFilePath, jsonData, 0o600); err != nil {
				return fmt.Errorf("failed to write IBM Cloud service endpoint override JSON file for bootstrap destroy: %w", err)
			}
			logrus.Debugf("generated ibm endpoint overrides file: %s", endpointsFilePath)
		}
	}

	// Get cluster profile for new FeatureGate access.  Blank is no longer an option, so default to
	// SelfManaged.
	clusterProfile := types.GetClusterProfileName()
	featureSets, ok := features.AllFeatureSets()[clusterProfile]
	if !ok {
		return fmt.Errorf("no feature sets for cluster profile %q", clusterProfile)
	}
	fg := featuregates.FeatureGateFromFeatureSets(featureSets, metadata.FeatureSet, metadata.CustomFeatureSet)

	provider, err := infra.ProviderForPlatform(platform, fg)
	if err != nil {
		return fmt.Errorf("error getting infrastructure provider: %w", err)
	}

	if err := provider.DestroyBootstrap(ctx, dir); err != nil {
		return fmt.Errorf("error destroying bootstrap resources %w", err)
	}

	return nil
}
