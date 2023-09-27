//go:build !altinfra
// +build !altinfra

package platform

import (
	"fmt"

	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform/stages/alibabacloud"
	"github.com/openshift/installer/pkg/terraform/stages/aws"
	"github.com/openshift/installer/pkg/terraform/stages/azure"
	"github.com/openshift/installer/pkg/terraform/stages/baremetal"
	"github.com/openshift/installer/pkg/terraform/stages/gcp"
	"github.com/openshift/installer/pkg/terraform/stages/ibmcloud"
	"github.com/openshift/installer/pkg/terraform/stages/libvirt"
	"github.com/openshift/installer/pkg/terraform/stages/nutanix"
	"github.com/openshift/installer/pkg/terraform/stages/openstack"
	"github.com/openshift/installer/pkg/terraform/stages/ovirt"
	"github.com/openshift/installer/pkg/terraform/stages/powervs"
	"github.com/openshift/installer/pkg/terraform/stages/vsphere"
	alibabacloudtypes "github.com/openshift/installer/pkg/types/alibabacloud"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	externaltypes "github.com/openshift/installer/pkg/types/external"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	libvirttypes "github.com/openshift/installer/pkg/types/libvirt"
	nonetypes "github.com/openshift/installer/pkg/types/none"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// ProviderForPlatform returns the stages to run to provision the infrastructure for the specified platform.
func ProviderForPlatform(platform, installDir string) ([]infrastructure.Stage, func() error, error) {
	var stages []infrastructure.Stage
	var cleanup func() error
	var err error
	switch platform {
	case alibabacloudtypes.Name:
		if stages, cleanup, err = alibabacloud.InitializeProvider(installDir); err != nil {
			return nil, nil, err
		}
	case awstypes.Name:
		if stages, cleanup, err = aws.InitializeProvider(installDir); err != nil {
			return nil, nil, err
		}
	case azuretypes.Name:
		if stages, cleanup, err = azure.InitializeProvider(installDir); err != nil {
			return nil, nil, err
		}
	case azuretypes.StackTerraformName:
		if stages, cleanup, err = azure.InitializeStackProvider(installDir); err != nil {
			return nil, nil, err
		}
	case baremetaltypes.Name:
		if stages, cleanup, err = baremetal.InitializeProvider(installDir); err != nil {
			return nil, nil, err
		}
	case gcptypes.Name:
		if stages, cleanup, err = gcp.InitializeProvider(installDir); err != nil {
			return nil, nil, err
		}
	case ibmcloudtypes.Name:
		if stages, cleanup, err = ibmcloud.InitializeProvider(installDir); err != nil {
			return nil, nil, err
		}
	case libvirttypes.Name:
		if stages, cleanup, err = libvirt.InitializeProvider(installDir); err != nil {
			return nil, nil, err
		}
	case nutanixtypes.Name:
		if stages, cleanup, err = nutanix.InitializeProvider(installDir); err != nil {
			return nil, nil, err
		}
	case powervstypes.Name:
		if stages, cleanup, err = powervs.InitializeProvider(installDir); err != nil {
			return nil, nil, err
		}
	case openstacktypes.Name:
		if stages, cleanup, err = openstack.InitializeProvider(installDir); err != nil {
			return nil, nil, err
		}
	case ovirttypes.Name:
		if stages, cleanup, err = ovirt.InitializeProvider(installDir); err != nil {
			return nil, nil, err
		}
	case vspheretypes.Name:
		if stages, cleanup, err = vsphere.InitializeProvider(installDir); err != nil {
			return nil, nil, err
		}
	case nonetypes.Name:
		// terraform is not used when the platform is "none"
		return []infrastructure.Stage{}, nil, nil
	case externaltypes.Name:
		// terraform is not used when the platform is "external"
		return []infrastructure.Stage{}, nil, nil
	default:
		panic(fmt.Sprintf("unsupported platform %q", platform))
	}
	return stages, cleanup, nil
}
