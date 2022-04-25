package machines

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	machinev1 "github.com/openshift/api/machine/v1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	baremetalapi "github.com/openshift/cluster-api-provider-baremetal/pkg/apis"
	baremetalprovider "github.com/openshift/cluster-api-provider-baremetal/pkg/apis/baremetal/v1alpha1"
	ibmcloudapi "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis"
	ibmcloudprovider "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1"
	libvirtapi "github.com/openshift/cluster-api-provider-libvirt/pkg/apis"
	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	ovirtproviderapi "github.com/openshift/cluster-api-provider-ovirt/pkg/apis"
	ovirtprovider "github.com/openshift/cluster-api-provider-ovirt/pkg/apis/ovirtprovider/v1beta1"
	nutanixapi "github.com/openshift/machine-api-provider-nutanix/pkg/apis"
	nutanixprovider "github.com/openshift/machine-api-provider-nutanix/pkg/apis/nutanixprovider/v1beta1"
	powervsapi "github.com/openshift/machine-api-provider-powervs/pkg/apis"
	powervsprovider "github.com/openshift/machine-api-provider-powervs/pkg/apis/powervsprovider/v1alpha1"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	openstackapi "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis"
	openstackprovider "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
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
	azclient "github.com/openshift/installer/pkg/client/azure"
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
)

const (
	// workerMachineSetFileName is the format string for constructing the worker MachineSet filenames.
	workerMachineSetFileName = "99_openshift-cluster-api_worker-machineset-%s.yaml"

	// workerUserDataFileName is the filename used for the worker user-data secret.
	workerUserDataFileName = "99_openshift-cluster-api_worker-user-data-secret.yaml"

	// decimalRootVolumeSize is the size in GB we use for some platforms.
	// See below.
	decimalRootVolumeSize = 120

	// powerOfTwoRootVolumeSize is the size in GB we use for other platforms.
	// The reasons for the specific choices between these two may boil down
	// to which section of code the person adding a platform was copy-pasting from.
	// https://github.com/openshift/openshift-docs/blob/main/modules/installation-requirements-user-infra.adoc#minimum-resource-requirements
	powerOfTwoRootVolumeSize = 128
)

var (
	workerMachineSetFileNamePattern = fmt.Sprintf(workerMachineSetFileName, "*")

	_ asset.WritableAsset = (*Worker)(nil)
)

func defaultAWSMachinePoolPlatform() awstypes.MachinePool {
	return awstypes.MachinePool{
		EC2RootVolume: awstypes.EC2RootVolume{
			Type: "gp3",
			Size: decimalRootVolumeSize,
		},
	}
}

func defaultLibvirtMachinePoolPlatform() libvirttypes.MachinePool {
	return libvirttypes.MachinePool{}
}

func defaultAzureMachinePoolPlatform() azuretypes.MachinePool {
	return azuretypes.MachinePool{
		OSDisk: azuretypes.OSDisk{
			DiskSizeGB: powerOfTwoRootVolumeSize,
			DiskType:   azuretypes.DefaultDiskType,
		},
	}
}

func defaultGCPMachinePoolPlatform() gcptypes.MachinePool {
	return gcptypes.MachinePool{
		InstanceType: "n1-standard-4",
		OSDisk: gcptypes.OSDisk{
			DiskSizeGB: powerOfTwoRootVolumeSize,
			DiskType:   "pd-ssd",
		},
	}
}

func defaultIBMCloudMachinePoolPlatform() ibmcloudtypes.MachinePool {
	return ibmcloudtypes.MachinePool{
		InstanceType: "bx2-4x16",
	}
}

func defaultOpenStackMachinePoolPlatform() openstacktypes.MachinePool {
	return openstacktypes.MachinePool{
		Zones: []string{""},
	}
}

func defaultBareMetalMachinePoolPlatform() baremetaltypes.MachinePool {
	return baremetaltypes.MachinePool{}
}

func defaultOvirtMachinePoolPlatform() ovirttypes.MachinePool {
	return ovirttypes.MachinePool{
		CPU: &ovirttypes.CPU{
			Cores:   4,
			Sockets: 1,
		},
		MemoryMB: 16348,
		OSDisk: &ovirttypes.Disk{
			SizeGB: decimalRootVolumeSize,
		},
		VMType:            ovirttypes.VMTypeServer,
		AutoPinningPolicy: ovirttypes.AutoPinningNone,
	}
}

