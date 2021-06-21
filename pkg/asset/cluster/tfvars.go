package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	coreosarch "github.com/coreos/stream-metadata-go/arch"
	gcpprovider "github.com/openshift/cluster-api-provider-gcp/pkg/apis/gcpprovider/v1beta1"
	kubevirtprovider "github.com/openshift/cluster-api-provider-kubevirt/pkg/apis/kubevirtprovider/v1alpha1"
	kubevirtutils "github.com/openshift/cluster-api-provider-kubevirt/pkg/utils"
	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	ovirtprovider "github.com/openshift/cluster-api-provider-ovirt/pkg/apis/ovirtprovider/v1beta1"
	vsphereprovider "github.com/openshift/machine-api-operator/pkg/apis/vsphereprovider/v1beta1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	awsprovider "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1beta1"
	azureprovider "sigs.k8s.io/cluster-api-provider-azure/pkg/apis/azureprovider/v1beta1"
	openstackprovider "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	baremetalbootstrap "github.com/openshift/installer/pkg/asset/ignition/bootstrap/baremetal"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	openstackconfig "github.com/openshift/installer/pkg/asset/installconfig/openstack"
	ovirtconfig "github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/openshiftinstall"
	"github.com/openshift/installer/pkg/asset/rhcos"
	rhcospkg "github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/tfvars"
	awstfvars "github.com/openshift/installer/pkg/tfvars/aws"
	azuretfvars "github.com/openshift/installer/pkg/tfvars/azure"
	baremetaltfvars "github.com/openshift/installer/pkg/tfvars/baremetal"
	gcptfvars "github.com/openshift/installer/pkg/tfvars/gcp"
	ibmcloudtfvars "github.com/openshift/installer/pkg/tfvars/ibmcloud"
	kubevirttfvars "github.com/openshift/installer/pkg/tfvars/kubevirt"
	libvirttfvars "github.com/openshift/installer/pkg/tfvars/libvirt"
	openstacktfvars "github.com/openshift/installer/pkg/tfvars/openstack"
	ovirttfvars "github.com/openshift/installer/pkg/tfvars/ovirt"
	vspheretfvars "github.com/openshift/installer/pkg/tfvars/vsphere"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/kubevirt"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/vsphere"
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
		new(rhcos.BootstrapImage),
		&bootstrap.Bootstrap{},
		&machine.Master{},
		&machines.Master{},
		&machines.Worker{},
		&baremetalbootstrap.IronicCreds{},
	}
}

