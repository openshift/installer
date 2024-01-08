// Package bootstrap uses Terraform to remove bootstrap resources.
package bootstrap

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/cluster/metadata"
	openstackasset "github.com/openshift/installer/pkg/asset/cluster/openstack"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/clusterapi"
	osp "github.com/openshift/installer/pkg/destroy/openstack"
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

	if sys := clusterapi.System(); sys.State() == clusterapi.SystemStateRunning {
		return destroyBoostrapMachine(ctx, sys.Client(), metadata)
	}

	platform := metadata.Platform()
	if platform == "" {
		return errors.New("no platform configured in metadata")
	}

	if platform == openstack.Name {
		if err := openstackasset.PreTerraform(); err != nil {
			return errors.Wrapf(err, "Failed to  initialize infrastructure")
		}

		imageName := metadata.InfraID + "-ignition"
		if err := osp.DeleteGlanceImage(imageName, metadata.OpenStack.Cloud); err != nil {
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

	fg := featuregates.FeatureGateFromFeatureSets(configv1.FeatureSets, metadata.FeatureSet, metadata.CustomFeatureSet)

	provider, err := infra.ProviderForPlatform(platform, fg)
	if err != nil {
		return fmt.Errorf("error getting infrastructure provider: %w", err)
	}

	if err := provider.DestroyBootstrap(dir); err != nil {
		return fmt.Errorf("error destroying bootstrap resources %w", err)
	}

	return nil
}

func destroyBoostrapMachine(ctx context.Context, c client.Client, metadata *types.ClusterMetadata) error {
	if err := c.Delete(ctx, &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      capiutils.GenerateBoostrapMachineName(metadata.InfraID),
			Namespace: capiutils.Namespace,
		},
	}); client.IgnoreNotFound(err) != nil {
		return fmt.Errorf("failed to delete bootstrap machine: %w", err)
	}
	return nil
}
