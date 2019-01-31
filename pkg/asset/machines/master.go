package machines

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	libvirtapi "github.com/openshift/cluster-api-provider-libvirt/pkg/apis"
	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1alpha1"
	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines/aws"
	"github.com/openshift/installer/pkg/asset/machines/libvirt"
	"github.com/openshift/installer/pkg/asset/machines/openstack"
	"github.com/openshift/installer/pkg/asset/rhcos"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
	libvirttypes "github.com/openshift/installer/pkg/types/libvirt"
	nonetypes "github.com/openshift/installer/pkg/types/none"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	awsapi "sigs.k8s.io/cluster-api-provider-aws/pkg/apis"
	awsprovider "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
	openstackapi "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis"
	openstackprovider "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"
)

// Master generates the machines for the `master` machine pool.
type Master struct {
	FileList []*asset.File
}

var (
	directory = "openshift"

	// MasterMachineFileName is the format string for constucting the master Machine filenames.
	MasterMachineFileName = "99_openshift-cluster-api_master-machines-%s.yaml"

	// MasterUserDataFileName is the filename used for the master user-data secret.
	MasterUserDataFileName = "99_openshift-cluster-api_master-user-data-secret.yaml"

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
func (m *Master) Generate(dependencies asset.Parents) error {
	clusterID := &installconfig.ClusterID{}
	installconfig := &installconfig.InstallConfig{}
	rhcosImage := new(rhcos.Image)
	mign := &machine.Master{}
	dependencies.Get(clusterID, installconfig, rhcosImage, mign)

	var err error
	machines := []machineapi.Machine{}
	ic := installconfig.Config
	pool := ic.ControlPlane
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
		machines, err = aws.Machines(clusterID.ClusterID, ic, pool, string(*rhcosImage), "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		aws.ConfigMasters(machines, ic.ObjectMeta.Name)
	case libvirttypes.Name:
		mpool := defaultLibvirtMachinePoolPlatform()
		mpool.Set(ic.Platform.Libvirt.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Libvirt)
		pool.Platform.Libvirt = &mpool
		machines, err = libvirt.Machines(clusterID.ClusterID, ic, pool, "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
	case nonetypes.Name:
		return nil
	case openstacktypes.Name:
		mpool := defaultOpenStackMachinePoolPlatform(ic.Platform.OpenStack.FlavorName)
		mpool.Set(ic.Platform.OpenStack.DefaultMachinePlatform)
		mpool.Set(pool.Platform.OpenStack)
		pool.Platform.OpenStack = &mpool

		machines, err = openstack.Machines(clusterID.ClusterID, ic, pool, string(*rhcosImage), "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		openstack.ConfigMasters(machines, ic.ObjectMeta.Name)
	default:
		return fmt.Errorf("invalid Platform")
	}

	userDataMap := map[string][]byte{"master-user-data": mign.File.Data}
	data, err := userDataList(userDataMap)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secret for master machines")
	}

	m.FileList = []*asset.File{{
		Filename: filepath.Join(directory, MasterUserDataFileName),
		Data:     data,
	}}

	count := len(machines)
	if count == 0 {
		return errors.New("at least one master machine must be configured")
	}

	padFormat := fmt.Sprintf("%%0%dd", len(fmt.Sprintf("%d", count)))
	for i, machine := range machines {
		data, err := yaml.Marshal(machine)
		if err != nil {
			return errors.Wrapf(err, "marshal master %d", i)
		}

		padded := fmt.Sprintf(padFormat, i)
		m.FileList = append(m.FileList, &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(MasterMachineFileName, padded)),
			Data:     data,
		})
	}

	return nil
}

// Files returns the files generated by the asset.
func (m *Master) Files() []*asset.File {
	return m.FileList
}

// Load reads the asset files from disk.
func (m *Master) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(filepath.Join(directory, MasterUserDataFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	m.FileList = []*asset.File{file}

	fileList, err := f.FetchByPattern(filepath.Join(directory, fmt.Sprintf(MasterMachineFileName, "*")))
	if err != nil {
		return true, err
	}

	if len(fileList) == 0 {
		return true, errors.Errorf("master machine manifests are required if you also provide %s", file.Filename)
	}

	m.FileList = append(m.FileList, fileList...)
	return true, nil
}

// Machines returns master Machine manifest YAML.
func (m *Master) Machines() [][]byte {
	machines := [][]byte{}
	userData := filepath.Join(directory, MasterUserDataFileName)
	for _, file := range m.FileList {
		if file.Filename == userData {
			continue
		}
		machines = append(machines, file.Data)
	}
	return machines
}

// StructuredMachines returns master Machine manifest structures.
func (m *Master) StructuredMachines() ([]machineapi.Machine, error) {
	scheme := runtime.NewScheme()
	awsapi.AddToScheme(scheme)
	libvirtapi.AddToScheme(scheme)
	openstackapi.AddToScheme(scheme)
	decoder := serializer.NewCodecFactory(scheme).UniversalDecoder(
		awsprovider.SchemeGroupVersion,
		libvirtprovider.SchemeGroupVersion,
		openstackprovider.SchemeGroupVersion,
	)

	machines := []machineapi.Machine{}
	for i, data := range m.Machines() {
		machine := &machineapi.Machine{}
		err := yaml.Unmarshal(data, &machine)
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
