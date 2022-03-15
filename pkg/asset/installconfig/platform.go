package installconfig

import (
	"fmt"
	"sort"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	alibabacloudconfig "github.com/openshift/installer/pkg/asset/installconfig/alibabacloud"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	azureconfig "github.com/openshift/installer/pkg/asset/installconfig/azure"
	baremetalconfig "github.com/openshift/installer/pkg/asset/installconfig/baremetal"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	ibmcloudconfig "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	libvirtconfig "github.com/openshift/installer/pkg/asset/installconfig/libvirt"
	openstackconfig "github.com/openshift/installer/pkg/asset/installconfig/openstack"
	ovirtconfig "github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	powervsconfig "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	vsphereconfig "github.com/openshift/installer/pkg/asset/installconfig/vsphere"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/alibabacloud"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// Platform is an asset that queries the user for the platform on which to install
// the cluster.
type platform struct {
	types.Platform
}

var _ asset.Asset = (*platform)(nil)

// Dependencies returns no dependencies.
func (a *platform) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate queries for input from the user.
func (a *platform) Generate(asset.Parents) error {
	platform, err := a.queryUserForPlatform()
	if err != nil {
		return err
	}

	switch platform {
	case alibabacloud.Name:
		a.AlibabaCloud, err = alibabacloudconfig.Platform()
	case aws.Name:
		a.AWS, err = awsconfig.Platform()
	case azure.Name:
		a.Azure, err = azureconfig.Platform()
	case baremetal.Name:
		a.BareMetal, err = baremetalconfig.Platform()
	case gcp.Name:
		a.GCP, err = gcpconfig.Platform()
	case ibmcloud.Name:
		a.IBMCloud, err = ibmcloudconfig.Platform()
	case libvirt.Name:
		a.Libvirt, err = libvirtconfig.Platform()
	case none.Name:
		a.None = &none.Platform{}
	case openstack.Name:
		a.OpenStack, err = openstackconfig.Platform()
	case ovirt.Name:
		a.Ovirt, err = ovirtconfig.Platform()
	case vsphere.Name:
		a.VSphere, err = vsphereconfig.Platform()
	case powervs.Name:
		a.PowerVS, err = powervsconfig.Platform()
	default:
		return fmt.Errorf("unknown platform type %q", platform)
	}

	return err
}

// Name returns the human-friendly name of the asset.
func (a *platform) Name() string {
	return "Platform"
}

func (a *platform) queryUserForPlatform() (platform string, err error) {
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Platform",
				Options: types.PlatformNames,
				Help:    "The platform on which the cluster will run.  For a full list of platforms, including those not supported by this wizard, see https://github.com/openshift/installer",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := ans.(core.OptionAnswer).Value
				i := sort.SearchStrings(types.PlatformNames, choice)
				if i == len(types.PlatformNames) || types.PlatformNames[i] != choice {
					return errors.Errorf("invalid platform %q", choice)
				}
				return nil
			}),
		},
	}, &platform)
	return
}

func (a *platform) CurrentName() string {
	return a.Platform.Name()
}
