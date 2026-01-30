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
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

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

	// Clean up bootstrap-only cluster resources (e.g., Konnectivity agent DaemonSet)
	// This runs after infrastructure is destroyed, so failures are warnings only.
	if err := deleteBootstrapClusterResources(ctx, dir); err != nil {
		logrus.Warnf("Failed to clean up bootstrap cluster resources: %v", err)
	}

	return nil
}

// deleteBootstrapClusterResources removes cluster resources that were only needed
// during bootstrap, such as the Konnectivity agent DaemonSet.
func deleteBootstrapClusterResources(ctx context.Context, dir string) error {
	kubeconfigPath := filepath.Join(dir, "auth", "kubeconfig")

	// Check if kubeconfig exists
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		logrus.Debug("Kubeconfig not found, skipping bootstrap cluster resource cleanup")
		return nil
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Delete the Konnectivity agent DaemonSet
	logrus.Info("Deleting bootstrap Konnectivity agent DaemonSet")
	err = client.AppsV1().DaemonSets("kube-system").Delete(ctx, "konnectivity-agent", metav1.DeleteOptions{})
	if err != nil {
		// Log but don't fail if the DaemonSet doesn't exist or can't be deleted
		logrus.Debugf("Failed to delete konnectivity-agent DaemonSet: %v", err)
	}

	// Delete the Konnectivity agent certificate Secret
	logrus.Info("Deleting bootstrap Konnectivity agent certificate Secret")
	err = client.CoreV1().Secrets("kube-system").Delete(ctx, "konnectivity-agent-certs", metav1.DeleteOptions{})
	if err != nil {
		// Log but don't fail if the Secret doesn't exist or can't be deleted
		logrus.Debugf("Failed to delete konnectivity-agent-certs Secret: %v", err)
	}

	return nil
}
