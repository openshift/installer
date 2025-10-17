package machines

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	baremetalhost "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1beta1"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/api/features"
	machinev1 "github.com/openshift/api/machine/v1"
	machinev1alpha1 "github.com/openshift/api/machine/v1alpha1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	baremetalapi "github.com/openshift/cluster-api-provider-baremetal/pkg/apis"
	baremetalprovider "github.com/openshift/cluster-api-provider-baremetal/pkg/apis/baremetal/v1alpha1"
	libvirtapi "github.com/openshift/cluster-api-provider-libvirt/pkg/apis"
	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	ovirtproviderapi "github.com/openshift/cluster-api-provider-ovirt/pkg/apis"
	ovirtprovider "github.com/openshift/cluster-api-provider-ovirt/pkg/apis/ovirtprovider/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines/aws"
	"github.com/openshift/installer/pkg/asset/machines/azure"
	"github.com/openshift/installer/pkg/asset/machines/baremetal"
	"github.com/openshift/installer/pkg/asset/machines/gcp"
	"github.com/openshift/installer/pkg/asset/machines/ibmcloud"
	"github.com/openshift/installer/pkg/asset/machines/machineconfig"
	"github.com/openshift/installer/pkg/asset/machines/nutanix"
	"github.com/openshift/installer/pkg/asset/machines/openstack"
	"github.com/openshift/installer/pkg/asset/machines/ovirt"
	"github.com/openshift/installer/pkg/asset/machines/powervs"
	"github.com/openshift/installer/pkg/asset/machines/vsphere"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/asset/rhcos"
	rhcosutils "github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	azuredefaults "github.com/openshift/installer/pkg/types/azure/defaults"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	externaltypes "github.com/openshift/installer/pkg/types/external"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	nonetypes "github.com/openshift/installer/pkg/types/none"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
	powervsdefaults "github.com/openshift/installer/pkg/types/powervs/defaults"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
	ibmcloudapi "github.com/openshift/machine-api-provider-ibmcloud/pkg/apis"
	ibmcloudprovider "github.com/openshift/machine-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1"
)

// Master generates the machines for the `master` machine pool.
type Master struct {
	UserDataFile           *asset.File
	MachineConfigFiles     []*asset.File
	MachineFiles           []*asset.File
	ControlPlaneMachineSet *asset.File
	IPClaimFiles           []*asset.File
	IPAddrFiles            []*asset.File

	// SecretFiles is used by the baremetal platform to register the
	// credential information for communicating with management
	// controllers on hosts.
	SecretFiles []*asset.File

	// NetworkConfigSecretFiles is used by the baremetal platform to
	// store the networking configuration per host
	NetworkConfigSecretFiles []*asset.File

	// HostFiles is the list of baremetal hosts provided in the
	// installer configuration.
	HostFiles []*asset.File

	// FencingCredentialsSecretFiles is a collection of secrets
	// that will be used to fence machines to enable safe recovery
	// in Two Node OpenShift with Fencing (TNF) deployments
	FencingCredentialsSecretFiles []*asset.File
}

const (
	directory = "openshift"

	// secretFileName is the format string for constructing the Secret
	// filenames for baremetal clusters.
	secretFileName = "99_openshift-cluster-api_host-bmc-secrets-%s.yaml"

	// networkConfigSecretFileName is the format string for constructing
	// the networking configuration Secret filenames for baremetal
	// clusters.
	networkConfigSecretFileName = "99_openshift-cluster-api_host-network-config-secrets-%s.yaml"

	// hostFileName is the format string for constucting the Host
	// filenames for baremetal clusters.
	hostFileName = "99_openshift-cluster-api_hosts-%s.yaml"

	// masterMachineFileName is the format string for constucting the
	// master Machine filenames.
	masterMachineFileName = "99_openshift-cluster-api_master-machines-%s.yaml"

	// masterUserDataFileName is the filename used for the master
	// user-data secret.
	masterUserDataFileName = "99_openshift-cluster-api_master-user-data-secret.yaml"

	// controlPlaneMachineSetFileName is the filename used for the control plane machine sets.
	controlPlaneMachineSetFileName = "99_openshift-machine-api_master-control-plane-machine-set.yaml"

	// ipClaimFileName is the filename used for the ip claims list.
	ipClaimFileName = "99_openshift-machine-api_claim-%s.yaml"

	// ipAddressFileName is the filename used for the ip addresses list.
	ipAddressFileName = "99_openshift-machine-api_address-%s.yaml"

	// fencingCredentialsSecretFileName is the format string for constructing the
	// secret filenames for Two Node OpenShift with Fencing (TNF) clusters.
	fencingCredentialsSecretFileName = "99_openshift-etcd_fencing-credentials-secrets-%s.yaml" // #nosec G101
)