func defaultVSphereMachinePoolPlatform() vspheretypes.MachinePool {
	return vspheretypes.MachinePool{
		NumCPUs:           2,
		NumCoresPerSocket: 2,
		MemoryMiB:         8192,
		OSDisk: vspheretypes.OSDisk{
			DiskSizeGB: decimalRootVolumeSize,
		},
	}
}

func defaultPowerVSMachinePoolPlatform() powervstypes.MachinePool {
	return powervstypes.MachinePool{
		Memory:     "32",
		Processors: "0.5",
		ProcType:   powervstypes.Shared,
		SysType:    "s922",
	}
}

func defaultNutanixMachinePoolPlatform() nutanixtypes.MachinePool {
	return nutanixtypes.MachinePool{
		NumCPUs:           2,
		NumCoresPerSocket: 2,
		MemoryMiB:         8192,
		OSDisk: nutanixtypes.OSDisk{
			DiskSizeMiB: decimalRootVolumeSize * 1024,
		},
	}
}

func awsDefaultWorkerMachineTypes(region string, arch types.Architecture) []string {
	classes := awsdefaults.InstanceClasses(region, arch)
	types := make([]string, len(classes))
	for i, c := range classes {
		types[i] = fmt.Sprintf("%s.large", c)
	}
	return types
}

// Worker generates the machinesets for `worker` machine pool.
type Worker struct {
	UserDataFile       *asset.File
	MachineConfigFiles []*asset.File
	MachineSetFiles    []*asset.File
}

// Name returns a human friendly name for the Worker Asset.
func (w *Worker) Name() string {
	return "Worker Machines"
}

// Dependencies returns all of the dependencies directly needed by the
// Worker asset
func (w *Worker) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		// PlatformCredsCheck just checks the creds (and asks, if needed)
		// We do not actually use it in this asset directly, hence
		// it is put in the dependencies but not fetched in Generate
		&installconfig.PlatformCredsCheck{},
		&installconfig.InstallConfig{},
		new(rhcos.Image),
		&machine.Worker{},
	}
}

