package machines

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

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

func defaultOpenStackMachinePoolPlatform(flavor string) openstacktypes.MachinePool {
	return openstacktypes.MachinePool{
		FlavorName: flavor,
	}
}

// Worker generates the machinesets for `worker` machine pool.
type Worker struct {
	MachineSetRaw     []byte
	UserDataSecretRaw []byte
}

var _ asset.Asset = (*Worker)(nil)

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

	var err error
	userDataMap := map[string][]byte{"worker-user-data": wign.File.Data}
	w.UserDataSecretRaw, err = userDataList(userDataMap)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secret for worker machines")
	}

	machineSets := []runtime.Object{}

	ic := installconfig.Config
	for _, pool := range ic.Compute {
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
			sets, err := aws.MachineSets(clusterID.ClusterID, ic, &pool, string(*rhcosImage), "worker", "worker-user-data")
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
			sets, err := libvirt.MachineSets(clusterID.ClusterID, ic, &pool, "worker", "worker-user-data")
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case nonetypes.Name:
		case openstacktypes.Name:
			mpool := defaultOpenStackMachinePoolPlatform(ic.Platform.OpenStack.FlavorName)
			mpool.Set(ic.Platform.OpenStack.DefaultMachinePlatform)
			mpool.Set(pool.Platform.OpenStack)
			pool.Platform.OpenStack = &mpool

			sets, err := openstack.MachineSets(clusterID.ClusterID, ic, &pool, string(*rhcosImage), "worker", "worker-user-data")
			if err != nil {
				return errors.Wrap(err, "failed to create master machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		default:
			return fmt.Errorf("invalid Platform")
		}
	}

	list := &metav1.List{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "List",
		},
		Items: make([]runtime.RawExtension, len(machineSets)),
	}
	for i, set := range machineSets {
		list.Items[i] = runtime.RawExtension{Object: set}
	}
	raw, err := yaml.Marshal(list)
	if err != nil {
		return errors.Wrap(err, "failed to marshal")
	}
	w.MachineSetRaw = raw

	return nil
}

func applyTemplateData(template *template.Template, templateData interface{}) []byte {
	buf := &bytes.Buffer{}
	if err := template.Execute(buf, templateData); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