var (
	secretFileNamePattern              = fmt.Sprintf(secretFileName, "*")
	networkConfigSecretFileNamePattern = fmt.Sprintf(networkConfigSecretFileName, "*")
	hostFileNamePattern                = fmt.Sprintf(hostFileName, "*")
	masterMachineFileNamePattern       = fmt.Sprintf(masterMachineFileName, "*")
	masterIPClaimFileNamePattern       = fmt.Sprintf(ipClaimFileName, "*master*")
	masterIPAddressFileNamePattern     = fmt.Sprintf(ipAddressFileName, "*master*")
	fencingCredentialsFilenamePattern  = fmt.Sprintf(fencingCredentialsSecretFileName, "*")

	_ asset.WritableAsset = (*Master)(nil)
)

// Name returns a human friendly name for the Master Asset.
func (m *Master) Name() string {
	return "Master Machines"
}

// Dependencies returns all of the dependencies directly needed by the
// Master asset
func (m *Master) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		// PlatformCredsCheck just checks the creds (and asks, if needed)
		// We do not actually use it in this asset directly, hence
		// it is put in the dependencies but not fetched in Generate
		&installconfig.PlatformCredsCheck{},
		&installconfig.InstallConfig{},
		new(rhcos.Image),
		&machine.Master{},
	}
}