// Generate generates the Worker asset.
func (w *Worker) Generate(dependencies asset.Parents) error {
	ctx := context.TODO()
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	rhcosImage := new(rhcos.Image)
	wign := &machine.Worker{}
	dependencies.Get(clusterID, installConfig, rhcosImage, wign)

	workerUserDataSecretName := "worker-user-data"

	machineConfigs := []*mcfgv1.MachineConfig{}
	machineSets := []runtime.Object{}
	var err error
	ic := installConfig.Config
	for _, pool := range ic.Compute {
		if pool.Hyperthreading == types.HyperthreadingDisabled {
			ignHT, err := machineconfig.ForHyperthreadingDisabled("worker")
			if err != nil {
				return errors.Wrap(err, "failed to create ignition for hyperthreading disabled for worker machines")
			}
			machineConfigs = append(machineConfigs, ignHT)
		}
		if ic.SSHKey != "" {
			ignSSH, err := machineconfig.ForAuthorizedKeys(ic.SSHKey, "worker")
			if err != nil {
				return errors.Wrap(err, "failed to create ignition for authorized SSH keys for worker machines")
			}
			machineConfigs = append(machineConfigs, ignSSH)
		}
		if ic.FIPS {
			ignFIPS, err := machineconfig.ForFIPSEnabled("worker")
			if err != nil {
				return errors.Wrap(err, "failed to create ignition for FIPS enabled for worker machines")
			}
			machineConfigs = append(machineConfigs, ignFIPS)
		}
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

			mpool := alibabacloudtypes.DefaultWorkerMachinePoolPlatform()
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
			sets, err := alibabacloud.MachineSets(
				clusterID.InfraID,
				ic,
				&pool,
				"worker",
				workerUserDataSecretName,
				installConfig.Config.Platform.AlibabaCloud.Tags,
				vswitchMaps,
			)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
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

			mpool := defaultAWSMachinePoolPlatform()

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
				mpool.InstanceType, err = aws.PreferredInstanceType(ctx, installConfig.AWS, awsDefaultWorkerMachineTypes(installConfig.Config.Platform.AWS.Region, installConfig.Config.ControlPlane.Architecture), mpool.Zones)
				if err != nil {
					logrus.Warn(errors.Wrap(err, "failed to find default instance type"))
					mpool.InstanceType = awsDefaultWorkerMachineTypes(installConfig.Config.Platform.AWS.Region, installConfig.Config.ControlPlane.Architecture)[0]
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
			sets, err := aws.MachineSets(
				clusterID.InfraID,
				installConfig.Config.Platform.AWS.Region,
				subnets,
				&pool,
				"worker",
				workerUserDataSecretName,
				installConfig.Config.Platform.AWS.UserTags,
			)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case azuretypes.Name:
			mpool := defaultAzureMachinePoolPlatform()
			mpool.InstanceType = azuredefaults.ComputeInstanceType(
				installConfig.Config.Platform.Azure.CloudName,
				installConfig.Config.Platform.Azure.Region,
			)
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

			if len(mpool.Zones) == 0 {
				azs, err := azure.AvailabilityZones(session, ic.Platform.Azure.Region, mpool.InstanceType)
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
			client := azclient.NewClient(session)
			hyperVGen, err := client.GetHyperVGenerationVersion(context.TODO(), mpool.InstanceType, mpool.OSDisk.DiskType, installConfig.Config.Platform.Azure.Region)
			if err != nil {
				return err
			}

			sets, err := azure.MachineSets(clusterID.InfraID, ic, &pool, string(*rhcosImage), "worker", workerUserDataSecretName, hyperVGen)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case baremetaltypes.Name:
			mpool := defaultBareMetalMachinePoolPlatform()
			mpool.Set(ic.Platform.BareMetal.DefaultMachinePlatform)
			mpool.Set(pool.Platform.BareMetal)
			pool.Platform.BareMetal = &mpool

			// Use managed user data secret, since images used by MachineSet
			// are always up to date
			workerUserDataSecretName = "worker-user-data-managed"
			sets, err := baremetal.MachineSets(clusterID.InfraID, ic, &pool, "", "worker", workerUserDataSecretName)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
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
			sets, err := gcp.MachineSets(clusterID.InfraID, ic, &pool, string(*rhcosImage), "worker", workerUserDataSecretName)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case ibmcloudtypes.Name:
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
			sets, err := ibmcloud.MachineSets(clusterID.InfraID, ic, &pool, "worker", workerUserDataSecretName)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case libvirttypes.Name:
			mpool := defaultLibvirtMachinePoolPlatform()
			mpool.Set(ic.Platform.Libvirt.DefaultMachinePlatform)
			mpool.Set(pool.Platform.Libvirt)
			pool.Platform.Libvirt = &mpool
			sets, err := libvirt.MachineSets(clusterID.InfraID, ic, &pool, "worker", workerUserDataSecretName)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case openstacktypes.Name:
			mpool := defaultOpenStackMachinePoolPlatform()
			mpool.Set(ic.Platform.OpenStack.DefaultMachinePlatform)
			mpool.Set(pool.Platform.OpenStack)
			pool.Platform.OpenStack = &mpool

			imageName, _ := rhcosutils.GenerateOpenStackImageName(string(*rhcosImage), clusterID.InfraID)

			sets, err := openstack.MachineSets(clusterID.InfraID, ic, &pool, imageName, "worker", workerUserDataSecretName, nil)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case vspheretypes.Name:
			mpool := defaultVSphereMachinePoolPlatform()
			mpool.Set(ic.Platform.VSphere.DefaultMachinePlatform)
			mpool.Set(pool.Platform.VSphere)
			pool.Platform.VSphere = &mpool
			templateName := clusterID.InfraID + "-rhcos"

			sets, err := vsphere.MachineSets(clusterID.InfraID, ic, &pool, templateName, "worker", workerUserDataSecretName)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case ovirttypes.Name:
			mpool := defaultOvirtMachinePoolPlatform()
			mpool.Set(ic.Platform.Ovirt.DefaultMachinePlatform)
			mpool.Set(pool.Platform.Ovirt)
			pool.Platform.Ovirt = &mpool

			imageName, _ := rhcosutils.GenerateOpenStackImageName(string(*rhcosImage), clusterID.InfraID)

			sets, err := ovirt.MachineSets(clusterID.InfraID, ic, &pool, imageName, "worker", workerUserDataSecretName)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects for ovirt provider")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case powervstypes.Name:
			mpool := defaultPowerVSMachinePoolPlatform()
			mpool.Set(ic.Platform.PowerVS.DefaultMachinePlatform)
			mpool.Set(pool.Platform.PowerVS)
			pool.Platform.PowerVS = &mpool
			sets, err := powervs.MachineSets(clusterID.InfraID, ic, &pool, "worker", "worker-user-data")
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects for powervs provider")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case nonetypes.Name:
		case nutanixtypes.Name:
			mpool := defaultNutanixMachinePoolPlatform()
			mpool.Set(ic.Platform.Nutanix.DefaultMachinePlatform)
			mpool.Set(pool.Platform.Nutanix)
			pool.Platform.Nutanix = &mpool
			imageName := nutanixtypes.RHCOSImageName(clusterID.InfraID)

			sets, err := nutanix.MachineSets(clusterID.InfraID, ic, &pool, imageName, "worker", workerUserDataSecretName)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		default:
			return fmt.Errorf("invalid Platform")
		}
	}

	data, err := userDataSecret(workerUserDataSecretName, wign.File.Data)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secret for worker machines")
	}
	w.UserDataFile = &asset.File{
		Filename: filepath.Join(directory, workerUserDataFileName),
		Data:     data,
	}

	w.MachineConfigFiles, err = machineconfig.Manifests(machineConfigs, "worker", directory)
	if err != nil {
		return errors.Wrap(err, "failed to create MachineConfig manifests for worker machines")
	}

	w.MachineSetFiles = make([]*asset.File, len(machineSets))
	padFormat := fmt.Sprintf("%%0%dd", len(fmt.Sprintf("%d", len(machineSets))))
	for i, machineSet := range machineSets {
		data, err := yaml.Marshal(machineSet)
		if err != nil {
			return errors.Wrapf(err, "marshal worker %d", i)
		}

		padded := fmt.Sprintf(padFormat, i)
		w.MachineSetFiles[i] = &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(workerMachineSetFileName, padded)),
			Data:     data,
		}
	}
	return nil
}

