package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	coreosarch "github.com/coreos/stream-metadata-go/arch"
	"github.com/ghodss/yaml"
	ibmcloudprovider "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1"
	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	ovirtprovider "github.com/openshift/cluster-api-provider-ovirt/pkg/apis/ovirtprovider/v1beta1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	baremetalbootstrap "github.com/openshift/installer/pkg/asset/ignition/bootstrap/baremetal"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	aztypes "github.com/openshift/installer/pkg/asset/installconfig/azure"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	ovirtconfig "github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	vsphereconfig "github.com/openshift/installer/pkg/asset/installconfig/vsphere"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/openshiftinstall"
	"github.com/openshift/installer/pkg/asset/rhcos"
	rhcospkg "github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/tfvars"
	alibabacloudtfvars "github.com/openshift/installer/pkg/tfvars/alibabacloud"
	awstfvars "github.com/openshift/installer/pkg/tfvars/aws"
	azuretfvars "github.com/openshift/installer/pkg/tfvars/azure"
	baremetaltfvars "github.com/openshift/installer/pkg/tfvars/baremetal"
	gcptfvars "github.com/openshift/installer/pkg/tfvars/gcp"
	ibmcloudtfvars "github.com/openshift/installer/pkg/tfvars/ibmcloud"
	libvirttfvars "github.com/openshift/installer/pkg/tfvars/libvirt"
	nutanixtfvars "github.com/openshift/installer/pkg/tfvars/nutanix"
	openstacktfvars "github.com/openshift/installer/pkg/tfvars/openstack"
	ovirttfvars "github.com/openshift/installer/pkg/tfvars/ovirt"
	powervstfvars "github.com/openshift/installer/pkg/tfvars/powervs"
	vspheretfvars "github.com/openshift/installer/pkg/tfvars/vsphere"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/alibabacloud"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/types/vsphere"
)

const (
	// TfVarsFileName is the filename for Terraform variables.
	TfVarsFileName = "terraform.tfvars.json"

	// TfPlatformVarsFileName is the name for platform-specific
	// Terraform variable files.
	//
	// https://www.terraform.io/docs/configuration/variables.html#variable-files
	TfPlatformVarsFileName = "terraform.platform.auto.tfvars.json"

	tfvarsAssetName = "Terraform Variables"
)