// Generate generates the Master asset.
//
//nolint:gocyclo
func (m *Master) Generate(ctx context.Context, dependencies asset.Parents) error {
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	rhcosImage := new(rhcos.Image)
	mign := &machine.Master{}
	dependencies.Get(clusterID, installConfig, rhcosImage, mign)

	masterUserDataSecretName := "master-user-data"

	ic := installConfig.Config

	pool := *ic.ControlPlane
	var err error
	machines := []machinev1beta1.Machine{}
	var ipClaims []ipamv1.IPAddressClaim
	var ipAddrs []ipamv1.IPAddress
	var controlPlaneMachineSet *machinev1.ControlPlaneMachineSet

	// Check if SNO topology is supported on this platform
	if pool.Replicas != nil && *pool.Replicas == 1 {
		bootstrapInPlace := false
		if ic.BootstrapInPlace != nil && ic.BootstrapInPlace.InstallationDisk != "" {
			bootstrapInPlace = true
		}
		if !supportedSingleNodePlatform(bootstrapInPlace, ic.Platform.Name()) {
			return fmt.Errorf("this install method does not support Single Node installation on platform %s", ic.Platform.Name())
		}
	}
	switch ic.Platform.Name() {
	case awstypes.Name:
		subnets, err := aws.MachineSubnetsByZones(ctx, installConfig, awstypes.ClusterNodeSubnetRole)
		if err != nil {
			return err
		}

		mpool := defaultAWSMachinePoolPlatform("master")

		osImage := strings.SplitN(rhcosImage.ControlPlane, ",", 2)
		osImageID := osImage[0]
		if len(osImage) == 2 {
			osImageID = "" // the AMI will be generated later on
		}
		mpool.AMIID = osImageID

		mpool.Set(ic.Platform.AWS.DefaultMachinePlatform)
		mpool.Set(pool.Platform.AWS)
		zoneDefaults := false
		if len(mpool.Zones) == 0 {
			if len(subnets) > 0 {
				for zone := range subnets {
					mpool.Zones = append(mpool.Zones, zone)
				}
				// Since zones are extracted from map keys, order is not guaranteed.
				// Thus, sort the zones by lexical order to ensure CAPI and MAPI machines
				// are distributed to zones in the same order.
				slices.Sort(mpool.Zones)
			} else {
				mpool.Zones, err = installConfig.AWS.AvailabilityZones(ctx)
				if err != nil {
					return err
				}
				zoneDefaults = true
			}
		}

		if mpool.InstanceType == "" {
			topology := configv1.HighlyAvailableTopologyMode
			if pool.Replicas != nil && *pool.Replicas == 1 {
				topology = configv1.SingleReplicaTopologyMode
			}
			mpool.InstanceType, err = aws.PreferredInstanceType(ctx, installConfig.AWS, awsdefaults.InstanceTypes(installConfig.Config.Platform.AWS.Region, installConfig.Config.ControlPlane.Architecture, topology), mpool.Zones)
			if err != nil {
				logrus.Warn(errors.Wrap(err, "failed to find default instance type"))
				mpool.InstanceType = awsdefaults.InstanceTypes(installConfig.Config.Platform.AWS.Region, installConfig.Config.ControlPlane.Architecture, topology)[0]
			}
		}

		// if the list of zones is the default we need to try to filter the list in case there are some zones where the instance might not be available
		if zoneDefaults {
			mpool.Zones, err = aws.FilterZonesBasedOnInstanceType(ctx, installConfig.AWS, mpool.InstanceType, mpool.Zones)
			if err != nil {
				logrus.Warn(errors.Wrap(err, "failed to filter zone list"))
			}
		}

		pool.Platform.AWS = &mpool
		machines, controlPlaneMachineSet, err = aws.Machines(
			clusterID.InfraID,
			installConfig.Config.Platform.AWS.Region,
			subnets,
			&pool,
			"master",
			masterUserDataSecretName,
			installConfig.Config.Platform.AWS.UserTags,
			awstypes.IsPublicOnlySubnetsEnabled(),
		)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		aws.ConfigMasters(machines, controlPlaneMachineSet, clusterID.InfraID, ic.Publish)
	case gcptypes.Name:
		mpool := defaultGCPMachinePoolPlatform(pool.Architecture)
		mpool.Set(ic.Platform.GCP.DefaultMachinePlatform)
		mpool.Set(pool.Platform.GCP)
		if len(mpool.Zones) == 0 {
			azs, err := gcp.ZonesForInstanceType(ic.Platform.GCP.ProjectID, ic.Platform.GCP.Region, mpool.InstanceType, ic.Platform.GCP.ServiceEndpoints)
			if err != nil {
				return errors.Wrap(err, "failed to fetch availability zones")
			}
			mpool.Zones = azs
		}
		pool.Platform.GCP = &mpool
		machines, controlPlaneMachineSet, err = gcp.Machines(clusterID.InfraID, ic, &pool, rhcosImage.ControlPlane, "master", masterUserDataSecretName)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		err := gcp.ConfigMasters(machines, controlPlaneMachineSet, clusterID.InfraID, ic.Publish)
		if err != nil {
			return err
		}

		// CAPG-based installs will use only backend services--no target pools,
		// so we don't want to include target pools in the control plane machineset.
		// TODO(padillon): once this feature gate is the default and we are
		// no longer using Terraform, we can update ConfigMasters not to populate this.
		if capiutils.IsEnabled(installConfig) {
			for _, machine := range machines {
				providerSpec, ok := machine.Spec.ProviderSpec.Value.Object.(*machinev1beta1.GCPMachineProviderSpec)
				if !ok {
					return errors.New("unable to convert ProviderSpec to GCPMachineProviderSpec")
				}
				providerSpec.TargetPools = nil
			}
			cpms := controlPlaneMachineSet.Spec.Template.OpenShiftMachineV1Beta1Machine.Spec.ProviderSpec.Value.Object
			providerSpec, ok := cpms.(*machinev1beta1.GCPMachineProviderSpec)
			if !ok {
				return errors.New("Unable to set target pools to control plane machine set")
			}
			providerSpec.TargetPools = nil
		}
	case ibmcloudtypes.Name:
		subnets := map[string]string{}
		if len(ic.Platform.IBMCloud.ControlPlaneSubnets) > 0 {
			subnetMetas, err := installConfig.IBMCloud.ControlPlaneSubnets(ctx)
			if err != nil {
				return err
			}
			for _, subnet := range subnetMetas {
				subnets[subnet.Zone] = subnet.Name
			}
		}
		mpool := defaultIBMCloudMachinePoolPlatform()
		mpool.Set(ic.Platform.IBMCloud.DefaultMachinePlatform)
		mpool.Set(pool.Platform.IBMCloud)
		if len(mpool.Zones) == 0 {
			azs, err := ibmcloud.AvailabilityZones(ic.Platform.IBMCloud.Region, ic.Platform.IBMCloud.ServiceEndpoints)
			if err != nil {
				return errors.Wrap(err, "failed to fetch availability zones")
			}
			mpool.Zones = azs
		}
		pool.Platform.IBMCloud = &mpool
		machines, err = ibmcloud.Machines(clusterID.InfraID, ic, subnets, &pool, "master", masterUserDataSecretName)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		// TODO: IBM: implement ConfigMasters() if needed
		// ibmcloud.ConfigMasters(machines, clusterID.InfraID, ic.Publish)
	case openstacktypes.Name:
		mpool := defaultOpenStackMachinePoolPlatform()
		mpool.Set(ic.Platform.OpenStack.DefaultMachinePlatform)
		mpool.Set(pool.Platform.OpenStack)
		pool.Platform.OpenStack = &mpool

		imageName, _ := rhcosutils.GenerateOpenStackImageName(rhcosImage.ControlPlane, clusterID.InfraID)

		machines, controlPlaneMachineSet, err = openstack.Machines(ctx, clusterID.InfraID, ic, &pool, imageName, "master", masterUserDataSecretName)
		if err != nil {
			return fmt.Errorf("failed to create master machine objects: %w", err)
		}
		openstack.ConfigMasters(machines, clusterID.InfraID)
	case azuretypes.Name:
		mpool := defaultAzureMachinePoolPlatform(installConfig.Config.Platform.Azure.CloudName)
		mpool.InstanceType = azuredefaults.ControlPlaneInstanceType(
			installConfig.Config.Platform.Azure.CloudName,
			installConfig.Config.Platform.Azure.Region,
			installConfig.Config.ControlPlane.Architecture,
		)
		mpool.OSDisk.DiskSizeGB = 1024
		if installConfig.Config.Platform.Azure.CloudName == azuretypes.StackCloud {
			mpool.OSDisk.DiskSizeGB = azuredefaults.AzurestackMinimumDiskSize
		}
		mpool.Set(ic.Platform.Azure.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Azure)

		client, err := installConfig.Azure.Client()
		if err != nil {
			return err
		}

		if len(mpool.Zones) == 0 {
			azs, err := client.GetAvailabilityZones(ctx, ic.Platform.Azure.Region, mpool.InstanceType)
			if err != nil {
				return errors.Wrap(err, "failed to fetch availability zones")
			}
			mpool.Zones = azs
			if len(azs) == 0 {
				// if no azs are given we set to []string{""} for convenience over later operations.
				// It means no-zoned for the machine API
				mpool.Zones = []string{""}
			}
		}

		if mpool.OSImage.Publisher != "" {
			img, ierr := client.GetMarketplaceImage(ctx, ic.Platform.Azure.Region, mpool.OSImage.Publisher, mpool.OSImage.Offer, mpool.OSImage.SKU, mpool.OSImage.Version)
			if ierr != nil {
				return fmt.Errorf("failed to fetch marketplace image: %w", ierr)
			}
			// Publisher is case-sensitive and matched against exactly. Also the
			// Plan's publisher might not be exactly the same as the Image's
			// publisher
			if img.Plan != nil && img.Plan.Publisher != nil {
				mpool.OSImage.Publisher = *img.Plan.Publisher
			}
		}
		pool.Platform.Azure = &mpool

		capabilities, err := client.GetVMCapabilities(ctx, mpool.InstanceType, installConfig.Config.Platform.Azure.Region)
		if err != nil {
			return err
		}
		session, err := installConfig.Azure.Session()
		if err != nil {
			return err
		}
		useImageGallery := installConfig.Azure.CloudName != azuretypes.StackCloud
		machines, controlPlaneMachineSet, err = azure.Machines(clusterID.InfraID, ic, &pool, rhcosImage.ControlPlane, "master", masterUserDataSecretName, capabilities, useImageGallery, session)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		err = azure.ConfigMasters(machines, controlPlaneMachineSet, clusterID.InfraID)
		if err != nil {
			return err
		}
	case baremetaltypes.Name:
		mpool := defaultBareMetalMachinePoolPlatform()
		mpool.Set(ic.Platform.BareMetal.DefaultMachinePlatform)
		mpool.Set(pool.Platform.BareMetal)
		pool.Platform.BareMetal = &mpool

		// Use managed user data secret, since we always have up to date images
		// available in the cluster
		masterUserDataSecretName = "master-user-data-managed"
		enabledCaps := installConfig.Config.GetEnabledCapabilities()
		if !enabledCaps.Has(configv1.ClusterVersionCapabilityMachineAPI) {
			break
		}
		machines, err = baremetal.Machines(clusterID.InfraID, ic, &pool, "master", masterUserDataSecretName)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}

		hostSettings, err := baremetal.Hosts(ic, machines, masterUserDataSecretName)
		if err != nil {
			return errors.Wrap(err, "failed to assemble host data")
		}

		hosts, err := createHostAssetFiles(hostSettings.Hosts, hostFileName)
		if err != nil {
			return err
		}
		m.HostFiles = append(m.HostFiles, hosts...)

		secrets, err := createSecretAssetFiles(hostSettings.Secrets, secretFileName)
		if err != nil {
			return err
		}
		m.SecretFiles = append(m.SecretFiles, secrets...)

		networkSecrets, err := createSecretAssetFiles(hostSettings.NetworkConfigSecrets, networkConfigSecretFileName)
		if err != nil {
			return err
		}
		m.NetworkConfigSecretFiles = append(m.NetworkConfigSecretFiles, networkSecrets...)

	case ovirttypes.Name:
		mpool := defaultOvirtMachinePoolPlatform()
		mpool.VMType = ovirttypes.VMTypeHighPerformance
		mpool.Set(ic.Platform.Ovirt.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Ovirt)
		pool.Platform.Ovirt = &mpool

		imageName, _ := rhcosutils.GenerateOpenStackImageName(rhcosImage.ControlPlane, clusterID.InfraID)

		machines, err = ovirt.Machines(clusterID.InfraID, ic, &pool, imageName, "master", masterUserDataSecretName)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects for ovirt provider")
		}
	case vspheretypes.Name:
		mpool := defaultVSphereMachinePoolPlatform()
		mpool.NumCPUs = 4
		mpool.NumCoresPerSocket = 4
		mpool.MemoryMiB = 16384
		mpool.Set(ic.Platform.VSphere.DefaultMachinePlatform)
		mpool.Set(pool.Platform.VSphere)

		// The machinepool has no zones defined, there are FailureDomains
		// This is a vSphere zonal installation. Generate machinepool zone
		// list.

		fdCount := int64(len(ic.Platform.VSphere.FailureDomains))
		var idx int64
		if len(mpool.Zones) == 0 && len(ic.VSphere.FailureDomains) != 0 {
			for i := int64(0); i < *(ic.ControlPlane.Replicas); i++ {
				idx = i
				if idx >= fdCount {
					idx = i % fdCount
				}
				mpool.Zones = append(mpool.Zones, ic.VSphere.FailureDomains[idx].Name)
			}
		}

		pool.Platform.VSphere = &mpool

		data, err := vsphere.Machines(clusterID.InfraID, ic, &pool, "master", masterUserDataSecretName)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		machines = data.Machines
		controlPlaneMachineSet = data.ControlPlaneMachineSet
		ipClaims = data.IPClaims
		ipAddrs = data.IPAddresses

		vsphere.ConfigMasters(machines, clusterID.InfraID)
	case powervstypes.Name:
		mpool := defaultPowerVSMachinePoolPlatform(ic)
		mpool.Set(ic.Platform.PowerVS.DefaultMachinePlatform)
		mpool.Set(pool.Platform.PowerVS)
		// Only the service instance is guaranteed to exist and be passed via the install config
		// The other two, we should standardize a name including the cluster id. At this point, all
		// we have are names.
		pool.Platform.PowerVS = &mpool
		machines, controlPlaneMachineSet, err = powervs.Machines(clusterID.InfraID, ic, &pool, "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		if err := powervs.ConfigMasters(machines, controlPlaneMachineSet, clusterID.InfraID, ic.Publish); err != nil {
			return errors.Wrap(err, "failed to to configure master machine objects")
		}
	case externaltypes.Name, nonetypes.Name:
	case nutanixtypes.Name:
		mpool := defaultNutanixMachinePoolPlatform()
		mpool.NumCPUs = 8
		mpool.Set(ic.Platform.Nutanix.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Nutanix)
		if err = mpool.ValidateConfig(ic.Platform.Nutanix, "master"); err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		pool.Platform.Nutanix = &mpool
		templateName := nutanixtypes.RHCOSImageName(ic.Platform.Nutanix, clusterID.InfraID)

		machines, controlPlaneMachineSet, err = nutanix.Machines(clusterID.InfraID, ic, &pool, templateName, "master", masterUserDataSecretName)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		nutanix.ConfigMasters(machines, clusterID.InfraID)
	default:
		return fmt.Errorf("invalid Platform")
	}

	data, err := UserDataSecret(masterUserDataSecretName, mign.File.Data)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secret for master machines")
	}

	m.UserDataFile = &asset.File{
		Filename: filepath.Join(directory, masterUserDataFileName),
		Data:     data,
	}

	machineConfigs := []*mcfgv1.MachineConfig{}
	if pool.Hyperthreading == types.HyperthreadingDisabled {
		ignHT, err := machineconfig.ForHyperthreadingDisabled("master")
		if err != nil {
			return errors.Wrap(err, "failed to create ignition for hyperthreading disabled for master machines")
		}
		machineConfigs = append(machineConfigs, ignHT)
	}
	if ic.SSHKey != "" {
		ignSSH, err := machineconfig.ForAuthorizedKeys(ic.SSHKey, "master")
		if err != nil {
			return errors.Wrap(err, "failed to create ignition for authorized SSH keys for master machines")
		}
		machineConfigs = append(machineConfigs, ignSSH)
	}
	if ic.FIPS {
		ignFIPS, err := machineconfig.ForFIPSEnabled("master")
		if err != nil {
			return errors.Wrap(err, "failed to create ignition for FIPS enabled for master machines")
		}
		machineConfigs = append(machineConfigs, ignFIPS)
	}
	if ic.Platform.Name() == powervstypes.Name {
		// always enable multipath for powervs.
		ignMultipath, err := machineconfig.ForMultipathEnabled("master")
		if err != nil {
			return errors.Wrap(err, "failed to create ignition for Multipath enabled for master machines")
		}
		machineConfigs = append(machineConfigs, ignMultipath)

		// set SMT level if specified for powervs.
		if pool.Platform.PowerVS.SMTLevel != "" {
			ignPowerSMT, err := machineconfig.ForPowerSMT("master", pool.Platform.PowerVS.SMTLevel)
			if err != nil {
				return errors.Wrap(err, "failed to create ignition for Power SMT for master machines")
			}
			machineConfigs = append(machineConfigs, ignPowerSMT)
		}

		if installConfig.Config.Publish == types.InternalPublishingStrategy &&
			(len(installConfig.Config.ImageDigestSources) > 0 || len(installConfig.Config.DeprecatedImageContentSources) > 0) {
			ignChrony, err := machineconfig.ForCustomNTP("master", powervsdefaults.DefaultNTPServer)
			if err != nil {
				return errors.Wrap(err, "failed to create ignition for custom NTP for master machines")
			}
			machineConfigs = append(machineConfigs, ignChrony)

			ignRoutes, err := machineconfig.ForExtraRoutes("master", powervsdefaults.DefaultExtraRoutes(), ic.MachineNetwork[0].CIDR.String())
			if err != nil {
				return errors.Wrap(err, "failed to create ignition for extra routes for master machines")
			}
			machineConfigs = append(machineConfigs, ignRoutes)
		}
	}
	// The maximum number of networks supported on ServiceNetwork is two, one IPv4 and one IPv6 network.
	// The cluster-network-operator handles the validation of this field.
	// Reference: https://github.com/openshift/cluster-network-operator/blob/fc3e0e25b4cfa43e14122bdcdd6d7f2585017d75/pkg/network/cluster_config.go#L45-L52
	if ic.Networking != nil && len(ic.Networking.ServiceNetwork) == 2 {
		// Only configure kernel args for dual-stack clusters.
		ignIPv6, err := machineconfig.ForDualStackAddresses("master")
		if err != nil {
			return errors.Wrap(err, "failed to create ignition to configure IPv6 for master machines")
		}
		machineConfigs = append(machineConfigs, ignIPv6)
	}

	if installConfig.Config.EnabledFeatureGates().Enabled(features.FeatureGateMultiDiskSetup) {
		for i, diskSetup := range installConfig.Config.ControlPlane.DiskSetup {
			var dataDisk any
			var diskName string

			switch diskSetup.Type {
			case types.Etcd:
				diskName = diskSetup.Etcd.PlatformDiskID
			case types.Swap:
				diskName = diskSetup.Etcd.PlatformDiskID
			case types.UserDefined:
				diskName = diskSetup.UserDefined.PlatformDiskID
			default:
				// We shouldn't get here, but just in case
				return errors.Errorf("disk setup type %s is not supported", diskSetup.Type)
			}

			switch ic.Platform.Name() {
			case azuretypes.Name:
				azureControlPlaneMachinePool := ic.ControlPlane.Platform.Azure

				if i < len(azureControlPlaneMachinePool.DataDisks) {
					dataDisk = azureControlPlaneMachinePool.DataDisks[i]
				}
			case vspheretypes.Name:
				vsphereControlPlaneMachinePool := ic.ControlPlane.Platform.VSphere
				for index, disk := range vsphereControlPlaneMachinePool.DataDisks {
					if disk.Name == diskName {
						dataDisk = vsphere.DiskInfo{
							Index: index,
							Disk:  disk,
						}
						break
					}
				}
			default:
				return errors.Errorf("disk setup for %s is not supported", ic.Platform.Name())
			}

			if dataDisk != nil {
				diskSetupIgn, err := NodeDiskSetup(installConfig, "master", diskSetup, dataDisk)
				if err != nil {
					return errors.Wrap(err, "failed to create ignition to setup disks for control plane")
				}
				machineConfigs = append(machineConfigs, diskSetupIgn)
			}
		}
	}

	m.MachineConfigFiles, err = machineconfig.Manifests(machineConfigs, "master", directory)
	if err != nil {
		return errors.Wrap(err, "failed to create MachineConfig manifests for master machines")
	}

	m.MachineFiles = make([]*asset.File, len(machines))
	if controlPlaneMachineSet != nil && *pool.Replicas > 1 {
		data, err := yaml.Marshal(controlPlaneMachineSet)
		if err != nil {
			return errors.Wrapf(err, "marshal control plane machine set")
		}
		m.ControlPlaneMachineSet = &asset.File{
			Filename: filepath.Join(directory, controlPlaneMachineSetFileName),
			Data:     data,
		}
	}

	m.IPClaimFiles = make([]*asset.File, len(ipClaims))
	for i, claim := range ipClaims {
		data, err := yaml.Marshal(claim)
		if err != nil {
			return errors.Wrapf(err, "unable to marshal ip claim %v", claim.Name)
		}

		m.IPClaimFiles[i] = &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(ipClaimFileName, claim.Name)),
			Data:     data,
		}
	}

	m.IPAddrFiles = make([]*asset.File, len(ipAddrs))
	for i, address := range ipAddrs {
		data, err := yaml.Marshal(address)
		if err != nil {
			return errors.Wrapf(err, "unable to marshal ip claim %v", address.Name)
		}

		m.IPAddrFiles[i] = &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(ipAddressFileName, address.Name)),
			Data:     data,
		}
	}

	padFormat := fmt.Sprintf("%%0%dd", len(fmt.Sprintf("%d", len(machines))))
	for i, machine := range machines {
		data, err := yaml.Marshal(machine)
		if err != nil {
			return errors.Wrapf(err, "marshal master %d", i)
		}

		padded := fmt.Sprintf(padFormat, i)
		m.MachineFiles[i] = &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(masterMachineFileName, padded)),
			Data:     data,
		}
	}

	// This is only used by Two Node OpenShift with Fencing (TNF)
	// The credentials are rendered into secrets that will be consumed by the
	// cluster-etcd operator to enable recovery via fencing
	if pool.Fencing != nil && len(pool.Fencing.Credentials) > 0 {
		credentials, err := gatherFencingCredentials(pool.Fencing.Credentials)
		if err != nil {
			return err
		}

		secrets, err := createSecretAssetFiles(credentials, fencingCredentialsSecretFileName)
		if err != nil {
			return fmt.Errorf("failed to gather fencing credentials for control plane hosts: %w", err)
		}
		m.FencingCredentialsSecretFiles = append(m.FencingCredentialsSecretFiles, secrets...)
	}

	return nil
}

