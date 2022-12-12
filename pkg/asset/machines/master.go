package machines

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	baremetalhost "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1"
	machinev1alpha1 "github.com/openshift/api/machine/v1alpha1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	baremetalapi "github.com/openshift/cluster-api-provider-baremetal/pkg/apis"
	baremetalprovider "github.com/openshift/cluster-api-provider-baremetal/pkg/apis/baremetal/v1alpha1"
	libvirtapi "github.com/openshift/cluster-api-provider-libvirt/pkg/apis"
	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	ovirtproviderapi "github.com/openshift/cluster-api-provider-ovirt/pkg/apis"
	ovirtprovider "github.com/openshift/cluster-api-provider-ovirt/pkg/apis/ovirtprovider/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/asset/machines/alibabacloud"
	"github.com/openshift/installer/pkg/asset/machines/aws"
	"github.com/openshift/installer/pkg/asset/machines/azure"
	"github.com/openshift/installer/pkg/asset/machines/baremetal"
	"github.com/openshift/installer/pkg/asset/machines/gcp"
	"github.com/openshift/installer/pkg/asset/machines/ibmcloud"
	"github.com/openshift/installer/pkg/asset/machines/libvirt"
	"github.com/openshift/installer/pkg/asset/machines/machineconfig"
	"github.com/openshift/installer/pkg/asset/machines/nutanix"
	"github.com/openshift/installer/pkg/asset/machines/openstack"
	"github.com/openshift/installer/pkg/asset/machines/ovirt"
	"github.com/openshift/installer/pkg/asset/machines/powervs"
	"github.com/openshift/installer/pkg/asset/machines/vsphere"
	"github.com/openshift/installer/pkg/asset/rhcos"
	rhcosutils "github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	alibabacloudtypes "github.com/openshift/installer/pkg/types/alibabacloud"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	azuredefaults "github.com/openshift/installer/pkg/types/azure/defaults"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	libvirttypes "github.com/openshift/installer/pkg/types/libvirt"
	nonetypes "github.com/openshift/installer/pkg/types/none"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
	ibmcloudapi "github.com/openshift/machine-api-provider-ibmcloud/pkg/apis"
	ibmcloudprovider "github.com/openshift/machine-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
)

// Master generates the machines for the `master` machine pool.
type Master struct {
	UserDataFile           *asset.File
	MachineConfigFiles     []*asset.File
	MachineFiles           []*asset.File
	ControlPlaneMachineSet *asset.File

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

	// masterUserDataFileName is the filename used for the control plane machine sets.
	controlPlaneMachineSetFileName = "99_openshift-machine-api_master-control-plane-machine-set.yaml"
)

