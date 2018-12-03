package machines

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
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
	awstypes "github.com/openshift/installer/pkg/types/aws"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
)

func defaultAWSMachinePoolPlatform() awstypes.MachinePool {
	return awstypes.MachinePool{
		InstanceType: "m4.large",
	}
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
		&installconfig.InstallConfig{},
		&machine.Worker{},
	}
}

// Generate generates the Worker asset.
func (w *Worker) Generate(dependencies asset.Parents) error {
	installconfig := &installconfig.InstallConfig{}
	wign := &machine.Worker{}
	dependencies.Get(installconfig, wign)

	var err error
	userDataMap := map[string][]byte{"worker-user-data": wign.File.Data}
	w.UserDataSecretRaw, err = userDataList(userDataMap)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secret for worker machines")
	}

	ic := installconfig.Config
	pool := workerPool(ic.Machines)
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
		sets, err := aws.MachineSets(ic, &pool, "worker", "worker-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create worker machine objects")
		}
		aws.ConfigWorkers(sets)

		list := listFromMachineSets(sets)
		raw, err := yaml.Marshal(list)
		if err != nil {
			return errors.Wrap(err, "failed to marshal")
		}
		w.MachineSetRaw = raw
	case "libvirt":
		sets, err := libvirt.MachineSets(ic, &pool, "worker", "worker-user-data")
		if err != nil {
			return errors.Wrap(err, "failed to create worker machine objects")
		}

		list := listFromMachineSets(sets)
		raw, err := yaml.Marshal(list)
		if err != nil {
			return errors.Wrap(err, "failed to marshal")
		}
		w.MachineSetRaw = raw
	case "openstack":
		numOfWorkers := int64(0)
		if pool.Replicas != nil {
			numOfWorkers = *pool.Replicas
		}
		config := openstack.Config{
			ClusterName: ic.ObjectMeta.Name,
			Replicas:    numOfWorkers,
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