func gatherFencingCredentials(credentials []*types.Credential) ([]corev1.Secret, error) {
	secrets := []corev1.Secret{}

	for _, credential := range credentials {
		secret := &corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Secret",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("fencing-credentials-%s", credential.HostName),
				Namespace: "openshift-etcd",
			},
			Data: map[string][]byte{
				"username":                []byte(credential.Username),
				"password":                []byte(credential.Password),
				"address":                 []byte(credential.Address),
				"certificateVerification": []byte(credential.CertificateVerification),
			},
		}

		secrets = append(secrets, *secret)
	}

	return secrets, nil
}

// Files returns the files generated by the asset.
func (m *Master) Files() []*asset.File {
	files := make([]*asset.File, 0, 1+len(m.MachineConfigFiles)+len(m.MachineFiles))
	if m.UserDataFile != nil {
		files = append(files, m.UserDataFile)
	}
	files = append(files, m.MachineConfigFiles...)
	// Hosts refer to secrets, so place the secrets before the hosts
	// to avoid unnecessary reconciliation errors.
	files = append(files, m.SecretFiles...)
	files = append(files, m.NetworkConfigSecretFiles...)
	// Machines are linked to hosts via the machineRef, so we create
	// the hosts first to ensure if the operator starts trying to
	// reconcile a machine it can pick up the related host.
	files = append(files, m.HostFiles...)
	files = append(files, m.MachineFiles...)
	if m.ControlPlaneMachineSet != nil {
		files = append(files, m.ControlPlaneMachineSet)
	}
	files = append(files, m.IPClaimFiles...)
	files = append(files, m.IPAddrFiles...)
	files = append(files, m.FencingCredentialsSecretFiles...)

	return files
}

