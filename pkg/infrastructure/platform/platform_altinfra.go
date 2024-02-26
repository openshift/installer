//go:build altinfra
// +build altinfra

package platform

import (
	"fmt"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/infrastructure/aws"
	awscapi "github.com/openshift/installer/pkg/infrastructure/aws/clusterapi"
	azurecapi "github.com/openshift/installer/pkg/infrastructure/azure"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	gcpcapi "github.com/openshift/installer/pkg/infrastructure/gcp/clusterapi"
	vspherecapi "github.com/openshift/installer/pkg/infrastructure/vsphere/clusterapi"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/featuregates"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// ProviderForPlatform returns the stages to run to provision the infrastructure for the specified platform.
func ProviderForPlatform(platform string, fg featuregates.FeatureGate) (infrastructure.Provider, error) {
	switch platform {
	case awstypes.Name:
		if fg.Enabled(configv1.FeatureGateClusterAPIInstall) {
			return clusterapi.InitializeProvider(&awscapi.Provider{}), nil
		}
		return aws.InitializeProvider(), nil
	case azuretypes.Name:
		if fg.Enabled(configv1.FeatureGateClusterAPIInstall) {
			return clusterapi.InitializeProvider(&azurecapi.Provider{}), nil
		}
		return nil, nil
	case gcptypes.Name:
		if fg.Enabled(configv1.FeatureGateClusterAPIInstall) {
			return clusterapi.InitializeProvider(gcpcapi.Provider{}), nil
		}
		return nil, nil
	case vspheretypes.Name:
		if fg.Enabled(configv1.FeatureGateClusterAPIInstall) {
			return clusterapi.InitializeProvider(vspherecapi.Provider{}), nil
		}
	}
	return nil, fmt.Errorf("platform %q is not supported in the altinfra Installer build", platform)
}
