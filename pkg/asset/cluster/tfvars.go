package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/user"

	igntypes "github.com/coreos/ignition/config/v2_2/types"
	gcpprovider "github.com/openshift/cluster-api-provider-gcp/pkg/apis/gcpprovider/v1beta1"
	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	awsprovider "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
	azureprovider "sigs.k8s.io/cluster-api-provider-azure/pkg/apis/azureprovider/v1alpha1"
	openstackprovider "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	azureconfig "github.com/openshift/installer/pkg/asset/installconfig/azure"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/tfvars"
	awstfvars "github.com/openshift/installer/pkg/tfvars/aws"
	azuretfvars "github.com/openshift/installer/pkg/tfvars/azure"
	baremetaltfvars "github.com/openshift/installer/pkg/tfvars/baremetal"
	gcptfvars "github.com/openshift/installer/pkg/tfvars/gcp"
	libvirttfvars "github.com/openshift/installer/pkg/tfvars/libvirt"
	openstacktfvars "github.com/openshift/installer/pkg/tfvars/openstack"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/version"
)

const (
	// TfVarsFileName is the filename for Terraform variables.
	TfVarsFileName = "terraform.tfvars.json"

	// TfPlatformVarsFileName is a template for platform-specific
	// Terraform variable files.
	//
	// https://www.terraform.io/docs/configuration/variables.html#variable-files
	TfPlatformVarsFileName = "terraform.%s.auto.tfvars.json"

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
		&machines.Worker{},
	}
}

// Generate generates the terraform.tfvars file.
func (t *TerraformVariables) Generate(parents asset.Parents) error {
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	bootstrapIgnAsset := &bootstrap.Bootstrap{}
	masterIgnAsset := &machine.Master{}
	mastersAsset := &machines.Master{}
	workersAsset := &machines.Worker{}
	rhcosImage := new(rhcos.Image)
	parents.Get(clusterID, installConfig, bootstrapIgnAsset, masterIgnAsset, mastersAsset, workersAsset, rhcosImage)

	platform := installConfig.Config.Platform.Name()
	switch platform {
	case none.Name, vsphere.Name:
		return errors.Errorf("cannot create the cluster because %q is a UPI platform", platform)
	}

	masterIgn := string(masterIgnAsset.Files()[0].Data)
	bootstrapIgn, err := injectInstallInfo(bootstrapIgnAsset.Files()[0].Data)
	if err != nil {
		return errors.Wrap(err, "unable to inject installation info")
	}

	masterCount := len(mastersAsset.MachineFiles)
	data, err := tfvars.TFVars(
		clusterID.InfraID,
		installConfig.Config.ClusterDomain(),
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

	switch platform {
	case aws.Name:
		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}
		masterConfigs := make([]*awsprovider.AWSMachineProviderConfig, len(masters))
		for i, m := range masters {
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*awsprovider.AWSMachineProviderConfig)
		}
		workers, err := workersAsset.MachineSets()
		if err != nil {
			return err
		}
		workerConfigs := make([]*awsprovider.AWSMachineProviderConfig, len(workers))
		for i, m := range workers {
			workerConfigs[i] = m.Spec.Template.Spec.ProviderSpec.Value.Object.(*awsprovider.AWSMachineProviderConfig)
		}
		data, err := awstfvars.TFVars(masterConfigs, workerConfigs)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	case azure.Name:
		sess, err := azureconfig.GetSession()
		if err != nil {
			return err
		}
		auth := azuretfvars.Auth{
			SubscriptionID: sess.Credentials.SubscriptionID,
			ClientID:       sess.Credentials.ClientID,
			ClientSecret:   sess.Credentials.ClientSecret,
			TenantID:       sess.Credentials.TenantID,
		}
		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}
		masterConfigs := make([]*azureprovider.AzureMachineProviderSpec, len(masters))
		for i, m := range masters {
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*azureprovider.AzureMachineProviderSpec)
		}
		data, err := azuretfvars.TFVars(
			auth,
			installConfig.Config.Azure.BaseDomainResourceGroupName,
			masterConfigs,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	case gcp.Name:
		sess, err := gcpconfig.GetSession(context.TODO())
		if err != nil {
			return err
		}
		auth := gcptfvars.Auth{
			ProjectID:      installConfig.Config.GCP.ProjectID,
			ServiceAccount: string(sess.Credentials.JSON),
		}
		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}
		masterConfigs := make([]*gcpprovider.GCPMachineProviderSpec, len(masters))
		for i, m := range masters {
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*gcpprovider.GCPMachineProviderSpec)
		}
		publicZone, err := gcpconfig.GetPublicZone(context.TODO(), installConfig.Config.GCP.ProjectID, installConfig.Config.BaseDomain)
		if err != nil {
			return errors.Wrapf(err, "failed to get GCP public zone")
		}
		data, err := gcptfvars.TFVars(
			auth,
			masterConfigs,
			publicZone.Name,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	case libvirt.Name:
		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}
		data, err = libvirttfvars.TFVars(
			masters[0].Spec.ProviderSpec.Value.Object.(*libvirtprovider.LibvirtMachineProviderConfig),
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
	case openstack.Name:
		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}
		data, err = openstacktfvars.TFVars(
			masters[0].Spec.ProviderSpec.Value.Object.(*openstackprovider.OpenstackProviderSpec),
			installConfig.Config.Platform.OpenStack.Region,
			installConfig.Config.Platform.OpenStack.ExternalNetwork,
			installConfig.Config.Platform.OpenStack.LbFloatingIP,
			installConfig.Config.Platform.OpenStack.APIVIP,
			installConfig.Config.Platform.OpenStack.DNSVIP,
			installConfig.Config.Platform.OpenStack.TrunkSupport,
			installConfig.Config.Platform.OpenStack.OctaviaSupport,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	case baremetal.Name:
		data, err = baremetaltfvars.TFVars(
			installConfig.Config.Platform.BareMetal.LibvirtURI,
			installConfig.Config.Platform.BareMetal.IronicURI,
			string(*rhcosImage),
			"baremetal",
			"provisioning",
			installConfig.Config.Platform.BareMetal.Hosts,
			installConfig.Config.Platform.BareMetal.Image,
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

// injectInstallInfo adds information about the installer and its invoker as a
// ConfigMap to the provided bootstrap Ignition config.
func injectInstallInfo(bootstrap []byte) (string, error) {
	config := &igntypes.Config{}
	if err := json.Unmarshal(bootstrap, &config); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal bootstrap Ignition config")
	}

	var invoker string
	if env, ok := os.LookupEnv("OPENSHIFT_INSTALL_INVOKER"); ok {
		invoker = env
	} else if user, err := user.Current(); err == nil {
		invoker = user.Username
	} else {
		logrus.Warnf("Unable to determine username: %v", err)
		invoker = "<unknown>"
	}

	config.Storage.Files = append(config.Storage.Files, ignition.FileFromString("/opt/openshift/manifests/openshift-install.yml", "root", 0644, fmt.Sprintf(`---
apiVersion: v1
kind: ConfigMap
metadata:
  name: openshift-install
  namespace: openshift-config
data:
  version: "%s"
  invoker: "%s"
`, version.Raw, invoker)))

	ign, err := json.Marshal(config)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal bootstrap Ignition config")
	}

	return string(ign), nil
}