// Load reads the asset files from disk.
func (m *Master) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(filepath.Join(directory, masterUserDataFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	m.UserDataFile = file

	m.MachineConfigFiles, err = machineconfig.Load(f, "master", directory)
	if err != nil {
		return true, err
	}

	var fileList []*asset.File

	fileList, err = f.FetchByPattern(filepath.Join(directory, secretFileNamePattern))
	if err != nil {
		return true, err
	}
	m.SecretFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, networkConfigSecretFileNamePattern))
	if err != nil {
		return true, err
	}
	m.NetworkConfigSecretFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, hostFileNamePattern))
	if err != nil {
		return true, err
	}
	m.HostFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, masterMachineFileNamePattern))
	if err != nil {
		return true, err
	}
	m.MachineFiles = fileList

	file, err = f.FetchByName(filepath.Join(directory, controlPlaneMachineSetFileName))
	if err != nil {
		// Throw the error only if the file was present, since UPI and baremetal
		// deployments do not use CPMS. We ignore this file if it's missing.
		if !os.IsNotExist(err) {
			return true, err
		}

		logrus.Debugf("CPMS file missing. Ignoring it while loading machine asset.")
	}
	m.ControlPlaneMachineSet = file

	fileList, err = f.FetchByPattern(filepath.Join(directory, masterIPClaimFileNamePattern))
	if err != nil {
		return true, err
	}
	m.IPClaimFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, masterIPAddressFileNamePattern))
	if err != nil {
		return true, err
	}
	m.IPAddrFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, fencingCredentialsFilenamePattern))
	if err != nil {
		// Throw the error only if the file was present, since fencing credentials
		// are apply only to Two Node OpenShift with Fencing (TNF) deployments.
		// All other deployments will ignore this file.
		if !os.IsNotExist(err) {
			return true, err
		}
	}
	m.FencingCredentialsSecretFiles = fileList

	return true, nil
}

