package tfvars

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	ovirtprovider "github.com/openshift/cluster-api-provider-ovirt/pkg/apis/ovirtprovider/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	baremetalbootstrap "github.com/openshift/installer/pkg/asset/ignition/bootstrap/baremetal"
	gcpbootstrap "github.com/openshift/installer/pkg/asset/ignition/bootstrap/gcp"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	aztypes "github.com/openshift/installer/pkg/asset/installconfig/azure"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	ibmcloudconfig "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	ovirtconfig "github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	powervsconfig "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/openshiftinstall"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/tfvars"
	awstfvars "github.com/openshift/installer/pkg/tfvars/aws"
	azuretfvars "github.com/openshift/installer/pkg/tfvars/azure"
	baremetaltfvars "github.com/openshift/installer/pkg/tfvars/baremetal"
	gcptfvars "github.com/openshift/installer/pkg/tfvars/gcp"
	ibmcloudtfvars "github.com/openshift/installer/pkg/tfvars/ibmcloud"
	nutanixtfvars "github.com/openshift/installer/pkg/tfvars/nutanix"
	openstacktfvars "github.com/openshift/installer/pkg/tfvars/openstack"
	ovirttfvars "github.com/openshift/installer/pkg/tfvars/ovirt"
	powervstfvars "github.com/openshift/installer/pkg/tfvars/powervs"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/types/vsphere"
	ibmcloudprovider "github.com/openshift/machine-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1"
)

