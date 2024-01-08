//go:build altinfra
// +build altinfra

package platform

import (
	"fmt"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/infrastructure/aws"
	awsinfra "github.com/openshift/installer/pkg/infrastructure/aws"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/featuregates"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// ProviderForPlatform returns the stages to run to provision the infrastructure for the specified platform.
func ProviderForPlatform(platform string, fg featuregates.FeatureGate) (infrastructure.Provider, error) {
	switch platform {
	case awstypes.Name:
		if fg.Enabled(configv1.FeatureGateClusterAPIInstall) {
			return clusterapi.InitializeProvider(awsinfra.CAPIInfraHelper{}), nil
		}
		return aws.InitializeProvider(), nil
	case azuretypes.Name:
		panic("not implemented")
		return nil, nil
	case vspheretypes.Name:
		panic("not implemented")
		return nil, nil
	}
	return nil, fmt.Errorf("platform %q is not supported in the altinfra Installer build", platform)
}
