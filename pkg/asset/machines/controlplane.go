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
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
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

// ControlPlane generates the machines for the control-plane machine pool.
type ControlPlane struct {
	FileList []*asset.File
}

var (
	directory = "openshift"

	// ControlPlaneMachineFileName is the format string for constucting the control-plane Machine filenames.
	ControlPlaneMachineFileName = "99_openshift-cluster-api_control-plane-machines-%s.yaml"

	// ControlPlaneUserDataFileName is the filename used for the control-plane user-data secret.
	ControlPlaneUserDataFileName = "99_openshift-cluster-api_control-plane-user-data-secret.yaml"

	_ asset.WritableAsset = (*ControlPlane)(nil)
)

// Name returns a human friendly name for the Control Plane Asset.
func (a *ControlPlane) Name() string {
	return "Control Plane Machines"
}

// Dependencies returns all of the dependencies directly needed by the
// ControlPlane asset
func (a *ControlPlane) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		// PlatformCredsCheck just checks the creds (and asks, if needed)
		// We do not actually use it in this asset directly, hence
		// it is put in the dependencies but not fetched in Generate
		&installconfig.PlatformCredsCheck{},
		&installconfig.InstallConfig{},
		new(rhcos.Image),
		&machine.ControlPlane{},
	}
}

// Generate generates the ControlPlane asset.
func (a *ControlPlane) Generate(dependencies asset.Parents) error {
	clusterID := &installconfig.ClusterID{}
	installconfig := &installconfig.InstallConfig{}
	rhcosImage := new(rhcos.Image)
	cpign := &machine.ControlPlane{}
	dependencies.Get(clusterID, installconfig, rhcosImage, cpign)

	var err error
	machines := []machineapi.Machine{}
	ic := installconfig.Config
	pool := ic.ControlPlane
	switch ic.Platform.Name() {
	case awstypes.Name:
		mpool := defaultAWSMachinePoolPlatform()
		mpool.InstanceType = "m4.xlarge"
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
		machines, err = aws.Machines(clusterID.ClusterID, ic, pool, string(*rhcosImage), types.ControlPlaneMachineRole, "control-plane-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create control plane machine objects")
		}
		aws.ConfigControlPlane(machines, ic.ObjectMeta.Name)
	case libvirttypes.Name:
		mpool := defaultLibvirtMachinePoolPlatform()
		mpool.Set(ic.Platform.Libvirt.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Libvirt)
		pool.Platform.Libvirt = &mpool
		machines, err = libvirt.Machines(clusterID.ClusterID, ic, pool, types.ControlPlaneMachineRole, "control-plane-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create control plane machine objects")
		}
	case nonetypes.Name:
		return nil
	case openstacktypes.Name:
		mpool := defaultOpenStackMachinePoolPlatform(ic.Platform.OpenStack.FlavorName)
		mpool.Set(ic.Platform.OpenStack.DefaultMachinePlatform)
		mpool.Set(pool.Platform.OpenStack)
		pool.Platform.OpenStack = &mpool

		machines, err = openstack.Machines(clusterID.ClusterID, ic, pool, string(*rhcosImage), types.ControlPlaneMachineRole, "control-plane-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create control plane machine objects")
		}
		openstack.ConfigControlPlane(machines, ic.ObjectMeta.Name)
	default:
		return fmt.Errorf("invalid Platform")
	}

	userDataMap := map[string][]byte{"control-plane-user-data": cpign.File.Data}
	data, err := userDataList(userDataMap)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secret for control plane machines")
	}

	a.FileList = []*asset.File{{
		Filename: filepath.Join(directory, ControlPlaneUserDataFileName),
		Data:     data,
	}}

	count := len(machines)
	if count == 0 {
		return errors.New("at least one control plane machine must be configured")
	}

	padFormat := fmt.Sprintf("%%0%dd", len(fmt.Sprintf("%d", count)))
	for i, machine := range machines {
		data, err := yaml.Marshal(machine)
		if err != nil {
			return errors.Wrapf(err, "marshal control plane %d", i)
		}

		padded := fmt.Sprintf(padFormat, i)
		a.FileList = append(a.FileList, &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(ControlPlaneMachineFileName, padded)),
			Data:     data,
		})
	}

	return nil
}

// Files returns the files generated by the asset.
func (a *ControlPlane) Files() []*asset.File {
	return a.FileList
}

// Load reads the asset files from disk.
func (a *ControlPlane) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(filepath.Join(directory, ControlPlaneUserDataFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	a.FileList = []*asset.File{file}

	fileList, err := f.FetchByPattern(filepath.Join(directory, fmt.Sprintf(ControlPlaneMachineFileName, "*")))
	if err != nil {
		return true, err
	}

	if len(fileList) == 0 {
		return true, errors.Errorf("control plane machine manifests are required if you also provide %s", file.Filename)
	}

	a.FileList = append(a.FileList, fileList...)
	return true, nil
}

// Machines returns control plane Machine manifest YAML.
func (a *ControlPlane) Machines() [][]byte {
	machines := [][]byte{}
	userData := filepath.Join(directory, ControlPlaneUserDataFileName)
	for _, file := range a.FileList {
		if file.Filename == userData {
			continue
		}
		machines = append(machines, file.Data)
	}
	return machines
}

// StructuredMachines returns control plane Machine manifest structures.
func (a *ControlPlane) StructuredMachines() ([]machineapi.Machine, error) {
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
	for i, data := range a.Machines() {
		machine := &machineapi.Machine{}
		err := yaml.Unmarshal(data, &machine)
		if err != nil {
			return machines, errors.Wrapf(err, "unmarshal control plane %d", i)
		}

		obj, _, err := decoder.Decode(machine.Spec.ProviderSpec.Value.Raw, nil, nil)
		if err != nil {
			return machines, errors.Wrapf(err, "unmarshal control plane %d", i)
		}

		machine.Spec.ProviderSpec.Value = &runtime.RawExtension{Object: obj}
		machines = append(machines, *machine)
	}

	return machines, nil
}