// Files returns the files generated by the asset.
func (w *Worker) Files() []*asset.File {
	files := make([]*asset.File, 0, 1+len(w.MachineConfigFiles)+len(w.MachineSetFiles))
	if w.UserDataFile != nil {
		files = append(files, w.UserDataFile)
	}
	files = append(files, w.MachineConfigFiles...)
	files = append(files, w.MachineSetFiles...)
	return files
}

// Load reads the asset files from disk.
func (w *Worker) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(filepath.Join(directory, workerUserDataFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	w.UserDataFile = file

	w.MachineConfigFiles, err = machineconfig.Load(f, "worker", directory)
	if err != nil {
		return true, err
	}

	fileList, err := f.FetchByPattern(filepath.Join(directory, workerMachineSetFileNamePattern))
	if err != nil {
		return true, err
	}

	w.MachineSetFiles = fileList
	return true, nil
}

// MachineSets returns MachineSet manifest structures.
func (w *Worker) MachineSets() ([]machinev1beta1.MachineSet, error) {
	scheme := runtime.NewScheme()
	baremetalapi.AddToScheme(scheme)
	nutanixapi.AddToScheme(scheme)
	ibmcloudapi.AddToScheme(scheme)
	libvirtapi.AddToScheme(scheme)
	openstackapi.AddToScheme(scheme)
	ovirtproviderapi.AddToScheme(scheme)
	powervsapi.AddToScheme(scheme)
	scheme.AddKnownTypes(machinev1beta1.SchemeGroupVersion,
		&machinev1beta1.AWSMachineProviderConfig{},
		&machinev1beta1.VSphereMachineProviderSpec{},
		&machinev1beta1.AzureMachineProviderSpec{},
		&machinev1beta1.GCPMachineProviderSpec{},
	)
	scheme.AddKnownTypes(machinev1.GroupVersion,
		&machinev1.AlibabaCloudMachineProviderConfig{},
	)
	machinev1beta1.AddToScheme(scheme)
	decoder := serializer.NewCodecFactory(scheme).UniversalDecoder(
		baremetalprovider.SchemeGroupVersion,
		ibmcloudprovider.SchemeGroupVersion,
		libvirtprovider.SchemeGroupVersion,
		nutanixprovider.SchemeGroupVersion,
		openstackprovider.SchemeGroupVersion,
		ovirtprovider.SchemeGroupVersion,
		powervsprovider.GroupVersion,
		machinev1beta1.SchemeGroupVersion,
		machinev1.GroupVersion,
	)

	machineSets := []machinev1beta1.MachineSet{}
	for i, file := range w.MachineSetFiles {
		machineSet := &machinev1beta1.MachineSet{}
		err := yaml.Unmarshal(file.Data, &machineSet)
		if err != nil {
			return machineSets, errors.Wrapf(err, "unmarshal worker %d", i)
		}

		obj, _, err := decoder.Decode(machineSet.Spec.Template.Spec.ProviderSpec.Value.Raw, nil, nil)
		if err != nil {
			return machineSets, errors.Wrapf(err, "unmarshal worker %d", i)
		}

		machineSet.Spec.Template.Spec.ProviderSpec.Value = &runtime.RawExtension{Object: obj}
		machineSets = append(machineSets, *machineSet)
	}

	return machineSets, nil
}