const (
	// GCPFirewallPermission is the role/permission to create or skip the creation of
	// firewall rules for GCP during an xpn installation.
	GCPFirewallPermission = "compute.firewalls.create"

	// TfVarsFileName is the filename for Terraform variables.
	TfVarsFileName = "terraform.tfvars.json"

	// TfPlatformVarsFileName is the name for platform-specific
	// Terraform variable files.
	//
	// https://www.terraform.io/docs/configuration/variables.html#variable-files
	TfPlatformVarsFileName = "terraform.platform.auto.tfvars.json"

	tfvarsAssetName = "Cluster Infrastructure Variables"
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

// Dependencies returns the dependency of the TerraformVariable.
func (t *TerraformVariables) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		&installconfig.InstallConfig{},
		new(rhcos.Image),
		new(rhcos.Release),
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
//
//nolint:gocyclo // legacy, pre-linter cyclomatic complexity
func (t *TerraformVariables) Generate(ctx context.Context, parents asset.Parents) error {
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	bootstrapIgnAsset := &bootstrap.Bootstrap{}
	masterIgnAsset := &machine.Master{}
	mastersAsset := &machines.Master{}
	workersAsset := &machines.Worker{}
	manifestsAsset := &manifests.Manifests{}
	rhcosImage := new(rhcos.Image)
	rhcosRelease := new(rhcos.Release)
	rhcosBootstrapImage := new(rhcos.BootstrapImage)
	ironicCreds := &baremetalbootstrap.IronicCreds{}
	parents.Get(clusterID, installConfig, bootstrapIgnAsset, masterIgnAsset, mastersAsset, workersAsset, manifestsAsset, rhcosImage, rhcosRelease, rhcosBootstrapImage, ironicCreds)

	platform := installConfig.Config.Platform.Name()
	switch platform {
	case external.Name, none.Name:
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

	lengthBootstrapFile := int64(len(bootstrapIgn))
	if installConfig.Config.Platform.Azure != nil && installConfig.Config.Platform.Azure.CustomerManagedKey != nil &&
		installConfig.Config.Platform.Azure.CustomerManagedKey.UserAssignedIdentityKey != "" {
		if lengthBootstrapFile%512 != 0 {
			lengthBootstrapFile = (((lengthBootstrapFile / 512) + 1) * 512)
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
		lengthBootstrapFile,
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

	numWorkers := int64(0)
	for _, worker := range installConfig.Config.Compute {
		numWorkers += ptr.Deref(worker.Replicas, 0)
	}

	switch platform {
	case aws.Name:
		var vpc string
		var privateSubnets []string
		var publicSubnets []string

		if len(installConfig.Config.Platform.AWS.VPC.Subnets) > 0 {
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
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*machinev1beta1.AWSMachineProviderConfig) //nolint:errcheck // legacy, pre-linter
		}
		workers, err := workersAsset.MachineSets()
		if err != nil {
			return err
		}

		workerConfigs := make([]*machinev1beta1.AWSMachineProviderConfig, len(workers))
		for i, m := range workers {
			workerConfigs[i] = m.Spec.Template.Spec.ProviderSpec.Value.Object.(*machinev1beta1.AWSMachineProviderConfig) //nolint:errcheck // legacy, pre-linter
		}
		osImage := strings.SplitN(rhcosImage.ControlPlane, ",", 2)
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

		var securityGroups []string
		if mp := installConfig.Config.AWS.DefaultMachinePlatform; mp != nil {
			securityGroups = mp.AdditionalSecurityGroupIDs
		}
		masterIAMRoleName := ""
		if mp := installConfig.Config.ControlPlane; mp != nil {
			awsMP := &aws.MachinePool{}
			awsMP.Set(installConfig.Config.AWS.DefaultMachinePlatform)
			awsMP.Set(mp.Platform.AWS)
			masterIAMRoleName = awsMP.IAMRole
			if len(awsMP.AdditionalSecurityGroupIDs) > 0 {
				securityGroups = awsMP.AdditionalSecurityGroupIDs
			}
		}

		// AWS Zones is used to determine which route table the edge zone will be associated.
		allZones, err := installConfig.AWS.AllZones(ctx)
		if err != nil {
			return err
		}

		data, err := awstfvars.TFVars(awstfvars.TFVarsSources{
			VPC:                       vpc,
			PrivateSubnets:            privateSubnets,
			PublicSubnets:             publicSubnets,
			AvailabilityZones:         allZones,
			InternalZone:              installConfig.Config.AWS.HostedZone,
			InternalZoneRole:          installConfig.Config.AWS.HostedZoneRole,
			Services:                  installConfig.Config.AWS.ServiceEndpoints,
			Publish:                   installConfig.Config.Publish,
			MasterConfigs:             masterConfigs,
			WorkerConfigs:             workerConfigs,
			AMIID:                     osImageID,
			AMIRegion:                 osImageRegion,
			IgnitionBucket:            bucket,
			IgnitionPresignedURL:      url,
			AdditionalTrustBundle:     installConfig.Config.AdditionalTrustBundle,
			MasterIAMRoleName:         masterIAMRoleName,
			WorkerIAMRoleName:         workerIAMRoleName,
			Architecture:              installConfig.Config.ControlPlane.Architecture,
			Proxy:                     installConfig.Config.Proxy,
			PreserveBootstrapIgnition: installConfig.Config.AWS.BestEffortDeleteIgnition,
			MasterSecurityGroups:      securityGroups,
			PublicIpv4Pool:            installConfig.Config.AWS.PublicIpv4Pool,
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
			UseMSI:                    session.AuthType == aztypes.ManagedIdentityAuth,
		}
		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}
		masterConfigs := make([]*machinev1beta1.AzureMachineProviderSpec, len(masters))
		for i, m := range masters {
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*machinev1beta1.AzureMachineProviderSpec) //nolint:errcheck // legacy, pre-linter
		}
		workers, err := workersAsset.MachineSets()
		if err != nil {
			return err
		}
		workerConfigs := make([]*machinev1beta1.AzureMachineProviderSpec, len(workers))
		for i, w := range workers {
			workerConfigs[i] = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*machinev1beta1.AzureMachineProviderSpec) //nolint:errcheck // legacy, pre-linter
		}
		client, err := installConfig.Azure.Client()
		if err != nil {
			return err
		}
		hyperVGeneration, err := client.GetHyperVGenerationVersion(ctx, masterConfigs[0].VMSize, masterConfigs[0].Location, "")
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

		managedKeys := azure.CustomerManagedKey{}
		if installConfig.Config.Azure.CustomerManagedKey != nil {
			managedKeys.KeyVault = azure.KeyVault{
				ResourceGroup: installConfig.Config.Azure.CustomerManagedKey.KeyVault.ResourceGroup,
				Name:          installConfig.Config.Azure.CustomerManagedKey.KeyVault.Name,
				KeyName:       installConfig.Config.Azure.CustomerManagedKey.KeyVault.KeyName,
			}
			managedKeys.UserAssignedIdentityKey = installConfig.Config.Azure.CustomerManagedKey.UserAssignedIdentityKey
		}

		lbPrivate := false
		if installConfig.Config.OperatorPublishingStrategy != nil {
			lbPrivate = installConfig.Config.OperatorPublishingStrategy.APIServer == "Internal"
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
				ImageURL:                        rhcosImage.ControlPlane,
				ImageRelease:                    rhcosRelease.GetAzureReleaseVersion(),
				PreexistingNetwork:              preexistingnetwork,
				Publish:                         installConfig.Config.Publish,
				OutboundType:                    installConfig.Config.Azure.OutboundType,
				BootstrapIgnStub:                bootstrapIgnStub,
				BootstrapIgnitionURLPlaceholder: bootstrapIgnURLPlaceholder,
				HyperVGeneration:                hyperVGeneration,
				VMArchitecture:                  installConfig.Config.ControlPlane.Architecture,
				InfrastructureName:              clusterID.InfraID,
				KeyVault:                        managedKeys.KeyVault,
				UserAssignedIdentityKey:         managedKeys.UserAssignedIdentityKey,
				LBPrivate:                       lbPrivate,
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
		sess, err := gcpconfig.GetSession(ctx)
		if err != nil {
			return err
		}

		auth := gcptfvars.Auth{
			ProjectID:        installConfig.Config.GCP.ProjectID,
			NetworkProjectID: installConfig.Config.GCP.NetworkProjectID,
			ServiceAccount:   string(sess.Credentials.JSON),
		}

		client, err := gcpconfig.NewClient(context.Background())
		if err != nil {
			return err
		}

		// In the case of a shared vpn, the firewall rules should only be created if the user has permissions to do so
		createFirewallRules := true
		if installConfig.Config.GCP.NetworkProjectID != "" {
			permissions, err := client.GetProjectPermissions(context.Background(), installConfig.Config.GCP.NetworkProjectID, []string{
				GCPFirewallPermission,
			})
			if err != nil {
				return err
			}
			createFirewallRules = permissions.Has(GCPFirewallPermission)

			if !createFirewallRules {
				logrus.Warnf("failed to find permission %s, skipping firewall rule creation", GCPFirewallPermission)
			}
		}

		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}
		masterConfigs := make([]*machinev1beta1.GCPMachineProviderSpec, len(masters))
		for i, m := range masters {
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*machinev1beta1.GCPMachineProviderSpec) //nolint:errcheck // legacy, pre-linter
		}
		workers, err := workersAsset.MachineSets()
		if err != nil {
			return err
		}
		// Based on the number of workers, we could have the following outcomes:
		// 1. compute replicas > 0, worker machinesets > 0, masters not schedulable, valid cluster
		// 2. compute replicas > 0, worker machinesets = 0, invalid cluster
		// 3. compute replicas = 0, masters schedulable, valid cluster
		if numWorkers != 0 && len(workers) == 0 {
			return fmt.Errorf("invalid configuration. No worker assets available for requested number of compute replicas (%d)", numWorkers)
		}
		if numWorkers == 0 && !mastersSchedulable {
			return fmt.Errorf("invalid configuration. No workers requested but masters are not schedulable")
		}

		workerConfigs := make([]*machinev1beta1.GCPMachineProviderSpec, len(workers))
		for i, w := range workers {
			workerConfigs[i] = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*machinev1beta1.GCPMachineProviderSpec) //nolint:errcheck // legacy, pre-linter
		}
		preexistingnetwork := installConfig.Config.GCP.Network != ""

		// Search the project for a dns zone with the specified base domain.
		// When the user has selected a custom DNS solution, the zones should be skipped.
		publicZoneName := ""
		privateZoneName := ""

		if installConfig.Config.GCP.UserProvisionedDNS != dns.UserProvisionedDNSEnabled {
			if installConfig.Config.Publish == types.ExternalPublishingStrategy {
				publicZone, err := client.GetDNSZone(ctx, installConfig.Config.GCP.ProjectID, installConfig.Config.BaseDomain, true)
				if err != nil {
					return errors.Wrapf(err, "failed to get GCP public zone")
				}
				publicZoneName = publicZone.Name
			}

			// Set the private zone
			privateZoneName, err = manifests.GetGCPPrivateZoneName(ctx, client, installConfig, clusterID.InfraID)
			if err != nil {
				return fmt.Errorf("failed to find gcp private dns zone: %w", err)
			}
		}

		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		url, err := gcpbootstrap.CreateSignedURL(clusterID.InfraID)
		if err != nil {
			return fmt.Errorf("failed to provision gcp bootstrap storage resources: %w", err)
		}

		shim, err := bootstrap.GenerateIgnitionShimWithCertBundleAndProxy(url, installConfig.Config.AdditionalTrustBundle, installConfig.Config.Proxy)
		if err != nil {
			return fmt.Errorf("failed to create gcp ignition shim: %w", err)
		}

		tags, err := gcpconfig.NewTagManager(client).GetUserTags(ctx,
			installConfig.Config.Platform.GCP.ProjectID,
			installConfig.Config.Platform.GCP.UserTags)
		if err != nil {
			return fmt.Errorf("failed to fetch user-defined tags: %w", err)
		}

		data, err := gcptfvars.TFVars(
			gcptfvars.TFVarsSources{
				Auth:                auth,
				MasterConfigs:       masterConfigs,
				WorkerConfigs:       workerConfigs,
				CreateFirewallRules: createFirewallRules,
				PreexistingNetwork:  preexistingnetwork,
				PublicZoneName:      publicZoneName,
				PrivateZoneName:     privateZoneName,
				PublishStrategy:     installConfig.Config.Publish,
				InfrastructureName:  clusterID.InfraID,
				UserProvisionedDNS:  installConfig.Config.GCP.UserProvisionedDNS == dns.UserProvisionedDNSEnabled,
				UserTags:            tags,
				IgnitionShim:        string(shim),
				PresignedURL:        url,
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
		meta := ibmcloudconfig.NewMetadata(installConfig.Config)
		client, err := meta.Client()
		if err != nil {
			return err
		}
		auth := ibmcloudtfvars.Auth{
			APIKey: client.GetAPIKey(),
		}

		// Get master and worker machine info
		masters, err := mastersAsset.Machines()
		if err != nil {
			return err
		}
		masterConfigs := make([]*ibmcloudprovider.IBMCloudMachineProviderSpec, len(masters))
		for i, m := range masters {
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*ibmcloudprovider.IBMCloudMachineProviderSpec) //nolint:errcheck // legacy, pre-linter
		}
		workers, err := workersAsset.MachineSets()
		if err != nil {
			return err
		}
		workerConfigs := make([]*ibmcloudprovider.IBMCloudMachineProviderSpec, len(workers))
		for i, w := range workers {
			workerConfigs[i] = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*ibmcloudprovider.IBMCloudMachineProviderSpec) //nolint:errcheck // legacy, pre-linter
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
		vpcPermitted := false

		if installConfig.Config.Publish == types.InternalPublishingStrategy {
			// Get DNSInstanceCRN from metadata
			dnsInstance, err := meta.DNSInstance(ctx)
			if err != nil {
				return err
			}
			if dnsInstance != nil {
				dnsID = dnsInstance.ID
			}
			// If the VPC already exists and the cluster is Private, check if the VPC is already a Permitted Network on DNS Instance
			if preexistingVPC {
				vpcPermitted, err = meta.IsVPCPermittedNetwork(ctx, installConfig.Config.Platform.IBMCloud.VPCName)
				if err != nil {
					return err
				}
			}
		} else {
			// Get CISInstanceCRN from metadata
			cisCRN, err = meta.CISInstanceCRN(ctx)
			if err != nil {
				return err
			}
		}

		// NOTE(cjschaef): If one or more ServiceEndpoint's are supplied, attempt to build the Terraform endpoint_file_path
		// https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/guides/custom-service-endpoints#file-structure-for-endpoints-file
		var endpointsJSONFile string
		// Set Terraform visibility mode if necessary
		terraformPrivateVisibility := false
		if len(installConfig.Config.Platform.IBMCloud.ServiceEndpoints) > 0 {
			// Determine if any endpoints require 'private' Terraform visibility mode (any contain 'private' or 'direct' for COS)
			// This is a requirement for the IBM Cloud Terraform provider, forcing 'public' or 'private' visibility mode.
			for _, endpoint := range installConfig.Config.Platform.IBMCloud.ServiceEndpoints {
				if strings.Contains(endpoint.URL, "private") || strings.Contains(endpoint.URL, "direct") {
					// If at least one endpoint is private (or direct) we expect to use Private visibility mode
					terraformPrivateVisibility = true
					break
				}
			}

			endpointData, err := ibmcloudtfvars.CreateEndpointJSON(installConfig.Config.Platform.IBMCloud.ServiceEndpoints, installConfig.Config.Platform.IBMCloud.Region)
			if err != nil {
				return err
			}
			// While service endpoints may not be empty, they may not be required for Terraform.
			// So, if we have not endpoint data, we don't need to generate the JSON override file.
			if endpointData != nil {
				// Add endpoint JSON data to list of generated files for Terraform
				t.FileList = append(t.FileList, &asset.File{
					Filename: ibmcloudtfvars.IBMCloudEndpointJSONFileName,
					Data:     endpointData,
				})
				endpointsJSONFile = ibmcloudtfvars.IBMCloudEndpointJSONFileName
			}
		}

		data, err = ibmcloudtfvars.TFVars(
			ibmcloudtfvars.TFVarsSources{
				Auth:                       auth,
				CISInstanceCRN:             cisCRN,
				DNSInstanceID:              dnsID,
				EndpointsJSONFile:          endpointsJSONFile,
				ImageURL:                   rhcosImage.ControlPlane,
				MasterConfigs:              masterConfigs,
				MasterDedicatedHosts:       masterDedicatedHosts,
				NetworkResourceGroupName:   installConfig.Config.Platform.IBMCloud.NetworkResourceGroupName,
				PreexistingVPC:             preexistingVPC,
				PublishStrategy:            installConfig.Config.Publish,
				ResourceGroupName:          installConfig.Config.Platform.IBMCloud.ResourceGroupName,
				TerraformPrivateVisibility: terraformPrivateVisibility,
				VPCPermitted:               vpcPermitted,
				WorkerConfigs:              workerConfigs,
				WorkerDedicatedHosts:       workerDedicatedHosts,
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
			ctx,
			installConfig,
			mastersAsset,
			workersAsset,
			rhcosImage.ControlPlane,
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
		data, err = baremetaltfvars.TFVars(
			installConfig.Config.Platform.BareMetal.LibvirtURI,
			string(*rhcosBootstrapImage),
			installConfig.Config.Platform.BareMetal.ExternalBridge,
			installConfig.Config.Platform.BareMetal.ExternalMACAddress,
			installConfig.Config.Platform.BareMetal.ProvisioningBridge,
			installConfig.Config.Platform.BareMetal.ProvisioningMACAddress,
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
			rhcosImage.ControlPlane,
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

		var (
			cisCRN, dnsCRN, vpcGatewayName, vpcSubnet string
			vpcPermitted, vpcGatewayAttached          bool
		)
		if len(installConfig.Config.PowerVS.VPCSubnets) > 0 {
			vpcSubnet = installConfig.Config.PowerVS.VPCSubnets[0]
		}
		switch installConfig.Config.Publish {
		case types.InternalPublishingStrategy:
			// Get DNSInstanceCRN from InstallConfig metadata
			dnsCRN, err = installConfig.PowerVS.DNSInstanceCRN(ctx)
			if err != nil {
				return err
			}

			// If the VPC already exists and the cluster is Private, check if the VPC is already a Permitted Network on DNS Instance
			if installConfig.Config.PowerVS.VPCName != "" {
				vpcPermitted, err = installConfig.PowerVS.IsVPCPermittedNetwork(ctx, installConfig.Config.Platform.PowerVS.VPCName, installConfig.Config.BaseDomain)
				if err != nil {
					return err
				}
				vpcGatewayName, vpcGatewayAttached, err = installConfig.PowerVS.GetExistingVPCGateway(ctx, installConfig.Config.Platform.PowerVS.VPCName, vpcSubnet)
				if err != nil {
					return err
				}
			}
		case types.ExternalPublishingStrategy:
			// Get CISInstanceCRN from InstallConfig metadata
			cisCRN, err = installConfig.PowerVS.CISInstanceCRN(ctx)
			if err != nil {
				return err
			}
		default:
			return errors.New("unknown publishing strategy")
		}

		masterConfigs := make([]*machinev1.PowerVSMachineProviderConfig, len(masters))
		for i, m := range masters {
			masterConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*machinev1.PowerVSMachineProviderConfig) //nolint:errcheck // legacy, pre-linter
		}

		client, err := powervsconfig.NewClient()
		if err != nil {
			return err
		}
		var (
			vpcRegion, vpcZone string
			vpc                *vpcv1.VPC
		)
		vpcName := installConfig.Config.PowerVS.VPCName
		if vpcName != "" {
			vpc, err = client.GetVPCByName(ctx, vpcName)
			if err != nil {
				return err
			}
			var crnElems = strings.SplitN(*vpc.CRN, ":", 8)
			vpcRegion = crnElems[5]
		} else {
			specified := installConfig.Config.PowerVS.VPCRegion
			if specified != "" {
				if powervs.ValidateVPCRegion(specified) {
					vpcRegion = specified
				} else {
					return errors.New("unknown VPC region")
				}
			} else if vpcRegion, err = powervs.VPCRegionForPowerVSRegion(installConfig.Config.PowerVS.Region); err != nil {
				return err
			}
		}
		if vpcSubnet != "" {
			var sn *vpcv1.Subnet
			sn, err = client.GetSubnetByName(ctx, vpcSubnet, vpcRegion)
			if err != nil {
				return err
			}
			vpcZone = *sn.Zone.Name
		} else {
			rand.New(rand.NewSource(time.Now().UnixNano()))           //nolint:gosec // we don't need a crypto secure number
			vpcZone = fmt.Sprintf("%s-%d", vpcRegion, rand.Intn(2)+1) //nolint:gosec // we don't need a crypto secure number
		}

		cpStanza := installConfig.Config.ControlPlane
		if cpStanza == nil || cpStanza.Platform.PowerVS == nil || cpStanza.Platform.PowerVS.SysType == "" {
			sysTypes, err := powervs.AvailableSysTypes(installConfig.Config.PowerVS.Region, installConfig.Config.PowerVS.Zone)
			if err != nil {
				return err
			}
			for i := range masters {
				masterConfigs[i].SystemType = sysTypes[0]
			}
		}

		attachedTG := ""
		tgConnectionVPCID := ""
		if installConfig.Config.PowerVS.ServiceInstanceGUID != "" {
			attachedTG, err = client.GetAttachedTransitGateway(ctx, installConfig.Config.PowerVS.ServiceInstanceGUID)
			if err != nil {
				return err
			}
			if attachedTG != "" && vpc != nil {
				tgConnectionVPCID, err = client.GetTGConnectionVPC(ctx, attachedTG, *vpc.ID)
				if err != nil {
					return err
				}
			}
		}

		// If a service instance GUID was passed in the install-config.yaml file, then
		// find the corresponding name for it.  Otherwise, we expect our Terraform to
		// dynamically create one.
		serviceInstanceName, err := client.ServiceInstanceGUIDToName(ctx, installConfig.Config.PowerVS.ServiceInstanceGUID)
		if err != nil {
			return err
		}

		osImage := strings.SplitN(rhcosImage.ControlPlane, "/", 2)
		data, err = powervstfvars.TFVars(
			powervstfvars.TFVarsSources{
				MasterConfigs:          masterConfigs,
				Region:                 installConfig.Config.Platform.PowerVS.Region,
				Zone:                   installConfig.Config.Platform.PowerVS.Zone,
				APIKey:                 APIKey,
				SSHKey:                 installConfig.Config.SSHKey,
				PowerVSResourceGroup:   installConfig.Config.PowerVS.PowerVSResourceGroup,
				ImageBucketName:        osImage[0],
				ImageBucketFileName:    osImage[1],
				VPCRegion:              vpcRegion,
				VPCZone:                vpcZone,
				VPCName:                vpcName,
				VPCSubnetName:          vpcSubnet,
				VPCPermitted:           vpcPermitted,
				VPCGatewayName:         vpcGatewayName,
				VPCGatewayAttached:     vpcGatewayAttached,
				CISInstanceCRN:         cisCRN,
				DNSInstanceCRN:         dnsCRN,
				PublishStrategy:        installConfig.Config.Publish,
				EnableSNAT:             len(installConfig.Config.DeprecatedImageContentSources) == 0 && len(installConfig.Config.ImageDigestSources) == 0,
				AttachedTransitGateway: attachedTG,
				TGConnectionVPCID:      tgConnectionVPCID,
				ServiceInstanceName:    serviceInstanceName,
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
		t.FileList = make([]*asset.File, 0)
		return nil
	case nutanix.Name:
		controlPlanes, err := mastersAsset.Machines()
		if err != nil {
			return errors.Wrapf(err, "error getting control plane machines")
		}
		controlPlaneConfigs := make([]*machinev1.NutanixMachineProviderConfig, len(controlPlanes))
		for i, c := range controlPlanes {
			controlPlaneConfigs[i] = c.Spec.ProviderSpec.Value.Object.(*machinev1.NutanixMachineProviderConfig) //nolint:errcheck // legacy, pre-linter
		}

		imgURI := rhcosImage.ControlPlane
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
