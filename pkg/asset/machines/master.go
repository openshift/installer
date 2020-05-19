package machines

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	baremetalapi "github.com/metal3-io/cluster-api-provider-baremetal/pkg/apis"
	baremetalprovider "github.com/metal3-io/cluster-api-provider-baremetal/pkg/apis/baremetal/v1alpha1"
	gcpapi "github.com/openshift/cluster-api-provider-gcp/pkg/apis"
	gcpprovider "github.com/openshift/cluster-api-provider-gcp/pkg/apis/gcpprovider/v1beta1"
	libvirtapi "github.com/openshift/cluster-api-provider-libvirt/pkg/apis"
	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	ovirtproviderapi "github.com/openshift/cluster-api-provider-ovirt/pkg/apis"
	ovirtprovider "github.com/openshift/cluster-api-provider-ovirt/pkg/apis/ovirtprovider/v1beta1"
	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	vsphereapi "github.com/openshift/machine-api-operator/pkg/apis/vsphereprovider"
	vsphereprovider "github.com/openshift/machine-api-operator/pkg/apis/vsphereprovider/v1beta1"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	awsapi "sigs.k8s.io/cluster-api-provider-aws/pkg/apis"
	awsprovider "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1beta1"
	azureapi "sigs.k8s.io/cluster-api-provider-azure/pkg/apis"
	azureprovider "sigs.k8s.io/cluster-api-provider-azure/pkg/apis/azureprovider/v1beta1"
	openstackapi "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis"
	openstackprovider "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines/aws"
	"github.com/openshift/installer/pkg/asset/machines/azure"
	"github.com/openshift/installer/pkg/asset/machines/baremetal"
	"github.com/openshift/installer/pkg/asset/machines/gcp"
	"github.com/openshift/installer/pkg/asset/machines/libvirt"
	"github.com/openshift/installer/pkg/asset/machines/machineconfig"
	"github.com/openshift/installer/pkg/asset/machines/openstack"
	"github.com/openshift/installer/pkg/asset/machines/ovirt"
	"github.com/openshift/installer/pkg/asset/machines/vsphere"
	"github.com/openshift/installer/pkg/asset/rhcos"
	rhcosutils "github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	azuredefaults "github.com/openshift/installer/pkg/types/azure/defaults"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	libvirttypes "github.com/openshift/installer/pkg/types/libvirt"
	nonetypes "github.com/openshift/installer/pkg/types/none"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// Master generates the machines for the `master` machine pool.
type Master struct {
	UserDataFile       *asset.File
	MachineConfigFiles []*asset.File
	MachineFiles       []*asset.File

	// SecretFiles is used by the baremetal platform to register the
	// credential information for communicating with management
	// controllers on hosts.
	SecretFiles []*asset.File

	// HostFiles is the list of baremetal hosts provided in the
	// installer configuration.
	HostFiles []*asset.File
}

const (
	directory = "openshift"

	// secretFileName is the format string for constructing the Secret
	// filenames for baremetal clusters.
	secretFileName = "99_openshift-cluster-api_host-bmc-secrets-%s.yaml"

	// hostFileName is the format string for constucting the Host
	// filenames for baremetal clusters.
	hostFileName = "99_openshift-cluster-api_hosts-%s.yaml"

	// masterMachineFileName is the format string for constucting the
	// master Machine filenames.
	masterMachineFileName = "99_openshift-cluster-api_master-machines-%s.yaml"

	// masterUserDataFileName is the filename used for the master
	// user-data secret.
	masterUserDataFileName = "99_openshift-cluster-api_master-user-data-secret.yaml"
)