// TerraformVariables depends on InstallConfig, Manifests,
// and Ignition to generate the terrafor.tfvars.
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
		&installconfig.PlatformProvisionCheck{},
		&manifests.Manifests{},
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
	manifestsAsset := &manifests.Manifests{}
	rhcosImage := new(rhcos.Image)
	rhcosBootstrapImage := new(rhcos.BootstrapImage)
	ironicCreds := &baremetalbootstrap.IronicCreds{}
	parents.Get(clusterID, installConfig, bootstrapIgnAsset, masterIgnAsset, mastersAsset, workersAsset, manifestsAsset, rhcosImage, rhcosBootstrapImage, ironicCreds)

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
	mastersSchedulable := false
	for _, f := range manifestsAsset.Files() {
		if f.Filename == manifests.SchedulerCfgFilename {
			schedulerConfig := configv1.Scheduler{}
			err = yaml.Unmarshal(f.Data, &schedulerConfig)
			if err != nil {
				return errors.Wrapf(err, "failed to unmarshall %s", manifests.SchedulerCfgFilename)
			}
			mastersSchedulable = schedulerConfig.Spec.MastersSchedulable
			break
		}
	}

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
		mastersSchedulable,
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
		masterConfigs := make([]*machinev1beta1.AWSMachineProviderConfig, len(masters))
		for i, m := range masters {
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*machinev1beta1.AWSMachineProviderConfig)
		}
		workers, err := workersAsset.MachineSets()
		if err != nil {
			return err
		}
		workerConfigs := make([]*machinev1beta1.AWSMachineProviderConfig, len(workers))
		for i, m := range workers {
			workerConfigs[i] = m.Spec.Template.Spec.ProviderSpec.Value.Object.(*machinev1beta1.AWSMachineProviderConfig)
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
			Proxy:                 installConfig.Config.Proxy,
		})
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: TfPlatformVarsFileName,
			Data:     data,
		})
	case azure.Name:
		session, err := installConfig.Azure.Session()
		if err != nil {
			return err
		}
		auth := azuretfvars.Auth{
			SubscriptionID:            session.Credentials.SubscriptionID,
			ClientID:                  session.Credentials.ClientID,
			ClientSecret:              session.Credentials.ClientSecret,
			TenantID:                  session.Credentials.TenantID,
			ClientCertificatePath:     session.Credentials.ClientCertificatePath,
			ClientCertificatePassword: session.Credentials.ClientCertificatePassword,
		}
		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}
		masterConfigs := make([]*machinev1beta1.AzureMachineProviderSpec, len(masters))
		for i, m := range masters {
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*machinev1beta1.AzureMachineProviderSpec)
		}
		workers, err := workersAsset.MachineSets()
		if err != nil {
			return err
		}
		workerConfigs := make([]*machinev1beta1.AzureMachineProviderSpec, len(workers))
		for i, w := range workers {
			workerConfigs[i] = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*machinev1beta1.AzureMachineProviderSpec)
		}
		client := aztypes.NewClient(session)
		hyperVGeneration, err := client.GetHyperVGenerationVersion(context.TODO(), masterConfigs[0].VMSize, masterConfigs[0].Location, "")
		if err != nil {
			return err
		}

		preexistingnetwork := installConfig.Config.Azure.VirtualNetwork != ""

		var bootstrapIgnStub, bootstrapIgnURLPlaceholder string
		if installConfig.Azure.CloudName == azure.StackCloud {
			// Due to the SAS created in Terraform to limit access to bootstrap ignition, we cannot know the URL in advance.
			// Instead, we will pass a placeholder string in the ignition to be replaced in TF once the value is known.
			bootstrapIgnURLPlaceholder = "BOOTSTRAP_IGNITION_URL_PLACEHOLDER"
			shim, err := bootstrap.GenerateIgnitionShimWithCertBundleAndProxy(bootstrapIgnURLPlaceholder, installConfig.Config.AdditionalTrustBundle, installConfig.Config.Proxy)
			if err != nil {
				return errors.Wrap(err, "failed to create stub Ignition config for bootstrap")
			}
			bootstrapIgnStub = string(shim)
		}

		data, err := azuretfvars.TFVars(
			azuretfvars.TFVarsSources{
				Auth:                            auth,
				CloudName:                       installConfig.Config.Azure.CloudName,
				ARMEndpoint:                     installConfig.Config.Azure.ARMEndpoint,
				ResourceGroupName:               installConfig.Config.Azure.ResourceGroupName,
				BaseDomainResourceGroupName:     installConfig.Config.Azure.BaseDomainResourceGroupName,
				MasterConfigs:                   masterConfigs,
				WorkerConfigs:                   workerConfigs,
				ImageURL:                        string(*rhcosImage),
				PreexistingNetwork:              preexistingnetwork,
				Publish:                         installConfig.Config.Publish,
				OutboundType:                    installConfig.Config.Azure.OutboundType,
				BootstrapIgnStub:                bootstrapIgnStub,
				BootstrapIgnitionURLPlaceholder: bootstrapIgnURLPlaceholder,
				HyperVGeneration:                hyperVGeneration,
			},
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: TfPlatformVarsFileName,
			Data:     data,
		})
	case gcp.Name:
		var publicZoneName string
		sess, err := gcpconfig.GetSession(ctx)
		if err != nil {
			return err
		}
		auth := gcptfvars.Auth{
			ProjectID:        installConfig.Config.GCP.ProjectID,
			NetworkProjectID: installConfig.Config.GCP.NetworkProjectID,
			ServiceAccount:   string(sess.Credentials.JSON),
		}

		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}
		masterConfigs := make([]*machinev1beta1.GCPMachineProviderSpec, len(masters))
		for i, m := range masters {
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*machinev1beta1.GCPMachineProviderSpec)
		}
		workers, err := workersAsset.MachineSets()
		if err != nil {
			return err
		}
		workerConfigs := make([]*machinev1beta1.GCPMachineProviderSpec, len(workers))
		for i, w := range workers {
			workerConfigs[i] = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*machinev1beta1.GCPMachineProviderSpec)
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
			Filename: TfPlatformVarsFileName,
			Data:     data,
		})
	case ibmcloud.Name:
		client, err := installConfig.IBMCloud.Client()
		if err != nil {
			return err
		}
		auth := ibmcloudtfvars.Auth{
			APIKey: client.APIKey,
		}

		// Get master and worker machine info
		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}
		masterConfigs := make([]*ibmcloudprovider.IBMCloudMachineProviderSpec, len(masters))
		for i, m := range masters {
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*ibmcloudprovider.IBMCloudMachineProviderSpec)
		}
		workers, err := workersAsset.MachineSets()
		if err != nil {
			return err
		}
		workerConfigs := make([]*ibmcloudprovider.IBMCloudMachineProviderSpec, len(workers))
		for i, w := range workers {
			workerConfigs[i] = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*ibmcloudprovider.IBMCloudMachineProviderSpec)
		}

		// Set existing network (boolean of whether one is being used)
		preexistingVPC := installConfig.Config.Platform.IBMCloud.GetVPCName() != ""

		// Set machine pool info
		var masterMachinePool ibmcloud.MachinePool
		var workerMachinePool ibmcloud.MachinePool
		if installConfig.Config.Platform.IBMCloud.DefaultMachinePlatform != nil {
			masterMachinePool.Set(installConfig.Config.Platform.IBMCloud.DefaultMachinePlatform)
			workerMachinePool.Set(installConfig.Config.Platform.IBMCloud.DefaultMachinePlatform)
		}
		if installConfig.Config.ControlPlane.Platform.IBMCloud != nil {
			masterMachinePool.Set(installConfig.Config.ControlPlane.Platform.IBMCloud)
		}
		if worker := installConfig.Config.WorkerMachinePool(); worker != nil {
			workerMachinePool.Set(worker.Platform.IBMCloud)
		}

		// Get master dedicated host info
		var masterDedicatedHosts []ibmcloudtfvars.DedicatedHost
		for _, dhost := range masterMachinePool.DedicatedHosts {
			if dhost.Name != "" {
				dh, err := client.GetDedicatedHostByName(ctx, dhost.Name, installConfig.Config.Platform.IBMCloud.Region)
				if err != nil {
					return err
				}
				masterDedicatedHosts = append(masterDedicatedHosts, ibmcloudtfvars.DedicatedHost{
					ID: *dh.ID,
				})
			} else {
				masterDedicatedHosts = append(masterDedicatedHosts, ibmcloudtfvars.DedicatedHost{
					Profile: dhost.Profile,
				})
			}
		}

		// Get worker dedicated host info
		var workerDedicatedHosts []ibmcloudtfvars.DedicatedHost
		for _, dhost := range workerMachinePool.DedicatedHosts {
			if dhost.Name != "" {
				dh, err := client.GetDedicatedHostByName(ctx, dhost.Name, installConfig.Config.Platform.IBMCloud.Region)
				if err != nil {
					return err
				}
				workerDedicatedHosts = append(workerDedicatedHosts, ibmcloudtfvars.DedicatedHost{
					ID: *dh.ID,
				})
			} else {
				workerDedicatedHosts = append(workerDedicatedHosts, ibmcloudtfvars.DedicatedHost{
					Profile: dhost.Profile,
				})
			}
		}

		var cisCRN, dnsID string

		if installConfig.Config.Publish == types.InternalPublishingStrategy {
			// Get DNSInstanceCRN from InstallConfig metadata
			dnsInstance, err := installConfig.IBMCloud.DNSInstance(ctx)
			if err != nil {
				return err
			}
			if dnsInstance != nil {
				dnsID = dnsInstance.ID
			}
		} else {
			// Get CISInstanceCRN from InstallConfig metadata
			cisCRN, err = installConfig.IBMCloud.CISInstanceCRN(ctx)
			if err != nil {
				return err
			}
		}

		data, err = ibmcloudtfvars.TFVars(
			ibmcloudtfvars.TFVarsSources{
				Auth:                 auth,
				CISInstanceCRN:       cisCRN,
				DNSInstanceID:        dnsID,
				ImageURL:             string(*rhcosImage),
				MasterConfigs:        masterConfigs,
				MasterDedicatedHosts: masterDedicatedHosts,
				PreexistingVPC:       preexistingVPC,
				PublishStrategy:      installConfig.Config.Publish,
				ResourceGroupName:    installConfig.Config.Platform.IBMCloud.ResourceGroupName,
				WorkerConfigs:        workerConfigs,
				WorkerDedicatedHosts: workerDedicatedHosts,
			},
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: TfPlatformVarsFileName,
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
			Filename: TfPlatformVarsFileName,
			Data:     data,
		})
	case openstack.Name:
		data, err = openstacktfvars.TFVars(
			installConfig,
			mastersAsset,
			workersAsset,
			string(*rhcosImage),
			clusterID,
			bootstrapIgn,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: TfPlatformVarsFileName,
			Data:     data,
		})
	case baremetal.Name:
		var imageCacheIP string
		if installConfig.Config.Platform.BareMetal.ProvisioningNetwork == baremetal.DisabledProvisioningNetwork {
			imageCacheIP = installConfig.Config.Platform.BareMetal.APIVIPs[0]
		} else {
			imageCacheIP = installConfig.Config.Platform.BareMetal.BootstrapProvisioningIP
		}

		data, err = baremetaltfvars.TFVars(
			*installConfig.Config.ControlPlane.Replicas,
			installConfig.Config.Platform.BareMetal.LibvirtURI,
			installConfig.Config.Platform.BareMetal.APIVIPs[0],
			imageCacheIP,
			string(*rhcosBootstrapImage),
			installConfig.Config.Platform.BareMetal.ExternalBridge,
			installConfig.Config.Platform.BareMetal.ExternalMACAddress,
			installConfig.Config.Platform.BareMetal.ProvisioningBridge,
			installConfig.Config.Platform.BareMetal.ProvisioningMACAddress,
			installConfig.Config.Platform.BareMetal.Hosts,
			mastersAsset.HostFiles,
			string(*rhcosImage),
			ironicCreds.Username,
			ironicCreds.Password,
			masterIgn,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: TfPlatformVarsFileName,
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
				return fmt.Errorf("failed to compute values for Engine platform, "+
					"there are multiple vNIC profiles. found %v vNIC profiles for network %s",
					len(profiles), installConfig.Config.Platform.Ovirt.NetworkName)
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
			Filename: TfPlatformVarsFileName,
			Data:     data,
		})
	case powervs.Name:
		APIKey, err := installConfig.PowerVS.APIKey(ctx)
		if err != nil {
			return err
		}

		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}

		// Get CISInstanceCRN from InstallConfig metadata
		crn, err := installConfig.PowerVS.CISInstanceCRN(ctx)
		if err != nil {
			return err
		}

		masterConfigs := make([]*machinev1.PowerVSMachineProviderConfig, len(masters))
		for i, m := range masters {
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*machinev1.PowerVSMachineProviderConfig)
		}

		var vpcSubnet string
		if len(installConfig.Config.PowerVS.VPCSubnets) > 0 {
			vpcSubnet = installConfig.Config.PowerVS.VPCSubnets[0]
		}

		osImage := strings.SplitN(string(*rhcosImage), "/", 2)
		data, err = powervstfvars.TFVars(
			powervstfvars.TFVarsSources{
				MasterConfigs:        masterConfigs,
				Region:               installConfig.Config.Platform.PowerVS.Region,
				Zone:                 installConfig.Config.Platform.PowerVS.Zone,
				APIKey:               APIKey,
				SSHKey:               installConfig.Config.SSHKey,
				PowerVSResourceGroup: installConfig.Config.PowerVS.PowerVSResourceGroup,
				ImageBucketName:      osImage[0],
				ImageBucketFileName:  osImage[1],
				NetworkName:          installConfig.Config.PowerVS.PVSNetworkName,
				VPCName:              installConfig.Config.PowerVS.VPCName,
				VPCSubnetName:        vpcSubnet,
				CloudConnectionName:  installConfig.Config.PowerVS.CloudConnectionName,
				CISInstanceCRN:       crn,
				PublishStrategy:      installConfig.Config.Publish,
			},
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: TfPlatformVarsFileName,
			Data:     data,
		})

	case vsphere.Name:
		networkZoneMap := make(map[string]string)
		var networkID string
		controlPlanes, err := mastersAsset.Machines()
		if err != nil {
			return err
		}

		controlPlaneConfigs := make([]*machinev1beta1.VSphereMachineProviderSpec, len(controlPlanes))
		for i, c := range controlPlanes {
			controlPlaneConfigs[i] = c.Spec.ProviderSpec.Value.Object.(*machinev1beta1.VSphereMachineProviderSpec)
		}

		vim25Client, _, cleanup, err := vsphereconfig.CreateVSphereClients(context.TODO(),
			installConfig.Config.VSphere.VCenter,
			installConfig.Config.VSphere.Username,
			installConfig.Config.VSphere.Password)
		if err != nil {
			return errors.Wrapf(err, "unable to connect to vCenter %s. Ensure provided information is correct and client certs have been added to system trust.", installConfig.Config.VSphere.VCenter)
		}
		defer cleanup()
		finder := vsphereconfig.NewFinder(vim25Client)

		/* Each deployment zone requires a template to be imported.
		 * The control plane machinepool might not have that zone assigned.
		 * The zone could be unused or just defined for compute machines
		 */
		for _, deploymentZone := range installConfig.Config.VSphere.DeploymentZones {
			var failureDomain vsphere.FailureDomain

			for _, fd := range installConfig.Config.VSphere.FailureDomains {
				if fd.Name == deploymentZone.FailureDomain {
					failureDomain = fd
				}
			}

			// Must use the Managed Object ID for a port group (e.g. dvportgroup-5258)
			// instead of the name since port group names aren't always unique in vSphere.
			// https://bugzilla.redhat.com/show_bug.cgi?id=1918005
			networkZoneMap[deploymentZone.Name], err = vsphereconfig.GetNetworkMoID(context.TODO(),
				vim25Client,
				finder,
				failureDomain.Topology.Datacenter,
				failureDomain.Topology.ComputeCluster,
				failureDomain.Topology.Networks[0])

			if err != nil {
				return errors.Wrap(err, "failed to get vSphere network ID")
			}
		}

		networkID, err = vsphereconfig.GetNetworkMoID(context.TODO(),
			vim25Client,
			finder,
			controlPlaneConfigs[0].Workspace.Datacenter,
			installConfig.Config.VSphere.Cluster,
			controlPlaneConfigs[0].Network.Devices[0].NetworkName)
		if err != nil {
			return errors.Wrap(err, "failed to get vSphere network ID")
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
				DiskType:            installConfig.Config.Platform.VSphere.DiskType,
				NetworkID:           networkID,

				NetworkZone:          networkZoneMap,
				InfraID:              clusterID.InfraID,
				InstallConfig:        installConfig,
				ControlPlaneMachines: controlPlanes,
			},
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: TfPlatformVarsFileName,
			Data:     data,
		})

	case alibabacloud.Name:
		client, err := installConfig.AlibabaCloud.Client()
		if err != nil {
			return errors.Wrapf(err, "failed to create new client use region %s", installConfig.Config.Platform.AlibabaCloud.Region)
		}
		bucket := fmt.Sprintf("%s-bootstrap", clusterID.InfraID)
		object := "bootstrap.ign"
		signURL, err := client.GetOSSObjectSignURL(bucket, object)
		if err != nil {
			return errors.Wrapf(err, "failed to get a presigned URL for OSS object %s", object)
		}

		auth := alibabacloudtfvars.Auth{
			AccessKey: client.AccessKeyID,
			SecretKey: client.AccessKeySecret,
		}

		masters, err := mastersAsset.Machines()
		if err != nil {
			return errors.Wrapf(err, "failed to get master machine info")
		}
		masterConfigs := make([]*machinev1.AlibabaCloudMachineProviderConfig, len(masters))
		for i, m := range masters {
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*machinev1.AlibabaCloudMachineProviderConfig)
		}
		workers, err := workersAsset.MachineSets()
		if err != nil {
			return errors.Wrapf(err, "failed to get worker machine info")
		}
		workerConfigs := make([]*machinev1.AlibabaCloudMachineProviderConfig, len(workers))
		for i, w := range workers {
			workerConfigs[i] = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*machinev1.AlibabaCloudMachineProviderConfig)
		}

		natGatewayZones, err := client.ListEnhanhcedNatGatewayAvailableZones()
		if err != nil {
			return errors.Wrapf(err, "failed to list avaliable zones for NAT gateway")
		}
		natGatewayZoneID := natGatewayZones.Zones[0].ZoneId

		vswitchIDs := []string{}
		if len(installConfig.Config.AlibabaCloud.VSwitchIDs) > 0 {
			vswitchIDs = installConfig.Config.AlibabaCloud.VSwitchIDs
		}
		data, err := alibabacloudtfvars.TFVars(
			alibabacloudtfvars.TFVarsSources{
				Auth:                  auth,
				VpcID:                 installConfig.Config.AlibabaCloud.VpcID,
				VSwitchIDs:            vswitchIDs,
				PrivateZoneID:         installConfig.Config.AlibabaCloud.PrivateZoneID,
				ResourceGroupID:       installConfig.Config.AlibabaCloud.ResourceGroupID,
				BaseDomain:            installConfig.Config.BaseDomain,
				NatGatewayZoneID:      natGatewayZoneID,
				MasterConfigs:         masterConfigs,
				WorkerConfigs:         workerConfigs,
				IgnitionBucket:        bucket,
				IgnitionPresignedURL:  signURL,
				AdditionalTrustBundle: installConfig.Config.AdditionalTrustBundle,
				Architecture:          installConfig.Config.ControlPlane.Architecture,
				Publish:               installConfig.Config.Publish,
				Proxy:                 installConfig.Config.Proxy,
			},
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: TfPlatformVarsFileName,
			Data:     data,
		})
	case nutanix.Name:
		if rhcosImage == nil {
			return errors.New("unable to retrieve rhcos image")
		}
		controlPlanes, err := mastersAsset.Machines()
		if err != nil {
			return errors.Wrapf(err, "error getting control plane machines")
		}
		controlPlaneConfigs := make([]*machinev1.NutanixMachineProviderConfig, len(controlPlanes))
		for i, c := range controlPlanes {
			controlPlaneConfigs[i] = c.Spec.ProviderSpec.Value.Object.(*machinev1.NutanixMachineProviderConfig)
		}

		imgURI := string(*rhcosImage)
		if installConfig.Config.Nutanix.ClusterOSImage != "" {
			imgURI = installConfig.Config.Nutanix.ClusterOSImage
		}
		data, err = nutanixtfvars.TFVars(
			nutanixtfvars.TFVarsSources{
				PrismCentralAddress:   installConfig.Config.Nutanix.PrismCentral.Endpoint.Address,
				Port:                  strconv.Itoa(int(installConfig.Config.Nutanix.PrismCentral.Endpoint.Port)),
				Username:              installConfig.Config.Nutanix.PrismCentral.Username,
				Password:              installConfig.Config.Nutanix.PrismCentral.Password,
				ImageURI:              imgURI,
				BootstrapIgnitionData: bootstrapIgn,
				ClusterID:             clusterID.InfraID,
				ControlPlaneConfigs:   controlPlaneConfigs,
			},
		)
		if err != nil {
			return errors.Wrapf(err, "failed to get %s Terraform variables", platform)
		}
		t.FileList = append(t.FileList, &asset.File{
			Filename: TfPlatformVarsFileName,
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

	switch file, err := f.FetchByName(TfPlatformVarsFileName); {
	case err == nil:
		t.FileList = append(t.FileList, file)
	case !os.IsNotExist(err):
		return false, err
	}

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
