package machines

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	libvirtapi "github.com/openshift/cluster-api-provider-libvirt/pkg/apis"
	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	awsapi "sigs.k8s.io/cluster-api-provider-aws/pkg/apis"
	awsprovider "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
	azureapi "sigs.k8s.io/cluster-api-provider-azure/pkg/apis"
	azureprovider "sigs.k8s.io/cluster-api-provider-azure/pkg/apis/azureprovider/v1alpha1"
	openstackapi "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis"
	openstackprovider "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines/aws"
	"github.com/openshift/installer/pkg/asset/machines/azure"
	"github.com/openshift/installer/pkg/asset/machines/libvirt"
	"github.com/openshift/installer/pkg/asset/machines/machineconfig"
	"github.com/openshift/installer/pkg/asset/machines/openstack"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	azuredefaults "github.com/openshift/installer/pkg/types/azure/defaults"
	libvirttypes "github.com/openshift/installer/pkg/types/libvirt"
	nonetypes "github.com/openshift/installer/pkg/types/none"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// Master generates the machines for the `master` machine pool.
type Master struct {
	UserDataFile       *asset.File
	MachineConfigFiles []*asset.File
	MachineFiles       []*asset.File
}

const (
	directory = "openshift"

	// masterMachineFileName is the format string for constucting the master Machine filenames.
	masterMachineFileName = "99_openshift-cluster-api_master-machines-%s.yaml"

	// masterUserDataFileName is the filename used for the master user-data secret.
	masterUserDataFileName = "99_openshift-cluster-api_master-user-data-secret.yaml"
)

var (
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

func awsDefaultMasterMachineType(installconfig *installconfig.InstallConfig) string {
	region := installconfig.Config.Platform.AWS.Region
	instanceClass := awsdefaults.InstanceClass(region)
	return fmt.Sprintf("%s.xlarge", instanceClass)
}

// Generate generates the Master asset.
func (m *Master) Generate(log *logrus.Entry, parents asset.Parents) error {
	clusterID := &installconfig.ClusterID{}
	installconfig := &installconfig.InstallConfig{}
	rhcosImage := new(rhcos.Image)
	mign := &machine.Master{}
	parents.Get(clusterID, installconfig, rhcosImage, mign)

	ic := installconfig.Config

	pool := ic.ControlPlane
	var err error
	machines := []machineapi.Machine{}
	switch ic.Platform.Name() {
	case awstypes.Name:
		mpool := defaultAWSMachinePoolPlatform()
		mpool.InstanceType = awsDefaultMasterMachineType(installconfig)
		mpool.Set(ic.Platform.AWS.DefaultMachinePlatform)
		mpool.Set(pool.Platform.AWS)
		if len(mpool.Zones) == 0 {
			azs, err := aws.AvailabilityZones(ic.Platform.AWS.Region)
			if err != nil {
				return errors.Wrap(err, "failed to fetch availability zones")
			}
			mpool.Zones = azs
		}
		pool.Platform.AWS = &mpool
		machines, err = aws.Machines(clusterID.InfraID, ic, pool, string(*rhcosImage), "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		aws.ConfigMasters(machines, clusterID.InfraID)
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

		machines, err = openstack.Machines(clusterID.InfraID, ic, pool, string(*rhcosImage), "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		openstack.ConfigMasters(machines, clusterID.InfraID)
	case azuretypes.Name:
		mpool := defaultAzureMachinePoolPlatform()
		mpool.InstanceType = azuredefaults.InstanceClass(installconfig.Config.Platform.Azure.Region)
		mpool.Set(ic.Platform.Azure.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Azure)
		pool.Platform.Azure = &mpool

		machines, err = azure.Machines(clusterID.InfraID, ic, pool, string(*rhcosImage), "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		azure.ConfigMasters(machines, clusterID.InfraID)
	case nonetypes.Name, vspheretypes.Name:
	default:
		return fmt.Errorf("invalid Platform")
	}

	userDataMap := map[string][]byte{"master-user-data": mign.File.Data}
	data, err := userDataList(userDataMap)
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

	fileList, err := f.FetchByPattern(filepath.Join(directory, masterMachineFileNamePattern))
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
	libvirtapi.AddToScheme(scheme)
	openstackapi.AddToScheme(scheme)
	decoder := serializer.NewCodecFactory(scheme).UniversalDecoder(
		awsprovider.SchemeGroupVersion,
		azureprovider.SchemeGroupVersion,
		libvirtprovider.SchemeGroupVersion,
		openstackprovider.SchemeGroupVersion,
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