var (
	secretFileNamePattern        = fmt.Sprintf(secretFileName, "*")
	hostFileNamePattern          = fmt.Sprintf(hostFileName, "*")
	masterMachineFileNamePattern = fmt.Sprintf(masterMachineFileName, "*")

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

func awsDefaultMasterMachineTypes(region string) []string {
	classes := awsdefaults.InstanceClasses(region)
	types := make([]string, len(classes))
	for i, c := range classes {
		types[i] = fmt.Sprintf("%s.xlarge", c)
	}
	return types
}

// Generate generates the Master asset.
func (m *Master) Generate(dependencies asset.Parents) error {
	ctx := context.TODO()
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	rhcosImage := new(rhcos.Image)
	mign := &machine.Master{}
	dependencies.Get(clusterID, installConfig, rhcosImage, mign)

	ic := installConfig.Config

	pool := ic.ControlPlane
	var err error
	machines := []machineapi.Machine{}
	switch ic.Platform.Name() {
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
			}
		}
		if mpool.InstanceType == "" {
			mpool.InstanceType, err = aws.PreferredInstanceType(ctx, installConfig.AWS, awsDefaultMasterMachineTypes(installConfig.Config.Platform.AWS.Region), mpool.Zones)
			if err != nil {
				logrus.Warn(errors.Wrap(err, "failed to find default instance type"))
				mpool.InstanceType = awsDefaultMasterMachineTypes(installConfig.Config.Platform.AWS.Region)[0]
			}
		}

		pool.Platform.AWS = &mpool
		machines, err = aws.Machines(
			clusterID.InfraID,
			installConfig.Config.Platform.AWS.Region,
			subnets,
			pool,
			"master",
			"master-user-data",
			installConfig.Config.Platform.AWS.UserTags,
		)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		aws.ConfigMasters(machines, clusterID.InfraID, ic.Publish)
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
		machines, err = gcp.Machines(clusterID.InfraID, ic, pool, string(*rhcosImage), "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		gcp.ConfigMasters(machines, clusterID.InfraID, ic.Publish)
	case libvirttypes.Name:
		mpool := defaultLibvirtMachinePoolPlatform()
		mpool.Set(ic.Platform.Libvirt.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Libvirt)
		pool.Platform.Libvirt = &mpool
		machines, err = libvirt.Machines(clusterID.InfraID, ic, pool, "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
	case openstacktypes.Name:
		mpool := defaultOpenStackMachinePoolPlatform(ic.Platform.OpenStack.FlavorName)
		mpool.Set(ic.Platform.OpenStack.DefaultMachinePlatform)
		mpool.Set(pool.Platform.OpenStack)
		pool.Platform.OpenStack = &mpool

		imageName, _ := rhcosutils.GenerateOpenStackImageName(string(*rhcosImage), clusterID.InfraID)

		machines, err = openstack.Machines(clusterID.InfraID, ic, pool, imageName, "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		openstack.ConfigMasters(machines, clusterID.InfraID)
	case azuretypes.Name:
		mpool := defaultAzureMachinePoolPlatform()
		mpool.InstanceType = azuredefaults.ControlPlaneInstanceType(installConfig.Config.Platform.Azure.Region)
		mpool.OSDisk.DiskSizeGB = 1024
		mpool.Set(ic.Platform.Azure.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Azure)
		if len(mpool.Zones) == 0 {
			azs, err := azure.AvailabilityZones(ic.Platform.Azure.Region, mpool.InstanceType)
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

		machines, err = azure.Machines(clusterID.InfraID, ic, pool, string(*rhcosImage), "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		azure.ConfigMasters(machines, clusterID.InfraID)
	case baremetaltypes.Name:
		mpool := defaultBareMetalMachinePoolPlatform()
		mpool.Set(ic.Platform.BareMetal.DefaultMachinePlatform)
		mpool.Set(pool.Platform.BareMetal)
		pool.Platform.BareMetal = &mpool

		machines, err = baremetal.Machines(clusterID.InfraID, ic, pool, string(*rhcosImage), "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}

		hostSettings, err := baremetal.Hosts(ic, machines)
		if err != nil {
			return errors.Wrap(err, "failed to assemble host data")
		}

		if len(hostSettings.Hosts) > 0 {
			m.HostFiles = make([]*asset.File, len(hostSettings.Hosts))
			padFormat := fmt.Sprintf("%%0%dd", len(fmt.Sprintf("%d", len(hostSettings.Hosts))))
			for i, host := range hostSettings.Hosts {
				data, err := yaml.Marshal(host)
				if err != nil {
					return errors.Wrapf(err, "marshal host %d", i)
				}

				padded := fmt.Sprintf(padFormat, i)
				m.HostFiles[i] = &asset.File{
					Filename: filepath.Join(directory, fmt.Sprintf(hostFileName, padded)),
					Data:     data,
				}
			}
		}

		if len(hostSettings.Secrets) > 0 {
			m.SecretFiles = make([]*asset.File, len(hostSettings.Secrets))
			padFormat := fmt.Sprintf("%%0%dd", len(fmt.Sprintf("%d", len(hostSettings.Secrets))))
			for i, secret := range hostSettings.Secrets {
				data, err := yaml.Marshal(secret)
				if err != nil {
					return errors.Wrapf(err, "marshal secret %d", i)
				}

				padded := fmt.Sprintf(padFormat, i)
				m.SecretFiles[i] = &asset.File{
					Filename: filepath.Join(directory, fmt.Sprintf(secretFileName, padded)),
					Data:     data,
				}
			}
		}
	case ovirttypes.Name:
		mpool := defaultOvirtMachinePoolPlatform()
		mpool.VMType = ovirttypes.VMTypeHighPerformance
		mpool.Set(ic.Platform.Ovirt.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Ovirt)
		pool.Platform.Ovirt = &mpool

		imageName, _ := rhcosutils.GenerateOpenStackImageName(string(*rhcosImage), clusterID.InfraID)

		machines, err = ovirt.Machines(clusterID.InfraID, ic, pool, imageName, "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects for ovirt provider")
		}
	case vspheretypes.Name:
		mpool := defaultVSphereMachinePoolPlatform()
		mpool.NumCPUs = 4
		mpool.MemoryMiB = 16384
		mpool.Set(ic.Platform.VSphere.DefaultMachinePlatform)
		mpool.Set(pool.Platform.VSphere)
		pool.Platform.VSphere = &mpool
		templateName := clusterID.InfraID + "-rhcos"

		machines, err = vsphere.Machines(clusterID.InfraID, ic, pool, templateName, "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		vsphere.ConfigMasters(machines, clusterID.InfraID)
	case nonetypes.Name:
	default:
		return fmt.Errorf("invalid Platform")
	}

	data, err := userDataSecret("master-user-data", mign.File.Data)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secret for master machines")
	}

	m.UserDataFile = &asset.File{
		Filename: filepath.Join(directory, masterUserDataFileName),
		Data:     data,
	}

	machineConfigs := []*mcfgv1.MachineConfig{}
	if pool.Hyperthreading == types.HyperthreadingDisabled {
		machineConfigs = append(machineConfigs, machineconfig.ForHyperthreadingDisabled("master"))
	}
	if ic.SSHKey != "" {
		machineConfigs = append(machineConfigs, machineconfig.ForAuthorizedKeys(ic.SSHKey, "master"))
	}
	if ic.FIPS {
		machineConfigs = append(machineConfigs, machineconfig.ForFIPSEnabled("master"))
	}

	m.MachineConfigFiles, err = machineconfig.Manifests(machineConfigs, "master", directory)
	if err != nil {
		return errors.Wrap(err, "failed to create MachineConfig manifests for master machines")
	}

	m.MachineFiles = make([]*asset.File, len(machines))
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
	// Machines are linked to hosts via the machineRef, so we create
	// the hosts first to ensure if the operator starts trying to
	// reconcile a machine it can pick up the related host.
	files = append(files, m.HostFiles...)
	files = append(files, m.MachineFiles...)
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

	return true, nil
}

// Machines returns master Machine manifest structures.
func (m *Master) Machines() ([]machineapi.Machine, error) {
	scheme := runtime.NewScheme()
	awsapi.AddToScheme(scheme)
	azureapi.AddToScheme(scheme)
	baremetalapi.AddToScheme(scheme)
	gcpapi.AddToScheme(scheme)
	libvirtapi.AddToScheme(scheme)
	openstackapi.AddToScheme(scheme)
	ovirtproviderapi.AddToScheme(scheme)
	vsphereapi.AddToScheme(scheme)
	decoder := serializer.NewCodecFactory(scheme).UniversalDecoder(
		awsprovider.SchemeGroupVersion,
		azureprovider.SchemeGroupVersion,
		baremetalprovider.SchemeGroupVersion,
		gcpprovider.SchemeGroupVersion,
		libvirtprovider.SchemeGroupVersion,
		openstackprovider.SchemeGroupVersion,
		vsphereprovider.SchemeGroupVersion,
		ovirtprovider.SchemeGroupVersion,
	)

	machines := []machineapi.Machine{}
	for i, file := range m.MachineFiles {
		machine := &machineapi.Machine{}
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
	if filename == masterUserDataFileName || filename == workerUserDataFileName {
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
