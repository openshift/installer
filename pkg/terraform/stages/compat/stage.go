package compat

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/terraform"
	gatherazure "github.com/openshift/installer/pkg/terraform/gather/azure"
	gatherbaremetal "github.com/openshift/installer/pkg/terraform/gather/baremetal"
	gatherkubevirt "github.com/openshift/installer/pkg/terraform/gather/kubevirt"
	gatherlibvirt "github.com/openshift/installer/pkg/terraform/gather/libvirt"
	gatheropenstack "github.com/openshift/installer/pkg/terraform/gather/openstack"
	gatherovirt "github.com/openshift/installer/pkg/terraform/gather/ovirt"
	gathervsphere "github.com/openshift/installer/pkg/terraform/gather/vsphere"
	"github.com/openshift/installer/pkg/types"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	kubevirttypes "github.com/openshift/installer/pkg/types/kubevirt"
	libvirttypes "github.com/openshift/installer/pkg/types/libvirt"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// PlatformStages are the stages to run to provision the infrastructure used the legacy compat procedures.
func PlatformStages(platform string) []terraform.Stage {
	return []terraform.Stage{stage{platform: platform}}
}

type stage struct {
	platform string
}

func (s stage) Name() string {
	return ""
}

func (s stage) StateFilename() string {
	return "terraform.tfstate"
}

func (s stage) OutputsFilename() string {
	return "outputs.tfvars.json"
}

func (s stage) DestroyWithBootstrap() bool {
	return true
}

func (s stage) Destroy(directory string, extraArgs []string) error {
	switch s.platform {
	case libvirttypes.Name:
		// First remove the bootstrap node from DNS
		if _, err := terraform.Apply(directory, s.platform, s, append(extraArgs, "-var=bootstrap_dns=false")...); err != nil {
			return errors.Wrap(err, "Terraform apply")
		}
	case ovirttypes.Name:
		extraArgs = append(extraArgs, "-target=module.template.ovirt_vm.tmp_import_vm")
		extraArgs = append(extraArgs, "-target=module.template.ovirt_image_transfer.releaseimage")
	}

	extraArgs = append(extraArgs, "-target=module.bootstrap")

	return errors.Wrap(terraform.Destroy(directory, s.platform, s, extraArgs...), "terraform destroy")
}

func (s stage) ExtractHostAddresses(directory string, config *types.InstallConfig) (string, int, []string, error) {
	tfStateFilePath := filepath.Join(directory, s.StateFilename())
	_, err := os.Stat(tfStateFilePath)
	if os.IsNotExist(err) {
		return "", 0, nil, nil
	}
	if err != nil {
		return "", 0, nil, err
	}

	tfstate, err := terraform.ReadState(tfStateFilePath)
	if err != nil {
		return "", 0, nil, errors.Wrapf(err, "failed to read state from %q", tfStateFilePath)
	}
	bootstrap, port, masters, err := extractHostAddresses(config, tfstate)
	return bootstrap, port, masters, errors.Wrapf(err, "failed to get bootstrap and control plane host addresses from %q", tfStateFilePath)
}

func extractHostAddresses(config *types.InstallConfig, tfstate *terraform.State) (bootstrap string, port int, masters []string, err error) {
	port = 22
	switch config.Platform.Name() {
	case azuretypes.Name:
		bootstrap, err = gatherazure.BootstrapIP(tfstate)
		if err != nil {
			return
		}
		masters, err = gatherazure.ControlPlaneIPs(tfstate)
		if err != nil {
			logrus.Error(err)
		}
	case baremetaltypes.Name:
		bootstrap = config.Platform.BareMetal.BootstrapProvisioningIP
		masters, err = gatherbaremetal.ControlPlaneIPs(config, tfstate)
		if err != nil {
			return
		}
	case libvirttypes.Name:
		bootstrap, err = gatherlibvirt.BootstrapIP(tfstate)
		if err != nil {
			return
		}
		masters, err = gatherlibvirt.ControlPlaneIPs(tfstate)
		if err != nil {
			logrus.Error(err)
		}
	case openstacktypes.Name:
		bootstrap, err = gatheropenstack.BootstrapIP(tfstate)
		if err != nil {
			return
		}
		masters, err = gatheropenstack.ControlPlaneIPs(tfstate)
		if err != nil {
			logrus.Error(err)
		}
	case ovirttypes.Name:
		bootstrap, err = gatherovirt.BootstrapIP(tfstate)
		if err != nil {
			return
		}
		masters, err = gatherovirt.ControlPlaneIPs(tfstate)
		if err != nil {
			logrus.Error(err)
		}

	case vspheretypes.Name:
		bootstrap, err = gathervsphere.BootstrapIP(config, tfstate)
		if err != nil {
			return
		}
		masters, err = gathervsphere.ControlPlaneIPs(config, tfstate)
		if err != nil {
			logrus.Error(err)
		}
	case kubevirttypes.Name:
		bootstrap, err = gatherkubevirt.BootstrapIP(tfstate)
		if err != nil {
			return
		}
		masters, err = gatherkubevirt.ControlPlaneIPs(tfstate)
		if err != nil {
			logrus.Error(err)
		}
	}
	return bootstrap, port, masters, nil
}
