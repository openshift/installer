package machines

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

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
)

func defaultAWSMachinePoolPlatform() awstypes.MachinePool {
	return awstypes.MachinePool{}
}

func defaultLibvirtMachinePoolPlatform() libvirttypes.MachinePool {
	return libvirttypes.MachinePool{}
}

func defaultOpenStackMachinePoolPlatform(flavor string) openstacktypes.MachinePool {
	return openstacktypes.MachinePool{
		FlavorName: flavor,
	}
}

func trunkSupportBoolean(trunkSupport string) (result bool) {
	if trunkSupport == "1" {
		result = true
	} else {
		result = false
	}
	return
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
		&installconfig.InstallConfig{},
		new(rhcos.Image),
		&machine.Worker{},
	}
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

	ic := installconfig.Config
	pool := workerPool(ic.Machines)
	switch ic.Platform.Name() {
	case awstypes.Name:
		mpool := defaultAWSMachinePoolPlatform()
		mpool.InstanceType = "m4.large"
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

		list := listFromMachineSets(sets)
		raw, err := yaml.Marshal(list)
		if err != nil {
			return errors.Wrap(err, "failed to marshal")
		}
		w.MachineSetRaw = raw
	case libvirttypes.Name:
		mpool := defaultLibvirtMachinePoolPlatform()
		mpool.Set(ic.Platform.Libvirt.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Libvirt)
		pool.Platform.Libvirt = &mpool
		sets, err := libvirt.MachineSets(clusterID.ClusterID, ic, &pool, "worker", "worker-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create worker machine objects")
		}

		list := listFromMachineSets(sets)
		raw, err := yaml.Marshal(list)
		if err != nil {
			return errors.Wrap(err, "failed to marshal")
		}
		w.MachineSetRaw = raw
	case nonetypes.Name:
	case openstacktypes.Name:
		numOfWorkers := int64(0)
		if pool.Replicas != nil {
			numOfWorkers = *pool.Replicas
		}
		config := openstack.Config{
			CloudName:   ic.Platform.OpenStack.Cloud,
			ClusterName: ic.ObjectMeta.Name,
			Replicas:    numOfWorkers,
			Image:       string(*rhcosImage),
			Region:      ic.Platform.OpenStack.Region,
			Machine:     defaultOpenStackMachinePoolPlatform(ic.Platform.OpenStack.FlavorName),
			Trunk:       trunkSupportBoolean(ic.Platform.OpenStack.TrunkSupport),
		}

		tags := map[string]string{
			"openshiftClusterID": clusterID.ClusterID,
		}
		config.Tags = tags

		config.Machine.Set(ic.Platform.OpenStack.DefaultMachinePlatform)
		config.Machine.Set(pool.Platform.OpenStack)

		w.MachineSetRaw = applyTemplateData(openstack.WorkerMachineSetTmpl, config)
	default:
		return fmt.Errorf("invalid Platform")
	}
	return nil
}

func workerPool(pools []types.MachinePool) types.MachinePool {
	for idx, pool := range pools {
		if pool.Name == "worker" {
			return pools[idx]
		}
	}
	return types.MachinePool{}
}

func applyTemplateData(template *template.Template, templateData interface{}) []byte {
	buf := &bytes.Buffer{}
	if err := template.Execute(buf, templateData); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func listFromMachineSets(objs []clusterapi.MachineSet) *metav1.List {
	list := &metav1.List{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "List",
		},
	}
	for idx := range objs {
		list.Items = append(list.Items, runtime.RawExtension{Object: &objs[idx]})
	}
	return list
}
