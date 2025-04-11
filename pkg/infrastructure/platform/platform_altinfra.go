//go:build altinfra
// +build altinfra

package platform

import (
	"fmt"

	"github.com/openshift/installer/pkg/infrastructure"
	awscapi "github.com/openshift/installer/pkg/infrastructure/aws/clusterapi"
	azurecapi "github.com/openshift/installer/pkg/infrastructure/azure"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	gcpcapi "github.com/openshift/installer/pkg/infrastructure/gcp/clusterapi"
	ibmcloudcapi "github.com/openshift/installer/pkg/infrastructure/ibmcloud/clusterapi"
	nutanixcapi "github.com/openshift/installer/pkg/infrastructure/nutanix/clusterapi"
	openstackcapi "github.com/openshift/installer/pkg/infrastructure/openstack/clusterapi"
	powervscapi "github.com/openshift/installer/pkg/infrastructure/powervs/clusterapi"
	vspherecapi "github.com/openshift/installer/pkg/infrastructure/vsphere/clusterapi"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/featuregates"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// ProviderForPlatform returns the stages to run to provision the infrastructure for the specified platform.
func ProviderForPlatform(platform string, fg featuregates.FeatureGate) (infrastructure.Provider, error) {
	switch platform {
	case awstypes.Name:
		return clusterapi.InitializeProvider(&awscapi.Provider{}), nil
	case azuretypes.Name:
		return clusterapi.InitializeProvider(&azurecapi.Provider{}), nil
	case gcptypes.Name:
		return clusterapi.InitializeProvider(gcpcapi.Provider{}), nil
	case ibmcloudtypes.Name:
		return clusterapi.InitializeProvider(ibmcloudcapi.Provider{}), nil
	case vspheretypes.Name:
		return clusterapi.InitializeProvider(vspherecapi.Provider{}), nil
	case powervstypes.Name:
		return clusterapi.InitializeProvider(powervscapi.Provider{}), nil
	case openstacktypes.Name:
		return clusterapi.InitializeProvider(openstackcapi.Provider{}), nil
	case nutanixtypes.Name:
		return clusterapi.InitializeProvider(nutanixcapi.Provider{}), nil
	}
	return nil, fmt.Errorf("platform %q is not supported in the altinfra Installer build", platform)
}
