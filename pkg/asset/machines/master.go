package machines

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

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
	MachinesRaw        []byte
	UserDataSecretsRaw []byte
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

	userDataContent := map[string][]byte{}
	for i, file := range mign.FileList {
		userDataContent[fmt.Sprintf("master-user-data-%d", i)] = file.Data
	}

	var err error
	m.UserDataSecretsRaw, err = userDataList(userDataContent)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secrets for master machines")
	}

	ic := installconfig.Config
	pool := masterPool(ic.Machines)
	numOfMasters := int64(0)
	if pool.Replicas != nil {
		numOfMasters = *pool.Replicas
	}

	switch ic.Platform.Name() {
	case "aws":
		config := aws.MasterConfig{}
		config.ClusterName = ic.ObjectMeta.Name
		config.Region = ic.Platform.AWS.Region
		config.Machine = defaultAWSMachinePoolPlatform()

		tags := map[string]string{
			"tectonicClusterID": ic.ClusterID,
		}
		for k, v := range ic.Platform.AWS.UserTags {
			tags[k] = v
		}
		config.Tags = tags

		config.Machine.Set(ic.Platform.AWS.DefaultMachinePlatform)
		config.Machine.Set(pool.Platform.AWS)

		ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
		defer cancel()
		ami, err := rhcos.AMI(ctx, rhcos.DefaultChannel, config.Region)
		if err != nil {
			return errors.Wrap(err, "failed to determine default AMI")
		}
		config.AMIID = ami
		azs, err := aws.AvailabilityZones(config.Region)
		if err != nil {
			return errors.Wrap(err, "failed to fetch availability zones")
		}

		for i := 0; i < int(numOfMasters); i++ {
			azIndex := i % len(azs)
			config.Instances = append(config.Instances, aws.MasterInstance{AvailabilityZone: azs[azIndex]})
		}

		m.MachinesRaw = applyTemplateData(aws.MasterMachineTmpl, config)
	case "libvirt":
		instances := []string{}
		for i := 0; i < int(numOfMasters); i++ {
			instances = append(instances, fmt.Sprintf("master-%d", i))
		}
		config := libvirt.MasterConfig{
			ClusterName: ic.ObjectMeta.Name,
			Instances:   instances,
			Platform:    *ic.Platform.Libvirt,
		}
		m.MachinesRaw = applyTemplateData(libvirt.MasterMachinesTmpl, config)
	case "openstack":
		instances := []string{}
		for i := 0; i < int(numOfMasters); i++ {
			instances = append(instances, fmt.Sprintf("master-%d", i))
		}
		config := openstack.MasterConfig{
			ClusterName: ic.ObjectMeta.Name,
			Instances:   instances,
			Image:       ic.Platform.OpenStack.BaseImage,
			Region:      ic.Platform.OpenStack.Region,
			Machine:     defaultOpenStackMachinePoolPlatform(),
		}

		tags := map[string]string{
			"tectonicClusterID": ic.ClusterID,
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