// Machines returns master Machine manifest structures.
func (m *Master) Machines() ([]machinev1beta1.Machine, error) {
	scheme := runtime.NewScheme()
	baremetalapi.AddToScheme(scheme)
	ibmcloudapi.AddToScheme(scheme)
	libvirtapi.AddToScheme(scheme)
	ovirtproviderapi.AddToScheme(scheme)
	scheme.AddKnownTypes(machinev1alpha1.GroupVersion,
		&machinev1alpha1.OpenstackProviderSpec{},
	)
	scheme.AddKnownTypes(machinev1beta1.SchemeGroupVersion,
		&machinev1beta1.AWSMachineProviderConfig{},
		&machinev1beta1.VSphereMachineProviderSpec{},
		&machinev1beta1.AzureMachineProviderSpec{},
		&machinev1beta1.GCPMachineProviderSpec{},
	)
	scheme.AddKnownTypes(machinev1.GroupVersion,
		&machinev1.NutanixMachineProviderConfig{},
		&machinev1.PowerVSMachineProviderConfig{},
		&machinev1.ControlPlaneMachineSet{},
	)

	machinev1beta1.AddToScheme(scheme)
	machinev1.Install(scheme)
	decoder := serializer.NewCodecFactory(scheme).UniversalDecoder(
		machinev1.GroupVersion,
		baremetalprovider.SchemeGroupVersion,
		ibmcloudprovider.SchemeGroupVersion,
		libvirtprovider.SchemeGroupVersion,
		machinev1alpha1.GroupVersion,
		machinev1beta1.SchemeGroupVersion,
		ovirtprovider.SchemeGroupVersion,
	)

	machines := []machinev1beta1.Machine{}
	for i, file := range m.MachineFiles {
		machine := &machinev1beta1.Machine{}
		err := yaml.Unmarshal(file.Data, &machine)
		if err != nil {
			return machines, errors.Wrapf(err, "unmarshal master %d", i)
		}

		obj, _, err := decoder.Decode(machine.Spec.ProviderSpec.Value.Raw, nil, nil)
		if err != nil {
			return machines, errors.Wrapf(err, "unmarshal master %d", i)
		}

		machine.Spec.ProviderSpec.Value = &runtime.RawExtension{Object: obj}
		machines = append(machines, *machine)
	}

	return machines, nil
}

