package external

import (
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// GetInfraPlatformSpec constructs ExternalPlatformSpec for the infrastructure spec.
func GetInfraPlatformSpec(ic *installconfig.InstallConfig) *configv1.ExternalPlatformSpec {
	icPlatformSpec := ic.Config.External

	return &configv1.ExternalPlatformSpec{
		PlatformName: icPlatformSpec.PlatformName,
	}
}

// GetInfraPlatformStatus constructs ExternalPlatformSpec for the infrastructure spec.
func GetInfraPlatformStatus() *configv1.ExternalPlatformStatus {
	return &configv1.ExternalPlatformStatus{
		CloudControllerManager: configv1.CloudControllerManagerStatus{
			State: configv1.CloudControllerManagerExternal,
		},
	}
}
