// Package openstack extracts OpenStack metadata from install
// configurations.
package openstack

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	rhcos_asset "github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

// SetTerraformEnvironment injects the environment variables required
// authenticate against the OpenStack API.
func SetTerraformEnvironment() error {
	// Terraform runs in a different directory but we want to allow people to
	// use clouds.yaml files in their local directory. Emulate this by setting
	// the necessary environment variable to point to this file if (a) the user
	// hasn't already set this environment variable and (b) there is actually
	// a local file
	if path := os.Getenv("OS_CLIENT_CONFIG_FILE"); path != "" {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to determine working directory: %w", err)
	}

	cloudsYAML := filepath.Join(cwd, "clouds.yaml")
	if _, err = os.Stat(cloudsYAML); err == nil {
		os.Setenv("OS_CLIENT_CONFIG_FILE", cloudsYAML)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("unable to determine if clouds.yaml exists: %w", err)
	}
	return nil
}

// PreTerraform performs any infrastructure initialization which must
// happen before Terraform creates the remaining infrastructure.
func PreTerraform(ctx context.Context, tfvarsFile *asset.File, installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID, rhcosImage *rhcos_asset.Image) error {
	if err := replaceBootstrapIgnitionInTFVars(ctx, tfvarsFile, installConfig, clusterID); err != nil {
		return err
	}

	// upload the corresponding image to Glance if rhcosImage contains a
	// URL. If rhcosImage contains a name, then that points to an existing
	// Glance image.
	if imageName, isURL := rhcos.GenerateOpenStackImageName(string(*rhcosImage), clusterID.InfraID); isURL {
		if err := uploadBaseImage(installConfig.Config.Platform.OpenStack.Cloud, rhcosImage, imageName, clusterID, installConfig.Config.Platform.OpenStack.ClusterOSImageProperties); err != nil {
			return err
		}
	}

	return SetTerraformEnvironment()
}

// PreCAPI performs any infrastructure initialization which must
// happen before CAPI creates the remaining infrastructure.
func PreCAPI(ctx context.Context, manifests *[]client.Object, installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID, rhcosImage *rhcos_asset.Image) error {
	if err := replaceBootstrapIgnitionInManifests(ctx, manifests, installConfig, clusterID); err != nil {
		return err
	}

	// upload the corresponding image to Glance if rhcosImage contains a
	// URL. If rhcosImage contains a name, then that points to an existing
	// Glance image.
	if imageName, isURL := rhcos.GenerateOpenStackImageName(string(*rhcosImage), clusterID.InfraID); isURL {
		if err := uploadBaseImage(installConfig.Config.Platform.OpenStack.Cloud, rhcosImage, imageName, clusterID, installConfig.Config.Platform.OpenStack.ClusterOSImageProperties); err != nil {
			return err
		}
	}

	return nil
}

// Metadata converts an install configuration to OpenStack metadata.
func Metadata(infraID string, config *types.InstallConfig) *openstack.Metadata {
	return &openstack.Metadata{
		Cloud: config.Platform.OpenStack.Cloud,
		Identifier: map[string]string{
			"openshiftClusterID": infraID,
		},
	}
}
