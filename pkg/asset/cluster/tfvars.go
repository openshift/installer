package cluster

import (
	"fmt"
	"os"

	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1alpha1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/tfvars"
	awstfvars "github.com/openshift/installer/pkg/tfvars/aws"
	libvirttfvars "github.com/openshift/installer/pkg/tfvars/libvirt"
	openstacktfvars "github.com/openshift/installer/pkg/tfvars/openstack"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	awsprovider "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
	openstackprovider "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"
)

const (
	// TfVarsFileName is the filename for Terraform variables.
	TfVarsFileName = "terraform.tfvars"

	// TfPlatformVarsFileName is a template for platform-specific
	// Terraform variable files.
	//
	// https://www.terraform.io/docs/configuration/variables.html#variable-files
	TfPlatformVarsFileName = "terraform.%s.auto.tfvars"

	tfvarsAssetName = "Terraform Variables"
)

// TerraformVariables depends on InstallConfig and
// Ignition to generate the terrafor.tfvars.
type TerraformVariables struct {
	FileList []*asset.File
}

var _ asset.WritableAsset = (*TerraformVariables)(nil)

// Name returns the human-friendly name of the asset.
func (t *TerraformVariables) Name() string {
	return tfvarsAssetName
}

// Dependencies returns the dependency of the TerraformVariable
func (t *TerraformVariables) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		&installconfig.InstallConfig{},
		new(rhcos.Image),
		&bootstrap.Bootstrap{},
		&machine.Master{},
		&machines.Master{},
	}
}

// Generate generates the terraform.tfvars file.
func (t *TerraformVariables) Generate(parents asset.Parents) error {
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	bootstrapIgnAsset := &bootstrap.Bootstrap{}
	masterIgnAsset := &machine.Master{}
	mastersAsset := &machines.Master{}
	rhcosImage := new(rhcos.Image)
	parents.Get(clusterID, installConfig, bootstrapIgnAsset, masterIgnAsset, mastersAsset, rhcosImage)

	bootstrapIgn := string(bootstrapIgnAsset.Files()[0].Data)
	masterIgn := string(masterIgnAsset.Files()[0].Data)

	masterCount := len(mastersAsset.Machines)
	if mastersAsset.Machines == nil {
		masterCount = len(mastersAsset.MachinesDeprecated)
	}
	data, err := tfvars.TFVars(
		clusterID.ClusterID,
		installConfig.Config.ObjectMeta.Name,
		installConfig.Config.BaseDomain,
		&installConfig.Config.Networking.MachineCIDR.IPNet,
		bootstrapIgn,
		masterIgn,
		masterCount,
	)
	if err != nil {
		return errors.Wrap(err, "failed to get Terraform variables")
	}
	t.FileList = []*asset.File{
		{
			Filename: TfVarsFileName,
			Data:     data,
		},
	}

	if masterCount == 0 {
		return errors.Errorf("master slice cannot be empty")
	}

	switch platform := installConfig.Config.Platform.Name(); platform {
	case aws.Name:
		data, err = awstfvars.TFVars(
			mastersAsset.Machines[0].Spec.ProviderSpec.Value.Object.(*awsprovider.AWSMachineProviderConfig),
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	case libvirt.Name:
		data, err = libvirttfvars.TFVars(
			mastersAsset.MachinesDeprecated[0].Spec.ProviderSpec.Value.Object.(*libvirtprovider.LibvirtMachineProviderConfig),
			string(*rhcosImage),
			&installConfig.Config.Networking.MachineCIDR.IPNet,
			installConfig.Config.Platform.Libvirt.Network.IfName,
			masterCount,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	case none.Name:
	case openstack.Name:
		data, err = openstacktfvars.TFVars(
			mastersAsset.MachinesDeprecated[0].Spec.ProviderSpec.Value.Object.(*openstackprovider.OpenstackProviderSpec),
			installConfig.Config.Platform.OpenStack.Region,
			installConfig.Config.Platform.OpenStack.ExternalNetwork,
			installConfig.Config.Platform.OpenStack.TrunkSupport,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	default:
		logrus.Warnf("unrecognized platform %s", platform)
	}

	return nil
}

// Files returns the files generated by the asset.
func (t *TerraformVariables) Files() []*asset.File {
	return t.FileList
}

// Load reads the terraform.tfvars from disk.
func (t *TerraformVariables) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(TfVarsFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	t.FileList = []*asset.File{file}

	fileList, err := f.FetchByPattern(fmt.Sprintf(TfPlatformVarsFileName, "*"))
	if err != nil {
		return false, err
	}
	t.FileList = append(t.FileList, fileList...)

	return true, nil
}
