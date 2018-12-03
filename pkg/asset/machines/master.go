package machines

import (
	"context"
	"fmt"
	"time"

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
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
)

// Master generates the machines for the `master` machine pool.
type Master struct {
	MachinesRaw       []byte
	UserDataSecretRaw []byte
}

var _ asset.Asset = (*Master)(nil)

// Name returns a human friendly name for the Master Asset.
func (m *Master) Name() string {
	return "Master Machines"
}

// Dependencies returns all of the dependencies directly needed by the
// Master asset
func (m *Master) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&machine.Master{},
	}
}

// Generate generates the Master asset.
func (m *Master) Generate(dependencies asset.Parents) error {
	installconfig := &installconfig.InstallConfig{}
	mign := &machine.Master{}
	dependencies.Get(installconfig, mign)

	var err error
	userDataMap := map[string][]byte{"master-user-data": mign.File.Data}
	m.UserDataSecretRaw, err = userDataList(userDataMap)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secret for master machines")
	}

	ic := installconfig.Config
	pool := masterPool(ic.Machines)
	switch ic.Platform.Name() {
	case "aws":
		mpool := defaultAWSMachinePoolPlatform()
		mpool.Set(ic.Platform.AWS.DefaultMachinePlatform)
		mpool.Set(pool.Platform.AWS)
		if mpool.AMIID == "" {
			ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
			defer cancel()
			ami, err := rhcos.AMI(ctx, rhcos.DefaultChannel, ic.Platform.AWS.Region)
			if err != nil {
				return errors.Wrap(err, "failed to determine default AMI")
			}
			mpool.AMIID = ami
		}
		if len(mpool.Zones) == 0 {
			azs, err := aws.AvailabilityZones(ic.Platform.AWS.Region)
			if err != nil {
				return errors.Wrap(err, "failed to fetch availability zones")
			}
			mpool.Zones = azs
		}
		pool.Platform.AWS = &mpool
		machines, err := aws.Machines(ic, &pool, "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		aws.ConfigMasters(machines, ic.ObjectMeta.Name)

		list := listFromMachines(machines)
		raw, err := yaml.Marshal(list)
		if err != nil {
			return errors.Wrap(err, "failed to marshal")
		}
		m.MachinesRaw = raw
	case "libvirt":
		machines, err := libvirt.Machines(ic, &pool, "master", "master-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}

		list := listFromMachines(machines)
		raw, err := yaml.Marshal(list)
		if err != nil {
			return errors.Wrap(err, "failed to marshal")
		}
		m.MachinesRaw = raw
	case "openstack":
		numOfMasters := int64(0)
		if pool.Replicas != nil {
			numOfMasters = *pool.Replicas
		}
		instances := []string{}
		for i := 0; i < int(numOfMasters); i++ {
			instances = append(instances, fmt.Sprintf("master-%d", i))
		}
		config := openstack.MasterConfig{
			ClusterName: ic.ObjectMeta.Name,
			Instances:   instances,
			Image:       ic.Platform.OpenStack.BaseImage,
			Region:      ic.Platform.OpenStack.Region,
			Machine:     defaultOpenStackMachinePoolPlatform(ic.Platform.OpenStack.FlavorName),
		}

		tags := map[string]string{
			"openshiftClusterID": ic.ClusterID,
		}
		config.Tags = tags

		config.Machine.Set(ic.Platform.OpenStack.DefaultMachinePlatform)
		config.Machine.Set(pool.Platform.OpenStack)

		m.MachinesRaw = applyTemplateData(openstack.MasterMachinesTmpl, config)
	default:
		return fmt.Errorf("invalid Platform")
	}
	return nil
}

func masterPool(pools []types.MachinePool) types.MachinePool {
	for idx, pool := range pools {
		if pool.Name == "master" {
			return pools[idx]
		}
	}
	return types.MachinePool{}
}

func listFromMachines(objs []clusterapi.Machine) *metav1.List {
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