// Generate generates the terraform.tfvars file.
func (t *TerraformVariables) Generate(parents asset.Parents) error {
	ctx := context.TODO()
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	bootstrapIgnAsset := &bootstrap.Bootstrap{}
	masterIgnAsset := &machine.Master{}
	mastersAsset := &machines.Master{}
	workersAsset := &machines.Worker{}
	rhcosImage := new(rhcos.Image)
	rhcosBootstrapImage := new(rhcos.BootstrapImage)
	ironicCreds := &baremetalbootstrap.IronicCreds{}
	parents.Get(clusterID, installConfig, bootstrapIgnAsset, masterIgnAsset, mastersAsset, workersAsset, rhcosImage, rhcosBootstrapImage, ironicCreds)

	platform := installConfig.Config.Platform.Name()
	switch platform {
	case none.Name:
		return errors.Errorf("cannot create the cluster because %q is a UPI platform", platform)
	}

	masterIgn := string(masterIgnAsset.Files()[0].Data)
	bootstrapIgn, err := injectInstallInfo(bootstrapIgnAsset.Files()[0].Data)
	if err != nil {
		return errors.Wrap(err, "unable to inject installation info")
	}

	var useIPv4, useIPv6 bool
	for _, network := range installConfig.Config.Networking.ServiceNetwork {
		if network.IP.To4() != nil {
			useIPv4 = true
		} else {
			useIPv6 = true
		}
	}

	machineV4CIDRs, machineV6CIDRs := []string{}, []string{}
	for _, network := range installConfig.Config.Networking.MachineNetwork {
		if network.CIDR.IPNet.IP.To4() != nil {
			machineV4CIDRs = append(machineV4CIDRs, network.CIDR.IPNet.String())
		} else {
			machineV6CIDRs = append(machineV6CIDRs, network.CIDR.IPNet.String())
		}
	}

	masterCount := len(mastersAsset.MachineFiles)
	data, err := tfvars.TFVars(
		clusterID.InfraID,
		installConfig.Config.ClusterDomain(),
		installConfig.Config.BaseDomain,
		machineV4CIDRs,
		machineV6CIDRs,
		useIPv4,
		useIPv6,
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
		var vpc string
		var privateSubnets []string
		var publicSubnets []string

		if len(installConfig.Config.Platform.AWS.Subnets) > 0 {
			subnets, err := installConfig.AWS.PrivateSubnets(ctx)
			if err != nil {
				return err
			}

			for id := range subnets {
				privateSubnets = append(privateSubnets, id)
			}

			subnets, err = installConfig.AWS.PublicSubnets(ctx)
			if err != nil {
				return err
			}

			for id := range subnets {
				publicSubnets = append(publicSubnets, id)
			}

			vpc, err = installConfig.AWS.VPC(ctx)
			if err != nil {
				return err
			}
		}

		sess, err := installConfig.AWS.Session(ctx)
		if err != nil {
			return err
		}
		object := "bootstrap.ign"
		bucket := fmt.Sprintf("%s-bootstrap", clusterID.InfraID)
		url, err := awsconfig.PresignedS3URL(sess, installConfig.Config.Platform.AWS.Region, bucket, object)
		if err != nil {
			return err
		}
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
		osImage := strings.SplitN(string(*rhcosImage), ",", 2)
		osImageID := osImage[0]
		osImageRegion := installConfig.Config.AWS.Region
		if len(osImage) == 2 {
			osImageRegion = osImage[1]
		}

		workerIAMRoleName := ""
		if mp := installConfig.Config.WorkerMachinePool(); mp != nil {
			awsMP := &aws.MachinePool{}
			awsMP.Set(installConfig.Config.AWS.DefaultMachinePlatform)
			awsMP.Set(mp.Platform.AWS)
			workerIAMRoleName = awsMP.IAMRole
		}

		masterIAMRoleName := ""
		if mp := installConfig.Config.ControlPlane; mp != nil {
			awsMP := &aws.MachinePool{}
			awsMP.Set(installConfig.Config.AWS.DefaultMachinePlatform)
			awsMP.Set(mp.Platform.AWS)
			masterIAMRoleName = awsMP.IAMRole
		}

		data, err := awstfvars.TFVars(awstfvars.TFVarsSources{
			VPC:                   vpc,
			PrivateSubnets:        privateSubnets,
			PublicSubnets:         publicSubnets,
			InternalZone:          installConfig.Config.AWS.HostedZone,
			Services:              installConfig.Config.AWS.ServiceEndpoints,
			Publish:               installConfig.Config.Publish,
			MasterConfigs:         masterConfigs,
			WorkerConfigs:         workerConfigs,
			AMIID:                 osImageID,
			AMIRegion:             osImageRegion,
			IgnitionBucket:        bucket,
			IgnitionPresignedURL:  url,
			AdditionalTrustBundle: installConfig.Config.AdditionalTrustBundle,
			MasterIAMRoleName:     masterIAMRoleName,
			WorkerIAMRoleName:     workerIAMRoleName,
			Architecture:          installConfig.Config.ControlPlane.Architecture,
		})
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	case azure.Name:
		session, err := installConfig.Azure.Session()
		if err != nil {
			return err
		}

		auth := azuretfvars.Auth{
			SubscriptionID: session.Credentials.SubscriptionID,
			ClientID:       session.Credentials.ClientID,
			ClientSecret:   session.Credentials.ClientSecret,
			TenantID:       session.Credentials.TenantID,
		}
		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}
		masterConfigs := make([]*azureprovider.AzureMachineProviderSpec, len(masters))
		for i, m := range masters {
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*azureprovider.AzureMachineProviderSpec)
		}
		workers, err := workersAsset.MachineSets()
		if err != nil {
			return err
		}
		workerConfigs := make([]*azureprovider.AzureMachineProviderSpec, len(workers))
		for i, w := range workers {
			workerConfigs[i] = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*azureprovider.AzureMachineProviderSpec)
		}

		preexistingnetwork := installConfig.Config.Azure.VirtualNetwork != ""
		data, err := azuretfvars.TFVars(
			azuretfvars.TFVarsSources{
				Auth:                        auth,
				CloudName:                   installConfig.Config.Azure.CloudName,
				ResourceGroupName:           installConfig.Config.Azure.ResourceGroupName,
				BaseDomainResourceGroupName: installConfig.Config.Azure.BaseDomainResourceGroupName,
				MasterConfigs:               masterConfigs,
				WorkerConfigs:               workerConfigs,
				ImageURL:                    string(*rhcosImage),
				PreexistingNetwork:          preexistingnetwork,
				Publish:                     installConfig.Config.Publish,
				OutboundType:                installConfig.Config.Azure.OutboundType,
			},
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	case gcp.Name:
		var publicZoneName string
		sess, err := gcpconfig.GetSession(ctx)
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
		workers, err := workersAsset.MachineSets()
		if err != nil {
			return err
		}
		workerConfigs := make([]*gcpprovider.GCPMachineProviderSpec, len(workers))
		for i, w := range workers {
			workerConfigs[i] = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*gcpprovider.GCPMachineProviderSpec)
		}
		if installConfig.Config.Publish == types.ExternalPublishingStrategy {
			publicZone, err := gcpconfig.GetPublicZone(ctx, installConfig.Config.GCP.ProjectID, installConfig.Config.BaseDomain)
			if err != nil {
				return errors.Wrapf(err, "failed to get GCP public zone")
			}
			publicZoneName = publicZone.Name
		}
		preexistingnetwork := installConfig.Config.GCP.Network != ""

		archName := coreosarch.RpmArch(string(installConfig.Config.ControlPlane.Architecture))
		st, err := rhcospkg.FetchCoreOSBuild(ctx)
		if err != nil {
			return err
		}
		streamArch, err := st.GetArchitecture(archName)
		if err != nil {
			return err
		}

		img := streamArch.Images.Gcp
		if img == nil {
			return fmt.Errorf("%s: No GCP build found", st.FormatPrefix(archName))
		}
		// For backwards compatibility, we generate this URL to the image (only applies to RHCOS, not FCOS/OKD)
		// right now.  It will only be used if nested virt or other licenses are enabled, which we
		// really should deprecate and remove - xref https://github.com/openshift/installer/pull/4696
		imageURL := fmt.Sprintf("https://storage.googleapis.com/rhcos/rhcos/%s.tar.gz", img.Name)
		data, err := gcptfvars.TFVars(
			gcptfvars.TFVarsSources{
				Auth:               auth,
				MasterConfigs:      masterConfigs,
				WorkerConfigs:      workerConfigs,
				ImageURI:           imageURL,
				ImageLicenses:      installConfig.Config.GCP.Licenses,
				PublicZoneName:     publicZoneName,
				PublishStrategy:    installConfig.Config.Publish,
				PreexistingNetwork: preexistingnetwork,
			},
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	case ibmcloud.Name:
		client, err := installConfig.IBMCloud.Client()
		if err != nil {
			return err
		}
		auth := ibmcloudtfvars.Auth{
			APIKey: client.Authenticator.ApiKey,
		}

		// TODO: IBM: Get master and worker machine info
		// masters, err := mastersAsset.Machines()
		// if err != nil {
		// 	return err
		// }
		// masterConfigs := make([]*ibmcloudprovider.IBMCloudMachineProviderSpec, len(masters))
		// for i, m := range masters {
		// 	masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*ibmcloudprovider.IBMCloudMachineProviderSpec)
		// }
		// workers, err := workersAsset.MachineSets()
		// if err != nil {
		// 	return err
		// }
		// workerConfigs := make([]*ibmcloudprovider.IBMCloudMachineProviderSpec, len(workers))
		// for i, w := range workers {
		// 	workerConfigs[i] = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*ibmcloudprovider.IBMCloudMachineProviderSpec)
		// }

		// TODO: IBM: Fetch config from masterConfig instead
		zones, err := client.GetVPCZonesForRegion(ctx, installConfig.Config.Platform.IBMCloud.Region)
		if err != nil {
			return err
		}

		// Get CISInstanceCRN from InstallConfig metadata
		crn, err := installConfig.IBMCloud.CISInstanceCRN(ctx)
		if err != nil {
			return err
		}

		data, err = ibmcloudtfvars.TFVars(
			ibmcloudtfvars.TFVarsSources{
				Auth:              auth,
				CISInstanceCRN:    crn,
				PublishStrategy:   installConfig.Config.Publish,
				ResourceGroupName: installConfig.Config.Platform.IBMCloud.ResourceGroupName,

				// TODO: IBM: Fetch config from masterConfig instead
				Region:                  installConfig.Config.Platform.IBMCloud.Region,
				MachineType:             "bx2d-4x16",
				MasterAvailabilityZones: zones,
				ImageURL:                string(*rhcosImage),
			},
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
		// convert options list to a list of mappings which can be consumed by terraform
		var dnsmasqoptions []map[string]string
		for _, option := range installConfig.Config.Platform.Libvirt.Network.DnsmasqOptions {
			dnsmasqoptions = append(dnsmasqoptions,
				map[string]string{
					"option_name":  option.Name,
					"option_value": option.Value})
		}

		data, err = libvirttfvars.TFVars(
			libvirttfvars.TFVarsSources{
				MasterConfig:   masters[0].Spec.ProviderSpec.Value.Object.(*libvirtprovider.LibvirtMachineProviderConfig),
				OsImage:        string(*rhcosImage),
				MachineCIDR:    &installConfig.Config.Networking.MachineNetwork[0].CIDR.IPNet,
				Bridge:         installConfig.Config.Platform.Libvirt.Network.IfName,
				MasterCount:    masterCount,
				Architecture:   installConfig.Config.ControlPlane.Architecture,
				DnsmasqOptions: dnsmasqoptions,
			},
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	case openstack.Name:
		cloud, err := openstackconfig.GetSession(installConfig.Config.Platform.OpenStack.Cloud)
		if err != nil {
			return errors.Wrap(err, "failed to get cloud config for openstack")
		}
		var caCert string
		// Get the ca-cert-bundle key if there is a value for cacert in clouds.yaml
		if caPath := cloud.CloudConfig.CACertFile; caPath != "" {
			caFile, err := ioutil.ReadFile(caPath)
			if err != nil {
				return errors.Wrap(err, "failed to read clouds.yaml ca-cert from disk")
			}
			caCert = string(caFile)
		}

		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}

		var masterSpecs []*openstackprovider.OpenstackProviderSpec
		for _, master := range masters {
			masterSpecs = append(masterSpecs, master.Spec.ProviderSpec.Value.Object.(*openstackprovider.OpenstackProviderSpec))
		}
		data, err = openstacktfvars.TFVars(
			masterSpecs,
			installConfig.Config.Platform.OpenStack.Cloud,
			installConfig.Config.Platform.OpenStack.ExternalNetwork,
			installConfig.Config.Platform.OpenStack.ExternalDNS,
			installConfig.Config.Platform.OpenStack.APIFloatingIP,
			installConfig.Config.Platform.OpenStack.IngressFloatingIP,
			installConfig.Config.Platform.OpenStack.APIVIP,
			installConfig.Config.Platform.OpenStack.IngressVIP,
			string(*rhcosImage),
			installConfig.Config.Platform.OpenStack.ClusterOSImageProperties,
			clusterID.InfraID,
			caCert,
			bootstrapIgn,
			installConfig.Config.ControlPlane.Platform.OpenStack,
			installConfig.Config.OpenStack.DefaultMachinePlatform,
			installConfig.Config.Platform.OpenStack.MachinesSubnet,
			installConfig.Config.Proxy,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	case baremetal.Name:
		var imageCacheIP string
		if installConfig.Config.Platform.BareMetal.ProvisioningNetwork == baremetal.DisabledProvisioningNetwork {
			imageCacheIP = installConfig.Config.Platform.BareMetal.APIVIP
		} else {
			imageCacheIP = installConfig.Config.Platform.BareMetal.BootstrapProvisioningIP
		}

		data, err = baremetaltfvars.TFVars(
			installConfig.Config.Platform.BareMetal.LibvirtURI,
			installConfig.Config.Platform.BareMetal.APIVIP,
			imageCacheIP,
			string(*rhcosBootstrapImage),
			installConfig.Config.Platform.BareMetal.ExternalBridge,
			installConfig.Config.Platform.BareMetal.ExternalMACAddress,
			installConfig.Config.Platform.BareMetal.ProvisioningBridge,
			installConfig.Config.Platform.BareMetal.ProvisioningMACAddress,
			installConfig.Config.Platform.BareMetal.Hosts,
			string(*rhcosImage),
			ironicCreds.Username,
			ironicCreds.Password,
			masterIgn,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	case ovirt.Name:
		config, err := ovirtconfig.NewConfig()
		if err != nil {
			return err
		}
		con, err := ovirtconfig.NewConnection()
		if err != nil {
			return err
		}
		defer con.Close()

		if installConfig.Config.Platform.Ovirt.VNICProfileID == "" {
			profiles, err := ovirtconfig.FetchVNICProfileByClusterNetwork(
				con,
				installConfig.Config.Platform.Ovirt.ClusterID,
				installConfig.Config.Platform.Ovirt.NetworkName)
			if err != nil {
				return errors.Wrapf(err, "failed to compute values for Engine platform")
			}
			if len(profiles) != 1 {
				return errors.Wrapf(err, "failed to compute values for Engine platform, there are multiple vNIC profiles.")
			}
			installConfig.Config.Platform.Ovirt.VNICProfileID = profiles[0].MustId()
		}

		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}

		data, err := ovirttfvars.TFVars(
			ovirttfvars.Auth{
				URL:      config.URL,
				Username: config.Username,
				Password: config.Password,
				Cafile:   config.CAFile,
				Cabundle: config.CABundle,
				Insecure: config.Insecure,
			},
			installConfig.Config.Platform.Ovirt.ClusterID,
			installConfig.Config.Platform.Ovirt.StorageDomainID,
			installConfig.Config.Platform.Ovirt.NetworkName,
			installConfig.Config.Platform.Ovirt.VNICProfileID,
			string(*rhcosImage),
			clusterID.InfraID,
			masters[0].Spec.ProviderSpec.Value.Object.(*ovirtprovider.OvirtMachineProviderSpec),
			installConfig.Config.Platform.Ovirt.AffinityGroups,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	case vsphere.Name:
		controlPlanes, err := mastersAsset.Machines()
		if err != nil {
			return err
		}
		controlPlaneConfigs := make([]*vsphereprovider.VSphereMachineProviderSpec, len(controlPlanes))
		for i, c := range controlPlanes {
			controlPlaneConfigs[i] = c.Spec.ProviderSpec.Value.Object.(*vsphereprovider.VSphereMachineProviderSpec)
		}

		// Set this flag to use an existing folder specified in the install-config. Otherwise, create one.
		preexistingFolder := installConfig.Config.Platform.VSphere.Folder != ""

		data, err = vspheretfvars.TFVars(
			vspheretfvars.TFVarsSources{
				ControlPlaneConfigs: controlPlaneConfigs,
				Username:            installConfig.Config.VSphere.Username,
				Password:            installConfig.Config.VSphere.Password,
				Cluster:             installConfig.Config.VSphere.Cluster,
				ImageURL:            string(*rhcosImage),
				PreexistingFolder:   preexistingFolder,
			},
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: fmt.Sprintf(TfPlatformVarsFileName, platform),
			Data:     data,
		})
	case kubevirt.Name:
		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}
		masterSpecs := make([]*kubevirtprovider.KubevirtMachineProviderSpec, len(masters))
		for i, m := range masters {
			masterSpecs[i] = m.Spec.ProviderSpec.Value.Object.(*kubevirtprovider.KubevirtMachineProviderSpec)
		}

		labels := kubevirtutils.BuildLabels(clusterID.InfraID)
		data, err := kubevirttfvars.TFVars(
			kubevirttfvars.TFVarsSources{
				MasterSpecs:     masterSpecs,
				ImageURL:        string(*rhcosImage),
				Namespace:       installConfig.Config.Kubevirt.Namespace,
				ResourcesLabels: labels,
			},
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

	cm, err := openshiftinstall.CreateInstallConfigMap("openshift-install")
	if err != nil {
		return "", errors.Wrap(err, "failed to generate openshift-install config")
	}

	config.Storage.Files = append(config.Storage.Files, ignition.FileFromString("/opt/openshift/manifests/openshift-install.yaml", "root", 0644, cm))

	ign, err := ignition.Marshal(config)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal bootstrap Ignition config")
	}

	return string(ign), nil
}
