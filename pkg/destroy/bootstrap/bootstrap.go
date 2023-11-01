// Package bootstrap uses Terraform to remove bootstrap resources.
package bootstrap

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset/cluster"
	openstackasset "github.com/openshift/installer/pkg/asset/cluster/openstack"
	osp "github.com/openshift/installer/pkg/destroy/openstack"
	infra "github.com/openshift/installer/pkg/infrastructure/platform"
	typesazure "github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/openstack"
)

// Destroy uses Terraform to remove bootstrap resources.
func Destroy(dir string) (err error) {
	metadata, err := cluster.LoadMetadata(dir)
	if err != nil {
		return err
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

	provider, err := infra.ProviderForPlatform(platform)
	if err != nil {
		return fmt.Errorf("error getting infrastructure provider: %w", err)
	}

	if err := provider.DestroyBootstrap(dir); err != nil {
		return fmt.Errorf("error destroying bootstrap resources %w", err)
	}

	return nil
}
