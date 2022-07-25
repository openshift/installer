// Package openstack extracts OpenStack metadata from install
// configurations.
package openstack

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

// PreTerraform performs any infrastructure initialization which must
// happen before Terraform creates the remaining infrastructure.
func PreTerraform(ctx context.Context, clusterID string, installConfig *installconfig.InstallConfig) error {
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
		return errors.Wrapf(err, "unable to determine working directory")
	}

	cloudsYAML := filepath.Join(cwd, "clouds.yaml")
	if _, err = os.Stat(cloudsYAML); err == nil {
		os.Setenv("OS_CLIENT_CONFIG_FILE", cloudsYAML)
	} else if !errors.Is(err, os.ErrNotExist) {
		return errors.Wrapf(err, "unable to determine if clouds.yaml exists")
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
