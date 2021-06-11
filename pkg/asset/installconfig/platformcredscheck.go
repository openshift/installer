package installconfig

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	kubevirtconfig "github.com/openshift/installer/pkg/asset/installconfig/kubevirt"
	openstackconfig "github.com/openshift/installer/pkg/asset/installconfig/openstack"
	ovirtconfig "github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	powervsconfig "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/kubevirt"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/types/vsphere"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PlatformCredsCheck is an asset that checks the platform credentials, asks for them or errors out if invalid
// the cluster.
type PlatformCredsCheck struct {
}

var _ asset.Asset = (*PlatformCredsCheck)(nil)

// Dependencies returns the dependencies for PlatformCredsCheck
func (a *PlatformCredsCheck) Dependencies() []asset.Asset {
	return []asset.Asset{
		&InstallConfig{},
	}
}

// Generate queries for input from the user.
func (a *PlatformCredsCheck) Generate(dependencies asset.Parents) error {
	ctx := context.TODO()
	ic := &InstallConfig{}
	dependencies.Get(ic)

	var err error
	platform := ic.Config.Platform.Name()
	switch platform {
	case aws.Name:
		_, err := ic.AWS.Session(ctx)
		if err != nil {
			return err
		}
	case gcp.Name:
		_, err = gcpconfig.GetSession(ctx)
		if err != nil {
			return errors.Wrap(err, "creating GCP session")
		}
        // meh ... is there a way to name this session ibmcloud? but have the installconfig Name as powervs?
	case powervs.Name:
		_, err = powervsconfig.GetSession()
		if err != nil {
			return errors.Wrap(err, "creating IBM Cloud session")
		}
	case openstack.Name:
		_, err = openstackconfig.GetSession(ic.Config.Platform.OpenStack.Cloud)
		if err != nil {
			return errors.Wrap(err, "creating OpenStack session")
		}
	case baremetal.Name, libvirt.Name, none.Name, vsphere.Name:
		// no creds to check
	case azure.Name:
		_, err = ic.Azure.Session()
		if err != nil {
			return errors.Wrap(err, "creating Azure session")
		}
	case ovirt.Name:
		con, err := ovirtconfig.NewConnection()
		if err != nil {
			return errors.Wrap(err, "creating Engine connection")
		}
		err = con.Test()
		if err != nil {
			return errors.Wrap(err, "testing Engine connection")
		}
	case kubevirt.Name:
		client, err := kubevirtconfig.NewClient()
		if err != nil {
			return errors.Wrap(err, "creating KubeVirt client")
		}
		// Test the connection to InfraCluster by calling ListVM API
		if _, err = client.ListVirtualMachine(context.Background(), ic.Config.Platform.Kubevirt.Namespace, metav1.ListOptions{}); err != nil {
			return errors.Wrap(err, "testing KubeVirt connection")
		}
	default:
		err = fmt.Errorf("unknown platform type %q", platform)
	}

	return err
}

// Name returns the human-friendly name of the asset.
func (a *PlatformCredsCheck) Name() string {
	return "Platform Credentials Check"
}
