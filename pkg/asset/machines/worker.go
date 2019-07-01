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
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	libvirttypes "github.com/openshift/installer/pkg/types/libvirt"
	nonetypes "github.com/openshift/installer/pkg/types/none"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

const (
	// workerMachineSetFileName is the format string for constructing the worker MachineSet filenames.
	workerMachineSetFileName = "99_openshift-cluster-api_worker-machineset-%s.yaml"

	// workerUserDataFileName is the filename used for the worker user-data secret.
	workerUserDataFileName = "99_openshift-cluster-api_worker-user-data-secret.yaml"
)

var (
	workerMachineSetFileNamePattern = fmt.Sprintf(workerMachineSetFileName, "*")

	_ asset.WritableAsset = (*Worker)(nil)
)

func defaultAWSMachinePoolPlatform() awstypes.MachinePool {
	return awstypes.MachinePool{
		EC2RootVolume: awstypes.EC2RootVolume{
			Type: "gp2",
			Size: 120,
		},
	}
}

func defaultLibvirtMachinePoolPlatform() libvirttypes.MachinePool {
	return libvirttypes.MachinePool{}
}

func defaultAzureMachinePoolPlatform() azuretypes.MachinePool {
	return azuretypes.MachinePool{}
}

func defaultOpenStackMachinePoolPlatform(flavor string) openstacktypes.MachinePool {
	return openstacktypes.MachinePool{
		FlavorName: flavor,
	}
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

func awsDefaultWorkerMachineType(installconfig *installconfig.InstallConfig) string {
	region := installconfig.Config.Platform.AWS.Region
	instanceClass := awsdefaults.InstanceClass(region)
	return fmt.Sprintf("%s.large", instanceClass)
}

// Generate generates the Worker asset.
func (w *Worker) Generate(dependencies asset.Parents) error {
	clusterID := &installconfig.ClusterID{}
	installconfig := &installconfig.InstallConfig{}
	rhcosImage := new(rhcos.Image)
	wign := &machine.Worker{}
	dependencies.Get(clusterID, installconfig, rhcosImage, wign)

	machineConfigs := []*mcfgv1.MachineConfig{}
	machineSets := []runtime.Object{}
	var err error
	ic := installconfig.Config
	for _, pool := range ic.Compute {
		if pool.Hyperthreading == types.HyperthreadingDisabled {
			machineConfigs = append(machineConfigs, machineconfig.ForHyperthreadingDisabled("worker"))
		}
		if ic.SSHKey != "" {
			machineConfigs = append(machineConfigs, machineconfig.ForAuthorizedKeys(ic.SSHKey, "worker"))
		}

		switch ic.Platform.Name() {
		case awstypes.Name:
			mpool := defaultAWSMachinePoolPlatform()
			mpool.InstanceType = awsDefaultWorkerMachineType(installconfig)
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
			sets, err := aws.MachineSets(clusterID.InfraID, ic, &pool, string(*rhcosImage), "worker", "worker-user-data")
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case gcptypes.Name:
		case libvirttypes.Name:
			mpool := defaultLibvirtMachinePoolPlatform()
			mpool.Set(ic.Platform.Libvirt.DefaultMachinePlatform)
			mpool.Set(pool.Platform.Libvirt)
			pool.Platform.Libvirt = &mpool
			sets, err := libvirt.MachineSets(clusterID.InfraID, ic, &pool, "worker", "worker-user-data")
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case openstacktypes.Name:
			mpool := defaultOpenStackMachinePoolPlatform(ic.Platform.OpenStack.FlavorName)
			mpool.Set(ic.Platform.OpenStack.DefaultMachinePlatform)
			mpool.Set(pool.Platform.OpenStack)
			pool.Platform.OpenStack = &mpool

			sets, err := openstack.MachineSets(clusterID.InfraID, ic, &pool, string(*rhcosImage), "worker", "worker-user-data")
			if err != nil {
				return errors.Wrap(err, "failed to create master machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case azuretypes.Name:
			mpool := defaultAzureMachinePoolPlatform()
			mpool.InstanceType = azuredefaults.InstanceClass(installconfig.Config.Platform.Azure.Region)
			mpool.Set(ic.Platform.Azure.DefaultMachinePlatform)
			mpool.Set(pool.Platform.Azure)
			pool.Platform.Azure = &mpool
			//TODO: add support for availibility zones
			sets, err := azure.MachineSets(clusterID.InfraID, ic, &pool, string(*rhcosImage), "worker", "worker-user-data")
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case nonetypes.Name, vspheretypes.Name:
		default:
			return fmt.Errorf("invalid Platform")
		}
	}

	userDataMap := map[string][]byte{"worker-user-data": wign.File.Data}
	data, err := userDataList(userDataMap)
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
func (w *Worker) MachineSets() ([]machineapi.MachineSet, error) {
	scheme := runtime.NewScheme()
	awsapi.AddToScheme(scheme)
	libvirtapi.AddToScheme(scheme)
	openstackapi.AddToScheme(scheme)
	azureapi.AddToScheme(scheme)
	decoder := serializer.NewCodecFactory(scheme).UniversalDecoder(
		awsprovider.SchemeGroupVersion,
		libvirtprovider.SchemeGroupVersion,
		openstackprovider.SchemeGroupVersion,
		azureprovider.SchemeGroupVersion,
	)

	machineSets := []machineapi.MachineSet{}
	for i, file := range w.MachineSetFiles {
		machineSet := &machineapi.MachineSet{}
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