// IsMachineManifest tests whether a file is a manifest that belongs to the
// Master Machines or Worker Machines asset.
func IsMachineManifest(file *asset.File) bool {
	if filepath.Dir(file.Filename) != directory {
		return false
	}
	filename := filepath.Base(file.Filename)
	if filename == masterUserDataFileName || filename == workerUserDataFileName || filename == controlPlaneMachineSetFileName {
		return true
	}
	if matched, err := machineconfig.IsManifest(filename); err != nil {
		panic(err)
	} else if matched {
		return true
	}
	for _, pattern := range []struct {
		Pattern string
		Type    string
	}{
		{Pattern: secretFileNamePattern, Type: "secret"},
		{Pattern: networkConfigSecretFileNamePattern, Type: "network config secret"},
		{Pattern: hostFileNamePattern, Type: "host"},
		{Pattern: masterMachineFileNamePattern, Type: "master machine"},
		{Pattern: workerMachineSetFileNamePattern, Type: "worker machineset"},
		{Pattern: masterIPAddressFileNamePattern, Type: "master ip address"},
		{Pattern: masterIPClaimFileNamePattern, Type: "master ip address claim"},
		{Pattern: fencingCredentialsFilenamePattern, Type: "fencing credentials secret"},
	} {
		if matched, err := filepath.Match(pattern.Pattern, filename); err != nil {
			panic(fmt.Sprintf("bad format for %s file name pattern", pattern.Type))
		} else if matched {
			return true
		}
	}
	return false
}

