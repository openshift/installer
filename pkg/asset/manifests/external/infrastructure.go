package external

import (
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types/external"
)

// GetInfraPlatformSpec constructs ExternalPlatformSpec for the infrastructure spec.
func GetInfraPlatformSpec(ic *installconfig.InstallConfig) *configv1.ExternalPlatformSpec {
	icPlatformSpec := ic.Config.External

	return &configv1.ExternalPlatformSpec{
		PlatformName: icPlatformSpec.PlatformName,
	}
}

// GetInfraPlatformStatus constructs ExternalPlatformSpec for the infrastructure spec.
func GetInfraPlatformStatus(ic *installconfig.InstallConfig) *configv1.ExternalPlatformStatus {
	icPlatformSpec := ic.Config.External

	var ccmState configv1.CloudControllerManagerState

	switch icPlatformSpec.CloudControllerManager {
	case external.CloudControllerManagerTypeExternal:
		ccmState = configv1.CloudControllerManagerExternal
	default:
		ccmState = configv1.CloudControllerManagerNone
	}

	return &configv1.ExternalPlatformStatus{
		CloudControllerManager: configv1.CloudControllerManagerStatus{
			State: ccmState,
		},
	}
}