var (
	secretFileNamePattern              = fmt.Sprintf(secretFileName, "*")
	networkConfigSecretFileNamePattern = fmt.Sprintf(networkConfigSecretFileName, "*")
	hostFileNamePattern                = fmt.Sprintf(hostFileName, "*")
	masterMachineFileNamePattern       = fmt.Sprintf(masterMachineFileName, "*")

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
func (m *Master) Generate(dependencies asset.Parents) error {
	ctx := context.TODO()
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
	var controlPlaneMachineSet *machinev1.ControlPlaneMachineSet
	switch ic.Platform.Name() {
	case alibabacloudtypes.Name:
		client, err := installConfig.AlibabaCloud.Client()
		if err != nil {
			return err
		}
		vswitchMaps, err := installConfig.AlibabaCloud.VSwitchMaps()
		if err != nil {
			return errors.Wrap(err, "failed to get VSwitchs map")
		}
		mpool := alibabacloudtypes.DefaultMasterMachinePoolPlatform()
		mpool.ImageID = string(*rhcosImage)
		mpool.Set(ic.Platform.AlibabaCloud.DefaultMachinePlatform)
		mpool.Set(pool.Platform.AlibabaCloud)
		if len(mpool.Zones) == 0 {
			if len(vswitchMaps) > 0 {
				for zone := range vswitchMaps {
					mpool.Zones = append(mpool.Zones, zone)
				}
			} else {
				azs, err := client.GetAvailableZonesByInstanceType(mpool.InstanceType)
				if err != nil || len(azs) == 0 {
					return errors.Wrap(err, "failed to fetch availability zones")
				}
				mpool.Zones = azs
			}
		}

		pool.Platform.AlibabaCloud = &mpool
		machines, err = alibabacloud.Machines(clusterID.InfraID, ic, &pool, "master", masterUserDataSecretName, installConfig.Config.Platform.AlibabaCloud.Tags, vswitchMaps)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
	case awstypes.Name:
		subnets := map[string]string{}
		if len(ic.Platform.AWS.Subnets) > 0 {
			subnetMeta, err := installConfig.AWS.PrivateSubnets(ctx)
			if err != nil {
				return err
			}
			for id, subnet := range subnetMeta {
				subnets[subnet.Zone] = id
			}
		}

		mpool := defaultAWSMachinePoolPlatform("master")

		osImage := strings.SplitN(string(*rhcosImage), ",", 2)
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
		)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		aws.ConfigMasters(machines, controlPlaneMachineSet, clusterID.InfraID, ic.Publish)
	case gcptypes.Name:
		mpool := defaultGCPMachinePoolPlatform()
		mpool.Set(ic.Platform.GCP.DefaultMachinePlatform)
		mpool.Set(pool.Platform.GCP)
		if len(mpool.Zones) == 0 {
			azs, err := gcp.AvailabilityZones(ic.Platform.GCP.ProjectID, ic.Platform.GCP.Region)
			if err != nil {
				return errors.Wrap(err, "failed to fetch availability zones")
			}
			mpool.Zones = azs
		}
		pool.Platform.GCP = &mpool
		machines, controlPlaneMachineSet, err = gcp.Machines(clusterID.InfraID, ic, &pool, string(*rhcosImage), "master", masterUserDataSecretName)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		err := gcp.ConfigMasters(machines, controlPlaneMachineSet, clusterID.InfraID, ic.Publish)
		if err != nil {
			return err
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
			azs, err := ibmcloud.AvailabilityZones(ic.Platform.IBMCloud.Region)
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
	case libvirttypes.Name:
		mpool := defaultLibvirtMachinePoolPlatform()
		mpool.Set(ic.Platform.Libvirt.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Libvirt)
		pool.Platform.Libvirt = &mpool
		machines, err = libvirt.Machines(clusterID.InfraID, ic, &pool, "master", masterUserDataSecretName)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
	case openstacktypes.Name:
		mpool := defaultOpenStackMachinePoolPlatform()
		mpool.Set(ic.Platform.OpenStack.DefaultMachinePlatform)
		mpool.Set(pool.Platform.OpenStack)
		pool.Platform.OpenStack = &mpool

		imageName, _ := rhcosutils.GenerateOpenStackImageName(string(*rhcosImage), clusterID.InfraID)

		machines, err = openstack.Machines(clusterID.InfraID, ic, &pool, imageName, "master", masterUserDataSecretName)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		openstack.ConfigMasters(machines, clusterID.InfraID)
	case azuretypes.Name:
		mpool := defaultAzureMachinePoolPlatform()
		mpool.InstanceType = azuredefaults.ControlPlaneInstanceType(
			installConfig.Config.Platform.Azure.CloudName,
			installConfig.Config.Platform.Azure.Region,
			installConfig.Config.ControlPlane.Architecture,
		)
		mpool.OSDisk.DiskSizeGB = 1024
		mpool.Set(ic.Platform.Azure.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Azure)

		session, err := installConfig.Azure.Session()
		if err != nil {
			return errors.Wrap(err, "failed to fetch session")
		}

		// Default to current subscription if one was not specified
		if mpool.OSDisk.DiskEncryptionSet != nil && mpool.OSDisk.DiskEncryptionSet.SubscriptionID == "" {
			mpool.OSDisk.DiskEncryptionSet.SubscriptionID = session.Credentials.SubscriptionID
		}

		client := icazure.NewClient(session)
		if len(mpool.Zones) == 0 {
			azs, err := client.GetAvailabilityZones(context.TODO(), ic.Platform.Azure.Region, mpool.InstanceType)
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

		pool.Platform.Azure = &mpool

		capabilities, err := client.GetVMCapabilities(context.TODO(), mpool.InstanceType, installConfig.Config.Platform.Azure.Region)
		if err != nil {
			return err
		}
		useImageGallery := installConfig.Azure.CloudName != azuretypes.StackCloud
		machines, controlPlaneMachineSet, err = azure.Machines(clusterID.InfraID, ic, &pool, string(*rhcosImage), "master", masterUserDataSecretName, capabilities, useImageGallery)
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
		machines, err = baremetal.Machines(clusterID.InfraID, ic, &pool, "master", masterUserDataSecretName)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}

		hostSettings, err := baremetal.Hosts(ic, machines)
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

		imageName, _ := rhcosutils.GenerateOpenStackImageName(string(*rhcosImage), clusterID.InfraID)

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
		templateName := clusterID.InfraID + "-rhcos"

		machines, err = vsphere.Machines(clusterID.InfraID, ic, &pool, templateName, "master", masterUserDataSecretName)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		vsphere.ConfigMasters(machines, clusterID.InfraID)
	case powervstypes.Name:
		mpool := defaultPowerVSMachinePoolPlatform()
		mpool.Set(ic.Platform.PowerVS.DefaultMachinePlatform)
		mpool.Set(pool.Platform.PowerVS)
		// Only the service instance is guaranteed to exist and be passed via the install config
		// The other two, we should standardize a name including the cluster id. At this point, all
		// we have are names.
		pool.Platform.PowerVS = &mpool
		machines, err = powervs.Machines(clusterID.InfraID, ic, &pool, "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		powervs.ConfigMasters(machines, clusterID.InfraID)
	case nonetypes.Name:
	case nutanixtypes.Name:
		mpool := defaultNutanixMachinePoolPlatform()
		mpool.NumCPUs = 8
		mpool.Set(ic.Platform.Nutanix.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Nutanix)
		pool.Platform.Nutanix = &mpool
		templateName := nutanixtypes.RHCOSImageName(clusterID.InfraID)

		machines, err = nutanix.Machines(clusterID.InfraID, ic, &pool, templateName, "master", masterUserDataSecretName)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		nutanix.ConfigMasters(machines, clusterID.InfraID)
	default:
		return fmt.Errorf("invalid Platform")
	}

	data, err := userDataSecret(masterUserDataSecretName, mign.File.Data)
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
	return nil
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
		if os.IsNotExist(err) {
			// Choosing to ignore the CPMS file if it does not exist since UPI does not need it.
			logrus.Debugf("CPMS file missing. Ignoring it while loading machine asset.")
			return true, nil
		}
		return true, err
	}
	m.ControlPlaneMachineSet = file

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
		&machinev1.AlibabaCloudMachineProviderConfig{},
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
	if matched, err := filepath.Match(masterMachineFileNamePattern, filename); err != nil {
		panic("bad format for master machine file name pattern")
	} else if matched {
		return true
	}
	if matched, err := filepath.Match(workerMachineSetFileNamePattern, filename); err != nil {
		panic("bad format for worker machine file name pattern")
	} else {
		return matched
	}
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