// IPAddresses returns IPAddress manifest structures.
func (m *Master) IPAddresses() ([]ipamv1.IPAddress, error) {
	ipAddresses := []ipamv1.IPAddress{}
	for i, file := range m.IPAddrFiles {
		logrus.Debugf("Attempting to load address %v.", file.Filename)
		address := &ipamv1.IPAddress{}
		err := yaml.Unmarshal(file.Data, &address)
		if err != nil {
			return ipAddresses, errors.Wrapf(err, "unable to unmarshal ip address %d", i)
		}

		ipAddresses = append(ipAddresses, *address)
	}

	return ipAddresses, nil
}

func createSecretAssetFiles(resources []corev1.Secret, fileName string) ([]*asset.File, error) {

	var objects []interface{}
	for _, r := range resources {
		objects = append(objects, r)
	}

	return createAssetFiles(objects, fileName)
}

func createHostAssetFiles(resources []baremetalhost.BareMetalHost, fileName string) ([]*asset.File, error) {

	var objects []interface{}
	for _, r := range resources {
		objects = append(objects, r)
	}

	return createAssetFiles(objects, fileName)
}

func createAssetFiles(objects []interface{}, fileName string) ([]*asset.File, error) {

	assetFiles := make([]*asset.File, len(objects))
	padFormat := fmt.Sprintf("%%0%dd", len(fmt.Sprintf("%d", len(objects))))
	for i, obj := range objects {
		data, err := yaml.Marshal(obj)
		if err != nil {
			return nil, errors.Wrapf(err, "marshal resource %d", i)
		}
		padded := fmt.Sprintf(padFormat, i)
		assetFiles[i] = &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(fileName, padded)),
			Data:     data,
		}
	}

	return assetFiles, nil
}

// IsFencingCredentialsFile is used during the creation of the ignition files
// to override its file mode to use a reduced permission set.
func IsFencingCredentialsFile(filepath string) (bool, error) {
	match, err := regexp.MatchString(fmt.Sprintf(fencingCredentialsSecretFileName, "[0-9]+"), filepath)
	if err != nil {
		return false, err
	}
	return match, nil
}

// supportedSingleNodePlatform indicates if the IPI Installer can be used to install SNO on
// a platform.
func supportedSingleNodePlatform(bootstrapInPlace bool, platformName string) bool {
	switch platformName {
	case awstypes.Name, gcptypes.Name, azuretypes.Name, powervstypes.Name, nonetypes.Name, ibmcloudtypes.Name:
		// Single node OpenShift installations supported without `bootstrapInPlace`
		return true
	case externaltypes.Name:
		// Single node OpenShift installations supported with `bootstrapInPlace`
		return bootstrapInPlace
	default:
		return false
	}
}
