// Package openstack extracts OpenStack metadata from install
// configurations.
package openstack

import (
	"context"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	rhcos_asset "github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/infrastructure/openstack/preprovision"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

// Metadata converts an install configuration to OpenStack metadata.
func Metadata(infraID string, config *types.InstallConfig) *openstack.Metadata {
	return &openstack.Metadata{
		Cloud: config.Platform.OpenStack.Cloud,
		Identifier: map[string]string{
			"openshiftClusterID": infraID,
		},
	}
}

// PreTerraform performs any infrastructure initialization which must
// happen before Terraform creates the remaining infrastructure.
func PreTerraform(ctx context.Context, tfvarsFile *asset.File, installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID, rhcosImage *rhcos_asset.Image) error {
	if !capiutils.IsEnabled(installConfig) {
		if err := preprovision.ReplaceBootstrapIgnitionInTFVars(ctx, tfvarsFile, installConfig, clusterID); err != nil {
			return err
		}

		if err := preprovision.TagVIPPorts(ctx, installConfig, clusterID.InfraID); err != nil {
			return err
		}

		// upload the corresponding image to Glance if rhcosImage contains a
		// URL. If rhcosImage contains a name, then that points to an existing
		// Glance image.
		if imageName, isURL := rhcos.GenerateOpenStackImageName(rhcosImage.ControlPlane, clusterID.InfraID); isURL {
			if err := preprovision.UploadBaseImage(ctx, installConfig.Config.Platform.OpenStack.Cloud, rhcosImage.ControlPlane, imageName, clusterID.InfraID, installConfig.Config.Platform.OpenStack.ClusterOSImageProperties); err != nil {
				return err
			}
		}
	}

	return preprovision.SetTerraformEnvironment()
}
