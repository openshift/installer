package machines

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines/aws"
	"github.com/openshift/installer/pkg/asset/machines/libvirt"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
)

func defaultAWSMachinePoolPlatform() types.AWSMachinePoolPlatform {
	return types.AWSMachinePoolPlatform{
		InstanceType: "t2.medium",
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
	w.UserDataSecretRaw, err = userData("worker-user-data", wign.File.Data)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secret for worker machines")
	}

	ic := installconfig.Config
	pool := workerPool(ic.Machines)
	numOfWorkers := int64(0)
	if pool.Replicas != nil {
		numOfWorkers = *pool.Replicas
	}

	switch ic.Platform.Name() {
	case "aws":
		config := aws.Config{
			ClusterName: ic.ObjectMeta.Name,
			Replicas:    numOfWorkers,
			Region:      ic.Platform.AWS.Region,
			Machine:     defaultAWSMachinePoolPlatform(),
		}

		tags := map[string]string{
			"tectonicClusterID": ic.ClusterID,
			// Info: https://github.com/openshift/cluster-api-provider-aws/issues/73
			// fmt.Sprintf("kubernetes.io/cluster/%s", ic.ObjectMeta.Name): "owned",
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

		w.MachineSetRaw = applyTemplateData(aws.WorkerMachineSetTmpl, config)
	case "libvirt":
		config := libvirt.Config{
			ClusterName: ic.ObjectMeta.Name,
			Replicas:    numOfWorkers,
			Platform:    *ic.Platform.Libvirt,
		}
		w.MachineSetRaw = applyTemplateData(libvirt.WorkerMachineSetTmpl, config)
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
